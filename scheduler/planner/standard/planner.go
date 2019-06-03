package standard

import (
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"github.com/AlexsJones/gravitywell/shared"
	logger "github.com/sirupsen/logrus"
)

type StandardPlanner struct {
	plan *Plan
}

func (s StandardPlanner) GeneratePlan(configuration *configuration.Configuration,
	flag configuration.CommandFlag, opt configuration.Options) (planner.IPlan, error) {

	s.plan = NewPlan(flag, opt)

	//Clusters -----------------------------------------------------------------------------------

	for _, clusterKinds := range configuration.ClusterKinds {

		for _, providers := range clusterKinds.Strategy {

			pcr := ProviderClusterReference{ProviderName: providers.Provider.Name}

			s.plan.providerClusterReference[providers.Provider.Name] = &pcr

			for _, clusters := range providers.Provider.Clusters {
				s.plan.clusterDeployments[pcr.ProviderName] = append(s.plan.clusterDeployments[pcr.ProviderName], clusters.Cluster)

				pcr.Dependencies = append(pcr.Dependencies, clusters.Cluster)
			}

		}
	}
	logger.Infof(shared.PrettyPrint(s.plan.clusterDeployments))
	//At this point if there are no clusters, we set a flag to tell the plan to only run applications
	if len(s.plan.clusterDeployments) == 0 {
		logger.Info("No clusters found to deploy - skipping")
		s.plan.shouldDeployClusters = false
	} else {
		logger.Info("Clusters found to deploy - sequencing")
		s.plan.shouldDeployClusters = true
	}

	// -------------------------------------------------------------------------------------------

	// Application -------------------------------------------------------------------------------
	for _, applicationKinds := range configuration.ApplicationKinds {
		for _, applicationLists := range applicationKinds.Strategy {

			for _, app := range applicationLists.Cluster.Applications {
				s.plan.clusterApplications[applicationLists.Cluster.Name] =
					append(s.plan.clusterApplications[applicationLists.Cluster.Name], app.Application)
			}
		}
	}
	// -------------------------------------------------------------------------------------------
	logger.Infof(shared.PrettyPrint(s.plan.clusterApplications))
	// -------------------------------------------------------------------------------------------
	return s.plan, nil
}
