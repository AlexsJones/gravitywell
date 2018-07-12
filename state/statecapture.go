package state

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

func NewCapture() *Capture {

	return &Capture{
		DeploymentState: make(map[string]Details),
	}
}
