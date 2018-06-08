package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execServiceResouce(k kubernetes.Interface, ss *v1.Service, namespace string, opts configuration.Options) (state.State, error) {
	color.Blue("Found service resource")
	ssclient := k.CoreV1().Services(namespace)

	if opts.DryRun {
		_, err := ssclient.Get(ss.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: Service resource %s does not exist\n", ss.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: Service resource %s exists\n", ss.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	_, err := ssclient.Create(ss)
	if err != nil {
		if !opts.TryUpdate {
			color.Cyan("Service already exists - Cowardly refusing to overwrite")
			return state.EDeploymentStateExists, err
		}
		_, err := ssclient.Update(ss)
		if err != nil {
			color.Red("Could not update service")
			return state.EDeploymentStateCantUpdate, err
		}
		color.Blue("Service updated")
	}
	return state.EDeploymentStateOkay, nil
}
