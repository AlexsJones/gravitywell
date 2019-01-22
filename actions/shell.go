package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/shell"
	log "github.com/Sirupsen/logrus"
	"path"
)

func ExecuteShellAction(action configuration.Action, opt configuration.Options, repoName string) {
	command, ok := action.Execute.Configuration["Command"]
	if !ok {
		log.Warn("Could not run the shell step as Command could not be found")
		return
	}

	p := path.Join(opt.TempVCSPath, repoName)

	tp, ok := action.Execute.Configuration["Path"]
	if ok {
		p = tp
	}

	log.Warn(fmt.Sprintf("Running shell command %s\n", command))
	if err := shell.ShellCommand(command, p, true); err != nil {
		log.Error(err.Error())
	}
}
