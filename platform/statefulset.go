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

func execStatefulSetResouce(k kubernetes.Interface, sts *v1beta1.StatefulSet, namespace string, opts configuration.Options) (state.State, error) {
	color.Blue("Found statefulset resource")
	stsclient := k.AppsV1beta1().StatefulSets(namespace)

	if opts.DryRun {
		_, err := stsclient.Get(sts.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", sts.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", sts.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	if opts.Redeploy {
		color.Blue("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		stsclient.Delete(sts.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
	}
	_, err := stsclient.Create(sts)
	if err != nil {
		if opts.TryUpdate {
			_, err := stsclient.UpdateStatus(sts)
			if err != nil {
				color.Red("Could not update Statefulset")
				return state.EDeploymentStateCantUpdate, err
			}
			color.Blue("Statefulset updated")
			return state.EDeploymentStateUpdated, nil
		}
	}
	color.Blue("Statefulset deployed")
	return state.EDeploymentStateOkay, nil
}
