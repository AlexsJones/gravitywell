package scheduler

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/fatih/color"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//Scheduler object ...
type Scheduler struct {
	configuration *configuration.Configuration
}

//NewScheduler object ...
func NewScheduler(conf *configuration.Configuration) (*Scheduler, error) {
	if conf == nil {
		return nil, errors.New("Invalid configuration")
	}
	return &Scheduler{
		configuration: conf}, nil
}

//Run a new scheduler based off of the current configuration
func (s *Scheduler) Run() error {

	for _, cluster := range s.configuration.Strategy {
		color.Yellow(fmt.Sprintf("Switching to cluster: %s\n", cluster.Cluster.Name))

		_, kiface, err := platform.GetKubeClient(cluster.Cluster.Name)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}

		nl, err := kiface.CoreV1().Namespaces().List(meta.ListOptions{})

		for _, i := range nl.Items {
			log.Println(i)
		}
	}

	return nil
}
