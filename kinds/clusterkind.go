package kinds

//ClusterKind ...
type ClusterKind struct {
	APIVersion string `yaml:"APIVersion"`
	Kind       string `yaml:"Kind"`
	Strategy   []struct {
		Provider Provider `yaml:"Provider"`
	} `yaml:"Strategy"`
}
type ProviderCluster struct {
	InitialNodeCount int               `yaml:"InitialNodeCount"`
	InitialNodeType  string            `yaml:"InitialNodeType"`
	Labels           map[string]string `yaml:"Labels"`
	FullName         string            `yaml:"FullName"`
	ShortName        string            `yaml:"ShortName"`
	Project          string            `yaml:"Project"`
	NodePools        []struct {
		NodePool struct {
			Count    int               `yaml:"Count"`
			Name     string            `yaml:"Name"`
			NodeType string            `yaml:"NodeType"`
			Labels   map[string]string `yaml:"Labels"`
		} `yaml:"NodePool"`
	} `yaml:"NodePools"`
	OauthScopes     string `yaml:"OauthScopes"`
	PostInstallHook []struct {
		Execute struct {
			Shell string `yaml:"Shell"`
			Path  string `yaml:"Path"`
		} `yaml:"Execute"`
	} `yaml:"PostInstallHook"`
	PostDeleteHooak []struct {
		Execute struct {
			Shell string `yaml:"Shell"`
			Path  string `yaml:"Path"`
		} `yaml:"Execute"`
	} `yaml:"PostDeleteHook"`
	Region string   `yaml:"Region"`
	Zones  []string `yaml:"Zones"`
}
type Provider struct {
	Clusters []struct {
		Cluster ProviderCluster `yaml:"Cluster"`
	} `yaml:"Clusters"`
	Name string `yaml:"Name"`
}
