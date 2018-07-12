package scheduler

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/fatih/color"
)

func process(opt configuration.Options, cluster configuration.Cluster) *state.Capture {

	stateCapture := state.NewCapture()
	stateCapture.ClusterName = cluster.Name
	//---------------------------------
	color.Cyan(fmt.Sprintf("Switching to cluster: %s\n", cluster.Name))
	restclient, k8siface, err := platform.GetKubeClient(cluster.Name)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	//---------------------------------
	for _, deployment := range cluster.Deployments {
		//---------------------------------
		//Generate name from repo
		var extension = filepath.Ext(deployment.Deployment.Git)
		var remoteVCSRepoName = deployment.Deployment.Git[0 : len(deployment.Deployment.Git)-len(extension)]
		splitStrings := strings.Split(remoteVCSRepoName, "/")
		remoteVCSRepoName = splitStrings[len(splitStrings)-1]

		if _, err := os.Stat(path.Join(opt.TempVCSPath, remoteVCSRepoName)); os.IsNotExist(err) {
			color.Yellow(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, path.Join(opt.TempVCSPath, remoteVCSRepoName)))
			gvcs := new(vcs.GitVCS)
			_, err = vcs.Fetch(gvcs, path.Join(opt.TempVCSPath, remoteVCSRepoName), deployment.Deployment.Git, opt.SSHKeyPath)
			if err != nil {
				color.Red(err.Error())
				stateCapture.DeploymentState[deployment.Deployment.Name] = state.Details{State: state.EDeploymentStateError}
				return stateCapture
			}
		} else {
			color.Yellow(fmt.Sprintf("Using existing repository %s", path.Join(opt.TempVCSPath, remoteVCSRepoName)))
		}
		//---------------------------------
		for _, a := range deployment.Deployment.Action {
			if a.Execute.Shell != "" {
				color.Yellow(fmt.Sprintf("Running shell command %s\n", a.Execute.Shell))
				if err := ShellCommand(a.Execute.Shell, path.Join(opt.TempVCSPath, remoteVCSRepoName), true); err != nil {
					color.Red(err.Error())
				}
			}
			//---------------------------------
			if a.Execute.Kubectl.Command == "" {
				color.Red("No Kubernetes create action to run")
			}
			//---------------------------------
			fileList := []string{}
			err := filepath.Walk(path.Join(opt.TempVCSPath, remoteVCSRepoName, a.Execute.Kubectl.Path), func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
			if err != nil {
				color.Red(err.Error())

			}
			for _, file := range fileList {
				color.Yellow(fmt.Sprintf("Attempting to deploy %s\n", file))
				if _, err = os.Stat(file); os.IsNotExist(err) {
					continue
				}
				if sa, _ := os.Stat(file); sa.IsDir() {
					continue
				}
				var stateResponse state.State
				color.Yellow(fmt.Sprintf("Running..."))
				if stateResponse, err = platform.DeployFromFile(restclient, k8siface, file, deployment.Deployment.Namespace, opt); err != nil {
					color.Red(err.Error())
				}
				var output = ""
				var hasError = false
				if err != nil {
					output = fmt.Sprintf("File: %s Namespace :%s Error: %s", file, deployment.Deployment.Namespace, err)
					hasError = true
				} else {
					output = fmt.Sprintf("File: %s Namespace :%s", file, deployment.Deployment.Namespace)
				}
				stateCapture.DeploymentState[deployment.Deployment.Name] = state.Details{State: stateResponse, HasDetail: true,
					Detail: output, HasError: hasError}
			}
			//---------------------------------
		}
	}
	return stateCapture
}
