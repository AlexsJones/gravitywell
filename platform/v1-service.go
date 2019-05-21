package platform

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	logger "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)
func execV1ServiceResouce(k kubernetes.Interface, objdep *v1.Service, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	logger.Info("Found service resource")
	ssclient := k.CoreV1().Services(namespace)

	exists := false
	_, err := ssclient.Get(objdep.Name, v12.GetOptions{})
	if err == nil {
		exists = true
	}

	if opts.DryRun {
		if exists == false {
			logger.Error(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			logger.Info(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	update := func() (state.State,error) {
		_, err := ssclient.Update(objdep)
		if err != nil {
			logger.Error("Could not update Service")
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info("Service updated")
		return state.EDeploymentStateUpdated, nil
	}
	del := func() (state.State,error) {
		if !exists {
			return state.EDeploymentStateDone,nil
		}
		logger.Info("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		err := ssclient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		if err != nil {
			return state.EDeploymentStateNotExists, err
		}
		for {
			_, err := ssclient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			logger.Info(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		return state.EDeploymentStateDone,nil
	}

	create := func() (state.State, error){
		if exists {
			return state.EDeploymentStateExists,nil
		}
		_, err := ssclient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Service deployed")
		return state.EDeploymentStateOkay, nil
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		if exists {
			if _,err := del(); err != nil {
				return state.EDeploymentStateError,err
			}
		}
		_, err = ssclient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Service resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Service deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {

		return create()
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {

		if opts.Force {
			if !exists {
				return create()
			}else {
			if _,err := del(); err != nil {
				return state.EDeploymentStateError,err
			}
			return update()
			}
		} else {
			return update()
		}
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := ssclient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
