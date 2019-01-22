package configuration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

type Action struct {
		Execute struct {
			Kind          string            `yaml:"Kind"`
			Configuration map[string]string `yaml:"Configuration"`
		} `yaml:"Execute"`
}

type Application struct {
	Name      string `yaml:"Name"`
	Namespace string `yaml:"Namespace"`
	Git       string `yaml:"Git"`
	Action    []Action `yaml:"Action"`
}

type ApplicationCluster struct {
	Name         string `yaml:"Name"`
	Applications []struct {
		Application Application `yaml:"Application"`
	} `yaml:"Applications"`
}

type ProviderCluster struct {
	InitialNodeCount int               `yaml:"InitialNodeCount"`
	InitialNodeType  string            `yaml:"InitialNodeType"`
	Labels           map[string]string `yaml:"Labels"`
	Name             string            `yaml:"Name"`
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

//ApplicationKind ...
type ApplicationKind struct {
	APIVersion string `yaml:"APIVersion"`
	Strategy   []struct {
		Cluster ApplicationCluster `yaml:"Cluster"`
	} `yaml:"Strategy"`
}

//ClusterKind ...
type ClusterKind struct {
	APIVersion string `yaml:"APIVersion"`
	Kind       string `yaml:"Kind"`
	Strategy   []struct {
		Provider Provider `yaml:"Provider"`
	} `yaml:"Strategy"`
}

//GravitywellKind ...
type GravitywellKind struct {
	APIVersion string `yaml:"APIVersion"`
	Kind       string `yaml:"Kind"`
}

type Configuration struct {
	ApplicationKinds []ApplicationKind
	ClusterKinds     []ClusterKind
}

func LoadConfigurationFromFile(path string, c *Configuration) error {
	log.Println(fmt.Sprintf("Loading %s", path))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	//for --- found in a file recurse this function
	appc := GravitywellKind{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		color.Red(fmt.Sprintf("%+v", err))
		os.Exit(1)
		return err
	}
	//Load specific kind
	switch appc.Kind {
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
		return nil, err
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
					color.Red(fmt.Sprintf("Skipping empty file %s", info.Name()))
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
	return conf, nil
}
