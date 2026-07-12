// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// ServerBindingsChanges represents changes made to an AsyncAPI ServerBindings object.
type ServerBindingsChanges struct {
	*PropertyChanges
	HTTPChanges      *HTTPServerBindingChanges  `json:"http,omitempty" yaml:"http,omitempty"`
	KafkaChanges     *KafkaServerBindingChanges `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	MQTTChanges      *MQTTServerBindingChanges  `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQSChanges       *SQSServerBindingChanges   `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	ExtensionChanges *ExtensionChanges          `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between ServerBindings objects.
func (s *ServerBindingsChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.HTTPChanges != nil {
		changes = append(changes, s.HTTPChanges.GetAllChanges()...)
	}
	if s.KafkaChanges != nil {
		changes = append(changes, s.KafkaChanges.GetAllChanges()...)
	}
	if s.MQTTChanges != nil {
		changes = append(changes, s.MQTTChanges.GetAllChanges()...)
	}
	if s.SQSChanges != nil {
		changes = append(changes, s.SQSChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *ServerBindingsChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.HTTPChanges != nil {
		c += s.HTTPChanges.TotalChanges()
	}
	if s.KafkaChanges != nil {
		c += s.KafkaChanges.TotalChanges()
	}
	if s.MQTTChanges != nil {
		c += s.MQTTChanges.TotalChanges()
	}
	if s.SQSChanges != nil {
		c += s.SQSChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *ServerBindingsChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.HTTPChanges != nil {
		c += s.HTTPChanges.TotalBreakingChanges()
	}
	if s.KafkaChanges != nil {
		c += s.KafkaChanges.TotalBreakingChanges()
	}
	if s.MQTTChanges != nil {
		c += s.MQTTChanges.TotalBreakingChanges()
	}
	if s.SQSChanges != nil {
		c += s.SQSChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareServerBindings compares two AsyncAPI ServerBindings objects and returns a pointer
// to ServerBindingsChanges, or nil if nothing changed.
func CompareServerBindings(l, r *lowasync.ServerBindings, configs ...*BreakingRulesConfig) *ServerBindingsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	sb := new(ServerBindingsChanges)

	compareNestedObject(l.HTTP, r.HTTP, lowasync.HTTPLabel,
		CompServerBindings, PropHTTP, &changes, configuredNestedCompare(config, CompareHTTPServerBinding), &sb.HTTPChanges, config)
	compareNestedObject(l.Kafka, r.Kafka, lowasync.KafkaLabel,
		CompServerBindings, PropKafka, &changes, configuredNestedCompare(config, CompareKafkaServerBinding), &sb.KafkaChanges, config)
	compareNestedObject(l.MQTT, r.MQTT, lowasync.MQTTLabel,
		CompServerBindings, PropMQTT, &changes, configuredNestedCompare(config, CompareMQTTServerBinding), &sb.MQTTChanges, config)
	compareNestedObject(l.SQS, r.SQS, lowasync.SQSLabel,
		CompServerBindings, PropSQS, &changes, configuredNestedCompare(config, CompareSQSServerBinding), &sb.SQSChanges, config)

	sb.ExtensionChanges = CheckExtensions(l, r)
	sb.PropertyChanges = NewPropertyChanges(changes)
	if sb.TotalChanges() <= 0 {
		return nil
	}
	return sb
}

// ChannelBindingsChanges represents changes made to an AsyncAPI ChannelBindings object.
type ChannelBindingsChanges struct {
	*PropertyChanges
	HTTPChanges      *HTTPChannelBindingChanges      `json:"http,omitempty" yaml:"http,omitempty"`
	WebSocketChanges *WebSocketChannelBindingChanges `json:"ws,omitempty" yaml:"ws,omitempty"`
	KafkaChanges     *KafkaChannelBindingChanges     `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQPChanges      *AMQPChannelBindingChanges      `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	SQSChanges       *SQSChannelBindingChanges       `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	ExtensionChanges *ExtensionChanges               `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between ChannelBindings objects.
func (c *ChannelBindingsChanges) GetAllChanges() []*Change {
	if c == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, c.Changes...)
	if c.HTTPChanges != nil {
		changes = append(changes, c.HTTPChanges.GetAllChanges()...)
	}
	if c.WebSocketChanges != nil {
		changes = append(changes, c.WebSocketChanges.GetAllChanges()...)
	}
	if c.KafkaChanges != nil {
		changes = append(changes, c.KafkaChanges.GetAllChanges()...)
	}
	if c.AMQPChanges != nil {
		changes = append(changes, c.AMQPChanges.GetAllChanges()...)
	}
	if c.SQSChanges != nil {
		changes = append(changes, c.SQSChanges.GetAllChanges()...)
	}
	if c.ExtensionChanges != nil {
		changes = append(changes, c.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (c *ChannelBindingsChanges) TotalChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalChanges()
	if c.HTTPChanges != nil {
		t += c.HTTPChanges.TotalChanges()
	}
	if c.WebSocketChanges != nil {
		t += c.WebSocketChanges.TotalChanges()
	}
	if c.KafkaChanges != nil {
		t += c.KafkaChanges.TotalChanges()
	}
	if c.AMQPChanges != nil {
		t += c.AMQPChanges.TotalChanges()
	}
	if c.SQSChanges != nil {
		t += c.SQSChanges.TotalChanges()
	}
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalChanges()
	}
	return t
}

// TotalBreakingChanges returns the number of breaking changes made.
func (c *ChannelBindingsChanges) TotalBreakingChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalBreakingChanges()
	if c.HTTPChanges != nil {
		t += c.HTTPChanges.TotalBreakingChanges()
	}
	if c.WebSocketChanges != nil {
		t += c.WebSocketChanges.TotalBreakingChanges()
	}
	if c.KafkaChanges != nil {
		t += c.KafkaChanges.TotalBreakingChanges()
	}
	if c.AMQPChanges != nil {
		t += c.AMQPChanges.TotalBreakingChanges()
	}
	if c.SQSChanges != nil {
		t += c.SQSChanges.TotalBreakingChanges()
	}
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// CompareChannelBindings compares two AsyncAPI ChannelBindings objects and returns a
// pointer to ChannelBindingsChanges, or nil if nothing changed.
func CompareChannelBindings(l, r *lowasync.ChannelBindings, configs ...*BreakingRulesConfig) *ChannelBindingsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	cb := new(ChannelBindingsChanges)

	compareNestedObject(l.HTTP, r.HTTP, lowasync.HTTPLabel,
		CompChannelBindings, PropHTTP, &changes, configuredNestedCompare(config, CompareHTTPChannelBinding), &cb.HTTPChanges, config)
	compareNestedObject(l.WebSocket, r.WebSocket, lowasync.WebSocketLabel,
		CompChannelBindings, PropWS, &changes, configuredNestedCompare(config, CompareWebSocketChannelBinding), &cb.WebSocketChanges, config)
	compareNestedObject(l.Kafka, r.Kafka, lowasync.KafkaLabel,
		CompChannelBindings, PropKafka, &changes, configuredNestedCompare(config, CompareKafkaChannelBinding), &cb.KafkaChanges, config)
	compareNestedObject(l.AMQP, r.AMQP, lowasync.AMQPLabel,
		CompChannelBindings, PropAMQP, &changes, configuredNestedCompare(config, CompareAMQPChannelBinding), &cb.AMQPChanges, config)
	compareNestedObject(l.SQS, r.SQS, lowasync.SQSLabel,
		CompChannelBindings, PropSQS, &changes, configuredNestedCompare(config, CompareSQSChannelBinding), &cb.SQSChanges, config)

	cb.ExtensionChanges = CheckExtensions(l, r)
	cb.PropertyChanges = NewPropertyChanges(changes)
	if cb.TotalChanges() <= 0 {
		return nil
	}
	return cb
}

// OperationBindingsChanges represents changes made to an AsyncAPI OperationBindings object.
type OperationBindingsChanges struct {
	*PropertyChanges
	HTTPChanges      *HTTPOperationBindingChanges  `json:"http,omitempty" yaml:"http,omitempty"`
	KafkaChanges     *KafkaOperationBindingChanges `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQPChanges      *AMQPOperationBindingChanges  `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTTChanges      *MQTTOperationBindingChanges  `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQSChanges       *SQSOperationBindingChanges   `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	ExtensionChanges *ExtensionChanges             `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OperationBindings objects.
func (o *OperationBindingsChanges) GetAllChanges() []*Change {
	if o == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, o.Changes...)
	if o.HTTPChanges != nil {
		changes = append(changes, o.HTTPChanges.GetAllChanges()...)
	}
	if o.KafkaChanges != nil {
		changes = append(changes, o.KafkaChanges.GetAllChanges()...)
	}
	if o.AMQPChanges != nil {
		changes = append(changes, o.AMQPChanges.GetAllChanges()...)
	}
	if o.MQTTChanges != nil {
		changes = append(changes, o.MQTTChanges.GetAllChanges()...)
	}
	if o.SQSChanges != nil {
		changes = append(changes, o.SQSChanges.GetAllChanges()...)
	}
	if o.ExtensionChanges != nil {
		changes = append(changes, o.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (o *OperationBindingsChanges) TotalChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalChanges()
	if o.HTTPChanges != nil {
		c += o.HTTPChanges.TotalChanges()
	}
	if o.KafkaChanges != nil {
		c += o.KafkaChanges.TotalChanges()
	}
	if o.AMQPChanges != nil {
		c += o.AMQPChanges.TotalChanges()
	}
	if o.MQTTChanges != nil {
		c += o.MQTTChanges.TotalChanges()
	}
	if o.SQSChanges != nil {
		c += o.SQSChanges.TotalChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (o *OperationBindingsChanges) TotalBreakingChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalBreakingChanges()
	if o.HTTPChanges != nil {
		c += o.HTTPChanges.TotalBreakingChanges()
	}
	if o.KafkaChanges != nil {
		c += o.KafkaChanges.TotalBreakingChanges()
	}
	if o.AMQPChanges != nil {
		c += o.AMQPChanges.TotalBreakingChanges()
	}
	if o.MQTTChanges != nil {
		c += o.MQTTChanges.TotalBreakingChanges()
	}
	if o.SQSChanges != nil {
		c += o.SQSChanges.TotalBreakingChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOperationBindings compares two AsyncAPI OperationBindings objects and returns a
// pointer to OperationBindingsChanges, or nil if nothing changed.
func CompareOperationBindings(l, r *lowasync.OperationBindings, configs ...*BreakingRulesConfig) *OperationBindingsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	ob := new(OperationBindingsChanges)

	compareNestedObject(l.HTTP, r.HTTP, lowasync.HTTPLabel,
		CompOperationBindings, PropHTTP, &changes, configuredNestedCompare(config, CompareHTTPOperationBinding), &ob.HTTPChanges, config)
	compareNestedObject(l.Kafka, r.Kafka, lowasync.KafkaLabel,
		CompOperationBindings, PropKafka, &changes, configuredNestedCompare(config, CompareKafkaOperationBinding), &ob.KafkaChanges, config)
	compareNestedObject(l.AMQP, r.AMQP, lowasync.AMQPLabel,
		CompOperationBindings, PropAMQP, &changes, configuredNestedCompare(config, CompareAMQPOperationBinding), &ob.AMQPChanges, config)
	compareNestedObject(l.MQTT, r.MQTT, lowasync.MQTTLabel,
		CompOperationBindings, PropMQTT, &changes, configuredNestedCompare(config, CompareMQTTOperationBinding), &ob.MQTTChanges, config)
	compareNestedObject(l.SQS, r.SQS, lowasync.SQSLabel,
		CompOperationBindings, PropSQS, &changes, configuredNestedCompare(config, CompareSQSOperationBinding), &ob.SQSChanges, config)

	ob.ExtensionChanges = CheckExtensions(l, r)
	ob.PropertyChanges = NewPropertyChanges(changes)
	if ob.TotalChanges() <= 0 {
		return nil
	}
	return ob
}

// MessageBindingsChanges represents changes made to an AsyncAPI MessageBindings object.
type MessageBindingsChanges struct {
	*PropertyChanges
	HTTPChanges      *HTTPMessageBindingChanges  `json:"http,omitempty" yaml:"http,omitempty"`
	KafkaChanges     *KafkaMessageBindingChanges `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQPChanges      *AMQPMessageBindingChanges  `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTTChanges      *MQTTMessageBindingChanges  `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQSChanges       *SQSMessageBindingChanges   `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	ExtensionChanges *ExtensionChanges           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between MessageBindings objects.
func (m *MessageBindingsChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.HTTPChanges != nil {
		changes = append(changes, m.HTTPChanges.GetAllChanges()...)
	}
	if m.KafkaChanges != nil {
		changes = append(changes, m.KafkaChanges.GetAllChanges()...)
	}
	if m.AMQPChanges != nil {
		changes = append(changes, m.AMQPChanges.GetAllChanges()...)
	}
	if m.MQTTChanges != nil {
		changes = append(changes, m.MQTTChanges.GetAllChanges()...)
	}
	if m.SQSChanges != nil {
		changes = append(changes, m.SQSChanges.GetAllChanges()...)
	}
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MessageBindingsChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.HTTPChanges != nil {
		c += m.HTTPChanges.TotalChanges()
	}
	if m.KafkaChanges != nil {
		c += m.KafkaChanges.TotalChanges()
	}
	if m.AMQPChanges != nil {
		c += m.AMQPChanges.TotalChanges()
	}
	if m.MQTTChanges != nil {
		c += m.MQTTChanges.TotalChanges()
	}
	if m.SQSChanges != nil {
		c += m.SQSChanges.TotalChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MessageBindingsChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.HTTPChanges != nil {
		c += m.HTTPChanges.TotalBreakingChanges()
	}
	if m.KafkaChanges != nil {
		c += m.KafkaChanges.TotalBreakingChanges()
	}
	if m.AMQPChanges != nil {
		c += m.AMQPChanges.TotalBreakingChanges()
	}
	if m.MQTTChanges != nil {
		c += m.MQTTChanges.TotalBreakingChanges()
	}
	if m.SQSChanges != nil {
		c += m.SQSChanges.TotalBreakingChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMessageBindings compares two AsyncAPI MessageBindings objects and returns a
// pointer to MessageBindingsChanges, or nil if nothing changed.
func CompareMessageBindings(l, r *lowasync.MessageBindings, configs ...*BreakingRulesConfig) *MessageBindingsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	mb := new(MessageBindingsChanges)

	compareNestedObject(l.HTTP, r.HTTP, lowasync.HTTPLabel,
		CompMessageBindings, PropHTTP, &changes, configuredNestedCompare(config, CompareHTTPMessageBinding), &mb.HTTPChanges, config)
	compareNestedObject(l.Kafka, r.Kafka, lowasync.KafkaLabel,
		CompMessageBindings, PropKafka, &changes, configuredNestedCompare(config, CompareKafkaMessageBinding), &mb.KafkaChanges, config)
	compareNestedObject(l.AMQP, r.AMQP, lowasync.AMQPLabel,
		CompMessageBindings, PropAMQP, &changes, configuredNestedCompare(config, CompareAMQPMessageBinding), &mb.AMQPChanges, config)
	compareNestedObject(l.MQTT, r.MQTT, lowasync.MQTTLabel,
		CompMessageBindings, PropMQTT, &changes, configuredNestedCompare(config, CompareMQTTMessageBinding), &mb.MQTTChanges, config)
	compareNestedObject(l.SQS, r.SQS, lowasync.SQSLabel,
		CompMessageBindings, PropSQS, &changes, configuredNestedCompare(config, CompareSQSMessageBinding), &mb.SQSChanges, config)

	mb.ExtensionChanges = CheckExtensions(l, r)
	mb.PropertyChanges = NewPropertyChanges(changes)
	if mb.TotalChanges() <= 0 {
		return nil
	}
	return mb
}

// HTTP Bindings

// HTTPServerBindingChanges represents changes made to an AsyncAPI HTTP Server Binding object.
type HTTPServerBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between HTTP Server Binding objects.
func (h *HTTPServerBindingChanges) GetAllChanges() []*Change {
	if h == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, h.Changes...)
	if h.ExtensionChanges != nil {
		changes = append(changes, h.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (h *HTTPServerBindingChanges) TotalChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalChanges()
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (h *HTTPServerBindingChanges) TotalBreakingChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalBreakingChanges()
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareHTTPServerBinding compares two AsyncAPI HTTP Server Binding objects and returns a
// pointer to HTTPServerBindingChanges, or nil if nothing changed. The binding carries no
// fields beyond extensions.
func CompareHTTPServerBinding(l, r *lowasync.HTTPServerBinding, configs ...*BreakingRulesConfig) *HTTPServerBindingChanges {
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	hb := new(HTTPServerBindingChanges)
	hb.ExtensionChanges = CheckExtensions(l, r)
	hb.PropertyChanges = NewPropertyChanges(changes)
	if hb.TotalChanges() <= 0 {
		return nil
	}
	return hb
}

// HTTPChannelBindingChanges represents changes made to an AsyncAPI HTTP Channel Binding object.
type HTTPChannelBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between HTTP Channel Binding objects.
func (h *HTTPChannelBindingChanges) GetAllChanges() []*Change {
	if h == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, h.Changes...)
	if h.ExtensionChanges != nil {
		changes = append(changes, h.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (h *HTTPChannelBindingChanges) TotalChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalChanges()
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (h *HTTPChannelBindingChanges) TotalBreakingChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalBreakingChanges()
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareHTTPChannelBinding compares two AsyncAPI HTTP Channel Binding objects and returns
// a pointer to HTTPChannelBindingChanges, or nil if nothing changed. The binding carries no
// fields beyond extensions.
func CompareHTTPChannelBinding(l, r *lowasync.HTTPChannelBinding, configs ...*BreakingRulesConfig) *HTTPChannelBindingChanges {
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	hb := new(HTTPChannelBindingChanges)
	hb.ExtensionChanges = CheckExtensions(l, r)
	hb.PropertyChanges = NewPropertyChanges(changes)
	if hb.TotalChanges() <= 0 {
		return nil
	}
	return hb
}

// HTTPOperationBindingChanges represents changes made to an AsyncAPI HTTP Operation Binding object.
type HTTPOperationBindingChanges struct {
	*PropertyChanges
	QueryChanges     *wcmodel.SchemaChanges `json:"query,omitempty" yaml:"query,omitempty"`
	ExtensionChanges *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between HTTP Operation Binding objects.
func (h *HTTPOperationBindingChanges) GetAllChanges() []*Change {
	if h == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, h.Changes...)
	if h.QueryChanges != nil {
		changes = append(changes, h.QueryChanges.GetAllChanges()...)
	}
	if h.ExtensionChanges != nil {
		changes = append(changes, h.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (h *HTTPOperationBindingChanges) TotalChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalChanges()
	if h.QueryChanges != nil {
		c += h.QueryChanges.TotalChanges()
	}
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (h *HTTPOperationBindingChanges) TotalBreakingChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalBreakingChanges()
	if h.QueryChanges != nil {
		c += h.QueryChanges.TotalBreakingChanges()
	}
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareHTTPOperationBinding compares two AsyncAPI HTTP Operation Binding objects and
// returns a pointer to HTTPOperationBindingChanges, or nil if nothing changed.
//
// The query field is a schema, so its comparison is delegated to libopenapi's
// what-changed schema comparator.
func CompareHTTPOperationBinding(l, r *lowasync.HTTPOperationBinding, configs ...*BreakingRulesConfig) *HTTPOperationBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompHTTPOperationBinding, PropMethod,
			l.Method.ValueNode, r.Method.ValueNode, lowasync.MethodLabel, &changes, l, r, config),
		NewPropertyCheck(CompHTTPOperationBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	hb := new(HTTPOperationBindingChanges)

	compareNestedObject(l.Query, r.Query, lowasync.QueryLabel,
		CompHTTPOperationBinding, PropQuery, &changes, wcmodel.CompareSchemas, &hb.QueryChanges, config)

	hb.ExtensionChanges = CheckExtensions(l, r)
	hb.PropertyChanges = NewPropertyChanges(changes)
	if hb.TotalChanges() <= 0 {
		return nil
	}
	return hb
}

// HTTPMessageBindingChanges represents changes made to an AsyncAPI HTTP Message Binding object.
type HTTPMessageBindingChanges struct {
	*PropertyChanges
	HeadersChanges   *wcmodel.SchemaChanges `json:"headers,omitempty" yaml:"headers,omitempty"`
	ExtensionChanges *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between HTTP Message Binding objects.
func (h *HTTPMessageBindingChanges) GetAllChanges() []*Change {
	if h == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, h.Changes...)
	if h.HeadersChanges != nil {
		changes = append(changes, h.HeadersChanges.GetAllChanges()...)
	}
	if h.ExtensionChanges != nil {
		changes = append(changes, h.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (h *HTTPMessageBindingChanges) TotalChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalChanges()
	if h.HeadersChanges != nil {
		c += h.HeadersChanges.TotalChanges()
	}
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (h *HTTPMessageBindingChanges) TotalBreakingChanges() int {
	if h == nil {
		return 0
	}
	c := h.PropertyChanges.TotalBreakingChanges()
	if h.HeadersChanges != nil {
		c += h.HeadersChanges.TotalBreakingChanges()
	}
	if h.ExtensionChanges != nil {
		c += h.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareHTTPMessageBinding compares two AsyncAPI HTTP Message Binding objects and returns
// a pointer to HTTPMessageBindingChanges, or nil if nothing changed.
//
// The headers field is a schema, so its comparison is delegated to libopenapi's
// what-changed schema comparator.
func CompareHTTPMessageBinding(l, r *lowasync.HTTPMessageBinding, configs ...*BreakingRulesConfig) *HTTPMessageBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompHTTPMessageBinding, PropStatusCode,
			l.StatusCode.ValueNode, r.StatusCode.ValueNode, lowasync.StatusCodeLabel, &changes, l, r, config),
		NewPropertyCheck(CompHTTPMessageBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	hb := new(HTTPMessageBindingChanges)

	compareNestedObject(l.Headers, r.Headers, lowasync.HeadersLabel,
		CompHTTPMessageBinding, PropHeaders, &changes, wcmodel.CompareSchemas, &hb.HeadersChanges, config)

	hb.ExtensionChanges = CheckExtensions(l, r)
	hb.PropertyChanges = NewPropertyChanges(changes)
	if hb.TotalChanges() <= 0 {
		return nil
	}
	return hb
}

// Kafka Bindings

// KafkaServerBindingChanges represents changes made to an AsyncAPI Kafka Server Binding object.
type KafkaServerBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Kafka Server Binding objects.
func (k *KafkaServerBindingChanges) GetAllChanges() []*Change {
	if k == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, k.Changes...)
	if k.ExtensionChanges != nil {
		changes = append(changes, k.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (k *KafkaServerBindingChanges) TotalChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalChanges()
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (k *KafkaServerBindingChanges) TotalBreakingChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalBreakingChanges()
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareKafkaServerBinding compares two AsyncAPI Kafka Server Binding objects and returns
// a pointer to KafkaServerBindingChanges, or nil if nothing changed.
func CompareKafkaServerBinding(l, r *lowasync.KafkaServerBinding, configs ...*BreakingRulesConfig) *KafkaServerBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompKafkaServerBinding, PropSchemaRegistryURL,
			l.SchemaRegistryURL.ValueNode, r.SchemaRegistryURL.ValueNode, lowasync.SchemaRegistryURLLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaServerBinding, PropSchemaRegistryVendor,
			l.SchemaRegistryVendor.ValueNode, r.SchemaRegistryVendor.ValueNode, lowasync.SchemaRegistryVendorLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaServerBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	kb := new(KafkaServerBindingChanges)
	kb.ExtensionChanges = CheckExtensions(l, r)
	kb.PropertyChanges = NewPropertyChanges(changes)
	if kb.TotalChanges() <= 0 {
		return nil
	}
	return kb
}

// KafkaChannelBindingChanges represents changes made to an AsyncAPI Kafka Channel Binding object.
type KafkaChannelBindingChanges struct {
	*PropertyChanges
	TopicConfigurationChanges *KafkaTopicConfigurationChanges `json:"topicConfiguration,omitempty" yaml:"topicConfiguration,omitempty"`
	ExtensionChanges          *ExtensionChanges               `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Kafka Channel Binding objects.
func (k *KafkaChannelBindingChanges) GetAllChanges() []*Change {
	if k == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, k.Changes...)
	if k.TopicConfigurationChanges != nil {
		changes = append(changes, k.TopicConfigurationChanges.GetAllChanges()...)
	}
	if k.ExtensionChanges != nil {
		changes = append(changes, k.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (k *KafkaChannelBindingChanges) TotalChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalChanges()
	if k.TopicConfigurationChanges != nil {
		c += k.TopicConfigurationChanges.TotalChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (k *KafkaChannelBindingChanges) TotalBreakingChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalBreakingChanges()
	if k.TopicConfigurationChanges != nil {
		c += k.TopicConfigurationChanges.TotalBreakingChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareKafkaChannelBinding compares two AsyncAPI Kafka Channel Binding objects and
// returns a pointer to KafkaChannelBindingChanges, or nil if nothing changed.
func CompareKafkaChannelBinding(l, r *lowasync.KafkaChannelBinding, configs ...*BreakingRulesConfig) *KafkaChannelBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompKafkaChannelBinding, PropTopic,
			l.Topic.ValueNode, r.Topic.ValueNode, lowasync.TopicLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaChannelBinding, PropPartitions,
			l.Partitions.ValueNode, r.Partitions.ValueNode, lowasync.PartitionsLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaChannelBinding, PropReplicas,
			l.Replicas.ValueNode, r.Replicas.ValueNode, lowasync.ReplicasLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaChannelBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	kb := new(KafkaChannelBindingChanges)

	compareNestedObject(l.TopicConfiguration, r.TopicConfiguration, lowasync.TopicConfigurationLabel,
		CompKafkaChannelBinding, PropTopicConfiguration, &changes,
		configuredNestedCompare(config, CompareKafkaTopicConfiguration), &kb.TopicConfigurationChanges, config)

	kb.ExtensionChanges = CheckExtensions(l, r)
	kb.PropertyChanges = NewPropertyChanges(changes)
	if kb.TotalChanges() <= 0 {
		return nil
	}
	return kb
}

// KafkaTopicConfigurationChanges represents changes made to an AsyncAPI Kafka Topic
// Configuration object.
type KafkaTopicConfigurationChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Kafka Topic Configuration objects.
func (k *KafkaTopicConfigurationChanges) GetAllChanges() []*Change {
	if k == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, k.Changes...)
	if k.ExtensionChanges != nil {
		changes = append(changes, k.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (k *KafkaTopicConfigurationChanges) TotalChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalChanges()
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (k *KafkaTopicConfigurationChanges) TotalBreakingChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalBreakingChanges()
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareKafkaTopicConfiguration compares two AsyncAPI Kafka Topic Configuration objects
// and returns a pointer to KafkaTopicConfigurationChanges, or nil if nothing changed.
func CompareKafkaTopicConfiguration(l, r *lowasync.KafkaTopicConfiguration, configs ...*BreakingRulesConfig) *KafkaTopicConfigurationChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompKafkaTopicConfiguration, PropRetentionMs,
			l.RetentionMs.ValueNode, r.RetentionMs.ValueNode, lowasync.RetentionMsLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropRetentionBytes,
			l.RetentionBytes.ValueNode, r.RetentionBytes.ValueNode, lowasync.RetentionBytesLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropDeleteRetentionMs,
			l.DeleteRetentionMs.ValueNode, r.DeleteRetentionMs.ValueNode, lowasync.DeleteRetentionMsLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropMaxMessageBytes,
			l.MaxMessageBytes.ValueNode, r.MaxMessageBytes.ValueNode, lowasync.MaxMessageBytesLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropConfluentKeySchemaValidation,
			l.ConfluentKeySchemaValidation.ValueNode, r.ConfluentKeySchemaValidation.ValueNode,
			lowasync.ConfluentKeySchemaValidationLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropConfluentKeySubjectNameStrategy,
			l.ConfluentKeySubjectNameStrategy.ValueNode, r.ConfluentKeySubjectNameStrategy.ValueNode,
			lowasync.ConfluentKeySubjectNameStrategyLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropConfluentValueSchemaValidation,
			l.ConfluentValueSchemaValidation.ValueNode, r.ConfluentValueSchemaValidation.ValueNode,
			lowasync.ConfluentValueSchemaValidationLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaTopicConfiguration, PropConfluentValueSubjectNameStrategy,
			l.ConfluentValueSubjectNameStrategy.ValueNode, r.ConfluentValueSubjectNameStrategy.ValueNode,
			lowasync.ConfluentValueSubjectNameStrategyLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ExtractStringValueSliceChangesWithRules(l.CleanupPolicy.Value, r.CleanupPolicy.Value,
		&changes, lowasync.CleanupPolicyLabel, CompKafkaTopicConfiguration, PropCleanupPolicy, config)

	kc := new(KafkaTopicConfigurationChanges)
	kc.ExtensionChanges = CheckExtensions(l, r)
	kc.PropertyChanges = NewPropertyChanges(changes)
	if kc.TotalChanges() <= 0 {
		return nil
	}
	return kc
}

// KafkaOperationBindingChanges represents changes made to an AsyncAPI Kafka Operation Binding object.
type KafkaOperationBindingChanges struct {
	*PropertyChanges
	GroupIDChanges   *wcmodel.SchemaChanges `json:"groupId,omitempty" yaml:"groupId,omitempty"`
	ClientIDChanges  *wcmodel.SchemaChanges `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ExtensionChanges *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Kafka Operation Binding objects.
func (k *KafkaOperationBindingChanges) GetAllChanges() []*Change {
	if k == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, k.Changes...)
	if k.GroupIDChanges != nil {
		changes = append(changes, k.GroupIDChanges.GetAllChanges()...)
	}
	if k.ClientIDChanges != nil {
		changes = append(changes, k.ClientIDChanges.GetAllChanges()...)
	}
	if k.ExtensionChanges != nil {
		changes = append(changes, k.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (k *KafkaOperationBindingChanges) TotalChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalChanges()
	if k.GroupIDChanges != nil {
		c += k.GroupIDChanges.TotalChanges()
	}
	if k.ClientIDChanges != nil {
		c += k.ClientIDChanges.TotalChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (k *KafkaOperationBindingChanges) TotalBreakingChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalBreakingChanges()
	if k.GroupIDChanges != nil {
		c += k.GroupIDChanges.TotalBreakingChanges()
	}
	if k.ClientIDChanges != nil {
		c += k.ClientIDChanges.TotalBreakingChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareKafkaOperationBinding compares two AsyncAPI Kafka Operation Binding objects and
// returns a pointer to KafkaOperationBindingChanges, or nil if nothing changed.
//
// The groupId and clientId fields are schemas, so their comparison is delegated to
// libopenapi's what-changed schema comparator.
func CompareKafkaOperationBinding(l, r *lowasync.KafkaOperationBinding, configs ...*BreakingRulesConfig) *KafkaOperationBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompKafkaOperationBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	kb := new(KafkaOperationBindingChanges)

	compareNestedObject(l.GroupID, r.GroupID, lowasync.GroupIDLabel,
		CompKafkaOperationBinding, PropGroupID, &changes, wcmodel.CompareSchemas, &kb.GroupIDChanges, config)
	compareNestedObject(l.ClientID, r.ClientID, lowasync.ClientIDLabel,
		CompKafkaOperationBinding, PropClientID, &changes, wcmodel.CompareSchemas, &kb.ClientIDChanges, config)

	kb.ExtensionChanges = CheckExtensions(l, r)
	kb.PropertyChanges = NewPropertyChanges(changes)
	if kb.TotalChanges() <= 0 {
		return nil
	}
	return kb
}

// KafkaMessageBindingChanges represents changes made to an AsyncAPI Kafka Message Binding object.
type KafkaMessageBindingChanges struct {
	*PropertyChanges
	KeyChanges       *wcmodel.SchemaChanges `json:"key,omitempty" yaml:"key,omitempty"`
	ExtensionChanges *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Kafka Message Binding objects.
func (k *KafkaMessageBindingChanges) GetAllChanges() []*Change {
	if k == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, k.Changes...)
	if k.KeyChanges != nil {
		changes = append(changes, k.KeyChanges.GetAllChanges()...)
	}
	if k.ExtensionChanges != nil {
		changes = append(changes, k.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (k *KafkaMessageBindingChanges) TotalChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalChanges()
	if k.KeyChanges != nil {
		c += k.KeyChanges.TotalChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (k *KafkaMessageBindingChanges) TotalBreakingChanges() int {
	if k == nil {
		return 0
	}
	c := k.PropertyChanges.TotalBreakingChanges()
	if k.KeyChanges != nil {
		c += k.KeyChanges.TotalBreakingChanges()
	}
	if k.ExtensionChanges != nil {
		c += k.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareKafkaMessageBinding compares two AsyncAPI Kafka Message Binding objects and
// returns a pointer to KafkaMessageBindingChanges, or nil if nothing changed.
//
// The key field is a schema, so its comparison is delegated to libopenapi's what-changed
// schema comparator.
func CompareKafkaMessageBinding(l, r *lowasync.KafkaMessageBinding, configs ...*BreakingRulesConfig) *KafkaMessageBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompKafkaMessageBinding, PropSchemaIDLocation,
			l.SchemaIDLocation.ValueNode, r.SchemaIDLocation.ValueNode, lowasync.SchemaIDLocationLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaMessageBinding, PropSchemaIDPayloadEncoding,
			l.SchemaIDPayloadEncoding.ValueNode, r.SchemaIDPayloadEncoding.ValueNode,
			lowasync.SchemaIDPayloadEncodingLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaMessageBinding, PropSchemaLookupStrategy,
			l.SchemaLookupStrategy.ValueNode, r.SchemaLookupStrategy.ValueNode,
			lowasync.SchemaLookupStrategyLabel, &changes, l, r, config),
		NewPropertyCheck(CompKafkaMessageBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	kb := new(KafkaMessageBindingChanges)

	compareNestedObject(l.Key, r.Key, lowasync.KeyLabel,
		CompKafkaMessageBinding, PropKey, &changes, wcmodel.CompareSchemas, &kb.KeyChanges, config)

	kb.ExtensionChanges = CheckExtensions(l, r)
	kb.PropertyChanges = NewPropertyChanges(changes)
	if kb.TotalChanges() <= 0 {
		return nil
	}
	return kb
}

// WebSocket Bindings

// WebSocketChannelBindingChanges represents changes made to an AsyncAPI WebSocket Channel
// Binding object.
type WebSocketChannelBindingChanges struct {
	*PropertyChanges
	QueryChanges     *wcmodel.SchemaChanges `json:"query,omitempty" yaml:"query,omitempty"`
	HeadersChanges   *wcmodel.SchemaChanges `json:"headers,omitempty" yaml:"headers,omitempty"`
	ExtensionChanges *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between WebSocket Channel Binding objects.
func (w *WebSocketChannelBindingChanges) GetAllChanges() []*Change {
	if w == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, w.Changes...)
	if w.QueryChanges != nil {
		changes = append(changes, w.QueryChanges.GetAllChanges()...)
	}
	if w.HeadersChanges != nil {
		changes = append(changes, w.HeadersChanges.GetAllChanges()...)
	}
	if w.ExtensionChanges != nil {
		changes = append(changes, w.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (w *WebSocketChannelBindingChanges) TotalChanges() int {
	if w == nil {
		return 0
	}
	c := w.PropertyChanges.TotalChanges()
	if w.QueryChanges != nil {
		c += w.QueryChanges.TotalChanges()
	}
	if w.HeadersChanges != nil {
		c += w.HeadersChanges.TotalChanges()
	}
	if w.ExtensionChanges != nil {
		c += w.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (w *WebSocketChannelBindingChanges) TotalBreakingChanges() int {
	if w == nil {
		return 0
	}
	c := w.PropertyChanges.TotalBreakingChanges()
	if w.QueryChanges != nil {
		c += w.QueryChanges.TotalBreakingChanges()
	}
	if w.HeadersChanges != nil {
		c += w.HeadersChanges.TotalBreakingChanges()
	}
	if w.ExtensionChanges != nil {
		c += w.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareWebSocketChannelBinding compares two AsyncAPI WebSocket Channel Binding objects
// and returns a pointer to WebSocketChannelBindingChanges, or nil if nothing changed.
//
// The query and headers fields are schemas, so their comparison is delegated to
// libopenapi's what-changed schema comparator.
func CompareWebSocketChannelBinding(l, r *lowasync.WebSocketChannelBinding, configs ...*BreakingRulesConfig) *WebSocketChannelBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompWebSocketChannelBinding, PropMethod,
			l.Method.ValueNode, r.Method.ValueNode, lowasync.MethodLabel, &changes, l, r, config),
		NewPropertyCheck(CompWebSocketChannelBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	wb := new(WebSocketChannelBindingChanges)

	compareNestedObject(l.Query, r.Query, lowasync.QueryLabel,
		CompWebSocketChannelBinding, PropQuery, &changes, wcmodel.CompareSchemas, &wb.QueryChanges, config)
	compareNestedObject(l.Headers, r.Headers, lowasync.HeadersLabel,
		CompWebSocketChannelBinding, PropHeaders, &changes, wcmodel.CompareSchemas, &wb.HeadersChanges, config)

	wb.ExtensionChanges = CheckExtensions(l, r)
	wb.PropertyChanges = NewPropertyChanges(changes)
	if wb.TotalChanges() <= 0 {
		return nil
	}
	return wb
}

// AMQP Bindings

// AMQPChannelBindingChanges represents changes made to an AsyncAPI AMQP Channel Binding object.
type AMQPChannelBindingChanges struct {
	*PropertyChanges
	ExchangeChanges  *AMQPExchangeChanges `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	QueueChanges     *AMQPQueueChanges    `json:"queue,omitempty" yaml:"queue,omitempty"`
	ExtensionChanges *ExtensionChanges    `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AMQP Channel Binding objects.
func (a *AMQPChannelBindingChanges) GetAllChanges() []*Change {
	if a == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, a.Changes...)
	if a.ExchangeChanges != nil {
		changes = append(changes, a.ExchangeChanges.GetAllChanges()...)
	}
	if a.QueueChanges != nil {
		changes = append(changes, a.QueueChanges.GetAllChanges()...)
	}
	if a.ExtensionChanges != nil {
		changes = append(changes, a.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (a *AMQPChannelBindingChanges) TotalChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalChanges()
	if a.ExchangeChanges != nil {
		c += a.ExchangeChanges.TotalChanges()
	}
	if a.QueueChanges != nil {
		c += a.QueueChanges.TotalChanges()
	}
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (a *AMQPChannelBindingChanges) TotalBreakingChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalBreakingChanges()
	if a.ExchangeChanges != nil {
		c += a.ExchangeChanges.TotalBreakingChanges()
	}
	if a.QueueChanges != nil {
		c += a.QueueChanges.TotalBreakingChanges()
	}
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareAMQPChannelBinding compares two AsyncAPI AMQP Channel Binding objects and returns
// a pointer to AMQPChannelBindingChanges, or nil if nothing changed.
func CompareAMQPChannelBinding(l, r *lowasync.AMQPChannelBinding, configs ...*BreakingRulesConfig) *AMQPChannelBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompAMQPChannelBinding, PropIs,
			l.Is.ValueNode, r.Is.ValueNode, lowasync.IsLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPChannelBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ab := new(AMQPChannelBindingChanges)

	compareNestedObject(l.Exchange, r.Exchange, lowasync.ExchangeLabel,
		CompAMQPChannelBinding, PropExchange, &changes, configuredNestedCompare(config, CompareAMQPExchange), &ab.ExchangeChanges, config)
	compareNestedObject(l.Queue, r.Queue, lowasync.QueueLabel,
		CompAMQPChannelBinding, PropQueue, &changes, configuredNestedCompare(config, CompareAMQPQueue), &ab.QueueChanges, config)

	ab.ExtensionChanges = CheckExtensions(l, r)
	ab.PropertyChanges = NewPropertyChanges(changes)
	if ab.TotalChanges() <= 0 {
		return nil
	}
	return ab
}

// AMQPExchangeChanges represents changes made to an AsyncAPI AMQP Exchange object.
type AMQPExchangeChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AMQP Exchange objects.
func (a *AMQPExchangeChanges) GetAllChanges() []*Change {
	if a == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, a.Changes...)
	if a.ExtensionChanges != nil {
		changes = append(changes, a.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (a *AMQPExchangeChanges) TotalChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (a *AMQPExchangeChanges) TotalBreakingChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalBreakingChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareAMQPExchange compares two AsyncAPI AMQP Exchange objects and returns a pointer to
// AMQPExchangeChanges, or nil if nothing changed.
func CompareAMQPExchange(l, r *lowasync.AMQPExchange, configs ...*BreakingRulesConfig) *AMQPExchangeChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompAMQPExchange, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPExchange, PropType,
			l.Type.ValueNode, r.Type.ValueNode, lowasync.TypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPExchange, PropDurable,
			l.Durable.ValueNode, r.Durable.ValueNode, lowasync.DurableLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPExchange, PropAutoDelete,
			l.AutoDelete.ValueNode, r.AutoDelete.ValueNode, lowasync.AutoDeleteLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPExchange, PropVHost,
			l.VHost.ValueNode, r.VHost.ValueNode, lowasync.VHostLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ae := new(AMQPExchangeChanges)
	ae.ExtensionChanges = CheckExtensions(l, r)
	ae.PropertyChanges = NewPropertyChanges(changes)
	if ae.TotalChanges() <= 0 {
		return nil
	}
	return ae
}

// AMQPQueueChanges represents changes made to an AsyncAPI AMQP Queue object.
type AMQPQueueChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AMQP Queue objects.
func (a *AMQPQueueChanges) GetAllChanges() []*Change {
	if a == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, a.Changes...)
	if a.ExtensionChanges != nil {
		changes = append(changes, a.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (a *AMQPQueueChanges) TotalChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (a *AMQPQueueChanges) TotalBreakingChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalBreakingChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareAMQPQueue compares two AsyncAPI AMQP Queue objects and returns a pointer to
// AMQPQueueChanges, or nil if nothing changed.
func CompareAMQPQueue(l, r *lowasync.AMQPQueue, configs ...*BreakingRulesConfig) *AMQPQueueChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompAMQPQueue, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPQueue, PropDurable,
			l.Durable.ValueNode, r.Durable.ValueNode, lowasync.DurableLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPQueue, PropExclusive,
			l.Exclusive.ValueNode, r.Exclusive.ValueNode, lowasync.ExclusiveLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPQueue, PropAutoDelete,
			l.AutoDelete.ValueNode, r.AutoDelete.ValueNode, lowasync.AutoDeleteLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPQueue, PropVHost,
			l.VHost.ValueNode, r.VHost.ValueNode, lowasync.VHostLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	aq := new(AMQPQueueChanges)
	aq.ExtensionChanges = CheckExtensions(l, r)
	aq.PropertyChanges = NewPropertyChanges(changes)
	if aq.TotalChanges() <= 0 {
		return nil
	}
	return aq
}

// AMQPOperationBindingChanges represents changes made to an AsyncAPI AMQP Operation Binding object.
type AMQPOperationBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AMQP Operation Binding objects.
func (a *AMQPOperationBindingChanges) GetAllChanges() []*Change {
	if a == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, a.Changes...)
	if a.ExtensionChanges != nil {
		changes = append(changes, a.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (a *AMQPOperationBindingChanges) TotalChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (a *AMQPOperationBindingChanges) TotalBreakingChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalBreakingChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareAMQPOperationBinding compares two AsyncAPI AMQP Operation Binding objects and
// returns a pointer to AMQPOperationBindingChanges, or nil if nothing changed.
func CompareAMQPOperationBinding(l, r *lowasync.AMQPOperationBinding, configs ...*BreakingRulesConfig) *AMQPOperationBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompAMQPOperationBinding, PropExpiration,
			l.Expiration.ValueNode, r.Expiration.ValueNode, lowasync.ExpirationLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropUserID,
			l.UserID.ValueNode, r.UserID.ValueNode, lowasync.UserIDLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropPriority,
			l.Priority.ValueNode, r.Priority.ValueNode, lowasync.PriorityLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropDeliveryMode,
			l.DeliveryMode.ValueNode, r.DeliveryMode.ValueNode, lowasync.DeliveryModeLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropMandatory,
			l.Mandatory.ValueNode, r.Mandatory.ValueNode, lowasync.MandatoryLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropTimestamp,
			l.Timestamp.ValueNode, r.Timestamp.ValueNode, lowasync.TimestampLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropAck,
			l.Ack.ValueNode, r.Ack.ValueNode, lowasync.AckLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPOperationBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ExtractStringValueSliceChangesWithRules(l.CC.Value, r.CC.Value,
		&changes, lowasync.CCLabel, CompAMQPOperationBinding, PropCC, config)
	ExtractStringValueSliceChangesWithRules(l.BCC.Value, r.BCC.Value,
		&changes, lowasync.BCCLabel, CompAMQPOperationBinding, PropBCC, config)

	ab := new(AMQPOperationBindingChanges)
	ab.ExtensionChanges = CheckExtensions(l, r)
	ab.PropertyChanges = NewPropertyChanges(changes)
	if ab.TotalChanges() <= 0 {
		return nil
	}
	return ab
}

// AMQPMessageBindingChanges represents changes made to an AsyncAPI AMQP Message Binding object.
type AMQPMessageBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AMQP Message Binding objects.
func (a *AMQPMessageBindingChanges) GetAllChanges() []*Change {
	if a == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, a.Changes...)
	if a.ExtensionChanges != nil {
		changes = append(changes, a.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (a *AMQPMessageBindingChanges) TotalChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (a *AMQPMessageBindingChanges) TotalBreakingChanges() int {
	if a == nil {
		return 0
	}
	c := a.PropertyChanges.TotalBreakingChanges()
	if a.ExtensionChanges != nil {
		c += a.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareAMQPMessageBinding compares two AsyncAPI AMQP Message Binding objects and returns
// a pointer to AMQPMessageBindingChanges, or nil if nothing changed.
func CompareAMQPMessageBinding(l, r *lowasync.AMQPMessageBinding, configs ...*BreakingRulesConfig) *AMQPMessageBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompAMQPMessageBinding, PropContentEncoding,
			l.ContentEncoding.ValueNode, r.ContentEncoding.ValueNode, lowasync.ContentEncodingLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPMessageBinding, PropMessageType,
			l.MessageType.ValueNode, r.MessageType.ValueNode, lowasync.MessageTypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompAMQPMessageBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ab := new(AMQPMessageBindingChanges)
	ab.ExtensionChanges = CheckExtensions(l, r)
	ab.PropertyChanges = NewPropertyChanges(changes)
	if ab.TotalChanges() <= 0 {
		return nil
	}
	return ab
}

// MQTT Bindings

// MQTTServerBindingChanges represents changes made to an AsyncAPI MQTT Server Binding object.
type MQTTServerBindingChanges struct {
	*PropertyChanges
	LastWillChanges  *MQTTLastWillChanges `json:"lastWill,omitempty" yaml:"lastWill,omitempty"`
	ExtensionChanges *ExtensionChanges    `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between MQTT Server Binding objects.
func (m *MQTTServerBindingChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.LastWillChanges != nil {
		changes = append(changes, m.LastWillChanges.GetAllChanges()...)
	}
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MQTTServerBindingChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.LastWillChanges != nil {
		c += m.LastWillChanges.TotalChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MQTTServerBindingChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.LastWillChanges != nil {
		c += m.LastWillChanges.TotalBreakingChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMQTTServerBinding compares two AsyncAPI MQTT Server Binding objects and returns a
// pointer to MQTTServerBindingChanges, or nil if nothing changed.
func CompareMQTTServerBinding(l, r *lowasync.MQTTServerBinding, configs ...*BreakingRulesConfig) *MQTTServerBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMQTTServerBinding, PropClientID,
			l.ClientID.ValueNode, r.ClientID.ValueNode, lowasync.ClientIDLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTServerBinding, PropCleanSession,
			l.CleanSession.ValueNode, r.CleanSession.ValueNode, lowasync.CleanSessionLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTServerBinding, PropKeepAlive,
			l.KeepAlive.ValueNode, r.KeepAlive.ValueNode, lowasync.KeepAliveLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTServerBinding, PropSessionExpiryInterval,
			l.SessionExpiryInterval.ValueNode, r.SessionExpiryInterval.ValueNode,
			lowasync.SessionExpiryIntervalLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTServerBinding, PropMaximumPacketSize,
			l.MaximumPacketSize.ValueNode, r.MaximumPacketSize.ValueNode,
			lowasync.MaximumPacketSizeLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTServerBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	mb := new(MQTTServerBindingChanges)

	compareNestedObject(l.LastWill, r.LastWill, lowasync.LastWillLabel,
		CompMQTTServerBinding, PropLastWill, &changes, configuredNestedCompare(config, CompareMQTTLastWill), &mb.LastWillChanges, config)

	mb.ExtensionChanges = CheckExtensions(l, r)
	mb.PropertyChanges = NewPropertyChanges(changes)
	if mb.TotalChanges() <= 0 {
		return nil
	}
	return mb
}

// MQTTLastWillChanges represents changes made to an AsyncAPI MQTT Last Will object.
type MQTTLastWillChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between MQTT Last Will objects.
func (m *MQTTLastWillChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MQTTLastWillChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MQTTLastWillChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMQTTLastWill compares two AsyncAPI MQTT Last Will objects and returns a pointer to
// MQTTLastWillChanges, or nil if nothing changed.
func CompareMQTTLastWill(l, r *lowasync.MQTTLastWill, configs ...*BreakingRulesConfig) *MQTTLastWillChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMQTTLastWill, PropTopic,
			l.Topic.ValueNode, r.Topic.ValueNode, lowasync.TopicLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTLastWill, PropQoS,
			l.QoS.ValueNode, r.QoS.ValueNode, lowasync.QoSLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTLastWill, PropMessage,
			l.Message.ValueNode, r.Message.ValueNode, lowasync.MessageLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTLastWill, PropRetain,
			l.Retain.ValueNode, r.Retain.ValueNode, lowasync.RetainLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	mw := new(MQTTLastWillChanges)
	mw.ExtensionChanges = CheckExtensions(l, r)
	mw.PropertyChanges = NewPropertyChanges(changes)
	if mw.TotalChanges() <= 0 {
		return nil
	}
	return mw
}

// MQTTOperationBindingChanges represents changes made to an AsyncAPI MQTT Operation Binding object.
type MQTTOperationBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between MQTT Operation Binding objects.
func (m *MQTTOperationBindingChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MQTTOperationBindingChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MQTTOperationBindingChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMQTTOperationBinding compares two AsyncAPI MQTT Operation Binding objects and
// returns a pointer to MQTTOperationBindingChanges, or nil if nothing changed.
func CompareMQTTOperationBinding(l, r *lowasync.MQTTOperationBinding, configs ...*BreakingRulesConfig) *MQTTOperationBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMQTTOperationBinding, PropQoS,
			l.QoS.ValueNode, r.QoS.ValueNode, lowasync.QoSLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTOperationBinding, PropRetain,
			l.Retain.ValueNode, r.Retain.ValueNode, lowasync.RetainLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTOperationBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	mb := new(MQTTOperationBindingChanges)
	mb.ExtensionChanges = CheckExtensions(l, r)
	mb.PropertyChanges = NewPropertyChanges(changes)
	if mb.TotalChanges() <= 0 {
		return nil
	}
	return mb
}

// MQTTMessageBindingChanges represents changes made to an AsyncAPI MQTT Message Binding object.
type MQTTMessageBindingChanges struct {
	*PropertyChanges
	CorrelationDataChanges *wcmodel.SchemaChanges `json:"correlationData,omitempty" yaml:"correlationData,omitempty"`
	ExtensionChanges       *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between MQTT Message Binding objects.
func (m *MQTTMessageBindingChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.CorrelationDataChanges != nil {
		changes = append(changes, m.CorrelationDataChanges.GetAllChanges()...)
	}
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MQTTMessageBindingChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.CorrelationDataChanges != nil {
		c += m.CorrelationDataChanges.TotalChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MQTTMessageBindingChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.CorrelationDataChanges != nil {
		c += m.CorrelationDataChanges.TotalBreakingChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMQTTMessageBinding compares two AsyncAPI MQTT Message Binding objects and returns
// a pointer to MQTTMessageBindingChanges, or nil if nothing changed.
//
// The correlationData field is a schema, so its comparison is delegated to libopenapi's
// what-changed schema comparator.
func CompareMQTTMessageBinding(l, r *lowasync.MQTTMessageBinding, configs ...*BreakingRulesConfig) *MQTTMessageBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMQTTMessageBinding, PropPayloadFormatIndicator,
			l.PayloadFormatIndicator.ValueNode, r.PayloadFormatIndicator.ValueNode,
			lowasync.PayloadFormatIndicatorLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTMessageBinding, PropContentType,
			l.ContentType.ValueNode, r.ContentType.ValueNode, lowasync.ContentTypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTMessageBinding, PropResponseTopic,
			l.ResponseTopic.ValueNode, r.ResponseTopic.ValueNode, lowasync.ResponseTopicLabel, &changes, l, r, config),
		NewPropertyCheck(CompMQTTMessageBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	mb := new(MQTTMessageBindingChanges)

	compareNestedObject(l.CorrelationData, r.CorrelationData, lowasync.CorrelationDataLabel,
		CompMQTTMessageBinding, PropCorrelationData, &changes, wcmodel.CompareSchemas, &mb.CorrelationDataChanges, config)

	mb.ExtensionChanges = CheckExtensions(l, r)
	mb.PropertyChanges = NewPropertyChanges(changes)
	if mb.TotalChanges() <= 0 {
		return nil
	}
	return mb
}

// SQS Bindings

// SQSServerBindingChanges represents changes made to an AsyncAPI SQS Server Binding object.
type SQSServerBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Server Binding objects.
func (s *SQSServerBindingChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSServerBindingChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSServerBindingChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSServerBinding compares two AsyncAPI SQS Server Binding objects and returns a
// pointer to SQSServerBindingChanges, or nil if nothing changed. The binding carries no
// fields beyond extensions.
func CompareSQSServerBinding(l, r *lowasync.SQSServerBinding, configs ...*BreakingRulesConfig) *SQSServerBindingChanges {
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	sb := new(SQSServerBindingChanges)
	sb.ExtensionChanges = CheckExtensions(l, r)
	sb.PropertyChanges = NewPropertyChanges(changes)
	if sb.TotalChanges() <= 0 {
		return nil
	}
	return sb
}

// SQSChannelBindingChanges represents changes made to an AsyncAPI SQS Channel Binding object.
type SQSChannelBindingChanges struct {
	*PropertyChanges
	QueueChanges           *SQSQueueChanges  `json:"queue,omitempty" yaml:"queue,omitempty"`
	DeadLetterQueueChanges *SQSQueueChanges  `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	ExtensionChanges       *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Channel Binding objects.
func (s *SQSChannelBindingChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.QueueChanges != nil {
		changes = append(changes, s.QueueChanges.GetAllChanges()...)
	}
	if s.DeadLetterQueueChanges != nil {
		changes = append(changes, s.DeadLetterQueueChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSChannelBindingChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.QueueChanges != nil {
		c += s.QueueChanges.TotalChanges()
	}
	if s.DeadLetterQueueChanges != nil {
		c += s.DeadLetterQueueChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSChannelBindingChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.QueueChanges != nil {
		c += s.QueueChanges.TotalBreakingChanges()
	}
	if s.DeadLetterQueueChanges != nil {
		c += s.DeadLetterQueueChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSChannelBinding compares two AsyncAPI SQS Channel Binding objects and returns a
// pointer to SQSChannelBindingChanges, or nil if nothing changed.
func CompareSQSChannelBinding(l, r *lowasync.SQSChannelBinding, configs ...*BreakingRulesConfig) *SQSChannelBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSChannelBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	sb := new(SQSChannelBindingChanges)

	compareNestedObject(l.Queue, r.Queue, lowasync.QueueLabel,
		CompSQSChannelBinding, PropQueue, &changes, configuredNestedCompare(config, CompareSQSQueue), &sb.QueueChanges, config)
	compareNestedObject(l.DeadLetterQueue, r.DeadLetterQueue, lowasync.DeadLetterQueueLabel,
		CompSQSChannelBinding, PropDeadLetterQueue, &changes, configuredNestedCompare(config, CompareSQSQueue), &sb.DeadLetterQueueChanges, config)

	sb.ExtensionChanges = CheckExtensions(l, r)
	sb.PropertyChanges = NewPropertyChanges(changes)
	if sb.TotalChanges() <= 0 {
		return nil
	}
	return sb
}

// SQSOperationBindingChanges represents changes made to an AsyncAPI SQS Operation Binding object.
type SQSOperationBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Operation Binding objects.
func (s *SQSOperationBindingChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSOperationBindingChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSOperationBindingChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSOperationBinding compares two AsyncAPI SQS Operation Binding objects and
// returns a pointer to SQSOperationBindingChanges, or nil if nothing changed.
//
// Queues have no natural identity key, so they are compared by hash; an in-place edit
// reports as one removal plus one addition.
func CompareSQSOperationBinding(l, r *lowasync.SQSOperationBinding, configs ...*BreakingRulesConfig) *SQSOperationBindingChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSOperationBinding, PropBindingVersion,
			l.BindingVersion.ValueNode, r.BindingVersion.ValueNode, lowasync.BindingVersionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	compareUnkeyedSlices(l.Queues.Value, r.Queues.Value, lowasync.QueuesLabel,
		CompSQSOperationBinding, PropQueues, &changes, config)

	sb := new(SQSOperationBindingChanges)
	sb.ExtensionChanges = CheckExtensions(l, r)
	sb.PropertyChanges = NewPropertyChanges(changes)
	if sb.TotalChanges() <= 0 {
		return nil
	}
	return sb
}

// SQSMessageBindingChanges represents changes made to an AsyncAPI SQS Message Binding object.
type SQSMessageBindingChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Message Binding objects.
func (s *SQSMessageBindingChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSMessageBindingChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSMessageBindingChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSMessageBinding compares two AsyncAPI SQS Message Binding objects and returns a
// pointer to SQSMessageBindingChanges, or nil if nothing changed. The binding carries no
// fields beyond extensions.
func CompareSQSMessageBinding(l, r *lowasync.SQSMessageBinding, configs ...*BreakingRulesConfig) *SQSMessageBindingChanges {
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	sb := new(SQSMessageBindingChanges)
	sb.ExtensionChanges = CheckExtensions(l, r)
	sb.PropertyChanges = NewPropertyChanges(changes)
	if sb.TotalChanges() <= 0 {
		return nil
	}
	return sb
}

// SQSQueueChanges represents changes made to an AsyncAPI SQS Queue object.
type SQSQueueChanges struct {
	*PropertyChanges
	RedrivePolicyChanges *SQSRedrivePolicyChanges `json:"redrivePolicy,omitempty" yaml:"redrivePolicy,omitempty"`
	PolicyChanges        *SQSPolicyChanges        `json:"policy,omitempty" yaml:"policy,omitempty"`
	ExtensionChanges     *ExtensionChanges        `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Queue objects.
func (s *SQSQueueChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.RedrivePolicyChanges != nil {
		changes = append(changes, s.RedrivePolicyChanges.GetAllChanges()...)
	}
	if s.PolicyChanges != nil {
		changes = append(changes, s.PolicyChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSQueueChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.RedrivePolicyChanges != nil {
		c += s.RedrivePolicyChanges.TotalChanges()
	}
	if s.PolicyChanges != nil {
		c += s.PolicyChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSQueueChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.RedrivePolicyChanges != nil {
		c += s.RedrivePolicyChanges.TotalBreakingChanges()
	}
	if s.PolicyChanges != nil {
		c += s.PolicyChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSQueue compares two AsyncAPI SQS Queue objects and returns a pointer to
// SQSQueueChanges, or nil if nothing changed.
//
// Tags are a free-form map, so they are compared structurally as a raw yaml node.
func CompareSQSQueue(l, r *lowasync.SQSQueue, configs ...*BreakingRulesConfig) *SQSQueueChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSQueue, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropARN,
			l.ARN.ValueNode, r.ARN.ValueNode, lowasync.ARNLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropFifoQueue,
			l.FifoQueue.ValueNode, r.FifoQueue.ValueNode, lowasync.FifoQueueLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropDeduplicationScope,
			l.DeduplicationScope.ValueNode, r.DeduplicationScope.ValueNode,
			lowasync.DeduplicationScopeLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropFifoThroughputLimit,
			l.FifoThroughputLimit.ValueNode, r.FifoThroughputLimit.ValueNode,
			lowasync.FifoThroughputLimitLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropDeliveryDelay,
			l.DeliveryDelay.ValueNode, r.DeliveryDelay.ValueNode, lowasync.DeliveryDelayLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropVisibilityTimeout,
			l.VisibilityTimeout.ValueNode, r.VisibilityTimeout.ValueNode,
			lowasync.VisibilityTimeoutLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropReceiveMessageWaitTime,
			l.ReceiveMessageWaitTime.ValueNode, r.ReceiveMessageWaitTime.ValueNode,
			lowasync.ReceiveMessageWaitTimeLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSQueue, PropMessageRetentionPeriod,
			l.MessageRetentionPeriod.ValueNode, r.MessageRetentionPeriod.ValueNode,
			lowasync.MessageRetentionPeriodLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	sq := new(SQSQueueChanges)

	compareNestedObject(l.RedrivePolicy, r.RedrivePolicy, lowasync.RedrivePolicyLabel,
		CompSQSQueue, PropRedrivePolicy, &changes, configuredNestedCompare(config, CompareSQSRedrivePolicy), &sq.RedrivePolicyChanges, config)
	compareNestedObject(l.Policy, r.Policy, lowasync.PolicyLabel,
		CompSQSQueue, PropPolicy, &changes, configuredNestedCompare(config, CompareSQSPolicy), &sq.PolicyChanges, config)

	compareRawNode(l.Tags, r.Tags, lowasync.TagsLabel, CompSQSQueue, PropTags, &changes, config)

	sq.ExtensionChanges = CheckExtensions(l, r)
	sq.PropertyChanges = NewPropertyChanges(changes)
	if sq.TotalChanges() <= 0 {
		return nil
	}
	return sq
}

// SQSIdentifierChanges represents changes made to an AsyncAPI SQS Identifier object.
type SQSIdentifierChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Identifier objects.
func (s *SQSIdentifierChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSIdentifierChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSIdentifierChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSIdentifier compares two AsyncAPI SQS Identifier objects and returns a pointer
// to SQSIdentifierChanges, or nil if nothing changed.
func CompareSQSIdentifier(l, r *lowasync.SQSIdentifier, configs ...*BreakingRulesConfig) *SQSIdentifierChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSIdentifier, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSIdentifier, PropARN,
			l.ARN.ValueNode, r.ARN.ValueNode, lowasync.ARNLabel, &changes, l, r, config),
		NewPropertyCheck(CompSQSIdentifier, PropFifoQueue,
			l.FifoQueue.ValueNode, r.FifoQueue.ValueNode, lowasync.FifoQueueLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	si := new(SQSIdentifierChanges)
	si.ExtensionChanges = CheckExtensions(l, r)
	si.PropertyChanges = NewPropertyChanges(changes)
	if si.TotalChanges() <= 0 {
		return nil
	}
	return si
}

// SQSRedrivePolicyChanges represents changes made to an AsyncAPI SQS Redrive Policy object.
type SQSRedrivePolicyChanges struct {
	*PropertyChanges
	DeadLetterQueueChanges *SQSIdentifierChanges `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	ExtensionChanges       *ExtensionChanges     `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Redrive Policy objects.
func (s *SQSRedrivePolicyChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.DeadLetterQueueChanges != nil {
		changes = append(changes, s.DeadLetterQueueChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSRedrivePolicyChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.DeadLetterQueueChanges != nil {
		c += s.DeadLetterQueueChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSRedrivePolicyChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.DeadLetterQueueChanges != nil {
		c += s.DeadLetterQueueChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSRedrivePolicy compares two AsyncAPI SQS Redrive Policy objects and returns a
// pointer to SQSRedrivePolicyChanges, or nil if nothing changed.
func CompareSQSRedrivePolicy(l, r *lowasync.SQSRedrivePolicy, configs ...*BreakingRulesConfig) *SQSRedrivePolicyChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSRedrivePolicy, PropMaxReceiveCount,
			l.MaxReceiveCount.ValueNode, r.MaxReceiveCount.ValueNode, lowasync.MaxReceiveCountLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	sr := new(SQSRedrivePolicyChanges)

	compareNestedObject(l.DeadLetterQueue, r.DeadLetterQueue, lowasync.DeadLetterQueueLabel,
		CompSQSRedrivePolicy, PropDeadLetterQueue, &changes, configuredNestedCompare(config, CompareSQSIdentifier), &sr.DeadLetterQueueChanges, config)

	sr.ExtensionChanges = CheckExtensions(l, r)
	sr.PropertyChanges = NewPropertyChanges(changes)
	if sr.TotalChanges() <= 0 {
		return nil
	}
	return sr
}

// SQSPolicyChanges represents changes made to an AsyncAPI SQS Policy object.
type SQSPolicyChanges struct {
	*PropertyChanges
	StatementChanges []*SQSPolicyStatementChanges `json:"statements,omitempty" yaml:"statements,omitempty"`
	ExtensionChanges *ExtensionChanges            `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Policy objects.
func (s *SQSPolicyChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	for _, sc := range s.StatementChanges {
		changes = append(changes, sc.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSPolicyChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	for _, sc := range s.StatementChanges {
		c += sc.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSPolicyChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	for _, sc := range s.StatementChanges {
		c += sc.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSPolicy compares two AsyncAPI SQS Policy objects and returns a pointer to
// SQSPolicyChanges, or nil if nothing changed.
//
// Statements are an ordered list without a natural identity key, so entries are paired
// positionally: statements sharing an index are compared field by field, and surplus
// entries on either side are reported as additions or removals.
func CompareSQSPolicy(l, r *lowasync.SQSPolicy, configs ...*BreakingRulesConfig) *SQSPolicyChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	sp := new(SQSPolicyChanges)
	sp.StatementChanges = compareIndexedSlices(l.Statements.Value, r.Statements.Value,
		lowasync.StatementsLabel, CompSQSPolicy, PropStatements, &changes, configuredNestedCompare(config, CompareSQSPolicyStatement), config)

	sp.ExtensionChanges = CheckExtensions(l, r)
	sp.PropertyChanges = NewPropertyChanges(changes)
	if sp.TotalChanges() <= 0 {
		return nil
	}
	return sp
}

// SQSPolicyStatementChanges represents changes made to an AsyncAPI SQS Policy Statement object.
type SQSPolicyStatementChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between SQS Policy Statement objects.
func (s *SQSPolicyStatementChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SQSPolicyStatementChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SQSPolicyStatementChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSQSPolicyStatement compares two AsyncAPI SQS Policy Statement objects and returns
// a pointer to SQSPolicyStatementChanges, or nil if nothing changed.
//
// Principal, action, resource and condition are free-form IAM values, so they are compared
// structurally as raw yaml nodes.
func CompareSQSPolicyStatement(l, r *lowasync.SQSPolicyStatement, configs ...*BreakingRulesConfig) *SQSPolicyStatementChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSQSPolicyStatement, PropEffect,
			l.Effect.ValueNode, r.Effect.ValueNode, lowasync.EffectLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	compareRawNode(l.Principal, r.Principal, lowasync.PrincipalLabel,
		CompSQSPolicyStatement, PropPrincipal, &changes, config)
	compareRawNode(l.Action, r.Action, lowasync.ActionLabel,
		CompSQSPolicyStatement, PropAction, &changes, config)
	compareRawNode(l.Resource, r.Resource, lowasync.ResourceLabel,
		CompSQSPolicyStatement, PropResource, &changes, config)
	compareRawNode(l.Condition, r.Condition, lowasync.ConditionLabel,
		CompSQSPolicyStatement, PropCondition, &changes, config)

	ss := new(SQSPolicyStatementChanges)
	ss.ExtensionChanges = CheckExtensions(l, r)
	ss.PropertyChanges = NewPropertyChanges(changes)
	if ss.TotalChanges() <= 0 {
		return nil
	}
	return ss
}
