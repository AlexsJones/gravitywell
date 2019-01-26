package actions

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"os"
)

func runGCPCreate(cmc *container.ClusterManagerClient, ctx context.Context,
	cluster kinds.ProviderCluster) error {

	var convertedNodePool []*containerpb.NodePool

	for _, model := range cluster.NodePools {
		nodePool := new(containerpb.NodePool)
		nodePool.Name = model.NodePool.Name
		nodePool.Config = new(containerpb.NodeConfig)
		nodePool.Config.MachineType = model.NodePool.NodeType
		nodePool.InitialNodeCount = int32(model.NodePool.Count)

		var labels = map[string]string{}

		if len(cluster.Labels) > 0 {
			for index, element := range cluster.Labels {
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

	return gcp.Create(cmc, ctx, cluster.Project,
		cluster.Region, cluster.ShortName,
		cluster.Zones,
		int32(cluster.InitialNodeCount),
		cluster.InitialNodeType,
		cluster.Labels,
		convertedNodePool)
}
func runGCPDelete(cmc *container.ClusterManagerClient, ctx context.Context,
	cluster kinds.ProviderCluster) error {

	return gcp.Delete(cmc, ctx, cluster.Project, cluster.Region, cluster.ShortName)

}
func GoogleCloudClusterProcessor(commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) {

	ctx := context.Background()
	cmc, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	create := func() {

		err := runGCPCreate(cmc, ctx, cluster)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}

			}
		}
	}
	delete := func() {
		err := runGCPDelete(cmc, ctx, cluster)
		if err != nil {
			color.Red(err.Error())
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHooak {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					color.Red(err.Error())
				}
			}
		}

	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		create()
	case configuration.Apply:
		create()
	case configuration.Replace:
		delete()
		create()
	case configuration.Delete:
		delete()
	}
}
