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
	log "github.com/Sirupsen/logrus"
)

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, cluster configuration.ApplicationCluster) *state.Capture {

		//TODO: Remove this state capture stuff, it's clunkys
	stateCapture := &state.Capture{
		ClusterName:     cluster.Name,
		DeploymentState: make(map[string]state.Details),
	}
	//---------------------------------
	log.Warn(fmt.Sprintf("Switching to cluster: %s\n", cluster.Name))
	restclient, k8siface, err := platform.GetKubeClient(cluster.Name)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	//---------------------------------
	for _, deployment := range cluster.Applications {
		log.Debug(fmt.Sprintf("Loading deployment %s\n", deployment.Application.Name))
		//---------------------------------
		//Generate name from repo
		var extension = filepath.Ext(deployment.Application.Git)
		var remoteVCSRepoName = deployment.Application.Git[0 : len(deployment.Application.Git)-len(extension)]
		splitStrings := strings.Split(remoteVCSRepoName, "/")
		remoteVCSRepoName = splitStrings[len(splitStrings)-1]

		if _, err := os.Stat(path.Join(opt.TempVCSPath, remoteVCSRepoName)); os.IsNotExist(err) {
			log.Debug(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, path.Join(opt.TempVCSPath, remoteVCSRepoName)))
			gvcs := new(vcs.GitVCS)
			_, err = vcs.Fetch(gvcs, path.Join(opt.TempVCSPath, remoteVCSRepoName), deployment.Application.Git, opt.SSHKeyPath)
			if err != nil {
				log.Error(err.Error())
				stateCapture.DeploymentState[deployment.Application.Name] = state.Details{State: state.EDeploymentStateError}
				return stateCapture
			}
		} else {
			log.Debug(fmt.Sprintf("Using existing repository %s", path.Join(opt.TempVCSPath, remoteVCSRepoName)))
		}
		//---------------------------------
		for _, a := range deployment.Application.Action {
			if a.Execute.Shell != "" {
				log.Warn(fmt.Sprintf("Running shell command %s\n", a.Execute.Shell))
				if err := ShellCommand(a.Execute.Shell, path.Join(opt.TempVCSPath, remoteVCSRepoName), true); err != nil {
					log.Error(err.Error())
				}
			}
			//---------------------------------
			fileList := []string{}
			err := filepath.Walk(path.Join(opt.TempVCSPath, remoteVCSRepoName, a.Execute.Kubectl.Path), func(path string, f os.FileInfo, err error) error {
				fileList = append(fileList, path)
				return nil
			})
			if err != nil {
				log.Error(err.Error())

			}
			err = platform.GenerateDeploymentPlan(restclient,
				k8siface, fileList, deployment.Application.Namespace, opt, commandFlag)
			if err != nil {
				log.Error(err.Error())
			}
			//---------------------------------
		}
	}
	return stateCapture
}
