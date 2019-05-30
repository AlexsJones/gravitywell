package minikube

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	logger "github.com/sirupsen/logrus"
	"strings"
)

func (g *MiniKubeProvider) Delete(cluster kinds.ProviderCluster) error {

	minikubeConnection := []string{"minikube delete"}

	if strings.ToLower(cluster.Name) != "minikube" {
		minikubeConnection = append(minikubeConnection, fmt.Sprintf("-p=%s", cluster.Name))
	}

	command := strings.Join(minikubeConnection, " ")

	logger.Info(fmt.Sprintf("Running shell command %s\n", command))

	if err := shell.ShellCommand(command, ".", true); err != nil {
		logger.Error(err.Error())

		return err
	}

	return nil
}
