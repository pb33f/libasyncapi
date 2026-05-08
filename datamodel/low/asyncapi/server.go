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

// Server represents a low-level AsyncAPI 3.0 Server object.
//
// An object representing a message broker, a server or any other kind of computer program capable of
// sending and/or receiving data.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverObject
type Server struct {
	Host            low.NodeReference[string]
	Protocol        low.NodeReference[string]
	ProtocolVersion low.NodeReference[string]
	Pathname        low.NodeReference[string]
	Description     low.NodeReference[string]
	Title           low.NodeReference[string]
	Summary         low.NodeReference[string]
	Variables       low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerVariable]]]
	Security        low.NodeReference[[]low.ValueReference[*SecurityScheme]]
	Tags            low.NodeReference[[]low.ValueReference[*Tag]]
	ExternalDocs    low.NodeReference[*ExternalDoc]
	Bindings        low.NodeReference[*ServerBindings]
	Extensions      *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode         *yaml.Node
	RootNode        *yaml.Node
	idx             *index.SpecIndex
	ctx             context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the Server object.
func (s *Server) GetRootNode() *yaml.Node {
	return s.RootNode
}

// GetKeyNode returns the key yaml node of the Server object.
func (s *Server) GetKeyNode() *yaml.Node {
	return s.KeyNode
}

// GetExtensions returns all extensions for Server.
func (s *Server) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return s.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (s *Server) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, s.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the Server object.
func (s *Server) GetIndex() *index.SpecIndex {
	return s.idx
}

// GetContext returns the context.Context instance used when building the Server object.
func (s *Server) GetContext() context.Context {
	return s.ctx
}

// Build extracts the Server object from the supplied root node.
func (s *Server) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	s.KeyNode = keyNode
	root = utils.NodeAlias(root)
	s.RootNode = root
	utils.CheckForMergeNodes(root)
	s.Reference = new(low.Reference)
	s.Nodes = low.ExtractNodes(ctx, root)
	s.Extensions = low.ExtractExtensions(root)
	s.idx = idx
	s.ctx = ctx

	// extract variables
	vars, vLabel, vValue, err := low.ExtractMap[*ServerVariable](ctx, VariablesLabel, root, idx)
	if err != nil {
		return err
	}
	if vars != nil {
		s.Variables = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[*ServerVariable]]]{
			Value:     vars,
			KeyNode:   vLabel,
			ValueNode: vValue,
		}
	}

	// extract tags
	tags, tLabel, tValue, err := low.ExtractArray[*Tag](ctx, TagsLabel, root, idx)
	if err != nil {
		return err
	}
	if tags != nil {
		s.Tags = low.NodeReference[[]low.ValueReference[*Tag]]{
			Value:     tags,
			KeyNode:   tLabel,
			ValueNode: tValue,
		}
	}

	// extract externalDocs
	extDocs, _ := low.ExtractObject[*ExternalDoc](ctx, ExternalDocsLabel, root, idx)
	s.ExternalDocs = extDocs

	// extract bindings
	bindings, _ := low.ExtractObject[*ServerBindings](ctx, BindingsLabel, root, idx)
	s.Bindings = bindings

	// extract security (array of security scheme references or inline)
	security, sLabel, sValue, err := low.ExtractArray[*SecurityScheme](ctx, SecurityLabel, root, idx)
	if err != nil {
		return err
	}
	if security != nil {
		s.Security = low.NodeReference[[]low.ValueReference[*SecurityScheme]]{
			Value:     security,
			KeyNode:   sLabel,
			ValueNode: sValue,
		}
	}

	return nil
}

// Hash returns a consistent SHA256 Hash of the Server object.
func (s *Server) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)

	if !s.Host.IsEmpty() {
		sb.WriteString(s.Host.Value)
		sb.WriteByte('|')
	}
	if !s.Protocol.IsEmpty() {
		sb.WriteString(s.Protocol.Value)
		sb.WriteByte('|')
	}
	if !s.ProtocolVersion.IsEmpty() {
		sb.WriteString(s.ProtocolVersion.Value)
		sb.WriteByte('|')
	}
	if !s.Pathname.IsEmpty() {
		sb.WriteString(s.Pathname.Value)
		sb.WriteByte('|')
	}
	if !s.Description.IsEmpty() {
		sb.WriteString(s.Description.Value)
		sb.WriteByte('|')
	}
	if !s.Title.IsEmpty() {
		sb.WriteString(s.Title.Value)
		sb.WriteByte('|')
	}
	if !s.Summary.IsEmpty() {
		sb.WriteString(s.Summary.Value)
		sb.WriteByte('|')
	}
	if s.Variables.Value != nil {
		for v := range orderedmap.SortAlpha(s.Variables.Value).ValuesFromOldest() {
			sb.WriteString(low.GenerateHashString(v.Value))
			sb.WriteByte('|')
		}
	}
	if s.Security.Value != nil {
		for _, sec := range s.Security.Value {
			sb.WriteString(low.GenerateHashString(sec.Value))
			sb.WriteByte('|')
		}
	}
	if s.Tags.Value != nil {
		for _, tag := range s.Tags.Value {
			sb.WriteString(low.GenerateHashString(tag.Value))
			sb.WriteByte('|')
		}
	}
	if !s.ExternalDocs.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.ExternalDocs.Value))
		sb.WriteByte('|')
	}
	if !s.Bindings.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.Bindings.Value))
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(s.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// ServerVariable represents a low-level AsyncAPI 3.0 Server Variable object.
//
// An object representing a Server Variable for server URL template substitution.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverVariableObject
type ServerVariable struct {
	Enum        low.NodeReference[[]low.ValueReference[string]]
	Default     low.NodeReference[string]
	Description low.NodeReference[string]
	Examples    low.NodeReference[[]low.ValueReference[string]]
	Extensions  *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode     *yaml.Node
	RootNode    *yaml.Node
	idx         *index.SpecIndex
	ctx         context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the ServerVariable object.
func (sv *ServerVariable) GetRootNode() *yaml.Node {
	return sv.RootNode
}

// GetKeyNode returns the key yaml node of the ServerVariable object.
func (sv *ServerVariable) GetKeyNode() *yaml.Node {
	return sv.KeyNode
}

// GetExtensions returns all extensions for ServerVariable.
func (sv *ServerVariable) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return sv.Extensions
}

// GetIndex returns the index.SpecIndex instance attached to the ServerVariable object.
func (sv *ServerVariable) GetIndex() *index.SpecIndex {
	return sv.idx
}

// GetContext returns the context.Context instance used when building the ServerVariable object.
func (sv *ServerVariable) GetContext() context.Context {
	return sv.ctx
}

// Build extracts the ServerVariable object from the supplied root node.
func (sv *ServerVariable) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	sv.KeyNode = keyNode
	root = utils.NodeAlias(root)
	sv.RootNode = root
	utils.CheckForMergeNodes(root)
	sv.Reference = new(low.Reference)
	sv.Nodes = low.ExtractNodes(ctx, root)
	sv.Extensions = low.ExtractExtensions(root)
	sv.idx = idx
	sv.ctx = ctx
	return nil
}

// Hash returns a consistent SHA256 Hash of the ServerVariable object.
func (sv *ServerVariable) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)

	if !sv.Default.IsEmpty() {
		sb.WriteString(sv.Default.Value)
		sb.WriteByte('|')
	}
	if !sv.Description.IsEmpty() {
		sb.WriteString(sv.Description.Value)
		sb.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(sv.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb.String()))
}

// ServerBindings represents a low-level AsyncAPI 3.0 Server Bindings object.
//
// Map of server binding objects where the keys are protocol names.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverBindingsObject
type ServerBindings struct {
	HTTP       low.NodeReference[*HTTPServerBinding]
	Kafka      low.NodeReference[*KafkaServerBinding]
	MQTT       low.NodeReference[*MQTTServerBinding]
	SQS        low.NodeReference[*SQSServerBinding]
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node.
func (sb *ServerBindings) GetRootNode() *yaml.Node {
	return sb.RootNode
}

// GetKeyNode returns the key yaml node.
func (sb *ServerBindings) GetKeyNode() *yaml.Node {
	return sb.KeyNode
}

// GetExtensions returns all extensions.
func (sb *ServerBindings) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return sb.Extensions
}

// GetIndex returns the index.SpecIndex instance.
func (sb *ServerBindings) GetIndex() *index.SpecIndex {
	return sb.idx
}

// GetContext returns the context.Context instance.
func (sb *ServerBindings) GetContext() context.Context {
	return sb.ctx
}

// Build extracts the ServerBindings object from the supplied root node.
func (sb *ServerBindings) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	sb.KeyNode = keyNode
	root = utils.NodeAlias(root)
	sb.RootNode = root
	utils.CheckForMergeNodes(root)
	sb.Reference = new(low.Reference)
	sb.Nodes = low.ExtractNodes(ctx, root)
	sb.Extensions = low.ExtractExtensions(root)
	sb.idx = idx
	sb.ctx = ctx

	http, _ := low.ExtractObject[*HTTPServerBinding](ctx, HTTPLabel, root, idx)
	sb.HTTP = http

	kafka, _ := low.ExtractObject[*KafkaServerBinding](ctx, KafkaLabel, root, idx)
	sb.Kafka = kafka

	mqtt, _ := low.ExtractObject[*MQTTServerBinding](ctx, MQTTLabel, root, idx)
	sb.MQTT = mqtt

	sqs, _ := low.ExtractObject[*SQSServerBinding](ctx, SQSLabel, root, idx)
	sb.SQS = sqs

	return nil
}

// Hash returns a consistent SHA256 Hash.
func (sb *ServerBindings) Hash() [32]byte {
	sb2 := low.GetStringBuilder()
	defer low.PutStringBuilder(sb2)
	if !sb.HTTP.IsEmpty() {
		sb2.WriteString(low.GenerateHashString(sb.HTTP.Value))
		sb2.WriteByte('|')
	}
	if !sb.Kafka.IsEmpty() {
		sb2.WriteString(low.GenerateHashString(sb.Kafka.Value))
		sb2.WriteByte('|')
	}
	if !sb.MQTT.IsEmpty() {
		sb2.WriteString(low.GenerateHashString(sb.MQTT.Value))
		sb2.WriteByte('|')
	}
	if !sb.SQS.IsEmpty() {
		sb2.WriteString(low.GenerateHashString(sb.SQS.Value))
		sb2.WriteByte('|')
	}
	for _, ext := range low.HashExtensions(sb.Extensions) {
		sb2.WriteString(ext)
		sb2.WriteByte('|')
	}
	return sha256.Sum256([]byte(sb2.String()))
}
