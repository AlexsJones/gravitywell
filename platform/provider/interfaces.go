package provider

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"github.com/AlexsJones/gravitywell/kinds"
)

type IProvider interface {
	Create(c *container.ClusterManagerClient, ctx context.Context, projectName string,
		locationName string, clusterName string, locations []string, initialNodeCount int32,
		initialNodeType string, clusterLabels map[string]string, nodePools []kinds.NodePool) error
	Delete(c *container.ClusterManagerClient, ctx context.Context, projectName string,
		locationName string, clusterName string) error
	List(c *container.ClusterManagerClient, ctx context.Context, projectName string) error
}

func Create(i IProvider,c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string, clusterName string, locations []string, initialNodeCount int32,
	initialNodeType string, clusterLabels map[string]string, nodePools []kinds.NodePool) error{

	return i.Create(c,ctx,projectName,
		locationName,clusterName,
		locations,initialNodeCount,
		initialNodeType,
		clusterLabels,nodePools)
}

func Delete(i IProvider, c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string, clusterName string) error {

		return i.Delete(c,ctx,projectName,locationName,clusterName)
}

func List(i IProvider, c *container.ClusterManagerClient, ctx context.Context, projectName string) error {
	return i.List(c,ctx,projectName)
}