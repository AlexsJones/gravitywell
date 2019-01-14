package platform

import (
	"errors"
	"fmt"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Job(k kubernetes.Interface, job *batchv1.Job, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found Job resource")
	client := k.BatchV1().Jobs(namespace)

	if opts.DryRun {
		_, err := client.Get(job.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: Job resource %s does not exist\n", job.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Debug(fmt.Sprintf("DRY-RUN: Job resource %s exists\n", job.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		client.Delete(job.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		_, err := client.Create(job)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Job resource %s due to %s", job.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Job deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := client.Create(job)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy Job resource %s due to %s", job.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Job deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := client.UpdateStatus(job)
		if err != nil {
			log.Error("Could not update Job")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("Job updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := client.Delete(job.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", job.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", job.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")

}
