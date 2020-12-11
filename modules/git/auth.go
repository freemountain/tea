// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"io/ioutil"
	"net/url"

	"code.gitea.io/tea/modules/utils"

	git_transport "github.com/go-git/go-git/v5/plumbing/transport"
	gogit_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	gogit_ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
)

type pwCallback = func(string) (string, error)

// GetAuthForURL returns the appropriate AuthMethod to be used in Push() / Pull()
// operations depending on the protocol, and prompts the user for credentials if
// necessary.
func GetAuthForURL(remoteURL *url.URL, authToken, keyFile string, passwordCallback pwCallback) (git_transport.AuthMethod, error) {
	switch remoteURL.Scheme {
	case "http", "https":
		// gitea supports push/pull via app token as username.
		return &gogit_http.BasicAuth{Password: "", Username: authToken}, nil

	case "ssh":
		// try to select right key via ssh-agent. if it fails, try to read a key manually
		user := remoteURL.User.Username()
		auth, err := gogit_ssh.DefaultAuthBuilder(user)
		if err != nil {
			signer, err2 := readSSHPrivKey(keyFile, passwordCallback)
			if err2 != nil {
				return nil, err2
			}
			auth = &gogit_ssh.PublicKeys{User: user, Signer: signer}
		}
		return auth, nil
	}
	return nil, fmt.Errorf("don't know how to handle url scheme %v", remoteURL.Scheme)
}

func readSSHPrivKey(keyFile string, passwordCallback pwCallback) (sig ssh.Signer, err error) {
	if keyFile != "" {
		keyFile, err = utils.AbsPathWithExpansion(keyFile)
	} else {
		keyFile, err = utils.AbsPathWithExpansion("~/.ssh/id_rsa")
	}
	if err != nil {
		return nil, err
	}
	sshKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	sig, err = ssh.ParsePrivateKey(sshKey)
	if _, ok := err.(*ssh.PassphraseMissingError); ok && passwordCallback != nil {
		// allow for up to 3 password attempts
		for i := 0; i < 3; i++ {
			var pass string
			pass, err = passwordCallback(keyFile)
			if err != nil {
				return nil, err
			}
			sig, err = ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(pass))
			if err == nil {
				break
			}
		}
	}
	return sig, err
}
