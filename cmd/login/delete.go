// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"errors"
	"log"

	"code.gitea.io/tea/modules/config"

	"github.com/urfave/cli/v2"
)

// CmdLoginDelete is a command to delete a login
var CmdLoginDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Remove a Gitea login",
	Description: `Remove a Gitea login`,
	ArgsUsage:   "<login name>",
	Action:      RunLoginDelete,
}

// RunLoginDelete runs the action of a login delete command
func RunLoginDelete(ctx *cli.Context) error {
	logins, err := config.GetLogins()
	if err != nil {
		log.Fatal(err)
	}

	var name string

	if len(ctx.Args().First()) != 0 {
		name = ctx.Args().First()
	} else if len(logins) == 1 {
		name = logins[0].Name
	} else {
		return errors.New("Please specify a login name")
	}

	return config.DeleteLogin(name)
}
