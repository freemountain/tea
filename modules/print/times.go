// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
)

// TrackedTimesList print list of tracked times to stdout
func TrackedTimesList(times []*gitea.TrackedTime, outputType string, fields []string, printTotal bool) {
	var printables = make([]printable, len(times))
	var totalDuration int64
	for i, t := range times {
		totalDuration += t.Time
		printables[i] = &printableTrackedTime{t, outputType}
	}
	t := tableFromItems(fields, printables)

	if printTotal {
		total := make([]string, len(fields))
		total[0] = "TOTAL"
		total[len(fields)-1] = formatDuration(totalDuration, outputType)
		t.addRowSlice(total)
	}

	t.print(outputType)
}

// TrackedTimeFields contains all available fields for printing of tracked times.
var TrackedTimeFields = []string{
	"id",
	"created",
	"repo",
	"issue",
	"user",
	"duration",
}

type printableTrackedTime struct {
	*gitea.TrackedTime
	outputFormat string
}

func (t printableTrackedTime) FormatField(field string) string {
	switch field {
	case "id":
		return fmt.Sprintf("%d", t.ID)
	case "created":
		return FormatTime(t.Created)
	case "repo":
		return t.Issue.Repository.FullName
	case "issue":
		return fmt.Sprintf("#%d", t.Issue.Index)
	case "user":
		return t.UserName
	case "duration":
		return formatDuration(t.Time, t.outputFormat)
	}
	return ""
}
