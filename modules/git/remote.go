// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"fmt"
	"net/url"

	"github.com/go-git/go-git/v5"
	git_config "github.com/go-git/go-git/v5/config"
)

// GetRemote tries to match a Remote of the repo via the given URL.
// Matching is based on the normalized URL, accepting different protocols.
func (r TeaRepo) GetRemote(remoteURL string) (*git.Remote, error) {
	repoURL, err := ParseURL(remoteURL)
	if err != nil {
		return nil, err
	}

	remotes, err := r.Remotes()
	if err != nil {
		return nil, err
	}
	for _, r := range remotes {
		for _, u := range r.Config().URLs {
			remoteURL, err := ParseURL(u)
			if err != nil {
				return nil, err
			}
			if remoteURL.Host == repoURL.Host && remoteURL.Path == repoURL.Path {
				return r, nil
			}
		}
	}

	return nil, nil
}

// GetOrCreateRemote tries to match a Remote of the repo via the given URL.
// If no match is found, a new Remote with `newRemoteName` is created.
// Matching is based on the normalized URL, accepting different protocols.
func (r TeaRepo) GetOrCreateRemote(remoteURL, newRemoteName string) (*git.Remote, error) {
	localRemote, err := r.GetRemote(remoteURL)
	if err != nil {
		return nil, err
	}

	// if no match found, create a new remote
	if localRemote == nil {
		localRemote, err = r.CreateRemote(&git_config.RemoteConfig{
			Name: newRemoteName,
			URLs: []string{remoteURL},
		})
		if err != nil {
			return nil, err
		}
	}

	return localRemote, nil
}

// TeaRemoteURL returns the first url entry for the given remote name
func (r TeaRepo) TeaRemoteURL(name string) (auth *url.URL, err error) {
	remote, err := r.Remote(name)
	if err != nil {
		return nil, err
	}
	urls := remote.Config().URLs
	if len(urls) == 0 {
		return nil, fmt.Errorf("remote %s has no URL configured", name)
	}
	return ParseURL(remote.Config().URLs[0])
}
