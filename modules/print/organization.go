// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// OrganizationDetails prints details of an org with formatting
func OrganizationDetails(org *gitea.Organization) {
	outputMarkdown(fmt.Sprintf(
		"# %s\n%s\n\n- Visibility: %s\n- Location: %s\n- Website: %s\n",
		org.UserName,
		org.Description,
		org.Visibility,
		org.Location,
		org.Website,
	), "")
}

// OrganizationsList prints a listing of the organizations
func OrganizationsList(organizations []*gitea.Organization, output string) {
	if len(organizations) == 0 {
		fmt.Println("No organizations found")
		return
	}

	t := tableWithHeader(
		"Name",
		"FullName",
		"Website",
		"Location",
		"Description",
	)

	for _, org := range organizations {
		t.addRow(
			org.UserName,
			org.FullName,
			org.Website,
			org.Location,
			org.Description,
		)
	}

	t.print(output)
}
