// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"crypto/sha256"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"go.yaml.in/yaml/v4"
)

// Message represents a low-level AsyncAPI 3.0 Message object.
//
// Describes a message received on a given channel and operation.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageObject
type Message struct {
	Headers       low.NodeReference[*base.SchemaProxy]
	Payload       low.NodeReference[*base.SchemaProxy]
	CorrelationID low.NodeReference[*CorrelationID]
	ContentType   low.NodeReference[string]
	Name          low.NodeReference[string]
	Title         low.NodeReference[string]
	Summary       low.NodeReference[string]
	Description   low.NodeReference[string]
	Tags          low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs  low.NodeReference[*ExternalDoc]
	Bindings      low.NodeReference[*MessageBindings]
	Examples      low.NodeReference[[]low.ValueReference[*MessageExample]]
	Traits        low.NodeReference[[]low.ValueReference[*MessageTrait]]
	Extensions    *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode       *yaml.Node
	RootNode      *yaml.Node
	idx           *index.SpecIndex
	ctx           context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Message object.
func (m *Message) GetRootNode() *yaml.Node {
	return m.RootNode
}

// GetKeyNode returns the key yaml node of the Message object.
func (m *Message) GetKeyNode() *yaml.Node {
	return m.KeyNode
}

// GetExtensions returns all extensions for Message.
func (m *Message) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return m.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (m *Message) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, m.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Message object.
func (m *Message) GetIndex() *index.SpecIndex {
	return m.idx
}

// GetContext returns the context.Context instance used when building the Message object.
func (m *Message) GetContext() context.Context {
	return m.ctx
}

// Build extracts the Message object from the supplied root node.
func (m *Message) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	m.KeyNode = keyNode
	root = utils.NodeAlias(root)
	m.RootNode = root
	utils.CheckForMergeNodes(root)
	m.Reference = new(low.Reference)
	m.Nodes = low.ExtractNodes(ctx, root)
	m.Extensions = low.ExtractExtensions(root)
	m.idx = idx
	m.ctx = ctx

	// extract headers schema
	headers, _ := low.ExtractObject[*base.SchemaProxy](ctx, HeadersLabel, root, idx)
	m.Headers = headers

	// extract payload schema
	payload, _ := low.ExtractObject[*base.SchemaProxy](ctx, PayloadLabel, root, idx)
	m.Payload = payload

	// extract correlationId
	corrId, _ := low.ExtractObject[*CorrelationID](ctx, CorrelationIDLabel, root, idx)
	m.CorrelationID = corrId

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
	}
	if tags != nil {
		m.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	m.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*MessageBindings](ctx, BindingsLabel, root, idx)
	m.Bindings = bindings

	// extract examples
	examples, exLabel, exValue, err := low.ExtractArray[*MessageExample](ctx, ExamplesLabel, root, idx)
	if err != nil {
		return err
	}
	if examples != nil {
		m.Examples = low.NodeReference[[]low.ValueReference[*MessageExample]]{
			Value:     examples,
			KeyNode:   exLabel,
			ValueNode: exValue,
		}
	}

	// extract traits
	traits, trLabel, trValue, err := low.ExtractArray[*MessageTrait](ctx, TraitsLabel, root, idx)
	if err != nil {
		return err
	}
	if traits != nil {
		m.Traits = low.NodeReference[[]low.ValueReference[*MessageTrait]]{
			Value:     traits,
			KeyNode:   trLabel,
			ValueNode: trValue,
		}
	}

	return nil
}

// Hash returns a consistent SHA256 Hash of the Message object.
func (m *Message) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)

	if !m.Headers.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.Headers.Value))
		sb.WriteByte('|')
	}
	if !m.Payload.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.Payload.Value))
		sb.WriteByte('|')
	}
	if !m.CorrelationID.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.CorrelationID.Value))
		sb.WriteByte('|')
	}
	if !m.ContentType.IsEmpty() {
		sb.WriteString(m.ContentType.Value)
		sb.WriteByte('|')
	}
	if !m.Name.IsEmpty() {
		sb.WriteString(m.Name.Value)
		sb.WriteByte('|')
	}
	if !m.Title.IsEmpty() {
		sb.WriteString(m.Title.Value)
		sb.WriteByte('|')
	}
	if !m.Summary.IsEmpty() {
		sb.WriteString(m.Summary.Value)
		sb.WriteByte('|')
	}
	if !m.Description.IsEmpty() {
		sb.WriteString(m.Description.Value)
		sb.WriteByte('|')
	}
	if m.Tags.Value != nil {
		for _, tag := range m.Tags.Value {
			sb.WriteString(low.GenerateHashString(tag.Value))
			sb.WriteByte('|')
		}
	}
	if !m.ExternalDocs.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.ExternalDocs.Value))
		sb.WriteByte('|')
	}
	if !m.Bindings.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.Bindings.Value))
		sb.WriteByte('|')
	}
	if m.Examples.Value != nil {
		for _, ex := range m.Examples.Value {
			sb.WriteString(low.GenerateHashString(ex.Value))
			sb.WriteByte('|')
		}
	}
	if m.Traits.Value != nil {
		for _, trait := range m.Traits.Value {
			sb.WriteString(low.GenerateHashString(trait.Value))
			sb.WriteByte('|')
		}
	}
	for _, ext := range low.HashExtensions(m.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// MessageBindings represents a low-level AsyncAPI 3.0 Message Bindings object.
//
// Map of message binding objects where the keys are protocol names.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageBindingsObject
type MessageBindings struct {
	HTTP       low.NodeReference[*HTTPMessageBinding]
	Kafka      low.NodeReference[*KafkaMessageBinding]
	AMQP       low.NodeReference[*AMQPMessageBinding]
	MQTT       low.NodeReference[*MQTTMessageBinding]
	SQS        low.NodeReference[*SQSMessageBinding]
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (mb *MessageBindings) GetRootNode() *yaml.Node {
	return mb.RootNode
}

// GetKeyNode returns the key yaml node.
func (mb *MessageBindings) GetKeyNode() *yaml.Node {
	return mb.KeyNode
}

// GetExtensions returns all extensions.
func (mb *MessageBindings) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return mb.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (mb *MessageBindings) GetIndex() *index.SpecIndex {
	return mb.idx
}

// GetContext returns the context.Context instance.
func (mb *MessageBindings) GetContext() context.Context {
	return mb.ctx
}

// Build extracts the MessageBindings object from the supplied root node.
func (mb *MessageBindings) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	mb.KeyNode = keyNode
	root = utils.NodeAlias(root)
	mb.RootNode = root
	utils.CheckForMergeNodes(root)
	mb.Reference = new(low.Reference)
	mb.Nodes = low.ExtractNodes(ctx, root)
	mb.Extensions = low.ExtractExtensions(root)
	mb.idx = idx
	mb.ctx = ctx

	http, _ := low.ExtractObject[*HTTPMessageBinding](ctx, HTTPLabel, root, idx)
	mb.HTTP = http

	kafka, _ := low.ExtractObject[*KafkaMessageBinding](ctx, KafkaLabel, root, idx)
	mb.Kafka = kafka

	amqp, _ := low.ExtractObject[*AMQPMessageBinding](ctx, AMQPLabel, root, idx)
	mb.AMQP = amqp

	mqtt, _ := low.ExtractObject[*MQTTMessageBinding](ctx, MQTTLabel, root, idx)
	mb.MQTT = mqtt

	sqs, _ := low.ExtractObject[*SQSMessageBinding](ctx, SQSLabel, root, idx)
	mb.SQS = sqs

	return nil
}

// Hash returns a consistent SHA256 Hash.
func (mb *MessageBindings) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !mb.HTTP.IsEmpty() {
		sb.WriteString(low.GenerateHashString(mb.HTTP.Value))
		sb.WriteByte('|')
	}
	if !mb.Kafka.IsEmpty() {
		sb.WriteString(low.GenerateHashString(mb.Kafka.Value))
		sb.WriteByte('|')
	}
	if !mb.AMQP.IsEmpty() {
		sb.WriteString(low.GenerateHashString(mb.AMQP.Value))
		sb.WriteByte('|')
	}
	if !mb.MQTT.IsEmpty() {
		sb.WriteString(low.GenerateHashString(mb.MQTT.Value))
		sb.WriteByte('|')
	}
	if !mb.SQS.IsEmpty() {
		sb.WriteString(low.GenerateHashString(mb.SQS.Value))
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(mb.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// MessageExample represents a low-level AsyncAPI 3.0 Message Example object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageExampleObject
type MessageExample struct {
	Headers    low.NodeReference[*yaml.Node]
	Payload    low.NodeReference[*yaml.Node]
	Name       low.NodeReference[string]
	Summary    low.NodeReference[string]
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (me *MessageExample) GetRootNode() *yaml.Node {
	return me.RootNode
}

// GetKeyNode returns the key yaml node.
func (me *MessageExample) GetKeyNode() *yaml.Node {
	return me.KeyNode
}

// GetExtensions returns all extensions.
func (me *MessageExample) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return me.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (me *MessageExample) GetIndex() *index.SpecIndex {
	return me.idx
}

// GetContext returns the context.Context instance.
func (me *MessageExample) GetContext() context.Context {
	return me.ctx
}

// Build extracts the MessageExample object from the supplied root node.
func (me *MessageExample) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	me.KeyNode = keyNode
	root = utils.NodeAlias(root)
	me.RootNode = root
	utils.CheckForMergeNodes(root)
	me.Reference = new(low.Reference)
	me.Nodes = low.ExtractNodes(ctx, root)
	me.Extensions = low.ExtractExtensions(root)
	me.idx = idx
	me.ctx = ctx
	return nil
}

// Hash returns a consistent SHA256 Hash.
func (me *MessageExample) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !me.Name.IsEmpty() {
		sb.WriteString(me.Name.Value)
		sb.WriteByte('|')
	}
	if !me.Summary.IsEmpty() {
		sb.WriteString(me.Summary.Value)
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(me.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// MessageTrait represents a low-level AsyncAPI 3.0 Message Trait object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageTraitObject
type MessageTrait struct {
	Headers       low.NodeReference[*base.SchemaProxy]
	CorrelationID low.NodeReference[*CorrelationID]
	ContentType   low.NodeReference[string]
	Name          low.NodeReference[string]
	Title         low.NodeReference[string]
	Summary       low.NodeReference[string]
	Description   low.NodeReference[string]
	Tags          low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs  low.NodeReference[*ExternalDoc]
	Bindings      low.NodeReference[*MessageBindings]
	Examples      low.NodeReference[[]low.ValueReference[*MessageExample]]
	Extensions    *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode       *yaml.Node
	RootNode      *yaml.Node
	idx           *index.SpecIndex
	ctx           context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (mt *MessageTrait) GetRootNode() *yaml.Node {
	return mt.RootNode
}

// GetKeyNode returns the key yaml node.
func (mt *MessageTrait) GetKeyNode() *yaml.Node {
	return mt.KeyNode
}

// GetExtensions returns all extensions.
func (mt *MessageTrait) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return mt.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (mt *MessageTrait) GetIndex() *index.SpecIndex {
	return mt.idx
}

// GetContext returns the context.Context instance.
func (mt *MessageTrait) GetContext() context.Context {
	return mt.ctx
}

// Build extracts the MessageTrait object from the supplied root node.
func (mt *MessageTrait) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	mt.KeyNode = keyNode
	root = utils.NodeAlias(root)
	mt.RootNode = root
	utils.CheckForMergeNodes(root)
	mt.Reference = new(low.Reference)
	mt.Nodes = low.ExtractNodes(ctx, root)
	mt.Extensions = low.ExtractExtensions(root)
	mt.idx = idx
	mt.ctx = ctx

	// extract headers schema
	headers, _ := low.ExtractObject[*base.SchemaProxy](ctx, HeadersLabel, root, idx)
	mt.Headers = headers

	// extract correlationId
	corrId, _ := low.ExtractObject[*CorrelationID](ctx, CorrelationIDLabel, root, idx)
	mt.CorrelationID = corrId

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
	}
	if tags != nil {
		mt.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	mt.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*MessageBindings](ctx, BindingsLabel, root, idx)
	mt.Bindings = bindings

	// extract examples
	examples, exLabel, exValue, err := low.ExtractArray[*MessageExample](ctx, ExamplesLabel, root, idx)
	if err != nil {
		return err
	}
	if examples != nil {
		mt.Examples = low.NodeReference[[]low.ValueReference[*MessageExample]]{
			Value:     examples,
			KeyNode:   exLabel,
			ValueNode: exValue,
		}
	}

	return nil
}

// Hash returns a consistent SHA256 Hash.
func (mt *MessageTrait) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !mt.ContentType.IsEmpty() {
		sb.WriteString(mt.ContentType.Value)
		sb.WriteByte('|')
	}
	if !mt.Name.IsEmpty() {
		sb.WriteString(mt.Name.Value)
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(mt.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// CorrelationID represents a low-level AsyncAPI 3.0 Correlation ID object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
type CorrelationID struct {
	Description low.NodeReference[string]
	Location    low.NodeReference[string]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode     *yaml.Node
	RootNode    *yaml.Node
	idx         *index.SpecIndex
	ctx         context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (ci *CorrelationID) GetRootNode() *yaml.Node {
	return ci.RootNode
}

// GetKeyNode returns the key yaml node.
func (ci *CorrelationID) GetKeyNode() *yaml.Node {
	return ci.KeyNode
}

// GetExtensions returns all extensions.
func (ci *CorrelationID) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return ci.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (ci *CorrelationID) GetIndex() *index.SpecIndex {
	return ci.idx
}

// GetContext returns the context.Context instance.
func (ci *CorrelationID) GetContext() context.Context {
	return ci.ctx
}

// Build extracts the CorrelationID object from the supplied root node.
func (ci *CorrelationID) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	ci.KeyNode = keyNode
	root = utils.NodeAlias(root)
	ci.RootNode = root
	utils.CheckForMergeNodes(root)
	ci.Reference = new(low.Reference)
	ci.Nodes = low.ExtractNodes(ctx, root)
	ci.Extensions = low.ExtractExtensions(root)
	ci.idx = idx
	ci.ctx = ctx
	return nil
}

// Hash returns a consistent SHA256 Hash.
func (ci *CorrelationID) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !ci.Description.IsEmpty() {
		sb.WriteString(ci.Description.Value)
		sb.WriteByte('|')
	}
	if !ci.Location.IsEmpty() {
		sb.WriteString(ci.Location.Value)
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(ci.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}
