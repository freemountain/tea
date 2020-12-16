// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package releases

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
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
func RunReleasesList(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	releases, _, err := ctx.Login.Client().ListReleases(ctx.Owner, ctx.Repo, gitea.ListReleasesOptions{
		ListOptions: ctx.GetListOptions(),
	})
	if err != nil {
		return err
	}

	print.ReleasesList(releases, ctx.Output)
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
