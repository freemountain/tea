// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"

	"code.gitea.io/sdk/gitea"
)

// PullDetails print an pull rendered to stdout
func PullDetails(pr *gitea.PullRequest, reviews []*gitea.PullReview) {
	base := pr.Base.Name
	head := pr.Head.Name
	if pr.Head.RepoID != pr.Base.RepoID {
		if pr.Head.Repository != nil {
			head = pr.Head.Repository.Owner.UserName + ":" + head
		} else {
			head = "delete:" + head
		}
	}

	out := fmt.Sprintf(
		"# #%d %s (%s)\n@%s created %s\t**%s** <- **%s**\n\n%s\n",
		pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		FormatTime(*pr.Created),
		base,
		head,
		pr.Body,
	)

	if len(reviews) != 0 {
		out += "\n"
		revMap := make(map[string]gitea.ReviewStateType)
		for _, review := range reviews {
			switch review.State {
			case gitea.ReviewStateApproved,
				gitea.ReviewStateRequestChanges,
				gitea.ReviewStateRequestReview:
				revMap[review.Reviewer.UserName] = review.State
			}
		}
		for k, v := range revMap {
			out += fmt.Sprintf("\n  @%s: %s", k, v)
		}
	}

	if pr.State == gitea.StateOpen && pr.Mergeable {
		out += "\nNo Conflicts"
	}

	outputMarkdown(out)
}

// PullsList prints a listing of pulls
func PullsList(prs []*gitea.PullRequest, output string) {
	t := tableWithHeader(
		"Index",
		"Title",
		"State",
		"Author",
		"Milestone",
		"Updated",
	)

	for _, pr := range prs {
		if pr == nil {
			continue
		}
		author := pr.Poster.FullName
		if len(author) == 0 {
			author = pr.Poster.UserName
		}
		mile := ""
		if pr.Milestone != nil {
			mile = pr.Milestone.Title
		}
		t.addRow(
			strconv.FormatInt(pr.Index, 10),
			pr.Title,
			string(pr.State),
			author,
			mile,
			FormatTime(*pr.Updated),
		)
	}

	t.print(output)
}
