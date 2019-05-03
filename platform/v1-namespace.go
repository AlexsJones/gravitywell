package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"github.com/google/logger"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1NamespaceResource(k kubernetes.Interface, objdep *v1.Namespace, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	color.Blue("Found Namespace resource")
	cmclient := k.CoreV1().Namespaces()

	if opts.DryRun {
		_, err := cmclient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("DRY-RUN: Namespace resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			logger.Info(fmt.Sprintf("DRY-RUN: Namespace resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		logger.Info("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = cmclient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := cmclient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			logger.Info(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		_, err := cmclient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Namespace resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create -------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := cmclient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Namespace resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Namespace deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply -------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := cmclient.Update(objdep)
		if err != nil {
			logger.Error("Could not update Namespace")
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info("Namespace updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := cmclient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
