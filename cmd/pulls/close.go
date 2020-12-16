// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdPullsClose closes a given open pull request
var CmdPullsClose = cli.Command{
	Name:        "close",
	Usage:       "Change state of a pull request to 'closed'",
	Description: `Change state of a pull request to 'closed'`,
	ArgsUsage:   "<pull index>",
	Action: func(ctx *cli.Context) error {
		var s = gitea.StateClosed
		return editPullState(ctx, gitea.EditPullRequestOption{State: &s})
	},
	Flags: flags.AllDefaultFlags,
}
