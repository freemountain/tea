// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	local_git "code.gitea.io/tea/modules/git"

	"code.gitea.io/sdk/gitea"
	"github.com/charmbracelet/glamour"
	"github.com/go-git/go-git/v5"
	git_config "github.com/go-git/go-git/v5/config"
	"github.com/urfave/cli/v2"
)

// CmdPulls is the main command to operate on PRs
var CmdPulls = cli.Command{
	Name:        "pulls",
	Aliases:     []string{"pull", "pr"},
	Usage:       "List, create, checkout and clean pull requests",
	Description: `List, create, checkout and clean pull requests`,
	ArgsUsage:   "[<pull index>]",
	Action:      runPulls,
	Flags:       IssuePRFlags,
	Subcommands: []*cli.Command{
		&CmdPullsList,
		&CmdPullsCheckout,
		&CmdPullsClean,
		&CmdPullsCreate,
	},
}

func runPulls(ctx *cli.Context) error {
	if ctx.Args().Len() == 1 {
		return runPullDetail(ctx, ctx.Args().First())
	}
	return runPullsList(ctx)
}

// CmdPullsList represents a sub command of issues to list pulls
var CmdPullsList = cli.Command{
	Name:        "ls",
	Usage:       "List pull requests of the repository",
	Description: `List pull requests of the repository`,
	Action:      runPullsList,
	Flags:       IssuePRFlags,
}

func runPullDetail(ctx *cli.Context, index string) error {
	login, owner, repo := initCommand()

	idx, err := argToIndex(index)
	if err != nil {
		return err
	}
	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}

	// TODO: use glamour once #181 is merged
	fmt.Printf("#%d %s\n%s created %s\n\n%s\n", pr.Index,
		pr.Title,
		pr.Poster.UserName,
		pr.Created.Format("2006-01-02 15:04:05"),
		pr.Body,
	)
	return nil
}

func runPullsList(ctx *cli.Context) error {
	login, owner, repo := initCommand()

	state := gitea.StateOpen
	switch ctx.String("state") {
	case "all":
		state = gitea.StateAll
	case "open":
		state = gitea.StateOpen
	case "closed":
		state = gitea.StateClosed
	}

	prs, _, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		State: state,
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	}

	var values [][]string

	if len(prs) == 0 {
		Output(outputValue, headers, values)
		return nil
	}

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		author := pr.Poster.FullName
		if len(author) == 0 {
			author = pr.Poster.UserName
		}
		mile := ""
		if pr.Milestone != nil {
			mile = pr.Milestone.Title
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(pr.Index, 10),
				pr.Title,
				string(pr.State),
				author,
				mile,
				pr.Updated.Format("2006-01-02 15:04:05"),
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}

// CmdPullsCheckout is a command to locally checkout the given PR
var CmdPullsCheckout = cli.Command{
	Name:        "checkout",
	Usage:       "Locally check out the given PR",
	Description: `Locally check out the given PR`,
	Action:      runPullsCheckout,
	ArgsUsage:   "<pull index>",
	Flags:       AllDefaultFlags,
}

func runPullsCheckout(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	if ctx.Args().Len() != 1 {
		log.Fatal("Must specify a PR index")
	}

	// fetch PR source-repo & -branch from gitea
	idx, err := argToIndex(ctx.Args().First())
	if err != nil {
		return err
	}
	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}
	remoteURL := pr.Head.Repository.CloneURL
	remoteBranchName := pr.Head.Ref

	// open local git repo
	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return nil
	}

	// verify related remote is in local repo, otherwise add it
	newRemoteName := fmt.Sprintf("pulls/%v", pr.Head.Repository.Owner.UserName)
	localRemote, err := localRepo.GetOrCreateRemote(remoteURL, newRemoteName)
	if err != nil {
		return err
	}

	localRemoteName := localRemote.Config().Name
	localBranchName := fmt.Sprintf("pulls/%v-%v", idx, remoteBranchName)

	// fetch remote
	fmt.Printf("Fetching PR %v (head %s:%s) from remote '%s'\n",
		idx, remoteURL, remoteBranchName, localRemoteName)

	url, err := local_git.ParseURL(localRemote.Config().URLs[0])
	if err != nil {
		return err
	}
	auth, err := local_git.GetAuthForURL(url, login.User, login.SSHKey)
	if err != nil {
		return err
	}

	err = localRemote.Fetch(&git.FetchOptions{Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(err)
	} else if err != nil {
		return err
	}

	// checkout local branch
	fmt.Printf("Creating branch '%s'\n", localBranchName)
	err = localRepo.TeaCreateBranch(localBranchName, remoteBranchName, localRemoteName)
	if err == git.ErrBranchExists {
		fmt.Println(err)
	} else if err != nil {
		return err
	}

	fmt.Printf("Checking out PR %v\n", idx)
	err = localRepo.TeaCheckout(localBranchName)

	return err
}

// CmdPullsClean removes the remote and local feature branches, if a PR is merged.
var CmdPullsClean = cli.Command{
	Name:        "clean",
	Usage:       "Deletes local & remote feature-branches for a closed pull request",
	Description: `Deletes local & remote feature-branches for a closed pull request`,
	ArgsUsage:   "<pull index>",
	Action:      runPullsClean,
	Flags: append([]cli.Flag{
		&cli.BoolFlag{
			Name:  "ignore-sha",
			Usage: "Find the local branch by name instead of commit hash (less precise)",
		},
	}, AllDefaultFlags...),
}

func runPullsClean(ctx *cli.Context) error {
	login, owner, repo := initCommand()
	if ctx.Args().Len() != 1 {
		return fmt.Errorf("Must specify a PR index")
	}

	// fetch PR source-repo & -branch from gitea
	idx, err := argToIndex(ctx.Args().First())
	if err != nil {
		return err
	}
	pr, _, err := login.Client().GetPullRequest(owner, repo, idx)
	if err != nil {
		return err
	}
	if pr.State == gitea.StateOpen {
		return fmt.Errorf("PR is still open, won't delete branches")
	}

	// IDEA: abort if PR.Head.Repository.CloneURL does not match login.URL?

	r, err := local_git.RepoForWorkdir()
	if err != nil {
		return err
	}

	// find a branch with matching sha or name, that has a remote matching the repo url
	var branch *git_config.Branch
	if ctx.Bool("ignore-sha") {
		branch, err = r.TeaFindBranchByName(pr.Head.Ref, pr.Head.Repository.CloneURL)
	} else {
		branch, err = r.TeaFindBranchBySha(pr.Head.Sha, pr.Head.Repository.CloneURL)
	}
	if err != nil {
		return err
	}
	if branch == nil {
		if ctx.Bool("ignore-sha") {
			return fmt.Errorf("Remote branch %s not found in local repo", pr.Head.Ref)
		}
		return fmt.Errorf(`Remote branch %s not found in local repo.
Either you don't track this PR, or the local branch has diverged from the remote.
If you still want to continue & are sure you don't loose any important commits,
call me again with the --ignore-sha flag`, pr.Head.Ref)
	}

	// prepare deletion of local branch:
	headRef, err := r.Head()
	if err != nil {
		return err
	}
	if headRef.Name().Short() == branch.Name {
		fmt.Printf("Checking out 'master' to delete local branch '%s'\n", branch.Name)
		err = r.TeaCheckout("master")
		if err != nil {
			return err
		}
	}

	// remove local & remote branch
	fmt.Printf("Deleting local branch %s and remote branch %s\n", branch.Name, pr.Head.Ref)
	url, err := r.TeaRemoteURL(branch.Remote)
	if err != nil {
		return err
	}
	auth, err := local_git.GetAuthForURL(url, login.User, login.SSHKey)
	if err != nil {
		return err
	}
	return r.TeaDeleteBranch(branch, pr.Head.Ref, auth)
}

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
	}, AllDefaultFlags...),
}

func runPullsCreate(ctx *cli.Context) error {
	login, ownerArg, repoArg := initCommand()
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
		owner, _ := getOwnerAndRepo(strings.TrimLeft(url.Path, "/"), "")
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

	in := fmt.Sprintf("# #%d %s (%s)\n%s created %s\n\n%s\n", pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		pr.Created.Format("2006-01-02 15:04:05"),
		pr.Body,
	)
	out, err := glamour.Render(in, getGlamourTheme())
	fmt.Print(out)

	fmt.Println(pr.HTMLURL)
	return err
}

func argToIndex(arg string) (int64, error) {
	if strings.HasPrefix(arg, "#") {
		arg = arg[1:]
	}
	return strconv.ParseInt(arg, 10, 64)
}
