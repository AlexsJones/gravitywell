package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	storagev1b1 "k8s.io/api/storage/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Beta1StorageResouce(k kubernetes.Interface, sc *storagev1b1.StorageClass, namespace string, opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Info("Found StorageClass resource")
	client := k.StorageV1beta1().StorageClasses()

	if opts.DryRun {
		_, err := client.Get(sc.Name, v12.GetOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("DRY-RUN: StorageClass resource %s does not exist\n", sc.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			log.Info(fmt.Sprintf("DRY-RUN: StorageClass resource %s exists\n", sc.Name))

			return state.EDeploymentStateExists, nil
		}
	}
	//Replace -------------------------------------------------------------------
	if commandFlag == configuration.Replace {
		log.Debug("Removing resource in preparation for redeploy")
		graceperiod := int64(0)
		_ = client.Delete(sc.Name, &meta_v1.DeleteOptions{GracePeriodSeconds: &graceperiod})
		for {
			_, err := client.Get(sc.Name, meta_v1.GetOptions{})
			if err != nil {
				break
			}
			time.Sleep(time.Second * 1)
		}
		_, err := client.Create(sc)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy StorageClass resource %s due to %s", sc.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("Deployment deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Create ---------------------------------------------------------------------
	if commandFlag == configuration.Create {
		_, err := client.Create(sc)
		if err != nil {
			log.Error(fmt.Sprintf("Could not deploy StorageClass resource %s due to %s", sc.Name, err.Error()))
			return state.EDeploymentStateError, err
		}
		log.Debug("StorageClass deployed")
		return state.EDeploymentStateOkay, nil
	}
	//Apply --------------------------------------------------------------------
	if commandFlag == configuration.Apply {
		_, err := client.Update(sc)
		if err != nil {
			log.Error("Could not update StorageClass")
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug("StorageClass updated")
		return state.EDeploymentStateUpdated, nil
	}
	//Delete -------------------------------------------------------------------
	if commandFlag == configuration.Delete {
		err := client.Delete(sc.Name, &meta_v1.DeleteOptions{})
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete %s", sc.Kind))
			return state.EDeploymentStateCantUpdate, err
		}
		log.Debug(fmt.Sprintf("%s deleted", sc.Kind))
		return state.EDeploymentStateOkay, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
