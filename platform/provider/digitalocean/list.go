package digitalocean

import (
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/digitalocean/godo"
)

func (g *DigitalOceanProvider) List(clusterp kinds.ProviderCluster) ([]*godo.KubernetesCluster, *godo.Response, error) {

	kls, resp, err := g.ClusterManagerClient.Kubernetes.List(g.Context,
		&godo.ListOptions{})
	if err != nil {
		return nil, resp, err
	}
	return kls, resp, nil
}
