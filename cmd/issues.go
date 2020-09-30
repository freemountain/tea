// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/issues"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdIssues represents to login a gitea server.
var CmdIssues = cli.Command{
	Name:        "issues",
	Usage:       "List, create and update issues",
	Description: "List, create and update issues",
	ArgsUsage:   "[<issue index>]",
	Action:      runIssues,
	Subcommands: []*cli.Command{
		&issues.CmdIssuesList,
		&issues.CmdIssuesCreate,
		&issues.CmdIssuesReopen,
		&issues.CmdIssuesClose,
	},
	Flags: flags.IssuePRFlags,
}

func runIssues(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runIssueDetail(ctx.Args().First())
	}
	return issues.RunIssuesList(ctx)
}

func runIssueDetail(index string) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	idx, err := utils.ArgToIndex(index)
	if err != nil {
		return err
	}
	issue, _, err := login.Client().GetIssue(owner, repo, idx)
	if err != nil {
		return err
	}
	print.IssueDetails(issue)
	return nil
}
