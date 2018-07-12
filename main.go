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
	tryUpdate := flag.Bool("try-update", false, "Try to update the resource if possible")
	ignoreList := flag.String("ignore-list", "", "A comma delimited list of clusters to ignore")
	sshkeypath := flag.String("ssh-key-path", "", "Provide to override default sshkey used")
	dryRun := flag.Bool("dry-run", false, "Run a dry run deployment to test what is deployment")
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

	var ignoreListAr []string
	if *ignoreList != "" {
		ignoreListAr = strings.Split(*ignoreList, ",")
	}

	if err := sh.Run(configuration.Options{VCS: "git", TempVCSPath: "./.gravitywell", APIVersion: "v1", SSHKeyPath: *sshkeypath,
		DryRun: *dryRun, TryUpdate: *tryUpdate, Redeploy: *redeploy, IgnoreList: ignoreListAr}); err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

}
