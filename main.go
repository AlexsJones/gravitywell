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
	logDir  = "logs"
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
	DryRun     bool   `short:"d" long:"dryrun" description:"Performs a dryrun."`
	FileName   string `short:"f" long:"filename" description:"filename to execute, also accepts a path."`
	SSHKeyPath string `short:"s" long:"sshkeypath" description:"Custom ssh key path."`
	MaxTimeout string `short:"m" long:"maxtimeout" description:"Max rollout time e.g. 60s or 1m"`
	Verbose    bool   `short:"v" long:"verbose" description:"Enable verbose logging"`
}

func Usage() {

	fmt.Println("...Usage...")

	fmt.Println("create/delete/replace/apply e.g. gravitywell create -f folder/")
	os.Exit(0)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logger.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	logger.SetLevel(logger.InfoLevel)
}

func main() {
	isEnabled := true
	isColorEnabled := true

	err := os.Setenv("GW_VERSION",fmt.Sprintf("Build version: %s",version))
	if err != nil {
		logger.Warning(err.Error())
	}
	banner.Init(os.Stdout, isEnabled, isColorEnabled, bytes.NewBufferString(b))
	//Parse Args-----------------------------------------------------------------------
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

	_, err = flags.ParseArgs(&Opts, os.Args)

	if err != nil {
		panic(err)
	}
	//----------------------------------------------------------------------------------
	conf, err := configuration.NewConfigurationFromPath(Opts.FileName)
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
