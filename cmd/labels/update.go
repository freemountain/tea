// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdLabelUpdate represents a sub command of labels to update label.
var CmdLabelUpdate = cli.Command{
	Name:        "update",
	Usage:       "Update a label",
	Description: `Update a label`,
	Action:      runLabelUpdate,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "label id",
		},
		&cli.StringFlag{
			Name:  "name",
			Usage: "label name",
		},
		&cli.StringFlag{
			Name:  "color",
			Usage: "label color value",
		},
		&cli.StringFlag{
			Name:  "description",
			Usage: "label description",
		},
	},
}

func runLabelUpdate(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	id := ctx.Int64("id")
	var pName, pColor, pDescription *string
	name := ctx.String("name")
	if name != "" {
		pName = &name
	}

	color := ctx.String("color")
	if color != "" {
		pColor = &color
	}

	description := ctx.String("description")
	if description != "" {
		pDescription = &description
	}

	var err error
	_, _, err = login.Client().EditLabel(owner, repo, id, gitea.EditLabelOption{
		Name:        pName,
		Color:       pColor,
		Description: pDescription,
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
