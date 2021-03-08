// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/utils"
)

// ResolveLabelNames returns a list of label IDs for a given list of label names
func ResolveLabelNames(client *gitea.Client, owner, repo string, labelNames []string) ([]int64, error) {
	labelIDs := make([]int64, len(labelNames))
	labels, _, err := client.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{})
	if err != nil {
		return nil, err
	}
	for _, l := range labels {
		if utils.Contains(labelNames, l.Name) {
			labelIDs = append(labelIDs, l.ID)
		}
	}
	return labelIDs, nil
}
