package standard

import "github.com/AlexsJones/gravitywell/scheduler/planner"

type Plan struct {
	statusWatcher planner.IPlanStatusWatcher
}

type PlanStatus struct {
	ShouldHalt bool
}

func (p *PlanStatus) Halt() bool {
	return p.ShouldHalt
}

func NewPlan() *Plan {
	return &Plan{
		statusWatcher: make(planner.IPlanStatusWatcher),
	}
}

func (p *Plan) run() {
	p.statusWatcher <- &PlanStatus{ShouldHalt: true}
}
func (p *Plan) Run() planner.IPlanStatusWatcher {

	go p.run()

	return p.statusWatcher
}
