package standard

import (
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
)

type StandardPlanner struct {
	plan *Plan
}

func (s StandardPlanner) GeneratePlan(configuration *configuration.Configuration) (planner.IPlan, error) {
	s.plan = NewPlan()

	for _, ck := range configuration.ClusterKinds {
		for _, clusterKind := range ck.Strategy {
			//A provider e.g. GCP/AWS/Azure
			for _, cluster := range clusterKind.Provider.Clusters {
				//A cluster

				//At this point we want to load applications into a cluster
				// cluster.Cluster.Applications
				//This operation requires us to recurse through all applications and map them back
				//wew hard work!
				//if err := s.mapApplicationsToCluster(&cluster.Cluster, opt); err != nil {
				//	color.Yellow(err.Error())
				//	continue
				//}
				//Let's add this cluster to our overall map of how we think about clusters
				// This means all the GCP clusters are keyed under the GCP provider
				//ster)
			}
		}
	}
	return s.plan, nil
}
