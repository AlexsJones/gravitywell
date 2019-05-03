package gcp

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"time"

	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func (g *GCPProvider) Create(clusterp kinds.ProviderCluster) error {

	//Convert generic node pool type into a specific GCP resource
	var convertedNodePool []*containerpb.NodePool

	for _, model := range clusterp.NodePools {
		nodePool := new(containerpb.NodePool)
		nodePool.Name = model.NodePool.Name
		nodePool.InitialNodeCount = int32(model.NodePool.Count)
		nodePool.Config = &containerpb.NodeConfig{
			MachineType: clusterp.InitialNodeType,
			OauthScopes: clusterp.OauthScopes,
		}
		var labels = map[string]string{}

		if len(clusterp.Labels) > 0 {
			for index, element := range clusterp.Labels {
				labels[index] = element
			}
		}

		if len(model.NodePool.Labels) > 0 {
			for index, element := range model.NodePool.Labels {
				labels[index] = element
			}
		}
		nodePool.Config.Labels = labels

		convertedNodePool = append(convertedNodePool, nodePool)
	}

	// ----------------------------------------------------------------
	var cluster *containerpb.Cluster
	if len(clusterp.NodePools) == 0 {

		cluster = &containerpb.Cluster{
			Name:             clusterp.ShortName,
			Locations:        clusterp.Zones,
			InitialNodeCount: int32(clusterp.InitialNodeCount),
			NodeConfig: &containerpb.NodeConfig{
				MachineType: clusterp.InitialNodeType,
				OauthScopes: clusterp.OauthScopes,
			},
			ResourceLabels: clusterp.Labels,
		}

	} else {
		cluster = &containerpb.Cluster{
			Name:           clusterp.ShortName,
			Locations:      clusterp.Zones,
			NodePools:      convertedNodePool,
			ResourceLabels: clusterp.Labels,
		}
	}
	clusterReq := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", clusterp.Project,
			clusterp.Region),
		Cluster: cluster,
	}

	clusterResponse, err := g.ClusterManagerClient.CreateCluster(g.Context, clusterReq)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	color.Blue(fmt.Sprintf("Started cluster build at %s", clusterResponse.StartTime))

	for {
		clust, err :=
			g.ClusterManagerClient.GetCluster(g.Context, &containerpb.GetClusterRequest{Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", clusterp.Project,
				clusterp.Region, clusterp.ShortName)})
		if err != nil {
			return err
		}
		if clust.GetStatus() == containerpb.Cluster_RUNNING {
			color.Green("Cluster running")
			return nil
		}
		time.Sleep(time.Second)
	}

}
