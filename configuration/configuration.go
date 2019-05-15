package configuration

import (
	"errors"
	"fmt"
	"github.com/AlexsJones/gravitywell/kinds"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
	logger.Info(fmt.Sprintf("Loading %s", path))
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	//for --- found in a file recurse this function
	appc := GravitywellKind{}
	err = yaml.Unmarshal(bytes, &appc)
	if err != nil {
		logger.Fatal(fmt.Sprintf("%+v: %s", err,path))
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
		logger.Info("Application kind found")
		c.ApplicationKinds = append(c.ApplicationKinds, appc)

	case "Cluster":
		appc := kinds.ClusterKind{}
		err = yaml.Unmarshal(bytes, &appc)
		if err != nil {
			return err
		}
		logger.Info("Cluster kind found")
		c.ClusterKinds = append(c.ClusterKinds, appc)
	case "ActionList":
		logger.Info("ActionList kind found")
	default:
		logger.Error("Kind not supported")
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
					logger.Error(fmt.Sprintf("Skipping empty file %s", info.Name()))
					return nil
				}
				if err := LoadConfigurationFromFile(path, conf); err != nil {
					logger.Error(fmt.Sprintf("%s",fmt.Sprintf(err.Error())))
				}
				return nil
			})
		if err != nil {
			return nil, err
		}
	case mode.IsRegular():
		if err := LoadConfigurationFromFile(path, conf); err != nil {
			logger.Error(fmt.Sprintf("%s %s",fmt.Sprintf(err.Error()),path))
		}
	}
	return conf, nil
}
