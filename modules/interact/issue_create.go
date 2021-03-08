// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package interact

import (
	"code.gitea.io/sdk/gitea"
	"code.gitea.io/tea/modules/config"
	"code.gitea.io/tea/modules/task"

	"github.com/AlecAivazis/survey/v2"
)

// CreateIssue interactively creates an issue
func CreateIssue(login *config.Login, owner, repo string) error {
	owner, repo, err := promptRepoSlug(owner, repo)
	if err != nil {
		return err
	}

	var opts gitea.CreateIssueOption
	if err := promptIssueProperties(login, owner, repo, &opts); err != nil {
		return err
	}

	return task.CreateIssue(login, owner, repo, opts)
}

func promptIssueProperties(login *config.Login, owner, repo string, o *gitea.CreateIssueOption) error {
	var milestoneName string
	var labels []string
	var err error

	selectableChan := make(chan (issueSelectables), 1)
	go fetchIssueSelectables(login, owner, repo, selectableChan)

	// title
	promptOpts := survey.WithValidator(survey.Required)
	promptI := &survey.Input{Message: "Issue title:", Default: o.Title}
	if err = survey.AskOne(promptI, &o.Title, promptOpts); err != nil {
		return err
	}

	// description
	promptD := &survey.Multiline{Message: "Issue description:", Default: o.Body}
	if err = survey.AskOne(promptD, &o.Body); err != nil {
		return err
	}

	// wait until selectables are fetched
	selectables := <-selectableChan
	if selectables.Err != nil {
		return selectables.Err
	}

	// skip remaining props if we don't have permission to set them
	if !selectables.Repo.Permissions.Push {
		return nil
	}

	// assignees
	if o.Assignees, err = promptMultiSelect("Assignees:", selectables.Collaborators, "[other]"); err != nil {
		return err
	}

	// milestone
	if len(selectables.MilestoneList) != 0 {
		if milestoneName, err = promptSelect("Milestone:", selectables.MilestoneList, "", "[none]"); err != nil {
			return err
		}
		o.Milestone = selectables.MilestoneMap[milestoneName]
	}

	// labels
	if len(selectables.LabelList) != 0 {
		promptL := &survey.MultiSelect{Message: "Labels:", Options: selectables.LabelList, VimMode: true, Default: o.Labels}
		if err := survey.AskOne(promptL, &labels); err != nil {
			return err
		}
		o.Labels = make([]int64, len(labels))
		for i, l := range labels {
			o.Labels[i] = selectables.LabelMap[l]
		}
	}

	// deadline
	if o.Deadline, err = promptDatetime("Due date:"); err != nil {
		return err
	}

	return nil
}

type issueSelectables struct {
	Repo          *gitea.Repository
	Collaborators []string
	MilestoneList []string
	MilestoneMap  map[string]int64
	LabelList     []string
	LabelMap      map[string]int64
	Err           error
}

func fetchIssueSelectables(login *config.Login, owner, repo string, done chan issueSelectables) {
	// TODO PERF make these calls concurrent
	r := issueSelectables{}
	c := login.Client()

	r.Repo, _, r.Err = c.GetRepo(owner, repo)
	if r.Err != nil {
		done <- r
		return
	}
	// we can set the following properties only if we have write access to the repo
	// so we fastpath this if not.
	if !r.Repo.Permissions.Push {
		done <- r
		return
	}

	// FIXME: this should ideally be ListAssignees(), https://github.com/go-gitea/gitea/issues/14856
	colabs, _, err := c.ListCollaborators(owner, repo, gitea.ListCollaboratorsOptions{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.Collaborators = make([]string, len(colabs)+1)
	r.Collaborators[0] = login.User
	for i, u := range colabs {
		r.Collaborators[i+1] = u.UserName
	}

	milestones, _, err := c.ListRepoMilestones(owner, repo, gitea.ListMilestoneOption{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.MilestoneMap = make(map[string]int64)
	r.MilestoneList = make([]string, len(milestones))
	for i, m := range milestones {
		r.MilestoneMap[m.Title] = m.ID
		r.MilestoneList[i] = m.Title
	}

	labels, _, err := c.ListRepoLabels(owner, repo, gitea.ListLabelsOptions{})
	if err != nil {
		r.Err = err
		done <- r
		return
	}
	r.LabelMap = make(map[string]int64)
	r.LabelList = make([]string, len(labels))
	for i, l := range labels {
		r.LabelMap[l.Name] = l.ID
		r.LabelList[i] = l.Name
	}

	done <- r
}
