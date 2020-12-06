// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/cmd/flags"
)

// OrganizationsList prints a listing of the organizations
func OrganizationsList(organizations []*gitea.Organization) {
	if len(organizations) == 0 {
		fmt.Println("No organizations found")
		return
	}

	headers := []string{
		"Name",
		"FullName",
		"Website",
		"Location",
		"Description",
	}

	var values [][]string

	for _, org := range organizations {
		values = append(
			values,
			[]string{
				org.UserName,
				org.FullName,
				org.Website,
				org.Location,
				org.Description,
			},
		)
	}

	OutputList(flags.GlobalOutputValue, headers, values)
}
