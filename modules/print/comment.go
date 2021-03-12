// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
)

// Comments renders a list of comments to stdout
func Comments(comments []*gitea.Comment) {
	var baseURL string
	if len(comments) != 0 {
		baseURL = comments[0].HTMLURL
	}

	var out = make([]string, len(comments))
	for i, c := range comments {
		out[i] = formatComment(c)
	}

	outputMarkdown(fmt.Sprintf(
		// this will become a heading by means of the first --- from a comment
		"Comments\n%s",
		strings.Join(out, "\n"),
	), baseURL)
}

// Comment renders a comment to stdout
func Comment(c *gitea.Comment) {
	outputMarkdown(formatComment(c), c.HTMLURL)
}

func formatComment(c *gitea.Comment) string {
	edited := ""
	if c.Updated.After(c.Created) {
		edited = fmt.Sprintf(" *(edited on %s)*", FormatTime(c.Updated))
	}
	return fmt.Sprintf(
		"---\n\n**@%s** wrote on %s%s:\n\n%s\n",
		c.Poster.UserName,
		FormatTime(c.Created),
		edited,
		c.Body,
	)
}
