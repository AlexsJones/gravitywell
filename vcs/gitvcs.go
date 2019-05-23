package vcs

import (
	logger "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
	"gopkg.in/src-d/go-git.v4"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

//GitVCS ...
type GitVCS struct {
}

//Fetch ...
func (g *GitVCS) Fetch(localpath string, remote string, reference string, keypath string) (string, error) {

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

	repo, err := git.PlainClone(localpath, false, &git.CloneOptions{
		URL:      remote,
		Progress: os.Stdout,
		Auth:     auth,
	})

	if err != nil {
		logger.Fatal(err)
	}

	if reference != "" {
		logger.Infof("Checking out %s:%s",remote,reference)
		tree, err := repo.Worktree()
		if err != nil {
			logger.Fatal(err)
		}
		err = tree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(reference),
		})
		if err != nil {
			logger.Fatal(err)
		}
	}

	return "", err
}

//Update ...
func (g *GitVCS) Update(localpath string) (string, error) {
	return g.Update(localpath)
}
