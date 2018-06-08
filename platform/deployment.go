package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execDeploymentResouce(k kubernetes.Interface, objdep *v1beta1.Deployment, namespace string, opts configuration.Options) (state.State, error) {
	color.Blue("Found deployment resource")

	deploymentClient := k.AppsV1beta1().Deployments(namespace)
	if opts.DryRun {
		_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	if opts.Redeploy {
		color.Blue("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
	}
	_, err := deploymentClient.Create(objdep)
	if err != nil {
		if opts.TryUpdate {
			_, err := deploymentClient.Update(objdep)
			if err != nil {
				color.Red("Deployment could not be updated")
				return state.EDeploymentStateUpdated, err
			}
		}
	}
	color.Blue("Deployment deployed")
	return state.EDeploymentStateOkay, nil
}
