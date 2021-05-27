package configuration

import (
	"testing"
)

func TestConfigurationBasic(t *testing.T) {

	conf, err := NewConfigurationFromPath("../examples/aws/templates", []string{})
	if err != nil {
		t.Fatal(err)
	}

	if len(conf.ClusterKinds) != 1 {
		t.Error()
	}
}
func TestConfigurationAdv(t *testing.T) {

	conf, err := NewConfigurationFromPath("../examples/common/templates", []string{})
	if err != nil {
		t.Fatal(err)
	}

	if len(conf.ApplicationKinds) != 3 {
		t.Error()
	}
}
