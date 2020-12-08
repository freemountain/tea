// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package task

import (
	"fmt"
	"log"
	"os"

	"code.gitea.io/sdk/gitea"
)

// LabelsExport save list of labels to disc
func LabelsExport(labels []*gitea.Label, path string) error {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, label := range labels {
		if _, err := fmt.Fprintf(f, "#%s %s\n", label.Color, label.Name); err != nil {
			return err
		}
	}
	return nil
}
