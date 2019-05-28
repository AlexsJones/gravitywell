package kinds

type Application struct {
	Name         string     `yaml:"Name"`
	Namespace    string     `yaml:"Namespace"`
	VCS struct {
		FileSystem string `yaml:"FileSystem"`
	Git          string     `yaml:"Git"`
	GitReference string     `yaml:"GitReference"`
	} `yaml:"VCS"`
	ActionList   ActionList `yaml:"ActionList"`
}

type ApplicationCluster struct {
	Name    string `yaml:"Name"`
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
