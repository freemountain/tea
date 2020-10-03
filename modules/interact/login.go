// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"

	"code.gitea.io/tea/modules/config"

	"github.com/AlecAivazis/survey/v2"
)

// CreateLogin create an login interactive
func CreateLogin() error {
	var name, token, user, passwd, sshKey, giteaURL string
	var insecure = false

	promptI := &survey.Input{Message: "URL of Gitea instance: "}
	if err := survey.AskOne(promptI, &giteaURL, survey.WithValidator(survey.Required)); err != nil {
		return err
	}
	giteaURL = strings.TrimSuffix(strings.TrimSpace(giteaURL), "/")
	if len(giteaURL) == 0 {
		fmt.Println("URL is required!")
		return nil
	}

	name, err := config.GenerateLoginName(giteaURL, "")
	if err != nil {
		return err
	}

	promptI = &survey.Input{Message: "Name of new Login [" + name + "]: "}
	if err := survey.AskOne(promptI, &name); err != nil {
		return err
	}

	var hasToken bool
	promptYN := &survey.Confirm{
		Message: "Do you have an access token?",
		Default: false,
	}
	if err = survey.AskOne(promptYN, &hasToken); err != nil {
		return err
	}

	if hasToken {
		promptI = &survey.Input{Message: "Token: "}
		if err := survey.AskOne(promptI, &token, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	} else {
		promptI = &survey.Input{Message: "Username: "}
		if err = survey.AskOne(promptI, &user, survey.WithValidator(survey.Required)); err != nil {
			return err
		}

		promptPW := &survey.Password{Message: "Password: "}
		if err = survey.AskOne(promptPW, &passwd, survey.WithValidator(survey.Required)); err != nil {
			return err
		}
	}

	var optSettings bool
	promptYN = &survey.Confirm{
		Message: "Set Optional settings: ",
		Default: false,
	}
	if err = survey.AskOne(promptYN, &optSettings); err != nil {
		return err
	}
	if optSettings {
		promptI = &survey.Input{Message: "SSH Key Path: "}
		if err := survey.AskOne(promptI, &sshKey); err != nil {
			return err
		}

		promptYN = &survey.Confirm{
			Message: "Allow Insecure connections: ",
			Default: false,
		}
		if err = survey.AskOne(promptYN, &insecure); err != nil {
			return err
		}
	}

	return config.AddLogin(name, token, user, passwd, sshKey, giteaURL, insecure)
}
