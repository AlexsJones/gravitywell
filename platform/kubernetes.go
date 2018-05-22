package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/fatih/color"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
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

func DeployFromFile(config *rest.Config, k kubernetes.Interface, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	d := yaml.NewYAMLOrJSONDecoder(f, 4096)
	dd := k.Discovery()
	apigroups, err := discovery.GetAPIGroupResources(dd)
	if err != nil {
		return err
	}
	restmapper := discovery.NewRESTMapper(apigroups, meta.InterfacesForUnstructured)

	for {
		ext := runtime.RawExtension{}
		if err := d.Decode(&ext); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println("raw: ", string(ext.Raw))
		versions := &runtime.VersionedObjects{}
		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, versions)
		fmt.Println("obj: ", obj)
		mapping, err := restmapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			log.Fatal(err)
		}
		restconfig := config
		restconfig.GroupVersion = &schema.GroupVersion{
			Group:   mapping.GroupVersionKind.Group,
			Version: mapping.GroupVersionKind.Version,
		}
		dclient, err := dynamic.NewClient(restconfig)
		if err != nil {
			log.Fatal(err)
		}
		apiresourcelist, err := dd.ServerResources()
		if err != nil {
			log.Fatal(err)
		}
		var myapiresource metav1.APIResource
		for _, apiresourcegroup := range apiresourcelist {
			if apiresourcegroup.GroupVersion == mapping.GroupVersionKind.Version {
				for _, apiresource := range apiresourcegroup.APIResources {
					if apiresource.Name == mapping.Resource && apiresource.Kind == mapping.GroupVersionKind.Kind {
						myapiresource = apiresource
					}
				}
			}
		}
		fmt.Println(myapiresource)
		var unstruct unstructured.Unstructured
		unstruct.Object = make(map[string]interface{})
		var blob interface{}
		if a := json.Unmarshal(ext.Raw, &blob); err != nil {
			color.Red(a.Error())
		}
		unstruct.Object = blob.(map[string]interface{})
		fmt.Println("unstruct:", unstruct)
		ns := "default"
		if md, ok := unstruct.Object["metadata"]; ok {
			metadata := md.(map[string]interface{})
			if internalns, ok := metadata["namespace"]; ok {
				ns = internalns.(string)
			}
		}
		res := dclient.Resource(&myapiresource, ns)
		fmt.Println(res)
		us, err := res.Create(&unstruct)
		if err != nil {
			color.Red(err.Error())
		}
		fmt.Println("unstruct response:", us)

	}
	return nil
}
