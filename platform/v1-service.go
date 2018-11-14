package platform

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1ServiceResouce(k kubernetes.Interface, ss *v1.Service, application configuration.Application,
	executionStep configuration.Execute,
	opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found service resource")
	ssclient := k.CoreV1().Services(application.Namespace)

	if opts.DryRun {
		_, err := ssclient.Get(ss.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Service resource %s does not exist\n", ss.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Service resource %s exists\n", ss.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		ssclient.Delete(ss.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := ssclient.Create(ss)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", ss.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Service deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := ssclient.Create(ss)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", ss.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Service deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := ssclient.Update(ss)
		if err != nil {
			log.Error("Could not update Service")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Service updated")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
