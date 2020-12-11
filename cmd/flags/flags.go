// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package flags

import (
	"code.gitea.io/sdk/gitea"
	"github.com/urfave/cli/v2"
)

// create global variables for global Flags to simplify
// access to the options without requiring cli.Context
var (
	// GlobalLoginValue contain value of --login|-l arg
	GlobalLoginValue string
	// GlobalRepoValue contain value of --repo|-r arg
	GlobalRepoValue string
	// GlobalOutputValue contain value of --output|-o arg
	GlobalOutputValue string
	// GlobalRemoteValue contain value of --remote|-R arg
	GlobalRemoteValue string
)

// LoginFlag provides flag to specify tea login profile
var LoginFlag = cli.StringFlag{
	Name:        "login",
	Aliases:     []string{"l"},
	Usage:       "Use a different Gitea login. Optional",
	Destination: &GlobalLoginValue,
}

// RepoFlag provides flag to specify repository
var RepoFlag = cli.StringFlag{
	Name:        "repo",
	Aliases:     []string{"r"},
	Usage:       "Override local repository path or gitea repository slug to interact with. Optional",
	Destination: &GlobalRepoValue,
}

// RemoteFlag provides flag to specify remote repository
var RemoteFlag = cli.StringFlag{
	Name:        "remote",
	Aliases:     []string{"R"},
	Usage:       "Discover Gitea login from remote. Optional",
	Destination: &GlobalRemoteValue,
}

// OutputFlag provides flag to specify output type
var OutputFlag = cli.StringFlag{
	Name:        "output",
	Aliases:     []string{"o"},
	Usage:       "Output format. (csv, simple, table, tsv, yaml)",
	Destination: &GlobalOutputValue,
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

// GetListOptions return ListOptions based on PaginationFlags
func GetListOptions(ctx *cli.Context) gitea.ListOptions {
	page := ctx.Int("page")
	limit := ctx.Int("limit")
	if limit != 0 && page == 0 {
		page = 1
	}
	return gitea.ListOptions{
		Page:     page,
		PageSize: limit,
	}
}
