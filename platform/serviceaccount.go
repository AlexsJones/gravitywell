package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execServiceAccountResouce(k kubernetes.Interface, cm *v1.ServiceAccount, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	color.Blue("Found Configmap resource")
	cmclient := k.CoreV1().ServiceAccounts(namespace)

	if opts.DryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: ServiceAccount resource %s does not exist\n", cm.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: ServiceAccount resource %s exists\n", cm.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	if opts.Redeploy || commandFlag == configuration.Replace {
		color.Blue("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
	}
	_, err := cmclient.Create(cm)
	if err != nil {
		if opts.TryUpdate || commandFlag == configuration.Apply {
			_, err := cmclient.Update(cm)
			if err != nil {
				color.Red("ServiceAccount could not be updated")
				return state.EDeploymentStateCantUpdate, err
			}
			color.Blue("ServiceAccount updated")
			return state.EDeploymentStateUpdated, nil
		}
	}
	color.Blue("ServiceAccount deployed")
	return state.EDeploymentStateOkay, nil
}
