// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libopenapi/index"
)

// schemaIndexView hides foreign-language content from the JSON-reference index.
// Only changed ancestors are copied. Source nodes are never modified.
func schemaIndexView(root *yaml.Node) (*yaml.Node, map[*yaml.Node]*yaml.Node) {
	originals := make(map[*yaml.Node]*yaml.Node)
	var visit func(*yaml.Node, []string) *yaml.Node
	visit = func(node *yaml.Node, path []string) *yaml.Node {
		foreign := false
		if multiFormatLocation(path) {
			format := pointerChild(node, "schemaFormat")
			foreign = format != nil && !schemaformat.IsJSON(format.Value)
		}
		var content []*yaml.Node
		for i, child := range node.Content {
			var next *yaml.Node
			if foreign && i%2 == 1 && node.Content[i-1].Value == "schema" && !(child.Kind == yaml.MappingNode && len(child.Content) == 2 && pointerChild(child, "$ref") != nil) {
				masked := *child
				masked.Kind, masked.Tag, masked.Value, masked.Content = yaml.ScalarNode, "!!str", "", nil
				masked.Alias = nil
				next = &masked
				originals[next] = child
			} else {
				childPath := path
				if node.Kind == yaml.MappingNode && i%2 == 1 {
					childPath = append(path, node.Content[i-1].Value)
				}
				next = visit(child, childPath)
			}
			if next != child && content == nil {
				content = append([]*yaml.Node(nil), node.Content...)
			}
			if content != nil {
				content[i] = next
			}
		}
		if content == nil {
			return node
		}
		out := *node
		out.Content = content
		originals[&out] = node
		return &out
	}
	return visit(root, nil), originals
}

func multiFormatLocation(path []string) bool {
	if len(path) == 3 && path[0] == "components" && path[1] == "schemas" {
		return true
	}
	if len(path) == 4 && path[0] == "components" && path[1] == "messages" && path[3] == "payload" {
		return true
	}
	if len(path) > 0 && path[0] == "components" {
		path = path[1:]
	}
	return len(path) == 5 && path[0] == "channels" && path[2] == "messages" && path[4] == "payload"
}

// restoreIndexSources publishes original nodes after reference scanning finishes.
func restoreIndexSources(rolodex *index.Rolodex, root *yaml.Node, originals map[*yaml.Node]*yaml.Node) {
	if len(originals) == 0 {
		return
	}
	for _, idx := range append(rolodex.GetIndexes(), rolodex.GetRootIndex()) {
		if idx == nil {
			continue
		}
		restore := func(ref *index.Reference) {
			if node := originals[ref.Node]; node != nil {
				ref.Node = node
			}
			if node := originals[ref.ParentNode]; node != nil {
				ref.ParentNode = node
			}
		}
		for _, ref := range idx.GetMappedReferences() {
			restore(ref)
		}
		for _, ref := range idx.GetAllReferences() {
			restore(ref)
		}
		for _, ref := range idx.GetRawReferencesSequenced() {
			restore(ref)
		}
	}
	rolodex.GetRootIndex().SetRootNode(root)
	rolodex.SetRootNode(root)
}
