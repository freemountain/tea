// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// CmdLoginEdit represents to login a gitea server.
var CmdLoginEdit = cli.Command{
	Name:        "edit",
	Aliases:     []string{"e"},
	Usage:       "Edit Gitea logins",
	Description: `Edit Gitea logins`,
	Action:      runLoginEdit,
	Flags:       []cli.Flag{&flags.OutputFlag},
}

func runLoginEdit(_ *cli.Context) error {
	return open.Start(config.GetConfigPath())
}
