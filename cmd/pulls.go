// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/pulls"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdPulls is the main command to operate on PRs
var CmdPulls = cli.Command{
	Name:        "pulls",
	Aliases:     []string{"pull", "pr"},
	Usage:       "List, create, checkout and clean pull requests",
	Description: `List, create, checkout and clean pull requests`,
	ArgsUsage:   "[<pull index>]",
	Action:      runPulls,
	Flags:       flags.IssuePRFlags,
	Subcommands: []*cli.Command{
		&pulls.CmdPullsList,
		&pulls.CmdPullsCheckout,
		&pulls.CmdPullsClean,
		&pulls.CmdPullsCreate,
	},
}

func runPulls(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runPullDetail(ctx, ctx.Args().First())
	}
	return pulls.RunPullsList(ctx)
}

func runPullDetail(cmd *cli.Context, index string) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	idx, err := utils.ArgToIndex(index)
	if err != nil {
		return err
	}

	client := ctx.Login.Client()
	pr, _, err := client.GetPullRequest(ctx.Owner, ctx.Repo, idx)
	if err != nil {
		return err
	}

	reviews, _, err := client.ListPullReviews(ctx.Owner, ctx.Repo, idx, gitea.ListPullReviewsOptions{})
	if err != nil {
		fmt.Printf("error while loading reviews: %v\n", err)
	}

	ci, _, err := client.GetCombinedStatus(ctx.Owner, ctx.Repo, pr.Head.Sha)
	if err != nil {
		fmt.Printf("error while loading CI: %v\n", err)
	}

	print.PullDetails(pr, reviews, ci)
	return nil
}
