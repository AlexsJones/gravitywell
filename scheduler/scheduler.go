package scheduler

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/fatih/color"
)

type Options struct {
	VCS        string
	APIVersion string
}

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
func (s *Scheduler) Run(opt Options) error {

	if opt.APIVersion != s.configuration.APIVersion {
		color.Red(fmt.Sprintf("Manifest is not supported by the current API: %s\n", opt.APIVersion))
		os.Exit(1)
	}
	for _, cluster := range s.configuration.Strategy {

		color.Yellow(fmt.Sprintf("Switching to cluster: %s\n", cluster.Cluster.Name))

		_, _, err := platform.GetKubeClient(cluster.Cluster.Name)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		color.Yellow("Fetching remote.")

	}

	return nil
}
