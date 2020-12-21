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

// CmdPullsReject requests changes to a PR
var CmdPullsReject = cli.Command{
	Name:        "reject",
	Usage:       "Request changes to a pull request",
	Description: "Request changes to a pull request",
	ArgsUsage:   "<pull index> <reason>",
	Action: func(cmd *cli.Context) error {
		ctx := context.InitCommand(cmd)
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

		if ctx.Args().Len() < 2 {
			return fmt.Errorf("Must specify a PR index and comment")
		}

		idx, err := utils.ArgToIndex(ctx.Args().First())
		if err != nil {
			return err
		}

		comment := strings.Join(ctx.Args().Tail(), " ")

		return task.CreatePullReview(ctx, idx, gitea.ReviewStateRequestChanges, comment, nil)
	},
	Flags: flags.AllDefaultFlags,
}
