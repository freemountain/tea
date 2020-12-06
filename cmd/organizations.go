// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/tea/cmd/organizations"

	"github.com/urfave/cli/v2"
)

// CmdOrgs represents handle organization
var CmdOrgs = cli.Command{
	Name:        "organizations",
	Aliases:     []string{"organization", "org"},
	Usage:       "List, create, delete organizations",
	Description: "Show organization details",
	ArgsUsage:   "[<organization>]",
	Action:      runOrganizations,
	Subcommands: []*cli.Command{
		&organizations.CmdOrganizationList,
	},
}

func runOrganizations(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runOrganizationDetail(ctx.Args().First())
	}
	return organizations.RunOrganizationList(ctx)
}

func runOrganizationDetail(path string) error {

	log.Fatal("Not yet implemented.")
	return nil
}
