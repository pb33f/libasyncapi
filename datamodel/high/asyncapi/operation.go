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

// Operation represents a high-level AsyncAPI 3.0 Operation object.
//
// Describes a publish or subscribe operation.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationObject
type Operation struct {
	Action       string                              `json:"action,omitempty" yaml:"action,omitempty"`
	Channel      *low.Reference                      `json:"channel,omitempty" yaml:"channel,omitempty"`
	Title        string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Messages     []*low.Reference                    `json:"messages,omitempty" yaml:"messages,omitempty"`
	Reply        *OperationReply                     `json:"reply,omitempty" yaml:"reply,omitempty"`
	Tags         []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *OperationBindings                  `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Traits       []*OperationTrait                   `json:"traits,omitempty" yaml:"traits,omitempty"`
	Security     []*SecurityScheme                   `json:"security,omitempty" yaml:"security,omitempty"`
	Extensions   *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low          *lowasync.Operation
}

// NewOperation creates a new high-level Operation instance from a low-level one.
func NewOperation(op *lowasync.Operation) *Operation {
	o := new(Operation)
	o.low = op
	o.Action = op.Action.Value
	o.Title = op.Title.Value
	o.Summary = op.Summary.Value
	o.Description = op.Description.Value

	if !op.Channel.IsEmpty() {
		o.Channel = op.Channel.Value
	}
	if op.Messages.Value != nil {
		for _, msg := range op.Messages.Value {
			o.Messages = append(o.Messages, msg.Value)
		}
	}
	if !op.Reply.IsEmpty() {
		o.Reply = NewOperationReply(op.Reply.Value)
	}
	if op.Tags.Value != nil {
		for _, tag := range op.Tags.Value {
			o.Tags = append(o.Tags, NewTag(tag.Value))
		}
	}
	if !op.ExternalDocs.IsEmpty() {
		o.ExternalDocs = NewExternalDoc(op.ExternalDocs.Value)
	}
	if !op.Bindings.IsEmpty() {
		o.Bindings = NewOperationBindings(op.Bindings.Value)
	}
	if op.Traits.Value != nil {
		for _, trait := range op.Traits.Value {
			o.Traits = append(o.Traits, NewOperationTrait(trait.Value))
		}
	}
	if op.Security.Value != nil {
		for _, sec := range op.Security.Value {
			o.Security = append(o.Security, NewSecurityScheme(sec.Value))
		}
	}
	if orderedmap.Len(op.Extensions) > 0 {
		o.Extensions = high.ExtractExtensions(op.Extensions)
	}
	return o
}

// GoLow returns the low-level Operation instance.
func (o *Operation) GoLow() *lowasync.Operation {
	return o.low
}

// GoLowUntyped returns the low-level Operation instance with no type.
func (o *Operation) GoLowUntyped() any {
	return o.low
}

// Render will return a YAML representation of the Operation object as a byte slice.
func (o *Operation) Render() ([]byte, error) {
	return yaml.Marshal(o)
}

// MarshalYAML will create a ready to render YAML representation of the Operation object.
func (o *Operation) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(o, o.low)
	return nb.Render(), nil
}

// OperationBindings represents a high-level AsyncAPI 3.0 Operation Bindings object.
type OperationBindings struct {
	HTTP       *HTTPOperationBinding               `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka      *KafkaOperationBinding              `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP       *AMQPOperationBinding               `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTT       *MQTTOperationBinding               `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS        *SQSOperationBinding                `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.OperationBindings
}

// NewOperationBindings creates a new high-level OperationBindings instance.
func NewOperationBindings(ob *lowasync.OperationBindings) *OperationBindings {
	b := new(OperationBindings)
	b.low = ob
	if !ob.HTTP.IsEmpty() {
		b.HTTP = NewHTTPOperationBinding(ob.HTTP.Value)
	}
	if !ob.Kafka.IsEmpty() {
		b.Kafka = NewKafkaOperationBinding(ob.Kafka.Value)
	}
	if !ob.AMQP.IsEmpty() {
		b.AMQP = NewAMQPOperationBinding(ob.AMQP.Value)
	}
	if !ob.MQTT.IsEmpty() {
		b.MQTT = NewMQTTOperationBinding(ob.MQTT.Value)
	}
	if !ob.SQS.IsEmpty() {
		b.SQS = NewSQSOperationBinding(ob.SQS.Value)
	}
	if orderedmap.Len(ob.Extensions) > 0 {
		b.Extensions = high.ExtractExtensions(ob.Extensions)
	}
	return b
}

// GoLow returns the low-level OperationBindings instance.
func (b *OperationBindings) GoLow() *lowasync.OperationBindings {
	return b.low
}

// GoLowUntyped returns the low-level OperationBindings instance with no type.
func (b *OperationBindings) GoLowUntyped() any {
	return b.low
}

// OperationTrait represents a high-level AsyncAPI 3.0 Operation Trait object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationTraitObject
type OperationTrait struct {
	Title        string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Security     []*SecurityScheme                   `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *OperationBindings                  `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Extensions   *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low          *lowasync.OperationTrait
}

// NewOperationTrait creates a new high-level OperationTrait instance from a low-level one.
func NewOperationTrait(trait *lowasync.OperationTrait) *OperationTrait {
	t := new(OperationTrait)
	t.low = trait
	t.Title = trait.Title.Value
	t.Summary = trait.Summary.Value
	t.Description = trait.Description.Value
	if trait.Security.Value != nil {
		for _, sec := range trait.Security.Value {
			t.Security = append(t.Security, NewSecurityScheme(sec.Value))
		}
	}
	if trait.Tags.Value != nil {
		for _, tag := range trait.Tags.Value {
			t.Tags = append(t.Tags, NewTag(tag.Value))
		}
	}
	if !trait.ExternalDocs.IsEmpty() {
		t.ExternalDocs = NewExternalDoc(trait.ExternalDocs.Value)
	}
	if !trait.Bindings.IsEmpty() {
		t.Bindings = NewOperationBindings(trait.Bindings.Value)
	}
	if orderedmap.Len(trait.Extensions) > 0 {
		t.Extensions = high.ExtractExtensions(trait.Extensions)
	}
	return t
}

// GoLow returns the low-level OperationTrait instance.
func (t *OperationTrait) GoLow() *lowasync.OperationTrait {
	return t.low
}

// GoLowUntyped returns the low-level OperationTrait instance with no type.
func (t *OperationTrait) GoLowUntyped() any {
	return t.low
}

// OperationReply represents a high-level AsyncAPI 3.0 Operation Reply object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyObject
type OperationReply struct {
	Address    *OperationReplyAddress              `json:"address,omitempty" yaml:"address,omitempty"`
	Channel    *low.Reference                      `json:"channel,omitempty" yaml:"channel,omitempty"`
	Messages   []*low.Reference                    `json:"messages,omitempty" yaml:"messages,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.OperationReply
}

// NewOperationReply creates a new high-level OperationReply instance from a low-level one.
func NewOperationReply(reply *lowasync.OperationReply) *OperationReply {
	r := new(OperationReply)
	r.low = reply
	if !reply.Address.IsEmpty() {
		r.Address = NewOperationReplyAddress(reply.Address.Value)
	}
	if !reply.Channel.IsEmpty() {
		r.Channel = reply.Channel.Value
	}
	if reply.Messages.Value != nil {
		for _, msg := range reply.Messages.Value {
			r.Messages = append(r.Messages, msg.Value)
		}
	}
	if orderedmap.Len(reply.Extensions) > 0 {
		r.Extensions = high.ExtractExtensions(reply.Extensions)
	}
	return r
}

// GoLow returns the low-level OperationReply instance.
func (r *OperationReply) GoLow() *lowasync.OperationReply {
	return r.low
}

// GoLowUntyped returns the low-level OperationReply instance with no type.
func (r *OperationReply) GoLowUntyped() any {
	return r.low
}

// OperationReplyAddress represents a high-level AsyncAPI 3.0 Operation Reply Address object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#operationReplyAddressObject
type OperationReplyAddress struct {
	Location    string                              `json:"location,omitempty" yaml:"location,omitempty"`
	Description string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Extensions  *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low         *lowasync.OperationReplyAddress
}

// NewOperationReplyAddress creates a new high-level OperationReplyAddress instance from a low-level one.
func NewOperationReplyAddress(addr *lowasync.OperationReplyAddress) *OperationReplyAddress {
	a := new(OperationReplyAddress)
	a.low = addr
	a.Location = addr.Location.Value
	a.Description = addr.Description.Value
	if orderedmap.Len(addr.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(addr.Extensions)
	}
	return a
}

// GoLow returns the low-level OperationReplyAddress instance.
func (a *OperationReplyAddress) GoLow() *lowasync.OperationReplyAddress {
	return a.low
}

// GoLowUntyped returns the low-level OperationReplyAddress instance with no type.
func (a *OperationReplyAddress) GoLowUntyped() any {
	return a.low
}
