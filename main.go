package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/scheduler"
	log "github.com/Sirupsen/logrus"
	"github.com/dimiro1/banner"
)

const ban string = `
{{.AnsiColor.Blue}}  ________                  .__  __                         .__  .__
{{.AnsiColor.Blue}} /  _____/___________ ___  _|__|/  |_ ___.__.__  _  __ ____ |  | |  |
{{.AnsiColor.Blue}}/   \  __\_  __ \__  \\  \/ /  \   __<   |  |\ \/ \/ // __ \|  | |  |
{{.AnsiColor.Blue}}\    \_\  \  | \// __ \\   /|  ||  |  \___  | \     /\  ___/|  |_|  |__
{{.AnsiColor.Blue}} \______  /__|  (____  /\_/ |__||__|  / ____|  \/\_/  \___  >____/____/
{{.AnsiColor.Blue}}        \/           \/               \/                  \/
{{.AnsiColor.Blue}}
{{.AnsiColor.Yellow}} Pull all of your kubernetes cluster configurations into one place.
`

func init() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {

	banner.Init(os.Stdout, true, true, bytes.NewBufferString(ban))
	redeploy := flag.Bool("redeploy", false, "Forces a delete and deploy overriding all kubectl commands. WARNING: Destructive")
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
		log.Warn(fmt.Sprintf("This is a very destructive action, are you sure [Y/N]?: "))
		text, _ := reader.ReadString('\n')
		trimmed := strings.Trim(text, "\n")
		if strings.Compare(trimmed, "Y") != 0 {

			os.Exit(0)
		}
	}

	conf, err := configuration.NewConfiguration(*config)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	sh, err := scheduler.NewScheduler(conf)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	var ignoreListAr []string
	if *ignoreList != "" {
		ignoreListAr = strings.Split(*ignoreList, ",")
	}

	if err := sh.Run(configuration.Options{VCS: "git", TempVCSPath: "./.gravitywell", APIVersion: "v1", SSHKeyPath: *sshkeypath,
		DryRun: *dryRun, TryUpdate: *tryUpdate, Redeploy: *redeploy, IgnoreList: ignoreListAr}); err != nil {
		log.Warn(err.Error())
		os.Exit(1)
	}

}
