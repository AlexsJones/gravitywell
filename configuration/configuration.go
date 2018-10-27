package configuration

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type Cluster struct {
	Name        string `yaml:"Name"`
	Applications []struct {
		Application struct {
			Name            string `yaml:"Name"`
			Namespace       string `yaml:"Namespace"`
			CreateNamespace bool   `yaml:"CreateNamespace"`
			Git             string `yaml:"Git"`
			Action          []struct {
				Execute struct {
					Shell   string `yaml:"Shell"`
					Kubectl struct {
						Path    string `yaml:"Path"`
						Type    string `yaml:"Type"`
						Command string `yaml:"Command"`
					} `yaml:"Kubectl"`
				} `yaml:"Execute"`
			} `yaml:"Action"`
		} `yaml:"Application"`
	} `yaml:"Applications"`
}

//ApplicationKind ...
type ApplicationKind struct {
	APIVersion string `yaml:"APIVersion"`
	Strategy   []struct {
		Cluster Cluster `yaml:"Cluster"`
	} `yaml:"Strategy"`
}

//ClusterKind ...
type ClusterKind struct {

}
//GravitywellKind ...
type GravitywellKind struct {
	APIVersion string `yaml:"APIVersion"`
	Kind string `yaml:"Kind"`
}

type Configuration struct {
	ApplicationKinds []ApplicationKind
	ClusterKinds []ClusterKind
}

func LoadConfigurationFromFile(path string, c *Configuration) error {
	log.Println(fmt.Sprintf("Loading %s", path))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	appc := GravitywellKind{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		return err
	}
	//Load specific kind
	switch appc.Kind{
	case "Application":
		appc := ApplicationKind{}
		err = yaml.Unmarshal(bytes, &appc)
		if err != nil {
			return err
		}
		color.Yellow("Application kind found")
		c.ApplicationKinds = append(c.ApplicationKinds, appc)
	case "Cluster":
		appc := ClusterKind{}
		err = yaml.Unmarshal(bytes, &appc)
		if err != nil {
			return err
		}
		color.Yellow("Cluster kind found")
		c.ClusterKinds = append(c.ClusterKinds, appc)
	default:
		color.Red("Kind not supported")
		return errors.New("kind not supported")
	}
	return nil
}

func NewConfigurationFromPath(path string) (*Configuration, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil,err
	}

	conf := &Configuration{}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		err := filepath.Walk(path,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.Size() == 0 {
					color.Red(fmt.Sprintf("Skipping empty file %s",info.Name()))
					return nil
				}
				LoadConfigurationFromFile(path, conf)
				return nil
			})
		if err != nil {
			return nil, err
		}
	case mode.IsRegular():
		LoadConfigurationFromFile(path, conf)
	}
	return conf,nil
}
