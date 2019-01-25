package configuration

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

//GravitywellKind ...
type GravitywellKind struct {
	APIVersion string `yaml:"APIVersion"`
	Kind       string `yaml:"Kind"`
}

type Configuration struct {
	ApplicationKinds []kinds.ApplicationKind
	ClusterKinds     []kinds.ClusterKind
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
		appc := kinds.ApplicationKind{}
		err = yaml.Unmarshal(bytes, &appc)
		if err != nil {
			return err
		}
		color.Yellow("Application kind found")
		c.ApplicationKinds = append(c.ApplicationKinds, appc)

	case "Cluster":
		appc := kinds.ClusterKind{}
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
