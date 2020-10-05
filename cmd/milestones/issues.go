// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package milestones

import (
	"fmt"
	"strconv"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestonesIssues represents a sub command of milestones to manage issue/pull of an milestone
var CmdMilestonesIssues = cli.Command{
	Name:        "issues",
	Aliases:     []string{"i"},
	Usage:       "manage issue/pull of an milestone",
	Description: "manage issue/pull of an milestone",
	ArgsUsage:   "<milestone name>",
	Action:      runMilestoneIssueList,
	Subcommands: []*cli.Command{
		&CmdMilestoneAddIssue,
		&CmdMilestoneRemoveIssue,
	},
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by issue state (all|open|closed)",
			DefaultText: "open",
		},
		&cli.StringFlag{
			Name:  "kind",
			Usage: "Filter by kind (issue|pull)",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// CmdMilestoneAddIssue represents a sub command of milestone issues to add an issue/pull to an milestone
var CmdMilestoneAddIssue = cli.Command{
	Name:        "add",
	Aliases:     []string{"a"},
	Usage:       "Add an issue/pull to an milestone",
	Description: "Add an issue/pull to an milestone",
	ArgsUsage:   "<milestone name> <issue/pull index>",
	Action:      runMilestoneIssueAdd,
	Flags:       flags.AllDefaultFlags,
}

// CmdMilestoneRemoveIssue represents a sub command of milestones to remove an issue/pull from an milestone
var CmdMilestoneRemoveIssue = cli.Command{
	Name:        "remove",
	Aliases:     []string{"r"},
	Usage:       "Remove an issue/pull to an milestone",
	Description: "Remove an issue/pull to an milestone",
	ArgsUsage:   "<milestone name> <issue/pull index>",
	Action:      runMilestoneIssueRemove,
	Flags:       flags.AllDefaultFlags,
}

func runMilestoneIssueList(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	kind := gitea.IssueTypeAll
	switch ctx.String("kind") {
	case "issue":
		kind = gitea.IssueTypeIssue
	case "pull":
		kind = gitea.IssueTypePull
	}

	fmt.Println(state)

	milestone := ctx.Args().First()
	// make sure milestone exist
	_, _, err := client.GetMilestoneByName(owner, repo, milestone)
	if err != nil {
		return err
	}

	issues, _, err := client.ListRepoIssues(owner, repo, gitea.ListIssueOption{
		ListOptions: flags.GetListOptions(ctx),
		Milestones:  []string{milestone},
		Type:        kind,
		State:       state,
	})
	if err != nil {
		return err
	}

	headers := []string{
		"Index",
		"State",
		"Kind",
		"Author",
		"Updated",
		"Title",
	}

	var values [][]string

	if len(issues) == 0 {
		print.OutputList(flags.GlobalOutputValue, headers, values)
		return nil
	}

	for _, issue := range issues {
		name := issue.Poster.FullName
		if len(name) == 0 {
			name = issue.Poster.UserName
		}
		kind := "Issue"
		if issue.PullRequest != nil {
			kind = "Pull"
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(issue.Index, 10),
				string(issue.State),
				kind,
				name,
				print.FormatTime(issue.Updated),
				issue.Title,
			},
		)
	}
	print.OutputList(flags.GlobalOutputValue, headers, values)
	return nil
}

func runMilestoneIssueAdd(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("need two arguments")
	}

	mileName := ctx.Args().Get(0)
	issueIndex := ctx.Args().Get(1)
	idx, err := utils.ArgToIndex(issueIndex)
	if err != nil {
		return err
	}

	// make sure milestone exist
	mile, _, err := client.GetMilestoneByName(owner, repo, mileName)
	if err != nil {
		return err
	}

	_, _, err = client.EditIssue(owner, repo, idx, gitea.EditIssueOption{
		Milestone: &mile.ID,
	})
	return err
}

func runMilestoneIssueRemove(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()
	if ctx.Args().Len() == 0 {
		return fmt.Errorf("need two arguments")
	}

	mileName := ctx.Args().Get(0)
	issueIndex := ctx.Args().Get(1)
	idx, err := utils.ArgToIndex(issueIndex)
	if err != nil {
		return err
	}

	issue, _, err := client.GetIssue(owner, repo, idx)
	if err != nil {
		return err
	}

	if issue.Milestone == nil {
		return fmt.Errorf("issue is not assigned to a milestone")
	}

	if issue.Milestone.Title != mileName {
		return fmt.Errorf("issue is not assigned to this milestone")
	}

	zero := int64(0)
	_, _, err = client.EditIssue(owner, repo, idx, gitea.EditIssueOption{
		Milestone: &zero,
	})
	return err
}
