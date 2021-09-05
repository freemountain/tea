// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package notifications

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/utils"
	"github.com/urfave/cli/v2"
)

// CmdNotificationsMarkRead represents a sub command of notifications to list read notifications
var CmdNotificationsMarkRead = cli.Command{
	Name:        "read",
	Aliases:     []string{"r"},
	Usage:       "Mark all filtered or a specific notification as read",
	Description: "Mark all filtered or a specific notification as read",
	ArgsUsage:   "[all | <notification id>]",
	Flags:       flags.NotificationFlags,
	Action: func(ctx *cli.Context) error {
		cmd := context.InitCommand(ctx)
		filter, err := flags.NotificationStateFlag.GetValues(ctx)
		if err != nil {
			return err
		}
		if !flags.NotificationStateFlag.IsSet() {
			filter = []string{string(gitea.NotifyStatusUnread)}
		}
		return markNotificationAs(cmd, filter, gitea.NotifyStatusRead)
	},
}

// CmdNotificationsMarkUnread will mark notifications as unread.
var CmdNotificationsMarkUnread = cli.Command{
	Name:        "unread",
	Aliases:     []string{"u"},
	Usage:       "Mark all filtered or a specific notification as unread",
	Description: "Mark all filtered or a specific notification as unread",
	ArgsUsage:   "[all | <notification id>]",
	Flags:       flags.NotificationFlags,
	Action: func(ctx *cli.Context) error {
		cmd := context.InitCommand(ctx)
		filter, err := flags.NotificationStateFlag.GetValues(ctx)
		if err != nil {
			return err
		}
		if !flags.NotificationStateFlag.IsSet() {
			filter = []string{string(gitea.NotifyStatusRead)}
		}
		return markNotificationAs(cmd, filter, gitea.NotifyStatusUnread)
	},
}

// CmdNotificationsMarkPinned will mark notifications as unread.
var CmdNotificationsMarkPinned = cli.Command{
	Name:        "pin",
	Aliases:     []string{"p"},
	Usage:       "Mark all filtered or a specific notification as pinned",
	Description: "Mark all filtered or a specific notification as pinned",
	ArgsUsage:   "[all | <notification id>]",
	Flags:       flags.NotificationFlags,
	Action: func(ctx *cli.Context) error {
		cmd := context.InitCommand(ctx)
		filter, err := flags.NotificationStateFlag.GetValues(ctx)
		if err != nil {
			return err
		}
		if !flags.NotificationStateFlag.IsSet() {
			filter = []string{string(gitea.NotifyStatusUnread)}
		}
		return markNotificationAs(cmd, filter, gitea.NotifyStatusPinned)
	},
}

// CmdNotificationsUnpin will mark pinned notifications as unread.
var CmdNotificationsUnpin = cli.Command{
	Name:        "unpin",
	Usage:       "Unpin all pinned or a specific notification",
	Description: "Marks all pinned or a specific notification as read",
	ArgsUsage:   "[all | <notification id>]",
	Flags:       flags.NotificationFlags,
	Action: func(ctx *cli.Context) error {
		cmd := context.InitCommand(ctx)
		filter := []string{string(gitea.NotifyStatusPinned)}
		// NOTE: we implicitly mark it as read, to match web UI semantics. marking as unread might be more useful?
		return markNotificationAs(cmd, filter, gitea.NotifyStatusRead)
	},
}

func markNotificationAs(cmd *context.TeaContext, filterStates []string, targetState gitea.NotifyStatus) (err error) {
	client := cmd.Login.Client()
	subject := cmd.Args().First()
	allRepos := cmd.Bool("mine")

	states := []gitea.NotifyStatus{}
	for _, s := range filterStates {
		states = append(states, gitea.NotifyStatus(s))
	}

	switch subject {
	case "", "all":
		opts := gitea.MarkNotificationOptions{Status: states, ToStatus: targetState}

		if allRepos {
			_, err = client.ReadNotifications(opts)
		} else {
			cmd.Ensure(context.CtxRequirement{RemoteRepo: true})
			_, err = client.ReadRepoNotifications(cmd.Owner, cmd.Repo, opts)
		}

		// TODO: print all affected notification subject URLs
		// (not supported by API currently, https://github.com/go-gitea/gitea/issues/16797)

	default:
		id, err := utils.ArgToIndex(subject)
		if err != nil {
			return err
		}
		_, err = client.ReadNotification(id, targetState)
		if err != nil {
			return err
		}

		n, _, err := client.GetNotification(id)
		if err != nil {
			return err
		}
		// FIXME: this is an API URL, we want to display a web ui link..
		fmt.Println(n.Subject.URL)
		return nil
	}

	return err
}
