// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdPullsClean removes the remote and local feature branches, if a PR is merged.
var CmdPullsClean = cli.Command{
	Name:        "clean",
	Usage:       "Deletes local & remote feature-branches for a closed pull request",
	Description: `Deletes local & remote feature-branches for a closed pull request`,
	ArgsUsage:   "<pull index>",
	Action:      runPullsClean,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:  "ignore-sha",
			Usage: "Find the local branch by name instead of commit hash (less precise)",
		},
	}, flags.AllDefaultFlags...),
}

func runPullsClean(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if ctx.Args().Len() != 1 {
		return fmt.Errorf("Must specify a PR index")
	}

	idx, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	return task.PullClean(login, owner, repo, idx, ctx.Bool("ignore-sha"), interact.PromptPassword)
}
