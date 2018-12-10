package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)


func Create(c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string,clusterName string, locations []string,initialNodeCount int32,
	initialNodeType string,
	nodePools []*containerpb.NodePool) {

	clusterReq := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectName,
			locationName),
		Cluster:&containerpb.Cluster{
			Name: clusterName,
			Locations: locations,
			NodePools: nodePools,
			InitialNodeCount: initialNodeCount,
			NodeConfig : &containerpb.NodeConfig{
				MachineType: initialNodeType,
			},
		},
	}

	clusterResponse, err:= c.CreateCluster(ctx,clusterReq)
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Blue(fmt.Sprintf("Started cluster build at %s",clusterResponse.StartTime))

}
