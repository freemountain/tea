// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"path"
	"strings"

	"code.gitea.io/tea/cmd/flags"
	"code.gitea.io/tea/modules/context"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli/v2"
)

// CmdOpen represents a sub command of issues to open issue on the web browser
var CmdOpen = cli.Command{
	Name:        "open",
	Aliases:     []string{"o"},
	Category:    catHelpers,
	Usage:       "Open something of the repository in web browser",
	Description: `Open something of the repository in web browser`,
	Action:      runOpen,
	Flags:       append([]cli.Flag{}, flags.LoginRepoFlags...),
}

func runOpen(cmd *cli.Context) error {
	ctx := context.InitCommand(cmd)
	ctx.Ensure(context.CtxRequirement{RemoteRepo: true})

	var suffix string
	number := ctx.Args().Get(0)
	switch {
	case strings.EqualFold(number, "issues"):
		suffix = "issues"
	case strings.EqualFold(number, "pulls"):
		suffix = "pulls"
	case strings.EqualFold(number, "releases"):
		suffix = "releases"
	case strings.EqualFold(number, "commits"):
		repo, err := local_git.RepoForWorkdir()
		if err != nil {
			return err
		}
		b, err := repo.Head()
		if err != nil {
			return err
		}
		name := b.Name()
		switch {
		case name.IsBranch():
			suffix = "commits/branch/" + name.Short()
		case name.IsTag():
			suffix = "commits/tag/" + name.Short()
		}
	case strings.EqualFold(number, "branches"):
		suffix = "branches"
	case strings.EqualFold(number, "wiki"):
		suffix = "wiki"
	case strings.EqualFold(number, "activity"):
		suffix = "activity"
	case strings.EqualFold(number, "settings"):
		suffix = "settings"
	case strings.EqualFold(number, "labels"):
		suffix = "labels"
	case strings.EqualFold(number, "milestones"):
		suffix = "milestones"
	case number != "":
		suffix = "issues/" + number
	default:
		suffix = number
	}

	u := path.Join(ctx.Login.URL, ctx.Owner, ctx.Repo, suffix)
	return open.Run(u)
}
