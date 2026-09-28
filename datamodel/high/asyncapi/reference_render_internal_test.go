// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
)

func TestPatchReferenceField_GuardsInvalidInputs(t *testing.T) {
	patchReferenceField(nil, "channel", nil, false)
	patchReferenceField(&yaml.Node{Kind: yaml.ScalarNode}, "channel", nil, false)

	root := &yaml.Node{Kind: yaml.MappingNode}
	patchReferenceField(root, "messages", []*low.Reference{nil}, true)
	assert.Empty(t, root.Content)
	patchReferenceField(root, "channel", []*low.Reference{nil}, false)
	assert.Empty(t, root.Content)
}
