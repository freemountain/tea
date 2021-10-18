// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"net/url"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	local_git "code.gitea.io/tea/modules/git"

	"github.com/go-git/go-git/v5"
	git_config "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

// RepoClone creates a local git clone in the given path, and sets up upstream remote
// for fork repos, for good usability with tea.
func RepoClone(
	path string,
	login *config.Login,
	repoOwner, repoName string,
	callback func(string) (string, error),
	depth int,
) (*local_git.TeaRepo, error) {

	repoMeta, _, err := login.Client().GetRepo(repoOwner, repoName)
	if err != nil {
		return nil, err
	}

	originURL, err := cloneURL(repoMeta, login)
	if err != nil {
		return nil, err
	}

	auth, err := local_git.GetAuthForURL(originURL, login.Token, login.SSHKey, callback)
	if err != nil {
		return nil, err
	}

	// default path behaviour as native git
	if path == "" {
		path = repoName
	}

	repo, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:             originURL.String(),
		Auth:            auth,
		Depth:           depth,
		InsecureSkipTLS: login.Insecure,
	})
	if err != nil {
		return nil, err
	}

	// set up upstream remote for forks
	if repoMeta.Fork && repoMeta.Parent != nil {
		upstreamURL, err := cloneURL(repoMeta.Parent, login)
		if err != nil {
			return nil, err
		}
		upstreamBranch := repoMeta.Parent.DefaultBranch
		repo.CreateRemote(&git_config.RemoteConfig{
			Name: "upstream",
			URLs: []string{upstreamURL.String()},
		})
		repoConf, err := repo.Config()
		if err != nil {
			return nil, err
		}
		if b, ok := repoConf.Branches[upstreamBranch]; ok {
			b.Remote = "upstream"
			b.Merge = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", upstreamBranch))
		}
		if err = repo.SetConfig(repoConf); err != nil {
			return nil, err
		}
	}

	return &local_git.TeaRepo{Repository: repo}, nil
}

func cloneURL(repo *gitea.Repository, login *config.Login) (*url.URL, error) {
	urlStr := repo.CloneURL
	if login.SSHKey != "" {
		urlStr = repo.SSHURL
	}
	return local_git.ParseURL(urlStr)
}
