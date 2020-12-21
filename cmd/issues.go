// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/issues"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdIssues represents to login a gitea server.
var CmdIssues = cli.Command{
	Name:        "issues",
	Aliases:     []string{"issue", "i"},
	Category:    catEntities,
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
	Flags: issues.CmdIssuesList.Flags,
}

func runIssues(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runIssueDetail(ctx, ctx.Args().First())
	}
	return issues.RunIssuesList(ctx)
}

func runIssueDetail(cmd *cli.Context, index string) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	idx, err := utils.ArgToIndex(index)
	if err != nil {
		return err
	}
	issue, _, err := ctx.Login.Client().GetIssue(ctx.Owner, ctx.Repo, idx)
	if err != nil {
		return err
	}
	print.IssueDetails(issue)
	return nil
}
