// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/milestones"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdMilestones represents to operate repositories milestones.
var CmdMilestones = cli.Command{
	Name:        "milestones",
	Aliases:     []string{"milestone", "ms"},
	Category:    catEntities,
	Usage:       "List and create milestones",
	Description: `List and create milestones`,
	ArgsUsage:   "[<milestone name>]",
	Action:      runMilestones,
	Subcommands: []*cli.Command{
		&milestones.CmdMilestonesList,
		&milestones.CmdMilestonesCreate,
		&milestones.CmdMilestonesClose,
		&milestones.CmdMilestonesDelete,
		&milestones.CmdMilestonesReopen,
		&milestones.CmdMilestonesIssues,
	},
	Flags: flags.AllDefaultFlags,
}

func runMilestones(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runMilestoneDetail(ctx, ctx.Args().First())
	}
	return milestones.RunMilestonesList(ctx)
}

func runMilestoneDetail(cmd *cli.Context, name string) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	client := ctx.Login.Client()

	milestone, _, err := client.GetMilestoneByName(ctx.Owner, ctx.Repo, name)
	if err != nil {
		return err
	}

	print.MilestoneDetails(milestone)
	return nil
}
