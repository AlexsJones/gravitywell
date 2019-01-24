package configuration

type GitRepo struct {
	GitCryptKey string `yaml:"GitCryptKey"`
	Url         string `yaml:"Url"`
	Branch      string `yaml:"Branch"`
	Path        string `yaml:"Path"`
}
