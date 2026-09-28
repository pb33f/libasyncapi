// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/low"
	lowbase "github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Components represents a high-level AsyncAPI 3.0 Components object.
//
// Holds a set of reusable objects for different aspects of the AsyncAPI specification.
// All objects defined within the components object will have no effect on the API
// unless they are explicitly referenced from properties outside the components object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#componentsObject
type Components struct {
	Schemas           *orderedmap.Map[string, *highbase.SchemaProxy]  `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Servers           *orderedmap.Map[string, *Server]                `json:"servers,omitempty" yaml:"servers,omitempty"`
	Channels          *orderedmap.Map[string, *Channel]               `json:"channels,omitempty" yaml:"channels,omitempty"`
	Operations        *orderedmap.Map[string, *Operation]             `json:"operations,omitempty" yaml:"operations,omitempty"`
	Messages          *orderedmap.Map[string, *Message]               `json:"messages,omitempty" yaml:"messages,omitempty"`
	SecuritySchemes   *orderedmap.Map[string, *SecurityScheme]        `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	ServerVariables   *orderedmap.Map[string, *ServerVariable]        `json:"serverVariables,omitempty" yaml:"serverVariables,omitempty"`
	Parameters        *orderedmap.Map[string, *Parameter]             `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	CorrelationIDs    *orderedmap.Map[string, *CorrelationID]         `json:"correlationIds,omitempty" yaml:"correlationIds,omitempty"`
	Replies           *orderedmap.Map[string, *OperationReply]        `json:"replies,omitempty" yaml:"replies,omitempty"`
	ReplyAddresses    *orderedmap.Map[string, *OperationReplyAddress] `json:"replyAddresses,omitempty" yaml:"replyAddresses,omitempty"`
	ExternalDocs      *orderedmap.Map[string, *ExternalDoc]           `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Tags              *orderedmap.Map[string, *Tag]                   `json:"tags,omitempty" yaml:"tags,omitempty"`
	OperationTraits   *orderedmap.Map[string, *OperationTrait]        `json:"operationTraits,omitempty" yaml:"operationTraits,omitempty"`
	MessageTraits     *orderedmap.Map[string, *MessageTrait]          `json:"messageTraits,omitempty" yaml:"messageTraits,omitempty"`
	ServerBindings    *orderedmap.Map[string, *ServerBindings]        `json:"serverBindings,omitempty" yaml:"serverBindings,omitempty"`
	ChannelBindings   *orderedmap.Map[string, *ChannelBindings]       `json:"channelBindings,omitempty" yaml:"channelBindings,omitempty"`
	OperationBindings *orderedmap.Map[string, *OperationBindings]     `json:"operationBindings,omitempty" yaml:"operationBindings,omitempty"`
	MessageBindings   *orderedmap.Map[string, *MessageBindings]       `json:"messageBindings,omitempty" yaml:"messageBindings,omitempty"`
	Extensions        *orderedmap.Map[string, *yaml.Node]             `json:"-" yaml:"-"`
	low               *lowasync.Components
}

// buildComponentMap is a generic helper for transforming low-level ordered maps to high-level ones.
func buildComponentMap[LV any, V any](
	lowMap *orderedmap.Map[low.KeyReference[string], low.ValueReference[LV]],
	constructor func(LV) V,
) *orderedmap.Map[string, V] {
	if lowMap == nil {
		return nil
	}
	result := orderedmap.New[string, V]()
	for k, v := range lowMap.FromOldest() {
		result.Set(k.Value, constructor(v.Value))
	}
	return result
}

// NewComponents creates a new high-level Components instance from a low-level one.
func NewComponents(comp *lowasync.Components) *Components {
	c := new(Components)
	c.low = comp

	// Schemas require special handling due to SchemaProxy wrapping
	if comp.Schemas.Value != nil {
		c.Schemas = orderedmap.New[string, *highbase.SchemaProxy]()
		for k, v := range comp.Schemas.Value.FromOldest() {
			c.Schemas.Set(k.Value, highbase.NewSchemaProxy(&low.NodeReference[*lowbase.SchemaProxy]{
				Value:     v.Value,
				ValueNode: v.ValueNode,
			}))
		}
	}

	// All other component types use the generic helper
	c.Servers = buildComponentMap(comp.Servers.Value, NewServer)
	c.Channels = buildComponentMap(comp.Channels.Value, NewChannel)
	c.Operations = buildComponentMap(comp.Operations.Value, NewOperation)
	c.Messages = buildComponentMap(comp.Messages.Value, NewMessage)
	c.SecuritySchemes = buildComponentMap(comp.SecuritySchemes.Value, NewSecurityScheme)
	c.ServerVariables = buildComponentMap(comp.ServerVariables.Value, NewServerVariable)
	c.Parameters = buildComponentMap(comp.Parameters.Value, NewParameter)
	c.CorrelationIDs = buildComponentMap(comp.CorrelationIDs.Value, NewCorrelationID)
	c.Replies = buildComponentMap(comp.Replies.Value, NewOperationReply)
	c.ReplyAddresses = buildComponentMap(comp.ReplyAddresses.Value, NewOperationReplyAddress)
	c.ExternalDocs = buildComponentMap(comp.ExternalDocs.Value, NewExternalDoc)
	c.Tags = buildComponentMap(comp.Tags.Value, NewTag)
	c.OperationTraits = buildComponentMap(comp.OperationTraits.Value, NewOperationTrait)
	c.MessageTraits = buildComponentMap(comp.MessageTraits.Value, NewMessageTrait)
	c.ServerBindings = buildComponentMap(comp.ServerBindings.Value, NewServerBindings)
	c.ChannelBindings = buildComponentMap(comp.ChannelBindings.Value, NewChannelBindings)
	c.OperationBindings = buildComponentMap(comp.OperationBindings.Value, NewOperationBindings)
	c.MessageBindings = buildComponentMap(comp.MessageBindings.Value, NewMessageBindings)

	if orderedmap.Len(comp.Extensions) > 0 {
		c.Extensions = high.ExtractExtensions(comp.Extensions)
	}
	return c
}

// GoLow returns the low-level Components instance.
func (c *Components) GoLow() *lowasync.Components {
	return c.low
}

// GoLowUntyped returns the low-level Components instance with no type.
func (c *Components) GoLowUntyped() any {
	return c.low
}
