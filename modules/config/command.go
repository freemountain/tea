// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	gogit "github.com/go-git/go-git/v5"
)

// InitCommand resolves the application context, and returns the active login, and if
// available the repo slug. It does this by reading the config file for logins, parsing
// the remotes of the .git repo specified in repoFlag or $PWD, and using overrides from
// command flags. If a local git repo can't be found, repo slug values are unset.
func InitCommand(repoFlag, loginFlag, remoteFlag string) (login *Login, owner string, reponame string) {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	var repoSlug string
	var repoPath string // empty means PWD
	var repoFlagPathExists bool

	// check if repoFlag can be interpreted as path to local repo.
	if len(repoFlag) != 0 {
		repoFlagPathExists, err = utils.PathExists(repoFlag)
		if err != nil {
			log.Fatal(err.Error())
		}
		if repoFlagPathExists {
			repoPath = repoFlag
		}
	}

	// try to read git repo & extract context, ignoring if PWD is not a repo
	login, repoSlug, err = contextFromLocalRepo(repoPath, remoteFlag)
	if err != nil && err != gogit.ErrRepositoryNotExists {
		log.Fatal(err.Error())
	}

	// if repoFlag is not a path, use it to override repoSlug
	if len(repoFlag) != 0 && !repoFlagPathExists {
		repoSlug = repoFlag
	}

	// override login from flag, or use default login if repo based detection failed
	if len(loginFlag) != 0 {
		login = GetLoginByName(loginFlag)
		if login == nil {
			log.Fatalf("Login name '%s' does not exist", loginFlag)
		}
	} else if login == nil {
		if login, err = GetDefaultLogin(); err != nil {
			log.Fatal(err.Error())
		}
	}

	// parse reposlug (owner falling back to login owner if reposlug contains only repo name)
	owner, reponame = utils.GetOwnerAndRepo(repoSlug, login.User)
	return
}

// contextFromLocalRepo discovers login & repo slug from the default branch remote of the given local repo
func contextFromLocalRepo(repoValue, remoteValue string) (*Login, string, error) {
	repo, err := git.RepoFromPath(repoValue)
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
