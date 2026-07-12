// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"hash/maphash"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"go.yaml.in/yaml/v4"
)

// Tag represents a low-level AsyncAPI 3.0 Tag object.
//
// Allows adding metadata to a single tag.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#tagObject
type Tag struct {
	Name         low.NodeReference[string]
	Description  low.NodeReference[string]
	ExternalDocs low.NodeReference[*ExternalDoc]
	Extensions   *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode      *yaml.Node
	RootNode     *yaml.Node
	idx          *index.SpecIndex
	ctx          context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Tag object.
func (t *Tag) GetRootNode() *yaml.Node {
	return t.RootNode
}

// GetKeyNode returns the key yaml node of the Tag object.
func (t *Tag) GetKeyNode() *yaml.Node {
	return t.KeyNode
}

// GetExtensions returns all extensions for Tag.
func (t *Tag) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return t.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (t *Tag) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, t.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Tag object.
func (t *Tag) GetIndex() *index.SpecIndex {
	return t.idx
}

// GetContext returns the context.Context instance used when building the Tag object.
func (t *Tag) GetContext() context.Context {
	return t.ctx
}

// Build extracts the Tag object from the supplied root node.
func (t *Tag) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	t.KeyNode = keyNode
	root = utils.NodeAlias(root)
	t.RootNode = root
	utils.CheckForMergeNodes(root)
	t.Reference = new(low.Reference)
	t.Nodes = low.ExtractNodes(ctx, root)
	t.Extensions = low.ExtractExtensions(root)
	t.idx = idx
	t.ctx = ctx

	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	t.ExternalDocs = extDocs

	return nil
}

// Hash returns a process-local content hash of the Tag object.
func (t *Tag) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !t.Name.IsEmpty() {
			h.WriteString(t.Name.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !t.Description.IsEmpty() {
			h.WriteString(t.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !t.ExternalDocs.IsEmpty() {
			h.WriteString(low.GenerateHashString(t.ExternalDocs.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(t.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
