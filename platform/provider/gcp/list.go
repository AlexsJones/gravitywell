package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func List(c *container.ClusterManagerClient, ctx context.Context, projectName string) {

	clusterReq := &containerpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/-", projectName),
	}

	clusterResponse, err:= c.ListClusters(ctx,clusterReq)
	if err != nil {
		// TODO: Handle error.
		color.Red(err.Error())
		return
	}

	for _, cluster := range clusterResponse.Clusters {
		req := &containerpb.ListNodePoolsRequest{
			Parent:fmt.Sprintf("projects/%s/locations/%s/clusters/%s",projectName,cluster.Location,cluster.Name),
		}


		color.Green(fmt.Sprintf("Cluster %s located in %s status: %s\n",cluster.Name, cluster.Location,cluster.Status))
		resp, err := c.ListNodePools(ctx, req)
		if err != nil {
			continue
		}
		for _, np := range resp.NodePools {
			color.Blue(fmt.Sprintf("\t%s %d * %s\n",np.Name,np.InitialNodeCount,np.Config.MachineType))

		}
	}
}