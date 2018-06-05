package vcs

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	git "gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

//GitVCS ...
type GitVCS struct {
}

//Fetch ...
func (g *GitVCS) Fetch(localpath string, remote string, keypath string) (string, error) {

	var p string

	if keypath == "" {
		p = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	} else {
		p = keypath
	}

	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	signer, err := ssh.ParsePrivateKey([]byte(buf))
	if err != nil {
		return "", err
	}
	auth := &gitssh.PublicKeys{
		User:   "git",
		Signer: signer,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
	}

	_, err = git.PlainClone(localpath, false, &git.CloneOptions{
		URL:      remote,
		Progress: os.Stdout,
		Auth:     auth,
	})
	return "", err
}

//Update ...
func (g *GitVCS) Update(localpath string) (string, error) {
	return g.Update(localpath)
}
