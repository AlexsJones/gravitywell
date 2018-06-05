package scheduler

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
)

//Options ...
type Options struct {
	VCS         string
	TempVCSPath string
	APIVersion  string
	Parallel    bool
	DryRun      bool
	TryUpdate   bool
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

func (s *Scheduler) printStatemap(m map[state.State]string) {

	var col func(string, ...interface{})

	for k, v := range m {
		if k == state.EDeploymentStateError {
			col = color.Red
		} else {
			col = color.Green
		}

		col(fmt.Sprintf("Deployment %s State => %s\n", v, state.Translate(k)))
	}
}

//Run a new scheduler based off of the current configuration
func (s *Scheduler) Run(opt Options) error {

	if opt.APIVersion != s.configuration.APIVersion {
		color.Red(fmt.Sprintf("Manifest is not supported by the current API: %s\n", opt.APIVersion))
		os.Exit(1)
	}
	//---------------------------------
	if _, err := os.Stat(opt.TempVCSPath); os.IsNotExist(err) {
		os.Mkdir(opt.TempVCSPath, 0777)
	} else {
		os.RemoveAll(opt.TempVCSPath)
		os.Mkdir(opt.TempVCSPath, 0777)
	}
	//---------------------------------
	if opt.Parallel {

		var wg sync.WaitGroup
		for _, cluster := range s.configuration.Strategy {
			wg.Add(1)
			go func(options Options, cluster configuration.Cluster) {
				stateMap := process(options, cluster)
				s.printStatemap(stateMap)
				wg.Done()
			}(opt, cluster.Cluster)

		}
		wg.Wait()
	} else {
		for _, cluster := range s.configuration.Strategy {
			stateMap := process(opt, cluster.Cluster)
			s.printStatemap(stateMap)
		}
	}
	return nil
}
