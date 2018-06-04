package platform

import (
	"fmt"

	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execDeploymentResouce(k kubernetes.Interface, objdep *v1beta1.Deployment, namespace string, dryRun bool, tryUpdate bool) error {
	color.Blue("Found deployment resource")
	deploymentClient := k.AppsV1beta1().Deployments(namespace)
	if dryRun {
		_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
		}
		return err
	}

	_, err := deploymentClient.Create(objdep)
	if err != nil {
		if !tryUpdate {
			color.Cyan("Deployment already exists - Cowardly refusing to overwrite")
			return err
		}
		_, err := deploymentClient.Update(objdep)
		if err != nil {
			color.Red("Deployment could not be updated")
			return err
		}
	}
	return nil
}
