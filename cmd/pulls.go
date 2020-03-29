// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strconv"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdPulls represents to login a gitea server.
var CmdPulls = cli.Command{
	Name:        "pulls",
	Usage:       "List open pull requests",
	Description: `List open pull requests`,
	Action:      runPulls,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by PR state (all|open|closed)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(state),
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"State",
		"Author",
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
				string(pr.State),
				name,
				pr.Updated.Format("2006-01-02 15:04:05"),
				pr.Title,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}
