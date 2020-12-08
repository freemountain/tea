// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
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
func RunPullsList(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	prs, _, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		State: state,
	})

	if err != nil {
		log.Fatal(err)
	}

	print.PullsList(prs, flags.GlobalOutputValue)
	return nil
}
