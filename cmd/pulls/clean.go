// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package pulls

import (
	"fmt"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	git_config "github.com/go-git/go-git/v5/config"
	"github.com/urfave/cli/v2"
)

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
	}, flags.AllDefaultFlags...),
}

func runPullsClean(ctx *cli.Context) error {
	login, owner, repo := config.InitCommand(flags.GlobalRepoValue, flags.GlobalLoginValue, flags.GlobalRemoteValue)
	if ctx.Args().Len() != 1 {
		return fmt.Errorf("Must specify a PR index")
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
