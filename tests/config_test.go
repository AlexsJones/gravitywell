package tests

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"testing"
)

func TestLoadMultipleFiles(t *testing.T) {

	conf, e := configuration.NewConfigurationFromPath("./test_dir")
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	if len(conf.ClusterKinds) == 0 {
		t.Fail()
	}
	if len(conf.ApplicationKinds) == 0 {
		t.Fail()
	}
}


func TestValues(t *testing.T) {

	conf, e := configuration.NewConfigurationFromPath("./test_dir")
	if e != nil {
		fmt.Println(e.Error())
		t.Fail()
	}
	if len(conf.ClusterKinds) != 1 {
		t.Fail()
	}

	if len(conf.ClusterKinds[0].Strategy) != 1 {
		t.Fail()
	}
	if len(conf.ClusterKinds[0].Strategy[0].Provider.Clusters) != 1 {
		t.Fail()
	}

}