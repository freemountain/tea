// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReposListFlags contains all flags needed for repo listing
var CmdReposListFlags = append([]cli.Flag{
	&cli.BoolFlag{
		Name:     "watched",
		Aliases:  []string{"w"},
		Required: false,
		Usage:    "List your watched repos instead",
	},
	&cli.BoolFlag{
		Name:     "starred",
		Aliases:  []string{"s"},
		Required: false,
		Usage:    "List your starred repos instead",
	},
	&printFieldsFlag,
	&typeFilterFlag,
	&flags.PaginationPageFlag,
	&flags.PaginationLimitFlag,
}, flags.LoginOutputFlags...)

// CmdReposList represents a sub command of repos to list them
var CmdReposList = cli.Command{
	Name:        "ls",
	Aliases:     []string{"list"},
	Usage:       "List repositories you have access to",
	Description: "List repositories you have access to",
	Action:      RunReposList,
	Flags:       CmdReposListFlags,
}

// RunReposList list repositories
func RunReposList(ctx *cli.Context) error {
	login, _, _ := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	typeFilter, err := getTypeFilter(ctx)
	if err != nil {
		return err
	}

	var rps []*gitea.Repository
	if ctx.Bool("starred") {
		user, _, err := client.GetMyUserInfo()
		if err != nil {
			return err
		}
		rps, _, err = client.SearchRepos(gitea.SearchRepoOptions{
			ListOptions:     flags.GetListOptions(ctx),
			StarredByUserID: user.ID,
		})
	} else if ctx.Bool("watched") {
		rps, _, err = client.GetMyWatchedRepos() // TODO: this does not expose pagination..
	} else {
		rps, _, err = client.ListMyRepos(gitea.ListReposOptions{
			ListOptions: flags.GetListOptions(ctx),
		})
	}

	if err != nil {
		return err
	}

	reposFiltered := rps
	if typeFilter != gitea.RepoTypeNone {
		reposFiltered = filterReposByType(rps, typeFilter)
	}

	print.ReposList(reposFiltered, flags.GlobalOutputValue, getFields(ctx))
	return nil
}

func filterReposByType(repos []*gitea.Repository, t gitea.RepoType) []*gitea.Repository {
	var filtered []*gitea.Repository
	for _, r := range repos {
		switch t {
		case gitea.RepoTypeFork:
			if !r.Fork {
				continue
			}
		case gitea.RepoTypeMirror:
			if !r.Mirror {
				continue
			}
		case gitea.RepoTypeSource:
			if r.Fork || r.Mirror {
				continue
			}
		}

		filtered = append(filtered, r)
	}
	return filtered
}
