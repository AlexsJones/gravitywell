package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Beta1IngressResouce(k kubernetes.Interface, ingress *v1beta1.Ingress, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found ingress resource")
	dsclient := k.Extensions().Ingresses(namespace)

	if opts.DryRun {
		_, err := dsclient.Get(ingress.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Ingress resource %s does not exist\n", ingress.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Ingress resource %s exists\n", ingress.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = dsclient.Delete(ingress.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := dsclient.Get(ingress.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_, err := dsclient.Create(ingress)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ingress resource %s due to %s", ingress.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Ingress deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := dsclient.Create(ingress)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ingress resource %s due to %s", ingress.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Ingress deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := dsclient.UpdateStatus(ingress)
		if err != nil {
			log.Error("Could not update Ingress")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Ingress updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := dsclient.Delete(ingress.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", ingress.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", ingress.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")

}
