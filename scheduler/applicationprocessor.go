package scheduler

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/AlexsJones/gravitywell/subprocessor"
	"github.com/AlexsJones/gravitywell/vcs"
	log "github.com/Sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, cluster configuration.ApplicationCluster) *state.Capture {

	//TODO: Remove this state capture stuff, it's clunkys
	stateCapture := &state.Capture{
		ClusterName:     cluster.Name,
		DeploymentState: make(map[string]state.Details),
	}
	//---------------------------------
	log.Info(fmt.Sprintf("Switching to cluster: %s\n", cluster.Name))
	restclient, k8siface, err := platform.GetKubeClient(cluster.Name)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	//---------------------------------

	//Batch actions per namespace
	groupedDeployments := make(map[string][]configuration.Application)

	for _, deployment := range cluster.Applications {

		groupedDeployments[deployment.Application.Namespace] = append(groupedDeployments[deployment.Application.Namespace], deployment.Application)
	}

	for key, value := range groupedDeployments {

		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", key, cluster.Name)
		for _, depl := range value {
			fmt.Printf("\t\t %s\n", depl.Name)
		}
	}
	//Execute deployments
	coordinator := subprocessor.NewCoordinator()
	var wg sync.WaitGroup
	go coordinator.Run()

	for key, value := range groupedDeployments {

		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", key, cluster.Name)
		for _, deployment := range value {
			wg.Add(1)
			coordinator.ResourceChannel <- subprocessor.Resource{
				Process: func() {
					log.Debug(fmt.Sprintf("Loading deployment %s\n", deployment.Name))
					//---------------------------------
					//Generate name from repo
					var extension = filepath.Ext(deployment.Git)
					var remoteVCSRepoName = deployment.Git[0 : len(deployment.Git)-len(extension)]
					splitStrings := strings.Split(remoteVCSRepoName, "/")
					remoteVCSRepoName = splitStrings[len(splitStrings)-1]

					if _, err := os.Stat(path.Join(opt.TempVCSPath, remoteVCSRepoName)); os.IsNotExist(err) {
						log.Debug(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, path.Join(opt.TempVCSPath, remoteVCSRepoName)))
						gvcs := new(vcs.GitVCS)
						_, err = vcs.Fetch(gvcs, path.Join(opt.TempVCSPath, remoteVCSRepoName), deployment.Git, opt.SSHKeyPath)
						if err != nil {
							log.Error(err.Error())
							stateCapture.DeploymentState[deployment.Name] = state.Details{State: state.EDeploymentStateError}

						}
					} else {
						log.Debug(fmt.Sprintf("Using existing repository %s", path.Join(opt.TempVCSPath, remoteVCSRepoName)))
					}
					// Run actions ---------------
					for _, a := range deployment.Action {

						// Switch action based on Kind
						switch strings.ToLower(a.Execute.Kind) {

						case "shell":

							command, ok := a.Execute.Configuration["Command"]
							if !ok {
								log.Warn("Could not run the shell step as Command could not be found")
								continue
							}

							p := path.Join(opt.TempVCSPath, remoteVCSRepoName)

							tp, ok := a.Execute.Configuration["Path"]
							if ok {
								p = tp
							}

							log.Warn(fmt.Sprintf("Running shell command %s\n", command))
							if err := ShellCommand(command, p, true); err != nil {
								log.Error(err.Error())
							}

						case "kubernetes":

							var deploymentPath = "."

							if tp, ok := a.Execute.Configuration["Path"]; ok && tp != "" {
								deploymentPath = tp

							}
							// Deploy -------------------------
							fileList := []string{}
							err := filepath.Walk(path.Join(opt.TempVCSPath,
								remoteVCSRepoName, deploymentPath),
								func(path string, f os.FileInfo, err error) error {
									fileList = append(fileList, path)
									return nil
								})
							if err != nil {
								log.Error(err.Error())

							}
							err = platform.GenerateDeploymentPlan(restclient,
								k8siface, fileList,
								deployment.Namespace, opt, commandFlag)
							if err != nil {
								log.Error(err.Error())
							}
							//---------------------------------
						}

					}

					wg.Done()
				},
			}
		}
	}
	wg.Wait()
	coordinator.Destroy()

	return stateCapture
}
