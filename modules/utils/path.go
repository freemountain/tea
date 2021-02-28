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
	return exists(fileName, false)
}

// DirExists returns whether the given file exists or not
func DirExists(path string) (bool, error) {
	return exists(path, true)
}

func exists(path string, expectDir bool) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		} else if err.(*os.PathError).Err.Error() == "not a directory" {
			// some middle segment of path is a file, cannot traverse
			// FIXME: catches error on linux; go does not provide a way to catch this properly..
			return false, nil
		}
		return false, err
	}
	isDir := f.IsDir()
	if isDir && !expectDir {
		return false, errors.New("A directory with the same name exists")
	} else if !isDir && expectDir {
		return false, errors.New("A file with the same name exists")
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
