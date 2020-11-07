// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"
	"log"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
)

// CmdPullsCheckout is a command to locally checkout the given PR
var CmdPullsCheckout = cli.Command{
	Name:        "checkout",
	Usage:       "Locally check out the given PR",
	Description: `Locally check out the given PR`,
	Action:      runPullsCheckout,
	ArgsUsage:   "<pull index>",
	Flags:       flags.AllDefaultFlags,
}

func runPullsCheckout(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if ctx.Args().Len() != 1 {
		log.Fatal("Must specify a PR index")
	}
	idx, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	localRepo, err := local_git.RepoForWorkdir()
	if err != nil {
		return err
	}

	localBranchName, remoteBranchName, newRemoteName, remoteURL, err :=
		gitConfigForPR(localRepo, login, owner, repo, idx)
	if err != nil {
		return err
	}

	// verify related remote is in local repo, otherwise add it
	localRemote, err := localRepo.GetOrCreateRemote(remoteURL, newRemoteName)
	if err != nil {
		return err
	}
	localRemoteName := localRemote.Config().Name

	// get auth & fetch remote
	fmt.Printf("Fetching PR %v (head %s:%s) from remote '%s'\n",
		idx, remoteURL, remoteBranchName, localRemoteName)
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
	err = localRepo.TeaCreateBranch(localBranchName, remoteBranchName, localRemoteName)
	if err == git.ErrBranchExists {
		fmt.Println("There may be changes since you last checked out, run `git pull` to get them.")
	} else if err != nil {
		return err
	}

	return localRepo.TeaCheckout(localBranchName)
}

func gitConfigForPR(repo *local_git.TeaRepo, login *config.Login, owner, repoName string, idx int64) (localBranch, remoteBranch, remoteName, remoteURL string, err error) {
	// fetch PR source-repo & -branch from gitea
	pr, _, err := login.Client().GetPullRequest(owner, repoName, idx)
	if err != nil {
		return
	}

	// test if we can pull via SSH, and configure git remote accordingly
	remoteURL = pr.Head.Repository.CloneURL
	keys, _, err := login.Client().ListMyPublicKeys(gitea.ListPublicKeysOptions{})
	if err != nil {
		return
	}
	if len(keys) != 0 {
		remoteURL = pr.Head.Repository.SSHURL
	}

	// try to find a matching existing branch, otherwise return branch in pulls/ namespace
	localBranch = fmt.Sprintf("pulls/%v-%v", idx, pr.Head.Ref)
	if b, _ := repo.TeaFindBranchBySha(pr.Head.Sha, remoteURL); b != nil {
		localBranch = b.Name
	}

	remoteBranch = pr.Head.Ref
	remoteName = fmt.Sprintf("pulls/%v", pr.Head.Repository.Owner.UserName)
	return
}
