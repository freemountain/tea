// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"github.com/urfave/cli"
)

// create global variables for global Flags to simplify
// access to the options without requiring cli.Context
var (
	loginValue  string
	repoValue   string
	outputValue string
	remoteValue string
)

// LoginFlag provides flag to specify tea login profile
var LoginFlag = cli.StringFlag{
	Name:        "login, l",
	Usage:       "Use a different Gitea login. Optional",
	Destination: &loginValue,
}

// RepoFlag provides flag to specify repository
var RepoFlag = cli.StringFlag{
	Name:        "repo, r",
	Usage:       "Repository to interact with. Optional",
	Destination: &repoValue,
}

// RemoteFlag provides flag to specify remote repository
var RemoteFlag = cli.StringFlag{
	Name:        "remote, R",
	Usage:       "Discover Gitea login from remote. Optional",
	Destination: &remoteValue,
}

// OutputFlag provides flag to specify output type
var OutputFlag = cli.StringFlag{
	Name:        "output, o",
	Usage:       "Output format. (csv, simple, table, tsv, yaml)",
	Destination: &outputValue,
}

// LoginOutputFlags defines login and output flags that should
// added to all subcommands and appended to the flags of the
// subcommand to work around issue and provide --login and --output:
// https://github.com/urfave/cli/issues/585
var LoginOutputFlags = []cli.Flag{
	LoginFlag,
	OutputFlag,
}

// LoginRepoFlags defines login and repo flags that should
// be used for all subcommands and appended to the flags of
// the subcommand to work around issue and provide --login and --repo:
// https://github.com/urfave/cli/issues/585
var LoginRepoFlags = []cli.Flag{
	LoginFlag,
	RepoFlag,
	RemoteFlag,
}

// AllDefaultFlags defines flags that should be available
// for all subcommands working with dedicated repositories
// to work around issue and provide --login, --repo and --output:
// https://github.com/urfave/cli/issues/585
var AllDefaultFlags = append([]cli.Flag{
	RepoFlag,
	RemoteFlag,
}, LoginOutputFlags...)

// initCommand returns repository and *Login based on flags
func initCommand() (*Login, string, string) {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("Unable to load config file " + yamlConfigPath)
	}

	var login *Login
	if loginValue == "" {
		login, err = getActiveLogin()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		login = getLoginByName(loginValue)
		if login == nil {
			log.Fatal("Login name " + loginValue + " does not exist")
		}
	}

	repoPath := repoValue
	if repoPath == "" {
		login, repoPath, err = curGitRepoPath()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	owner, repo := splitRepo(repoPath)
	return login, owner, repo
}

// initCommandLoginOnly return *Login based on flags
func initCommandLoginOnly() *Login {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed ", yamlConfigPath)
	}

	var login *Login

	login = getLoginByName(loginValue)
	if login == nil {
		log.Fatal("indicated login name ", loginValue, " does not exist")
	}

	return login
}
