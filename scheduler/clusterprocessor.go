package scheduler

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"fmt"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

func runGCPCreate(cmc *container.ClusterManagerClient,ctx context.Context,
	cluster configuration.ProviderCluster) error {

	var convertedNodePool []*containerpb.NodePool

	 for _, model := range cluster.NodePools {
			 nodePool := new(containerpb.NodePool)
			 nodePool.Name = model.NodePool.Name
			 nodePool.Config = new(containerpb.NodeConfig)
			 nodePool.Config.MachineType = model.NodePool.NodeType
			 nodePool.InitialNodeCount = int32(model.NodePool.Count)
		 }
	return gcp.Create(cmc,ctx,cluster.Project,
		cluster.Region, cluster.Name,
		cluster.Zones,
	int32(cluster.InitialNodeCount),
	cluster.InitialNodeType,
		convertedNodePool)
}
func runGCPDelete(cmc *container.ClusterManagerClient,ctx context.Context,
	cluster configuration.ProviderCluster) error {

		return gcp.Delete(cmc,ctx,cluster.Project,cluster.Region,cluster.Name)

}
func ClusterProcessor(commandFlag configuration.CommandFlag,
	provider configuration.Provider) {

	if provider.Name == "" {
		log.Warn("Provider requires a name")
		os.Exit(1)
	}
	switch strings.ToLower(provider.Name) {

	case "google cloud platform":
		ctx := context.Background()
		cmc, err := container.NewClusterManagerClient(ctx)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		// Run Command ------------------------------------------------------------------
		if commandFlag == configuration.Create || commandFlag == configuration.Apply {

			for _, cluster := range provider.Clusters {
				err := runGCPCreate(cmc,ctx,cluster.Cluster)
				if err != nil {
					color.Red(err.Error())
				}
				// Run post install -----------------------------------------------------
				for _, executeCommand := range cluster.Cluster.PostInstallHook {
					if executeCommand.Execute.Shell != "" {
						err := ShellCommand(executeCommand.Execute.Shell,
							executeCommand.Execute.Path,false)
						if err != nil {
							color.Red(err.Error())
						}
						}
				}
			}
		}
		if commandFlag == configuration.Delete {
			for _, cluster := range provider.Clusters {
				err := runGCPDelete(cmc,ctx,cluster.Cluster)
				if err != nil {
					color.Red(err.Error())
					continue
				}
				// Run post delete -----------------------------------------------------
				for _, executeCommand := range cluster.Cluster.PostDeleteHooak {
					if executeCommand.Execute.Shell != "" {
						err := ShellCommand(executeCommand.Execute.Shell,
							executeCommand.Execute.Path,false)
						if err != nil {
							color.Red(err.Error())
						}
					}
				}
			}
		}

	case "amazon web services":
		log.Warn("Amazon Web Services not yet supported")
		os.Exit(1)
	default:
		log.Warn(fmt.Sprintf("Platform %s not supported", provider.Name))
	}
}