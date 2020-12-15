# comparing git forge commandline interfaces

[tea]: https://gitea.com/gitea/tea
[sip]: https://gitea.com/jolheiser/sip
[gitlab]: https://github.com/makkes/gitlab-cli
[glab]: https://github.com/profclems/glab
[gh]: https://cli.github.com

last update: 2020-12-11

## general
/ | [tea][tea] | [sip][sip] | [gitlab][gitlab] | [gh][gh]
-----------------------|:-----:|:-----:|:-----:|:-----:
forge|gitea|gitea|gitlab|github
official forge support|✓|✘|✘|✓
dev status|adding features|maintenance||
platform|any|any|any|any

## philosophy
/ | [tea][tea] | [sip][sip] | [gitlab][gitlab] | [gh][gh]
-----------------------|:-----:|:-----:|:-----:|:-----:
aims to replace git cli|✘|||✓
works with decentralization in mind|✓|✓|✓|✘
per-repo setup needed|✘||✓|✘
workflow helpers|✓|||
interactive mode |[(✓)](https://gitea.com/gitea/tea/issues?type=all&state=open&labels=&milestone=0&assignee=0&q=interactive)|✘| |✓
programmatic mode|✓|||✓
machine readable output|✓|||
follows XDG spec|✓|||

## features
/ | [tea][tea] | [sip][sip] | [gitlab][gitlab] | [gh][gh]
-----------------------|:-----:|:-----:|:-----:|:-----:
open web UI|✓|||
search repos|✓|||
search issues|✘|✓||
textual item search filter syntax|✘|✓||
CRUD repos|[(✓)](https://gitea.com/gitea/tea/issues/239)|||
CRUD issues|[(✓)](https://gitea.com/gitea/tea/issues/229)|||
CRUD milestones|[(✓)](https://gitea.com/gitea/tea/issues/246)|||
CRUD releases|✓|||
CRUD labels|✓|||
CRUD PRs|✓|||
CRUD time tracking|✓|||x
CRUD orgs|[(✓)](https://gitea.com/gitea/tea/issues/287)|||
create PRs from local repo|✓|||
create PRs from remote repo|✓|||
code review|[u](https://gitea.com/gitea/tea/issues/131)|||
merge PRs||||
read comments|[u](https://gitea.com/gitea/tea/issues/172)|||
post comments||||
manage CI|✘|✘|✓|
manage notifications|[(✓)]()|||
administration|[u](https://gitea.com/gitea/tea/issues/161)|✘||✘
markdown rendering|✓|||✓
issue import/export|[u](https://gitea.com/gitea/tea/issues/132)|||
checkout PRs|✓|||

- ✓: supported
- (✓): partial support
- u: upcoming
- ✘: not supported
- ?: unknown
