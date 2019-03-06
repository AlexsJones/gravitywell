package vcs

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlexsJones/gravitywell/configuration"
	log "github.com/Sirupsen/logrus"
)

func nameForRepository(repoUrl string) string {
	extension := filepath.Ext(repoUrl)
	remoteVCSRepoName := repoUrl[0 : len(repoUrl)-len(extension)]
	splitStrings := strings.Split(remoteVCSRepoName, "/")
	return splitStrings[len(splitStrings)-1]
}

func FetchRepo(remote string, opt configuration.Options, branch string) (string, error) {
	remoteVCSRepoName := nameForRepository(remote)
	if _, err := os.Stat(path.Join(opt.TempVCSPath, remoteVCSRepoName)); !os.IsNotExist(err) {
		log.Debug(fmt.Sprintf("Using existing repository %s", path.Join(opt.TempVCSPath, remoteVCSRepoName)))
		return remoteVCSRepoName, nil
	}
	log.Debug(fmt.Sprintf("Fetching deployment %s into %s\n", remoteVCSRepoName, path.Join(opt.TempVCSPath, remoteVCSRepoName)))
	gvcs := new(GitVCS)
	_, err := gvcs.Fetch(path.Join(opt.TempVCSPath, remoteVCSRepoName), remote, opt.SSHKeyPath, branch)
	return remoteVCSRepoName, err
}
