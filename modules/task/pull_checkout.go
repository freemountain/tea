// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/go-git/go-git/v5"
)

// PullCheckout checkout current workdir to the head branch of specified pull request
func PullCheckout(login *config.Login, repoOwner, repoName string, index int64) error {
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

	// test if we can pull via SSH, and configure git remote accordingly
	remoteURL := pr.Head.Repository.CloneURL
	keys, _, err := client.ListMyPublicKeys(gitea.ListPublicKeysOptions{})
	if err != nil {
		return err
	}
	if len(keys) != 0 {
		remoteURL = pr.Head.Repository.SSHURL
	}

	// try to find a matching existing branch, otherwise return branch in pulls/ namespace
	localBranchName := fmt.Sprintf("pulls/%v-%v", index, pr.Head.Ref)
	if b, _ := localRepo.TeaFindBranchBySha(pr.Head.Sha, remoteURL); b != nil {
		localBranchName = b.Name
	}

	newRemoteName := fmt.Sprintf("pulls/%v", pr.Head.Repository.Owner.UserName)

	// verify related remote is in local repo, otherwise add it
	localRemote, err := localRepo.GetOrCreateRemote(remoteURL, newRemoteName)
	if err != nil {
		return err
	}
	localRemoteName := localRemote.Config().Name

	// get auth & fetch remote
	fmt.Printf("Fetching PR %v (head %s:%s) from remote '%s'\n", index, remoteURL, pr.Head.Ref, localRemoteName)
	url, err := local_git.ParseURL(remoteURL)
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
	err = localRepo.TeaCreateBranch(localBranchName, pr.Head.Ref, localRemoteName)
	if err == git.ErrBranchExists {
		fmt.Println("There may be changes since you last checked out, run `git pull` to get them.")
	} else if err != nil {
		return err
	}

	return localRepo.TeaCheckout(localBranchName)
}
