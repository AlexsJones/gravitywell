package digitalocean

import (
	"context"
	"github.com/digitalocean/godo"
)

type DigitalOceanProvider struct {
	Context              context.Context
	ClusterManagerClient *godo.Client
}
