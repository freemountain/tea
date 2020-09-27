// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli/v2"
)

// CmdReleases represents to login a gitea server.
var CmdReleases = cli.Command{
	Name:        "release",
	Aliases:     []string{"releases"},
	Usage:       "Manage releases",
	Description: "Manage releases",
	Action:      runReleases,
	Subcommands: []*cli.Command{
		&CmdReleaseList,
		&CmdReleaseCreate,
		&CmdReleaseDelete,
		&CmdReleaseEdit,
	},
	Flags: AllDefaultFlags,
}

// CmdReleaseList represents a sub command of Release to list releases
var CmdReleaseList = cli.Command{
	Name:        "ls",
	Usage:       "List Releases",
	Description: "List Releases",
	Action:      runReleases,
	Flags: append([]cli.Flag{
		&PaginationPageFlag,
		&PaginationLimitFlag,
	}, AllDefaultFlags...),
}

func runReleases(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	releases, _, err := login.Client().ListReleases(owner, repo, gitea.ListReleasesOptions{ListOptions: getListOptions(ctx)})
	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Tag-Name",
		"Title",
		"Published At",
		"Status",
		"Tar URL",
	}

	var values [][]string

	if len(releases) == 0 {
		Output(outputValue, headers, values)
		return nil
	}

	for _, release := range releases {
		status := "released"
		if release.IsDraft {
			status = "draft"
		} else if release.IsPrerelease {
			status = "prerelease"
		}
		values = append(
			values,
			[]string{
				release.TagName,
				release.Title,
				release.PublishedAt.Format("2006-01-02 15:04:05"),
				status,
				release.TarURL,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}

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
	}, AllDefaultFlags...),
}

func runReleaseCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand()

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

// CmdReleaseDelete represents a sub command of Release to delete a release
var CmdReleaseDelete = cli.Command{
	Name:        "delete",
	Usage:       "Delete a release",
	Description: `Delete a release`,
	ArgsUsage:   "<release tag>",
	Action:      runReleaseDelete,
	Flags:       AllDefaultFlags,
}

func runReleaseDelete(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	client := login.Client()

	tag := ctx.Args().First()
	if len(tag) == 0 {
		fmt.Println("Release tag needed to delete")
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
	return err
}

func getReleaseByTag(owner, repo, tag string, client *gitea.Client) (*gitea.Release, error) {
	rl, _, err := client.ListReleases(owner, repo, gitea.ListReleasesOptions{})
	if err != nil {
		return nil, err
	}
	if len(rl) == 0 {
		fmt.Println("Repo does not have any release")
		return nil, nil
	}
	for _, r := range rl {
		if r.TagName == tag {
			return r, nil
		}
	}
	fmt.Println("Release tag does not exist")
	return nil, nil
}

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
	}, AllDefaultFlags...),
}

func runReleaseEdit(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	client := login.Client()

	tag := ctx.Args().First()
	if len(tag) == 0 {
		fmt.Println("Release tag needed to edit")
		return nil
	}

	release, err := getReleaseByTag(owner, repo, tag, client)
	if err != nil {
		return err
	}
	if release == nil {
		return nil
	}

	var isDraft, isPre *bool
	bTrue := true
	bFalse := false
	if ctx.IsSet("draft") {
		isDraft = &bFalse
		if strings.ToLower(ctx.String("draft"))[:1] == "t" {
			isDraft = &bTrue
		}
	}
	if ctx.IsSet("prerelease") {
		isPre = &bFalse
		if strings.ToLower(ctx.String("prerelease"))[:1] == "t" {
			isPre = &bTrue
		}
	}

	_, _, err = client.EditRelease(owner, repo, release.ID, gitea.EditReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      isDraft,
		IsPrerelease: isPre,
	})
	return err
}
