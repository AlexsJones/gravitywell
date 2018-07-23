package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execDeploymentResouce(k kubernetes.Interface, objdep *v1beta1.Deployment, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found deployment resource")

	deploymentClient := k.AppsV1beta1().Deployments(namespace)
	if opts.DryRun {
		_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	if opts.Redeploy || commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		if err := deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod}); err != nil {
			log.Error(err.Error())
		}
	}
	_, err := deploymentClient.Create(objdep)
	if err != nil {
		if opts.TryUpdate || commandFlag == configuration.Apply {
			_, err := deploymentClient.Update(objdep)
			if err != nil {
				log.Error("Deployment could not be updated")
				return state.EDeploymentStateUpdated, err
			}
		}
	}
	log.Debug("Deployment deployed")
	return state.EDeploymentStateOkay, nil
}
