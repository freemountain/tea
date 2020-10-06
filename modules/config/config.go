// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	"github.com/adrg/xdg"
	"gopkg.in/yaml.v2"
)

// LocalConfig represents local configurations
type LocalConfig struct {
	Logins []Login `yaml:"logins"`
}

var (
	// Config contain if loaded local tea config
	Config LocalConfig
)

// GetConfigPath return path to tea config file
func GetConfigPath() string {
	configFilePath, err := xdg.ConfigFile("tea/config.yml")

	var exists bool
	if err != nil {
		exists = false
	} else {
		exists, _ = utils.PathExists(configFilePath)
	}

	// fallback to old config if no new one exists
	if !exists {
		file := filepath.Join(xdg.Home, ".tea", "tea.yml")
		exists, _ = utils.PathExists(file)
		if exists {
			return file
		}
	}

	if err != nil {
		log.Fatal("unable to get or create config file")
	}

	return configFilePath
}

// LoadConfig load config into global Config var
func LoadConfig() error {
	ymlPath := GetConfigPath()
	exist, _ := utils.FileExist(ymlPath)
	if exist {
		bs, err := ioutil.ReadFile(ymlPath)
		if err != nil {
			return fmt.Errorf("Failed to read config file: %s", ymlPath)
		}

		err = yaml.Unmarshal(bs, &Config)
		if err != nil {
			return fmt.Errorf("Failed to parse contents of config file: %s", ymlPath)
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
