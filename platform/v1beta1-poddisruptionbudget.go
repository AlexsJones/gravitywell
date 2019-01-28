package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	v1polbeta "k8s.io/api/policy/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Beta1PodDisruptionBudgetResouce(k kubernetes.Interface, pdb *v1polbeta.PodDisruptionBudget, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found PodDisruptionBudget resource")
	pdbclient := k.PolicyV1beta1().PodDisruptionBudgets(namespace)

	if opts.DryRun {
		_, err := pdbclient.Get(pdb.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s does not exist\n", pdb.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Info(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s exists\n", pdb.Name))
			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = pdbclient.Delete(pdb.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := pdbclient.Get(pdb.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_, err := pdbclient.Create(pdb)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy PodDisruptionBudget resource %s due to %s", pdb.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("PodDisruptionBudget deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := pdbclient.Create(pdb)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy PodDisruptionBudget resource %s due to %s", pdb.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("PodDisruptionBudget deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := pdbclient.Update(pdb)
		if err != nil {
			log.Error("Could not update PodDisruptionBudget")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("PodDisruptionBudget updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := pdbclient.Delete(pdb.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", pdb.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", pdb.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
