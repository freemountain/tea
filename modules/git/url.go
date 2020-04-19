// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	protocolRe = regexp.MustCompile("^[a-zA-Z_+-]+://")
)

// URLParser represents a git URL parser
type URLParser struct {
}

// Parse parses the git URL
func (p *URLParser) Parse(rawURL string) (u *url.URL, err error) {
	if !protocolRe.MatchString(rawURL) &&
		strings.Contains(rawURL, ":") &&
		// not a Windows path
		!strings.Contains(rawURL, "\\") {
		rawURL = "ssh://" + strings.Replace(rawURL, ":", "/", 1)
	}

	u, err = url.Parse(rawURL)
	if err != nil {
		return
	}

	if u.Scheme == "git+ssh" {
		u.Scheme = "ssh"
	}

	if strings.HasPrefix(u.Path, "//") {
		u.Path = strings.TrimPrefix(u.Path, "/")
	}

	// .git suffix is optional and breaks normalization
	if strings.HasSuffix(u.Path, ".git") {
		u.Path = strings.TrimSuffix(u.Path, ".git")
	}

	return
}

// ParseURL parses URL string and return URL struct
func ParseURL(rawURL string) (u *url.URL, err error) {
	p := &URLParser{}
	return p.Parse(rawURL)
}
