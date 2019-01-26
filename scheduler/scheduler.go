package scheduler

import (
	"errors"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"github.com/AlexsJones/gravitywell/scheduler/planner/standard"
	"log"
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

	stdplnr := standard.StandardPlanner{}

	plan, err := planner.GeneratePlan(stdplnr, s.configuration, commandFlag)
	if err != nil {
		log.Fatal(err)
	}

	statusWatcher := plan.Run()
	for {
		select {
		case msg := <-statusWatcher:
			if msg.Halt() {
				//Halting
				log.Fatal("Received halt")
			}

		}
	}

	//composedRunBook := make(map[configuration.Provider][]configuration.ProviderCluster)
	//
	////1. For manifests with clusters & applications
	//for _, ck := range s.configuration.ClusterKinds {
	//	for _, clusterKind := range ck.Strategy {
	//		//A provider e.g. GCP/AWS/Azure
	//		for _, cluster := range clusterKind.Provider.Clusters {
	//			//A cluster
	//
	//			//At this point we want to load applications into a cluster
	//			// cluster.Cluster.Applications
	//			//This operation requires us to recurse through all applications and map them back
	//			//wew hard work!
	//			if err := s.mapApplicationsToCluster(&cluster.Cluster, opt); err != nil {
	//				color.Yellow(err.Error())
	//				continue
	//			}
	//			//Let's add this cluster to our overall map of how we think about clusters
	//			// This means all the GCP clusters are keyed under the GCP provider
	//			composedRunBook[clusterKind.Provider] = append(composedRunBook[clusterKind.Provider], cluster.Cluster)
	//		}
	//	}
	//}
	//
	//if len(s.configuration.ClusterKinds) == 0 {
	//	//2. If there are no clusters then we need to infer they already exist within the kubecontext
	//
	//} else {
	//	//3. Alternatively we run the composed deployment plan
	//}
	//configuredDeployments := make(map[configuration.ClusterKind][]configuration.Application)
	//
	//
	//
	//
	//
	//for _, applicationKind := range s.configuration.ApplicationKinds {
	//
	//	if opt.APIVersion != applicationKind.APIVersion {
	//		log.Error(fmt.Sprintf("Manifest is not supported by the current API: %s\n", opt.APIVersion))
	//		continue
	//	}
	//	//---------------------------------
	//	for _, cluster := range applicationKind.Strategy {
	//
	//		//The required clusterName
	//		providerCluster, err := findProviderCluster(cluster.Cluster.Name)
	//		if err != nil {
	//			color.Red(fmt.Sprintf("Applications found without a valid correlating: cluster %s was not found.",cluster.Cluster.Name)
	//			log.Fatal(err)
	//		}
	//		for _, deployment := range cluster.Cluster.Applications {
	//
	//			//Having found an application we now need to correlate it to a cluster
	//			configuredDeployments[providerCluster] =
	//				append(configuredDeployments[providerCluster],
	//				deployment.Application)
	//		}
	//	}
	//
	//	//---------------------------------
	//
	//	//for _, cluster := range applicationKind.Strategy {
	//	//
	//	//	ApplicationProcessor(commandFlag, opt, cluster.Cluster)
	//	//}
	//	//----------------------------------
	//
	//}
	return nil
}
