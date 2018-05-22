package platform

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
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
func DeployFromFile(config *rest.Config, k kubernetes.Interface, path string, namespace string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, _ := decode(raw, nil, nil)

	fmt.Printf("%++v\n\n", obj.GetObjectKind())

	switch obj.(type) {
	case *v1beta1.Deployment:
		color.Blue("Found deployment resource")
		objdep := obj.(*v1beta1.Deployment)
		deploymentClient := k.AppsV1beta1().Deployments(namespace)
		_, err := deploymentClient.Create(objdep)
		if err != nil {
			color.Blue("Deployment already exists")

		}
	case *v1beta1.StatefulSet:
		color.Blue("Found statefulset resource")
		sts := obj.(*v1beta1.StatefulSet)
		stsclient := k.AppsV1beta1().StatefulSets(namespace)
		_, err := stsclient.Create(sts)
		if err != nil {
			color.Blue("Statefulset already exists")
		}
	case *v1.Service:
		color.Blue("Found service resource")
		ss := obj.(*v1.Service)
		ssclient := k.CoreV1().Services(namespace)
		_, err := ssclient.Create(ss)
		if err != nil {
			color.Blue("Service already exists")
		}
	case *v1.ConfigMap:
		color.Blue("Found Configmap resource")
		cm := obj.(*v1.ConfigMap)
		cmlicnet := k.CoreV1().ConfigMaps(namespace)
		_, err := cmlicnet.Create(cm)
		if err != nil {
			color.Blue("Configmap already exists")
		}
	default:
		color.Red("Unable to convert API resource")
	}

	return nil
}
