// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      RunIssuesList,
	Flags:       flags.IssuePRFlags,
}

// RunIssuesList list issues
func RunIssuesList(ctx *cli.Context) error {
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

	issues, _, err := login.Client().ListRepoIssues(owner, repo, gitea.ListIssueOption{
		ListOptions: flags.GetListOptions(ctx),
		State:       state,
		Type:        gitea.IssueTypeIssue,
	})

	if err != nil {
		log.Fatal(err)
	}

	print.IssuesList(issues, flags.GlobalOutputValue)
	return nil
}
