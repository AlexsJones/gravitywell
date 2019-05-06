package kinds

type Application struct {
	Name       string     `yaml:"Name"`
	Namespace  string     `yaml:"Namespace"`
	Git        string     `yaml:"Git"`
	ActionList ActionList `yaml:"ActionList"`
}

type ApplicationCluster struct {
	ShortName    string `yaml:"ShortName"`
	FullName     string `yaml:"FullName"`
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
