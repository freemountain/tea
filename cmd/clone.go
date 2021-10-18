// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdRepoClone represents a sub command of repos to create a local copy
var CmdRepoClone = cli.Command{
	Name:    "clone",
	Aliases: []string{"C"},
	Usage:   "Clone a repository locally",
	Description: `Clone a repository locally, without a local git installation required.
The repo slug can be specified in different formats:
	gitea/tea
	tea
	gitea.com/gitea/tea
	git@gitea.com:gitea/tea
	https://gitea.com/gitea/tea
	ssh://gitea.com:22/gitea/tea
When a host is specified in the repo-slug, it will override the login specified with --login.
	`,
	Category:  catHelpers,
	Action:    runRepoClone,
	ArgsUsage: "<repo-slug> [target dir]",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "depth",
			Aliases: []string{"d"},
			Usage:   "num commits to fetch, defaults to all",
		},
		&flags.LoginFlag,
	},
}

func runRepoClone(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)

	args := ctx.Args()
	if args.Len() < 1 {
		return cli.ShowCommandHelp(cmd, "clone")
	}
	dir := args.Get(1)

	var (
		login *config.Login = ctx.Login
		owner string        = ctx.Login.User
		repo  string
	)

	// parse first arg as repo specifier
	repoSlug := args.Get(0)
	url, err := git.ParseURL(repoSlug)
	if err != nil {
		return err
	}

	owner, repo = utils.GetOwnerAndRepo(url.Path, login.User)
	if url.Host != "" {
		login = config.GetLoginByHost(url.Host)
		if login == nil {
			return fmt.Errorf("No login configured matching host '%s', run `tea login add` first", url.Host)
		}
	}

	_, err = task.RepoClone(
		dir,
		login,
		owner,
		repo,
		interact.PromptPassword,
		ctx.Int("depth"),
	)

	return err
}
