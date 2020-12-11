// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"
	"strings"

	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"golang.org/x/crypto/ssh"
)

// Login represents a login to a gitea server, you even could add multiple logins for one gitea server
type Login struct {
	Name    string `yaml:"name"`
	URL     string `yaml:"url"`
	Token   string `yaml:"token"`
	Default bool   `yaml:"default"`
	SSHHost string `yaml:"ssh_host"`
	// optional path to the private key
	SSHKey   string `yaml:"ssh_key"`
	Insecure bool   `yaml:"insecure"`
	// User is username from gitea
	User string `yaml:"user"`
	// Created is auto created unix timestamp
	Created int64 `yaml:"created"`
}

// GetDefaultLogin return the default login
func GetDefaultLogin() (*Login, error) {
	if len(Config.Logins) == 0 {
		return nil, errors.New("No available login")
	}
	for _, l := range Config.Logins {
		if l.Default {
			return &l, nil
		}
	}

	return &Config.Logins[0], nil
}

// GetLoginByName get login by name
func GetLoginByName(name string) *Login {
	for _, l := range Config.Logins {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

// GenerateLoginName generates a name string based on instance URL & adds username if the result is not unique
func GenerateLoginName(url, user string) (string, error) {
	parsedURL, err := utils.NormalizeURL(url)
	if err != nil {
		return "", err
	}
	name := parsedURL.Host

	// append user name if login name already exists
	if len(user) != 0 {
		for _, l := range Config.Logins {
			if l.Name == name {
				name += "_" + user
				break
			}
		}
	}

	return name, nil
}

// DeleteLogin delete a login by name
func DeleteLogin(name string) error {
	var idx = -1
	for i, l := range Config.Logins {
		if l.Name == name {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("can not delete login '%s', does not exist", name)
	}

	Config.Logins = append(Config.Logins[:idx], Config.Logins[idx+1:]...)

	return SaveConfig()
}

// Client returns a client to operate Gitea API
func (l *Login) Client() *gitea.Client {
	httpClient := &http.Client{}
	if l.Insecure {
		cookieJar, _ := cookiejar.New(nil)

		httpClient = &http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}}
	}

	client, err := gitea.NewClient(l.URL,
		gitea.SetToken(l.Token),
		gitea.SetHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// GetSSHHost returns SSH host name
func (l *Login) GetSSHHost() string {
	if l.SSHHost != "" {
		return l.SSHHost
	}

	u, err := url.Parse(l.URL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// FindSSHKey retrieves the ssh keys registered in gitea, and tries to find
// a matching private key in ~/.ssh/. If no match is found, path is empty.
func (l *Login) FindSSHKey() (string, error) {
	// get keys registered on gitea instance
	keys, _, err := l.Client().ListMyPublicKeys(gitea.ListPublicKeysOptions{})
	if err != nil || len(keys) == 0 {
		return "", err
	}

	// enumerate ~/.ssh/*.pub files
	glob, err := utils.AbsPathWithExpansion("~/.ssh/*.pub")
	if err != nil {
		return "", err
	}
	localPubkeyPaths, err := filepath.Glob(glob)
	if err != nil {
		return "", err
	}

	// parse each local key with present privkey & compare fingerprints to online keys
	for _, pubkeyPath := range localPubkeyPaths {
		var pubkeyFile []byte
		pubkeyFile, err = ioutil.ReadFile(pubkeyPath)
		if err != nil {
			continue
		}
		fields := strings.Split(string(pubkeyFile), " ")
		if len(fields) < 2 { // first word is key type, second word is key material
			continue
		}

		var keymaterial []byte
		keymaterial, err = base64.StdEncoding.DecodeString(fields[1])
		if err != nil {
			continue
		}

		var pubkey ssh.PublicKey
		pubkey, err = ssh.ParsePublicKey(keymaterial)
		if err != nil {
			continue
		}

		privkeyPath := strings.TrimSuffix(pubkeyPath, ".pub")
		var exists bool
		exists, err = utils.FileExist(privkeyPath)
		if err != nil || !exists {
			continue
		}

		// if pubkey fingerprints match, return path to corresponding privkey.
		fingerprint := ssh.FingerprintSHA256(pubkey)
		for _, key := range keys {
			if fingerprint == key.Fingerprint {
				return privkeyPath, nil
			}
		}
	}

	return "", err
}
