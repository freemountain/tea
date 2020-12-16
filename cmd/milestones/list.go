// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestonesList represents a sub command of milestones to list milestones
var CmdMilestonesList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List milestones of the repository",
	Description: `List milestones of the repository`,
	Action:      RunMilestonesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by milestone state (all|open|closed)",
			DefaultText: "open",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunMilestonesList list milestones
func RunMilestonesList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	client := ctx.Login.Client()
	milestones, _, err := client.ListRepoMilestones(ctx.Owner, ctx.Repo, gitea.ListMilestoneOption{
		ListOptions: ctx.GetListOptions(),
		State:       state,
	})

	if err != nil {
		return err
	}

	print.MilestonesList(milestones, ctx.Output, state)
	return nil
}
