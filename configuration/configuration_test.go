package configuration

import "testing"

func TestConfigurationBasic(t *testing.T) {

	conf, err := NewConfigurationFromPath("../examples")
	if err != nil {
		t.Fatal(err)
	}

	if len(conf.ClusterKinds) != 1 {
		t.Error()
	}
	if len(conf.ApplicationKinds) != 2 {
		t.Error()
	}
}
