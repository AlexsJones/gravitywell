package scheduler

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/fatih/color"
)

//Options ...
type Options struct {
	VCS         string
	TempVCSPath string
	APIVersion  string
	Parallel    bool
}

//Scheduler object ...
type Scheduler struct {
	configuration *configuration.Configuration
}

//NewScheduler object ...
func NewScheduler(conf *configuration.Configuration) (*Scheduler, error) {
	if conf == nil {
		return nil, errors.New("Invalid configuration")
	}
	return &Scheduler{
		configuration: conf}, nil
}

//Run a new scheduler based off of the current configuration
func (s *Scheduler) Run(opt Options) error {

	if opt.APIVersion != s.configuration.APIVersion {
		color.Red(fmt.Sprintf("Manifest is not supported by the current API: %s\n", opt.APIVersion))
		os.Exit(1)
	}
	//---------------------------------
	if _, err := os.Stat(opt.TempVCSPath); os.IsNotExist(err) {
		os.Mkdir(opt.TempVCSPath, 0777)
	} else {
		os.RemoveAll(opt.TempVCSPath)
		os.Mkdir(opt.TempVCSPath, 0777)
	}
	//---------------------------------
	for _, cluster := range s.configuration.Strategy {
		//---------------------------------
		color.Yellow(fmt.Sprintf("Switching to cluster: %s\n", cluster.Cluster.Name))
		restclient, k8siface, err := platform.GetKubeClient(cluster.Cluster.Name)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		if opt.Parallel {
			color.Blue("Parallel cluster runs not implemented yet.")
		}
		for _, deployment := range cluster.Cluster.Deployments {
			//---------------------------------
			color.Yellow(fmt.Sprintf("Fetching deployment %s into %s\n", deployment.Deployment.Name, path.Join(opt.TempVCSPath, deployment.Deployment.Name)))
			if opt.VCS != "git" {
				color.Red("Only supporting git as VCS currently. Sorry.")
				os.Exit(1)
			}
			//---------------------------------
			gvcs := new(vcs.GitVCS)
			_, err = vcs.Fetch(gvcs, path.Join(opt.TempVCSPath, deployment.Deployment.Name), deployment.Deployment.Git)
			if err != nil {
				color.Red(err.Error())
				os.Exit(1)
			}
			//---------------------------------
			for _, a := range deployment.Deployment.Action {
				if a.Execute.Shell != "" {
					ShellCommand(a.Execute.Shell, path.Join(opt.TempVCSPath, cluster.Cluster.Name), false)
				}
				//---------------------------------
				if a.Execute.Kubectl.Command == "" {
					color.Red("No Kubernetes create action to run")
					continue
				}
				//---------------------------------
				fileList := []string{}
				err := filepath.Walk(path.Join(opt.TempVCSPath, deployment.Deployment.Name, a.Execute.Kubectl.Path), func(path string, f os.FileInfo, err error) error {

					fileList = append(fileList, path)
					return nil
				})

				if err != nil {
					color.Red(err.Error())
					os.Exit(1)
				}
				for _, file := range fileList {
					color.Yellow(fmt.Sprintf("Attempting to deploy %s\n", file))
					if _, err := os.Stat(file); os.IsNotExist(err) {
						continue

					}
					if sa, _ := os.Stat(file); sa.IsDir() {
						continue
					}

					if err := platform.DeployFromFile(restclient, k8siface, file, deployment.Deployment.Namespace); err != nil {
						color.Red(err.Error())
						return err
					}
				}
				//---------------------------------

			}
		}

	}

	return nil
}
