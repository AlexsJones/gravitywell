package platform

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/_vendor-20181210205236/github.com/fatih/color"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"github.com/jpillora/backoff"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)

func execV1Beta1StatefulSetResouce(k kubernetes.Interface, objdep *v1beta1.StatefulSet, namespace string, opts configuration.Options,
	commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) (state.State, error) {
	log.Debug("Found statefulset resource")
	stsclient := k.AppsV1beta1().StatefulSets(namespace)

	awaitReady := func() error {

		color.Yellow("Awaiting readiness...")
		b := &backoff.Backoff{
			Min:    10 * time.Second,
			Max:    60 * time.Second,
			Jitter: true,
		}
		for {
			stsResponse, err := stsclient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				return errors.New("failed to get deployment")
			}
			if stsResponse.Status.ReadyReplicas >= stsResponse.Status.Replicas {
				return nil
			}
			log.Debug(fmt.Sprintf("Awaiting deployment replica roll out %d/%d",
				stsResponse.Status.ReadyReplicas,
				stsResponse.Status.Replicas))

			time.Sleep(b.Duration())
			if b.Attempt() >= 3 {
				return errors.New("max retry attempts hit")
			}
		}
	}

	if opts.DryRun {
		_, err := stsclient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = stsclient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := stsclient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			log.Debug(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		_, err := stsclient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy objdep resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, err
			}
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := stsclient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy objdep resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, err
			}
		}
		log.Debug("Statefulset deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := stsclient.UpdateStatus(objdep)
		if err != nil {
			log.Error("Could not update Statefulset")
			return state.EDeploymentStateCantUpdate, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, err
			}
		}
		log.Debug("Statefulset updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := stsclient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
