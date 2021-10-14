// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"fmt"
	"strings"
	"time"

	"code.gitea.io/tea/modules/utils"
	"github.com/AlecAivazis/survey/v2"
	"github.com/araddon/dateparse"
)

// Multiline represents options for a prompt that expects multiline input
type Multiline struct {
	Message   string
	Default   string
	Syntax    string
	UseEditor bool
}

// NewMultiline creates a prompt that switches between the inline multiline text
// and a texteditor based prompt
func NewMultiline(opts Multiline) (prompt survey.Prompt) {
	if opts.UseEditor {
		prompt = &survey.Editor{
			Message:  opts.Message,
			Default:  opts.Default,
			FileName: "*." + opts.Syntax,
		}
	} else {
		prompt = &survey.Multiline{Message: opts.Message, Default: opts.Default}
	}
	return
}

// PromptPassword asks for a password and blocks until input was made.
func PromptPassword(name string) (pass string, err error) {
	promptPW := &survey.Password{Message: name + " password:"}
	err = survey.AskOne(promptPW, &pass, survey.WithValidator(survey.Required))
	return
}

// promptRepoSlug interactively prompts for a Gitea repository or returns the current one
func promptRepoSlug(defaultOwner, defaultRepo string) (owner, repo string, err error) {
	prompt := "Target repo:"
	defaultVal := ""
	required := true
	if len(defaultOwner) != 0 && len(defaultRepo) != 0 {
		defaultVal = fmt.Sprintf("%s/%s", defaultOwner, defaultRepo)
		required = false
	}
	var repoSlug string

	owner = defaultOwner
	repo = defaultRepo

	err = survey.AskOne(
		&survey.Input{
			Message: prompt,
			Default: defaultVal,
		},
		&repoSlug,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if !required && len(str) == 0 {
					return nil
				}
				split := strings.Split(str, "/")
				if len(split) != 2 || len(split[0]) == 0 || len(split[1]) == 0 {
					return fmt.Errorf("must follow the <owner>/<repo> syntax")
				}
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)

	if err == nil && len(repoSlug) != 0 {
		repoSlugSplit := strings.Split(repoSlug, "/")
		owner = repoSlugSplit[0]
		repo = repoSlugSplit[1]
	}
	return
}

// promptDatetime prompts for a date or datetime string.
// Supports all formats understood by araddon/dateparse.
func promptDatetime(prompt string) (val *time.Time, err error) {
	var input string
	err = survey.AskOne(
		&survey.Input{Message: prompt},
		&input,
		survey.WithValidator(func(input interface{}) error {
			if str, ok := input.(string); ok {
				if len(str) == 0 {
					return nil
				}
				t, err := dateparse.ParseAny(str)
				if err != nil {
					return err
				}
				val = &t
			} else {
				return fmt.Errorf("invalid result type")
			}
			return nil
		}),
	)
	return
}

// promptSelect creates a generic multiselect prompt, with processing of custom values.
func promptMultiSelect(prompt string, options []string, customVal string) ([]string, error) {
	var selection []string
	promptA := &survey.MultiSelect{
		Message: prompt,
		Options: makeSelectOpts(options, customVal, ""),
		VimMode: true,
	}
	if err := survey.AskOne(promptA, &selection); err != nil {
		return nil, err
	}
	return promptCustomVal(prompt, customVal, selection)
}

// promptSelect creates a generic select prompt, with processing of custom values or none-option.
func promptSelect(prompt string, options []string, customVal, noneVal string) (string, error) {
	var selection string
	promptA := &survey.Select{
		Message: prompt,
		Options: makeSelectOpts(options, customVal, noneVal),
		VimMode: true,
		Default: noneVal,
	}
	if err := survey.AskOne(promptA, &selection); err != nil {
		return "", err
	}
	if noneVal != "" && selection == noneVal {
		return "", nil
	}
	if customVal != "" {
		sel, err := promptCustomVal(prompt, customVal, []string{selection})
		if err != nil {
			return "", err
		}
		selection = sel[0]
	}
	return selection, nil
}

// makeSelectOpts adds cusotmVal & noneVal to opts if set.
func makeSelectOpts(opts []string, customVal, noneVal string) []string {
	if customVal != "" {
		opts = append(opts, customVal)
	}
	if noneVal != "" {
		opts = append(opts, noneVal)
	}
	return opts
}

// promptCustomVal checks if customVal is present in selection, and prompts
// for custom input to add to the selection instead.
func promptCustomVal(prompt, customVal string, selection []string) ([]string, error) {
	// check for custom value & prompt again with text input
	// HACK until https://github.com/AlecAivazis/survey/issues/339 is implemented
	if otherIndex := utils.IndexOf(selection, customVal); otherIndex != -1 {
		var customAssignees string
		promptA := &survey.Input{Message: prompt, Help: "comma separated list"}
		if err := survey.AskOne(promptA, &customAssignees); err != nil {
			return nil, err
		}
		selection = append(selection[:otherIndex], selection[otherIndex+1:]...)
		selection = append(selection, strings.Split(customAssignees, ",")...)
	}
	return selection, nil
}
