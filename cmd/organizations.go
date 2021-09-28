// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"code.gitea.io/tea/cmd/organizations"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"github.com/urfave/cli/v2"
)

// CmdOrgs represents handle organization
var CmdOrgs = cli.Command{
	Name:        "organizations",
	Aliases:     []string{"organization", "org"},
	Category:    catEntities,
	Usage:       "List, create, delete organizations",
	Description: "Show organization details",
	ArgsUsage:   "[<organization>]",
	Action:      runOrganizations,
	Subcommands: []*cli.Command{
		&organizations.CmdOrganizationList,
		&organizations.CmdOrganizationCreate,
		&organizations.CmdOrganizationDelete,
	},
	Flags: organizations.CmdOrganizationList.Flags,
}

func runOrganizations(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	if ctx.Args().Len() == 1 {
		return runOrganizationDetail(ctx)
	}
	return organizations.RunOrganizationList(cmd)
}

func runOrganizationDetail(ctx *context.TeaContext) error {
	org, _, err := ctx.Login.Client().GetOrg(ctx.Args().First())
	if err != nil {
		return err
	}

	print.OrganizationDetails(org)
	return nil
}
