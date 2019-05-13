package minikube

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/google/logger"
	"strings"
)

func (g *MiniKubeProvider) Delete(cluster kinds.ProviderCluster) error {

	minikubeConnection := []string{"minikube delete"}

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

		return err
	}

	return nil
}
