// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdPullsApprove approves a PR
var CmdPullsApprove = cli.Command{
	Name:        "approve",
	Aliases:     []string{"lgtm", "a"},
	Usage:       "Approve a pull request",
	Description: "Approve a pull request",
	ArgsUsage:   "<pull index> [<comment>]",
	Action: func(cmd *cli.Context) error {
		ctx := context.InitCommand(cmd)
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

		if ctx.Args().Len() == 0 {
			return fmt.Errorf("Must specify a PR index")
		}

		idx, err := utils.ArgToIndex(ctx.Args().First())
		if err != nil {
			return err
		}

		comment := strings.Join(ctx.Args().Tail(), " ")

		return task.CreatePullReview(ctx, idx, gitea.ReviewStateApproved, comment, nil)
	},
	Flags: flags.AllDefaultFlags,
}
