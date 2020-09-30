// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package issues

import (
	"code.gitea.io/tea/cmd/flags"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdIssuesReopen represents a sub command of issues to open an issue
var CmdIssuesReopen = cli.Command{
	Name:        "reopen",
	Aliases:     []string{"open"},
	Usage:       "Change state of an issue to 'open'",
	Description: `Change state of an issue to 'open'`,
	ArgsUsage:   "<issue index>",
	Action: func(ctx *cli.Context) error {
		var s = gitea.StateOpen
		return editIssueState(ctx, gitea.EditIssueOption{State: &s})
	},
	Flags: flags.AllDefaultFlags,
}
