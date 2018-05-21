package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	"github.com/fatih/color"
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
	sh, err := scheduler.NewScheduler(conf)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	if err := sh.Run(scheduler.Options{VCS: "git", TempVCSPath: "./staging", APIVersion: "v1"}); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

}
