// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package times

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimesList represents a sub command of times to list them
var CmdTrackedTimesList = cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Action:  RunTimesList,
	Usage:   "List tracked times on issues & pulls",
	Description: `List tracked times, across repos, or on a single repo or issue:
- given a username all times on a repo by that user are shown,
- given a issue index with '#' prefix, all times on that issue are listed,
- given --mine, your times are listed across all repositories.
Depending on your permissions on the repository, only your own tracked times might be listed.`,
	ArgsUsage: "[username | #issue]",

	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "from",
			Aliases: []string{"f"},
			Usage:   "Show only times tracked after this date",
		},
		&cli.StringFlag{
			Name:    "until",
			Aliases: []string{"u"},
			Usage:   "Show only times tracked before this date",
		},
		&cli.BoolFlag{
			Name:    "total",
			Aliases: []string{"t"},
			Usage:   "Print the total duration at the end",
		},
		&cli.BoolFlag{
			Name:    "mine",
			Aliases: []string{"m"},
			Usage:   "Show all times tracked by you across all repositories (overrides command arguments)",
		},
		&cli.StringFlag{
			Name: "fields",
			Usage: fmt.Sprintf(`Comma-separated list of fields to print. Available values:
			%s
		`, strings.Join(print.TrackedTimeFields, ",")),
		},
	}, flags.AllDefaultFlags...),
}

// RunTimesList list repositories
func RunTimesList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	client := ctx.Login.Client()

	var times []*gitea.TrackedTime
	var err error
	var from, until time.Time
	var fields []string

	if ctx.IsSet("from") {
		from, err = dateparse.ParseLocal(ctx.String("from"))
		if err != nil {
			return err
		}
	}
	if ctx.IsSet("until") {
		until, err = dateparse.ParseLocal(ctx.String("until"))
		if err != nil {
			return err
		}
	}

	opts := gitea.ListTrackedTimesOptions{Since: from, Before: until}

	user := ctx.Args().First()
	if ctx.Bool("mine") {
		times, _, err = client.GetMyTrackedTimes()
		fields = []string{"created", "repo", "issue", "duration"}
	} else if user == "" {
		// get all tracked times on the repo
		times, _, err = client.ListRepoTrackedTimes(ctx.Owner, ctx.Repo, opts)
		fields = []string{"created", "issue", "user", "duration"}
	} else if strings.HasPrefix(user, "#") {
		// get all tracked times on the specified issue
		issue, err := utils.ArgToIndex(user)
		if err != nil {
			return err
		}
		times, _, err = client.ListIssueTrackedTimes(ctx.Owner, ctx.Repo, issue, opts)
		fields = []string{"created", "user", "duration"}
	} else {
		// get all tracked times by the specified user
		opts.User = user
		times, _, err = client.ListRepoTrackedTimes(ctx.Owner, ctx.Repo, opts)
		fields = []string{"created", "issue", "duration"}
	}

	if err != nil {
		return err
	}

	if ctx.IsSet("fields") {
		if fields, err = flags.GetFields(cmd, print.TrackedTimeFields); err != nil {
			return err
		}
	}

	print.TrackedTimesList(times, ctx.Output, fields, ctx.Bool("total"))
	return nil
}
