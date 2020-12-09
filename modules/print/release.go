// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"code.gitea.io/sdk/gitea"
)

// ReleasesList prints a listing of releases
func ReleasesList(releases []*gitea.Release, output string) {
	t := tableWithHeader(
		"Tag-Name",
		"Title",
		"Published At",
		"Status",
		"Tar URL",
	)

	for _, release := range releases {
		status := "released"
		if release.IsDraft {
			status = "draft"
		} else if release.IsPrerelease {
			status = "prerelease"
		}
		t.addRow(
			release.TagName,
			release.Title,
			FormatTime(release.PublishedAt),
			status,
			release.TarURL,
		)
	}

	t.print(output)
}
