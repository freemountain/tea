// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	go_git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// GetRepoReference returns the current repository's current branch or tag
func GetRepoReference(p string) (*plumbing.Reference, error) {
	gitPath, err := go_git.PlainOpenWithOptions(p, &go_git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}

	return gitPath.Head()
}
