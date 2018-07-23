package scheduler

import (
	"bufio"
	"fmt"
	"os/exec"

	log "github.com/Sirupsen/logrus"
)

//ShellCommand ...
func ShellCommand(command string, path string, validated bool) error {
	cmd := exec.Command("bash", "-c", command)
	if path != "" {
		cmd.Dir = path
	}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	scanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
	go func() {
		for errScanner.Scan() {
			fmt.Printf("%s\n", errScanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	} else {
		if validated {
			log.Debug("Successful")
		}
	}
	return nil
}
