// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"go.yaml.in/yaml/v4"
)

// Message represents a high-level AsyncAPI 3.0 Message object.
//
// Describes a message received on a given channel and operation.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageObject
type Message struct {
	Headers       *highbase.SchemaProxy               `json:"headers,omitempty" yaml:"headers,omitempty"`
	Payload       *highbase.SchemaProxy               `json:"payload,omitempty" yaml:"payload,omitempty"`
	CorrelationID *CorrelationID                      `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	ContentType   string                              `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Name          string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Title         string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Summary       string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Tags          []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs  *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings      *MessageBindings                    `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Examples      []*MessageExample                   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Traits        []*MessageTrait                     `json:"traits,omitempty" yaml:"traits,omitempty"`
	Extensions    *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low           *lowasync.Message
}

// NewMessage creates a new high-level Message instance from a low-level one.
func NewMessage(msg *lowasync.Message) *Message {
	m := new(Message)
	m.low = msg
	m.ContentType = msg.ContentType.Value
	m.Name = msg.Name.Value
	m.Title = msg.Title.Value
	m.Summary = msg.Summary.Value
	m.Description = msg.Description.Value

	if !msg.Headers.IsEmpty() {
		m.Headers = highbase.NewSchemaProxy(&msg.Headers)
	}
	if !msg.Payload.IsEmpty() {
		m.Payload = highbase.NewSchemaProxy(&msg.Payload)
	}
	if !msg.CorrelationID.IsEmpty() {
		m.CorrelationID = NewCorrelationID(msg.CorrelationID.Value)
	}
	if msg.Tags.Value != nil {
		for _, tag := range msg.Tags.Value {
			m.Tags = append(m.Tags, NewTag(tag.Value))
		}
	}
	if !msg.ExternalDocs.IsEmpty() {
		m.ExternalDocs = NewExternalDoc(msg.ExternalDocs.Value)
	}
	if !msg.Bindings.IsEmpty() {
		m.Bindings = NewMessageBindings(msg.Bindings.Value)
	}
	if msg.Examples.Value != nil {
		for _, ex := range msg.Examples.Value {
			m.Examples = append(m.Examples, NewMessageExample(ex.Value))
		}
	}
	if msg.Traits.Value != nil {
		for _, trait := range msg.Traits.Value {
			m.Traits = append(m.Traits, NewMessageTrait(trait.Value))
		}
	}
	if orderedmap.Len(msg.Extensions) > 0 {
		m.Extensions = high.ExtractExtensions(msg.Extensions)
	}
	return m
}

// GoLow returns the low-level Message instance.
func (m *Message) GoLow() *lowasync.Message {
	return m.low
}

// GoLowUntyped returns the low-level Message instance with no type.
func (m *Message) GoLowUntyped() any {
	return m.low
}

// Render will return a YAML representation of the Message object as a byte slice.
func (m *Message) Render() ([]byte, error) {
	return yaml.Marshal(m)
}

// MarshalYAML will create a ready to render YAML representation of the Message object.
func (m *Message) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(m, m.low)
	return nb.Render(), nil
}

// MessageBindings represents a high-level AsyncAPI 3.0 Message Bindings object.
type MessageBindings struct {
	HTTP       *HTTPMessageBinding                 `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka      *KafkaMessageBinding                `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP       *AMQPMessageBinding                 `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTT       *MQTTMessageBinding                 `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS        *SQSMessageBinding                  `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.MessageBindings
}

// NewMessageBindings creates a new high-level MessageBindings instance.
func NewMessageBindings(mb *lowasync.MessageBindings) *MessageBindings {
	b := new(MessageBindings)
	b.low = mb
	if !mb.HTTP.IsEmpty() {
		b.HTTP = NewHTTPMessageBinding(mb.HTTP.Value)
	}
	if !mb.Kafka.IsEmpty() {
		b.Kafka = NewKafkaMessageBinding(mb.Kafka.Value)
	}
	if !mb.AMQP.IsEmpty() {
		b.AMQP = NewAMQPMessageBinding(mb.AMQP.Value)
	}
	if !mb.MQTT.IsEmpty() {
		b.MQTT = NewMQTTMessageBinding(mb.MQTT.Value)
	}
	if !mb.SQS.IsEmpty() {
		b.SQS = NewSQSMessageBinding(mb.SQS.Value)
	}
	if orderedmap.Len(mb.Extensions) > 0 {
		b.Extensions = high.ExtractExtensions(mb.Extensions)
	}
	return b
}

// GoLow returns the low-level MessageBindings instance.
func (b *MessageBindings) GoLow() *lowasync.MessageBindings {
	return b.low
}

// GoLowUntyped returns the low-level MessageBindings instance with no type.
func (b *MessageBindings) GoLowUntyped() any {
	return b.low
}

// MessageExample represents a high-level AsyncAPI 3.0 Message Example object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageExampleObject
type MessageExample struct {
	Headers    *yaml.Node                          `json:"headers,omitempty" yaml:"headers,omitempty"`
	Payload    *yaml.Node                          `json:"payload,omitempty" yaml:"payload,omitempty"`
	Name       string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Summary    string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.MessageExample
}

// NewMessageExample creates a new high-level MessageExample instance from a low-level one.
func NewMessageExample(ex *lowasync.MessageExample) *MessageExample {
	e := new(MessageExample)
	e.low = ex
	e.Headers = ex.Headers.Value
	e.Payload = ex.Payload.Value
	e.Name = ex.Name.Value
	e.Summary = ex.Summary.Value
	if orderedmap.Len(ex.Extensions) > 0 {
		e.Extensions = high.ExtractExtensions(ex.Extensions)
	}
	return e
}

// GoLow returns the low-level MessageExample instance.
func (e *MessageExample) GoLow() *lowasync.MessageExample {
	return e.low
}

// GoLowUntyped returns the low-level MessageExample instance with no type.
func (e *MessageExample) GoLowUntyped() any {
	return e.low
}

// MessageTrait represents a high-level AsyncAPI 3.0 Message Trait object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#messageTraitObject
type MessageTrait struct {
	Headers       *highbase.SchemaProxy               `json:"headers,omitempty" yaml:"headers,omitempty"`
	CorrelationID *CorrelationID                      `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	ContentType   string                              `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Name          string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Title         string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Summary       string                              `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Tags          []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs  *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings      *MessageBindings                    `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Examples      []*MessageExample                   `json:"examples,omitempty" yaml:"examples,omitempty"`
	Extensions    *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low           *lowasync.MessageTrait
}

// NewMessageTrait creates a new high-level MessageTrait instance from a low-level one.
func NewMessageTrait(trait *lowasync.MessageTrait) *MessageTrait {
	t := new(MessageTrait)
	t.low = trait
	t.ContentType = trait.ContentType.Value
	t.Name = trait.Name.Value
	t.Title = trait.Title.Value
	t.Summary = trait.Summary.Value
	t.Description = trait.Description.Value

	if !trait.Headers.IsEmpty() {
		t.Headers = highbase.NewSchemaProxy(&trait.Headers)
	}
	if !trait.CorrelationID.IsEmpty() {
		t.CorrelationID = NewCorrelationID(trait.CorrelationID.Value)
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
		t.Bindings = NewMessageBindings(trait.Bindings.Value)
	}
	if trait.Examples.Value != nil {
		for _, ex := range trait.Examples.Value {
			t.Examples = append(t.Examples, NewMessageExample(ex.Value))
		}
	}
	if orderedmap.Len(trait.Extensions) > 0 {
		t.Extensions = high.ExtractExtensions(trait.Extensions)
	}
	return t
}

// GoLow returns the low-level MessageTrait instance.
func (t *MessageTrait) GoLow() *lowasync.MessageTrait {
	return t.low
}

// GoLowUntyped returns the low-level MessageTrait instance with no type.
func (t *MessageTrait) GoLowUntyped() any {
	return t.low
}

// CorrelationID represents a high-level AsyncAPI 3.0 Correlation ID object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#correlationIdObject
type CorrelationID struct {
	Description string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Location    string                              `json:"location,omitempty" yaml:"location,omitempty"`
	Extensions  *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low         *lowasync.CorrelationID
}

// NewCorrelationID creates a new high-level CorrelationID instance from a low-level one.
func NewCorrelationID(cid *lowasync.CorrelationID) *CorrelationID {
	c := new(CorrelationID)
	c.low = cid
	c.Description = cid.Description.Value
	c.Location = cid.Location.Value
	if orderedmap.Len(cid.Extensions) > 0 {
		c.Extensions = high.ExtractExtensions(cid.Extensions)
	}
	return c
}

// GoLow returns the low-level CorrelationID instance.
func (c *CorrelationID) GoLow() *lowasync.CorrelationID {
	return c.low
}

// GoLowUntyped returns the low-level CorrelationID instance with no type.
func (c *CorrelationID) GoLowUntyped() any {
	return c.low
}
