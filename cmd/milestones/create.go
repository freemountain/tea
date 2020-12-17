// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"time"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/araddon/dateparse"
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
			Name:    "deadline",
			Aliases: []string{"expires", "x"},
			Usage:   "set milestone deadline (default is no due date)",
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

	date := ctx.String("deadline")
	deadline := &time.Time{}
	if date != "" {
		t, err := dateparse.ParseAny(date)
		if err == nil {
			return err
		}
		deadline = &t
	}

	state := gitea.StateOpen
	if ctx.String("state") == "closed" {
		state = gitea.StateClosed
	}

	if ctx.NumFlags() == 0 {
		return interact.CreateMilestone(ctx.Login, ctx.Owner, ctx.Repo)
	}

	return task.CreateMilestone(
		ctx.Login,
		ctx.Owner,
		ctx.Repo,
		ctx.String("title"),
		ctx.String("description"),
		deadline,
		state,
	)
}
