// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestonesReopen represents a sub command of milestones to open an milestone
var CmdMilestonesReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of an milestone to 'open'",
	Description: `Change state of an milestone to 'open'`,
	ArgsUsage:   "<milestone name>",
	Action: func(ctx *cli.Context) error {
		return editMilestoneStatus(ctx, false)
	},
	Flags: flags.AllDefaultFlags,
}

func editMilestoneStatus(cmd *cli.Context, close bool) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	client := ctx.Login.Client()

	state := gitea.StateOpen
	if close {
		state = gitea.StateClosed
	}
	_, _, err := client.EditMilestoneByName(ctx.Owner, ctx.Repo, ctx.Args().First(), gitea.EditMilestoneOption{
		State: &state,
		Title: ctx.Args().First(),
	})

	return err
}
