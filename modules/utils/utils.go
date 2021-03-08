// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

// Contains checks containment
func Contains(haystack []string, needle string) bool {
	return IndexOf(haystack, needle) != -1
}

// IndexOf returns the index of first occurrence of needle in haystack
func IndexOf(haystack []string, needle string) int {
	for i, s := range haystack {
		if s == needle {
			return i
		}
	}
	return -1
}
