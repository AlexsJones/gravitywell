package standard

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"github.com/AlexsJones/gravitywell/scheduler/subprocessor"
	"github.com/fatih/color"
	logger "github.com/sirupsen/logrus"
	"os"
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
	coordinator          *subprocessor.Coordinator
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
		coordinator:              subprocessor.NewCoordinator(),
	}
}

func (p *Plan) processApplications(clusterName string, applications []kinds.Application) {
	//Deploy cluster applications
	for _, application := range applications {

		f := func() {
			logger.Info(fmt.Sprintf("Running deployment of %s for cluster %s", application.Name, clusterName))
			actions.ApplicationProcessor(p.commandFlag, p.opt, clusterName, application)
		}
		if p.opt.Parallelism {
			p.coordinator.ResourceChannel <- subprocessor.Resource{
				Process: f}
		} else {
			f()
		}
	}
}
func (p *Plan) clusterFirstDeploymentPlan() {

	for k, _ := range p.providerClusterReference {
		//Cloud provider name
		logger.Info(p.providerClusterReference[k].ProviderName)

		switch strings.ToLower(p.providerClusterReference[k].ProviderName) {

		case "minikube":
			config, err := actions.NewMinikubeConfig()
			if err != nil {
				logger.Fatal(err)
			}
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				err = actions.MinikubeClusterProcessor(config, p.commandFlag, clusters)
				if err != nil {
					logger.Fatal(err)
				}

				if p.commandFlag == configuration.Delete {
					logger.Info("Cluster deleted will not continue")
					os.Exit(0)
				}
				//Deploy cluster applications
				p.processApplications(clusters.Name, p.clusterApplications[clusters.Name])
			}
		case "amazon web services":
			//Configure session
			config, err := actions.NewAmazonWebServicesConfig()
			if err != nil {
				logger.Fatal(err)
			}
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				err = actions.AmazonWebServicesClusterProcessor(config, p.commandFlag, clusters)
				if err != nil {
					logger.Fatal(err)
				}

				if p.commandFlag == configuration.Delete {
					logger.Info("Cluster deleted will not continue")
					os.Exit(0)
				}
				//Deploy cluster applications
				p.processApplications(clusters.Name, p.clusterApplications[clusters.Name])
			}
		case "google cloud platform":
			//Configure session
			config, err := actions.NewGoogleCloudPlatformConfig()
			if err != nil {
				logger.Fatal(err)
			}
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				err := actions.GoogleCloudPlatformClusterProcessor(config, p.commandFlag, clusters)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				if p.commandFlag == configuration.Delete {
					logger.Info("Cluster deleted will not continue")
					os.Exit(0)
				}
				//Deploy cluster applications
				p.processApplications(clusters.Name, p.clusterApplications[clusters.Name])
			}
		case "digital ocean":
			//Configure session
			config, err := actions.NewDigitalOceanConfig()
			if err != nil {
				logger.Fatal(err)
			}
			for _, clusters := range p.providerClusterReference[k].Dependencies {
				//Deploy cluster
				err := actions.DigitalOceanClusterProcessor(config, p.commandFlag, clusters)
				if err != nil {
					logger.Fatal(err)
					os.Exit(1)
				}

				if p.commandFlag == configuration.Delete {
					logger.Info("Cluster deleted will not continue")
					os.Exit(0)
				}
				//Deploy cluster applications
				p.processApplications(clusters.Name, p.clusterApplications[clusters.Name])
			}
		default:
			logger.Warning(fmt.Sprintf("Provider %s unsupported", p.providerClusterReference[k].ProviderName))
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

	// 1. Check whether to deploy clusters
	if p.shouldDeployClusters {
		// 2. deploy cluster first then applications.
		color.Yellow("Running deployment sequence: ClusterFirst")
		p.clusterFirstDeploymentPlan()
	} else {
		// 3. Deploy applications if no cluster has been found
		color.Yellow("Running deployment sequence: ApplicationFirst")
		p.applicationFirstDeploymentPlan()
	}

	p.statusWatcher <- &PlanStatus{ShouldHalt: true}
}
func (p *Plan) Run() planner.IPlanStatusWatcher {

	go p.run()

	if p.opt.Parallelism {
		logger.Info("Enabled parallelism")
		go p.coordinator.Run()
	}

	return p.statusWatcher
}
