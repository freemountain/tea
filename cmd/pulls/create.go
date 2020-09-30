// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"
	"log"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

// CmdPullsCreate creates a pull request
var CmdPullsCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a pull-request",
	Description: "Create a pull-request",
	Action:      runPullsCreate,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:  "head",
			Usage: "Set head branch (default is current one)",
		},
		&cli.StringFlag{
			Name:    "base",
			Aliases: []string{"b"},
			Usage:   "Set base branch (default is default branch)",
		},
		&cli.StringFlag{
			Name:    "title",
			Aliases: []string{"t"},
			Usage:   "Set title of pull (default is head branch name)",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"d"},
			Usage:   "Set body of new pull",
		},
	}, flags.AllDefaultFlags...),
}

func runPullsCreate(ctx *cli.Context) error {
	login, ownerArg, repoArg := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	client := login.Client()

	repo, _, err := client.GetRepo(ownerArg, repoArg)
	if err != nil {
		log.Fatal("could not fetch repo meta: ", err)
	}

	// open local git repo
	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		log.Fatal("could not open local repo: ", err)
	}

	// push if possible
	log.Println("git push")
	err = localRepo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Printf("Error occurred during 'git push':\n%s\n", err.Error())
	}

	base := ctx.String("base")
	// default is default branch
	if len(base) == 0 {
		base = repo.DefaultBranch
	}

	head := ctx.String("head")
	// default is current one
	if len(head) == 0 {
		headBranch, err := localRepo.Head()
		if err != nil {
			log.Fatal(err)
		}
		sha := headBranch.Hash().String()

		remote, err := localRepo.TeaFindBranchRemote("", sha)
		if err != nil {
			log.Fatal("could not determine remote for current branch: ", err)
		}

		if remote == nil {
			// if no remote branch is found for the local hash, we abort:
			// user has probably not configured a remote for the local branch,
			// or local branch does not represent remote state.
			log.Fatal("no matching remote found for this branch. try git push -u <remote> <branch>")
		}

		branchName, err := localRepo.TeaGetCurrentBranchName()
		if err != nil {
			log.Fatal(err)
		}

		url, err := local_git.ParseURL(remote.Config().URLs[0])
		if err != nil {
			log.Fatal(err)
		}
		owner, _ := config.GetOwnerAndRepo(strings.TrimLeft(url.Path, "/"), "")
		head = fmt.Sprintf("%s:%s", owner, branchName)
	}

	title := ctx.String("title")
	// default is head branch name
	if len(title) == 0 {
		title = head
		if strings.Contains(title, ":") {
			title = strings.SplitN(title, ":", 2)[1]
		}
		title = strings.Replace(title, "-", " ", -1)
		title = strings.Replace(title, "_", " ", -1)
		title = strings.Title(strings.ToLower(title))
	}
	// title is required
	if len(title) == 0 {
		fmt.Printf("Title is required")
		return nil
	}

	pr, _, err := client.CreatePullRequest(ownerArg, repoArg, gitea.CreatePullRequestOption{
		Head:  head,
		Base:  base,
		Title: title,
		Body:  ctx.String("description"),
	})

	if err != nil {
		log.Fatalf("could not create PR from %s to %s:%s: %s", head, ownerArg, base, err)
	}

	print.PullDetails(pr)

	fmt.Println(pr.HTMLURL)
	return err
}
