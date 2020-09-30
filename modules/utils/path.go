// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// PathExists returns whether the given file or directory exists or not
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FileExist returns whether the given file exists or not
func FileExist(fileName string) (bool, error) {
	f, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if f.IsDir() {
		return false, errors.New("A directory with the same name exists")
	}
	return true, nil
}

// AbsPathWithExpansion expand path beginning with "~/" to absolute path
func AbsPathWithExpansion(p string) (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	if p == "~" {
		return u.HomeDir, nil
	} else if strings.HasPrefix(p, "~/") {
		return filepath.Join(u.HomeDir, p[2:]), nil
	} else {
		return filepath.Abs(p)
	}
}
