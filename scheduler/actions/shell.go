package actions

import (
	"fmt"
	"github.com/AlexsJones/gravitywell/configuration"
	"github.com/AlexsJones/gravitywell/kinds"
	"github.com/AlexsJones/gravitywell/scheduler/actions/shell"
	"github.com/google/logger"
	"path"
)

func ExecuteShellAction(action kinds.Execute, opt configuration.Options, repoName string) {
	command, ok := action.Configuration["Command"]
	if !ok {
		logger.Warning("Could not run the shell step as Command could not be found")
		return
	}

	p := path.Join(opt.TempVCSPath, repoName)

	tp, ok := action.Configuration["Path"]
	if ok {
		p = tp
	}

	logger.Warning(fmt.Sprintf("Running shell command %s\n", command))
	if err := shell.ShellCommand(command, p, true); err != nil {
		logger.Error(err.Error())
	}
}
