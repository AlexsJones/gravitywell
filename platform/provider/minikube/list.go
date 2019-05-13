package minikube

import (
	"github.com/AlexsJones/gravitywell/kinds"
)

func (g *MiniKubeProvider) List(clusterp kinds.ProviderCluster) error {

	//clusterReq := &containerpb.ListClustersRequest{
	//	Parent: fmt.Sprintf("projects/%s/locations/-", clusterp.Project),
	//}
	//
	//clusterResponse, err := g.ClusterManagerClient.ListClusters(g.Context, clusterReq)
	//if err != nil {
	//	// TODO: Handle error.
	//	color.Red(err.Error())
	//	return err
	//}
	//
	//for _, cluster := range clusterResponse.Clusters {
	//	req := &containerpb.ListNodePoolsRequest{
	//		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", clusterp.Project, cluster.Location, cluster.Name),
	//	}
	//
	//	color.Green(fmt.Sprintf("Cluster %s located in %s status: %s\n", cluster.Name, cluster.Location, cluster.Status))
	//	resp, err := g.ClusterManagerClient.ListNodePools(g.Context, req)
	//	if err != nil {
	//		continue
	//	}
	//	for _, np := range resp.NodePools {
	//		color.Blue(fmt.Sprintf("\t%s %d * %s\n", np.Name, np.InitialNodeCount, np.Config.MachineType))
	//
	//	}
	//}
	return nil
}
