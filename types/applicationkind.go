package types

type Action struct {
	Execute struct {
		Kind          string            `yaml:"Kind"`
		Configuration map[string]string `yaml:"Configuration"`
	} `yaml:"Execute"`
}

type Application struct {
	Name      string   `yaml:"Name"`
	Namespace string   `yaml:"Namespace"`
	Git       string   `yaml:"Git"`
	Action    []Action `yaml:"Action"`
}

type ApplicationCluster struct {
	Name         string `yaml:"Name"`
	Applications []struct {
		Application Application `yaml:"Application"`
	} `yaml:"Applications"`
}

//ApplicationKind ...
type ApplicationKind struct {
	APIVersion string `yaml:"APIVersion"`
	Strategy   []struct {
		Cluster ApplicationCluster `yaml:"Cluster"`
	} `yaml:"Strategy"`
}
