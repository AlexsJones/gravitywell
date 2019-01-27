package planner

import "github.com/AlexsJones/gravitywell/configuration"

//IPlanStatus ------------------------------------------------------------------------
type IPlanStatus interface {
	Halt() bool
}

func Halt(i IPlanStatus) bool {
	return i.Halt()
}

type IPlanStatusWatcher chan IPlanStatus

//IPlan ------------------------------------------------------------------------------
type IPlan interface {
	Run() IPlanStatusWatcher
}

func Run(i IPlan) {
	i.Run()
}

//IPlanner ---------------------------------------------------------------------------
type IPlanner interface {
	GeneratePlan(configuration *configuration.Configuration, flag configuration.CommandFlag, opt configuration.Options) (IPlan, error)
}

func GeneratePlan(i IPlanner, configuration *configuration.Configuration, flag configuration.CommandFlag, opt configuration.Options) (IPlan, error) {

	return i.GeneratePlan(configuration, flag, opt)
}
