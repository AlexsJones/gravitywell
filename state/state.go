package state

type State int

const (
	EDeploymentStateNil State = iota
	EDeploymentStateOkay
	EDeploymentStateError
	EDeploymentStateExists
	EDeploymentStateNotExists
	EDeploymentStateCantUpdate
	EDeploymentStatDone
)
