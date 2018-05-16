package vcs

import "errors"

type GitVCS struct {
}

func (g *GitVCS) Fetch(v IVCS, remote string) (string, error) {

	return "", errors.New("Fetch failure")
}
