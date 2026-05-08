// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
	"go.yaml.in/yaml/v4"
)

// Channel represents a high-level AsyncAPI 3.0 Channel object.
//
// Describes a shared communication channel.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#channelObject
type Channel struct {
	Address      *string                             `json:"address,omitempty" yaml:"address,omitempty"`
	Messages     *orderedmap.Map[string, *Message]   `json:"messages,omitempty" yaml:"messages,omitempty"`
	Title        string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Servers      []*low.Reference                    `json:"servers,omitempty" yaml:"servers,omitempty"`
	Parameters   *orderedmap.Map[string, *Parameter] `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Tags         []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *ChannelBindings                    `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Extensions   *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low          *lowasync.Channel
}

// NewChannel creates a new high-level Channel instance from a low-level one.
func NewChannel(channel *lowasync.Channel) *Channel {
	c := new(Channel)
	c.low = channel
	c.Address = channel.Address.Value
	c.Title = channel.Title.Value
	c.Summary = channel.Summary.Value
	c.Description = channel.Description.Value

	if channel.Messages.Value != nil {
		c.Messages = low.FromReferenceMapWithFunc(channel.Messages.Value, NewMessage)
	}
	if channel.Servers.Value != nil {
		for _, srv := range channel.Servers.Value {
			c.Servers = append(c.Servers, srv.Value)
		}
	}
	if channel.Parameters.Value != nil {
		c.Parameters = low.FromReferenceMapWithFunc(channel.Parameters.Value, NewParameter)
	}
	if channel.Tags.Value != nil {
		for _, tag := range channel.Tags.Value {
			c.Tags = append(c.Tags, NewTag(tag.Value))
		}
	}
	if !channel.ExternalDocs.IsEmpty() {
		c.ExternalDocs = NewExternalDoc(channel.ExternalDocs.Value)
	}
	if !channel.Bindings.IsEmpty() {
		c.Bindings = NewChannelBindings(channel.Bindings.Value)
	}
	if orderedmap.Len(channel.Extensions) > 0 {
		c.Extensions = high.ExtractExtensions(channel.Extensions)
	}
	return c
}

// GoLow returns the low-level Channel instance.
func (c *Channel) GoLow() *lowasync.Channel {
	return c.low
}

// GoLowUntyped returns the low-level Channel instance with no type.
func (c *Channel) GoLowUntyped() any {
	return c.low
}

// Render will return a YAML representation of the Channel object as a byte slice.
func (c *Channel) Render() ([]byte, error) {
	return yaml.Marshal(c)
}

// MarshalYAML will create a ready to render YAML representation of the Channel object.
func (c *Channel) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(c, c.low)
	return nb.Render(), nil
}

// ChannelBindings represents a high-level AsyncAPI 3.0 Channel Bindings object.
type ChannelBindings struct {
	HTTP       *HTTPChannelBinding                 `json:"http,omitempty" yaml:"http,omitempty"`
	WebSocket  *WebSocketChannelBinding            `json:"ws,omitempty" yaml:"ws,omitempty"`
	Kafka      *KafkaChannelBinding                `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP       *AMQPChannelBinding                 `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	SQS        *SQSChannelBinding                  `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.ChannelBindings
}

// NewChannelBindings creates a new high-level ChannelBindings instance.
func NewChannelBindings(cb *lowasync.ChannelBindings) *ChannelBindings {
	b := new(ChannelBindings)
	b.low = cb
	if !cb.HTTP.IsEmpty() {
		b.HTTP = NewHTTPChannelBinding(cb.HTTP.Value)
	}
	if !cb.WebSocket.IsEmpty() {
		b.WebSocket = NewWebSocketChannelBinding(cb.WebSocket.Value)
	}
	if !cb.Kafka.IsEmpty() {
		b.Kafka = NewKafkaChannelBinding(cb.Kafka.Value)
	}
	if !cb.AMQP.IsEmpty() {
		b.AMQP = NewAMQPChannelBinding(cb.AMQP.Value)
	}
	if !cb.SQS.IsEmpty() {
		b.SQS = NewSQSChannelBinding(cb.SQS.Value)
	}
	if orderedmap.Len(cb.Extensions) > 0 {
		b.Extensions = high.ExtractExtensions(cb.Extensions)
	}
	return b
}

// GoLow returns the low-level ChannelBindings instance.
func (b *ChannelBindings) GoLow() *lowasync.ChannelBindings {
	return b.low
}

// GoLowUntyped returns the low-level ChannelBindings instance with no type.
func (b *ChannelBindings) GoLowUntyped() any {
	return b.low
}

// Parameter represents a high-level AsyncAPI 3.0 Parameter object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#parameterObject
type Parameter struct {
	Enum        []string                            `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string                              `json:"default,omitempty" yaml:"default,omitempty"`
	Description string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Examples    []string                            `json:"examples,omitempty" yaml:"examples,omitempty"`
	Location    string                              `json:"location,omitempty" yaml:"location,omitempty"`
	Extensions  *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low         *lowasync.Parameter
}

// NewParameter creates a new high-level Parameter instance from a low-level one.
func NewParameter(param *lowasync.Parameter) *Parameter {
	p := new(Parameter)
	p.low = param
	p.Default = param.Default.Value
	p.Description = param.Description.Value
	p.Location = param.Location.Value
	if param.Enum.Value != nil {
		for _, e := range param.Enum.Value {
			p.Enum = append(p.Enum, e.Value)
		}
	}
	if param.Examples.Value != nil {
		for _, ex := range param.Examples.Value {
			p.Examples = append(p.Examples, ex.Value)
		}
	}
	if orderedmap.Len(param.Extensions) > 0 {
		p.Extensions = high.ExtractExtensions(param.Extensions)
	}
	return p
}

// GoLow returns the low-level Parameter instance.
func (p *Parameter) GoLow() *lowasync.Parameter {
	return p.low
}

// GoLowUntyped returns the low-level Parameter instance with no type.
func (p *Parameter) GoLowUntyped() any {
	return p.low
}

// Render will return a YAML representation of the Parameter object as a byte slice.
func (p *Parameter) Render() ([]byte, error) {
	return yaml.Marshal(p)
}

// MarshalYAML will create a ready to render YAML representation of the Parameter object.
func (p *Parameter) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(p, p.low)
	return nb.Render(), nil
}
