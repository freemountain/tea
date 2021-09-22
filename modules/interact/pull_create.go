// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreatePull interactively creates a PR
func CreatePull(ctx *context.TeaContext) (err error) {
	var base, head string

	// owner, repo
	if ctx.Owner, ctx.Repo, err = promptRepoSlug(ctx.Owner, ctx.Repo); err != nil {
		return err
	}

	// base
	if base, err = task.GetDefaultPRBase(ctx.Login, ctx.Owner, ctx.Repo); err != nil {
		return err
	}
	promptI := &survey.Input{Message: "Target branch:", Default: base}
	if err := survey.AskOne(promptI, &base); err != nil {
		return err
	}

	// head
	var headOwner, headBranch string
	promptOpts := survey.WithValidator(survey.Required)

	if ctx.LocalRepo != nil {
		headOwner, headBranch, err = task.GetDefaultPRHead(ctx.LocalRepo)
		if err == nil {
			promptOpts = nil
		}
	}
	promptI = &survey.Input{Message: "Source repo owner:", Default: headOwner}
	if err := survey.AskOne(promptI, &headOwner); err != nil {
		return err
	}
	promptI = &survey.Input{Message: "Source branch:", Default: headBranch}
	if err := survey.AskOne(promptI, &headBranch, promptOpts); err != nil {
		return err
	}

	head = task.GetHeadSpec(headOwner, headBranch, ctx.Owner)

	opts := gitea.CreateIssueOption{Title: task.GetDefaultPRTitle(head)}
	if err = promptIssueProperties(ctx.Login, ctx.Owner, ctx.Repo, &opts); err != nil {
		return err
	}

	return task.CreatePull(
		ctx,
		base,
		head,
		&opts)
}
