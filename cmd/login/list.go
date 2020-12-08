// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdLoginList represents to login a gitea server.
var CmdLoginList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List Gitea logins",
	Description: `List Gitea logins`,
	Action:      RunLoginList,
	Flags:       []cli.Flag{&flags.OutputFlag},
}

// RunLoginList list all logins
func RunLoginList(ctx *cli.Context) error {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	print.LoginsList(config.Config.Logins, flags.GlobalOutputValue)
	return nil
}
