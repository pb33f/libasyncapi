// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"errors"
	"fmt"
	"hash/maphash"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"go.yaml.in/yaml/v4"
)

// Channel represents a low-level AsyncAPI 3.0 Channel object.
//
// Describes a shared communication channel.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#channelObject
type Channel struct {
	Address      low.NodeReference[*string] // pointer to allow null
	Messages     low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Message]]]
	Title        low.NodeReference[string]
	Summary      low.NodeReference[string]
	Description  low.NodeReference[string]
	Servers      low.NodeReference[[]low.ValueReference[*low.Reference]] // references only
	Parameters   low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Parameter]]]
	Tags         low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs low.NodeReference[*ExternalDoc]
	Bindings     low.NodeReference[*ChannelBindings]
	Extensions   *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode      *yaml.Node
	RootNode     *yaml.Node
	idx          *index.SpecIndex
	ctx          context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Channel object.
func (c *Channel) GetRootNode() *yaml.Node {
	return c.RootNode
}

// GetKeyNode returns the key yaml node of the Channel object.
func (c *Channel) GetKeyNode() *yaml.Node {
	return c.KeyNode
}

// GetExtensions returns all extensions for Channel.
func (c *Channel) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return c.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (c *Channel) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, c.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Channel object.
func (c *Channel) GetIndex() *index.SpecIndex {
	return c.idx
}

// GetContext returns the context.Context instance used when building the Channel object.
func (c *Channel) GetContext() context.Context {
	return c.ctx
}

// Build extracts the Channel object from the supplied root node.
func (c *Channel) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	c.KeyNode = keyNode
	root = utils.NodeAlias(root)
	c.RootNode = root
	utils.CheckForMergeNodes(root)
	c.Reference = new(low.Reference)
	c.Nodes = low.ExtractNodes(ctx, root)
	c.Extensions = low.ExtractExtensions(root)
	c.idx = idx
	c.ctx = ctx

	// extract address (can be null)
	_, addrLabel, addrValue := utils.FindKeyNodeFullTop(AddressLabel, root.Content)
	if addrValue != nil {
		c.Nodes.Store(addrLabel.Line, addrLabel)
		if addrValue.Tag == "!!null" || addrValue.Value == "null" || addrValue.Value == "~" {
			c.Address = low.NodeReference[*string]{
				Value:     nil,
				KeyNode:   addrLabel,
				ValueNode: addrValue,
			}
		} else {
			addr := addrValue.Value
			c.Address = low.NodeReference[*string]{
				Value:     &addr,
				KeyNode:   addrLabel,
				ValueNode: addrValue,
			}
		}
	}

	// extract messages
	msgs, mLabel, mValue, err := low.ExtractMap[*Message](ctx, MessagesLabel, root, idx)
	if err != nil {
		return err
	}
	if msgs != nil {
		c.Messages = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Message]]]{
			Value:     msgs,
			KeyNode:   mLabel,
			ValueNode: mValue,
		}
	}

	// extract parameters
	params, pLabel, pValue, err := low.ExtractMap[*Parameter](ctx, ParametersLabel, root, idx)
	if err != nil {
		return err
	}
	if params != nil {
		c.Parameters = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*Parameter]]]{
			Value:     params,
			KeyNode:   pLabel,
			ValueNode: pValue,
		}
	}

	// extract servers (array of references - each element is a mapping with $ref)
	// The format is:
	//   servers:
	//     - $ref: '#/servers/...'
	var buildErrs []error
	_, sLabel, sValue := utils.FindKeyNodeFullTop(ServersLabel, root.Content)
	if sValue != nil && sValue.Kind != yaml.SequenceNode {
		buildErrs = append(buildErrs, fmt.Errorf("channel servers must be an array of Reference Objects, line %d, column %d",
			sValue.Line, sValue.Column))
	}
	if sValue != nil && sValue.Kind == yaml.SequenceNode {
		var serverRefs []low.ValueReference[*low.Reference]
		for _, sn := range sValue.Content {
			ref := referenceFromNode(sn)
			if ref == nil {
				buildErrs = append(buildErrs, fmt.Errorf("channel servers entry must be a Reference Object containing a non-empty $ref string, line %d, column %d",
					sn.Line, sn.Column))
				continue
			}
			serverRefs = append(serverRefs, low.ValueReference[*low.Reference]{
				Value:     ref,
				ValueNode: sn,
			})
		}
		if len(serverRefs) > 0 {
			c.Servers = low.NodeReference[[]low.ValueReference[*low.Reference]]{
				Value:     serverRefs,
				KeyNode:   sLabel,
				ValueNode: sValue,
			}
		}
	}

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		buildErrs = append(buildErrs, err)
	}
	if tags != nil {
		c.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	c.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*ChannelBindings](ctx, BindingsLabel, root, idx)
	c.Bindings = bindings

	return errors.Join(buildErrs...)
}

func referenceFromNode(node *yaml.Node) *low.Reference {
	node = utils.NodeAlias(node)
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value != RefLabel {
			continue
		}
		value := utils.NodeAlias(node.Content[i+1])
		if value == nil || value.Kind != yaml.ScalarNode || value.Value == "" {
			return nil
		}
		ref := new(low.Reference)
		ref.SetReference(value.Value, value)
		return ref
	}
	return nil
}

// Hash returns a process-local content hash of the Channel object.
func (c *Channel) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if c.Address.Value != nil {
			h.WriteString(*c.Address.Value)
			h.WriteByte(low.HASH_PIPE)
		} else if c.Address.ValueNode != nil {
			// an explicit `address: null` (address unknown) is semantically distinct
			// from an absent address, so it must hash differently.
			h.WriteString("\x00null")
			h.WriteByte(low.HASH_PIPE)
		}
		if !c.Title.IsEmpty() {
			h.WriteString(c.Title.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !c.Summary.IsEmpty() {
			h.WriteString(c.Summary.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !c.Description.IsEmpty() {
			h.WriteString(c.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if c.Messages.Value != nil {
			for v := range orderedmap.SortAlpha(c.Messages.Value).ValuesFromOldest() {
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
		if c.Servers.Value != nil {
			for _, srv := range c.Servers.Value {
				if srv.Value != nil {
					h.WriteString(srv.Value.GetReference())
					h.WriteByte(low.HASH_PIPE)
				}
			}
		}
		if c.Tags.Value != nil {
			for _, tag := range c.Tags.Value {
				h.WriteString(low.GenerateHashString(tag.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !c.ExternalDocs.IsEmpty() {
			h.WriteString(low.GenerateHashString(c.ExternalDocs.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !c.Bindings.IsEmpty() {
			h.WriteString(low.GenerateHashString(c.Bindings.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(c.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// ChannelBindings represents a low-level AsyncAPI 3.0 Channel Bindings object.
//
// Map of channel binding objects where the keys are protocol names.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#channelBindingsObject
type ChannelBindings struct {
	HTTP       low.NodeReference[*HTTPChannelBinding]
	WebSocket  low.NodeReference[*WebSocketChannelBinding]
	Kafka      low.NodeReference[*KafkaChannelBinding]
	AMQP       low.NodeReference[*AMQPChannelBinding]
	SQS        low.NodeReference[*SQSChannelBinding]
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (cb *ChannelBindings) GetRootNode() *yaml.Node {
	return cb.RootNode
}

// GetKeyNode returns the key yaml node.
func (cb *ChannelBindings) GetKeyNode() *yaml.Node {
	return cb.KeyNode
}

// GetExtensions returns all extensions.
func (cb *ChannelBindings) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return cb.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (cb *ChannelBindings) GetIndex() *index.SpecIndex {
	return cb.idx
}

// GetContext returns the context.Context instance.
func (cb *ChannelBindings) GetContext() context.Context {
	return cb.ctx
}

// Build extracts the ChannelBindings object from the supplied root node.
func (cb *ChannelBindings) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	cb.KeyNode = keyNode
	root = utils.NodeAlias(root)
	cb.RootNode = root
	utils.CheckForMergeNodes(root)
	cb.Reference = new(low.Reference)
	cb.Nodes = low.ExtractNodes(ctx, root)
	cb.Extensions = low.ExtractExtensions(root)
	cb.idx = idx
	cb.ctx = ctx

	http, _ := low.ExtractObject[*HTTPChannelBinding](ctx, HTTPLabel, root, idx)
	cb.HTTP = http

	ws, _ := low.ExtractObject[*WebSocketChannelBinding](ctx, WebSocketLabel, root, idx)
	cb.WebSocket = ws

	kafka, _ := low.ExtractObject[*KafkaChannelBinding](ctx, KafkaLabel, root, idx)
	cb.Kafka = kafka

	amqp, _ := low.ExtractObject[*AMQPChannelBinding](ctx, AMQPLabel, root, idx)
	cb.AMQP = amqp

	sqs, _ := low.ExtractObject[*SQSChannelBinding](ctx, SQSLabel, root, idx)
	cb.SQS = sqs

	return nil
}

// Hash returns a process-local content hash.
func (cb *ChannelBindings) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !cb.HTTP.IsEmpty() {
			h.WriteString(low.GenerateHashString(cb.HTTP.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !cb.WebSocket.IsEmpty() {
			h.WriteString(low.GenerateHashString(cb.WebSocket.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !cb.Kafka.IsEmpty() {
			h.WriteString(low.GenerateHashString(cb.Kafka.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !cb.AMQP.IsEmpty() {
			h.WriteString(low.GenerateHashString(cb.AMQP.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !cb.SQS.IsEmpty() {
			h.WriteString(low.GenerateHashString(cb.SQS.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(cb.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// Parameter represents a low-level AsyncAPI 3.0 Parameter object.
//
// Describes a parameter included in a channel address.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#parameterObject
type Parameter struct {
	Enum        low.NodeReference[[]low.ValueReference[string]]
	Default     low.NodeReference[string]
	Description low.NodeReference[string]
	Examples    low.NodeReference[[]low.ValueReference[string]]
	Location    low.NodeReference[string]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode     *yaml.Node
	RootNode    *yaml.Node
	idx         *index.SpecIndex
	ctx         context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Parameter object.
func (p *Parameter) GetRootNode() *yaml.Node {
	return p.RootNode
}

// GetKeyNode returns the key yaml node of the Parameter object.
func (p *Parameter) GetKeyNode() *yaml.Node {
	return p.KeyNode
}

// GetExtensions returns all extensions for Parameter.
func (p *Parameter) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return p.Extensions
}

// GetIndex returns the index.SpecIndex instance attached to the Parameter object.
func (p *Parameter) GetIndex() *index.SpecIndex {
	return p.idx
}

// GetContext returns the context.Context instance used when building the Parameter object.
func (p *Parameter) GetContext() context.Context {
	return p.ctx
}

// Build extracts the Parameter object from the supplied root node.
func (p *Parameter) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	p.KeyNode = keyNode
	root = utils.NodeAlias(root)
	p.RootNode = root
	utils.CheckForMergeNodes(root)
	p.Reference = new(low.Reference)
	p.Nodes = low.ExtractNodes(ctx, root)
	p.Extensions = low.ExtractExtensions(root)
	p.idx = idx
	p.ctx = ctx
	return nil
}

// Hash returns a process-local content hash of the Parameter object.
func (p *Parameter) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if p.Enum.Value != nil {
			for _, v := range p.Enum.Value {
				h.WriteString(v.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !p.Default.IsEmpty() {
			h.WriteString(p.Default.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !p.Description.IsEmpty() {
			h.WriteString(p.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if p.Examples.Value != nil {
			for _, v := range p.Examples.Value {
				h.WriteString(v.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !p.Location.IsEmpty() {
			h.WriteString(p.Location.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(p.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
