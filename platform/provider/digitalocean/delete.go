package digitalocean

import (
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/shared"
	logger "github.com/sirupsen/logrus"
)

func (g *DigitalOceanProvider) Delete(clusterp kinds.ProviderCluster) error {

	clusters, _, err := g.List(clusterp)
	if err != nil {
		return err
	}

	for _, cluster := range clusters {
		if cluster.Name == clusterp.Name {
			//Found Id
			resp, err := g.ClusterManagerClient.Kubernetes.Delete(g.Context,
				cluster.ID)
			if err != nil {
				return err
			}
			logger.Info(shared.PrettyPrint(resp))
		}
	}
	return nil
}
