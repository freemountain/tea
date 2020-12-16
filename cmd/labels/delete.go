// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"code.gitea.io/tea/modules/context"

	"github.com/urfave/cli/v2"
)

// CmdLabelDelete represents a sub command of labels to delete label.
var CmdLabelDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete a label",
	Description: `Delete a label`,
	Action:      runLabelDelete,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "id",
			Usage: "label id",
		},
	},
}

func runLabelDelete(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	_, err := ctx.Login.Client().DeleteLabel(ctx.Owner, ctx.Repo, ctx.Int64("id"))
	return err
}
