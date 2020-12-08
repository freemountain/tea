// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// MilestoneDetails print an milestone formatted to stdout
func MilestoneDetails(milestone *gitea.Milestone) {
	fmt.Printf("%s\n",
		milestone.Title,
	)
	if len(milestone.Description) != 0 {
		fmt.Printf("\n%s\n", milestone.Description)
	}
	if milestone.Deadline != nil && !milestone.Deadline.IsZero() {
		fmt.Printf("\nDeadline: %s\n", FormatTime(*milestone.Deadline))
	}
}

// MilestonesList prints a listing of milestones
func MilestonesList(miles []*gitea.Milestone, output string, state gitea.StateType) {

	headers := []string{
		"Title",
	}
	if state == gitea.StateAll {
		headers = append(headers, "State")
	}
	headers = append(headers,
		"Open/Closed Issues",
		"DueDate",
	)

	var values [][]string

	for _, m := range miles {
		var deadline = ""

		if m.Deadline != nil && !m.Deadline.IsZero() {
			deadline = FormatTime(*m.Deadline)
		}

		item := []string{
			m.Title,
		}
		if state == gitea.StateAll {
			item = append(item, string(m.State))
		}
		item = append(item,
			fmt.Sprintf("%d/%d", m.OpenIssues, m.ClosedIssues),
			deadline,
		)

		values = append(values, item)
	}
	outputList(output, headers, values)
}
