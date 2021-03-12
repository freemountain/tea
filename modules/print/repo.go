// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
)

// ReposList prints a listing of the repos
func ReposList(repos []*gitea.Repository, output string, fields []string) {
	var printables = make([]printable, len(repos))
	for i, r := range repos {
		printables[i] = &printableRepo{r}
	}
	t := tableFromItems(fields, printables)
	t.print(output)
}

// RepoDetails print an repo formatted to stdout
func RepoDetails(repo *gitea.Repository, topics []string) {
	title := "# " + repo.FullName
	if repo.Mirror {
		title += " (mirror)"
	}
	if repo.Fork {
		title += " (fork)"
	}
	if repo.Archived {
		title += " (archived)"
	}
	if repo.Empty {
		title += " (empty)"
	}
	title += "\n"

	var desc string
	if len(repo.Description) != 0 {
		desc = fmt.Sprintf("*%s*\n\n", repo.Description)
	}

	stats := fmt.Sprintf(
		"Issues: %d, Stars: %d, Forks: %d, Size: %s\n",
		repo.OpenIssues,
		repo.Stars,
		repo.Forks,
		formatSize(int64(repo.Size)),
	)

	// NOTE: for mirrors, this is the time the mirror was last fetched..
	updated := fmt.Sprintf(
		"Updated: %s (%s ago)\n",
		repo.Updated.Format("2006-01-02 15:04"),
		time.Now().Sub(repo.Updated).Truncate(time.Minute),
	)

	urls := fmt.Sprintf(
		"- Browse:\t%s\n- Clone:\t%s\n",
		repo.HTMLURL,
		repo.SSHURL,
	)
	if len(repo.Website) != 0 {
		urls += fmt.Sprintf("- Web:\t%s\n", repo.Website)
	}

	perm := fmt.Sprintf(
		"- Permission:\t%s\n",
		formatPermission(repo.Permissions),
	)

	var tops string
	if len(topics) != 0 {
		tops = fmt.Sprintf("- Topics:\t%s\n", strings.Join(topics, ", "))
	}

	outputMarkdown(fmt.Sprintf(
		"%s%s\n%s\n%s%s%s%s",
		title,
		desc,
		stats,
		updated,
		urls,
		perm,
		tops,
	), repo.HTMLURL)
}

// RepoFields are the available fields to print with ReposList()
var RepoFields = []string{
	"description",
	"forks",
	"id",
	"name",
	"owner",
	"stars",
	"ssh",
	"updated",
	"url",
	"permission",
	"type",
}

type printableRepo struct{ *gitea.Repository }

func (x printableRepo) FormatField(field string) string {
	switch field {
	case "description":
		return x.Description
	case "forks":
		return fmt.Sprintf("%d", x.Forks)
	case "id":
		return x.FullName
	case "name":
		return x.Name
	case "owner":
		return x.Owner.UserName
	case "stars":
		return fmt.Sprintf("%d", x.Stars)
	case "ssh":
		return x.SSHURL
	case "updated":
		return FormatTime(x.Updated)
	case "url":
		return x.HTMLURL
	case "permission":
		return formatPermission(x.Permissions)
	case "type":
		if x.Fork {
			return "fork"
		}
		if x.Mirror {
			return "mirror"
		}
		return "source"
	}
	return ""
}
