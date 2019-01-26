package standard

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"log"
	"strings"
)

//Generate a map e.g.
// Google Cloud platform:< Cluster Data >
type ProviderClusterReference struct {
	ProviderName string
	Dependencies []kinds.ProviderCluster
}

type Plan struct {
	statusWatcher            planner.IPlanStatusWatcher
	providerClusterReference map[string]*ProviderClusterReference
	clusterDeployments       map[string][]kinds.ProviderCluster
	clusterApplications      map[string][]kinds.Application
	commandFlag              configuration.CommandFlag
}

type PlanStatus struct {
	ShouldHalt bool
}

func (p *PlanStatus) Halt() bool {
	return p.ShouldHalt
}

func NewPlan(flag configuration.CommandFlag) *Plan {
	return &Plan{
		statusWatcher:            make(planner.IPlanStatusWatcher),
		providerClusterReference: make(map[string]*ProviderClusterReference),
		clusterDeployments:       make(map[string][]kinds.ProviderCluster),
		clusterApplications:      make(map[string][]kinds.Application),
		commandFlag:              flag,
	}
}

func (p *Plan) run() {
	//In the beginning ...
	for k, _ := range p.providerClusterReference {
		//Cloud provider name
		log.Println(p.providerClusterReference[k].ProviderName)

		switch strings.ToLower(p.providerClusterReference[k].ProviderName) {
		case "google cloud platform":
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				actions.GoogleCloudClusterProcessor(p.commandFlag, clusters)
			}
		default:
			log.Fatal(fmt.Sprintf("Provider %s unsupported", p.providerClusterReference[k].ProviderName))
			p.statusWatcher <- &PlanStatus{ShouldHalt: true}
		}
	}

	p.statusWatcher <- &PlanStatus{ShouldHalt: true}
}
func (p *Plan) Run() planner.IPlanStatusWatcher {

	go p.run()

	return p.statusWatcher
}
