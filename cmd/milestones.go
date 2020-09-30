// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/milestones"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdMilestones represents to operate repositories milestones.
var CmdMilestones = cli.Command{
	Name:        "milestones",
	Aliases:     []string{"ms", "mile"},
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
		return runMilestoneDetail(ctx.Args().First())
	}
	return milestones.RunMilestonesList(ctx)
}

func runMilestoneDetail(name string) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	milestone, _, err := client.GetMilestoneByName(owner, repo, name)
	if err != nil {
		return err
	}

	print.MilestoneDetails(milestone)
	return nil
}
