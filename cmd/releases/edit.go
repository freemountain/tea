// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package releases

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReleaseEdit represents a sub command of Release to edit releases
var CmdReleaseEdit = cli.Command{
	Name:        "edit",
	Usage:       "Edit a release",
	Description: `Edit a release`,
	ArgsUsage:   "<release tag>",
	Action:      runReleaseEdit,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "tag",
			Usage: "Change Tag",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "Change Target",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Change Title",
		},
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"n"},
			Usage:   "Change Notes",
		},
		&cli.StringFlag{
			Name:        "draft",
			Aliases:     []string{"d"},
			Usage:       "Mark as Draft [True/false]",
			DefaultText: "true",
		},
		&cli.StringFlag{
			Name:        "prerelease",
			Aliases:     []string{"p"},
			Usage:       "Mark as Pre-Release [True/false]",
			DefaultText: "true",
		},
	}, flags.AllDefaultFlags...),
}

func runReleaseEdit(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})
	client := ctx.Login.Client()

	tag := ctx.Args().First()
	if len(tag) == 0 {
		fmt.Println("Release tag needed to edit")
		return nil
	}

	release, err := getReleaseByTag(ctx.Owner, ctx.Repo, tag, client)
	if err != nil {
		return err
	}
	if release == nil {
		return nil
	}

	var isDraft, isPre *bool
	if ctx.IsSet("draft") {
		isDraft = gitea.OptionalBool(strings.ToLower(ctx.String("draft"))[:1] == "t")
	}
	if ctx.IsSet("prerelease") {
		isPre = gitea.OptionalBool(strings.ToLower(ctx.String("prerelease"))[:1] == "t")
	}

	_, _, err = client.EditRelease(ctx.Owner, ctx.Repo, release.ID, gitea.EditReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      isDraft,
		IsPrerelease: isPre,
	})
	return err
}
