package platform

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"github.com/jpillora/backoff"
	logger "github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)

func execV1Beta1StatefulSetResouce(k kubernetes.Interface, objdep *v1beta1.StatefulSet, namespace string,
	opts configuration.Options, commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) (state.State, error) {
	name := "StatefulSet"

	client := k.AppsV1beta1().StatefulSets(namespace)

	exists := false
	_, err := client.Get(objdep.Name, meta_v1.GetOptions{})
	if err == nil {
		exists = true
	}

	if opts.DryRun {
		if exists == false {
			logger.Error(fmt.Sprintf("DRY-RUN: %s resource %s does not exist\n", name, objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			logger.Info(fmt.Sprintf("DRY-RUN: %s resource %s exists\n", name, objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	// ----------------------------------------------------------------------------------------------------------------
	awaitReady := func() error {

		color.Yellow("Awaiting readiness...")
		b := &backoff.Backoff{
			Min:    10 * time.Second,
			Max:    opts.MaxBackOffDuration,
			Jitter: true,
		}
		for {
			stsResponse, err := client.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				return errors.New("failed to get deployment")
			}
			if stsResponse.Status.ReadyReplicas >= stsResponse.Status.CurrentReplicas {
				return nil
			}
			logger.Info(fmt.Sprintf("Awaiting deployment replica roll out %d/%d",
				stsResponse.Status.ReadyReplicas,
				stsResponse.Status.CurrentReplicas))

			time.Sleep(b.Duration())
			if b.Attempt() >= 3 {
				return errors.New("max retry attempts hit")
			}
		}
	}
	create := func() (state.State, error) {
		if exists {
			return state.EDeploymentStateExists, nil
		}
		_, err := client.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy %s resource %s due to %s", name, objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		logger.Info(fmt.Sprintf("%s deployed", name))
		return state.EDeploymentStateOkay, nil
	}
	update := func() (state.State, error) {
		if !exists {
			return create()
		}
		_, err := client.Update(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not update %s", name))
			return state.EDeploymentStateCantUpdate, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		logger.Info(fmt.Sprintf("%s updated", name))
		return state.EDeploymentStateUpdated, nil
	}
	del := func() (state.State, error) {
		if !exists {
			return state.EDeploymentStateDone, nil
		}
		logger.Info("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		err := client.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		if err != nil {
			return state.EDeploymentStateNotExists, err
		}
		for {
			_, err := client.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			logger.Info(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		return state.EDeploymentStateDone, nil
	}
	// ----------------------------------------------------------------------------------------------------------------

	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {

		return create()
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {

		if !exists {
			return create()
		} else {
			return update()
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		if exists {
			if _, err := del(); err != nil {
				return state.EDeploymentStateError, err
			}
		}
		_, err = client.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy %s resource %s due to %s", name, objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		logger.Info(fmt.Sprintf("%s deployed", name))
		return state.EDeploymentStateOkay, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := client.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New(fmt.Sprintf("no kubectl command given to %s", name))
}
