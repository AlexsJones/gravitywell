package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	v1rbac "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1RbacRoleResouce(k kubernetes.Interface, cm *v1rbac.Role, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found Role resource")
	cmclient := k.RbacV1().Roles(namespace)

	if opts.DryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Role resource %s does not exist\n", cm.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Info(fmt.Sprintf("DRY-RUN: Role resource %s exists\n", cm.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := cmclient.Get(cm.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Role resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Role resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Role deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := cmclient.Update(cm)
		if err != nil {
			log.Error("Could not update Role")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Role updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", cm.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", cm.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
