// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package flags

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/context"
	"code.gitea.io/tea/modules/task"
	"code.gitea.io/tea/modules/utils"

	"github.com/araddon/dateparse"
	"github.com/urfave/cli/v2"
)

// LoginFlag provides flag to specify tea login profile
var LoginFlag = cli.StringFlag{
	Name:    "login",
	Aliases: []string{"l"},
	Usage:   "Use a different Gitea Login. Optional",
}

// RepoFlag provides flag to specify repository
var RepoFlag = cli.StringFlag{
	Name:    "repo",
	Aliases: []string{"r"},
	Usage:   "Override local repository path or gitea repository slug to interact with. Optional",
}

// RemoteFlag provides flag to specify remote repository
var RemoteFlag = cli.StringFlag{
	Name:    "remote",
	Aliases: []string{"R"},
	Usage:   "Discover Gitea login from remote. Optional",
}

// OutputFlag provides flag to specify output type
var OutputFlag = cli.StringFlag{
	Name:    "output",
	Aliases: []string{"o"},
	Usage:   "Output format. (csv, simple, table, tsv, yaml)",
}

// StateFlag provides flag to specify issue/pr state, defaulting to "open"
var StateFlag = cli.StringFlag{
	Name:        "state",
	Usage:       "Filter by state (all|open|closed)",
	DefaultText: "open",
}

// PaginationPageFlag provides flag for pagination options
var PaginationPageFlag = cli.StringFlag{
	Name:    "page",
	Aliases: []string{"p"},
	Usage:   "specify page, default is 1",
}

// PaginationLimitFlag provides flag for pagination options
var PaginationLimitFlag = cli.StringFlag{
	Name:    "limit",
	Aliases: []string{"lm"},
	Usage:   "specify limit of items per page",
}

// LoginOutputFlags defines login and output flags that should
// added to all subcommands and appended to the flags of the
// subcommand to work around issue and provide --login and --output:
// https://github.com/urfave/cli/issues/585
var LoginOutputFlags = []cli.Flag{
	&LoginFlag,
	&OutputFlag,
}

// LoginRepoFlags defines login and repo flags that should
// be used for all subcommands and appended to the flags of
// the subcommand to work around issue and provide --login and --repo:
// https://github.com/urfave/cli/issues/585
var LoginRepoFlags = []cli.Flag{
	&LoginFlag,
	&RepoFlag,
	&RemoteFlag,
}

// AllDefaultFlags defines flags that should be available
// for all subcommands working with dedicated repositories
// to work around issue and provide --login, --repo and --output:
// https://github.com/urfave/cli/issues/585
var AllDefaultFlags = append([]cli.Flag{
	&RepoFlag,
	&RemoteFlag,
}, LoginOutputFlags...)

// IssuePRFlags defines flags that should be available on issue & pr listing flags.
var IssuePRFlags = append([]cli.Flag{
	&StateFlag,
	&PaginationPageFlag,
	&PaginationLimitFlag,
}, AllDefaultFlags...)

// IssuePREditFlags defines flags for properties of issues and PRs
var IssuePREditFlags = append([]cli.Flag{
	&cli.StringFlag{
		Name:    "title",
		Aliases: []string{"t"},
	},
	&cli.StringFlag{
		Name:    "description",
		Aliases: []string{"d"},
	},
	&cli.StringFlag{
		Name:    "assignees",
		Aliases: []string{"a"},
		Usage:   "Comma-separated list of usernames to assign",
	},
	&cli.StringFlag{
		Name:    "labels",
		Aliases: []string{"L"},
		Usage:   "Comma-separated list of labels to assign",
	},
	&cli.StringFlag{
		Name:    "deadline",
		Aliases: []string{"D"},
		Usage:   "Deadline timestamp to assign",
	},
	&cli.StringFlag{
		Name:    "milestone",
		Aliases: []string{"m"},
		Usage:   "Milestone to assign",
	},
}, LoginRepoFlags...)

// GetIssuePREditFlags parses all IssuePREditFlags
func GetIssuePREditFlags(ctx *context.TeaContext) (*gitea.CreateIssueOption, error) {
	opts := gitea.CreateIssueOption{
		Title:     ctx.String("title"),
		Body:      ctx.String("description"),
		Assignees: strings.Split(ctx.String("assignees"), ","),
	}
	var err error

	date := ctx.String("deadline")
	if date != "" {
		t, err := dateparse.ParseAny(date)
		if err != nil {
			return nil, err
		}
		opts.Deadline = &t
	}

	client := ctx.Login.Client()

	labelNames := strings.Split(ctx.String("labels"), ",")
	if len(labelNames) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		if opts.Labels, err = task.ResolveLabelNames(client, ctx.Owner, ctx.Repo, labelNames); err != nil {
			return nil, err
		}
	}

	if milestoneName := ctx.String("milestone"); len(milestoneName) != 0 {
		if client == nil {
			client = ctx.Login.Client()
		}
		ms, _, err := client.GetMilestoneByName(ctx.Owner, ctx.Repo, milestoneName)
		if err != nil {
			return nil, fmt.Errorf("Milestone '%s' not found", milestoneName)
		}
		opts.Milestone = ms.ID
	}

	return &opts, nil
}

// FieldsFlag generates a flag selecting printable fields.
// To retrieve the value, use GetFields()
func FieldsFlag(availableFields, defaultFields []string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "fields",
		Aliases: []string{"f"},
		Usage: fmt.Sprintf(`Comma-separated list of fields to print. Available values:
			%s
		`, strings.Join(availableFields, ",")),
		Value: strings.Join(defaultFields, ","),
	}
}

// GetFields parses the values provided in a fields flag, and
// optionally validates against valid values.
func GetFields(ctx *cli.Context, validFields []string) ([]string, error) {
	selection := strings.Split(ctx.String("fields"), ",")
	if validFields != nil {
		for _, field := range selection {
			if !utils.Contains(validFields, field) {
				return nil, fmt.Errorf("Invalid field '%s'", field)
			}
		}
	}
	return selection, nil
}
