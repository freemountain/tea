// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdNotifications is the main command to operate with notifications
var CmdNotifications = cli.Command{
	Name:        "notifications",
	Usage:       "show notifications",
	Description: "show notifications, by default based of the current repo and unread one",
	Action:      runNotifications,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all notifications of related gitea instance",
		},
		/* // not supported jet
		&cli.BoolFlag{
			Name:    "read",
			Aliases: []string{"rd"},
			Usage:   "show read notifications instead unread",
		},
		*/
		&cli.IntFlag{
			Name:    "page",
			Aliases: []string{"p"},
			Usage:   "specify page, default is 1",
			Value:   1,
		},
		&cli.IntFlag{
			Name:    "limit",
			Aliases: []string{"lm"},
			Usage:   "specify limit of items per page",
		},
	}, AllDefaultFlags...),
}

func runNotifications(ctx *cli.Context) error {
	var news []*gitea.NotificationThread
	var err error

	listOpts := gitea.ListOptions{
		Page:     ctx.Int("page"),
		PageSize: ctx.Int("limit"),
	}

	if ctx.Bool("all") {
		login := initCommandLoginOnly()
		news, err = login.Client().ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
		})
	} else {
		login, owner, repo := initCommand()
		news, err = login.Client().ListRepoNotifications(owner, repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Type",
		"Index",
		"Title",
	}
	if ctx.Bool("all") {
		headers = append(headers, "Repository")
	}

	var values [][]string

	for _, n := range news {
		if n.Subject == nil {
			continue
		}
		// if pull or Issue get Index
		var index string
		if n.Subject.Type == "Issue" || n.Subject.Type == "Pull" {
			index = n.Subject.URL
			urlParts := strings.Split(n.Subject.URL, "/")
			if len(urlParts) != 0 {
				index = urlParts[len(urlParts)-1]
			}
			index = "#" + index
		}

		item := []string{n.Subject.Type, index, n.Subject.Title}
		if ctx.Bool("all") {
			item = append(item, n.Repository.FullName)
		}
		values = append(values, item)
	}

	if len(values) != 0 {
		Output(outputValue, headers, values)
	}
	return nil
}
