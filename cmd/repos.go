// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Show repositories details",
	Description: "Show repositories details",
	ArgsUsage:   "[<repo owner>/<repo name>]",
	Action:      runRepos,
	Subcommands: []*cli.Command{
		&CmdReposList,
		&CmdRepoCreate,
	},
	Flags: LoginOutputFlags,
}

// CmdReposList represents a sub command of repos to list them
var CmdReposList = cli.Command{
	Name:        "ls",
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      runReposList,
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
	}, LoginOutputFlags...),
}

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
	}, LoginOutputFlags...),
}

func runRepos(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runRepoDetail(ctx, ctx.Args().First())
	}
	return runReposList(ctx)
}

// runReposList list repositories
func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly()
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
		OwnerID:    ownerID,
		IsPrivate:  isPrivate,
		IsArchived: isArchived,
		Type:       mode,
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
	Output(outputValue, headers, values)

	return nil
}

func runRepoDetail(_ *cli.Context, path string) error {
	login := initCommandLoginOnly()
	client := login.Client()
	repoOwner, repoName := getOwnerAndRepo(path, login.User)
	repo, _, err := client.GetRepo(repoOwner, repoName)
	if err != nil {
		return err
	}
	topics, _, err := client.ListRepoTopics(repo.Owner.UserName, repo.Name, gitea.ListRepoTopicsOptions{})
	if err != nil {
		return err
	}

	output := repo.FullName
	if repo.Mirror {
		output += " (mirror)"
	}
	if repo.Fork {
		output += " (fork)"
	}
	if repo.Archived {
		output += " (archived)"
	}
	if repo.Empty {
		output += " (empty)"
	}
	output += "\n"
	if len(topics) != 0 {
		output += "Topics: " + strings.Join(topics, ", ") + "\n"
	}
	output += "\n"
	output += repo.Description + "\n\n"
	output += fmt.Sprintf(
		"Open Issues: %d, Stars: %d, Forks: %d, Size: %s\n\n",
		repo.OpenIssues,
		repo.Stars,
		repo.Forks,
		utils.FormatSize(int64(repo.Size)),
	)

	fmt.Print(output)
	return nil
}

func runRepoCreate(ctx *cli.Context) error {
	login := initCommandLoginOnly()
	client := login.Client()
	var (
		repo *gitea.Repository
		err  error
	)
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
	}
	if len(ctx.String("owner")) != 0 {
		repo, _, err = client.CreateOrgRepo(ctx.String("owner"), opts)
	} else {
		repo, _, err = client.CreateRepo(opts)
	}
	if err != nil {
		return err
	}
	if err = runRepoDetail(ctx, repo.FullName); err != nil {
		return err
	}
	fmt.Printf("%s\n", repo.HTMLURL)
	return nil
}
