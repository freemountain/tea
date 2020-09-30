// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesClose represents a sub command of issues to close an issue
var CmdIssuesClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of an issue to 'closed'",
	Description: `Change state of an issue to 'closed'`,
	ArgsUsage:   "<issue index>",
	Action: func(ctx *cli.Context) error {
		var s = gitea.StateClosed
		return editIssueState(ctx, gitea.EditIssueOption{State: &s})
	},
	Flags: flags.AllDefaultFlags,
}

// editIssueState abstracts the arg parsing to edit the given issue
func editIssueState(ctx *cli.Context, opts gitea.EditIssueOption) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if ctx.Args().Len() == 0 {
		log.Fatal(ctx.Command.ArgsUsage)
	}

	index, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	_, _, err = login.Client().EditIssue(owner, repo, index, opts)
	// TODO: print (short)IssueDetails
	return err
}
