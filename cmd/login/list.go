// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"fmt"
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
	Action:      runLoginList,
	Flags:       []cli.Flag{&flags.OutputFlag},
}

func runLoginList(ctx *cli.Context) error {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Name",
		"URL",
		"SSHHost",
		"User",
		"Default",
	}

	var values [][]string

	for _, l := range config.Config.Logins {
		values = append(values, []string{
			l.Name,
			l.URL,
			l.GetSSHHost(),
			l.User,
			fmt.Sprint(l.Default),
		})
	}

	print.OutputList(flags.GlobalOutputValue, headers, values)

	return nil
}
