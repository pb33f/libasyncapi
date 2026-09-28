// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/utils"
)

// patchReferenceField replaces NodeBuilder's structural rendering of low.Reference with
// proper AsyncAPI Reference Objects. low.Reference stores its state in unexported fields,
// so generic YAML encoding otherwise emits an empty mapping.
func patchReferenceField(root *yaml.Node, label string, refs []*low.Reference, sequence bool) {
	if root == nil || root.Kind != yaml.MappingNode || len(refs) == 0 {
		return
	}
	var rendered *yaml.Node
	if sequence {
		rendered = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, ref := range refs {
			if ref != nil && ref.GetReference() != "" {
				rendered.Content = append(rendered.Content, utils.CreateRefNode(ref.GetReference()))
			}
		}
		if len(rendered.Content) == 0 {
			return
		}
	} else {
		ref := refs[0]
		if ref == nil || ref.GetReference() == "" {
			return
		}
		rendered = utils.CreateRefNode(ref.GetReference())
	}

	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == label {
			root.Content[i+1] = rendered
			return
		}
	}
}
