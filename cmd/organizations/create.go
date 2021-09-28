// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package organizations

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdOrganizationCreate represents a sub command of organizations to delete a given user organization
var CmdOrganizationCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create an organization",
	Description: "Create an organization",
	Action:      RunOrganizationCreate,
	ArgsUsage:   "<organization name>",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
		},
		&cli.StringFlag{
			Name:    "website",
			Aliases: []string{"w"},
		},
		&cli.StringFlag{
			Name:    "location",
			Aliases: []string{"L"},
		},
		&cli.StringFlag{
			Name:    "visibility",
			Aliases: []string{"v"},
		},
		&cli.BoolFlag{
			Name: "repo-admins-can-change-team-access",
		},
		&flags.LoginFlag,
	},
}

// RunOrganizationCreate sets up a new organization
func RunOrganizationCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)

	if ctx.Args().Len() < 1 {
		return fmt.Errorf("You have to specify the organization name you want to create")
	}

	var visibility gitea.VisibleType
	switch ctx.String("visibility") {
	case "", "public":
		visibility = gitea.VisibleTypePublic
	case "private":
		visibility = gitea.VisibleTypePrivate
	case "limited":
		visibility = gitea.VisibleTypeLimited
	default:
		return fmt.Errorf("unknown visibility '%s'", ctx.String("visibility"))
	}

	org, _, err := ctx.Login.Client().CreateOrg(gitea.CreateOrgOption{
		Name: ctx.Args().First(),
		// FullName: , // not really meaningful for orgs (not displayed in webui, use description instead?)
		Description:               ctx.String("description"),
		Website:                   ctx.String("website"),
		Location:                  ctx.String("location"),
		RepoAdminChangeTeamAccess: ctx.Bool("repo-admins-can-change-team-access"),
		Visibility:                visibility,
	})
	if err != nil {
		return err
	}

	print.OrganizationDetails(org)

	return err
}
