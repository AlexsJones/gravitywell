package actions

import (
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/platform/provider"
	"github.com/AlexsJones/gravitywell/platform/provider/minikube"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	logger "github.com/sirupsen/logrus"
)


func NewMinikubeConfig() (*minikube.MiniKubeProvider, error) {

	return nil, nil
}
func MinikubeClusterProcessor(minikubeprovider *minikube.MiniKubeProvider,
	commandFlag configuration.CommandFlag, cluster kinds.ProviderCluster) error {

	create := func() error {

		err := provider.Create(minikubeprovider, cluster)

		if err != nil {
			return err
		}
		// Run post install -----------------------------------------------------
		for _, executeCommand := range cluster.PostInstallHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}

			}
		}
		return nil
	}
	delete := func() error {
		err := provider.Delete(minikubeprovider, cluster)
		if err != nil {
			return err
		}
		// Run post delete -----------------------------------------------------
		for _, executeCommand := range cluster.PostDeleteHook {
			if executeCommand.Execute.Shell != "" {
				err := shell.ShellCommand(executeCommand.Execute.Shell,
					executeCommand.Execute.Path, false)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	// Run Command ------------------------------------------------------------------
	switch commandFlag {
	case configuration.Create:
		return create()
	case configuration.Apply:
		logger.Info("Ignoring apply on cluster - no such option")
		return nil
	case configuration.Replace:
		if err := delete(); err != nil {
			return err
		}
		return create()
	case configuration.Delete:
		return delete()
	}
	return nil
}
