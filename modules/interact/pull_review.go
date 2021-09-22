// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"os"

	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
)

var reviewStates = map[string]gitea.ReviewStateType{
	"approve":         gitea.ReviewStateApproved,
	"comment":         gitea.ReviewStateComment,
	"request changes": gitea.ReviewStateRequestChanges,
}
var reviewStateOptions = []string{"comment", "request changes", "approve"}

// ReviewPull interactively reviews a PR
func ReviewPull(ctx *context.TeaContext, idx int64) error {
	var state gitea.ReviewStateType
	var comment string
	var codeComments []gitea.CreatePullReviewComment
	var err error

	// codeComments
	var reviewDiff bool
	promptDiff := &survey.Confirm{Message: "Review / comment the diff?", Default: true}
	if err = survey.AskOne(promptDiff, &reviewDiff); err != nil {
		return err
	}
	if reviewDiff {
		if codeComments, err = DoDiffReview(ctx, idx); err != nil {
			fmt.Printf("Error during diff review: %s\n", err)
		}
		fmt.Printf("Found %d code comments in your review\n", len(codeComments))
	}

	// state
	var stateString string
	promptState := &survey.Select{Message: "Your assessment:", Options: reviewStateOptions, VimMode: true}
	if err = survey.AskOne(promptState, &stateString); err != nil {
		return err
	}
	state = reviewStates[stateString]

	// comment
	var promptOpts survey.AskOpt
	if (state == gitea.ReviewStateComment && len(codeComments) == 0) || state == gitea.ReviewStateRequestChanges {
		promptOpts = survey.WithValidator(survey.Required)
	}
	err = survey.AskOne(&survey.Multiline{Message: "Concluding comment:"}, &comment, promptOpts)
	if err != nil {
		return err
	}

	return task.CreatePullReview(ctx, idx, state, comment, codeComments)
}

// DoDiffReview (1) fetches & saves diff in tempfile, (2) starts $VISUAL or $EDITOR to comment on diff,
// (3) parses resulting file into code comments.
func DoDiffReview(ctx *context.TeaContext, idx int64) ([]gitea.CreatePullReviewComment, error) {
	tmpFile, err := task.SavePullDiff(ctx, idx)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile)

	if err = task.OpenFileInEditor(tmpFile); err != nil {
		return nil, err
	}

	return task.ParseDiffComments(tmpFile)
}
