// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

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

	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

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
		news, _, err = login.Client().ListNotifications(gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	} else {
		news, _, err = login.Client().ListRepoNotifications(owner, repo, gitea.ListNotificationOptions{
			ListOptions: listOpts,
			Status:      status,
		})
	}
	if err != nil {
		log.Fatal(err)
	}

	print.NotificationsList(news, flags.GlobalOutputValue, ctx.Bool("all"))
	return nil
}
