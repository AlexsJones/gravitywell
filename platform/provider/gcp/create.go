package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"time"
)


func Create(c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string,clusterName string, locations []string,initialNodeCount int32,
	initialNodeType string,
	nodePools []*containerpb.NodePool) error {

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
		return err
	}
	color.Blue(fmt.Sprintf("Started cluster build at %s",clusterResponse.StartTime))

	for {
		clust, err :=
			c.GetCluster(ctx,&containerpb.GetClusterRequest{Name:
			fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectName,
				locationName,clusterName)})
		if err != nil {
			return err
		}
		if clust.GetStatus() == containerpb.Cluster_RUNNING {
			color.Green("Cluster running")
			return nil
		}
		time.Sleep(time.Second)
	}

}
