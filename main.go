package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	"github.com/fatih/color"
)

const (
	defaultvcs = "git"
	supportAPI = "v1"
)

func main() {
	_ = flag.Bool("parallel", false, "Run deployments in parallel")
	config := flag.String("config", "", "Configuration path")
	flag.Parse()

	conf, err := configuration.NewConfiguration(*config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if conf.APIVersion != supportAPI {
		color.Red(fmt.Sprintf("Manifest is not supported by the current API: %s\n", supportAPI))
		os.Exit(1)
	}
	sh, err := scheduler.NewScheduler(conf)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	if err := sh.Run(); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

}
