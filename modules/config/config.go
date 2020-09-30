// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	"gopkg.in/yaml.v2"
)

// LocalConfig represents local configurations
type LocalConfig struct {
	Logins []Login `yaml:"logins"`
}

var (
	// Config contain if loaded local tea config
	Config         LocalConfig
	yamlConfigPath string
)

// TODO: do not use init function to detect the tea configuration, use GetConfigPath()
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

// GetConfigPath return path to tea config file
func GetConfigPath() string {
	return yamlConfigPath
}

// LoadConfig load config into global Config var
func LoadConfig() error {
	ymlPath := GetConfigPath()
	exist, _ := utils.FileExist(ymlPath)
	if exist {
		fmt.Println("Found config file", ymlPath)
		bs, err := ioutil.ReadFile(ymlPath)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bs, &Config)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveConfig save config from global Config var into config file
func SaveConfig() error {
	ymlPath := GetConfigPath()
	bs, err := yaml.Marshal(Config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ymlPath, bs, 0660)
}

func curGitRepoPath(repoValue, remoteValue string) (*Login, string, error) {
	var err error
	var repo *git.TeaRepo
	if len(repoValue) == 0 {
		repo, err = git.RepoForWorkdir()
	} else {
		repo, err = git.RepoFromPath(repoValue)
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

	for _, l := range Config.Logins {
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

// GetOwnerAndRepo return repoOwner and repoName
// based on relative path and default owner (if not in path)
func GetOwnerAndRepo(repoPath, user string) (string, string) {
	if len(repoPath) == 0 {
		return "", ""
	}
	p := strings.Split(repoPath, "/")
	if len(p) >= 2 {
		return p[0], p[1]
	}
	return user, repoPath
}
