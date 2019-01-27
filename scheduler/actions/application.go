package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/vcs"
	log "github.com/Sirupsen/logrus"
	"strings"
)

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, clusterName string, application kinds.Application) {

	executeDeployment(application, opt, clusterName, commandFlag)
}

func groupDeploymentsPerNamespace(cluster kinds.ApplicationCluster) map[string][]kinds.Application {
	groupedDeployments := make(map[string][]kinds.Application)
	for _, deployment := range cluster.Applications {
		groupedDeployments[deployment.Application.Namespace] = append(groupedDeployments[deployment.Application.Namespace], deployment.Application)
	}

	for namespace, deployments := range groupedDeployments {
		fmt.Printf("Deployments for namespace %s on cluster %s\t\t\n", namespace, cluster.ShortName)
		for _, depl := range deployments {
			fmt.Printf("\t\t %s\n", depl.Name)
		}
	}
	return groupedDeployments
}

func executeDeployment(deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag) {
	log.Debug(fmt.Sprintf("Loading deployment %s\n", deployment.Name))

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt)
	if err != nil {
		log.Error(err.Error())

		return
	}

	for _, a := range deployment.Action {
		switch strings.ToLower(a.Execute.Kind) {
		case "shell":
			ExecuteShellAction(a, opt, remoteVCSRepoName)

		case "kubernetes":
			ExecuteKubernetesAction(a, clusterName, deployment, commandFlag, opt, remoteVCSRepoName)
		}

	}
}
