// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/login"

	"github.com/urfave/cli/v2"
)

// CmdLogout represents to logout a gitea server.
var CmdLogout = cli.Command{
	Name:        "logout",
	Usage:       "Log out from a Gitea server",
	Description: `Log out from a Gitea server`,
	ArgsUsage:   "<login name>",
	Action:      login.RunLoginDelete,
}
