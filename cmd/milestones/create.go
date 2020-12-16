// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestonesCreate represents a sub command of milestones to create milestone
var CmdMilestonesCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create an milestone on repository",
	Description: `Create an milestone on repository`,
	Action:      runMilestonesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "milestone title to create",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "milestone description to create",
		},
		&cli.StringFlag{
			Name:        "state",
			Usage:       "set milestone state (default is open)",
			DefaultText: "open",
		},
	}, flags.AllDefaultFlags...),
}

func runMilestonesCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	title := ctx.String("title")
	if len(title) == 0 {
		fmt.Printf("Title is required\n")
		return nil
	}

	state := gitea.StateOpen
	if ctx.String("state") == "closed" {
		state = gitea.StateClosed
	}

	mile, _, err := ctx.Login.Client().CreateMilestone(ctx.Owner, ctx.Repo, gitea.CreateMilestoneOption{
		Title:       title,
		Description: ctx.String("description"),
		State:       state,
	})
	if err != nil {
		return err
	}

	print.MilestoneDetails(mile)
	return nil
}
