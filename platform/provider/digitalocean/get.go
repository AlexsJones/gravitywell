package digitalocean

import (
	"github.com/digitalocean/godo"
)

func (g *DigitalOceanProvider) Get(id string) (*godo.KubernetesCluster, *godo.Response, error) {

	kls, resp, err := g.ClusterManagerClient.Kubernetes.Get(g.Context,
		id)
	if err != nil {
		return nil, resp, err
	}
	return kls, resp, nil
}
