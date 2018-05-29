package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Cluster struct {
	Name        string `yaml:"Name"`
	Deployments []struct {
		Deployment struct {
			Name      string `yaml:"Name"`
			Namespace string `yaml:"Namespace"`
			Git       string `yaml:"Git"`
			Action    []struct {
				Execute struct {
					Shell   string `yaml:"Shell"`
					Kubectl struct {
						Path    string `yaml:"Path"`
						Type    string `yaml:"Type"`
						Command string `yaml:"Command"`
					} `yaml:"Kubectl"`
				} `yaml:"Execute"`
			} `yaml:"Action"`
		} `yaml:"Deployment"`
	} `yaml:"Deployments"`
}

//Configuration ...
type Configuration struct {
	APIVersion string `yaml:"APIVersion"`
	Strategy   []struct {
		Cluster Cluster `yaml:"Cluster"`
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
