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

// GenerateToken creates a new token when given BasicAuth credentials
func (l *Login) GenerateToken(user, pass string) (string, error) {
	client := l.Client()
	gitea.SetBasicAuth(user, pass)(client)

	host, _ := os.Hostname()
	tl, _, err := client.ListAccessTokens(gitea.ListAccessTokensOptions{})
	if err != nil {
		return "", err
	}
	tokenName := host + "-tea"

	for i := range tl {
		if tl[i].Name == tokenName {
			tokenName += time.Now().Format("2006-01-02_15-04-05")
			break
		}
	}

	t, _, err := client.CreateAccessToken(gitea.CreateAccessTokenOption{Name: tokenName})
	return t.Token, err
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

	serverURL, err := utils.NormalizeURL(giteaURL)
	if err != nil {
		log.Fatal("Unable to parse URL", err)
	}

	err = LoadConfig()
	if err != nil {
		log.Fatal("Unable to load config file " + yamlConfigPath)
	}

	login := Login{
		Name:     name,
		URL:      serverURL.String(),
		Token:    token,
		Insecure: insecure,
		SSHKey:   sshKey,
	}

	if len(token) == 0 {
		login.Token, err = login.GenerateToken(user, passwd)
		if err != nil {
			log.Fatal(err)
		}
	}

	u, _, err := login.Client().GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}
	login.User = u.UserName

	if len(login.Name) == 0 {
		login.Name, err = GenerateLoginName(giteaURL, login.User)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = addLoginToConfig(login)
	if err != nil {
		log.Fatal(err)
	}

	err = SaveConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Login as %s on %s successful. Added this login as %s\n", login.User, login.URL, login.Name)

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
			return errors.New("Login for this URL and token already exists")
		}
	}

	if len(login.SSHHost) == 0 {
		u, err := url.Parse(login.URL)
		if err != nil {
			return err
		}
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

	owner, repo := utils.GetOwnerAndRepo(repoValue, login.User)
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
