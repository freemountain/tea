// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"

	"github.com/urfave/cli/v2"
)

// CmdPullsCreate creates a pull request
var CmdPullsCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a pull-request",
	Description: "Create a pull-request",
	Action:      runPullsCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "head",
			Usage: "Set head branch (default is current one)",
		},
		&cli.StringFlag{
			Name:    "base",
			Aliases: []string{"b"},
			Usage:   "Set base branch (default is default branch)",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Set title of pull (default is head branch name)",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "Set body of new pull",
		},
	}, flags.AllDefaultFlags...),
}

func runPullsCreate(ctx *cli.Context) error {
	login, ownerArg, repoArg := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	// no args -> interactive mode
	if ctx.NumFlags() == 0 {
		return interact.CreatePull(login, ownerArg, repoArg)
	}

	// else use args to create PR
	return task.CreatePull(
		login,
		ownerArg,
		repoArg,
		ctx.String("base"),
		ctx.String("head"),
		ctx.String("title"),
		ctx.String("description"),
	)
}
