package scheduler

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
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
func (s *Scheduler) Run(commandFlag configuration.CommandFlag,
	opt configuration.Options) error {

	//Cluster ...
	for _, ck := range s.configuration.ClusterKinds {
		for _, clusterKind := range ck.Strategy {
			ClusterProcessor(commandFlag, clusterKind.Provider)
		}
	}
	//Application ...
	for _, applicationKind := range s.configuration.ApplicationKinds {

		if opt.APIVersion != applicationKind.APIVersion {
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

		for _, cluster := range applicationKind.Strategy {
			stateMap := ApplicationProcessor(commandFlag, opt, cluster.Cluster)
			allClusterStates = append(allClusterStates, stateMap)
		}
		//----------------------------------
		for _, stateCapture := range allClusterStates {

			stateCapture.Print()
		}
	}
	return nil
}
