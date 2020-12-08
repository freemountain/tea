// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strconv"
	"time"

	"code.gitea.io/sdk/gitea"
)

func formatDuration(seconds int64, outputType string) string {
	switch outputType {
	case "yaml":
	case "csv":
		return fmt.Sprint(seconds)
	}
	return time.Duration(1e9 * seconds).String()
}

// TrackedTimesList print list of tracked times to stdout
func TrackedTimesList(times []*gitea.TrackedTime, outputType string, from, until time.Time, printTotal bool) {
	var outputValues [][]string
	var totalDuration int64

	for _, t := range times {
		if !from.IsZero() && from.After(t.Created) {
			continue
		}
		if !until.IsZero() && until.Before(t.Created) {
			continue
		}

		totalDuration += t.Time

		outputValues = append(
			outputValues,
			[]string{
				FormatTime(t.Created),
				"#" + strconv.FormatInt(t.Issue.Index, 10),
				t.UserName,
				formatDuration(t.Time, outputType),
			},
		)
	}

	if printTotal {
		outputValues = append(outputValues, []string{
			"TOTAL", "", "", formatDuration(totalDuration, outputType),
		})
	}

	headers := []string{
		"Created",
		"Issue",
		"User",
		"Duration",
	}
	outputList(outputType, headers, outputValues)
}
