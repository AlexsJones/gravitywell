package minikube

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	logger "github.com/sirupsen/logrus"
	"strings"
)

func (m *MiniKubeProvider) Create(cluster kinds.ProviderCluster) error {

	minikubeConnection := []string{"minikube start"}

	if cluster.NodeConfiguration.Memory != 0 {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("--memory=%d", cluster.NodeConfiguration.Memory))
	}

	if cluster.NodeConfiguration.CPU != 0 {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("--cpus=%d", cluster.NodeConfiguration.CPU))
	}

	if cluster.NodeConfiguration.VMDriver != "" {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("--vm-driver=%s", cluster.NodeConfiguration.VMDriver))
	}

	if cluster.NodeConfiguration.DiskSize != "" {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("--disk-size=%s", cluster.NodeConfiguration.DiskSize))
	}

	if len(cluster.NodeConfiguration.ExtraConfiguration.ApiserverEnableAdmissionPlugins) != 0 {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("--extra-config=apiserver.enable-admission-plugins=%s",
			strings.Join(cluster.NodeConfiguration.ExtraConfiguration.ApiserverEnableAdmissionPlugins, ",")))
	}

	if strings.ToLower(cluster.Name) != "minikube" {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("-p=%s", cluster.Name))
	}

	command := strings.Join(minikubeConnection, " ")

	logger.Info(fmt.Sprintf("Running shell command %s\n", command))

	if err := shell.ShellCommand(command, ".", true); err != nil {
		logger.Info(err.Error())
	}

	return nil
}
