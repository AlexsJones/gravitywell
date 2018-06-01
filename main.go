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
	parallel := flag.Bool("parallel", false, "Run cluster scope deployments in parallel")
	dryRun := flag.Bool("dryrun", false, "Run a dry run deployment to test what is deployment")
	config := flag.String("config", "", "Configuration path")
	flag.Parse()

	if *config == "" {
		return
	}
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

	if err := sh.Run(scheduler.Options{VCS: "git", TempVCSPath: "./staging", APIVersion: "v1", Parallel: *parallel, DryRun: *dryRun}); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

}
