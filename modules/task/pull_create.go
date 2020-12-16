// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"github.com/go-git/go-git/v5"
)

// CreatePull creates a PR in the given repo and prints the result
func CreatePull(login *config.Login, repoOwner, repoName, base, head, title, description string) error {

	// open local git repo
	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return fmt.Errorf("Could not open local repo: %s", err)
	}

	// push if possible
	fmt.Println("git push")
	err = localRepo.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("Error occurred during 'git push':\n%s\n", err.Error())
	}

	// default is default branch
	if len(base) == 0 {
		base, err = GetDefaultPRBase(login, repoOwner, repoName)
		if err != nil {
			return err
		}
	}

	// default is current one
	if len(head) == 0 {
		headOwner, headBranch, err := GetDefaultPRHead(localRepo)
		if err != nil {
			return err
		}

		head = GetHeadSpec(headOwner, headBranch, repoOwner)
	}

	// head & base may not be the same
	if head == base {
		return fmt.Errorf("can't create PR from %s to %s", head, base)
	}

	// default is head branch name
	if len(title) == 0 {
		title = GetDefaultPRTitle(head)
	}
	// title is required
	if len(title) == 0 {
		return fmt.Errorf("Title is required")
	}

	pr, _, err := login.Client().CreatePullRequest(repoOwner, repoName, gitea.CreatePullRequestOption{
		Head:  head,
		Base:  base,
		Title: title,
		Body:  description,
	})

	if err != nil {
		return fmt.Errorf("Could not create PR from %s to %s:%s: %s", head, repoOwner, base, err)
	}

	print.PullDetails(pr, nil, nil)

	fmt.Println(pr.HTMLURL)

	return err
}

// GetDefaultPRBase retrieves the default base branch for the given repo
func GetDefaultPRBase(login *config.Login, owner, repo string) (string, error) {
	meta, _, err := login.Client().GetRepo(owner, repo)
	if err != nil {
		return "", fmt.Errorf("could not fetch repo meta: %s", err)
	}
	return meta.DefaultBranch, nil
}

// GetDefaultPRHead uses the currently checked out branch, checks if
// a remote currently holds the commit it points to, extracts the owner
// from its URL, and assembles the result to a valid head spec for gitea.
func GetDefaultPRHead(localRepo *local_git.TeaRepo) (owner, branch string, err error) {
	headBranch, err := localRepo.Head()
	if err != nil {
		return
	}
	sha := headBranch.Hash().String()

	remote, err := localRepo.TeaFindBranchRemote("", sha)
	if err != nil {
		err = fmt.Errorf("could not determine remote for current branch: %s", err)
		return
	}

	if remote == nil {
		// if no remote branch is found for the local hash, we abort:
		// user has probably not configured a remote for the local branch,
		// or local branch does not represent remote state.
		err = fmt.Errorf("no matching remote found for this branch. try git push -u <remote> <branch>")
		return
	}

	branch, err = localRepo.TeaGetCurrentBranchName()
	if err != nil {
		return
	}

	url, err := local_git.ParseURL(remote.Config().URLs[0])
	if err != nil {
		return
	}
	owner, _ = utils.GetOwnerAndRepo(strings.TrimLeft(url.Path, "/"), "")
	return
}

// GetHeadSpec creates a head string as expected by gitea API
func GetHeadSpec(owner, branch, baseOwner string) string {
	if len(owner) != 0 && owner != baseOwner {
		return fmt.Sprintf("%s:%s", owner, branch)
	}
	return branch
}

// GetDefaultPRTitle transforms a string like a branchname to a readable text
func GetDefaultPRTitle(head string) string {
	title := head
	if strings.Contains(title, ":") {
		title = strings.SplitN(title, ":", 2)[1]
	}
	title = strings.Replace(title, "-", " ", -1)
	title = strings.Replace(title, "_", " ", -1)
	title = strings.Title(strings.ToLower(title))
	return title
}
