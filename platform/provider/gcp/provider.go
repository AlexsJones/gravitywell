package gcp

import (
	"cloud.google.com/go/container/apiv1"
	"context"
)

type GCPProvider struct {
	Context              context.Context
	ClusterManagerClient *container.ClusterManagerClient
}
