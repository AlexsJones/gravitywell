package vcs

import (
	"fmt"
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
func (g *GitVCS) Fetch(localpath string, remote string, keypath string, branch string) (string, error) {

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
		URL:        remote,
		RemoteName: branch,
		Progress:   os.Stdout,
		Auth:       auth,
	})
	return "", err
}

func (g *GitVCS) Add(localpath string, remote string, keypath string, files []string) error {

	repo, err := git.PlainOpen(localpath)
	worktree, err := repo.Worktree()
	for _, file := range files {

		fmt.Sprintf("%s", file)
		_, err = worktree.Add(file)
		if err != nil {
			return err
		}
	}
	return nil
}

//Commit ...
func (g *GitVCS) Commit(localpath string, remote string, keypath string, message string) error {

	repo, err := git.PlainOpen(localpath)
	worktree, err := repo.Worktree()
	_, err = worktree.Commit(message, &git.CommitOptions{})
	// 	Author: &object.Signature{
	// 		Name:  "John Doe",
	// 		Email: "john@doe.org",
	// 		When:  time.Now(),
	// 	},
	// })
	return err
}

//Push ...
func (g *GitVCS) Push(localpath string, remote string, keypath string) error {

	repo, err := git.PlainOpen(localpath)
	err = repo.Push(&git.PushOptions{})
	return err
}

//Update ...
func (g *GitVCS) Update(localpath string) (string, error) {
	return g.Update(localpath)
}
