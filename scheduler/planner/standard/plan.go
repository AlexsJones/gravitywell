package standard

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/_vendor-20181210205236/github.com/fatih/color"
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
	opt                      configuration.Options

	//Sequence control
	shouldDeployClusters bool
}

type PlanStatus struct {
	ShouldHalt bool
}

func (p *PlanStatus) Halt() bool {
	return p.ShouldHalt
}

func NewPlan(flag configuration.CommandFlag, opt configuration.Options) *Plan {
	return &Plan{
		statusWatcher:            make(planner.IPlanStatusWatcher),
		providerClusterReference: make(map[string]*ProviderClusterReference),
		clusterDeployments:       make(map[string][]kinds.ProviderCluster),
		clusterApplications:      make(map[string][]kinds.Application),
		commandFlag:              flag,
		opt:                      opt,
	}
}

func (p *Plan) clusterFirstDeploymentPlan() {

	for k, _ := range p.providerClusterReference {
		//Cloud provider name
		log.Println(p.providerClusterReference[k].ProviderName)

		switch strings.ToLower(p.providerClusterReference[k].ProviderName) {
		case "google cloud platform":
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				actions.GoogleCloudClusterProcessor(p.commandFlag, clusters)

				//Deploy cluster applications
				for _, application := range p.clusterApplications[clusters.FullName] {
					color.Yellow(fmt.Sprintf("Running deployment of %s for cluster %s", application.Name, clusters.FullName))
					actions.ApplicationProcessor(p.commandFlag, p.opt, clusters.FullName, application)
				}
			}
		default:
			log.Fatal(fmt.Sprintf("Provider %s unsupported", p.providerClusterReference[k].ProviderName))
			p.statusWatcher <- &PlanStatus{ShouldHalt: true}
		}
	}
}
func (p *Plan) applicationFirstDeploymentPlan() {

	for clusterFullName, _ := range p.clusterApplications {

		for _, application := range p.clusterApplications[clusterFullName] {

			color.Yellow(fmt.Sprintf("Running deployment of %s for cluster %s", application.Name, clusterFullName))
			actions.ApplicationProcessor(p.commandFlag, p.opt, clusterFullName, application)
		}
	}
}
func (p *Plan) run() {

	//0. Check whether to deploy clusters
	if p.shouldDeployClusters {
		//1. deploy cluster first then applications.
		color.Yellow("Running deployment sequence: ClusterFirst")
		p.clusterFirstDeploymentPlan()
	} else {
		//2. Deploy applications if no cluster has been found
		color.Yellow("Running deployment sequence: ApplicationFirst")
		p.applicationFirstDeploymentPlan()
	}
	p.statusWatcher <- &PlanStatus{ShouldHalt: true}
}
func (p *Plan) Run() planner.IPlanStatusWatcher {

	go p.run()

	return p.statusWatcher
}
