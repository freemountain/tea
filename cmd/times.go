// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/times"
	"github.com/urfave/cli/v2"
)

// CmdTrackedTimes represents the command to operate repositories' times.
var CmdTrackedTimes = cli.Command{
	Name:    "times",
	Aliases: []string{"time"},
	Usage:   "Operate on tracked times of a repository's issues & pulls",
	Description: `Operate on tracked times of a repository's issues & pulls.
		 Depending on your permissions on the repository, only your own tracked
		 times might be listed.`,
	ArgsUsage: "[username | #issue]",
	Action:    times.RunTimesList,
	Subcommands: []*cli.Command{
		&times.CmdTrackedTimesAdd,
		&times.CmdTrackedTimesDelete,
		&times.CmdTrackedTimesReset,
		&times.CmdTrackedTimesList,
	},
}
