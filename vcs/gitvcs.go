package vcs

import (
	"os"

	git "gopkg.in/src-d/go-git.v4"
)

//GitVCS ...
type GitVCS struct {
}

//Fetch ...
func (g *GitVCS) Fetch(localpath string, remote string) (string, error) {

	_, err := git.PlainClone(localpath, false, &git.CloneOptions{
		URL:      remote,
		Progress: os.Stdout,
	})
	return "", err
}

//Update ...
func (g *GitVCS) Update(localpath string) (string, error) {
	return g.Update(localpath)
}
