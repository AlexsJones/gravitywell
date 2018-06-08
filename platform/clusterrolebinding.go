package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	v1rbac "k8s.io/api/rbac/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execClusterRoleBindingResouce(k kubernetes.Interface, cm *v1rbac.ClusterRoleBinding, namespace string, opts configuration.Options) (state.State, error) {
	color.Blue("Found ClusterRoleBinding resource")
	cmclient := k.RbacV1().ClusterRoleBindings()

	if opts.DryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: ClusterRoleBinding resource %s does not exist\n", cm.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: ClusterRoleBinding resource %s exists\n", cm.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	_, err := cmclient.Create(cm)
	if err != nil {
		if !opts.TryUpdate {
			color.Cyan("ClusterRoleBinding already exists - Cowardly refusing to overwrite")
			return state.EDeploymentStateExists, err
		}
		_, err := cmclient.Update(cm)
		if err != nil {
			color.Red("ClusterRoleBinding could not be updated")
			return state.EDeploymentStateCantUpdate, err
		}
		color.Blue("ClusterRoleBinding updated")
	}
	return state.EDeploymentStateOkay, nil
}
