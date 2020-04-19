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
	"github.com/urfave/cli/v2"
	"gopkg.in/src-d/go-git.v4"
	git_config "gopkg.in/src-d/go-git.v4/config"
)

// CmdPulls is the main command to operate on PRs
var CmdPulls = cli.Command{
	Name:        "pulls",
	Aliases:     []string{"pull", "pr"},
	Usage:       "List open pull requests",
	Description: `List open pull requests`,
	Action:      runPulls,
	Flags: append([]cli.Flag{
		&cli.StringFlag{
			Name:        "state",
			Usage:       "Filter by PR state (all|open|closed)",
			DefaultText: "open",
		},
	}, AllDefaultFlags...),
	Subcommands: []*cli.Command{
		&CmdPullsCheckout,
		&CmdPullsClean,
	},
}

func runPulls(ctx *cli.Context) error {
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

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(state),
	})

	if err != nil {
		log.Fatal(err)
	}

	headers := []string{
		"Index",
		"State",
		"Author",
		"Updated",
		"Title",
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
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		values = append(
			values,
			[]string{
				strconv.FormatInt(pr.Index, 10),
				string(pr.State),
				name,
				pr.Updated.Format("2006-01-02 15:04:05"),
				pr.Title,
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
	pr, err := login.Client().GetPullRequest(owner, repo, idx)
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
	pr, err := login.Client().GetPullRequest(owner, repo, idx)
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

func argToIndex(arg string) (int64, error) {
	if strings.HasPrefix(arg, "#") {
		arg = arg[1:]
	}
	return strconv.ParseInt(arg, 10, 64)
}
