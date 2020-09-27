// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdMilestones represents to operate repositories milestones.
var CmdMilestones = cli.Command{
	Name:        "milestones",
	Aliases:     []string{"ms", "mile"},
	Usage:       "List and create milestones",
	Description: `List and create milestones`,
	ArgsUsage:   "[<milestone name>]",
	Action:      runMilestones,
	Subcommands: []*cli.Command{
		&CmdMilestonesList,
		&CmdMilestonesCreate,
		&CmdMilestonesClose,
		&CmdMilestonesDelete,
		&CmdMilestonesReopen,
		&CmdMilestonesIssues,
	},
	Flags: AllDefaultFlags,
}

// CmdMilestonesList represents a sub command of milestones to list milestones
var CmdMilestonesList = cli.Command{
	Name:        "ls",
	Usage:       "List milestones of the repository",
	Description: `List milestones of the repository`,
	Action:      runMilestonesList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by milestone state (all|open|closed)",
			DefaultText: "open",
		},
		&PaginationPageFlag,
		&PaginationLimitFlag,
	}, AllDefaultFlags...),
}

func runMilestones(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runMilestoneDetail(ctx, ctx.Args().First())
	}
	return runMilestonesList(ctx)
}

func runMilestoneDetail(ctx *cli.Context, name string) error {
	login, owner, repo := initCommand()
	client := login.Client()

	milestone, _, err := client.GetMilestoneByName(owner, repo, name)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n",
		milestone.Title,
	)
	if len(milestone.Description) != 0 {
		fmt.Printf("\n%s\n", milestone.Description)
	}
	if milestone.Deadline != nil && !milestone.Deadline.IsZero() {
		fmt.Printf("\nDeadline: %s\n", milestone.Deadline.Format("2006-01-02 15:04:05"))
	}
	return nil
}

func runMilestonesList(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "closed":
		state = gitea.StateClosed
	}

	milestones, _, err := login.Client().ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{
		ListOptions: getListOptions(ctx),
		State:       state,
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Title",
	}
	if state == gitea.StateAll {
		headers = append(headers, "State")
	}
	headers = append(headers,
		"Open/Closed Issues",
		"DueDate",
	)

	var values [][]string

	for _, m := range milestones {
		var deadline = ""

		if m.Deadline != nil && !m.Deadline.IsZero() {
			deadline = m.Deadline.Format("2006-01-02 15:04:05")
		}

		item := []string{
			m.Title,
		}
		if state == gitea.StateAll {
			item = append(item, string(m.State))
		}
		item = append(item,
			fmt.Sprintf("%d/%d", m.OpenIssues, m.ClosedIssues),
			deadline,
		)

		values = append(values, item)
	}
	Output(outputValue, headers, values)

	return nil
}

// CmdMilestonesCreate represents a sub command of milestones to create milestone
var CmdMilestonesCreate = cli.Command{
	Name:        "create",
	Usage:       "Create an milestone on repository",
	Description: `Create an milestone on repository`,
	Action:      runMilestonesCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "milestone title to create",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "milestone description to create",
		},
		&cli.StringFlag{
			Name:        "state",
			Usage:       "set milestone state (default is open)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
}

func runMilestonesCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	title := ctx.String("title")
	if len(title) == 0 {
		fmt.Printf("Title is required\n")
		return nil
	}

	state := gitea.StateOpen
	if ctx.String("state") == "closed" {
		state = gitea.StateClosed
	}

	mile, _, err := login.Client().CreateMilestone(owner, repo, gitea.CreateMilestoneOption{
		Title:       title,
		Description: ctx.String("description"),
		State:       state,
	})
	if err != nil {
		log.Fatal(err)
	}

	return runMilestoneDetail(ctx, mile.Title)
}

// CmdMilestonesClose represents a sub command of milestones to close an milestone
var CmdMilestonesClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of an milestone to 'closed'",
	Description: `Change state of an milestone to 'closed'`,
	ArgsUsage:   "<milestone name>",
	Action: func(ctx *cli.Context) error {
		if ctx.Bool("force") {
			return deleteMilestone(ctx)
		}
		return editMilestoneStatus(ctx, true)
	},
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "force",
			Aliases: []string{"f"},
			Usage:   "delete milestone",
		},
	}, AllDefaultFlags...),
}

func editMilestoneStatus(ctx *cli.Context, close bool) error {
	login, owner, repo := initCommand()
	client := login.Client()

	state := gitea.StateOpen
	if close {
		state = gitea.StateClosed
	}
	_, _, err := client.EditMilestoneByName(owner, repo, ctx.Args().First(), gitea.EditMilestoneOption{
		State: &state,
		Title: ctx.Args().First(),
	})

	return err
}

// CmdMilestonesDelete represents a sub command of milestones to delete an milestone
var CmdMilestonesDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "delete a milestone",
	Description: "delete a milestone",
	ArgsUsage:   "<milestone name>",
	Action:      deleteMilestone,
	Flags:       AllDefaultFlags,
}

func deleteMilestone(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	client := login.Client()

	_, err := client.DeleteMilestoneByName(owner, repo, ctx.Args().First())
	return err
}

// CmdMilestonesReopen represents a sub command of milestones to open an milestone
var CmdMilestonesReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of an milestone to 'open'",
	Description: `Change state of an milestone to 'open'`,
	ArgsUsage:   "<milestone name>",
	Action: func(ctx *cli.Context) error {
		return editMilestoneStatus(ctx, false)
	},
	Flags: AllDefaultFlags,
}
