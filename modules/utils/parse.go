// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"net/url"
	"strconv"
	"strings"
)

// ArgToIndex take issue/pull index as string and return int64
func ArgToIndex(arg string) (int64, error) {
	if strings.HasPrefix(arg, "#") {
		arg = arg[1:]
	}
	return strconv.ParseInt(arg, 10, 64)
}

// NormalizeURL normalizes the input with a protocol
func NormalizeURL(raw string) (*url.URL, error) {
	var prefix string
	if !strings.HasPrefix(raw, "http") {
		prefix = "https://"
	}
	return url.Parse(strings.TrimSuffix(prefix+raw, "/"))
}

// GetOwnerAndRepo return repoOwner and repoName
// based on relative path and default owner (if not in path)
func GetOwnerAndRepo(repoPath, user string) (string, string) {
	if len(repoPath) == 0 {
		return "", ""
	}
	p := strings.Split(strings.TrimLeft(repoPath, "/"), "/")
	if len(p) >= 2 {
		return p[0], p[1]
	}
	return user, repoPath
}
