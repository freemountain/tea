// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package releases

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReleaseCreate represents a sub command of Release to create release
var CmdReleaseCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a release",
	Description: `Create a release`,
	Action:      runReleaseCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "tag",
			Usage: "Tag name",
		},
		&cli.StringFlag{
			Name:  "target",
			Usage: "Target refs, branch name or commit id",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Release title",
		},
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"n"},
			Usage:   "Release notes",
		},
		&cli.BoolFlag{
			Name:    "draft",
			Aliases: []string{"d"},
			Usage:   "Is a draft",
		},
		&cli.BoolFlag{
			Name:    "prerelease",
			Aliases: []string{"p"},
			Usage:   "Is a pre-release",
		},
		&cli.StringSliceFlag{
			Name:    "asset",
			Aliases: []string{"a"},
			Usage:   "List of files to attach",
		},
	}, flags.AllDefaultFlags...),
}

func runReleaseCreate(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	release, resp, err := login.Client().CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      ctx.Bool("draft"),
		IsPrerelease: ctx.Bool("prerelease"),
	})

	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusConflict {
			fmt.Println("error: There already is a release for this tag")
			return nil
		}
		log.Fatal(err)
	}

	for _, asset := range ctx.StringSlice("asset") {
		var file *os.File

		if file, err = os.Open(asset); err != nil {
			log.Fatal(err)
		}

		filePath := filepath.Base(asset)

		if _, _, err = login.Client().CreateReleaseAttachment(owner, repo, release.ID, file, filePath); err != nil {
			file.Close()
			log.Fatal(err)
		}

		file.Close()
	}

	return nil
}
