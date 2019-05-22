package main

import (
	"bytes"
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	"github.com/dimiro1/banner"
	"github.com/jessevdk/go-flags"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var (
	version = "dev"
)

var b = `
{{ .AnsiColor.Red }}                               .__  __                         .__  .__   
{{ .AnsiColor.Blue }}        ________________ ___  _|__|/  |_ ___.__.__  _  __ ____ |  | |  |  
{{ .AnsiColor.Yellow }}       / ___\_  __ \__  \\  \/ /  \   __<   |  |\ \/ \/ // __ \|  | |  |  
{{ .AnsiColor.Green }}      / /_/  >  | \// __ \\   /|  ||  |  \___  | \     /\  ___/|  |_|  |__
{{ .AnsiColor.Magenta }}      \___  /|__|  (____  /\_/ |__||__|  / ____|  \/\_/  \___  >____/____/
{{ .AnsiColor.Cyan }}     /_____/            \/               \/                  \/           
{{ .AnsiColor.Default }}
{{ .Env "GW_VERSION" }}
`

var Opts struct {
	DryRun       bool     `short:"d" long:"dryrun" description:"Performs a dryrun."`
	FileName     string   `short:"f" long:"filename" description:"filename to execute, also accepts a path."required:"yes"`
	SSHKeyPath   string   `short:"s" long:"sshkeypath" description:"Custom ssh key path."`
	MaxTimeout   string   `short:"m" long:"maxtimeout" description:"Max rollout time e.g. 60s or 1m"`
	Verbose      bool     `short:"v" long:"verbose" description:"Enable verbose logging"`
	Force        bool     `short:"n" long:"force" description:"Force services to apply even if immutable"`
	IgnoreFilter []string `short:"i" long:"ignore" description:"Ignore excepts any partial string to test and ignore paths/directories with e.g. --ignore=cluster --ignore=actionlist"`
}

func Usage() {

	fmt.Println("the required command `[create|delete|apply|replace]' was not specified")
	os.Exit(0)
}

func init() {

	logger.SetOutput(os.Stdout)

	logger.SetLevel(logger.InfoLevel)
}

func main() {
	isEnabled := true
	isColorEnabled := true

	err := os.Setenv("GW_VERSION", fmt.Sprintf("Build version: %s", version))
	if err != nil {
		fmt.Printf(err.Error())
	}
	banner.Init(os.Stdout, isEnabled, isColorEnabled, bytes.NewBufferString(b))
	//Parse Args-----------------------------------------------------------------------
	//Pull the command out of the flags
	if len(os.Args) < 2 {
		Usage()
	}

	if os.Args[1] == "" {
		Usage()
	}
	command := strings.ToLower(os.Args[1])
	_, err = flags.ParseArgs(&Opts, os.Args)

	if err != nil {
		os.Exit(0)
	}
	//----------------------------------------------------------------------------------
	conf, err := configuration.NewConfigurationFromPath(Opts.FileName, Opts.IgnoreFilter)
	if err != nil {
		logger.Fatalf(err.Error())
	}

	sh, err := scheduler.NewScheduler(conf)
	if err != nil {
		logger.Fatalf(err.Error())
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
		DryRun:             Opts.DryRun,
		Force:              Opts.Force,
		IgnoreFilter:       Opts.IgnoreFilter,
	}

	if _, err := os.Stat(cf.TempVCSPath); os.IsNotExist(err) {
		err = os.Mkdir(cf.TempVCSPath, 0777)
	} else {
		err = os.RemoveAll(cf.TempVCSPath)
		err = os.Mkdir(cf.TempVCSPath, 0777)
	}
	if err :=
		sh.Run(commandFlag, cf); err != nil {
		logger.Fatalf(err.Error())
	}
}
