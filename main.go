package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AlexsJones/asana/configuration"
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
	log.Println(conf)
}
