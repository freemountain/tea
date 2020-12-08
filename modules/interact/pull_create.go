// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreatePull interactively creates a PR
func CreatePull(login *config.Login, owner, repo string) error {
	var base, head, title, description string

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// base
	baseBranch, err := task.GetDefaultPRBase(login, owner, repo)
	if err != nil {
		return err
	}
	promptI := &survey.Input{Message: "Target branch [" + baseBranch + "]:"}
	if err := survey.AskOne(promptI, &base); err != nil {
		return err
	}
	if len(base) == 0 {
		base = baseBranch
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
	var headOwnerInput, headBranchInput string
	promptI = &survey.Input{Message: "Source repo owner [" + headOwner + "]:"}
	if err := survey.AskOne(promptI, &headOwnerInput); err != nil {
		return err
	}
	if len(headOwnerInput) != 0 {
		headOwner = headOwnerInput
	}
	promptI = &survey.Input{Message: "Source branch [" + headBranch + "]:"}
	if err := survey.AskOne(promptI, &headBranchInput, promptOpts); err != nil {
		return err
	}
	if len(headBranchInput) != 0 {
		headBranch = headBranchInput
	}

	head = task.GetHeadSpec(headOwner, headBranch, owner)

	// title
	title = task.GetDefaultPRTitle(head)
	promptOpts = survey.WithValidator(survey.Required)
	if len(title) != 0 {
		promptOpts = nil
	}
	promptI = &survey.Input{Message: "PR title [" + title + "]:"}
	if err := survey.AskOne(promptI, &title, promptOpts); err != nil {
		return err
	}

	// description
	promptM := &survey.Multiline{Message: "PR description:"}
	if err := survey.AskOne(promptM, &description); err != nil {
		return err
	}

	return task.CreatePull(
		login,
		owner,
		repo,
		base,
		head,
		title,
		description)
}

func promptRepoSlug(defaultOwner, defaultRepo string) (owner, repo string, err error) {
	prompt := "Target repo:"
	required := true
	if len(defaultOwner) != 0 && len(defaultRepo) != 0 {
		prompt = fmt.Sprintf("Target repo [%s/%s]:", defaultOwner, defaultRepo)
		required = false
	}
	var repoSlug string

	owner = defaultOwner
	repo = defaultRepo

	err = survey.AskOne(
		&survey.Input{Message: prompt},
		&repoSlug,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if !required && len(str) == 0 {
					return nil
				}
				split := strings.Split(str, "/")
				if len(split) != 2 || len(split[0]) == 0 || len(split[1]) == 0 {
					return fmt.Errorf("must follow the <owner>/<repo> syntax")
				}
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)

	if err == nil && len(repoSlug) != 0 {
		repoSlugSplit := strings.Split(repoSlug, "/")
		owner = repoSlugSplit[0]
		repo = repoSlugSplit[1]
	}
	return
}
