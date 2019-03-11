package gcp

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"time"

	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func (g *GCPProvider)Create(projectName string,
	locationName string, clusterName string, locations []string, initialNodeCount int32,
	initialNodeType string, clusterLabels map[string]string,
	nodePools []kinds.NodePool) error {

	//Convert generic node pool type into a specific GCP resource
	var convertedNodePool []*containerpb.NodePool

	for _, model := range nodePools {
		nodePool := new(containerpb.NodePool)
		nodePool.Name = model.Name
		nodePool.Config = new(containerpb.NodeConfig)
		nodePool.Config.MachineType = model.NodeType
		nodePool.InitialNodeCount = int32(model.Count)

		var labels = map[string]string{}

		if len(clusterLabels) > 0 {
			for index, element := range clusterLabels {
				labels[index] = element
			}
		}

		if len(model.Labels) > 0 {
			for index, element := range model.Labels {
				labels[index] = element
			}
		}
		nodePool.Config.Labels = labels

		convertedNodePool = append(convertedNodePool, nodePool)
	}

	// ----------------------------------------------------------------
	var cluster *containerpb.Cluster
	if len(nodePools) == 0 {

		cluster = &containerpb.Cluster{
			Name:             clusterName,
			Locations:        locations,
			InitialNodeCount: initialNodeCount,
			NodeConfig: &containerpb.NodeConfig{
				MachineType: initialNodeType,
			},
			ResourceLabels: clusterLabels,
		}

	} else {
		cluster = &containerpb.Cluster{
			Name:           clusterName,
			Locations:      locations,
			NodePools:      convertedNodePool,
			ResourceLabels: clusterLabels,
		}
	}
	clusterReq := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", projectName,
			locationName),
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
			g.ClusterManagerClient.GetCluster(g.Context, &containerpb.GetClusterRequest{Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectName,
				locationName, clusterName)})
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
