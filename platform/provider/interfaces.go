package provider

import (
	"github.com/AlexsJones/gravitywell/kinds"
)

type IProvider interface {
	Create(projectName string,
		locationName string, clusterName string, locations []string, initialNodeCount int32,
		initialNodeType string, clusterLabels map[string]string, nodePools []kinds.NodePool) error
	Delete(projectName string, locationName string, clusterName string) error
	List(projectName string) error
}

func Create(i IProvider,  projectName string,
	locationName string, clusterName string, locations []string, initialNodeCount int32,
	initialNodeType string, clusterLabels map[string]string, nodePools []kinds.NodePool) error{

	return i.Create(projectName,
		locationName,clusterName,
		locations,initialNodeCount,
		initialNodeType,
		clusterLabels,nodePools)
}

func Delete(i IProvider, projectName string, locationName string, clusterName string) error {

	return i.Delete(projectName,locationName,clusterName)
}

func List(i IProvider,  projectName string) error {
	return i.List(projectName)
}