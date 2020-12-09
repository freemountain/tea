// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/tea/modules/config"
)

// LoginDetails print login entry to stdout
func LoginDetails(login *config.Login, output string) {
	in := fmt.Sprintf("# %s\n\n[@%s](%s/%s)\n",
		login.Name,
		login.User,
		strings.TrimSuffix(login.URL, "/"),
		login.User,
	)
	if len(login.SSHKey) != 0 {
		in += fmt.Sprintf("\nSSH Key: '%s' via %s\n",
			login.SSHKey,
			login.SSHHost,
		)
	}
	in += fmt.Sprintf("\nCreated: %s", time.Unix(login.Created, 0).Format(time.RFC822))

	outputMarkdown(in)
}

// LoginsList prints a listing of logins
func LoginsList(logins []config.Login, output string) {
	t := tableWithHeader(
		"Name",
		"URL",
		"SSHHost",
		"User",
		"Default",
	)

	for _, l := range logins {
		t.addRow(
			l.Name,
			l.URL,
			l.GetSSHHost(),
			l.User,
			fmt.Sprint(l.Default),
		)
	}

	t.print(output)
}
