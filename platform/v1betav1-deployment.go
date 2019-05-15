package platform

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	logger "github.com/sirupsen/logrus"
	"github.com/jpillora/backoff"
	v1betav1 "k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)

func execV1Betav1DeploymentResouce(k kubernetes.Interface, objdep *v1betav1.Deployment,
	namespace string, opts configuration.Options,
	commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) (state.State, error) {
	logger.Info("Found deployment resource")

	deploymentClient := k.ExtensionsV1beta1().Deployments(namespace)

	awaitReady := func() error {
		color.Yellow("Awaiting readiness...")
		b := &backoff.Backoff{
			Min:    10 * time.Second,
			Max:    opts.MaxBackOffDuration,
			Jitter: true,
		}
		for {
			stsResponse, err := deploymentClient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				return errors.New("failed to get deployment")
			}
			if stsResponse.Status.ReadyReplicas >= stsResponse.Status.Replicas {
				return nil
			}
			logger.Info(fmt.Sprintf("Awaiting deployment replica roll out %d/%d",
				stsResponse.Status.ReadyReplicas,
				stsResponse.Status.Replicas))

			time.Sleep(b.Duration())
			if b.Attempt() >= 3 {
				return errors.New("max retry attempts hit")
			}
		}
	}

	if opts.DryRun {
		_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			logger.Info(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		logger.Info("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := deploymentClient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			logger.Info(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		_, err := deploymentClient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Deployment resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		logger.Info("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := deploymentClient.Create(objdep)
		if err != nil {
			logger.Error(fmt.Sprintf("Could not deploy Deployment resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		logger.Info("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := deploymentClient.Update(objdep)
		if err != nil {
			logger.Error("Could not update Deployment")
			return state.EDeploymentStateCantUpdate, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		logger.Info("Deployment updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			logger.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		logger.Info(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
