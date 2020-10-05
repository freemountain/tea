// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

// OutputMarkdown prints markdown to stdout, formatted for terminals.
// If the input could not be parsed, it is printed unformatted, the error
// is returned anyway.
func OutputMarkdown(markdown string) error {
	out, err := glamour.Render(markdown, "auto")
	if err != nil {
		fmt.Printf(markdown)
		return err
	}
	fmt.Print(out)
	return nil
}
