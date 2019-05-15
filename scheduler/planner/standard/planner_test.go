package standard

import (
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler/planner"
	"testing"
)

func TestPlanner(t *testing.T) {

	conf, err := configuration.NewConfigurationFromPath("../../../example-gcp/templates")
	if err != nil {
		t.Fatal(err)
	}
	stdplnr := StandardPlanner{}

	_, err = planner.GeneratePlan(stdplnr, conf, configuration.Create,
		configuration.Options{})
	if err != nil {
		t.Fatal(err)
	}

}
