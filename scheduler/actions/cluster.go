package actions

import (
	"context"
	"os"

	container "cloud.google.com/go/container/apiv1"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/AlexsJones/gravitywell/platform/provider/gcp"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/AlexsJones/gravitywell/vault"
	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
)

func runGCPCreate(cmc *container.ClusterManagerClient, ctx context.Context,
	cluster kinds.ProviderCluster) (string, string, error) {

	var convertedNodePool []*containerpb.NodePool

	for _, model := range cluster.NodePools {
		nodePool := new(containerpb.NodePool)
		nodePool.Name = model.NodePool.Name
		nodePool.Config = new(containerpb.NodeConfig)
		nodePool.Config.MachineType = model.NodePool.NodeType
		nodePool.InitialNodeCount = int32(model.NodePool.Count)

		var labels = map[string]string{}
		cluster.Labels["project"] = cluster.Project
		cluster.Labels["region"] = cluster.Region
		cluster.Labels["cluster"] = cluster.ShortName

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

	return gcp.Create(cmc, ctx, cluster, convertedNodePool)
}
func runGCPDelete(cmc *container.ClusterManagerClient, ctx context.Context,
	cluster kinds.ProviderCluster) error {

	return gcp.Delete(cmc, ctx, cluster.Project, cluster.Region, cluster.ShortName)

}
func GoogleCloudClusterProcessor(commandFlag configuration.CommandFlag,
	cluster kinds.ProviderCluster, opt configuration.Options) {

	ctx := context.Background()
	cmc, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	create := func() {

		clusterEndpoint, clusterCertCa, err := runGCPCreate(cmc, ctx, cluster)
		if err != nil {
			color.Red(err.Error())
		}
		cluster.Endpoint = clusterEndpoint
		cluster.CertCa = clusterCertCa

		if err := platform.SetK8SContext(cluster); err != nil {
			color.Red(err.Error())
		}

		if err := vault.SetVaultConfiguration(opt, cluster); err != nil {
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
		if err := runGCPDelete(cmc, ctx, cluster); err != nil {
			color.Red(err.Error())
		}

		if err := platform.UnSetK8SContext(cluster); err != nil {
			color.Red(err.Error())
		}

		if err = vault.UnSetVaultUrl(cluster); err != nil {
			color.Red(err.Error())
		}
		if err = vault.UnSetVaultGit(opt, cluster); err != nil {
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
