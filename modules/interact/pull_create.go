// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreatePull interactively creates a PR
func CreatePull(login *config.Login, owner, repo string) error {
	var base, head string

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// base
	base, err = task.GetDefaultPRBase(login, owner, repo)
	if err != nil {
		return err
	}
	promptI := &survey.Input{Message: "Target branch:", Default: base}
	if err := survey.AskOne(promptI, &base); err != nil {
		return err
	}

	// head
	localRepo, err := git.RepoForWorkdir()
	if err != nil {
		return err
	}
	promptOpts := survey.WithValidator(survey.Required)
	headOwner, headBranch, err := task.GetDefaultPRHead(localRepo)
	if err == nil {
		promptOpts = nil
	}
	promptI = &survey.Input{Message: "Source repo owner:", Default: headOwner}
	if err := survey.AskOne(promptI, &headOwner); err != nil {
		return err
	}
	promptI = &survey.Input{Message: "Source branch:", Default: headBranch}
	if err := survey.AskOne(promptI, &headBranch, promptOpts); err != nil {
		return err
	}

	head = task.GetHeadSpec(headOwner, headBranch, owner)

	opts := gitea.CreateIssueOption{Title: task.GetDefaultPRTitle(head)}
	if err = promptIssueProperties(login, owner, repo, &opts); err != nil {
		return err
	}

	return task.CreatePull(
		login,
		owner,
		repo,
		base,
		head,
		&opts)
}
