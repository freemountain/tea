// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"log"

	"code.gitea.io/tea/modules/config"

	"github.com/urfave/cli/v2"
)

// CmdLogout represents to logout a gitea server.
var CmdLogout = cli.Command{
	Name:        "logout",
	Usage:       "Log out from a Gitea server",
	Description: `Log out from a Gitea server`,
	Action:      runLogout,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Login name to remove",
		},
	},
}

func runLogout(ctx *cli.Context) error {
	logins, err := config.GetLogins()
	if err != nil {
		log.Fatal(err)
	}

	var name string

	if ctx.IsSet("name") {
		name = ctx.String("name")
	} else if len(ctx.Args().First()) != 0 {
		name = ctx.Args().First()
	} else if len(logins) == 1 {
		name = logins[0].Name
	} else {
		return errors.New("Please specify a login name")
	}

	return config.DeleteLogin(name)
}
