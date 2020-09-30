// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/pulls"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

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
	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}

	// TODO: use glamour once #181 is merged
	fmt.Printf("#%d %s\n%s created %s\n\n%s\n", pr.Index,
		pr.Title,
		pr.Poster.UserName,
		pr.Created.Format("2006-01-02 15:04:05"),
		pr.Body,
	)
	return nil
}
