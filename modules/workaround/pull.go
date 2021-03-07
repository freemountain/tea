// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package workaround

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// FixPullHeadSha is a workaround for https://github.com/go-gitea/gitea/issues/12675
// When no head sha is available, this is because the branch got deleted in the base repo.
// pr.Head.Ref points in this case not to the head repo branch name, but the base repo ref,
// which stays available to resolve the commit sha.
func FixPullHeadSha(client *gitea.Client, pr *gitea.PullRequest) error {
	owner := pr.Base.Repository.Owner.UserName
	repo := pr.Base.Repository.Name
	if pr.Head != nil && pr.Head.Sha == "" {
		refs, _, err := client.GetRepoRefs(owner, repo, pr.Head.Ref)
		if err != nil {
			return err
		} else if len(refs) == 0 {
			return fmt.Errorf("unable to resolve PR ref '%s'", pr.Head.Ref)
		}
		pr.Head.Sha = refs[0].Object.SHA
	}
	return nil
}
