// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"hash/maphash"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
)

// ExternalDoc represents a low-level AsyncAPI 3.0 External Documentation object.
//
// Allows referencing an external resource for extended documentation.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#externalDocumentationObject
type ExternalDoc struct {
	Description low.NodeReference[string]
	URL         low.NodeReference[string]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode     *yaml.Node
	RootNode    *yaml.Node
	idx         *index.SpecIndex
	ctx         context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the ExternalDoc object.
func (e *ExternalDoc) GetRootNode() *yaml.Node {
	return e.RootNode
}

// GetKeyNode returns the key yaml node of the ExternalDoc object.
func (e *ExternalDoc) GetKeyNode() *yaml.Node {
	return e.KeyNode
}

// GetExtensions returns all extensions for ExternalDoc.
func (e *ExternalDoc) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return e.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (e *ExternalDoc) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, e.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the ExternalDoc object.
func (e *ExternalDoc) GetIndex() *index.SpecIndex {
	return e.idx
}

// GetContext returns the context.Context instance used when building the ExternalDoc object.
func (e *ExternalDoc) GetContext() context.Context {
	return e.ctx
}

// Build extracts the ExternalDoc object from the supplied root node.
func (e *ExternalDoc) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	e.KeyNode = keyNode
	root = utils.NodeAlias(root)
	e.RootNode = root
	utils.CheckForMergeNodes(root)
	e.Reference = new(low.Reference)
	e.Nodes = low.ExtractNodes(ctx, root)
	e.Extensions = low.ExtractExtensions(root)
	e.idx = idx
	e.ctx = ctx
	return nil
}

// Hash returns a process-local content hash of the ExternalDoc object.
func (e *ExternalDoc) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !e.Description.IsEmpty() {
			h.WriteString(e.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !e.URL.IsEmpty() {
			h.WriteString(e.URL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(e.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
