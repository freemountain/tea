// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	"github.com/muesli/termenv"
	"gopkg.in/yaml.v2"
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

// Config reprensents local configurations
type Config struct {
	Logins []Login `yaml:"logins"`
}

var (
	config         Config
	yamlConfigPath string
)

func init() {
	homeDir, err := utils.Home()
	if err != nil {
		log.Fatal("Retrieve home dir failed")
	}

	dir := filepath.Join(homeDir, ".tea")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal("Init tea config dir " + dir + " failed")
	}

	yamlConfigPath = filepath.Join(dir, "tea.yml")
}

func getGlamourTheme() string {
	if termenv.HasDarkBackground() {
		return "dark"
	}
	return "light"
}

func getOwnerAndRepo(repoPath, user string) (string, string) {
	if len(repoPath) == 0 {
		return "", ""
	}
	p := strings.Split(repoPath, "/")
	if len(p) >= 2 {
		return p[0], p[1]
	}
	return user, repoPath
}

func getDefaultLogin() (*Login, error) {
	if len(config.Logins) == 0 {
		return nil, errors.New("No available login")
	}
	for _, l := range config.Logins {
		if l.Default {
			return &l, nil
		}
	}

	return &config.Logins[0], nil
}

func getLoginByName(name string) *Login {
	for _, l := range config.Logins {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

func addLogin(login Login) error {
	for _, l := range config.Logins {
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
	config.Logins = append(config.Logins, login)

	return nil
}

func isFileExist(fileName string) (bool, error) {
	f, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if f.IsDir() {
		return false, errors.New("A directory with the same name exists")
	}
	return true, nil
}

func loadConfig(ymlPath string) error {
	exist, _ := isFileExist(ymlPath)
	if exist {
		Println("Found config file", ymlPath)
		bs, err := ioutil.ReadFile(ymlPath)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bs, &config)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveConfig(ymlPath string) error {
	bs, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ymlPath, bs, 0660)
}

func curGitRepoPath(path string) (*Login, string, error) {
	var err error
	var repo *git.TeaRepo
	if len(path) == 0 {
		repo, err = git.RepoForWorkdir()
	} else {
		repo, err = git.RepoFromPath(path)
	}
	if err != nil {
		return nil, "", err
	}
	gitConfig, err := repo.Config()
	if err != nil {
		return nil, "", err
	}

	// if no remote
	if len(gitConfig.Remotes) == 0 {
		return nil, "", errors.New("No remote(s) found in this Git repository")
	}

	// if only one remote exists
	if len(gitConfig.Remotes) >= 1 && len(remoteValue) == 0 {
		for remote := range gitConfig.Remotes {
			remoteValue = remote
		}
		if len(gitConfig.Remotes) > 1 {
			// if master branch is present, use it as the default remote
			masterBranch, ok := gitConfig.Branches["master"]
			if ok {
				if len(masterBranch.Remote) > 0 {
					remoteValue = masterBranch.Remote
				}
			}
		}
	}

	remoteConfig, ok := gitConfig.Remotes[remoteValue]
	if !ok || remoteConfig == nil {
		return nil, "", errors.New("Remote " + remoteValue + " not found in this Git repository")
	}

	for _, l := range config.Logins {
		for _, u := range remoteConfig.URLs {
			p, err := git.ParseURL(strings.TrimSpace(u))
			if err != nil {
				return nil, "", fmt.Errorf("Git remote URL parse failed: %s", err.Error())
			}
			if strings.EqualFold(p.Scheme, "http") || strings.EqualFold(p.Scheme, "https") {
				if strings.HasPrefix(u, l.URL) {
					ps := strings.Split(p.Path, "/")
					path := strings.Join(ps[len(ps)-2:], "/")
					return &l, strings.TrimSuffix(path, ".git"), nil
				}
			} else if strings.EqualFold(p.Scheme, "ssh") {
				if l.GetSSHHost() == strings.Split(p.Host, ":")[0] {
					return &l, strings.TrimLeft(strings.TrimSuffix(p.Path, ".git"), "/"), nil
				}
			}
		}
	}

	return nil, "", errors.New("No Gitea login found. You might want to specify --repo (and --login) to work outside of a repository")
}
