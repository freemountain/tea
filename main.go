// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tea is command line tool for Gitea.
package main // import "code.gitea.io/tea"

import (
	"fmt"
	"os"
	"strings"

	"code.gitea.io/tea/cmd"

	"github.com/urfave/cli/v2"
)

// Version holds the current tea version
var Version = "development"

// Tags holds the build tags used
var Tags = ""

func main() {
	// make parsing tea --version easier, by printing /just/ the version string
	cli.VersionPrinter = func(c *cli.Context) { fmt.Fprintln(c.App.Writer, c.App.Version) }

	app := cli.NewApp()
	app.Name = "tea"
	app.Usage = "command line tool to interact with Gitea"
	app.Description = appDescription
	app.CustomAppHelpTemplate = helpTemplate
	app.Version = Version + formatBuiltWith(Tags)
	app.Commands = []*cli.Command{
		&cmd.CmdLogin,
		&cmd.CmdLogout,
		&cmd.CmdAutocomplete,

		&cmd.CmdIssues,
		&cmd.CmdPulls,
		&cmd.CmdLabels,
		&cmd.CmdMilestones,
		&cmd.CmdReleases,
		&cmd.CmdTrackedTimes,
		&cmd.CmdOrgs,
		&cmd.CmdRepos,
		&cmd.CmdAddComment,

		&cmd.CmdOpen,
		&cmd.CmdNotifications,
	}
	app.EnableBashCompletion = true
	err := app.Run(os.Args)
	if err != nil {
		// app.Run already exits for errors implementing ErrorCoder,
		// so we only handle generic errors with code 1 here.
		fmt.Fprintf(app.ErrWriter, "Error: %v\n", err)
		os.Exit(1)
	}
}

func formatBuiltWith(Tags string) string {
	if len(Tags) == 0 {
		return ""
	}

	return " built with: " + strings.Replace(Tags, " ", ", ", -1)
}

var appDescription = `tea is a productivity helper for Gitea.  It can be used to manage most entities on one
or multiple Gitea instances and provides local helpers like 'tea pull checkout'.
tea makes use of context provided by the repository in $PWD if available, but is still
usable independently of $PWD. Configuration is persisted in $XDG_CONFIG_HOME/tea.
`

var helpTemplate = bold(`
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}`) + `
   {{if .Version}}{{if not .HideVersion}}version {{.Version}}{{end}}{{end}}

 USAGE
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .Commands}} command [subcommand] [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Description}}

 DESCRIPTION
   {{.Description | nindent 3 | trim}}{{end}}{{if .VisibleCommands}}

 COMMANDS{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{else}}{{range .VisibleCommands}}
   {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}

 OPTIONS
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}

 EXAMPLES
   tea login add                       # add a login once to get started

   tea pulls                           # list open pulls for the repo in $PWD
   tea pulls --repo $HOME/foo          # list open pulls for the repo in $HOME/foo
   tea pulls --remote upstream         # list open pulls for the repo pointed at by
                                       # your local "upstream" git remote
   # list open pulls for any gitea repo at the given login instance
   tea pulls --repo gitea/tea --login gitea.com

   tea milestone issues 0.7.0          # view open issues for milestone '0.7.0'
   tea issue 189                       # view contents of issue 189
   tea open 189                        # open web ui for issue 189
   tea open milestones                 # open web ui for milestones

   # send gitea desktop notifications every 5 minutes (bash + libnotify)
   while :; do tea notifications --all -o simple | xargs -i notify-send {}; sleep 300; done

 ABOUT
   Written & maintained by The Gitea Authors.
   If you find a bug or want to contribute, we'll welcome you at https://gitea.com/gitea/tea.
   More info about Gitea itself on https://gitea.io.
`

func bold(t string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", t)
}
