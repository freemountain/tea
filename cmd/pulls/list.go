// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"log"
	"strconv"

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

	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	var values [][]string

	if len(prs) == 0 {
		print.OutputList(flags.GlobalOutputValue, headers, values)
		return nil
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		author := pr.Poster.FullName
		if len(author) == 0 {
			author = pr.Poster.UserName
		}
		mile := ""
		if pr.Milestone != nil {
			mile = pr.Milestone.Title
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(pr.Index, 10),
				pr.Title,
				string(pr.State),
				author,
				mile,
				print.FormatTime(*pr.Updated),
			},
		)
	}
	print.OutputList(flags.GlobalOutputValue, headers, values)

	return nil
}
