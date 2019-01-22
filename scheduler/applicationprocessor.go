package scheduler

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/AlexsJones/gravitywell/subprocessor"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/AlexsJones/gravitywell/actions"
	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

func clientForCluster(clusterName string) (*rest.Config, kubernetes.Interface) {
	log.Info(fmt.Sprintf("Switching to cluster: %s\n", clusterName))
	restclient, k8siface, err := platform.GetKubeClient(clusterName)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	return restclient, k8siface
}

func groupDeploymentsPerNamespace(cluster configuration.ApplicationCluster) map[string][]configuration.Application {
	groupedDeployments := make(map[string][]configuration.Application)
	for _, deployment := range cluster.Applications {
		groupedDeployments[deployment.Application.Namespace] = append(groupedDeployments[deployment.Application.Namespace], deployment.Application)
	}

	for namespace, deployments := range groupedDeployments {
		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", namespace, cluster.Name)
		for _, depl := range deployments {
			fmt.Printf("\t\t %s\n", depl.Name)
		}
	}
	return groupedDeployments
}

func executeDeployment(deployment configuration.Application, opt configuration.Options, stateCapture *state.Capture, clusterName string, commandFlag configuration.CommandFlag) {
	log.Debug(fmt.Sprintf("Loading deployment %s\n", deployment.Name))

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt)
	if err != nil {
		log.Error(err.Error())
		stateCapture.DeploymentState[deployment.Name] = state.Details{State: state.EDeploymentStateError}
		return
	}

	for _, a := range deployment.Action {
		switch strings.ToLower(a.Execute.Kind) {
		case "shell":
			actions.ExecuteShellAction(a, opt, remoteVCSRepoName)

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
			restclient, k8siface := clientForCluster(clusterName)
			err = platform.GenerateDeploymentPlan(restclient,
				k8siface, fileList,
				deployment.Namespace, opt, commandFlag)
			if err != nil {
				log.Error(err.Error())
			}
			//---------------------------------
		}

	}
}

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, cluster configuration.ApplicationCluster) *state.Capture {

	//TODO: Remove this state capture stuff, it's clunkys
	stateCapture := &state.Capture{
		ClusterName:     cluster.Name,
		DeploymentState: make(map[string]state.Details),
	}
	//Batch actions per namespace
	deploymentsPerNamespace := groupDeploymentsPerNamespace(cluster)

	//Execute deployments
	coordinator := subprocessor.NewCoordinator()
	var wg sync.WaitGroup
	go coordinator.Run()

	for namespace, deployments := range deploymentsPerNamespace {
		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", namespace, cluster.Name)

		for _, deploy := range deployments {
			fmt.Printf("Deploying:\t\t %s\n", deploy.Name)
			wg.Add(1)
			coordinator.ResourceChannel <- subprocessor.Resource{
				Process: func() {
					executeDeployment(deploy, opt, stateCapture, cluster.Name, commandFlag)
					wg.Done()
				},
			}
		}
	}
	wg.Wait()
	coordinator.Destroy()

	return stateCapture
}
