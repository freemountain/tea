// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/interact"
	"code.gitea.io/tea/modules/print"
	"code.gitea.io/tea/modules/utils"

	"code.gitea.io/sdk/gitea"
	"github.com/AlecAivazis/survey/v2"
	"github.com/urfave/cli/v2"
)

// CmdAddComment is the main command to operate with notifications
var CmdAddComment = cli.Command{
	Name:        "comment",
	Aliases:     []string{"c"},
	Category:    catEntities,
	Usage:       "Add a comment to an issue / pr",
	Description: "Add a comment to an issue / pr",
	ArgsUsage:   "<issue / pr index> [<comment body>]",
	Action:      runAddComment,
	Flags:       flags.AllDefaultFlags,
}

func runAddComment(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	args := ctx.Args()
	if args.Len() == 0 {
		return fmt.Errorf("Please specify issue / pr index")
	}

	idx, err := utils.ArgToIndex(ctx.Args().First())
	if err != nil {
		return err
	}

	body := strings.Join(ctx.Args().Tail(), " ")
	if interact.IsStdinPiped() {
		// custom solution until https://github.com/AlecAivazis/survey/issues/328 is fixed
		if bodyStdin, err := ioutil.ReadAll(ctx.App.Reader); err != nil {
			return err
		} else if len(bodyStdin) != 0 {
			body = strings.Join([]string{body, string(bodyStdin)}, "\n\n")
		}
	} else if len(body) == 0 {
		if err = survey.AskOne(interact.NewMultiline(interact.Multiline{
			Message:   "Comment:",
			Syntax:    "md",
			UseEditor: config.GetPreferences().Editor,
		}), &body); err != nil {
			return err
		}
	}

	if len(body) == 0 {
		return fmt.Errorf("No comment body provided")
	}

	client := ctx.Login.Client()
	comment, _, err := client.CreateIssueComment(ctx.Owner, ctx.Repo, idx, gitea.CreateIssueCommentOption{
		Body: body,
	})
	if err != nil {
		return err
	}

	print.Comment(comment)

	return nil
}
