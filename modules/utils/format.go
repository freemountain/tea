// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import "fmt"

// FormatSize get kb in int and return string
func FormatSize(kb int64) string {
	if kb < 1024 {
		return fmt.Sprintf("%d Kb", kb)
	}
	mb := kb / 1024
	if mb < 1024 {
		return fmt.Sprintf("%d Mb", mb)
	}
	gb := mb / 1024
	if gb < 1024 {
		return fmt.Sprintf("%d Gb", gb)
	}
	return fmt.Sprintf("%d Tb", gb/1024)
}
