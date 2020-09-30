// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"github.com/urfave/cli/v2"
)

// CmdMilestonesDelete represents a sub command of milestones to delete an milestone
var CmdMilestonesDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "delete a milestone",
	Description: "delete a milestone",
	ArgsUsage:   "<milestone name>",
	Action:      deleteMilestone,
	Flags:       flags.AllDefaultFlags,
}

func deleteMilestone(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	_, err := client.DeleteMilestoneByName(owner, repo, ctx.Args().First())
	return err
}
