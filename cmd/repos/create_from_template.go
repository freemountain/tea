// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepoCreateFromTemplate represents a sub command of repos to generate one from a template repo
var CmdRepoCreateFromTemplate = cli.Command{
	Name:        "create-from-template",
	Aliases:     []string{"ct"},
	Usage:       "Create a repository based on an existing template",
	Description: "Create a repository based on an existing template",
	Action:      runRepoCreateFromTemplate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "template",
			Aliases:  []string{"t"},
			Required: true,
			Usage:    "source template to copy from",
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Required: true,
			Usage:    "name of new repo",
		},
		&cli.StringFlag{
			Name:    "owner",
			Aliases: []string{"O"},
			Usage:   "name of repo owner",
		},
		&cli.BoolFlag{
			Name:  "private",
			Usage: "make new repo private",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"desc"},
			Usage:   "add custom description to repo",
		},
		&cli.BoolFlag{
			Name:  "content",
			Value: true,
			Usage: "copy git content from template",
		},
		&cli.BoolFlag{
			Name:  "githooks",
			Value: true,
			Usage: "copy git hooks from template",
		},
		&cli.BoolFlag{
			Name:  "avatar",
			Value: true,
			Usage: "copy repo avatar from template",
		},
		&cli.BoolFlag{
			Name:  "labels",
			Value: true,
			Usage: "copy repo labels from template",
		},
		&cli.BoolFlag{
			Name:  "topics",
			Value: true,
			Usage: "copy topics from template",
		},
		&cli.BoolFlag{
			Name:  "webhooks",
			Usage: "copy webhooks from template",
		},
	}, flags.LoginOutputFlags...),
}

func runRepoCreateFromTemplate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()

	templateOwner, templateRepo := utils.GetOwnerAndRepo(ctx.String("template"), ctx.Login.User)
	owner := ctx.Login.User
	if ctx.IsSet("owner") {
		owner = ctx.String("owner")
	}

	opts := gitea.CreateRepoFromTemplateOption{
		Name:        ctx.String("name"),
		Owner:       owner,
		Description: ctx.String("description"),
		Private:     ctx.Bool("private"),
		GitContent:  ctx.Bool("content"),
		GitHooks:    ctx.Bool("githooks"),
		Avatar:      ctx.Bool("avatar"),
		Labels:      ctx.Bool("labels"),
		Topics:      ctx.Bool("topics"),
		Webhooks:    ctx.Bool("webhooks"),
	}

	repo, _, err := client.CreateRepoFromTemplate(templateOwner, templateRepo, opts)
	if err != nil {
		return err
	}

	topics, _, err := client.ListRepoTopics(repo.Owner.UserName, repo.Name, gitea.ListRepoTopicsOptions{})
	if err != nil {
		return err
	}
	print.RepoDetails(repo, topics)

	fmt.Printf("%s\n", repo.HTMLURL)
	return nil
}
