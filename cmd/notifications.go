// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdNotifications is the main command to operate with notifications
var CmdNotifications = cli.Command{
	Name:        "notifications",
	Aliases:     []string{"notification", "notif"},
	Usage:       "Show notifications",
	Description: "Show notifications, by default based of the current repo and unread one",
	Action:      runNotifications,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "show all notifications of related gitea instance",
		},
		&cli.BoolFlag{
			Name:    "read",
			Aliases: []string{"rd"},
			Usage:   "show read notifications instead unread",
		},
		&cli.BoolFlag{
			Name:    "pinned",
			Aliases: []string{"pd"},
			Usage:   "show pinned notifications instead unread",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

func runNotifications(ctx *cli.Context) error {
	var news []*gitea.NotificationThread
	var err error

	listOpts := flags.GetListOptions(ctx)
	if listOpts.Page == 0 {
		listOpts.Page = 1
	}

	var status []gitea.NotifyStatus
	if ctx.Bool("read") {
		status = []gitea.NotifyStatus{gitea.NotifyStatusRead}
	}
	if ctx.Bool("pinned") {
		status = append(status, gitea.NotifyStatusPinned)
	}

	if ctx.Bool("all") {
		login := config.InitCommandLoginOnly(flags.GlobalLoginValue)
		news, _, err = login.Client().ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	} else {
		login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
		news, _, err = login.Client().ListRepoNotifications(owner, repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
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
		print.OutputList(flags.GlobalOutputValue, headers, values)
	}
	return nil
}
