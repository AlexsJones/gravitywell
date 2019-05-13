package minikube

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/google/logger"
	"strings"
)

func (m *MiniKubeProvider) Create(cluster kinds.ProviderCluster) error {

	minikubeConnection := []string{"minikube start"}

	if cluster.NodeConfiguration.Memory != 0 {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("--memory=%d",cluster.NodeConfiguration.Memory))
	}

	if cluster.NodeConfiguration.CPU != 0 {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("--cpus=%d",cluster.NodeConfiguration.CPU))
	}

	if cluster.NodeConfiguration.VMDriver != "" {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("--vm-driver=%s",cluster.NodeConfiguration.VMDriver))
	}

	if cluster.NodeConfiguration.DiskSize != "" {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("--disk-size=%s",cluster.NodeConfiguration.DiskSize))
	}

	if len(cluster.NodeConfiguration.ExtraConfiguration.ApiserverEnableAdmissionPlugins) != 0 {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("--extra-config=apiserver.enable-admission-plugins=%s",
			strings.Join(cluster.NodeConfiguration.ExtraConfiguration.ApiserverEnableAdmissionPlugins,",")))
	}

	if strings.ToLower(cluster.FullName) != strings.ToLower(cluster.ShortName) {
		logger.Error("Minikube requires FullName & ShortName to match")
		return errors.New("minikube Cluster name invalid")
	}

	if strings.ToLower(cluster.FullName) != "minikube" || strings.ToLower(cluster.ShortName) != "minikube" {
		minikubeConnection = append(minikubeConnection,fmt.Sprintf("-p=%s",cluster.FullName))
	}

	command := strings.Join(minikubeConnection," ")

	logger.Info(fmt.Sprintf("Running shell command %s\n", command))

	if err := shell.ShellCommand(command, ".", true); err != nil {
		logger.Error(err.Error())
	}

	return nil
}
