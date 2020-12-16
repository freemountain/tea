// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"

	"github.com/urfave/cli/v2"
)

// CmdIssuesCreate represents a sub command of issues to create issue
var CmdIssuesCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create an issue on repository",
	Description: `Create an issue on repository`,
	Action:      runIssuesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "issue title to create",
		},
		&cli.StringFlag{
			Name:    "body",
			Aliases: []string{"b"},
			Usage:   "issue body to create",
		},
	}, flags.LoginRepoFlags...),
}

func runIssuesCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	if ctx.NumFlags() == 0 {
		return interact.CreateIssue(ctx.Login, ctx.Owner, ctx.Repo)
	}

	return task.CreateIssue(
		ctx.Login,
		ctx.Owner,
		ctx.Repo,
		ctx.String("title"),
		ctx.String("body"),
	)
}
