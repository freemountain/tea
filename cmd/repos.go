// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdRepos represents to login a gitea server.
var CmdRepos = cli.Command{
	Name:        "repos",
	Usage:       "Operate with repositories",
	Description: `Operate with repositories`,
	Action:      runReposList,
	Subcommands: []cli.Command{
		CmdReposList,
	},
	Flags: LoginOutputFlags,
}

// CmdReposList represents a sub command of issues to list issues
var CmdReposList = cli.Command{
	Name:        "ls",
	Usage:       "List available repositories",
	Description: `List available repositories`,
	Action:      runReposList,
	Flags: append([]cli.Flag{
		cli.StringFlag{
			Name:  "mode",
			Usage: "Filter listed repositories based on mode, optional - fork, mirror, source",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "Filter listed repositories based on organization, optional",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Filter listed repositories absed on user, optional",
		},
	}, LoginOutputFlags...),
}

// runReposList list repositories
func runReposList(ctx *cli.Context) error {
	login := initCommandLoginOnly()

	mode := ctx.String("mode")
	org := ctx.String("org")
	user := ctx.String("user")

	var rps []*gitea.Repository
	var err error

	if org != "" {
		rps, err = login.Client().ListOrgRepos(org)
	} else if user != "" {
		rps, err = login.Client().ListUserRepos(user)
	} else {
		rps, err = login.Client().ListMyRepos()
	}
	if err != nil {
		log.Fatal(err)
	}

	var repos []*gitea.Repository
	if mode == "" {
		repos = rps
	} else if mode == "fork" {
		for _, rp := range rps {
			if rp.Fork == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "mirror" {
		for _, rp := range rps {
			if rp.Mirror == true {
				repos = append(repos, rp)
			}
		}
	} else if mode == "source" {
		for _, rp := range rps {
			if rp.Mirror != true && rp.Fork != true {
				repos = append(repos, rp)
			}
		}
	} else {
		log.Fatal("Unknown mode: ", mode, "\nUse one of the following:\n- fork\n- mirror\n- source\n")
		return nil
	}

	if len(rps) == 0 {
		log.Fatal("No repositories found", rps)
		return nil
	}

	headers := []string{
		"Name",
		"Type",
		"SSH",
		"Owner",
	}
	var values [][]string

	for _, rp := range repos {
		var mode = "source"
		if rp.Fork {
			mode = "fork"
		}
		if rp.Mirror {
			mode = "mirror"
		}

		values = append(
			values,
			[]string{
				rp.FullName,
				mode,
				rp.SSHURL,
				rp.Owner.UserName,
			},
		)
	}
	Output(outputValue, headers, values)

	return nil
}
