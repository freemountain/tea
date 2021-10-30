// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/users"
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v2"
)

// CmdOrgs represents handle organization
var CmdUsers = cli.Command{
	Name:        "users",
	Aliases:     []string{"users"},
	Category:    catEntities,
	Usage:       "List, create, delete users",
	Description: "Show user details",
	ArgsUsage:   "[<user>]",
	Action:      runUsers,
	Subcommands: []*cli.Command{
		&users.CmdUserList,
	},
	//Flags: organizations.CmdOrganizationList.Flags,
}

func runUsers(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	if ctx.Args().Len() == 1 {
		//	return runOrganizationDetail(ctx)
	}
	return users.RunUserList(cmd)
}
