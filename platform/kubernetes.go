package platform

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	v1polbeta "k8s.io/api/policy/v1beta1"
	v1rbac "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//This is required for gcp auth provider scope
)

// GetKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func GetKubeClient(context string) (*rest.Config, kubernetes.Interface, error) {
	config, err := configForContext(context)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string) (*rest.Config, error) {
	config, err := getConfig(context).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// getConfig returns a Kubernetes client config for a given context.
func getConfig(context string) clientcmd.ClientConfig {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	rules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}

	if context != "" {
		overrides.CurrentContext = context
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
}

//DeployFromFile ...
func DeployFromFile(config *rest.Config, k kubernetes.Interface, path string, namespace string, opts configuration.Options) (state.State, error) {
	f, err := os.Open(path)
	if err != nil {
		return state.EDeploymentStateError, err
	}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return state.EDeploymentStateError, err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, _ := decode(raw, nil, nil)

	log.Printf("%++v\n\n", obj.GetObjectKind())

	var response state.State
	var e error
	switch obj.(type) {
	case *v1beta1.Deployment:
		response, e = execDeploymentResouce(k, obj.(*v1beta1.Deployment), namespace, opts)
	case *v1beta1.StatefulSet:
		response, e = execStatefulSetResouce(k, obj.(*v1beta1.StatefulSet), namespace, opts)
	case *v1.Service:
		response, e = execServiceResouce(k, obj.(*v1.Service), namespace, opts)
	case *v1.ConfigMap:
		response, e = execConfigMapResouce(k, obj.(*v1.ConfigMap), namespace, opts)
	case *v1polbeta.PodDisruptionBudget:
		response, e = execPodDisruptionBudgetResouce(k, obj.(*v1polbeta.PodDisruptionBudget), namespace, opts)
	case *v1.ServiceAccount:
		response, e = execServiceAccountResouce(k, obj.(*v1.ServiceAccount), namespace, opts)
	case *v1rbac.ClusterRoleBinding:

	default:
		color.Red("Unable to convert API resource")
	}

	return response, e
}
