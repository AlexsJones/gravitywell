package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlexsJones/ashara/configuration"
	"github.com/fatih/color"
)

const (
	defaultvcs = "git"
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

	for _, d := range conf.Strategy {
		color.Yellow(fmt.Sprintf("Attempting to fetch %s\n", d.Deployment.Name))

	}
}
