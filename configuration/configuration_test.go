package configuration

import "testing"

func TestConfigurationBasic(t *testing.T) {

	conf, err := NewConfigurationFromPath("../example-aws/templates")
	if err != nil {
		t.Fatal(err)
	}

	if len(conf.ClusterKinds) != 1 {
		t.Error()
	}
}
func TestConfigurationAdv(t *testing.T) {

	conf, err := NewConfigurationFromPath("../example-gcp/templates")
	if err != nil {
		t.Fatal(err)
	}

	if len(conf.ClusterKinds) != 1 {
		t.Error()
	}
	if len(conf.ApplicationKinds) != 4 {
		t.Error()
	}
}
