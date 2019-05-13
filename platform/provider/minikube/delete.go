package minikube

import (
	"github.com/AlexsJones/gravitywell/kinds"
)

func (g *MiniKubeProvider) Delete(clusterp kinds.ProviderCluster) error {

	//clusterReq := &containerpb.DeleteClusterRequest{
	//
	//	Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", clusterp.Project,
	//		clusterp.Region, clusterp.ShortName),
	//}
	//
	//clusterResponse, err := g.ClusterManagerClient.DeleteCluster(g.Context, clusterReq)
	//if err != nil {
	//	color.Red(err.Error())
	//	return err
	//}
	//color.Blue(fmt.Sprintf("Started cluster deletion at %s", clusterResponse.StartTime))
	//
	//for {
	//	_, err :=
	//		g.ClusterManagerClient.GetCluster(g.Context,
	//			&containerpb.GetClusterRequest{Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", clusterp.Project,
	//				clusterp.Region, clusterp.ShortName)})
	//
	//	if err != nil {
	//		//I know this looks awful but you need to test if the cluster is alive
	//		return nil
	//	}
	//	time.Sleep(time.Second)
	//}
	return nil
}
