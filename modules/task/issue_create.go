// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/print"
)

// CreateIssue creates an issue in the given repo and prints the result
func CreateIssue(login *config.Login, repoOwner, repoName, title, description string) error {

	// title is required
	if len(title) == 0 {
		return fmt.Errorf("Title is required")
	}

	issue, _, err := login.Client().CreateIssue(repoOwner, repoName, gitea.CreateIssueOption{
		Title: title,
		Body:  description,
		// TODO:
		//Assignee  string   `json:"assignee"`
		//Assignees []string `json:"assignees"`
		//Deadline *time.Time `json:"due_date"`
		//Milestone int64 `json:"milestone"`
		//Labels []int64 `json:"labels"`
		//Closed bool    `json:"closed"`
	})

	if err != nil {
		return fmt.Errorf("could not create issue: %s", err)
	}

	print.IssueDetails(issue)

	fmt.Println(issue.HTMLURL)

	return nil
}
