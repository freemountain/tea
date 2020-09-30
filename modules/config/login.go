// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"

	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
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
	// optional gitea username
	User string `yaml:"user"`
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

// AddLogin add login to config ( global var & file)
func AddLogin(name, token, user, passwd, sshKey, giteaURL string, insecure bool) error {

	if len(giteaURL) == 0 {
		log.Fatal("You have to input Gitea server URL")
	}
	if len(token) == 0 && (len(user)+len(passwd)) == 0 {
		log.Fatal("No token set")
	} else if len(user) != 0 && len(passwd) == 0 {
		log.Fatal("No password set")
	} else if len(user) == 0 && len(passwd) != 0 {
		log.Fatal("No user set")
	}

	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{}
	if insecure {
		cookieJar, _ := cookiejar.New(nil)
		httpClient = &http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}}
	}
	client, err := gitea.NewClient(giteaURL,
		gitea.SetToken(token),
		gitea.SetBasicAuth(user, passwd),
		gitea.SetHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatal(err)
	}

	u, _, err := client.GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}

	if len(token) == 0 {
		// create token
		host, _ := os.Hostname()
		tl, _, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{})
		if err != nil {
			return err
		}
		tokenName := host + "-tea"
		for i := range tl {
			if tl[i].Name == tokenName {
				tokenName += time.Now().Format("2006-01-02_15-04-05")
				break
			}
		}
		t, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{Name: tokenName})
		if err != nil {
			return err
		}
		token = t.Token
	}

	fmt.Println("Login successful! Login name " + u.UserName)

	if len(name) == 0 {
		parsedURL, err := url.Parse(giteaURL)
		if err != nil {
			return err
		}
		name = strings.ReplaceAll(strings.Title(parsedURL.Host), ".", "")
		for _, l := range Config.Logins {
			if l.Name == name {
				name += "_" + u.UserName
				break
			}
		}
	}

	err = addLoginToConfig(Login{
		Name:     name,
		URL:      giteaURL,
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
		User:     u.UserName,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = SaveConfig()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// addLoginToConfig add a login to global Config var
func addLoginToConfig(login Login) error {
	for _, l := range Config.Logins {
		if l.Name == login.Name {
			if l.URL == login.URL && l.Token == login.Token {
				return nil
			}
			return errors.New("Login name has already been used")
		}
		if l.URL == login.URL && l.Token == login.Token {
			return errors.New("URL has been added")
		}
	}

	u, err := url.Parse(login.URL)
	if err != nil {
		return err
	}

	if login.SSHHost == "" {
		login.SSHHost = u.Hostname()
	}
	Config.Logins = append(Config.Logins, login)

	return nil
}

// InitCommand returns repository and *Login based on flags
func InitCommand(repoValue, loginValue, remoteValue string) (*Login, string, string) {
	var login *Login

	err := LoadConfig()
	if err != nil {
		log.Fatal("load config file failed ", yamlConfigPath)
	}

	if login, err = GetDefaultLogin(); err != nil {
		log.Fatal(err.Error())
	}

	exist, err := utils.PathExists(repoValue)
	if err != nil {
		log.Fatal(err.Error())
	}

	if exist || len(repoValue) == 0 {
		login, repoValue, err = curGitRepoPath(repoValue, remoteValue)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if loginValue != "" {
		login = GetLoginByName(loginValue)
		if login == nil {
			log.Fatal("Login name " + loginValue + " does not exist")
		}
	}

	owner, repo := GetOwnerAndRepo(repoValue, login.User)
	return login, owner, repo
}

// InitCommandLoginOnly return *Login based on flags
func InitCommandLoginOnly(loginValue string) *Login {
	err := LoadConfig()
	if err != nil {
		log.Fatal("load config file failed ", yamlConfigPath)
	}

	var login *Login
	if loginValue == "" {
		login, err = GetDefaultLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = GetLoginByName(loginValue)
		if login == nil {
			log.Fatal("Login name " + loginValue + " does not exist")
		}
	}

	return login
}
