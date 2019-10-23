// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdLabels represents to operate repositories' labels.
var CmdLabels = cli.Command{
	Name:        "labels",
	Usage:       "Operate with labels of the repository",
	Description: `Operate with labels of the repository`,
	Action:      runLabels,
	Subcommands: []cli.Command{
		CmdLabelCreate,
		CmdLabelUpdate,
		CmdLabelDelete,
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repository, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "save, s",
			Usage: "Save all the labels as a file",
		},
	},
}

func runLabels(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	labels, err := login.Client().ListRepoLabels(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	if len(labels) == 0 {
		fmt.Println("No Labels")
		return nil
	}

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
			fmt.Fprintf(os.Stdout, "%d #%s %s\n", label.ID, label.Color, label.Name)
		}
	}

	return nil
}

// CmdLabelCreate represents a sub command of labels to create label.
var CmdLabelCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a label in repository",
	Description: `Create a label in repository`,
	Action:      runLabelCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name",
			Usage: "label name",
		},
		cli.StringFlag{
			Name:  "color",
			Usage: "label color value",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "label description",
		},
		cli.StringFlag{
			Name:  "file",
			Usage: "indicate a label file",
		},
	},
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

func runLabelCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	labelFile := ctx.String("file")
	var err error
	if len(labelFile) == 0 {
		_, err = login.Client().CreateLabel(owner, repo, gitea.CreateLabelOption{
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
		// FIXME: if Gitea's API support create multiple labels once, we should move to that API.
		for scanner.Scan() {
			line := scanner.Text()
			color, name, description := splitLabelLine(line)
			if color == "" || name == "" {
				log.Printf("Line %d ignored because lack of enough fields: %s\n", i, line)
			} else {
				_, err = login.Client().CreateLabel(owner, repo, gitea.CreateLabelOption{
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

// CmdLabelUpdate represents a sub command of labels to update label.
var CmdLabelUpdate = cli.Command{
	Name:        "update",
	Usage:       "Update a label in repository",
	Description: `Update a label in repository`,
	Action:      runLabelUpdate,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "id",
			Usage: "label id",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "label name",
		},
		cli.StringFlag{
			Name:  "color",
			Usage: "label color value",
		},
		cli.StringFlag{
			Name:  "description",
			Usage: "label description",
		},
	},
}

func runLabelUpdate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

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
	_, err = login.Client().EditLabel(owner, repo, id, gitea.EditLabelOption{
		Name:        pName,
		Color:       pColor,
		Description: pDescription,
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// CmdLabelDelete represents a sub command of labels to delete label.
var CmdLabelDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a label in repository",
	Description: `Delete a label in repository`,
	Action:      runLabelCreate,
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "id",
			Usage: "label id",
		},
	},
}

func runLabelDelete(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	err := login.Client().DeleteLabel(owner, repo, ctx.Int64("id"))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
