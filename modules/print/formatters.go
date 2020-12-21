// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"time"

	"code.gitea.io/sdk/gitea"
	"github.com/muesli/termenv"
)

// formatSize get kb in int and return string
func formatSize(kb int64) string {
	if kb < 1024 {
		return fmt.Sprintf("%d Kb", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%d Mb", mb)
	}
	gb := mb / 1024
	if gb < 1024 {
		return fmt.Sprintf("%d Gb", gb)
	}
	return fmt.Sprintf("%d Tb", gb/1024)
}

// FormatTime give a date-time in local timezone if available
func FormatTime(t time.Time) string {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return t.Format("2006-01-02 15:04 UTC")
	}
	return t.In(location).Format("2006-01-02 15:04")
}

func formatDuration(seconds int64, outputType string) string {
	if isMachineReadable(outputType) {
		return fmt.Sprint(seconds)
	}
	return time.Duration(1e9 * seconds).String()
}

func formatLabel(label *gitea.Label, allowColor bool, text string) string {
	colorProfile := termenv.Ascii
	if allowColor {
		colorProfile = termenv.EnvColorProfile()
	}
	if len(text) == 0 {
		text = label.Name
	}
	styled := termenv.String(text)
	styled = styled.Foreground(colorProfile.Color("#" + label.Color))
	return fmt.Sprint(styled)
}

func formatPermission(p *gitea.Permission) string {
	if p.Admin {
		return "admin"
	} else if p.Push {
		return "write"
	}
	return "read"
}

func formatUserName(u *gitea.User) string {
	if len(u.FullName) == 0 {
		return u.UserName
	}
	return u.FullName
}
