package main

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	log "github.com/Sirupsen/logrus"
	"github.com/jessevdk/go-flags"
	"os"
	"strings"
	"time"
)

var (
	version = "dev"
)

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

var Opts struct {
	DryRun     bool   `short:"d" long:"dryrun" description:"Performs a dryrun."`
	FileName   string `short:"f" long:"filename" description:"filename to execute, also accepts a path."`
	SSHKeyPath string `short:"s" long:"sshkeypath" description:"Custom ssh key path."`
	MaxTimeout string `short:"m" long:"maxtimeout" description:"Max rollout time e.g. 60s or 1m"`
}

func Usage() {

	fmt.Println("...Usage...")

	fmt.Println("create/delete/replace/apply e.g. gravitywell create -f folder/")
	os.Exit(0)
}
func main() {
	args := os.Args
	var command = ""
	if len(args) == 2 && args[1] == "version" {
		fmt.Println(version)
		os.Exit(0)
	}
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
	case "delete":
		commandFlag = configuration.Delete
	default:
		fmt.Println("Command not recognised.")
		os.Exit(1)
	}

	defaultMaxTimeout := time.Second * 60

	d, err := time.ParseDuration(Opts.MaxTimeout)
	if err == nil {
		defaultMaxTimeout = d
	}

	cf := configuration.Options{VCS: "git",
		TempVCSPath:        "./.gravitywell",
		APIVersion:         "v1",
		SSHKeyPath:         Opts.SSHKeyPath,
		MaxBackOffDuration: defaultMaxTimeout,
	}

	if _, err := os.Stat(cf.TempVCSPath); os.IsNotExist(err) {
		err = os.Mkdir(cf.TempVCSPath, 0777)
	} else {
		err = os.RemoveAll(cf.TempVCSPath)
		err = os.Mkdir(cf.TempVCSPath, 0777)
	}
	if err :=
		sh.Run(commandFlag, cf); err != nil {
		log.Warn(err.Error())
		os.Exit(1)
	}

}
