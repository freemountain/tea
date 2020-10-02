// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"log"
	"net/http"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReposList represents a sub command of repos to list them
var CmdReposList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      RunReposList,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "mode",
			Aliases:  []string{"m"},
			Required: false,
			Usage:    "Filter by mode: fork, mirror, source",
		},
		&cli.StringFlag{
			Name:     "owner",
			Aliases:  []string{"O"},
			Required: false,
			Usage:    "Filter by owner",
		},
		&cli.StringFlag{
			Name:     "private",
			Required: false,
			Usage:    "Filter private repos (true|false)",
		},
		&cli.StringFlag{
			Name:     "archived",
			Required: false,
			Usage:    "Filter archived repos (true|false)",
		},
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.LoginOutputFlags...),
}

// RunReposList list repositories
func RunReposList(ctx *cli.Context) error {
	login := config.InitCommandLoginOnly(flags.GlobalLoginValue)
	client := login.Client()

	var ownerID int64
	if ctx.IsSet("owner") {
		// test if owner is a organisation
		org, resp, err := client.GetOrg(ctx.String("owner"))
		if err != nil {
			if resp == nil || resp.StatusCode != http.StatusNotFound {
				return err
			}
			// if owner is no org, its a user
			user, _, err := client.GetUserInfo(ctx.String("owner"))
			if err != nil {
				return err
			}
			ownerID = user.ID
		} else {
			ownerID = org.ID
		}
	} else {
		me, _, err := client.GetMyUserInfo()
		if err != nil {
			return err
		}
		ownerID = me.ID
	}

	var isArchived *bool
	if ctx.IsSet("archived") {
		archived := strings.ToLower(ctx.String("archived"))[:1] == "t"
		isArchived = &archived
	}

	var isPrivate *bool
	if ctx.IsSet("private") {
		private := strings.ToLower(ctx.String("private"))[:1] == "t"
		isArchived = &private
	}

	mode := gitea.RepoTypeNone
	switch ctx.String("mode") {
	case "fork":
		mode = gitea.RepoTypeFork
	case "mirror":
		mode = gitea.RepoTypeMirror
	case "source":
		mode = gitea.RepoTypeSource
	}

	rps, _, err := client.SearchRepos(gitea.SearchRepoOptions{
		ListOptions: flags.GetListOptions(ctx),
		OwnerID:     ownerID,
		IsPrivate:   isPrivate,
		IsArchived:  isArchived,
		Type:        mode,
	})
	if err != nil {
		return err
	}

	if len(rps) == 0 {
		log.Fatal("No repositories found", rps)
		return nil
	}

	headers := []string{
		"Name",
		"Type",
		"SSH",
		"Owner",
	}
	var values [][]string

	for _, rp := range rps {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}

		values = append(
			values,
			[]string{
				rp.FullName,
				mode,
				rp.SSHURL,
				rp.Owner.UserName,
			},
		)
	}
	print.OutputList(flags.GlobalOutputValue, headers, values)

	return nil
}
