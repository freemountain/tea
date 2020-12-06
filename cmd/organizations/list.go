// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package organizations

import (
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdOrganizationList represents a sub command of organizations to list users organizations
var CmdOrganizationList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List Organizations",
	Description: "List users organizations",
	Action:      RunOrganizationList,
	Flags: append([]cli.Flag{
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunOrganizationList list user organizations
func RunOrganizationList(ctx *cli.Context) error {
	//TODO: Reconsider the usage InitCommandLoginOnly related to #200
	login := config.InitCommandLoginOnly(flags.GlobalLoginValue)

	client := login.Client()

	userOrganizations, _, err := client.ListUserOrgs(login.User, gitea.ListOrgsOptions{ListOptions: flags.GetListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	print.OrganizationsList(userOrganizations)

	return nil
}
