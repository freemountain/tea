// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdPullsReopen reopens a given closed pull request
var CmdPullsReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of a pull request to 'open'",
	Description: `Change state of a pull request to 'open'`,
	ArgsUsage:   "<pull index>",
	Action: func(ctx *cli.Context) error {
		var s = gitea.StateOpen
		return editPullState(ctx, gitea.EditPullRequestOption{State: &s})
	},
	Flags: flags.AllDefaultFlags,
}
