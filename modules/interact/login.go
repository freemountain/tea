// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/config"
)

// CreateLogin create an login interactive
func CreateLogin() error {
	var stdin, name, token, user, passwd, sshKey, giteaURL string
	var insecure = false

	fmt.Print("URL of Gitea instance: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	giteaURL = strings.TrimSpace(stdin)
	if len(giteaURL) == 0 {
		fmt.Println("URL is required!")
		return nil
	}

	name, err := config.GenerateLoginName(giteaURL, "")
	if err != nil {
		return err
	}

	fmt.Print("Name of new Login [" + name + "]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(strings.TrimSpace(stdin)) != 0 {
		name = strings.TrimSpace(stdin)
	}

	fmt.Print("Do you have a token [Yes/no]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(stdin) != 0 && strings.ToLower(stdin[:1]) == "n" {
		fmt.Print("Username: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		user = strings.TrimSpace(stdin)

		fmt.Print("Password: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		passwd = strings.TrimSpace(stdin)
	} else {
		fmt.Print("Token: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		token = strings.TrimSpace(stdin)
	}

	fmt.Print("Set Optional settings [yes/No]: ")
	if _, err := fmt.Scanln(&stdin); err != nil {
		stdin = ""
	}
	if len(stdin) != 0 && strings.ToLower(stdin[:1]) == "y" {
		fmt.Print("SSH Key Path: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		sshKey = strings.TrimSpace(stdin)

		fmt.Print("Allow Insecure connections  [yes/No]: ")
		if _, err := fmt.Scanln(&stdin); err != nil {
			stdin = ""
		}
		insecure = len(stdin) != 0 && strings.ToLower(stdin[:1]) == "y"
	}

	return config.AddLogin(name, token, user, passwd, sshKey, giteaURL, insecure)
}
