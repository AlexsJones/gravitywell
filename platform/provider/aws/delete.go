package aws

import (
	"cloud.google.com/go/container/apiv1"
	"context"
)

func (AWSProvider)Delete(c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string,
	clusterName string) error {

		return nil
}
