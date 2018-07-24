package platform

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	v1rbac "k8s.io/api/rbac/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execClusterRoleBindingResouce(k kubernetes.Interface, cm *v1rbac.ClusterRoleBinding, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found ClusterRoleBinding resource")
	cmclient := k.RbacV1().ClusterRoleBindings()

	if opts.DryRun {
		_, err := cmclient.Get(cm.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: ClusterRoleBinding resource %s does not exist\n", cm.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Info(fmt.Sprintf("DRY-RUN: ClusterRoleBinding resource %s exists\n", cm.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		cmclient.Delete(cm.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ClusterRoleBinding resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := cmclient.Create(cm)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy ClusterRoleBinding resource %s due to %s", cm.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("ClusterRoleBinding deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := cmclient.Update(cm)
		if err != nil {
			log.Error("Could not update ClusterRoleBinding")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("ClusterRoleBinding updated")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
