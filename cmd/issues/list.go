// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesList represents a sub command of issues to list issues
var CmdIssuesList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List issues of the repository",
	Description: `List issues of the repository`,
	Action:      RunIssuesList,
	Flags: append([]cli.Flag{
		flags.FieldsFlag(print.IssueFields, []string{
			"index", "title", "state", "author", "milestone", "labels",
		}),
	}, flags.IssuePRFlags...),
}

// RunIssuesList list issues
func RunIssuesList(cmd *cli.Context) error {
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

	issues, _, err := ctx.Login.Client().ListRepoIssues(ctx.Owner, ctx.Repo, gitea.ListIssueOption{
		ListOptions: ctx.GetListOptions(),
		State:       state,
		Type:        gitea.IssueTypeIssue,
	})

	if err != nil {
		return err
	}

	fields, err := flags.GetFields(cmd, print.IssueFields)
	if err != nil {
		return err
	}

	print.IssuesPullsList(issues, ctx.Output, fields)
	return nil
}
