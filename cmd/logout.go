// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"log"
	"os"

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
			Name:  "name, n",
			Usage: "Login name to remove",
		},
	},
}

func runLogout(ctx *cli.Context) error {
	var name string
	if len(os.Args) == 3 {
		name = os.Args[2]
	} else if ctx.IsSet("name") {
		name = ctx.String("name")
	} else {
		return errors.New("Please specify a login name")
	}

	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("Unable to load config file " + yamlConfigPath)
	}

	var idx = -1
	for i, l := range config.Logins {
		if l.Name == name {
			idx = i
			break
		}
	}
	if idx > -1 {
		config.Logins = append(config.Logins[:idx], config.Logins[idx+1:]...)
		err = saveConfig(yamlConfigPath)
		if err != nil {
			log.Fatal("Unable to save config file " + yamlConfigPath)
		}
	}

	return nil
}
