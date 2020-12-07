// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package times

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimesList represents a sub command of times to list them
var CmdTrackedTimesList = cli.Command{
	Name:    "ls",
	Aliases: []string{"list"},
	Action:  RunTimesList,
	Usage:   "Operate on tracked times of a repository's issues & pulls",
	Description: `Operate on tracked times of a repository's issues & pulls.
		 Depending on your permissions on the repository, only your own tracked
		 times might be listed.`,
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
	}, flags.AllDefaultFlags...),
}

// RunTimesList list repositories
func RunTimesList(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	var times []*gitea.TrackedTime
	var err error

	user := ctx.Args().First()
	fmt.Println(ctx.Command.ArgsUsage)
	if user == "" {
		// get all tracked times on the repo
		times, _, err = client.GetRepoTrackedTimes(owner, repo)
	} else if strings.HasPrefix(user, "#") {
		// get all tracked times on the specified issue
		issue, err := utils.ArgToIndex(user)
		if err != nil {
			return err
		}
		times, _, err = client.ListTrackedTimes(owner, repo, issue, gitea.ListTrackedTimesOptions{})
	} else {
		// get all tracked times by the specified user
		times, _, err = client.GetUserTrackedTimes(owner, repo, user)
	}

	if err != nil {
		return err
	}

	var from, until time.Time
	if ctx.String("from") != "" {
		from, err = dateparse.ParseLocal(ctx.String("from"))
		if err != nil {
			return err
		}
	}
	if ctx.String("until") != "" {
		until, err = dateparse.ParseLocal(ctx.String("until"))
		if err != nil {
			return err
		}
	}

	print.TrackedTimesList(times, flags.GlobalOutputValue, from, until, ctx.Bool("total"))
	return nil
}
