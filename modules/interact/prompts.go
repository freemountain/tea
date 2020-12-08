// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"github.com/AlecAivazis/survey/v2"
)

// PromptPassword asks for a password and blocks until input was made.
func PromptPassword(name string) (pass string, err error) {
	promptPW := &survey.Password{Message: name + " password:"}
	err = survey.AskOne(promptPW, &pass, survey.WithValidator(survey.Required))
	return
}
