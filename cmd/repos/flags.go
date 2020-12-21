// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

var typeFilterFlag = cli.StringFlag{
	Name:     "type",
	Aliases:  []string{"T"},
	Required: false,
	Usage:    "Filter by type: fork, mirror, source",
}

func getTypeFilter(ctx *cli.Context) (filter gitea.RepoType, err error) {
	t := ctx.String("type")
	filter = gitea.RepoTypeNone
	switch t {
	case "":
		filter = gitea.RepoTypeNone
	case "fork":
		filter = gitea.RepoTypeFork
	case "mirror":
		filter = gitea.RepoTypeMirror
	case "source":
		filter = gitea.RepoTypeSource
	default:
		err = fmt.Errorf("invalid repo type '%s'. valid: fork, mirror, source", t)
	}
	return
}
