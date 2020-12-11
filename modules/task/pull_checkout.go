// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"

	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/go-git/go-git/v5"
)

// PullCheckout checkout current workdir to the head branch of specified pull request
func PullCheckout(login *config.Login, repoOwner, repoName string, index int64, callback func(string) (string, error)) error {
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

	remoteURL := pr.Head.Repository.CloneURL
	if len(login.SSHKey) != 0 {
		// login.SSHKey is nonempty, if user specified a key manually or we automatically
		// found a matching private key on this machine during login creation.
		// this means, we are very likely to have a working ssh setup.
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

	// checkout local branch
	err = localRepo.TeaCreateBranch(localBranchName, pr.Head.Ref, localRemoteName)
	if err == nil {
		fmt.Printf("Created branch '%s'\n", localBranchName)
	} else if err == git.ErrBranchExists {
		fmt.Println("There may be changes since you last checked out, run `git pull` to get them.")
	} else if err != nil {
		return err
	}

	return localRepo.TeaCheckout(localBranchName)
}
