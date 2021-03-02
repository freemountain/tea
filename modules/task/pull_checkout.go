// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"

	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/go-git/go-git/v5"
	git_plumbing "github.com/go-git/go-git/v5/plumbing"
)

// PullCheckout checkout current workdir to the head branch of specified pull request
func PullCheckout(login *config.Login, repoOwner, repoName string, forceCreateBranch bool, index int64, callback func(string) (string, error)) error {
	client := login.Client()

	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return err
	}

	// fetch PR source-localRepo & -branch from gitea
	pr, _, err := client.GetPullRequest(repoOwner, repoName, index)
	if err != nil {
		return err
	}
	remoteDeleted := pr.Head.Ref == fmt.Sprintf("refs/pull/%d/head", pr.Index)
	if remoteDeleted {
		return fmt.Errorf("Can't checkout: remote head branch was already deleted")
	}

	remoteURL := pr.Head.Repository.CloneURL
	if len(login.SSHKey) != 0 {
		// login.SSHKey is nonempty, if user specified a key manually or we automatically
		// found a matching private key on this machine during login creation.
		// this means, we are very likely to have a working ssh setup.
		remoteURL = pr.Head.Repository.SSHURL
	}

	newRemoteName := fmt.Sprintf("pulls/%v", pr.Head.Repository.Owner.UserName)

	// verify related remote is in local repo, otherwise add it
	localRemote, err := localRepo.GetOrCreateRemote(remoteURL, newRemoteName)
	if err != nil {
		return err
	}
	localRemoteName := localRemote.Config().Name

	// get auth & fetch remote via its configured protocol
	url, err := localRepo.TeaRemoteURL(localRemoteName)
	if err != nil {
		return err
	}
	auth, err := local_git.GetAuthForURL(url, login.Token, login.SSHKey, callback)
	if err != nil {
		return err
	}
	fmt.Printf("Fetching PR %v (head %s:%s) from remote '%s'\n", index, url, pr.Head.Ref, localRemoteName)
	err = localRemote.Fetch(&git.FetchOptions{Auth: auth})
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(err)
	} else if err != nil {
		return err
	}

	var info string
	var checkoutRef git_plumbing.ReferenceName

	if b, _ := localRepo.TeaFindBranchBySha(pr.Head.Sha, remoteURL); b != nil {
		// if a matching branch exists, use that
		checkoutRef = git_plumbing.NewBranchReferenceName(b.Name)
		info = fmt.Sprintf("Found matching local branch %s, checking it out", checkoutRef.Short())
	} else if forceCreateBranch {
		// create a branch if wanted
		localBranchName := fmt.Sprintf("pulls/%v-%v", index, pr.Head.Ref)
		checkoutRef = git_plumbing.NewBranchReferenceName(localBranchName)
		if err = localRepo.TeaCreateBranch(localBranchName, pr.Head.Ref, localRemoteName); err == nil {
			info = fmt.Sprintf("Created branch '%s'\n", localBranchName)
		} else if err == git.ErrBranchExists {
			info = "There may be changes since you last checked out, run `git pull` to get them."
		} else {
			return err
		}
	} else {
		// use the remote tracking branch
		checkoutRef = git_plumbing.NewRemoteReferenceName(localRemoteName, pr.Head.Ref)
		info = fmt.Sprintf(
			"Checking out remote tracking branch %s. To make changes, create a new branch:\n  git checkout %s",
			checkoutRef.String(), pr.Head.Ref)
	}

	fmt.Println(info)
	return localRepo.TeaCheckout(checkoutRef)
}
