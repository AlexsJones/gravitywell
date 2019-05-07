package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/google/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"strings"
)

func ApplicationProcessor(commandFlag configuration.CommandFlag,
	opt configuration.Options, clusterName string, application kinds.Application) {

	executeDeployment(application, opt, clusterName, commandFlag)
}

func selectAndExecute(execute kinds.Execute, deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag, repoName string) {

	switch strings.ToLower(execute.Kind) {
	case "shell":
		ExecuteShellAction(execute, opt, repoName)

	case "kubernetes":
		ExecuteKubernetesAction(execute, clusterName, deployment, commandFlag, opt, repoName)
	}
}
func loadActionList(path string) kinds.ActionList {

	logger.Info(fmt.Sprintf("Loading %s", path))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}
	appc := kinds.ActionList{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		logger.Error(err)
	}
	return appc
}

func executeDeployment(deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag) {
	logger.Info(fmt.Sprintf("Loading deployment %s\n", deployment.Name))

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt)
	if err != nil {
		logger.Error(err.Error())

		return
	}
	//1. Run inline action lists first

	for _, a := range deployment.ActionList.Executions {

		selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
	}

	//2. Run local path action lists second
	if deployment.ActionList.LocalPath != "" {

		appc := loadActionList(deployment.ActionList.LocalPath)

		for _, a := range appc.Executions {

			selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
		}
	}
	//3. Run remote action lists last
	if deployment.ActionList.RemotePath != "" {

		appc := loadActionList(path.Join(opt.TempVCSPath, deployment.Name, deployment.ActionList.RemotePath))

		for _, a := range appc.Executions {

			selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
		}
	}

}
