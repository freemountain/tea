// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"time"

	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
)

// CreateMilestone creates a milestone in the given repo and prints the result
func CreateMilestone(login *config.Login, repoOwner, repoName, title, description string, deadline *time.Time, state gitea.StateType) error {

	// title is required
	if len(title) == 0 {
		return fmt.Errorf("Title is required")
	}

	mile, _, err := login.Client().CreateMilestone(repoOwner, repoName, gitea.CreateMilestoneOption{
		Title:       title,
		Description: description,
		Deadline:    deadline,
		State:       state,
	})
	if err != nil {
		return err
	}

	print.MilestoneDetails(mile)
	return nil
}
