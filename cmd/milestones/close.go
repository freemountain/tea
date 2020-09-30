// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"code.gitea.io/tea/cmd/flags"

	"github.com/urfave/cli/v2"
)

// CmdMilestonesClose represents a sub command of milestones to close an milestone
var CmdMilestonesClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of an milestone to 'closed'",
	Description: `Change state of an milestone to 'closed'`,
	ArgsUsage:   "<milestone name>",
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("force") {
			return deleteMilestone(ctx)
		}
		return editMilestoneStatus(ctx, true)
	},
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "delete milestone",
		},
	}, flags.AllDefaultFlags...),
}
