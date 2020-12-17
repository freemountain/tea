// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"time"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
	"github.com/araddon/dateparse"
)

// CreateMilestone interactively creates a milestone
func CreateMilestone(login *config.Login, owner, repo string) error {
	var title, description, dueDate string
	var deadline *time.Time

	// owner, repo
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Milestone title:"}
	if err := survey.AskOne(promptI, &title, promptOpts); err != nil {
		return err
	}

	// description
	promptM := &survey.Multiline{Message: "Milestone description:"}
	if err := survey.AskOne(promptM, &description); err != nil {
		return err
	}

	// deadline
	promptI = &survey.Input{Message: "Milestone deadline [no due date]:"}
	err = survey.AskOne(
		promptI,
		&dueDate,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if len(str) == 0 {
					return nil
				}
				t, err := dateparse.ParseAny(str)
				if err != nil {
					return err
				}
				deadline = &t
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)

	if err != nil {
		return err
	}

	return task.CreateMilestone(
		login,
		owner,
		repo,
		title,
		description,
		deadline,
		gitea.StateOpen)
}
