// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	var title, description string

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Issue title:"}
	if err := survey.AskOne(promptI, &title, promptOpts); err != nil {
		return err
	}

	// description
	promptM := &survey.Multiline{Message: "Issue description:"}
	if err := survey.AskOne(promptM, &description); err != nil {
		return err
	}

	return task.CreateIssue(
		login,
		owner,
		repo,
		title,
		description)
}
