// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdPullsList represents a sub command of issues to list pulls
var CmdPullsList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List pull requests of the repository",
	Description: `List pull requests of the repository`,
	Action:      RunPullsList,
	Flags:       flags.IssuePRFlags,
}

// RunPullsList return list of pulls
func RunPullsList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	prs, _, err := ctx.Login.Client().ListRepoPullRequests(ctx.Owner, ctx.Repo, gitea.ListPullRequestsOptions{
		State: state,
	})

	if err != nil {
		return err
	}

	print.PullsList(prs, ctx.Output)
	return nil
}
