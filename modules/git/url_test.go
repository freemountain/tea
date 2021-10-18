// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseUrl(t *testing.T) {
	u, err := ParseURL("ssh://git@gitea.com:3000/gitea/tea")
	assert.NoError(t, err)
	assert.Equal(t, "gitea.com:3000", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "/gitea/tea", u.Path)

	u, err = ParseURL("https://gitea.com/gitea/tea")
	assert.NoError(t, err)
	assert.Equal(t, "gitea.com", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/gitea/tea", u.Path)

	u, err = ParseURL("git@gitea.com:gitea/tea")
	assert.NoError(t, err)
	assert.Equal(t, "gitea.com", u.Host)
	assert.Equal(t, "ssh", u.Scheme)
	assert.Equal(t, "/gitea/tea", u.Path)

	u, err = ParseURL("gitea.com/gitea/tea")
	assert.NoError(t, err)
	assert.Equal(t, "gitea.com", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/gitea/tea", u.Path)

	u, err = ParseURL("foo/bar")
	assert.NoError(t, err)
	assert.Equal(t, "", u.Host)
	assert.Equal(t, "", u.Scheme)
	assert.Equal(t, "foo/bar", u.Path)

	u, err = ParseURL("/foo/bar")
	assert.NoError(t, err)
	assert.Equal(t, "", u.Host)
	assert.Equal(t, "https", u.Scheme)
	assert.Equal(t, "/foo/bar", u.Path)

	// this case is unintuitive, but to ambiguous to be handled differently
	u, err = ParseURL("gitea.com")
	assert.NoError(t, err)
	assert.Equal(t, "", u.Host)
	assert.Equal(t, "", u.Scheme)
	assert.Equal(t, "gitea.com", u.Path)
}
