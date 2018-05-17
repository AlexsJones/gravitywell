package scheduler

import (
	"errors"

	"github.com/AlexsJones/ashara/configuration"
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

//Design a new scheduler based off of the current configuration
func (s *Scheduler) Design() error {

	return nil
}
