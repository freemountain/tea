// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"code.gitea.io/tea/modules/context"

	"code.gitea.io/sdk/gitea"
	unidiff "gitea.com/noerw/unidiff-comments"
)

var diffReviewHelp = `# This is the current diff of PR #%d on %s.
# To add code comments, just insert a line inside the diff with your comment,
# prefixed with '# '. For example:
#
# - foo: string,
# - bar: string,
# + foo: int,
# # This is a code comment
# + bar: int,

`

// CreatePullReview submits a review for a PR
func CreatePullReview(ctx *context.TeaContext, idx int64, status gitea.ReviewStateType, comment string, codeComments []gitea.CreatePullReviewComment) error {
	c := ctx.Login.Client()

	review, _, err := c.CreatePullReview(ctx.Owner, ctx.Repo, idx, gitea.CreatePullReviewOptions{
		State:    status,
		Body:     comment,
		Comments: codeComments,
	})
	if err != nil {
		return err
	}

	fmt.Println(review.HTMLURL)
	return nil
}

// SavePullDiff fetches the diff of a pull request and stores it as a temporary file.
// The path to the file is returned.
func SavePullDiff(ctx *context.TeaContext, idx int64) (string, error) {
	diff, _, err := ctx.Login.Client().GetPullRequestDiff(ctx.Owner, ctx.Repo, idx)
	if err != nil {
		return "", err
	}
	writer, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("pull-%d-review-*.diff", idx))
	if err != nil {
		return "", err
	}
	defer writer.Close()

	// add a help header before the actual diff
	if _, err = fmt.Fprintf(writer, diffReviewHelp, idx, ctx.RepoSlug); err != nil {
		return "", err
	}

	if _, err = writer.Write(diff); err != nil {
		return "", err
	}
	return writer.Name(), nil
}

// ParseDiffComments reads a diff, extracts comments from it & returns them in a gitea compatible struct
func ParseDiffComments(diffFile string) ([]gitea.CreatePullReviewComment, error) {
	reader, err := os.Open(diffFile)
	if err != nil {
		return nil, fmt.Errorf("couldn't load diff: %s", err)
	}
	defer reader.Close()

	changeset, err := unidiff.ReadChangeset(reader)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse patch: %s", err)
	}

	var comments []gitea.CreatePullReviewComment
	for _, file := range changeset.Diffs {
		for _, c := range file.LineComments {
			comment := gitea.CreatePullReviewComment{
				Body: c.Text,
				Path: c.Anchor.Path,
			}
			comment.Path = strings.TrimPrefix(comment.Path, "a/")
			comment.Path = strings.TrimPrefix(comment.Path, "b/")
			switch c.Anchor.LineType {
			case "ADDED":
				comment.NewLineNum = c.Anchor.Line
			case "REMOVED", "CONTEXT":
				comment.OldLineNum = c.Anchor.Line
			}
			comments = append(comments, comment)
		}
	}

	return comments, nil
}

// OpenFileInEditor opens filename in a text editor, and blocks until the editor terminates.
func OpenFileInEditor(filename string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			fmt.Println("No $VISUAL or $EDITOR env is set, defaulting to vim")
			editor = "vi"
		}
	}

	// Get the full executable path for the editor.
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
