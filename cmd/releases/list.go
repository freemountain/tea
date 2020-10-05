// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package releases

import (
	"fmt"
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReleaseList represents a sub command of Release to list releases
var CmdReleaseList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List Releases",
	Description: "List Releases",
	Action:      RunReleasesList,
	Flags: append([]cli.Flag{
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.AllDefaultFlags...),
}

// RunReleasesList list releases
func RunReleasesList(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)

	releases, _, err := login.Client().ListReleases(owner, repo, gitea.ListReleasesOptions{ListOptions: flags.GetListOptions(ctx)})
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
		print.OutputList(flags.GlobalOutputValue, headers, values)
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
				print.FormatTime(release.PublishedAt),
				status,
				release.TarURL,
			},
		)
	}
	print.OutputList(flags.GlobalOutputValue, headers, values)

	return nil
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
