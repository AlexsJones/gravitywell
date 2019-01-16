package platform

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1StatefulSetResouce(k kubernetes.Interface, sts *appsv1.StatefulSet, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found statefulset resource")
	stsclient := k.AppsV1().StatefulSets(namespace)

	if opts.DryRun {
		_, err := stsclient.Get(sts.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", sts.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", sts.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		stsclient.Delete(sts.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := stsclient.Create(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy sts resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := stsclient.Create(sts)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy sts resource %s due to %s", sts.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := stsclient.UpdateStatus(sts)
		if err != nil {
			log.Error("Could not update Statefulset")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Statefulset updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := stsclient.Delete(sts.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", sts.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", sts.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
