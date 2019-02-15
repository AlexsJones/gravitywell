package platform

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"github.com/jpillora/backoff"
	"k8s.io/api/apps/v1beta2"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"time"
)

func execV1Beta2DeploymentResouce(k kubernetes.Interface, objdep *v1beta2.Deployment,
	namespace string, opts configuration.Options,
	commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) (state.State, error) {
	log.Debug("Found deployment resource")

	deploymentClient := k.AppsV1beta2().Deployments(namespace)

	awaitReady := func() error {

		b := &backoff.Backoff{
			Min:    10 * time.Second,
			Max:    60 * time.Second,
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
		_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := deploymentClient.Get(objdep.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
			log.Debug(fmt.Sprintf("Awaiting deletion of %s", objdep.Name))
		}
		_, err := deploymentClient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Deployment resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		log.Debug("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := deploymentClient.Create(objdep)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Deployment resource %s due to %s", objdep.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		log.Debug("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := deploymentClient.Update(objdep)
		if err != nil {
			log.Error("Could not update Deployment")
			return state.EDeploymentStateCantUpdate, err
		}
		if shouldAwaitDeployment {
			if err := awaitReady(); err != nil {
				return state.EDeploymentStateError, nil
			}
		}
		log.Debug("Deployment updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := deploymentClient.Delete(objdep.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", objdep.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", objdep.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
