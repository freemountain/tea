// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
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
func RunMilestonesList(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	milestones, _, err := login.Client().ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{
		ListOptions: flags.GetListOptions(ctx),
		State:       state,
	})

	if err != nil {
		log.Fatal(err)
	}

	print.MilestonesList(milestones, flags.GlobalOutputValue, state)
	return nil
}
