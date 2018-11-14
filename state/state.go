package state

type State int

const (
	EDeploymentStateNil State = iota
	EDeploymentStateOkay
	EDeploymentStateError
	EDeploymentStateExists
	EDeploymentStateNotExists
	EDeploymentStateCantUpdate
	EDeploymentStateDone
	EDeploymentStateUpdated
	EDeploymentStatePatched
)

//Translate enum to string
func Translate(i State) string {
	switch i {
	case EDeploymentStateNil:
		return "Nil"
	case EDeploymentStateOkay:
		return "Ok"
	case EDeploymentStateError:
		return "Error"
	case EDeploymentStateExists:
		return "Exists"
	case EDeploymentStateNotExists:
		return "Doesn't exist"
	case EDeploymentStateCantUpdate:
		return "Immutable/Can't update"
	case EDeploymentStateDone:
		return "Done"
	case EDeploymentStateUpdated:
		return "Updated"
	case EDeploymentStatePatched:
		return "Patched"
	}
	return "N/A"
}
