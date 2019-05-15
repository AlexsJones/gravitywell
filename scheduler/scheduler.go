package scheduler

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"github.com/AlexsJones/gravitywell/scheduler/planner/standard"
	logger "github.com/sirupsen/logrus"
)

//Scheduler object ...
type Scheduler struct {
	configuration *configuration.Configuration
}

//NewScheduler object ...
func NewScheduler(conf *configuration.Configuration) (*Scheduler, error) {
	if conf == nil {
		return nil, errors.New("invalid configuration")
	}
	return &Scheduler{
		configuration: conf}, nil
}

//Run a new scheduler based off of the current configuration
func (s *Scheduler) Run(commandFlag configuration.CommandFlag,
	opt configuration.Options) error {

	stdplnr := standard.StandardPlanner{}

	plan, err := planner.GeneratePlan(stdplnr, s.configuration, commandFlag, opt)
	if err != nil {
		logger.Fatal(err)
	}

	statusWatcher := plan.Run()
	for {
		select {
		case msg := <-statusWatcher:
			if msg.Halt() {
				//Halting
				fmt.Println("Received halt")
				return nil
			}
		}
	}
}
