// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/tea/modules/config"

	"github.com/charmbracelet/glamour"
)

// LoginDetails print login entry to stdout
func LoginDetails(login *config.Login) {

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

	out, err := glamour.Render(in, getGlamourTheme())
	if err != nil {
		// TODO: better Error handling
		fmt.Printf("Error:\n%v\n\n", err)
		return
	}
	fmt.Print(out)
}
