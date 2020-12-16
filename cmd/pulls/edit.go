// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"

	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// editPullState abstracts the arg parsing to edit the given pull request
func editPullState(cmd *cli.Context, opts gitea.EditPullRequestOption) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("Please provide a Pull Request index")
	}

	index, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	pr, _, err := ctx.Login.Client().EditPullRequest(ctx.Owner, ctx.Repo, index, opts)
	if err != nil {
		return err
	}

	print.PullDetails(pr, nil, nil)
	return nil
}
