// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package times

import (
	"fmt"
	"log"
	"strconv"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdTrackedTimesDelete is a sub command of CmdTrackedTimes, and removes time from an issue
var CmdTrackedTimesDelete = cli.Command{
	Name:      "delete",
	Aliases:   []string{"rm"},
	Usage:     "Delete a single tracked time on an issue",
	UsageText: "tea times delete <issue> <time ID>",
	Action:    runTrackedTimesDelete,
	Flags:     flags.LoginRepoFlags,
}

func runTrackedTimesDelete(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	if err := client.CheckServerVersionConstraint(">= 1.11"); err != nil {
		return err
	}

	if ctx.Args().Len() < 2 {
		return fmt.Errorf("No issue or time ID specified.\nUsage:\t%s", ctx.Command.UsageText)
	}

	issue, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		log.Fatal(err)
	}

	timeID, err := strconv.ParseInt(ctx.Args().Get(1), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.DeleteTime(owner, repo, issue, timeID)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
