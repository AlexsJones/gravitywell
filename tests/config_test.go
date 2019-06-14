package tests

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"testing"
)

func TestLoadMultipleFiles(t *testing.T) {

	conf, e := configuration.NewConfigurationFromPath("./test_dir", []string{})
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
