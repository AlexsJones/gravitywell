package configuration

type Vault struct {
	Url         string  `yaml:"Url"`
	Path        string  `yaml:"Path"`
	Description string  `yaml:"Description"`
	Repo        GitRepo `yaml:"Repo"`
}
