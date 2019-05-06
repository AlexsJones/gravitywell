package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/vcs"
	"github.com/fatih/color"
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

func executeDeployment(deployment kinds.Application, opt configuration.Options,
	clusterName string, commandFlag configuration.CommandFlag) {
	logger.Info(fmt.Sprintf("Loading deployment %s\n", deployment.Name))

	remoteVCSRepoName, err := vcs.FetchRepo(deployment.Git, opt)
	if err != nil {
		logger.Error(err.Error())

		return
	}
	//Run inline action list if neither a local path nor a remote is defined
	if deployment.ActionList.LocalPath == "" && deployment.ActionList.RemotePath == "" {
		for _, a := range deployment.ActionList.Executions {

			selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
		}
		return
	}
	//Retrieve local/remote action lists - relative to the caller ApplicationKind
	if deployment.ActionList.LocalPath != "" && deployment.ActionList.RemotePath != "" {
		color.Yellow("Both local and remote action lists have been defined. Will prioritise local")
	}

	var actionListLoadPath = ""

	if deployment.ActionList.RemotePath != "" {
		//Prepend the repository local directory reference
		actionListLoadPath = path.Join(opt.TempVCSPath, deployment.Name, deployment.ActionList.RemotePath)
	}
	if deployment.ActionList.LocalPath != "" {
		actionListLoadPath = deployment.ActionList.LocalPath
	}

	if actionListLoadPath == "" {
		logger.Error("No action lists defined for %s/%s at %s", deployment.Name, deployment.Namespace, deployment.Git)
	}

	logger.Info("Using action list path %s", actionListLoadPath)

	logger.Info(fmt.Sprintf("Loading %s", actionListLoadPath))
	bytes, err := ioutil.ReadFile(actionListLoadPath)
	if err != nil {
		logger.Error(err)
	}
	appc := kinds.ActionList{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		logger.Error(err)
	}
	for _, a := range appc.Executions {

		selectAndExecute(a, deployment, opt, clusterName, commandFlag, remoteVCSRepoName)
	}
}
