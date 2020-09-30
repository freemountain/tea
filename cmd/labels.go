// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/cmd/labels"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v2"
)

// CmdLabels represents to operate repositories' labels.
var CmdLabels = cli.Command{
	Name:        "labels",
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

	headers := []string{
		"Index",
		"Color",
		"Name",
		"Description",
	}

	var values [][]string

	labels, _, err := login.Client().ListRepoLabels(owner, repo, gitea.ListLabelsOptions{ListOptions: flags.GetListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	if len(labels) == 0 {
		print.OutputList(flags.GlobalOutputValue, headers, values)
		return nil
	}

	p := termenv.ColorProfile()

	fPath := ctx.String("save")
	if len(fPath) > 0 {
		f, err := os.Create(fPath)
		if err != nil {
			return err
		}
		defer f.Close()

		for _, label := range labels {
			fmt.Fprintf(f, "#%s %s\n", label.Color, label.Name)
		}
	} else {
		for _, label := range labels {
			color := termenv.String(label.Color)

			values = append(
				values,
				[]string{
					strconv.FormatInt(label.ID, 10),
					fmt.Sprint(color.Background(p.Color("#" + label.Color))),
					label.Name,
					label.Description,
				},
			)
		}
		print.OutputList(flags.GlobalOutputValue, headers, values)
	}

	return nil
}
