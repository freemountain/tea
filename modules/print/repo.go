// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
)

// ReposList prints a listing of the repos
func ReposList(rps []*gitea.Repository) {
	if len(rps) == 0 {
		fmt.Println("No repositories found")
		return
	}

	headers := []string{
		"Name",
		"Type",
		"SSH",
		"Owner",
	}
	var values [][]string

	for _, rp := range rps {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}

		values = append(
			values,
			[]string{
				rp.FullName,
				mode,
				rp.SSHURL,
				rp.Owner.UserName,
			},
		)
	}

	OutputList(flags.GlobalOutputValue, headers, values)
}

// RepoDetails print an repo formatted to stdout
func RepoDetails(repo *gitea.Repository, topics []string) {
	output := repo.FullName
	if repo.Mirror {
		output += " (mirror)"
	}
	if repo.Fork {
		output += " (fork)"
	}
	if repo.Archived {
		output += " (archived)"
	}
	if repo.Empty {
		output += " (empty)"
	}
	output += "\n"
	if len(topics) != 0 {
		output += "Topics: " + strings.Join(topics, ", ") + "\n"
	}
	output += "\n"
	output += repo.Description + "\n\n"
	output += fmt.Sprintf(
		"Open Issues: %d, Stars: %d, Forks: %d, Size: %s\n\n",
		repo.OpenIssues,
		repo.Stars,
		repo.Forks,
		formatSize(int64(repo.Size)),
	)

	fmt.Print(output)
}
