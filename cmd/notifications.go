// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/notifications"

	"github.com/urfave/cli/v2"
)

// CmdNotifications is the main command to operate with notifications
var CmdNotifications = cli.Command{
	Name:        "notifications",
	Aliases:     []string{"notification", "n"},
	Category:    catHelpers,
	Usage:       "Show notifications",
	Description: "Show notifications, by default based on the current repo if available",
	Action:      notifications.RunNotificationsList,
	Subcommands: []*cli.Command{
		&notifications.CmdNotificationsList,
		&notifications.CmdNotificationsMarkRead,
		&notifications.CmdNotificationsMarkUnread,
		&notifications.CmdNotificationsMarkPinned,
		&notifications.CmdNotificationsUnpin,
	},
	Flags: notifications.CmdNotificationsList.Flags,
}
