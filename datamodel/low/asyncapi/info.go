// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"hash/maphash"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
)

// Info represents a low-level AsyncAPI 3.0 Info object.
//
// The Info object provides metadata about the API. The metadata can be used by clients if needed,
// and can be presented in editing or documentation generation tools.
//
// Unlike OpenAPI Info, AsyncAPI Info includes tags and externalDocs fields, and these can
// contain references ($ref).
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#infoObject
type Info struct {
	Title          low.NodeReference[string]
	Version        low.NodeReference[string]
	Description    low.NodeReference[string]
	TermsOfService low.NodeReference[string]
	Contact        low.NodeReference[*base.Contact]
	License        low.NodeReference[*base.License]
	Tags           low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Tag]]]
	ExternalDocs   low.NodeReference[*ExternalDoc]
	Extensions     *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode        *yaml.Node
	RootNode       *yaml.Node
	idx            *index.SpecIndex
	ctx            context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Info object.
func (i *Info) GetRootNode() *yaml.Node {
	return i.RootNode
}

// GetKeyNode returns the key yaml node of the Info object.
func (i *Info) GetKeyNode() *yaml.Node {
	return i.KeyNode
}

// GetExtensions returns all extensions for Info.
func (i *Info) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return i.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (i *Info) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, i.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Info object.
func (i *Info) GetIndex() *index.SpecIndex {
	return i.idx
}

// GetContext returns the context.Context instance used when building the Info object.
func (i *Info) GetContext() context.Context {
	return i.ctx
}

// Build extracts the Info object from the supplied root node.
func (i *Info) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	i.KeyNode = keyNode
	root = utils.NodeAlias(root)
	i.RootNode = root
	utils.CheckForMergeNodes(root)
	i.Reference = new(low.Reference)
	i.Nodes = low.ExtractNodes(ctx, root)
	i.Extensions = low.ExtractExtensions(root)
	i.idx = idx
	i.ctx = ctx

	contact, _ := low.ExtractObject[*base.Contact](ctx, ContactLabel, root, idx)
	i.Contact = contact

	lic, _ := low.ExtractObject[*base.License](ctx, LicenseLabel, root, idx)
	i.License = lic

	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	i.ExternalDocs = extDocs

	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
	}
	if tags != nil {
		tagMap := orderedmap.New[low.KeyReference[string], low.ValueReference[*Tag]]()
		for _, tag := range tags {
			key := low.KeyReference[string]{
				Value:   tag.Value.Name.Value,
				KeyNode: tag.ValueNode,
			}
			tagMap.Set(key, tag)
		}
		i.Tags = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Tag]]]{
			Value:     tagMap,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	return nil
}

// Hash returns a process-local content hash of the Info object.
func (i *Info) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !i.Title.IsEmpty() {
			h.WriteString(i.Title.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.Version.IsEmpty() {
			h.WriteString(i.Version.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.Description.IsEmpty() {
			h.WriteString(i.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.TermsOfService.IsEmpty() {
			h.WriteString(i.TermsOfService.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.Contact.IsEmpty() {
			h.WriteString(low.GenerateHashString(i.Contact.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.License.IsEmpty() {
			h.WriteString(low.GenerateHashString(i.License.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !i.ExternalDocs.IsEmpty() {
			h.WriteString(low.GenerateHashString(i.ExternalDocs.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if i.Tags.Value != nil {
			for v := range orderedmap.SortAlpha(i.Tags.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		for _, ext := range low.HashExtensions(i.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
