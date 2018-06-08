package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	"github.com/fatih/color"
)

func main() {
	redeploy := flag.Bool("redeploy", false, "Forces a delete and deploy WARNING: Destructive")
	parallel := flag.Bool("parallel", false, "Run cluster scope deployments in parallel - best not to use if pulling parallel from the same git repo")
	tryUpdate := flag.Bool("tryupdate", false, "Try to update the resource if possible")
	sshkeypath := flag.String("sshkeypath", "", "Provide to override default sshkey used")
	dryRun := flag.Bool("dryrun", false, "Run a dry run deployment to test what is deployment")
	config := flag.String("config", "", "Configuration path")
	flag.Parse()

	if *config == "" {
		return
	}

	if *redeploy {
		reader := bufio.NewReader(os.Stdin)
		color.Red(fmt.Sprintf("This is a very destructive action, are you sure [Y/N]?: "))
		text, _ := reader.ReadString('\n')
		trimmed := strings.Trim(text, "\n")
		if strings.Compare(trimmed, "Y") != 0 {

			os.Exit(0)
		}
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

	if err := sh.Run(configuration.Options{VCS: "git", TempVCSPath: "./staging", APIVersion: "v1", SSHKeyPath: *sshkeypath, Parallel: *parallel, DryRun: *dryRun, TryUpdate: *tryUpdate, Redeploy: *redeploy}); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

}
