package configuration

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Strategy struct {
}

//Configuration is the YAML structure representation
type Configuration struct {
	APIVersion string    `yaml:"APIVersion"`
	Strategy   *Strategy `yaml:"Strategy"`
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
