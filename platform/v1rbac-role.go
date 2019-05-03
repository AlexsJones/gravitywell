package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/google/logger"
	v1rbac "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1RbacRoleResouce(k kubernetes.Interface, objdep *v1rbac.Role, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	logger.Info("Found Role resource")
	cmclient := k.RbacV1().Roles(namespace)

	if opts.DryRun {
		_, err := cmclient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("DRY-RUN: Role resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			logger.Info(fmt.Sprintf("DRY-RUN: Role resource %s exists\n", objdep.Name))
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
			logger.Error(fmt.Sprintf("Could not deploy Role resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := cmclient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Role resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info("Role deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := cmclient.Update(objdep)
		if err != nil {
			logger.Error("Could not update Role")
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info("Role updated")
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
