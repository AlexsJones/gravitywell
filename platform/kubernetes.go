package platform

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AlexsJones/kubebuilder/src/data"
	"k8s.io/api/core/v1"
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

//Kubernetes ...
type Kubernetes struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

//NewKubernetes object
func NewKubernetes(masterURL string, inclusterConfig bool) (*Kubernetes, error) {

	//InCluster...
	var config *rest.Config
	var err error
	if inclusterConfig {
		fmt.Println("Using in cluster configuration")
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Using out of cluster configuration")
		kubeconfig := filepath.Join(func() string {
			if h := os.Getenv("HOME"); h != "" {
				return h
			}
			return os.Getenv("USERPROFILE")
		}(), ".kube", "config")
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Kubernetes{clientset: clientset, config: config}, nil
}

//ValidateService from deserialisation of YAML
func (k *Kubernetes) ValidateService(build *data.BuildDefinition) (bool, error) {

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
func (k *Kubernetes) ValidateDeployment(build *data.BuildDefinition) (bool, error) {
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
func (k *Kubernetes) ValidateIngress(build *data.BuildDefinition) (bool, error) {
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
func (k *Kubernetes) CreateNamespace(namespace string) (*v1.Namespace, error) {
	if ns, err := k.GetNamespace(namespace); err == nil {
		fmt.Println("Found existing namespace")
		return ns, err
	}
	ns := &v1.Namespace{}

	ns.SetName(namespace)

	ns, err := k.clientset.CoreV1().Namespaces().Create(ns)
	return ns, err
}

//GetNamespace within kubernetes
func (k *Kubernetes) GetNamespace(namespace string) (*v1.Namespace, error) {

	ns, err := k.clientset.CoreV1().Namespaces().Get(namespace, meta.GetOptions{})
	return ns, err
}

//CreateDeployment ...
func (k *Kubernetes) CreateDeployment(build *data.BuildDefinition) (*beta.Deployment, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Deployment), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Deployment, err))
		return nil, err
	}

	deploymentClient := k.clientset.ExtensionsV1beta1().Deployments(build.Kubernetes.Namespace)

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
func (k *Kubernetes) CreateService(build *data.BuildDefinition) (*v1.Service, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Service), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Service, err))
		return nil, err
	}

	serviceClient := k.clientset.CoreV1().Services(build.Kubernetes.Namespace)

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
func (k *Kubernetes) CreateIngress(build *data.BuildDefinition) (*beta.Ingress, error) {

	deserializer := serializer.NewCodecFactory(clientsetscheme.Scheme).UniversalDeserializer()
	obj, _, err := deserializer.Decode([]byte(build.Kubernetes.Ingress), nil, nil)

	if err != nil {
		fmt.Println(fmt.Sprintf("could not decode yaml: %s\n%s", build.Kubernetes.Ingress, err))
		return nil, err
	}

	ingressClient := k.clientset.ExtensionsV1beta1().Ingresses(build.Kubernetes.Namespace)

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
func (k *Kubernetes) GetIngressLoadBalancerIPAddress(ingress *beta.Ingress, t time.Duration) (string, error) {

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
