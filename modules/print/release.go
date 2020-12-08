// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"code.gitea.io/sdk/gitea"
)

// ReleasesList prints a listing of releases
func ReleasesList(releases []*gitea.Release, output string) {
	var values [][]string
	headers := []string{
		"Tag-Name",
		"Title",
		"Published At",
		"Status",
		"Tar URL",
	}

	if len(releases) == 0 {
		outputList(output, headers, values)
		return
	}

	for _, release := range releases {
		status := "released"
		if release.IsDraft {
			status = "draft"
		} else if release.IsPrerelease {
			status = "prerelease"
		}
		values = append(
			values,
			[]string{
				release.TagName,
				release.Title,
				FormatTime(release.PublishedAt),
				status,
				release.TarURL,
			},
		)
	}

	outputList(output, headers, values)
}
