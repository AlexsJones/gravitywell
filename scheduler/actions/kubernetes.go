package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform"
	logger "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func clientForCluster(clusterName string) (*rest.Config, kubernetes.Interface) {
	logger.Info(fmt.Sprintf("Switching to cluster: %s\n", clusterName))
	restclient, k8siface, err := platform.GetKubeClient(clusterName)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	return restclient, k8siface
}

func ExecuteKubernetesAction(action kinds.Execute, clusterInformation struct {
	ClusterName string
	ClusterRegion string
	ClusterProjectName string
	ClusterProviderName string
},
	deployment kinds.Application,
	commandFlag configuration.CommandFlag, opt configuration.Options, repoName string) {
	var deploymentPath = "."
	shouldAwaitDeployment := false
	if tp, ok := action.Configuration["Path"]; ok && tp != "" {
		deploymentPath = tp
	}
	if tp, ok := action.Configuration["AwaitDeployment"]; ok && tp != "" {
		b, err := strconv.ParseBool(tp)
		if err != nil {
			logger.Error(err.Error())
		}
		shouldAwaitDeployment = b
	}

	fileList := []string{}
	err := filepath.Walk(path.Join(opt.TempVCSPath,
		repoName, deploymentPath),
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if f.IsDir(){
				logger.Info("Ignoring directory %s",fmt.Sprintf(path))
				return nil
			}
			fileList = append(fileList, path)
			return nil
		})
	if err != nil {
		logger.Fatal(err.Error())
	}
	var clusterName = clusterInformation.ClusterName
	//WARNING: Cloud provider sensitive behaviour is required here to determine the correct context name.

	switch strings.ToLower(clusterInformation.ClusterProviderName) {
	case "google cloud platform":
		clusterName = fmt.Sprintf("gke_%s_%s_%s",clusterInformation.ClusterProjectName,
			clusterInformation.ClusterRegion, clusterInformation.ClusterName)
	case "amazon web services":

	}

	_, k8siface := clientForCluster(clusterName)
	err = platform.GenerateDeploymentPlan(
		k8siface, fileList,
		deployment.Namespace, opt, commandFlag, shouldAwaitDeployment)
	if err != nil {
		logger.Error(err.Error())
	}
	//---------------------------------
}
