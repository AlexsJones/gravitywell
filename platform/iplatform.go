package platform

import (
	"time"

	"github.com/AlexsJones/kubebuilder/src/data"
	"k8s.io/api/core/v1"
	beta "k8s.io/api/extensions/v1beta1"
)

//IPlatform interface for container platform
type IPlatform interface {
	ValidateDeployment(build *data.BuildDefinition) (bool, error)
	ValidateService(build *data.BuildDefinition) (bool, error)
	ValidateIngress(build *data.BuildDefinition) (bool, error)
	CreateNamespace(string) (*v1.Namespace, error)
	CreateDeployment(build *data.BuildDefinition) (*beta.Deployment, error)
	CreateService(build *data.BuildDefinition) (*v1.Service, error)
	CreateIngress(build *data.BuildDefinition) (*beta.Ingress, error)
	GetIngressLoadBalancerIPAddress(ingress *beta.Ingress, t time.Duration) (string, error)
}

//ValidateDeployment from deserialisation of YAML
func ValidateDeployment(i IPlatform, build *data.BuildDefinition) (bool, error) {
	return i.ValidateDeployment(build)
}

//ValidateService from deserialisation of YAML
func ValidateService(i IPlatform, build *data.BuildDefinition) (bool, error) {
	return i.ValidateService(build)
}

//ValidateIngress from deserialisation of YAML
func ValidateIngress(i IPlatform, build *data.BuildDefinition) (bool, error) {
	return i.ValidateIngress(build)
}

//CreateNamespace within the platform cluster
func CreateNamespace(i IPlatform, namespace string) (*v1.Namespace, error) {
	return i.CreateNamespace(namespace)
}

//CreateDeployment within the platform cluster
func CreateDeployment(i IPlatform, build *data.BuildDefinition) (*beta.Deployment, error) {
	return i.CreateDeployment(build)
}

//CreateService within the platform cluster
func CreateService(i IPlatform, build *data.BuildDefinition) (*v1.Service, error) {
	return i.CreateService(build)
}

//CreateIngress within the platform cluster
func CreateIngress(i IPlatform, build *data.BuildDefinition) (*beta.Ingress, error) {
	return i.CreateIngress(build)
}

//GetIngressLoadBalancerIPAddress with a timeout for the action
func GetIngressLoadBalancerIPAddress(i IPlatform, ingress *beta.Ingress, t time.Duration) (string, error) {
	return i.GetIngressLoadBalancerIPAddress(ingress, t)
}
