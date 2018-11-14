package platform

import (
	"encoding/json"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	log "github.com/Sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execV1Beta1StatefulSetResouce(k kubernetes.Interface, sts *v1beta1.StatefulSet,
	application configuration.Application,
	executionStep configuration.Execute,
	opts configuration.Options, commandFlag configuration.CommandFlag) (state.State, error) {
	log.Debug("Found statefulset resource")
	stsclient := k.AppsV1beta1().StatefulSets(application.Namespace)

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
	//Patch --------------------------------------------------------------------
	if commandFlag == configuration.Patch {

		type patchStringValue struct {
			Op    string `json:"op"`
			Path  string `json:"path"`
			Value string `json:"value"`
		}

		if executionStep.Kubectl.Patch.Op == "" {
			log.Warn("Missing Op value for patch")
			return state.EDeploymentStateError, errors.New("Missing Op value for patch")
		}
		if executionStep.Kubectl.Patch.Path == "" {
			log.Warn("Missing Path value for patch")
			return state.EDeploymentStateError, errors.New("Missing Path value for patch")
		}
		if executionStep.Kubectl.Patch.Value == "" {
			log.Warn("Missing Value for patch")
			return state.EDeploymentStateError, errors.New("Missing Value for patch")
		}

		payload := []patchStringValue{{
			Op:    executionStep.Kubectl.Patch.Op,
			Path:  executionStep.Kubectl.Patch.Path,
			Value: executionStep.Kubectl.Patch.Value,
		}}

		payloadBytes, _ := json.Marshal(payload)

		_, err := stsclient.Patch(sts.Name, types.JSONPatchType, payloadBytes)
		if err != nil {
			log.Error("Could not update Deployment")
			return state.EDeploymentStateCantUpdate, err
		}

		log.Debug("Deployment patched")
		return state.EDeploymentStateUpdated, nil
	}
	return state.EDeploymentStateNil, errors.New("No kubectl command")
}
