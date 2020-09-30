// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/interact"

	"github.com/urfave/cli/v2"
)

// CmdLoginAdd represents to login a gitea server.
var CmdLoginAdd = cli.Command{
	Name:        "add",
	Usage:       "Add a Gitea login",
	Description: `Add a Gitea login`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Usage:   "Login name",
		},
		&cli.StringFlag{
			Name:     "url",
			Aliases:  []string{"u"},
			Value:    "https://try.gitea.io",
			EnvVars:  []string{"GITEA_SERVER_URL"},
			Usage:    "Server URL",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "token",
			Aliases: []string{"t"},
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_TOKEN"},
			Usage:   "Access token. Can be obtained from Settings > Applications",
		},
		&cli.StringFlag{
			Name:    "user",
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_USER"},
			Usage:   "User for basic auth (will create token)",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"pwd"},
			Value:   "",
			EnvVars: []string{"GITEA_SERVER_PASSWORD"},
			Usage:   "Password for basic auth (will create token)",
		},
		&cli.StringFlag{
			Name:    "ssh-key",
			Aliases: []string{"s"},
			Usage:   "Path to a SSH key to use for pull/push operations",
		},
		&cli.BoolFlag{
			Name:    "insecure",
			Aliases: []string{"i"},
			Usage:   "Disable TLS verification",
		},
	},
	Action: runLoginAdd,
}

func runLoginAdd(ctx *cli.Context) error {
	// TODO: if no args -> interactive
	return config.AddLogin(
		ctx.String("name"),
		ctx.String("token"),
		ctx.String("user"),
		ctx.String("password"),
		ctx.String("ssh-key"),
		ctx.String("url"),
		ctx.Bool("insecure"))
}

// RunLoginAddInteractive create login interactive
// TODO: should be unexported
func RunLoginAddInteractive(_ *cli.Context) error {
	return interact.CreateLogin()
}
