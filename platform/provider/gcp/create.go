package gcp

import (
	"context"
	"fmt"
	"time"

	container "cloud.google.com/go/container/apiv1"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func Create(c *container.ClusterManagerClient, ctx context.Context, cluster configuration.ProviderCluster, nodePools []*containerpb.NodePool) (string, string, error) {

	var CloudCluster *containerpb.Cluster
	if len(nodePools) == 0 {

		CloudCluster = &containerpb.Cluster{
			Name:             cluster.Name,
			Locations:        cluster.Zones,
			InitialNodeCount: int32(cluster.InitialNodeCount),
			NodeConfig: &containerpb.NodeConfig{
				MachineType: cluster.InitialNodeType,
			},
			ResourceLabels: cluster.Labels,
		}

	} else {
		CloudCluster = &containerpb.Cluster{
			Name:           cluster.Name,
			Locations:      cluster.Zones,
			NodePools:      nodePools,
			ResourceLabels: cluster.Labels,
		}
	}
	CloudClusterReq := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", cluster.Project,
			cluster.Region),
		Cluster: CloudCluster,
	}

	CloudClusterResponse, err := c.CreateCluster(ctx, CloudClusterReq)
	if err != nil {
		color.Red(err.Error())

		clust, err := c.GetCluster(ctx, &containerpb.GetClusterRequest{Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", cluster.Project, cluster.Region, cluster.Name)})
		if err != nil {
			return "", "", err
		}
		if clust.GetStatus() == containerpb.Cluster_RUNNING {
			color.Green("Cluster running")
			return clust.Endpoint, clust.MasterAuth.ClusterCaCertificate, nil
		}
		return "", "", err
	}
	color.Blue(fmt.Sprintf("Started cluster build at %s", CloudClusterResponse.StartTime))

	for {
		clust, err :=
			c.GetCluster(ctx, &containerpb.GetClusterRequest{Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", cluster.Project,
				cluster.Region, cluster.Name)})
		if err != nil {
			return "", "", err
		}
		if clust.GetStatus() == containerpb.Cluster_RUNNING {
			color.Green("Cluster running")
			return clust.Endpoint, clust.MasterAuth.ClusterCaCertificate, nil
		}
		time.Sleep(time.Second)
	}

}
