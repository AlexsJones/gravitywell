package aws

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"github.com/AlexsJones/gravitywell/kinds"
)


func (AWSProvider) Create(c *container.ClusterManagerClient, ctx context.Context, projectName string,
	locationName string, clusterName string, locations []string, initialNodeCount int32,
	initialNodeType string, clusterLabels map[string]string,
	nodePools []kinds.NodePool) error {
return nil
}
