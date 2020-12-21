// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdPullsReview starts an interactive review session
var CmdPullsReview = cli.Command{
	Name:        "review",
	Usage:       "Interactively review a pull request",
	Description: "Interactively review a pull request",
	ArgsUsage:   "<pull index>",
	Action: func(cmd *cli.Context) error {
		ctx := context.InitCommand(cmd)
		ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

		if ctx.Args().Len() != 1 {
			return fmt.Errorf("Must specify a PR index")
		}

		idx, err := utils.ArgToIndex(ctx.Args().First())
		if err != nil {
			return err
		}

		return interact.ReviewPull(ctx, idx)
	},
	Flags: flags.AllDefaultFlags,
}
