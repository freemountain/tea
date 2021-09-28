// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepoCreate represents a sub command of repos to create one
var CmdRepoCreate = cli.Command{
	Name:        "create",
	Aliases:     []string{"c"},
	Usage:       "Create a repository",
	Description: "Create a repository",
	Action:      runRepoCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{""},
			Required: true,
			Usage:    "name of new repo",
		},
		&cli.StringFlag{
			Name:     "owner",
			Aliases:  []string{"O"},
			Required: false,
			Usage:    "name of repo owner",
		},
		&cli.BoolFlag{
			Name:     "private",
			Required: false,
			Value:    false,
			Usage:    "make repo private",
		},
		&cli.StringFlag{
			Name:     "description",
			Aliases:  []string{"desc"},
			Required: false,
			Usage:    "add description to repo",
		},
		&cli.BoolFlag{
			Name:     "init",
			Required: false,
			Value:    false,
			Usage:    "initialize repo",
		},
		&cli.StringFlag{
			Name:     "labels",
			Required: false,
			Usage:    "name of label set to add",
		},
		&cli.StringFlag{
			Name:     "gitignores",
			Aliases:  []string{"git"},
			Required: false,
			Usage:    "list of gitignore templates (need --init)",
		},
		&cli.StringFlag{
			Name:     "license",
			Required: false,
			Usage:    "add license (need --init)",
		},
		&cli.StringFlag{
			Name:     "readme",
			Required: false,
			Usage:    "use readme template (need --init)",
		},
		&cli.StringFlag{
			Name:     "branch",
			Required: false,
			Usage:    "use custom default branch (need --init)",
		},
		&cli.BoolFlag{
			Name:  "template",
			Usage: "make repo a template repo",
		},
		&cli.StringFlag{
			Name:  "trustmodel",
			Usage: "select trust model (committer,collaborator,collaborator+committer)",
		},
	}, flags.LoginOutputFlags...),
}

func runRepoCreate(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	client := ctx.Login.Client()
	var (
		repo       *gitea.Repository
		err        error
		trustmodel gitea.TrustModel
	)

	if ctx.IsSet("trustmodel") {
		switch ctx.String("trustmodel") {
		case "committer":
			trustmodel = gitea.TrustModelCommitter
		case "collaborator":
			trustmodel = gitea.TrustModelCollaborator
		case "collaborator+committer":
			trustmodel = gitea.TrustModelCollaboratorCommitter
		default:
			return fmt.Errorf("unknown trustmodel type '%s'", ctx.String("trustmodel"))
		}
	}

	opts := gitea.CreateRepoOption{
		Name:          ctx.String("name"),
		Description:   ctx.String("description"),
		Private:       ctx.Bool("private"),
		AutoInit:      ctx.Bool("init"),
		IssueLabels:   ctx.String("labels"),
		Gitignores:    ctx.String("gitignores"),
		License:       ctx.String("license"),
		Readme:        ctx.String("readme"),
		DefaultBranch: ctx.String("branch"),
		Template:      ctx.Bool("template"),
		TrustModel:    trustmodel,
	}
	if len(ctx.String("owner")) != 0 {
		repo, _, err = client.CreateOrgRepo(ctx.String("owner"), opts)
	} else {
		repo, _, err = client.CreateRepo(opts)
	}
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
