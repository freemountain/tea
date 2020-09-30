// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/repos"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Show repositories details",
	Description: "Show repositories details",
	ArgsUsage:   "[<repo owner>/<repo name>]",
	Action:      runRepos,
	Subcommands: []*cli.Command{
		&repos.CmdReposList,
		&repos.CmdRepoCreate,
	},
	Flags: flags.LoginOutputFlags,
}

func runRepos(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runRepoDetail(ctx.Args().First())
	}
	return repos.RunReposList(ctx)
}

func runRepoDetail(path string) error {
	login := config.InitCommandLoginOnly(flags.GlobalLoginValue)
	client := login.Client()
	repoOwner, repoName := config.GetOwnerAndRepo(path, login.User)
	repo, _, err := client.GetRepo(repoOwner, repoName)
	if err != nil {
		return err
	}
	topics, _, err := client.ListRepoTopics(repo.Owner.UserName, repo.Name, gitea.ListRepoTopicsOptions{})
	if err != nil {
		return err
	}

	print.RepoDetails(repo, topics)
	return nil
}
