package platform

import (
	"fmt"

	"github.com/fatih/color"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execConfigMapResouce(k kubernetes.Interface, cm *v1.ConfigMap, namespace string, dryRun bool, tryUpdate bool) error {
	color.Blue("Found Configmap resource")
	cmclient := k.CoreV1().ConfigMaps(namespace)

	if dryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: Configmap resource %s does not exist\n", cm.Name))
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: Configmap resource %s exists\n", cm.Name))
		}
		return err
	}

	_, err := cmclient.Create(cm)
	if err != nil {
		if !tryUpdate {
			color.Cyan("Configmap already exists - Cowardly refusing to overwrite")
			return err
		}
		_, err := cmclient.Update(cm)
		if err != nil {
			color.Red("Configmap could not be updated")
			return err
		}
		color.Blue("Configmap updated")
	}
	return nil
}
