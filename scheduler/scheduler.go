package scheduler

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/platform"
	"github.com/fatih/color"
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
func (s *Scheduler) Run(kubernetes *platform.Kubernetes) error {

	for _, cluster := range s.configuration.Strategy {
		color.Yellow(fmt.Sprintf("Switching to cluster: %s\n", cluster.Cluster.Name))
	}

	return nil
}
