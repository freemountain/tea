// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/sdk/gitea"
)

// IssueDetails print an issue rendered to stdout
func IssueDetails(issue *gitea.Issue) {
	outputMarkdown(fmt.Sprintf(
		"# #%d %s (%s)\n@%s created %s\n\n%s\n",
		issue.Index,
		issue.Title,
		issue.State,
		issue.Poster.UserName,
		FormatTime(issue.Created),
		issue.Body,
	))
}

// IssuesList prints a listing of issues
func IssuesList(issues []*gitea.Issue, output string) {
	t := tableWithHeader(
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	)

	for _, issue := range issues {
		author := issue.Poster.FullName
		if len(author) == 0 {
			author = issue.Poster.UserName
		}
		mile := ""
		if issue.Milestone != nil {
			mile = issue.Milestone.Title
		}
		t.addRow(
			strconv.FormatInt(issue.Index, 10),
			issue.Title,
			string(issue.State),
			author,
			mile,
			FormatTime(issue.Updated),
		)
	}
	t.print(output)
}

// IssuesPullsList prints a listing of issues & pulls
// TODO combine with IssuesList
func IssuesPullsList(issues []*gitea.Issue, output string) {
	t := tableWithHeader(
		"Index",
		"State",
		"Kind",
		"Author",
		"Updated",
		"Title",
	)

	for _, issue := range issues {
		name := issue.Poster.FullName
		if len(name) == 0 {
			name = issue.Poster.UserName
		}
		kind := "Issue"
		if issue.PullRequest != nil {
			kind = "Pull"
		}
		t.addRow(
			strconv.FormatInt(issue.Index, 10),
			string(issue.State),
			kind,
			name,
			FormatTime(issue.Updated),
			issue.Title,
		)
	}

	t.print(output)
}
