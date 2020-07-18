# Changelog

## [v0.4.0](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=1264) - 2020-07-18

* FEATURES
  * Add notifications subcomand (#148)
  * Add subcomand 'pulls create' (#144)
* BUGFIXES
  * Fix Login Detection By Repo Param (#151)
  * Fix Login List Output (#150)
  * Fix --ssh-key Option (#135)
* ENHANCEMENTS
  * Subcomand Login Show List By Default (#152)
* BUILD
  * Migrate src-d/go-git to go-git/go-git (#128)
  * Migrate gitea-sdk to v0.12.0 (#133)
  * Migrate yaml lib (#130)
  * Add gitea-vet (#121)

## [v0.3.1](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=1265) - 2020-06-15

* BUGFIXES
  * --ssh-key should be string not bool (#135) (#137)
  * modules/git: fix dropped error (#127)
  * Issues details: add missing newline (#126)

## [v0.3.0](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=1227) - 2020-04-22

* FEATURES
  * Add `tea pulls [checkout | clean]` commands (#93 #97 #107) (#105)
  * Add `tea open` (#101)
  * Add `tea issues [open|close]` commands (#99)
* ENHANCEMENTS
  * Ignore PRs for `tea issues` (#111)
  * Add --state flag filter to issue & PR lists (#100)

## [v0.2.0](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=538) - 2020-03-06
* FEATURES
  * Add `tea times` command (#54)
* ENHANCEMENTS
  * Upgrade urfave/cli to v2 version (#85)
  * Add --remote flag to add/create subcommands (#77)
* BUILD
  * Upgrade gitea/go-sdk to 2020-01-03 (#81)
  * Update stretchr/testify v1.3.0 -> v1.4.0 (#83)
  * Improve makefile to enable goproxy when go get tools (#98)

## [v0.1.2](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=59) - 2019-11-15
* BUILD
  * Fix typo in drone (#75)

## [v0.1.1](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=59) - 2019-11-15
* FEATURES
  * Add repos subcommand (#65)
* ENHANCEMENTS
  * Minor improvements to command-line language (#66)

## [v0.1.0](https://gitea.com/gitea/tea/pulls?q=&type=all&state=closed&milestone=59) - 2019-10-28
* BREAKING
  * Changed git config determination to go-git (#41) [continue #45] (#62)
* FEATURES
  * Add labels commands (#36)
* BUGFIXES
  * Fix out -o flag (#53)
  * Fix log formatting, refactor flag definition in cmd/labels.go (#52)
* ENHANCEMENTS
  * List label description (#60)
  * Use Different Remote Repos (#58)
  * Unified output (#14) (#40)
  * Added global appendable Flags (#12) (#39)
* BUILD
  * Change .drone.yml to new format (#33)
* DOCS
  * Add install guide from brew on README (#61)
