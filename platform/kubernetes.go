package platform

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/state"
	"github.com/google/logger"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/apps/v1beta1"
	"k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchbeta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	v1betav1 "k8s.io/api/extensions/v1beta1"
	v1polbeta "k8s.io/api/policy/v1beta1"
	v1rbac "k8s.io/api/rbac/v1"
	storagev1b1 "k8s.io/api/storage/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
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

//GenerateDeploymentPlan
func GenerateDeploymentPlan(k kubernetes.Interface,
	files []string, namespace string, opts configuration.Options,
	commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) error {

	var kubernetesResources []runtime.Object
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			logger.Warning(fmt.Sprintf("Could not open from file %s", file))
			continue
		}
		raw, err := ioutil.ReadAll(f)
		if err != nil {
			logger.Warning(fmt.Sprintf("Could not read from file %s", file))
			continue
		}
		yamldelimiter := regexp.MustCompile(`(\A|\n)---`)
		documents := yamldelimiter.Split(string(raw), -1)
		for _, doc := range documents {
			//Decode into kubernetes object
			decode := scheme.Codecs.UniversalDeserializer().Decode
			obj, kind, err := decode([]byte(doc), nil, nil)
			if err != nil {
				logger.Warning(fmt.Sprintf("%s : %s", err.Error(), file))
				continue
			}
			log.Printf("Decoded Kind: %s", kind.String())

			kubernetesResources = append(kubernetesResources, obj)
		}
	}

	//TODO: Deployment plan printing

	if len(kubernetesResources) == 0 {
		return errors.New("no resources within file list")
	}

	//TODO: Run X resource first

	//Run namespace first
	out := 0
	for _, resource := range kubernetesResources {
		gvk := resource.GetObjectKind().GroupVersionKind()

		switch strings.ToLower(gvk.Kind) {
		case "namespace":
			//Remove the namespace from the array and run first
			_, err := execV1NamespaceResource(k, resource.(*v1.Namespace), namespace, opts, commandFlag)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
		default:
			kubernetesResources[out] = resource
			out++
		}
	}
	kubernetesResources = kubernetesResources[:out]

	//Run all other resources
	for _, resource := range kubernetesResources {

		s, err := DeployFromObject(k, resource, namespace, opts, commandFlag, shouldAwaitDeployment)
		if err != nil {
			logger.Error(fmt.Sprintf("%s : %s", err.Error(), resource.GetObjectKind().GroupVersionKind().Kind))
			continue
		}

		switch s {
		case state.EDeploymentStateError:
			color.Red(fmt.Sprintf("%s STATE: %s",
				resource.GetObjectKind().GroupVersionKind().Kind, state.Translate(s)))
		default:
			color.Green(fmt.Sprintf("%s STATE: %s",
				resource.GetObjectKind().GroupVersionKind().Kind, state.Translate(s)))
		}
	}
	return nil
}

//DeployFromObject ...
func DeployFromObject(k kubernetes.Interface, obj runtime.Object,
	namespace string, opts configuration.Options,
	commandFlag configuration.CommandFlag, shouldAwaitDeployment bool) (state.State, error) {

	var response state.State
	var e error
	switch obj.(type) {
	case *v1.Pod:
		response, e = execV1Pod(k, obj.(*v1.Pod), namespace, opts, commandFlag)
	case *batchbeta1.CronJob:
		response, e = execV1Beta1CronJob(k, obj.(*batchbeta1.CronJob), namespace, opts, commandFlag)
	case *batchv1.Job:
		response, e = execV1Job(k, obj.(*batchv1.Job), namespace, opts, commandFlag)
	case *v1betav1.Deployment:
		response, e = execV1Betav1DeploymentResouce(k, obj.(*v1betav1.Deployment),
			namespace, opts, commandFlag, shouldAwaitDeployment)
	case *v1beta1.Deployment:
		response, e = execV1Beta1DeploymentResouce(k, obj.(*v1beta1.Deployment),
			namespace, opts, commandFlag, shouldAwaitDeployment)
	case *v1beta2.Deployment:
		response, e = execV1Beta2DeploymentResouce(k, obj.(*v1beta2.Deployment),
			namespace, opts, commandFlag, shouldAwaitDeployment)
	case *v1beta1.StatefulSet:
		response, e = execV1Beta1StatefulSetResouce(k, obj.(*v1beta1.StatefulSet), namespace, opts, commandFlag, shouldAwaitDeployment)
	case *appsv1.StatefulSet:
		response, e = execV1StatefulSetResouce(k, obj.(*appsv1.StatefulSet),
			namespace, opts, commandFlag, shouldAwaitDeployment)
	case *v1.Secret:
		response, e = execV1SecretResouce(k, obj.(*v1.Secret), namespace, opts, commandFlag)
	case *v1.Service:
		response, e = execV1ServiceResouce(k, obj.(*v1.Service), namespace, opts, commandFlag)
	case *v1.ConfigMap:
		response, e = execV1ConfigMapResouce(k, obj.(*v1.ConfigMap), namespace, opts, commandFlag)
	case *v1polbeta.PodDisruptionBudget:
		response, e = execV1Beta1PodDisruptionBudgetResouce(k, obj.(*v1polbeta.PodDisruptionBudget), namespace, opts, commandFlag)
	case *v1.ServiceAccount:
		response, e = execV1ServiceAccountResouce(k, obj.(*v1.ServiceAccount), namespace, opts, commandFlag)
	//V1 RBAC
	case *v1rbac.ClusterRoleBinding:
		response, e = execV1RbacClusterRoleBindingResouce(k, obj.(*v1rbac.ClusterRoleBinding), namespace, opts, commandFlag)
	case *v1rbac.Role:
		response, e = execV1RbacRoleResouce(k, obj.(*v1rbac.Role), namespace, opts, commandFlag)
	case *v1rbac.RoleBinding:
		response, e = exec1VRbacRoleBindingResouce(k, obj.(*v1rbac.RoleBinding), namespace, opts, commandFlag)
	case *v1rbac.ClusterRole:
		response, e = execV1AuthClusterRoleResouce(k, obj.(*v1rbac.ClusterRole), namespace, opts, commandFlag)
	case *v1betav1.DaemonSet:
		response, e = execV1Beta1DaemonSetResouce(k, obj.(*v1betav1.DaemonSet), namespace, opts, commandFlag)
	case *v1betav1.Ingress:
		response, e = execV1Beta1IngressResouce(k, obj.(*v1betav1.Ingress), namespace, opts, commandFlag)
	case *storagev1b1.StorageClass:
		response, e = execV1Beta1StorageResouce(k, obj.(*storagev1b1.StorageClass), namespace, opts, commandFlag)
	default:
		logger.Error("Unable to convert API resource:", obj.GetObjectKind().GroupVersionKind())
	}

	return response, e
}
