package standard

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"log"
)

type StandardPlanner struct {
	plan *Plan
}

func (s StandardPlanner) GeneratePlan(configuration *configuration.Configuration) (planner.IPlan, error) {
	s.plan = NewPlan()

	//Clusters -----------------------------------------------------------------------------------
	//Generate a map e.g.
	// Google Cloud platform:< Cluster Data >
	type ProviderClusterReference struct {
		ProviderName string
		Dependencies []kinds.IKind
	}

	clusterSignPosts := make(map[string]*ProviderClusterReference)
	clusterDeployments := make(map[string][]kinds.IKind)

	for _, clusterKinds := range configuration.ClusterKinds {

		for _, providers := range clusterKinds.Strategy {

			pcr := ProviderClusterReference{ProviderName: providers.Provider.Name}

			clusterSignPosts[providers.Provider.Name] = &pcr

			for _, clusters := range providers.Provider.Clusters {
				clusterDeployments[pcr.ProviderName] = append(clusterDeployments[pcr.ProviderName], clusters)
			}

		}
	}
	log.Printf(fmt.Sprintf("%+v", clusterDeployments))
	// -------------------------------------------------------------------------------------------

	// Application -------------------------------------------------------------------------------

	// -------------------------------------------------------------------------------------------
	return s.plan, nil
}
