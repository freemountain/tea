// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package times

import (
	"fmt"
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

	"github.com/urfave/cli/v2"
)

// CmdTrackedTimesReset is a subcommand of CmdTrackedTimes, and
// clears all tracked times on an issue.
var CmdTrackedTimesReset = cli.Command{
	Name:      "reset",
	Usage:     "Reset tracked time on an issue",
	UsageText: "tea times reset <issue>",
	Action:    runTrackedTimesReset,
	Flags:     flags.LoginRepoFlags,
}

func runTrackedTimesReset(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	if err := client.CheckServerVersionConstraint(">= 1.11"); err != nil {
		return err
	}

	if ctx.Args().Len() != 1 {
		return fmt.Errorf("No issue specified.\nUsage:\t%s", ctx.Command.UsageText)
	}

	issue, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		log.Fatal(err)
	}

	_, err = client.ResetIssueTime(owner, repo, issue)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
