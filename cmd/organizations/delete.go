// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package organizations

import (
	"log"

	"code.gitea.io/tea/modules/context"
	"github.com/urfave/cli/v2"
)

// CmdOrganizationDelete represents a sub command of organizations to delete a given user organization
var CmdOrganizationDelete = cli.Command{
	Name:        "delete",
	Aliases:     []string{"rm"},
	Usage:       "Delete users Organizations",
	Description: "Delete users organizations",
	ArgsUsage:   "<organization name>",
	Action:      RunOrganizationDelete,
}

// RunOrganizationDelete delete user organization
func RunOrganizationDelete(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)

	client := ctx.Login.Client()

	if ctx.Args().Len() < 1 {
		log.Fatal("You have to specify the organization name you want to delete.")
		return nil
	}

	response, err := client.DeleteOrg(ctx.Args().First())
	if response != nil && response.StatusCode == 404 {
		log.Fatal("The given organization does not exist.")
		return nil
	}

	return err
}
