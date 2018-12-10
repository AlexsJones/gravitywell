package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func Delete(c *container.ClusterManagerClient, ctx context.Context,projectName string,
	locationName string,
	clusterName string) {

	clusterReq := &containerpb.DeleteClusterRequest{

		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectName,locationName,clusterName),
	}

	clusterResponse, err:= c.DeleteCluster(ctx,clusterReq)
	if err != nil {
		color.Red(err.Error())
		return
	}
	color.Blue(fmt.Sprintf("Started cluster deletion at %s",clusterResponse.StartTime))

}