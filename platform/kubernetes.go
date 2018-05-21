package platform

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlexsJones/kubebuilder/src/data"
	appsbetav1 "k8s.io/api/apps/v1beta1"
	v1beta "k8s.io/api/batch/v1beta1"
	"k8s.io/api/core/v1"
	apibetav1 "k8s.io/api/extensions/v1beta1"
	beta "k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
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

//ValidateService from deserialisation of YAML
func ValidateService(build *data.BuildDefinition) (bool, error) {

	//TODO: NEEDS some proper checks here not just deserialisation
	fmt.Println("Attempting deserialisation of YAML")

	d := scheme.Codecs.UniversalDeserializer()
	_, _, err := d.Decode([]byte(build.Kubernetes.Service), nil, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Service, err))
		return false, err
	}
	fmt.Println("Service YAML okay")
	return true, nil
}

//ValidateDeployment from deserialisation of YAML
func ValidateDeployment(k kubernetes.Interface, build *data.BuildDefinition) (bool, error) {
	//TODO: NEEDS some proper checks here not just deserialisation
	fmt.Println("Attempting deserialisation of YAML")

	d := scheme.Codecs.UniversalDeserializer()
	_, _, err := d.Decode([]byte(build.Kubernetes.Deployment), nil, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Deployment, err))
		return false, err
	}
	fmt.Println("Deployment YAML okay")
	return true, nil
}

//ValidateIngress from deserialisation of YAML
func ValidateIngress(k kubernetes.Interface, build *data.BuildDefinition) (bool, error) {
	//TODO: NEEDS some proper checks here not just deserialisation
	fmt.Println("Attempting deserialisation of YAML")

	d := scheme.Codecs.UniversalDeserializer()
	_, _, err := d.Decode([]byte(build.Kubernetes.Ingress), nil, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Ingress, err))
		return false, err
	}
	fmt.Println("Ingress YAML okay")
	return true, nil
}

//CreateNamespace within kubernetes
func CreateNamespace(k kubernetes.Interface, namespace string) (*v1.Namespace, error) {
	if ns, err := GetNamespace(k, namespace); err == nil {
		fmt.Println("Found existing namespace")
		return ns, err
	}
	ns := &v1.Namespace{}

	ns.SetName(namespace)
	ns, err := k.CoreV1().Namespaces().Create(ns)
	return ns, err
}

//GetNamespace within kubernetes
func GetNamespace(k kubernetes.Interface, namespace string) (*v1.Namespace, error) {

	ns, err := k.CoreV1().Namespaces().Get(namespace, meta.GetOptions{})
	return ns, err
}

//CreateDeployment ...
func CreateDeployment(k kubernetes.Interface, build *data.BuildDefinition) (*beta.Deployment, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Deployment), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Deployment, err))
		return nil, err
	}

	deploymentClient := k.ExtensionsV1beta1().Deployments(build.Kubernetes.Namespace)

	deployment, err := deploymentClient.Create(obj.(*beta.Deployment))
	if err != nil {
		fmt.Println("Trying to update existing deployment....")

		deployment, err = deploymentClient.Update(obj.(*beta.Deployment))
		if err == nil {
			fmt.Println("Updated existing deployment")
			return deployment, nil
		}

		return nil, err
	}
	return deployment, nil
}

//CreateService ...
func CreateService(k kubernetes.Interface, build *data.BuildDefinition) (*v1.Service, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Service), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Service, err))
		return nil, err
	}

	serviceClient := k.CoreV1().Services(build.Kubernetes.Namespace)

	service, err := serviceClient.Create(obj.(*v1.Service))
	if err != nil {
		fmt.Println("Trying to update existing service....")

		service, err = serviceClient.Update(obj.(*v1.Service))
		if err == nil {
			fmt.Println("Updated existing service")
			return service, nil
		}

		return nil, err
	}

	return service, nil
}

//CreateIngress ...
func CreateIngress(k kubernetes.Interface, build *data.BuildDefinition) (*beta.Ingress, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Ingress), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Ingress, err))
		return nil, err
	}

	ingressClient := k.ExtensionsV1beta1().Ingresses(build.Kubernetes.Namespace)

	ingress, err := ingressClient.Create(obj.(*beta.Ingress))
	if err != nil {
		fmt.Println("Trying to update existing ingress....")

		ingress, err = ingressClient.Update(obj.(*beta.Ingress))
		if err == nil {
			fmt.Println("Updated existing ingess")
			return ingress, nil
		}

		return nil, err
	}
	return ingress, nil
}

//GetIngressLoadBalancerIPAddress ...
func GetIngressLoadBalancerIPAddress(k kubernetes.Interface, ingress *beta.Ingress, t time.Duration) (string, error) {

	start := time.Now()
	for {

		elapsed := time.Since(start)
		if elapsed > t {
			return "", errors.New("Too much time has elapsed waiting for load balancer")
		}
		if len(ingress.Status.LoadBalancer.Ingress) > 0 {
			if ingress.Status.LoadBalancer.Ingress[0].IP != "" || len(ingress.Status.LoadBalancer.Ingress[0].IP) > 0 {
				return ingress.Status.LoadBalancer.Ingress[0].IP, nil
			}
		}
		time.Sleep(time.Second)
		fmt.Println("Waiting...")
	}
}

//GetNamespaces within kubernetes
func GetNamespaces(k kubernetes.Interface) (*v1.NamespaceList, error) {

	nl, err := k.CoreV1().Namespaces().List(meta.ListOptions{})

	return nl, err
}

//GetPods within kubernetes
func GetPods(k kubernetes.Interface, namespace string) (*v1.PodList, error) {

	nl, err := k.CoreV1().Pods(namespace).List(meta.ListOptions{})

	return nl, err
}

//GetServices within kubernetes
func GetServices(k kubernetes.Interface, namespace string) (*v1.ServiceList, error) {

	nl, err := k.CoreV1().Services(namespace).List(meta.ListOptions{})

	return nl, err
}

//GetDeployments within kubernetes
func GetDeployments(k kubernetes.Interface, namespace string) (*apibetav1.DeploymentList, error) {

	nl, err := k.ExtensionsV1beta1().Deployments(namespace).List(meta.ListOptions{})

	return nl, err
}

//GetStatefulSets within kubernetes
func GetStatefulSets(k kubernetes.Interface, namespace string) (*appsbetav1.StatefulSetList, error) {
	nl, err := k.AppsV1beta1().StatefulSets(namespace).List(meta.ListOptions{})
	return nl, err
}

//GetCronJobs within kubernetes
func GetCronJobs(k kubernetes.Interface, namespace string) (*v1beta.CronJobList, error) {
	nl, err := k.BatchV1beta1().CronJobs(namespace).List(meta.ListOptions{})
	return nl, err
}
