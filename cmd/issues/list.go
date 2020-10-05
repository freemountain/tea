// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"log"
	"strconv"

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

	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	var values [][]string

	if len(issues) == 0 {
		print.OutputList(flags.GlobalOutputValue, headers, values)
		return nil
	}

	for _, issue := range issues {
		author := issue.Poster.FullName
		if len(author) == 0 {
			author = issue.Poster.UserName
		}
		mile := ""
		if issue.Milestone != nil {
			mile = issue.Milestone.Title
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(issue.Index, 10),
				issue.Title,
				string(issue.State),
				author,
				mile,
				print.FormatTime(issue.Updated),
			},
		)
	}
	print.OutputList(flags.GlobalOutputValue, headers, values)

	return nil
}
