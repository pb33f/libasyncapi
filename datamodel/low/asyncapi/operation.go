// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"errors"
	"fmt"
	"hash/maphash"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
)

// Operation represents a low-level AsyncAPI 3.0 Operation object.
//
// Describes a specific operation.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationObject
type Operation struct {
	Action       low.NodeReference[string]
	Channel      low.NodeReference[*low.Reference] // reference only
	Title        low.NodeReference[string]
	Summary      low.NodeReference[string]
	Description  low.NodeReference[string]
	Security     low.NodeReference[[]low.ValueReference[*SecurityScheme]]
	Tags         low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs low.NodeReference[*ExternalDoc]
	Bindings     low.NodeReference[*OperationBindings]
	Traits       low.NodeReference[[]low.ValueReference[*OperationTrait]]
	Messages     low.NodeReference[[]low.ValueReference[*low.Reference]] // array of refs
	Reply        low.NodeReference[*OperationReply]
	Extensions   *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode      *yaml.Node
	RootNode     *yaml.Node
	idx          *index.SpecIndex
	ctx          context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Operation object.
func (o *Operation) GetRootNode() *yaml.Node {
	return o.RootNode
}

// GetKeyNode returns the key yaml node of the Operation object.
func (o *Operation) GetKeyNode() *yaml.Node {
	return o.KeyNode
}

// GetExtensions returns all extensions for Operation.
func (o *Operation) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return o.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (o *Operation) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, o.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Operation object.
func (o *Operation) GetIndex() *index.SpecIndex {
	return o.idx
}

// GetContext returns the context.Context instance used when building the Operation object.
func (o *Operation) GetContext() context.Context {
	return o.ctx
}

// Build extracts the Operation object from the supplied root node.
func (o *Operation) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	o.KeyNode = keyNode
	root = utils.NodeAlias(root)
	o.RootNode = root
	utils.CheckForMergeNodes(root)
	o.Reference = new(low.Reference)
	o.Nodes = low.ExtractNodes(ctx, root)
	o.Extensions = low.ExtractExtensions(root)
	o.idx = idx
	o.ctx = ctx
	var buildErrs []error

	// extract channel reference (channel is always a $ref in AsyncAPI operations)
	// The channel field format is:
	//   channel:
	//     $ref: '#/channels/...'
	// which means chanValue is a mapping node containing the $ref key
	_, chanLabel, chanValue := utils.FindKeyNodeFullTop(ChannelLabel, root.Content)
	if chanValue != nil {
		o.Nodes.Store(chanLabel.Line, chanLabel)
		ref := referenceFromNode(chanValue)
		if ref != nil {
			o.Channel = low.NodeReference[*low.Reference]{
				Value:     ref,
				KeyNode:   chanLabel,
				ValueNode: chanValue,
			}
		} else {
			buildErrs = append(buildErrs, fmt.Errorf("operation channel must be a Reference Object containing a non-empty $ref string, line %d, column %d",
				chanValue.Line, chanValue.Column))
		}
	}

	// extract reply
	reply, err := low.ExtractObject[*OperationReply](ctx, ReplyLabel, root, idx)
	if err != nil {
		buildErrs = append(buildErrs, err)
	}
	o.Reply = reply

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		buildErrs = append(buildErrs, err)
	}
	if tags != nil {
		o.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	o.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*OperationBindings](ctx, BindingsLabel, root, idx)
	o.Bindings = bindings

	// extract traits
	traits, trLabel, trValue, err := low.ExtractArray[*OperationTrait](ctx, TraitsLabel, root, idx)
	if err != nil {
		buildErrs = append(buildErrs, err)
	}
	if traits != nil {
		o.Traits = low.NodeReference[[]low.ValueReference[*OperationTrait]]{
			Value:     traits,
			KeyNode:   trLabel,
			ValueNode: trValue,
		}
	}

	// extract security
	security, sLabel, sValue, err := low.ExtractArray[*SecurityScheme](ctx, SecurityLabel, root, idx)
	if err != nil {
		buildErrs = append(buildErrs, err)
	}
	if security != nil {
		o.Security = low.NodeReference[[]low.ValueReference[*SecurityScheme]]{
			Value:     security,
			KeyNode:   sLabel,
			ValueNode: sValue,
		}
	}

	// extract messages (array of references - each element is a mapping with $ref)
	// The format is:
	//   messages:
	//     - $ref: '#/channels/.../messages/...'
	_, msgsLabel, msgsValue := utils.FindKeyNodeFullTop(MessagesLabel, root.Content)
	if msgsValue != nil && msgsValue.Kind != yaml.SequenceNode {
		buildErrs = append(buildErrs, fmt.Errorf("operation messages must be an array of Reference Objects, line %d, column %d",
			msgsValue.Line, msgsValue.Column))
	}
	if msgsValue != nil && msgsValue.Kind == yaml.SequenceNode {
		o.Nodes.Store(msgsLabel.Line, msgsLabel)
		var refs []low.ValueReference[*low.Reference]
		for _, msgNode := range msgsValue.Content {
			ref := referenceFromNode(msgNode)
			if ref == nil {
				buildErrs = append(buildErrs, fmt.Errorf("operation messages entry must be a Reference Object containing a non-empty $ref string, line %d, column %d",
					msgNode.Line, msgNode.Column))
				continue
			}
			refs = append(refs, low.ValueReference[*low.Reference]{
				Value:     ref,
				ValueNode: msgNode,
			})
		}
		if len(refs) > 0 {
			o.Messages = low.NodeReference[[]low.ValueReference[*low.Reference]]{
				Value:     refs,
				KeyNode:   msgsLabel,
				ValueNode: msgsValue,
			}
		}
	}

	return errors.Join(buildErrs...)
}

// Hash returns a process-local content hash of the Operation object.
func (o *Operation) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !o.Action.IsEmpty() {
			h.WriteString(o.Action.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if o.Channel.Value != nil {
			h.WriteString(o.Channel.Value.GetReference())
			h.WriteByte(low.HASH_PIPE)
		}
		if !o.Title.IsEmpty() {
			h.WriteString(o.Title.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !o.Summary.IsEmpty() {
			h.WriteString(o.Summary.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !o.Description.IsEmpty() {
			h.WriteString(o.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if o.Messages.Value != nil {
			for _, msg := range o.Messages.Value {
				if msg.Value != nil {
					h.WriteString(msg.Value.GetReference())
					h.WriteByte(low.HASH_PIPE)
				}
			}
		}
		if !o.Reply.IsEmpty() {
			h.WriteString(low.GenerateHashString(o.Reply.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if o.Security.Value != nil {
			for _, sec := range o.Security.Value {
				h.WriteString(low.GenerateHashString(sec.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if o.Tags.Value != nil {
			for _, tag := range o.Tags.Value {
				h.WriteString(low.GenerateHashString(tag.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !o.ExternalDocs.IsEmpty() {
			h.WriteString(low.GenerateHashString(o.ExternalDocs.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !o.Bindings.IsEmpty() {
			h.WriteString(low.GenerateHashString(o.Bindings.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if o.Traits.Value != nil {
			for _, trait := range o.Traits.Value {
				h.WriteString(low.GenerateHashString(trait.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		for _, ext := range low.HashExtensions(o.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OperationBindings represents a low-level AsyncAPI 3.0 Operation Bindings object.
//
// Map of operation binding objects where the keys are protocol names.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationBindingsObject
type OperationBindings struct {
	HTTP       low.NodeReference[*HTTPOperationBinding]
	Kafka      low.NodeReference[*KafkaOperationBinding]
	AMQP       low.NodeReference[*AMQPOperationBinding]
	MQTT       low.NodeReference[*MQTTOperationBinding]
	SQS        low.NodeReference[*SQSOperationBinding]
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (ob *OperationBindings) GetRootNode() *yaml.Node {
	return ob.RootNode
}

// GetKeyNode returns the key yaml node.
func (ob *OperationBindings) GetKeyNode() *yaml.Node {
	return ob.KeyNode
}

// GetExtensions returns all extensions.
func (ob *OperationBindings) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return ob.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (ob *OperationBindings) GetIndex() *index.SpecIndex {
	return ob.idx
}

// GetContext returns the context.Context instance.
func (ob *OperationBindings) GetContext() context.Context {
	return ob.ctx
}

// Build extracts the OperationBindings object from the supplied root node.
func (ob *OperationBindings) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	ob.KeyNode = keyNode
	root = utils.NodeAlias(root)
	ob.RootNode = root
	utils.CheckForMergeNodes(root)
	ob.Reference = new(low.Reference)
	ob.Nodes = low.ExtractNodes(ctx, root)
	ob.Extensions = low.ExtractExtensions(root)
	ob.idx = idx
	ob.ctx = ctx

	http, _ := low.ExtractObject[*HTTPOperationBinding](ctx, HTTPLabel, root, idx)
	ob.HTTP = http

	kafka, _ := low.ExtractObject[*KafkaOperationBinding](ctx, KafkaLabel, root, idx)
	ob.Kafka = kafka

	amqp, _ := low.ExtractObject[*AMQPOperationBinding](ctx, AMQPLabel, root, idx)
	ob.AMQP = amqp

	mqtt, _ := low.ExtractObject[*MQTTOperationBinding](ctx, MQTTLabel, root, idx)
	ob.MQTT = mqtt

	sqs, _ := low.ExtractObject[*SQSOperationBinding](ctx, SQSLabel, root, idx)
	ob.SQS = sqs

	return nil
}

// Hash returns a process-local content hash.
func (ob *OperationBindings) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !ob.HTTP.IsEmpty() {
			h.WriteString(low.GenerateHashString(ob.HTTP.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ob.Kafka.IsEmpty() {
			h.WriteString(low.GenerateHashString(ob.Kafka.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ob.AMQP.IsEmpty() {
			h.WriteString(low.GenerateHashString(ob.AMQP.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ob.MQTT.IsEmpty() {
			h.WriteString(low.GenerateHashString(ob.MQTT.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ob.SQS.IsEmpty() {
			h.WriteString(low.GenerateHashString(ob.SQS.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(ob.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OperationTrait represents a low-level AsyncAPI 3.0 Operation Trait object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationTraitObject
type OperationTrait struct {
	Title        low.NodeReference[string]
	Summary      low.NodeReference[string]
	Description  low.NodeReference[string]
	Security     low.NodeReference[[]low.ValueReference[*SecurityScheme]]
	Tags         low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs low.NodeReference[*ExternalDoc]
	Bindings     low.NodeReference[*OperationBindings]
	Extensions   *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode      *yaml.Node
	RootNode     *yaml.Node
	idx          *index.SpecIndex
	ctx          context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (ot *OperationTrait) GetRootNode() *yaml.Node {
	return ot.RootNode
}

// GetKeyNode returns the key yaml node.
func (ot *OperationTrait) GetKeyNode() *yaml.Node {
	return ot.KeyNode
}

// GetExtensions returns all extensions.
func (ot *OperationTrait) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return ot.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (ot *OperationTrait) GetIndex() *index.SpecIndex {
	return ot.idx
}

// GetContext returns the context.Context instance.
func (ot *OperationTrait) GetContext() context.Context {
	return ot.ctx
}

// Build extracts the OperationTrait object from the supplied root node.
func (ot *OperationTrait) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	ot.KeyNode = keyNode
	root = utils.NodeAlias(root)
	ot.RootNode = root
	utils.CheckForMergeNodes(root)
	ot.Reference = new(low.Reference)
	ot.Nodes = low.ExtractNodes(ctx, root)
	ot.Extensions = low.ExtractExtensions(root)
	ot.idx = idx
	ot.ctx = ctx

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
	}
	if tags != nil {
		ot.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	ot.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*OperationBindings](ctx, BindingsLabel, root, idx)
	ot.Bindings = bindings

	// extract security
	security, sLabel, sValue, err := low.ExtractArray[*SecurityScheme](ctx, SecurityLabel, root, idx)
	if err != nil {
		return err
	}
	if security != nil {
		ot.Security = low.NodeReference[[]low.ValueReference[*SecurityScheme]]{
			Value:     security,
			KeyNode:   sLabel,
			ValueNode: sValue,
		}
	}

	return nil
}

// Hash returns a process-local content hash.
func (ot *OperationTrait) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !ot.Title.IsEmpty() {
			h.WriteString(ot.Title.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ot.Summary.IsEmpty() {
			h.WriteString(ot.Summary.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ot.Description.IsEmpty() {
			h.WriteString(ot.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if ot.Security.Value != nil {
			for _, sec := range ot.Security.Value {
				h.WriteString(low.GenerateHashString(sec.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if ot.Tags.Value != nil {
			for _, tag := range ot.Tags.Value {
				h.WriteString(low.GenerateHashString(tag.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !ot.ExternalDocs.IsEmpty() {
			h.WriteString(low.GenerateHashString(ot.ExternalDocs.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ot.Bindings.IsEmpty() {
			h.WriteString(low.GenerateHashString(ot.Bindings.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(ot.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OperationReply represents a low-level AsyncAPI 3.0 Operation Reply object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyObject
type OperationReply struct {
	Address    low.NodeReference[*OperationReplyAddress]
	Channel    low.NodeReference[*low.Reference]                       // reference only
	Messages   low.NodeReference[[]low.ValueReference[*low.Reference]] // array of refs
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (or *OperationReply) GetRootNode() *yaml.Node {
	return or.RootNode
}

// GetKeyNode returns the key yaml node.
func (or *OperationReply) GetKeyNode() *yaml.Node {
	return or.KeyNode
}

// GetExtensions returns all extensions.
func (or *OperationReply) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return or.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (or *OperationReply) GetIndex() *index.SpecIndex {
	return or.idx
}

// GetContext returns the context.Context instance.
func (or *OperationReply) GetContext() context.Context {
	return or.ctx
}

// Build extracts the OperationReply object from the supplied root node.
func (or *OperationReply) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	or.KeyNode = keyNode
	root = utils.NodeAlias(root)
	or.RootNode = root
	utils.CheckForMergeNodes(root)
	or.Reference = new(low.Reference)
	or.Nodes = low.ExtractNodes(ctx, root)
	or.Extensions = low.ExtractExtensions(root)
	or.idx = idx
	or.ctx = ctx
	var buildErrs []error

	// extract address
	addr, _ := low.ExtractObject[*OperationReplyAddress](ctx, AddressLabel, root, idx)
	or.Address = addr

	// extract channel reference (channel is always a $ref in AsyncAPI reply)
	// The channel field format is: channel: $ref: '#/...'
	_, chanLabel, chanValue := utils.FindKeyNodeFullTop(ChannelLabel, root.Content)
	if chanValue != nil {
		or.Nodes.Store(chanLabel.Line, chanLabel)
		ref := referenceFromNode(chanValue)
		if ref != nil {
			or.Channel = low.NodeReference[*low.Reference]{
				Value:     ref,
				KeyNode:   chanLabel,
				ValueNode: chanValue,
			}
		} else {
			buildErrs = append(buildErrs, fmt.Errorf("operation reply channel must be a Reference Object containing a non-empty $ref string, line %d, column %d",
				chanValue.Line, chanValue.Column))
		}
	}

	// extract messages (array of references - each element is a mapping with $ref)
	// The format is: messages: - $ref: '#/...'
	_, msgsLabel, msgsValue := utils.FindKeyNodeFullTop(MessagesLabel, root.Content)
	if msgsValue != nil && msgsValue.Kind != yaml.SequenceNode {
		buildErrs = append(buildErrs, fmt.Errorf("operation reply messages must be an array of Reference Objects, line %d, column %d",
			msgsValue.Line, msgsValue.Column))
	}
	if msgsValue != nil && msgsValue.Kind == yaml.SequenceNode {
		or.Nodes.Store(msgsLabel.Line, msgsLabel)
		var refs []low.ValueReference[*low.Reference]
		for _, msgNode := range msgsValue.Content {
			ref := referenceFromNode(msgNode)
			if ref == nil {
				buildErrs = append(buildErrs, fmt.Errorf("operation reply messages entry must be a Reference Object containing a non-empty $ref string, line %d, column %d",
					msgNode.Line, msgNode.Column))
				continue
			}
			refs = append(refs, low.ValueReference[*low.Reference]{
				Value:     ref,
				ValueNode: msgNode,
			})
		}
		if len(refs) > 0 {
			or.Messages = low.NodeReference[[]low.ValueReference[*low.Reference]]{
				Value:     refs,
				KeyNode:   msgsLabel,
				ValueNode: msgsValue,
			}
		}
	}

	return errors.Join(buildErrs...)
}

// Hash returns a process-local content hash.
func (or *OperationReply) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !or.Address.IsEmpty() {
			h.WriteString(low.GenerateHashString(or.Address.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if or.Channel.Value != nil {
			h.WriteString(or.Channel.Value.GetReference())
			h.WriteByte(low.HASH_PIPE)
		}
		if or.Messages.Value != nil {
			for _, msg := range or.Messages.Value {
				if msg.Value != nil {
					h.WriteString(msg.Value.GetReference())
					h.WriteByte(low.HASH_PIPE)
				}
			}
		}
		for _, ext := range low.HashExtensions(or.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OperationReplyAddress represents a low-level AsyncAPI 3.0 Operation Reply Address object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyAddressObject
type OperationReplyAddress struct {
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
func (ora *OperationReplyAddress) GetRootNode() *yaml.Node {
	return ora.RootNode
}

// GetKeyNode returns the key yaml node.
func (ora *OperationReplyAddress) GetKeyNode() *yaml.Node {
	return ora.KeyNode
}

// GetExtensions returns all extensions.
func (ora *OperationReplyAddress) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return ora.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (ora *OperationReplyAddress) GetIndex() *index.SpecIndex {
	return ora.idx
}

// GetContext returns the context.Context instance.
func (ora *OperationReplyAddress) GetContext() context.Context {
	return ora.ctx
}

// Build extracts the OperationReplyAddress object from the supplied root node.
func (ora *OperationReplyAddress) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	ora.KeyNode = keyNode
	root = utils.NodeAlias(root)
	ora.RootNode = root
	utils.CheckForMergeNodes(root)
	ora.Reference = new(low.Reference)
	ora.Nodes = low.ExtractNodes(ctx, root)
	ora.Extensions = low.ExtractExtensions(root)
	ora.idx = idx
	ora.ctx = ctx
	return nil
}

// Hash returns a process-local content hash.
func (ora *OperationReplyAddress) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !ora.Description.IsEmpty() {
			h.WriteString(ora.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ora.Location.IsEmpty() {
			h.WriteString(ora.Location.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(ora.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
