// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repos

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/print"

	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// printFieldsFlag provides a selection of fields to print
var printFieldsFlag = cli.StringFlag{
	Name:    "fields",
	Aliases: []string{"f"},
	Usage: fmt.Sprintf(`Comma-separated list of fields to print. Available values:
		%s
	 `, strings.Join(print.RepoFields, ",")),
	Value: "owner,name,type,ssh",
}

func getFields(ctx *cli.Context) []string {
	return strings.Split(ctx.String("fields"), ",")
}

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
