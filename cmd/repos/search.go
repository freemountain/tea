// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdReposSearch represents a sub command of repos to find them
var CmdReposSearch = cli.Command{
	Name:        "search",
	Aliases:     []string{"s"},
	Usage:       "Find any repo on an Gitea instance",
	Description: "Find any repo on an Gitea instance",
	ArgsUsage:   "[<search term>]",
	Action:      runReposSearch,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			// TODO: it might be nice to search for topics as an ADDITIONAL filter.
			// for that, we'd probably need to make multiple queries and UNION the results.
			Name:     "topic",
			Aliases:  []string{"t"},
			Required: false,
			Usage:    "Search for term in repo topics instead of name",
		},
		&typeFilterFlag,
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
		flags.FieldsFlag(print.RepoFields, []string{
			"owner", "name", "type", "ssh",
		}),
		&flags.PaginationPageFlag,
		&flags.PaginationLimitFlag,
	}, flags.LoginOutputFlags...),
}

func runReposSearch(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()

	var ownerID int64
	if ctx.IsSet("owner") {
		// test if owner is a organisation
		org, _, err := client.GetOrg(ctx.String("owner"))
		if err != nil {
			// HACK: the client does not return a response on 404, so we can't check res.StatusCode
			if err.Error() != "404 Not Found" {
				return fmt.Errorf("Could not find owner: %s", err)
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
	}

	var isArchived *bool
	if ctx.IsSet("archived") {
		archived := strings.ToLower(ctx.String("archived"))[:1] == "t"
		isArchived = &archived
	}

	var isPrivate *bool
	if ctx.IsSet("private") {
		private := strings.ToLower(ctx.String("private"))[:1] == "t"
		isPrivate = &private
	}

	mode, err := getTypeFilter(cmd)
	if err != nil {
		return err
	}

	var keyword string
	if ctx.Args().Present() {
		keyword = strings.Join(ctx.Args().Slice(), " ")
	}

	user, _, err := client.GetMyUserInfo()
	if err != nil {
		return err
	}

	rps, _, err := client.SearchRepos(gitea.SearchRepoOptions{
		ListOptions:          ctx.GetListOptions(),
		OwnerID:              ownerID,
		IsPrivate:            isPrivate,
		IsArchived:           isArchived,
		Type:                 mode,
		Keyword:              keyword,
		KeywordInDescription: true,
		KeywordIsTopic:       ctx.Bool("topic"),
		PrioritizedByOwnerID: user.ID,
	})
	if err != nil {
		return err
	}

	fields, err := flags.GetFields(cmd, nil)
	if err != nil {
		return err
	}
	print.ReposList(rps, ctx.Output, fields)
	return nil
}
