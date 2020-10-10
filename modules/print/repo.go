// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"log"
	"strings"
	"time"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
)

type rp = *gitea.Repository
type fieldFormatter = func(*gitea.Repository) string

var (
	fieldFormatters map[string]fieldFormatter

	// RepoFields are the available fields to print with ReposList()
	RepoFields []string
)

func init() {
	fieldFormatters = map[string]fieldFormatter{
		"description": func(r rp) string { return r.Description },
		"forks":       func(r rp) string { return fmt.Sprintf("%d", r.Forks) },
		"id":          func(r rp) string { return r.FullName },
		"name":        func(r rp) string { return r.Name },
		"owner":       func(r rp) string { return r.Owner.UserName },
		"stars":       func(r rp) string { return fmt.Sprintf("%d", r.Stars) },
		"ssh":         func(r rp) string { return r.SSHURL },
		"updated":     func(r rp) string { return FormatTime(r.Updated) },
		"url":         func(r rp) string { return r.HTMLURL },
		"permission": func(r rp) string {
			if r.Permissions.Admin {
				return "admin"
			} else if r.Permissions.Push {
				return "write"
			}
			return "read"
		},
		"type": func(r rp) string {
			if r.Fork {
				return "fork"
			}
			if r.Mirror {
				return "mirror"
			}
			return "source"
		},
	}

	for f := range fieldFormatters {
		RepoFields = append(RepoFields, f)
	}
}

// ReposList prints a listing of the repos
func ReposList(repos []*gitea.Repository, fields []string) {
	if len(repos) == 0 {
		fmt.Println("No repositories found")
		return
	}

	if len(fields) == 0 {
		fmt.Println("No fields to print")
		return
	}

	formatters := make([]fieldFormatter, len(fields))
	values := make([][]string, len(repos))

	// find field format functions by header name
	for i, f := range fields {
		if formatter, ok := fieldFormatters[strings.ToLower(f)]; ok {
			formatters[i] = formatter
		} else {
			log.Fatalf("invalid field '%s'", f)
		}
	}

	// extract values from each repo and store them in 2D table
	for i, repo := range repos {
		values[i] = make([]string, len(formatters))
		for j, format := range formatters {
			values[i][j] = format(repo)
		}
	}

	OutputList(flags.GlobalOutputValue, fields, values)
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
		fieldFormatters["permission"](repo),
	)

	var tops string
	if len(topics) != 0 {
		tops = fmt.Sprintf("- Topics:\t%s\n", strings.Join(topics, ", "))
	}

	OutputMarkdown(fmt.Sprintf(
		"%s%s\n%s\n%s%s%s%s",
		title,
		desc,
		stats,
		updated,
		urls,
		perm,
		tops,
	))
}
