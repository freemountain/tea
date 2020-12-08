// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/pulls"
	"code.gitea.io/tea/modules/config"
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
		return runPullDetail(ctx.Args().First())
	}
	return pulls.RunPullsList(ctx)
}

func runPullDetail(index string) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	idx, err := utils.ArgToIndex(index)
	if err != nil {
		return err
	}

	client := login.Client()
	pr, _, err := client.GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}

	reviews, _, err := client.ListPullReviews(owner, repo, idx, gitea.ListPullReviewsOptions{})
	if err != nil {
		fmt.Printf("error while loading reviews: %v\n", err)
	}

	print.PullDetails(pr, reviews)
	return nil
}
