package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/core/v1"
	meta1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1SecretResouce(k kubernetes.Interface, secret *v1.Secret, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found Secret resource")
	client := k.CoreV1().Secrets(namespace)

	if opts.DryRun {
		_, err := client.Get(secret.Name, meta1.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Secret resource %s does not exist\n", secret.Name))
			return state.EDeploymentStateNotExists, err
		}
		log.Info(fmt.Sprintf("DRY-RUN: Secret resource %s exists\n", secret.Name))
		return state.EDeploymentStateExists, nil
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = client.Delete(secret.Name, &meta1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := client.Get(secret.Name, meta1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_, err := client.Create(secret)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Secret resource %s due to %s", secret.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Secret deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := client.Create(secret)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Secret resource %s due to %s", secret.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Secret deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := client.Update(secret)
		if err != nil {
			log.Error("Could not update Secret")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Secret updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := client.Delete(secret.Name, &meta1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", secret.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", secret.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
