// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesCreate represents a sub command of issues to create issue
var CmdIssuesCreate = cli.Command{
	Name:        "create",
	Usage:       "Create an issue on repository",
	Description: `Create an issue on repository`,
	Action:      runIssuesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "issue title to create",
		},
		&cli.StringFlag{
			Name:    "body",
			Aliases: []string{"b"},
			Usage:   "issue body to create",
		},
	}, flags.LoginRepoFlags...),
}

func runIssuesCreate(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	_, _, err := login.Client().CreateIssue(owner, repo, gitea.CreateIssueOption{
		Title: ctx.String("title"),
		Body:  ctx.String("body"),
		// TODO:
		//Assignee  string   `json:"assignee"`
		//Assignees []string `json:"assignees"`
		//Deadline *time.Time `json:"due_date"`
		//Milestone int64 `json:"milestone"`
		//Labels []int64 `json:"labels"`
		//Closed bool    `json:"closed"`
	})

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Print IssueDetails

	return nil
}
