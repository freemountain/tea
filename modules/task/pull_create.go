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
)

// CreatePull creates a PR in the given repo and prints the result
func CreatePull(login *config.Login, repoOwner, repoName, base, head string, opts *gitea.CreateIssueOption) error {
	// open local git repo
	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return fmt.Errorf("Could not open local repo: %s", err)
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
	if len(opts.Title) == 0 {
		opts.Title = GetDefaultPRTitle(head)
	}
	// title is required
	if len(opts.Title) == 0 {
		return fmt.Errorf("Title is required")
	}

	pr, _, err := login.Client().CreatePullRequest(repoOwner, repoName, gitea.CreatePullRequestOption{
		Head:      head,
		Base:      base,
		Title:     opts.Title,
		Body:      opts.Body,
		Assignees: opts.Assignees,
		Labels:    opts.Labels,
		Milestone: opts.Milestone,
		Deadline:  opts.Deadline,
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

// GetDefaultPRHead uses the currently checked out branch, tries to find a remote
// that has a branch with the same name, and extracts the owner from its URL.
// If no remote matches, owner is empty, meaning same as head repo owner.
func GetDefaultPRHead(localRepo *local_git.TeaRepo) (owner, branch string, err error) {
	if branch, err = localRepo.TeaGetCurrentBranchName(); err != nil {
		return
	}

	remote, err := localRepo.TeaFindBranchRemote(branch, "")
	if err != nil {
		err = fmt.Errorf("could not determine remote for current branch: %s", err)
		return
	}

	if remote == nil {
		// if no remote branch is found for the local branch,
		// we leave owner empty, meaning "use same repo as head" to gitea.
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
