package provider

import (
	"github.com/AlexsJones/gravitywell/kinds"
)

type IProvider interface {
	Create(cluster kinds.ProviderCluster) error

	Delete(cluster kinds.ProviderCluster) error

	List(cluster kinds.ProviderCluster) error
}

func Create(i IProvider, cluster kinds.ProviderCluster) error {

	return i.Create(cluster)
}

func Delete(i IProvider, cluster kinds.ProviderCluster) error {

	return i.Delete(cluster)
}

func List(i IProvider, cluster kinds.ProviderCluster) error {
	return i.List(cluster)
}
