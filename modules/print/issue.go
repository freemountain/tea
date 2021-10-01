// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/enescakir/emoji"
)

// IssueDetails print an issue rendered to stdout
func IssueDetails(issue *gitea.Issue, reactions []*gitea.Reaction) {
	out := fmt.Sprintf(
		"# #%d %s (%s)\n@%s created %s\n\n%s\n",
		issue.Index,
		issue.Title,
		issue.State,
		issue.Poster.UserName,
		FormatTime(issue.Created),
		issue.Body,
	)

	if len(reactions) > 0 {
		out += fmt.Sprintf("\n---\n\n%s\n", formatReactions(reactions))
	}

	outputMarkdown(out, issue.HTMLURL)
}

func formatReactions(reactions []*gitea.Reaction) string {
	reactionCounts := make(map[string]uint16)
	for _, r := range reactions {
		reactionCounts[r.Reaction] += 1
	}

	reactionStrings := make([]string, 0, len(reactionCounts))
	for reaction, count := range reactionCounts {
		reactionStrings = append(reactionStrings, fmt.Sprintf("%dx :%s:", count, reaction))
	}

	return emoji.Parse(strings.Join(reactionStrings, "  |  "))
}

// IssuesPullsList prints a listing of issues & pulls
func IssuesPullsList(issues []*gitea.Issue, output string, fields []string) {
	printIssues(issues, output, fields)
}

// IssueFields are all available fields to print with IssuesList()
var IssueFields = []string{
	"index",
	"state",
	"kind",
	"author",
	"author-id",
	"url",

	"title",
	"body",

	"created",
	"updated",
	"deadline",

	"assignees",
	"milestone",
	"labels",
	"comments",
}

func printIssues(issues []*gitea.Issue, output string, fields []string) {
	labelMap := map[int64]string{}
	var printables = make([]printable, len(issues))
	machineReadable := isMachineReadable(output)

	for i, x := range issues {
		// pre-serialize labels for performance
		for _, label := range x.Labels {
			if _, ok := labelMap[label.ID]; !ok {
				labelMap[label.ID] = formatLabel(label, !machineReadable, "")
			}
		}
		// store items with printable interface
		printables[i] = &printableIssue{x, &labelMap}
	}

	t := tableFromItems(fields, printables, machineReadable)
	t.print(output)
}

type printableIssue struct {
	*gitea.Issue
	formattedLabels *map[int64]string
}

func (x printableIssue) FormatField(field string, machineReadable bool) string {
	switch field {
	case "index":
		return fmt.Sprintf("%d", x.Index)
	case "state":
		return string(x.State)
	case "kind":
		if x.PullRequest != nil {
			return "Pull"
		}
		return "Issue"
	case "author":
		return formatUserName(x.Poster)
	case "author-id":
		return x.Poster.UserName
	case "url":
		return x.HTMLURL
	case "title":
		return x.Title
	case "body":
		return x.Body
	case "created":
		return FormatTime(x.Created)
	case "updated":
		return FormatTime(x.Updated)
	case "deadline":
		if x.Deadline == nil {
			return ""
		}
		return FormatTime(*x.Deadline)
	case "milestone":
		if x.Milestone != nil {
			return x.Milestone.Title
		}
		return ""
	case "labels":
		var labels = make([]string, len(x.Labels))
		for i, l := range x.Labels {
			labels[i] = (*x.formattedLabels)[l.ID]
		}
		return strings.Join(labels, " ")
	case "assignees":
		var assignees = make([]string, len(x.Assignees))
		for i, a := range x.Assignees {
			assignees[i] = formatUserName(a)
		}
		return strings.Join(assignees, " ")
	case "comments":
		return fmt.Sprintf("%d", x.Comments)
	}
	return ""
}
