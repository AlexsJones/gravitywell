package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execStatefulSetResouce(k kubernetes.Interface, sts *v1beta1.StatefulSet, namespace string, dryRun bool, tryUpdate bool) (state.State, error) {
	color.Blue("Found statefulset resource")
	stsclient := k.AppsV1beta1().StatefulSets(namespace)

	if dryRun {
		_, err := stsclient.Get(sts.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", sts.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", sts.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	_, err := stsclient.Create(sts)
	if err != nil {
		if !tryUpdate {
			color.Cyan("StatefulSet already exists - Cowardly refusing to overwrite")
			return state.EDeploymentStateExists, err
		}
		_, err := stsclient.UpdateStatus(sts)
		if err != nil {
			color.Red("Could not update Statefulset")
			return state.EDeploymentStateCantUpdate, err
		}
		color.Blue("Statefulset updated")
	}
	return state.EDeploymentStateOkay, nil
}
