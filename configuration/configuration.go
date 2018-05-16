package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

//Configuration generated from https://mengzhuo.github.io/yaml-to-go/
type Configuration struct {
	APIVersion string `yaml:"APIVersion"`
	Strategy   []struct {
		Deployment struct {
			Cluster string `yaml:"Cluster"`
			Name    string `yaml:"Name"`
			Git     string `yaml:"Git"`
			Action  []struct {
				Execute struct {
					Shell   string `yaml:"shell"`
					Kubectl struct {
						Create string `yaml:"create"`
					} `yaml:"kubectl"`
				} `yaml:"Execute"`
			} `yaml:"Action"`
		} `yaml:"Deployment"`
	} `yaml:"Strategy"`
}

//NewConfiguration creates a deserialised yaml object
func NewConfiguration(conf string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		return nil, err
	}
	c := Configuration{}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		log.Printf("Failed to validate syntax: %s", conf)
		return nil, err
	}
	return &c, nil
}
