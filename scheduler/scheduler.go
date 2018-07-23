package scheduler

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
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
func (s *Scheduler) Run(opt configuration.Options) error {

	if opt.APIVersion != s.configuration.APIVersion {
		log.Error(fmt.Sprintf("Manifest is not supported by the current API: %s\n", opt.APIVersion))
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
	var allClusterStates []*state.Capture

	for _, cluster := range s.configuration.Strategy {
		for _, ignoreItem := range opt.IgnoreList {
			if ignoreItem == cluster.Cluster.Name {
				log.Warn(fmt.Sprintf("Ignoring cluster %s\n", cluster.Cluster.Name))
				continue
			}
		}
		stateMap := process(opt, cluster.Cluster)
		allClusterStates = append(allClusterStates, stateMap)
	}
	//----------------------------------
	var col func(string, ...interface{})
	for _, stateCapture := range allClusterStates {
		for k, v := range stateCapture.DeploymentState {
			col = color.Green
			if v.State == state.EDeploymentStateError {
				col = color.Red
			}
			if v.State == state.EDeploymentStateNotExists {
				col = color.Red
			}
			col(fmt.Sprintf("Cluster %s Deployment %s State => %s\n", stateCapture.ClusterName, k, state.Translate(v.State)))
			if v.HasDetail && v.HasError {
				color.Cyan(fmt.Sprintf("\t %s\n", v.Detail))
			}
		}
	}
	return nil
}
