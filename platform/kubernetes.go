package platform

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	"k8s.io/api/apps/v1beta1"
	v1 "k8s.io/api/core/v1"
	v1polbeta "k8s.io/api/policy/v1beta1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
func DeployFromFile(config *rest.Config, k kubernetes.Interface, path string, namespace string, dryRun bool) error {
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

	log.Printf("%++v\n\n", obj.GetObjectKind())

	switch obj.(type) {
	case *v1beta1.Deployment:
		color.Blue("Found deployment resource")
		objdep := obj.(*v1beta1.Deployment)
		deploymentClient := k.AppsV1beta1().Deployments(namespace)

		if dryRun {
			_, err := deploymentClient.Get(objdep.Name, v12.GetOptions{})
			if err != nil {
				color.Red(fmt.Sprintf("DRY-RUN: Deployment resource %s does not exist\n", objdep.Name))
			} else {
				color.Blue(fmt.Sprintf("DRY-RUN: Deployment resource %s exists\n", objdep.Name))
			}
			return err
		}

		_, err := deploymentClient.Create(objdep)
		if err != nil {
			color.Blue("Deployment already exists")
			_, err := deploymentClient.Update(objdep)
			if err != nil {
				color.Red("Deployment could not be updated")
				return err
			}
			color.Blue("Deployment updated")
		}
	case *v1beta1.StatefulSet:
		color.Blue("Found statefulset resource")
		sts := obj.(*v1beta1.StatefulSet)
		stsclient := k.AppsV1beta1().StatefulSets(namespace)

		if dryRun {
			_, err := stsclient.Get(sts.Name, v12.GetOptions{})
			if err != nil {
				color.Red(fmt.Sprintf("DRY-RUN: StatefulSet resource %s does not exist\n", sts.Name))
			} else {
				color.Blue(fmt.Sprintf("DRY-RUN: StatefulSet resource %s exists\n", sts.Name))
			}
			return err
		}

		_, err := stsclient.Create(sts)
		if err != nil {
			color.Blue("Statefulset already exists")
			_, err := stsclient.UpdateStatus(sts)
			if err != nil {
				color.Red("Could not update Statefulset")
				return err
			}
			color.Blue("Statefulset updated")
		}
	case *v1.Service:
		color.Blue("Found service resource")
		ss := obj.(*v1.Service)
		ssclient := k.CoreV1().Services(namespace)

		if dryRun {
			_, err := ssclient.Get(ss.Name, v12.GetOptions{})
			if err != nil {
				color.Red(fmt.Sprintf("DRY-RUN: Service resource %s does not exist\n", ss.Name))
			} else {
				color.Blue(fmt.Sprintf("DRY-RUN: Service resource %s exists\n", ss.Name))
			}
			return err
		}

		_, err := ssclient.Create(ss)
		if err != nil {
			color.Blue("Service already exists")
			_, err := ssclient.Update(ss)
			if err != nil {
				color.Red("Could not update service")
				return err
			}
			color.Blue("Service updated")
		}
	case *v1.ConfigMap:
		color.Blue("Found Configmap resource")
		cm := obj.(*v1.ConfigMap)
		cmclient := k.CoreV1().ConfigMaps(namespace)

		if dryRun {
			_, err := cmclient.Get(cm.Name, v12.GetOptions{})
			if err != nil {
				color.Red(fmt.Sprintf("DRY-RUN: Configmap resource %s does not exist\n", cm.Name))
			} else {
				color.Blue(fmt.Sprintf("DRY-RUN: Configmap resource %s exists\n", cm.Name))
			}
			return err
		}

		_, err := cmclient.Create(cm)
		if err != nil {
			color.Blue("Configmap already exists")
			_, err := cmclient.Update(cm)
			if err != nil {
				color.Red("Configmap could not be updated")
				return err
			}
			color.Blue("Configmap updated")
		}
	case *v1polbeta.PodDisruptionBudget:
		color.Blue("Found PodDisruptionBudget resource")
		pdb := obj.(*v1polbeta.PodDisruptionBudget)
		pdbclient := k.PolicyV1beta1().PodDisruptionBudgets(namespace)

		if dryRun {
			_, err := pdbclient.Get(pdb.Name, v12.GetOptions{})
			if err != nil {
				color.Red(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s does not exist\n", pdb.Name))
			} else {
				color.Blue(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s exists\n", pdb.Name))
			}
			return err
		}

		_, err := pdbclient.Create(pdb)
		if err != nil {
			color.Blue("PodDisruptionBudget already exists")
			_, err := pdbclient.Update(pdb)
			if err != nil {
				color.Red("PodDisruptionBudget could not be updated")
				return err
			}
			color.Blue("Configmap updated")
		}
	default:
		color.Red("Unable to convert API resource")
	}

	return nil
}
