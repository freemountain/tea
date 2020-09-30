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

	// fetch PR source-repo & -branch from gitea
	idx, err := utils.ArgToIndex(ctx.Args().First())
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
