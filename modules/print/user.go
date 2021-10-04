// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// UserDetails print a formatted user to stdout
func UserDetails(user *gitea.User) {
	title := "# " + user.UserName
	if user.IsAdmin {
		title += " (admin)"
	}
	if !user.IsActive {
		title += " (disabled)"
	}
	if user.Restricted {
		title += " (restricted)"
	}
	if user.ProhibitLogin {
		title += " (login prohibited)"
	}
	title += "\n"

	var desc string
	if len(user.Description) != 0 {
		desc = fmt.Sprintf("*%s*\n\n", user.Description)
	}
	var website string
	if len(user.Website) != 0 {
		website = fmt.Sprintf("%s\n\n", user.Website)
	}

	stats := fmt.Sprintf(
		"Follower Count: %d, Following Count: %d, Starred Repos: %d\n",
		user.FollowerCount,
		user.FollowingCount,
		user.StarredRepoCount,
	)

	outputMarkdown(fmt.Sprintf(
		"%s%s\n%s\n%s",
		title,
		desc,
		website,
		stats,
	), "")
}

// UserList prints a listing of the users
func UserList(user []*gitea.User, output string, fields []string) {
	var printables = make([]printable, len(user))
	for i, u := range user {
		printables[i] = &printableUser{u}
	}
	t := tableFromItems(fields, printables, isMachineReadable(output))
	t.print(output)
}

// UserFields are the available fields to print with UserList()
var UserFields = []string{
	"id",
	"login",
	"full_name",
	"email",
	"avatar_url",
	"language",
	"is_admin",
	"restricted",
	"prohibit_login",
	"location",
	"website",
	"description",
	"visibility",
}

type printableUser struct{ *gitea.User }

func (x printableUser) FormatField(field string, machineReadable bool) string {
	switch field {
	case "id":
		return fmt.Sprintf("%d", x.ID)
	case "login":
		if x.IsAdmin {
			return fmt.Sprintf("%s (admin)", x.UserName)
		}
		if !x.IsActive {
			return fmt.Sprintf("%s (disabled)", x.UserName)
		}
		if x.Restricted {
			return fmt.Sprintf("%s (restricted)", x.UserName)
		}
		if x.ProhibitLogin {
			return fmt.Sprintf("%s (login prohibited)", x.UserName)
		}
		return x.UserName
	case "full_name":
		return x.FullName
	case "email":
		return x.Email
	case "avatar_url":
		return x.AvatarURL
	case "language":
		return x.Language
	case "is_admin":
		return formatBoolean(x.IsAdmin, !machineReadable)
	case "restricted":
		return formatBoolean(x.Restricted, !machineReadable)
	case "prohibit_login":
		return formatBoolean(x.ProhibitLogin, !machineReadable)
	case "location":
		return x.Location
	case "website":
		return x.Website
	case "description":
		return x.Description
	case "visibility":
		return string(x.Visibility)
	}
	return ""
}
