// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"errors"
	"hash/maphash"
	"sync"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
)

// Components represents a low-level AsyncAPI 3.0 Components object.
//
// Holds a set of reusable objects for different aspects of the AsyncAPI specification.
// All objects defined within the components object will have no effect on the API
// unless they are explicitly referenced from properties outside the components object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#componentsObject
type Components struct {
	Schemas           low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*base.SchemaProxy]]]
	Servers           low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Server]]]
	Channels          low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Channel]]]
	Operations        low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Operation]]]
	Messages          low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Message]]]
	SecuritySchemes   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*SecurityScheme]]]
	ServerVariables   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerVariable]]]
	Parameters        low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Parameter]]]
	CorrelationIDs    low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*CorrelationID]]]
	Replies           low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationReply]]]
	ReplyAddresses    low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationReplyAddress]]]
	ExternalDocs      low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ExternalDoc]]]
	Tags              low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Tag]]]
	OperationTraits   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationTrait]]]
	MessageTraits     low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MessageTrait]]]
	ServerBindings    low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerBindings]]]
	ChannelBindings   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ChannelBindings]]]
	OperationBindings low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationBindings]]]
	MessageBindings   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MessageBindings]]]
	Extensions        *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode           *yaml.Node
	RootNode          *yaml.Node
	idx               *index.SpecIndex
	ctx               context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Components object.
func (c *Components) GetRootNode() *yaml.Node {
	return c.RootNode
}

// GetKeyNode returns the key yaml node of the Components object.
func (c *Components) GetKeyNode() *yaml.Node {
	return c.KeyNode
}

// GetExtensions returns all extensions for Components.
func (c *Components) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return c.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (c *Components) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, c.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Components object.
func (c *Components) GetIndex() *index.SpecIndex {
	return c.idx
}

// GetContext returns the context.Context instance used when building the Components object.
func (c *Components) GetContext() context.Context {
	return c.ctx
}

// FindSchema attempts to locate a SchemaProxy from 'schemas' with a specific name.
func (c *Components) FindSchema(name string) *low.ValueReference[*base.SchemaProxy] {
	return low.FindItemInOrderedMap[*base.SchemaProxy](name, c.Schemas.Value)
}

// FindServer attempts to locate a Server from 'servers' with a specific name.
func (c *Components) FindServer(name string) *low.ValueReference[*Server] {
	return low.FindItemInOrderedMap[*Server](name, c.Servers.Value)
}

// FindChannel attempts to locate a Channel from 'channels' with a specific name.
func (c *Components) FindChannel(name string) *low.ValueReference[*Channel] {
	return low.FindItemInOrderedMap[*Channel](name, c.Channels.Value)
}

// FindOperation attempts to locate an Operation from 'operations' with a specific name.
func (c *Components) FindOperation(name string) *low.ValueReference[*Operation] {
	return low.FindItemInOrderedMap[*Operation](name, c.Operations.Value)
}

// FindMessage attempts to locate a Message from 'messages' with a specific name.
func (c *Components) FindMessage(name string) *low.ValueReference[*Message] {
	return low.FindItemInOrderedMap[*Message](name, c.Messages.Value)
}

// FindSecurityScheme attempts to locate a SecurityScheme from 'securitySchemes' with a specific name.
func (c *Components) FindSecurityScheme(name string) *low.ValueReference[*SecurityScheme] {
	return low.FindItemInOrderedMap[*SecurityScheme](name, c.SecuritySchemes.Value)
}

// Build extracts the Components object from the supplied root node.
func (c *Components) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	c.KeyNode = keyNode
	root = utils.NodeAlias(root)
	c.RootNode = root
	utils.CheckForMergeNodes(root)
	c.Reference = new(low.Reference)
	c.Nodes = low.ExtractNodes(ctx, root)
	c.Extensions = low.ExtractExtensions(root)
	c.idx = idx
	c.ctx = ctx

	var wg sync.WaitGroup
	var errMu sync.Mutex
	var buildErrs []error

	captureErr := func(err error) {
		if err != nil {
			errMu.Lock()
			buildErrs = append(buildErrs, err)
			errMu.Unlock()
		}
	}

	// extract all 19 component types in parallel
	wg.Add(19)

	go func() {
		defer wg.Done()
		schemas, sLabel, sValue, err := low.ExtractMap[*base.SchemaProxy](ctx, SchemasLabel, root, idx)
		captureErr(err)
		if schemas != nil {
			c.Schemas = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*base.SchemaProxy]]]{
				Value: schemas, KeyNode: sLabel, ValueNode: sValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		servers, sLabel, sValue, err := low.ExtractMap[*Server](ctx, ServersLabel, root, idx)
		captureErr(err)
		if servers != nil {
			c.Servers = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Server]]]{
				Value: servers, KeyNode: sLabel, ValueNode: sValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		channels, cLabel, cValue, err := low.ExtractMap[*Channel](ctx, ChannelsLabel, root, idx)
		captureErr(err)
		if channels != nil {
			c.Channels = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Channel]]]{
				Value: channels, KeyNode: cLabel, ValueNode: cValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		ops, oLabel, oValue, err := low.ExtractMap[*Operation](ctx, OperationsLabel, root, idx)
		captureErr(err)
		if ops != nil {
			c.Operations = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Operation]]]{
				Value: ops, KeyNode: oLabel, ValueNode: oValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		msgs, mLabel, mValue, err := low.ExtractMap[*Message](ctx, MessagesLabel, root, idx)
		captureErr(err)
		if msgs != nil {
			c.Messages = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Message]]]{
				Value: msgs, KeyNode: mLabel, ValueNode: mValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		ss, ssLabel, ssValue, err := low.ExtractMap[*SecurityScheme](ctx, SecuritySchemesLabel, root, idx)
		captureErr(err)
		if ss != nil {
			c.SecuritySchemes = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*SecurityScheme]]]{
				Value: ss, KeyNode: ssLabel, ValueNode: ssValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		sv, svLabel, svValue, err := low.ExtractMap[*ServerVariable](ctx, ServerVariablesLabel, root, idx)
		captureErr(err)
		if sv != nil {
			c.ServerVariables = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerVariable]]]{
				Value: sv, KeyNode: svLabel, ValueNode: svValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		params, pLabel, pValue, err := low.ExtractMap[*Parameter](ctx, ParametersLabel, root, idx)
		captureErr(err)
		if params != nil {
			c.Parameters = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Parameter]]]{
				Value: params, KeyNode: pLabel, ValueNode: pValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		cids, cLabel, cValue, err := low.ExtractMap[*CorrelationID](ctx, CorrelationIDsLabel, root, idx)
		captureErr(err)
		if cids != nil {
			c.CorrelationIDs = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*CorrelationID]]]{
				Value: cids, KeyNode: cLabel, ValueNode: cValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		reps, rLabel, rValue, err := low.ExtractMap[*OperationReply](ctx, RepliesLabel, root, idx)
		captureErr(err)
		if reps != nil {
			c.Replies = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationReply]]]{
				Value: reps, KeyNode: rLabel, ValueNode: rValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		ras, raLabel, raValue, err := low.ExtractMap[*OperationReplyAddress](ctx, ReplyAddressesLabel, root, idx)
		captureErr(err)
		if ras != nil {
			c.ReplyAddresses = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationReplyAddress]]]{
				Value: ras, KeyNode: raLabel, ValueNode: raValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		eds, edLabel, edValue, err := low.ExtractMap[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
		captureErr(err)
		if eds != nil {
			c.ExternalDocs = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ExternalDoc]]]{
				Value: eds, KeyNode: edLabel, ValueNode: edValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		tags, tLabel, tValue, err := low.ExtractMap[*Tag](ctx, TagsLabel, root, idx)
		captureErr(err)
		if tags != nil {
			c.Tags = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Tag]]]{
				Value: tags, KeyNode: tLabel, ValueNode: tValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		ots, otLabel, otValue, err := low.ExtractMap[*OperationTrait](ctx, OperationTraitsLabel, root, idx)
		captureErr(err)
		if ots != nil {
			c.OperationTraits = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationTrait]]]{
				Value: ots, KeyNode: otLabel, ValueNode: otValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		mts, mtLabel, mtValue, err := low.ExtractMap[*MessageTrait](ctx, MessageTraitsLabel, root, idx)
		captureErr(err)
		if mts != nil {
			c.MessageTraits = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MessageTrait]]]{
				Value: mts, KeyNode: mtLabel, ValueNode: mtValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		sbs, sbLabel, sbValue, err := low.ExtractMap[*ServerBindings](ctx, ServerBindingsLabel, root, idx)
		captureErr(err)
		if sbs != nil {
			c.ServerBindings = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerBindings]]]{
				Value: sbs, KeyNode: sbLabel, ValueNode: sbValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		cbs, cbLabel, cbValue, err := low.ExtractMap[*ChannelBindings](ctx, ChannelBindingsLabel, root, idx)
		captureErr(err)
		if cbs != nil {
			c.ChannelBindings = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ChannelBindings]]]{
				Value: cbs, KeyNode: cbLabel, ValueNode: cbValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		obs, obLabel, obValue, err := low.ExtractMap[*OperationBindings](ctx, OperationBindingsLabel, root, idx)
		captureErr(err)
		if obs != nil {
			c.OperationBindings = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*OperationBindings]]]{
				Value: obs, KeyNode: obLabel, ValueNode: obValue,
			}
		}
	}()

	go func() {
		defer wg.Done()
		mbs, mbLabel, mbValue, err := low.ExtractMap[*MessageBindings](ctx, MessageBindingsLabel, root, idx)
		captureErr(err)
		if mbs != nil {
			c.MessageBindings = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*MessageBindings]]]{
				Value: mbs, KeyNode: mbLabel, ValueNode: mbValue,
			}
		}
	}()

	wg.Wait()
	return errors.Join(buildErrs...)
}

// Hash returns a process-local content hash of the Components object.
func (c *Components) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if c.Schemas.Value != nil {
			for v := range orderedmap.SortAlpha(c.Schemas.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Servers.Value != nil {
			for v := range orderedmap.SortAlpha(c.Servers.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Channels.Value != nil {
			for v := range orderedmap.SortAlpha(c.Channels.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Operations.Value != nil {
			for v := range orderedmap.SortAlpha(c.Operations.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Messages.Value != nil {
			for v := range orderedmap.SortAlpha(c.Messages.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.SecuritySchemes.Value != nil {
			for v := range orderedmap.SortAlpha(c.SecuritySchemes.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.ServerVariables.Value != nil {
			for v := range orderedmap.SortAlpha(c.ServerVariables.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Parameters.Value != nil {
			for v := range orderedmap.SortAlpha(c.Parameters.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.CorrelationIDs.Value != nil {
			for v := range orderedmap.SortAlpha(c.CorrelationIDs.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Replies.Value != nil {
			for v := range orderedmap.SortAlpha(c.Replies.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.ReplyAddresses.Value != nil {
			for v := range orderedmap.SortAlpha(c.ReplyAddresses.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.ExternalDocs.Value != nil {
			for v := range orderedmap.SortAlpha(c.ExternalDocs.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.Tags.Value != nil {
			for v := range orderedmap.SortAlpha(c.Tags.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.OperationTraits.Value != nil {
			for v := range orderedmap.SortAlpha(c.OperationTraits.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.MessageTraits.Value != nil {
			for v := range orderedmap.SortAlpha(c.MessageTraits.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.ServerBindings.Value != nil {
			for v := range orderedmap.SortAlpha(c.ServerBindings.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.ChannelBindings.Value != nil {
			for v := range orderedmap.SortAlpha(c.ChannelBindings.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.OperationBindings.Value != nil {
			for v := range orderedmap.SortAlpha(c.OperationBindings.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if c.MessageBindings.Value != nil {
			for v := range orderedmap.SortAlpha(c.MessageBindings.Value).ValuesFromOldest() {
				h.WriteString(low.GenerateHashString(v.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		for _, ext := range low.HashExtensions(c.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
