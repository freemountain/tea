// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package users

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdOrganizationList represents a sub command of organizations to list users organizations
var CmdUserList = cli.Command{
	Name:        "list",
	Aliases:     []string{"ls"},
	Usage:       "List Users",
	Description: "List users",
	Action:      RunUserList,
	Flags: append([]cli.Flag{
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunOrganizationList list user organizations
func RunUserList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()

	users, _, err := client.AdminListUsers(gitea.AdminListUsersOptions{
		ListOptions: ctx.GetListOptions(),
	})

	if err != nil {
		return err
	}

	print.UserList(users, ctx.Output, print.UserFields)

	return nil
}
