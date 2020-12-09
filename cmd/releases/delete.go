// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package releases

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"github.com/urfave/cli/v2"
)

// CmdReleaseDelete represents a sub command of Release to delete a release
var CmdReleaseDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a release",
	Description: `Delete a release`,
	ArgsUsage:   "<release tag>",
	Action:      runReleaseDelete,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:    "confirm",
			Aliases: []string{"y"},
			Usage:   "Confirm deletion (required)",
		},
		&cli.BoolFlag{
			Name:  "delete-tag",
			Usage: "Also delete the git tag for this release",
		},
	}, flags.AllDefaultFlags...),
}

func runReleaseDelete(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	tag := ctx.Args().First()
	if len(tag) == 0 {
		fmt.Println("Release tag needed to delete")
		return nil
	}

	if !ctx.Bool("confirm") {
		fmt.Println("Are you sure? Please confirm with -y or --confirm.")
		return nil
	}

	release, err := getReleaseByTag(owner, repo, tag, client)
	if err != nil {
		return err
	}
	if release == nil {
		return nil
	}

	_, err = client.DeleteRelease(owner, repo, release.ID)
	if err != nil {
		return err
	}

	if ctx.Bool("delete-tag") {
		_, err = client.DeleteReleaseTag(owner, repo, tag)
		return err
	}

	return nil
}
