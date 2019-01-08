package state

import (
	"fmt"
	"github.com/fatih/color"
)

type Details struct {
	State     State
	HasDetail bool
	HasError  bool
	Detail    string
}

type Capture struct {
	ClusterName     string
	DeploymentState map[string]Details
}

var col func(string, ...interface{})

func (c *Capture) Print() {
	for k, v := range c.DeploymentState {
		col = color.Green
		if v.State == EDeploymentStateError {
			col = color.Red
		}
		if v.State == EDeploymentStateNotExists {
			col = color.Red
		}
		col(fmt.Sprintf("Cluster %s Deployment %s State => %s\n", c.ClusterName, k, Translate(v.State)))
		if v.HasDetail && v.HasError {
			color.Cyan(fmt.Sprintf("\t %s\n", v.Detail))
		}
	}
}
