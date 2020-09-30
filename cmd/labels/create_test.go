// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package labels

import (
	"bufio"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLabelLine(t *testing.T) {
	const labels = `#ededed in progress
#fbca04 kind/breaking ; breaking label
#fc2929 kind/bug
#c5def5 kind/deployment ; deployment label
#000000 in progress ; in progress label
`

	scanner := bufio.NewScanner(strings.NewReader(labels))
	var i = 1
	for scanner.Scan() {
		line := scanner.Text()
		color, name, description := splitLabelLine(line)

		switch i {
		case 1:
			assert.EqualValues(t, "#ededed", color)
			assert.EqualValues(t, "in progress", name)
		case 2:
			assert.EqualValues(t, "#fbca04", color)
			assert.EqualValues(t, "kind/breaking", name)
			assert.EqualValues(t, "breaking label", description)
		case 3:
			assert.EqualValues(t, "#fc2929", color)
			assert.EqualValues(t, "kind/bug", name)
		case 4:
			assert.EqualValues(t, "#c5def5", color)
			assert.EqualValues(t, "kind/deployment", name)
			assert.EqualValues(t, "deployment label", description)
		case 5:
			assert.EqualValues(t, "#000000", color)
			assert.EqualValues(t, "in progress", name)
			assert.EqualValues(t, "in progress label", description)
		}

		i++
	}
}
