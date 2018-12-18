package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"time"
)

func Delete(c *container.ClusterManagerClient, ctx context.Context,projectName string,
	locationName string,
	clusterName string) error {

	clusterReq := &containerpb.DeleteClusterRequest{

		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectName,locationName,clusterName),
	}

	clusterResponse, err:= c.DeleteCluster(ctx,clusterReq)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	color.Blue(fmt.Sprintf("Started cluster deletion at %s",clusterResponse.StartTime))

	for {
		_, err :=
			c.GetCluster(ctx,&containerpb.GetClusterRequest{Name:
			fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectName,
				locationName,clusterName)})


		if err != nil {
			//I know this looks awful but you need to test if the cluster is alive
			return nil
		}
		time.Sleep(time.Second)
	}
}