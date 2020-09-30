// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"bufio"
	"log"
	"os"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdLabelCreate represents a sub command of labels to create label.
var CmdLabelCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a label",
	Description: `Create a label`,
	Action:      runLabelCreate,
	Flags: []cli.Flag{
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
		&cli.StringFlag{
			Name:  "file",
			Usage: "indicate a label file",
		},
	},
}

func runLabelCreate(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	labelFile := ctx.String("file")
	var err error
	if len(labelFile) == 0 {
		_, _, err = login.Client().CreateLabel(owner, repo, gitea.CreateLabelOption{
			Name:        ctx.String("name"),
			Color:       ctx.String("color"),
			Description: ctx.String("description"),
		})
	} else {
		f, err := os.Open(labelFile)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		var i = 1
		for scanner.Scan() {
			line := scanner.Text()
			color, name, description := splitLabelLine(line)
			if color == "" || name == "" {
				log.Printf("Line %d ignored because lack of enough fields: %s\n", i, line)
			} else {
				_, _, err = login.Client().CreateLabel(owner, repo, gitea.CreateLabelOption{
					Name:        name,
					Color:       color,
					Description: description,
				})
			}

			i++
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func splitLabelLine(line string) (string, string, string) {
	fields := strings.SplitN(line, ";", 2)
	var color, name, description string
	if len(fields) < 1 {
		return "", "", ""
	} else if len(fields) >= 2 {
		description = strings.TrimSpace(fields[1])
	}
	fields = strings.Fields(fields[0])
	if len(fields) <= 0 {
		return "", "", ""
	}
	color = fields[0]
	if len(fields) == 2 {
		name = fields[1]
	} else if len(fields) > 2 {
		name = strings.Join(fields[1:], " ")
	}
	return color, name, description
}
