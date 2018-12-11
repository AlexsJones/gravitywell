package main

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

var Opts struct {
	DryRun bool `short:"d" long:"dryrun" description:"Performs a dryrun."`
	FileName string `short:"f" long:"filename" description:"filename to execute, also accepts a path."`
	SSHKeyPath string `short:"s" long:"sshkeypath" description:"Custom ssh key path."`
}
func Usage() {

	os.Exit(0)
}
func main() {

	args := os.Args
	var command = ""
	if len(args) <= 2 {
		Usage()
	}

	if args[1] == "" {
		Usage()
	}
	command = strings.ToLower(args[1])
	if command == "" {
		Usage()
	}
	args = args[2:len(args)]

	_, err := flags.ParseArgs(&Opts, os.Args)
	if err != nil {
		panic(err)
	}

	conf, err := configuration.NewConfigurationFromPath(Opts.FileName)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	sh, err := scheduler.NewScheduler(conf)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	var commandFlag configuration.CommandFlag
	switch command {
	case "create":
		commandFlag = configuration.Create
	case "apply":
		commandFlag = configuration.Apply
	case "replace":
		commandFlag = configuration.Replace
	default:
		fmt.Println("Command not recognised.")
		os.Exit(1)
	}

	if err :=
		sh.Run(commandFlag, configuration.Options{VCS: "git",
	TempVCSPath: "./.gravitywell",
	APIVersion: "v1",
	SSHKeyPath: Opts.SSHKeyPath,
	});
	err != nil {
		log.Warn(err.Error())
		os.Exit(1)
	}

}
