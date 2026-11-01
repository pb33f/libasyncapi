// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libasyncapi/internal/yamlparse"
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
			format := pointerChild(node, SchemaFormatLabel)
			foreign = format != nil && !schemaformat.IsJSON(format.Value)
		}
		var content []*yaml.Node
		for i, child := range node.Content {
			var next *yaml.Node
			if foreign && i%2 == 1 && node.Content[i-1].Value == SchemaLabel && !(child.Kind == yaml.MappingNode && len(child.Content) == 2 && pointerChild(child, RefLabel) != nil) {
				masked := *child
				masked.Kind, masked.Tag, masked.Value, masked.Content = yaml.ScalarNode, yamlparse.StringTag, "", nil
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
	if len(path) > 0 && path[0] == ComponentsLabel {
		path = path[1:]
		switch {
		case len(path) == 2 && path[0] == SchemasLabel:
			return true
		case len(path) == 3 && path[0] == MessageTraitsLabel:
			return path[2] == HeadersLabel
		case len(path) > 2 && path[0] == MessagesLabel:
			return messageSchemaField(path[2:])
		}
	}
	return len(path) > 4 && path[0] == ChannelsLabel && path[2] == MessagesLabel && messageSchemaField(path[4:])
}

// messageSchemaField reports whether a path below a message is a schema field.
// Sequences add no path segment, so every trait in a list shares one path.
func messageSchemaField(path []string) bool {
	if len(path) == 2 {
		return path[0] == TraitsLabel && path[1] == HeadersLabel
	}
	return len(path) == 1 && (path[0] == PayloadLabel || path[0] == HeadersLabel)
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
		for _, ref := range idx.GetAllComponentSchemas() {
			restore(ref)
		}
		for _, refs := range [][]*index.Reference{idx.GetAllInlineSchemas(), idx.GetAllInlineSchemaObjects(), idx.GetAllReferenceSchemas()} {
			for _, ref := range refs {
				restore(ref)
			}
		}
	}
	rolodex.GetRootIndex().SetRootNode(root)
	rolodex.SetRootNode(root)
}
