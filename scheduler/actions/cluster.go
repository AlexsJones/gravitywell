package actions

import (
	"cloud.google.com/go/container/apiv1"
	"context"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"os"
)


func GoogleCloudClusterProcessor(commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster) {

	gcpProviderClient := &gcp.GCPProvider{}

	ctx := context.Background()
	cmc, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	create := func() {

		err := provider.Create(gcpProviderClient,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName,
			cluster.Zones,
			int32(cluster.InitialNodeCount),
			cluster.InitialNodeType,
			cluster.Labels, cluster.NodePools)

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
		err := provider.Delete(gcpProviderClient,cmc, ctx, cluster.Project,
			cluster.Region, cluster.ShortName)
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
