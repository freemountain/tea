// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"code.gitea.io/tea/modules/utils"

	git_transport "github.com/go-git/go-git/v5/plumbing/transport"
	gogit_http "github.com/go-git/go-git/v5/plumbing/transport/http"
	gogit_ssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// GetAuthForURL returns the appropriate AuthMethod to be used in Push() / Pull()
// operations depending on the protocol, and prompts the user for credentials if
// necessary.
func GetAuthForURL(remoteURL *url.URL, httpUser, keyFile string) (auth git_transport.AuthMethod, err error) {
	user := remoteURL.User.Username()

	switch remoteURL.Scheme {
	case "https":
		if httpUser != "" {
			user = httpUser
		}
		if user == "" {
			user, err = promptUser(remoteURL.Host)
			if err != nil {
				return nil, err
			}
		}
		pass, isSet := remoteURL.User.Password()
		if !isSet {
			pass, err = promptPass(remoteURL.Host)
			if err != nil {
				return nil, err
			}
		}
		auth = &gogit_http.BasicAuth{Password: pass, Username: user}

	case "ssh":
		// try to select right key via ssh-agent. if it fails, try to read a key manually
		auth, err = gogit_ssh.DefaultAuthBuilder(user)
		if err != nil {
			signer, err := readSSHPrivKey(keyFile)
			if err != nil {
				return nil, err
			}
			auth = &gogit_ssh.PublicKeys{User: user, Signer: signer}
		}

	default:
		return nil, fmt.Errorf("don't know how to handle url scheme %v", remoteURL.Scheme)
	}

	return auth, nil
}

func readSSHPrivKey(keyFile string) (sig ssh.Signer, err error) {
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
	if err != nil {
		pass, err := promptPass(keyFile)
		if err != nil {
			return nil, err
		}
		sig, err = ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(pass))
		if err != nil {
			return nil, err
		}
	}
	return sig, err
}

func promptUser(domain string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s username: ", domain)
	username, err := reader.ReadString('\n')
	return strings.TrimSpace(username), err
}

func promptPass(domain string) (string, error) {
	fmt.Printf("%s password: ", domain)
	pass, err := terminal.ReadPassword(0)
	return string(pass), err
}
