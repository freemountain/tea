// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// PullDetails print an pull rendered to stdout
func PullDetails(pr *gitea.PullRequest) {
	OutputMarkdown(fmt.Sprintf(
		"# #%d %s (%s)\n%s created %s\n\n%s\n",
		pr.Index,
		pr.Title,
		pr.State,
		pr.Poster.UserName,
		FormatTime(*pr.Created),
		pr.Body,
	))
}
