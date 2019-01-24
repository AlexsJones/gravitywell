package scheduler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/AlexsJones/gravitywell/actions"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/AlexsJones/gravitywell/subprocessor"
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
	//Batch actions per namespace
	deploymentsPerNamespace := groupDeploymentsPerNamespace(cluster)

	//Execute deployments
	coordinator := subprocessor.NewCoordinator()
	var wg sync.WaitGroup
	go coordinator.Run()

	for namespace, deployments := range deploymentsPerNamespace {
		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", namespace, cluster.Name)

		wg.Add(1)
		coordinator.ResourceChannel <- subprocessor.Resource{
			Process: func() {
				for _, deploy := range deployments {
					fmt.Printf("Deploying:\t\t %s\n", deploy.Name)
					executeDeployment(deploy, opt, stateCapture, cluster.Name, commandFlag)
				}
				wg.Done()
			},
		}
	}
	wg.Wait()
	coordinator.Destroy()

	return stateCapture
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

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt, "master")
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
			actions.ExecuteKubernetesAction(a, clusterName, deployment, commandFlag, opt, remoteVCSRepoName)
		}

	}
}
