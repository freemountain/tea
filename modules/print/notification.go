// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
)

// NotificationsList prints a listing of notification threads
func NotificationsList(news []*gitea.NotificationThread, output string, showRepository bool) {
	headers := []string{
		"ID",
		"Status",
		"Type",
		"State",
		"Index",
		"Title",
	}
	if showRepository {
		headers = append(headers, "Repository")
	}

	t := table{headers: headers}

	for _, n := range news {
		if n.Subject == nil {
			continue
		}
		// if pull or Issue get Index
		var index string
		if n.Subject.Type == "Issue" || n.Subject.Type == "Pull" {
			index = n.Subject.URL
			urlParts := strings.Split(n.Subject.URL, "/")
			if len(urlParts) != 0 {
				index = urlParts[len(urlParts)-1]
			}
			index = "#" + index
		}

		status := "read"
		if n.Pinned {
			status = "pinned"
		} else if n.Unread {
			status = "unread"
		}

		item := []string{
			fmt.Sprint(n.ID),
			status,
			string(n.Subject.Type),
			string(n.Subject.State),
			index,
			n.Subject.Title,
		}
		if showRepository {
			item = append(item, n.Repository.FullName)
		}
		t.addRowSlice(item)
	}

	if t.Len() != 0 {
		t.print(output)
	}
}
