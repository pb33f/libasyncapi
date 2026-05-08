// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"crypto/sha256"

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

	// extract servers (array of references)
	_, sLabel, sValue := utils.FindKeyNodeFullTop(ServersLabel, root.Content)
	if sValue != nil && sValue.Kind == yaml.SequenceNode {
		var serverRefs []low.ValueReference[*low.Reference]
		for _, sn := range sValue.Content {
			ref := new(low.Reference)
			ref.SetReference(sn.Value, sn)
			serverRefs = append(serverRefs, low.ValueReference[*low.Reference]{
				Value:     ref,
				ValueNode: sn,
			})
		}
		c.Servers = low.NodeReference[[]low.ValueReference[*low.Reference]]{
			Value:     serverRefs,
			KeyNode:   sLabel,
			ValueNode: sValue,
		}
	}

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
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

	return nil
}

// Hash returns a consistent SHA256 Hash of the Channel object.
func (c *Channel) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)

	if c.Address.Value != nil {
		sb.WriteString(*c.Address.Value)
		sb.WriteByte('|')
	}
	if !c.Title.IsEmpty() {
		sb.WriteString(c.Title.Value)
		sb.WriteByte('|')
	}
	if !c.Summary.IsEmpty() {
		sb.WriteString(c.Summary.Value)
		sb.WriteByte('|')
	}
	if !c.Description.IsEmpty() {
		sb.WriteString(c.Description.Value)
		sb.WriteByte('|')
	}
	if c.Messages.Value != nil {
		for v := range orderedmap.SortAlpha(c.Messages.Value).ValuesFromOldest() {
			sb.WriteString(low.GenerateHashString(v.Value))
			sb.WriteByte('|')
		}
	}
	if c.Parameters.Value != nil {
		for v := range orderedmap.SortAlpha(c.Parameters.Value).ValuesFromOldest() {
			sb.WriteString(low.GenerateHashString(v.Value))
			sb.WriteByte('|')
		}
	}
	if c.Servers.Value != nil {
		for _, srv := range c.Servers.Value {
			if srv.Value != nil {
				sb.WriteString(srv.Value.GetReference())
				sb.WriteByte('|')
			}
		}
	}
	if c.Tags.Value != nil {
		for _, tag := range c.Tags.Value {
			sb.WriteString(low.GenerateHashString(tag.Value))
			sb.WriteByte('|')
		}
	}
	if !c.ExternalDocs.IsEmpty() {
		sb.WriteString(low.GenerateHashString(c.ExternalDocs.Value))
		sb.WriteByte('|')
	}
	if !c.Bindings.IsEmpty() {
		sb.WriteString(low.GenerateHashString(c.Bindings.Value))
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(c.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
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

// Hash returns a consistent SHA256 Hash.
func (cb *ChannelBindings) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !cb.HTTP.IsEmpty() {
		sb.WriteString(low.GenerateHashString(cb.HTTP.Value))
		sb.WriteByte('|')
	}
	if !cb.WebSocket.IsEmpty() {
		sb.WriteString(low.GenerateHashString(cb.WebSocket.Value))
		sb.WriteByte('|')
	}
	if !cb.Kafka.IsEmpty() {
		sb.WriteString(low.GenerateHashString(cb.Kafka.Value))
		sb.WriteByte('|')
	}
	if !cb.AMQP.IsEmpty() {
		sb.WriteString(low.GenerateHashString(cb.AMQP.Value))
		sb.WriteByte('|')
	}
	if !cb.SQS.IsEmpty() {
		sb.WriteString(low.GenerateHashString(cb.SQS.Value))
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(cb.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
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

// Hash returns a consistent SHA256 Hash of the Parameter object.
func (p *Parameter) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)

	if !p.Default.IsEmpty() {
		sb.WriteString(p.Default.Value)
		sb.WriteByte('|')
	}
	if !p.Description.IsEmpty() {
		sb.WriteString(p.Description.Value)
		sb.WriteByte('|')
	}
	if !p.Location.IsEmpty() {
		sb.WriteString(p.Location.Value)
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(p.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}
