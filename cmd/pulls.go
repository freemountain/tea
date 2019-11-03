// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strconv"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdPulls represents to login a gitea server.
var CmdPulls = cli.Command{
	Name:        "pulls",
	Usage:       "List open pull requests",
	Description: `List open pull requests`,
	Action:      runPulls,
	Flags:       AllDefaultFlags,
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"Name",
		"Updated",
		"Title",
	}

	var values [][]string

	if len(prs) == 0 {
		Output(outputValue, headers, values)
		return nil
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(pr.Index, 10),
				name,
				pr.Updated.Format("2006-01-02 15:04:05"),
				pr.Title,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}
