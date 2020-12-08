// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/labels"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdLabels represents to operate repositories' labels.
var CmdLabels = cli.Command{
	Name:        "labels",
	Aliases:     []string{"label"},
	Usage:       "Manage issue labels",
	Description: `Manage issue labels`,
	Action:      runLabels,
	Subcommands: []*cli.Command{
		&labels.CmdLabelCreate,
		&labels.CmdLabelUpdate,
		&labels.CmdLabelDelete,
	},
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   "Save all the labels as a file",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

func runLabels(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	labels, _, err := login.Client().ListRepoLabels(owner, repo, gitea.ListLabelsOptions{ListOptions: flags.GetListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	if ctx.IsSet("save") {
		return task.LabelsExport(labels, ctx.String("save"))
	}

	print.LabelsList(labels, flags.GlobalOutputValue)
	return nil
}
