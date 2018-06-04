package platform

import (
	"fmt"

	"github.com/fatih/color"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execServiceResouce(k kubernetes.Interface, ss *v1.Service, namespace string, dryRun bool, tryUpdate bool) error {
	color.Blue("Found service resource")
	ssclient := k.CoreV1().Services(namespace)

	if dryRun {
		_, err := ssclient.Get(ss.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: Service resource %s does not exist\n", ss.Name))
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: Service resource %s exists\n", ss.Name))
		}
		return err
	}

	_, err := ssclient.Create(ss)
	if err != nil {
		if !tryUpdate {
			color.Cyan("Service already exists - Cowardly refusing to overwrite")
			return err
		}
		_, err := ssclient.Update(ss)
		if err != nil {
			color.Red("Could not update service")
			return err
		}
		color.Blue("Service updated")
	}
	return nil
}
