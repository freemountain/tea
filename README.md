# <img alt='' src='https://gitea.com/repo-avatars/550-80a3a8c2ab0e2c2d69f296b7f8582485' height="40"/> *T E A*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![Release](https://raster.shields.io/badge/dynamic/json.svg?label=release&url=https://gitea.com/api/v1/repos/gitea/tea/releases&query=$[0].tag_name)](https://gitea.com/gitea/tea/releases) [![Build Status](https://drone.gitea.com/api/badges/gitea/tea/status.svg)](https://drone.gitea.com/gitea/tea) [![Join the chat at https://img.shields.io/discord/322538954119184384.svg](https://img.shields.io/discord/322538954119184384.svg)](https://discord.gg/Gitea) [![Go Report Card](https://goreportcard.com/badge/code.gitea.io/tea)](https://goreportcard.com/report/code.gitea.io/tea) [![GoDoc](https://godoc.org/code.gitea.io/tea?status.svg)](https://godoc.org/code.gitea.io/tea)

### The official CLI for Gitea

![demo gif](./demo.gif)

```
   tea - command line tool to interact with Gitea
   version 0.7.0-preview

 USAGE
   tea command [subcommand] [command options] [arguments...]

 DESCRIPTION
   tea is a productivity helper for Gitea.  It can be used to manage most entities on one
   or multiple Gitea instances and provides local helpers like 'tea pull checkout'.
   tea makes use of context provided by the repository in $PWD if available, but is still
   usable independently of $PWD. Configuration is persisted in $XDG_CONFIG_HOME/tea.

 COMMANDS
   help, h  Shows a list of commands or help for one command
   ENTITIES:
     issues, issue, i                  List, create and update issues
     pulls, pull, pr                   Manage and checkout pull requests
     labels, label                     Manage issue labels
     milestones, milestone, ms         List and create milestones
     releases, release, r              Manage releases
     times, time, t                    Operate on tracked times of a repository's issues & pulls
     organizations, organization, org  List, create, delete organizations
     repos, repo                       Show repository details
   HELPERS:
     open, o                         Open something of the repository in web browser
     notifications, notification, n  Show notifications
   SETUP:
     logins, login                  Log in to a Gitea server
     logout                         Log out from a Gitea server
     shellcompletion, autocomplete  Install shell completion for tea

 OPTIONS
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

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
```

- [Compare features with other git forge CLIs](./FEATURE-COMPARISON.md)
- tea uses [code.gitea.io/sdk](https://code.gitea.io/sdk) and interacts with the Gitea API.

## Installation

There are different ways to get `tea`:

1. Install via your system package manager:
    - macOS via `brew` (gitea-maintained):
      ```sh
      brew tap gitea/tap https://gitea.com/gitea/homebrew-gitea
      brew install tea
      ```
    - arch linux ([gitea-tea](https://aur.archlinux.org/packages/gitea-tea), thirdparty)
    - alpine linux ([tea](https://pkgs.alpinelinux.org/packages?name=tea&branch=edge), thirdparty)

2. Use the prebuilt binaries from [dl.gitea.io](https://dl.gitea.io/tea/)

3. Install from source (go 1.13 or newer is required):
    ```sh
    go get code.gitea.io/tea
    go install code.gitea.io/tea
    ```

4. Docker (thirdparty): [tgerczei/tea](https://hub.docker.com/r/tgerczei/tea)

## Compilation

Make sure you have installed a current go version.
To compile the sources yourself run the following:

```sh
git clone https://gitea.com/gitea/tea.git
cd tea
make
```

## Contributing

Fork -> Patch -> Push -> Pull Request

- `make test` run testsuite
- `make vet`  run checks (check the order of imports; preventing failure on CI pipeline beforehand)
- `make vendor` when adding new dependencies
- ... (for other development tasks, check the `Makefile`)

**Please** read the [CONTRIBUTING](CONTRIBUTING.md) documentation, it will tell you about internal structures and concepts.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the
full license text.
