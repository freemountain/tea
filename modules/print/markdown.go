// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package print

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"golang.org/x/crypto/ssh/terminal"
)

// outputMarkdown prints markdown to stdout, formatted for terminals.
// If the input could not be parsed, it is printed unformatted, the error
// is returned anyway.
func outputMarkdown(markdown string, baseURL string) error {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithBaseURL(baseURL),
		glamour.WithWordWrap(getWordWrap()),
	)
	if err != nil {
		fmt.Printf(markdown)
		return err
	}

	out, err := renderer.Render(markdown)
	if err != nil {
		fmt.Printf(markdown)
		return err
	}
	fmt.Print(out)
	return nil
}

// stolen from https://github.com/charmbracelet/glow/blob/e9d728c/main.go#L152-L165
func getWordWrap() int {
	fd := int(os.Stdout.Fd())
	width := 80
	if terminal.IsTerminal(fd) {
		if w, _, err := terminal.GetSize(fd); err == nil {
			width = w
		}
	}
	if width > 120 {
		width = 120
	}
	return width
}
