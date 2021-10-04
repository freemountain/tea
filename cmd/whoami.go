// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdWhoami represents the command to show current logged in user
var CmdWhoami = cli.Command{
	Name:        "whoami",
	Category:    catSetup,
	Description: `For debugging purposes, show the user that is currently logged in.`,
	Usage:       "Show current logged in user",
	Action: func(cmd *cli.Context) error {
		ctx := context.InitCommand(cmd)
		client := ctx.Login.Client()
		user, _, _ := client.GetMyUserInfo()
		print.UserDetails(user)
		return nil
	},
}
