// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"os"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/print"
	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/crypto/ssh/terminal"
)

// ShowCommentsMaybeInteractive fetches & prints comments, depending on the --comments flag.
// If that flag is unset, and output is not piped, prompts the user first.
func ShowCommentsMaybeInteractive(ctx *context.TeaContext, idx int64, totalComments int) error {
	if ctx.Bool("comments") {
		opts := gitea.ListIssueCommentOptions{ListOptions: ctx.GetListOptions()}
		c := ctx.Login.Client()
		comments, _, err := c.ListIssueComments(ctx.Owner, ctx.Repo, idx, opts)
		if err != nil {
			return err
		}
		print.Comments(comments)
	} else if isInteractive() && !ctx.IsSet("comments") {
		// if we're interactive, but --comments hasn't been explicitly set to false
		if err := ShowCommentsPaginated(ctx, idx, totalComments); err != nil {
			fmt.Printf("error while loading comments: %v\n", err)
		}
	}
	return nil
}

// ShowCommentsPaginated prompts if issue/pr comments should be shown and continues to do so.
func ShowCommentsPaginated(ctx *context.TeaContext, idx int64, totalComments int) error {
	c := ctx.Login.Client()
	opts := gitea.ListIssueCommentOptions{ListOptions: ctx.GetListOptions()}
	prompt := "show comments?"
	commentsLoaded := 0

	// paginated fetch
	// NOTE: as of gitea 1.13, pagination is not provided by this endpoint, but handles
	// this function gracefully anyways.
	for {
		loadComments := false
		confirm := survey.Confirm{Message: prompt, Default: true}
		if err := survey.AskOne(&confirm, &loadComments); err != nil {
			return err
		} else if !loadComments {
			break
		} else {
			if comments, _, err := c.ListIssueComments(ctx.Owner, ctx.Repo, idx, opts); err != nil {
				return err
			} else if len(comments) != 0 {
				print.Comments(comments)
				commentsLoaded += len(comments)
			}
			if commentsLoaded >= totalComments {
				break
			}
			opts.ListOptions.Page++
			prompt = "load more?"
		}
	}
	return nil
}

// IsInteractive checks if the output is piped, but NOT if the session is run interactively..
func isInteractive() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
