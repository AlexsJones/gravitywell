package digitalocean

import (
	"github.com/AlexsJones/gravitywell/kinds"
	logger "github.com/sirupsen/logrus"
)

func (g *DigitalOceanProvider) Delete(clusterp kinds.ProviderCluster) error {

	response, err := g.ClusterManagerClient.Kubernetes.Delete(g.Context,
		clusterp.Name)
	if err != nil {
		logger.Error(response)
		return err
	}
	return nil
}
