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

// HTTP Server Binding

// HTTPServerBinding represents a high-level HTTP Server Binding.
type HTTPServerBinding struct {
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.HTTPServerBinding
}

// NewHTTPServerBinding creates a new high-level HTTPServerBinding.
func NewHTTPServerBinding(b *lowasync.HTTPServerBinding) *HTTPServerBinding {
	h := new(HTTPServerBinding)
	h.low = b
	if orderedmap.Len(b.Extensions) > 0 {
		h.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return h
}

// GoLow returns the low-level HTTPServerBinding.
func (h *HTTPServerBinding) GoLow() *lowasync.HTTPServerBinding { return h.low }

// GoLowUntyped returns the low-level HTTPServerBinding with no type.
func (h *HTTPServerBinding) GoLowUntyped() any { return h.low }

// HTTP Channel Binding

// HTTPChannelBinding represents a high-level HTTP Channel Binding.
type HTTPChannelBinding struct {
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.HTTPChannelBinding
}

// NewHTTPChannelBinding creates a new high-level HTTPChannelBinding.
func NewHTTPChannelBinding(b *lowasync.HTTPChannelBinding) *HTTPChannelBinding {
	h := new(HTTPChannelBinding)
	h.low = b
	if orderedmap.Len(b.Extensions) > 0 {
		h.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return h
}

// GoLow returns the low-level HTTPChannelBinding.
func (h *HTTPChannelBinding) GoLow() *lowasync.HTTPChannelBinding { return h.low }

// GoLowUntyped returns the low-level HTTPChannelBinding with no type.
func (h *HTTPChannelBinding) GoLowUntyped() any { return h.low }

// HTTP Operation Binding

// HTTPOperationBinding represents a high-level HTTP Operation Binding.
type HTTPOperationBinding struct {
	Method         string                              `json:"method,omitempty" yaml:"method,omitempty"`
	Query          *highbase.SchemaProxy               `json:"query,omitempty" yaml:"query,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.HTTPOperationBinding
}

// NewHTTPOperationBinding creates a new high-level HTTPOperationBinding.
func NewHTTPOperationBinding(b *lowasync.HTTPOperationBinding) *HTTPOperationBinding {
	h := new(HTTPOperationBinding)
	h.low = b
	h.Method = b.Method.Value
	h.BindingVersion = b.BindingVersion.Value
	if !b.Query.IsEmpty() {
		h.Query = highbase.NewSchemaProxy(&b.Query)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		h.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return h
}

// GoLow returns the low-level HTTPOperationBinding.
func (h *HTTPOperationBinding) GoLow() *lowasync.HTTPOperationBinding { return h.low }

// GoLowUntyped returns the low-level HTTPOperationBinding with no type.
func (h *HTTPOperationBinding) GoLowUntyped() any { return h.low }

// HTTP Message Binding

// HTTPMessageBinding represents a high-level HTTP Message Binding.
type HTTPMessageBinding struct {
	Headers        *highbase.SchemaProxy               `json:"headers,omitempty" yaml:"headers,omitempty"`
	StatusCode     int                                 `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.HTTPMessageBinding
}

// NewHTTPMessageBinding creates a new high-level HTTPMessageBinding.
func NewHTTPMessageBinding(b *lowasync.HTTPMessageBinding) *HTTPMessageBinding {
	h := new(HTTPMessageBinding)
	h.low = b
	h.StatusCode = b.StatusCode.Value
	h.BindingVersion = b.BindingVersion.Value
	if !b.Headers.IsEmpty() {
		h.Headers = highbase.NewSchemaProxy(&b.Headers)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		h.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return h
}

// GoLow returns the low-level HTTPMessageBinding.
func (h *HTTPMessageBinding) GoLow() *lowasync.HTTPMessageBinding { return h.low }

// GoLowUntyped returns the low-level HTTPMessageBinding with no type.
func (h *HTTPMessageBinding) GoLowUntyped() any { return h.low }

// SQS Server Binding

// SQSServerBinding represents a high-level SQS Server Binding.
type SQSServerBinding struct {
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.SQSServerBinding
}

// NewSQSServerBinding creates a new high-level SQSServerBinding.
func NewSQSServerBinding(b *lowasync.SQSServerBinding) *SQSServerBinding {
	s := new(SQSServerBinding)
	s.low = b
	if orderedmap.Len(b.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSServerBinding.
func (s *SQSServerBinding) GoLow() *lowasync.SQSServerBinding { return s.low }

// GoLowUntyped returns the low-level SQSServerBinding with no type.
func (s *SQSServerBinding) GoLowUntyped() any { return s.low }

// SQS Channel Binding

// SQSChannelBinding represents a high-level SQS Channel Binding.
type SQSChannelBinding struct {
	Queue           *SQSQueue                           `json:"queue,omitempty" yaml:"queue,omitempty"`
	DeadLetterQueue *SQSQueue                           `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	BindingVersion  string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions      *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low             *lowasync.SQSChannelBinding
}

// NewSQSChannelBinding creates a new high-level SQSChannelBinding.
func NewSQSChannelBinding(b *lowasync.SQSChannelBinding) *SQSChannelBinding {
	s := new(SQSChannelBinding)
	s.low = b
	s.BindingVersion = b.BindingVersion.Value
	if !b.Queue.IsEmpty() {
		s.Queue = NewSQSQueue(b.Queue.Value)
	}
	if !b.DeadLetterQueue.IsEmpty() {
		s.DeadLetterQueue = NewSQSQueue(b.DeadLetterQueue.Value)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSChannelBinding.
func (s *SQSChannelBinding) GoLow() *lowasync.SQSChannelBinding { return s.low }

// GoLowUntyped returns the low-level SQSChannelBinding with no type.
func (s *SQSChannelBinding) GoLowUntyped() any { return s.low }

// SQS Operation Binding

// SQSOperationBinding represents a high-level SQS Operation Binding.
type SQSOperationBinding struct {
	Queues         []*SQSQueue                         `json:"queues,omitempty" yaml:"queues,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.SQSOperationBinding
}

// NewSQSOperationBinding creates a new high-level SQSOperationBinding.
func NewSQSOperationBinding(b *lowasync.SQSOperationBinding) *SQSOperationBinding {
	s := new(SQSOperationBinding)
	s.low = b
	s.BindingVersion = b.BindingVersion.Value
	if b.Queues.Value != nil {
		for _, queue := range b.Queues.Value {
			s.Queues = append(s.Queues, NewSQSQueue(queue.Value))
		}
	}
	if orderedmap.Len(b.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSOperationBinding.
func (s *SQSOperationBinding) GoLow() *lowasync.SQSOperationBinding { return s.low }

// GoLowUntyped returns the low-level SQSOperationBinding with no type.
func (s *SQSOperationBinding) GoLowUntyped() any { return s.low }

// SQS Message Binding

// SQSMessageBinding represents a high-level SQS Message Binding.
type SQSMessageBinding struct {
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.SQSMessageBinding
}

// NewSQSMessageBinding creates a new high-level SQSMessageBinding.
func NewSQSMessageBinding(b *lowasync.SQSMessageBinding) *SQSMessageBinding {
	s := new(SQSMessageBinding)
	s.low = b
	if orderedmap.Len(b.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSMessageBinding.
func (s *SQSMessageBinding) GoLow() *lowasync.SQSMessageBinding { return s.low }

// GoLowUntyped returns the low-level SQSMessageBinding with no type.
func (s *SQSMessageBinding) GoLowUntyped() any { return s.low }

// SQSQueue represents high-level SQS queue configuration.
type SQSQueue struct {
	Name                   string                              `json:"name,omitempty" yaml:"name,omitempty"`
	ARN                    string                              `json:"arn,omitempty" yaml:"arn,omitempty"`
	FifoQueue              bool                                `json:"fifoQueue,omitempty" yaml:"fifoQueue,omitempty"`
	DeduplicationScope     string                              `json:"deduplicationScope,omitempty" yaml:"deduplicationScope,omitempty"`
	FifoThroughputLimit    string                              `json:"fifoThroughputLimit,omitempty" yaml:"fifoThroughputLimit,omitempty"`
	DeliveryDelay          int                                 `json:"deliveryDelay,omitempty" yaml:"deliveryDelay,omitempty"`
	VisibilityTimeout      int                                 `json:"visibilityTimeout,omitempty" yaml:"visibilityTimeout,omitempty"`
	ReceiveMessageWaitTime int                                 `json:"receiveMessageWaitTime,omitempty" yaml:"receiveMessageWaitTime,omitempty"`
	MessageRetentionPeriod int                                 `json:"messageRetentionPeriod,omitempty" yaml:"messageRetentionPeriod,omitempty"`
	RedrivePolicy          *SQSRedrivePolicy                   `json:"redrivePolicy,omitempty" yaml:"redrivePolicy,omitempty"`
	Policy                 *SQSPolicy                          `json:"policy,omitempty" yaml:"policy,omitempty"`
	Tags                   map[string]string                   `json:"tags,omitempty" yaml:"tags,omitempty"`
	Extensions             *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                    *lowasync.SQSQueue
}

// NewSQSQueue creates a new high-level SQSQueue.
func NewSQSQueue(q *lowasync.SQSQueue) *SQSQueue {
	s := new(SQSQueue)
	s.low = q
	s.Name = q.Name.Value
	s.ARN = q.ARN.Value
	s.FifoQueue = q.FifoQueue.Value
	s.DeduplicationScope = q.DeduplicationScope.Value
	s.FifoThroughputLimit = q.FifoThroughputLimit.Value
	s.DeliveryDelay = q.DeliveryDelay.Value
	s.VisibilityTimeout = q.VisibilityTimeout.Value
	s.ReceiveMessageWaitTime = q.ReceiveMessageWaitTime.Value
	s.MessageRetentionPeriod = q.MessageRetentionPeriod.Value
	if !q.RedrivePolicy.IsEmpty() {
		s.RedrivePolicy = NewSQSRedrivePolicy(q.RedrivePolicy.Value)
	}
	if !q.Policy.IsEmpty() {
		s.Policy = NewSQSPolicy(q.Policy.Value)
	}
	if !q.Tags.IsEmpty() {
		s.Tags = yamlStringMap(q.Tags.Value)
	}
	if orderedmap.Len(q.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(q.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSQueue.
func (s *SQSQueue) GoLow() *lowasync.SQSQueue { return s.low }

// GoLowUntyped returns the low-level SQSQueue with no type.
func (s *SQSQueue) GoLowUntyped() any { return s.low }

// SQSIdentifier represents a high-level SQS queue reference.
type SQSIdentifier struct {
	Name       string                              `json:"name,omitempty" yaml:"name,omitempty"`
	ARN        string                              `json:"arn,omitempty" yaml:"arn,omitempty"`
	FifoQueue  bool                                `json:"fifoQueue,omitempty" yaml:"fifoQueue,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.SQSIdentifier
}

// NewSQSIdentifier creates a new high-level SQSIdentifier.
func NewSQSIdentifier(q *lowasync.SQSIdentifier) *SQSIdentifier {
	s := new(SQSIdentifier)
	s.low = q
	s.Name = q.Name.Value
	s.ARN = q.ARN.Value
	s.FifoQueue = q.FifoQueue.Value
	if orderedmap.Len(q.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(q.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSIdentifier.
func (s *SQSIdentifier) GoLow() *lowasync.SQSIdentifier { return s.low }

// GoLowUntyped returns the low-level SQSIdentifier with no type.
func (s *SQSIdentifier) GoLowUntyped() any { return s.low }

// SQSRedrivePolicy represents high-level SQS redrive policy configuration.
type SQSRedrivePolicy struct {
	DeadLetterQueue *SQSIdentifier                      `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	MaxReceiveCount int                                 `json:"maxReceiveCount,omitempty" yaml:"maxReceiveCount,omitempty"`
	Extensions      *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low             *lowasync.SQSRedrivePolicy
}

// NewSQSRedrivePolicy creates a new high-level SQSRedrivePolicy.
func NewSQSRedrivePolicy(p *lowasync.SQSRedrivePolicy) *SQSRedrivePolicy {
	s := new(SQSRedrivePolicy)
	s.low = p
	s.MaxReceiveCount = p.MaxReceiveCount.Value
	if !p.DeadLetterQueue.IsEmpty() {
		s.DeadLetterQueue = NewSQSIdentifier(p.DeadLetterQueue.Value)
	}
	if orderedmap.Len(p.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(p.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSRedrivePolicy.
func (s *SQSRedrivePolicy) GoLow() *lowasync.SQSRedrivePolicy { return s.low }

// GoLowUntyped returns the low-level SQSRedrivePolicy with no type.
func (s *SQSRedrivePolicy) GoLowUntyped() any { return s.low }

// SQSPolicy represents a high-level SQS queue policy.
type SQSPolicy struct {
	Statements []*SQSPolicyStatement               `json:"statements,omitempty" yaml:"statements,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.SQSPolicy
}

// NewSQSPolicy creates a new high-level SQSPolicy.
func NewSQSPolicy(p *lowasync.SQSPolicy) *SQSPolicy {
	s := new(SQSPolicy)
	s.low = p
	if p.Statements.Value != nil {
		for _, statement := range p.Statements.Value {
			s.Statements = append(s.Statements, NewSQSPolicyStatement(statement.Value))
		}
	}
	if orderedmap.Len(p.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(p.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSPolicy.
func (s *SQSPolicy) GoLow() *lowasync.SQSPolicy { return s.low }

// GoLowUntyped returns the low-level SQSPolicy with no type.
func (s *SQSPolicy) GoLowUntyped() any { return s.low }

// SQSPolicyStatement represents a high-level SQS queue policy statement.
type SQSPolicyStatement struct {
	Effect     string                              `json:"effect,omitempty" yaml:"effect,omitempty"`
	Principal  *yaml.Node                          `json:"principal,omitempty" yaml:"principal,omitempty"`
	Action     *yaml.Node                          `json:"action,omitempty" yaml:"action,omitempty"`
	Resource   *yaml.Node                          `json:"resource,omitempty" yaml:"resource,omitempty"`
	Condition  *yaml.Node                          `json:"condition,omitempty" yaml:"condition,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.SQSPolicyStatement
}

// NewSQSPolicyStatement creates a new high-level SQSPolicyStatement.
func NewSQSPolicyStatement(p *lowasync.SQSPolicyStatement) *SQSPolicyStatement {
	s := new(SQSPolicyStatement)
	s.low = p
	s.Effect = p.Effect.Value
	s.Principal = p.Principal.Value
	s.Action = p.Action.Value
	s.Resource = p.Resource.Value
	s.Condition = p.Condition.Value
	if orderedmap.Len(p.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(p.Extensions)
	}
	return s
}

// GoLow returns the low-level SQSPolicyStatement.
func (s *SQSPolicyStatement) GoLow() *lowasync.SQSPolicyStatement { return s.low }

// GoLowUntyped returns the low-level SQSPolicyStatement with no type.
func (s *SQSPolicyStatement) GoLowUntyped() any { return s.low }

func yamlStringMap(node *yaml.Node) map[string]string {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}
	values := make(map[string]string)
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i]
		value := node.Content[i+1]
		if key.Kind == yaml.ScalarNode && value.Kind == yaml.ScalarNode {
			values[key.Value] = value.Value
		}
	}
	if len(values) == 0 {
		return nil
	}
	return values
}

// Kafka Server Binding

// KafkaServerBinding represents a high-level Kafka Server Binding.
type KafkaServerBinding struct {
	SchemaRegistryURL    string                              `json:"schemaRegistryUrl,omitempty" yaml:"schemaRegistryUrl,omitempty"`
	SchemaRegistryVendor string                              `json:"schemaRegistryVendor,omitempty" yaml:"schemaRegistryVendor,omitempty"`
	BindingVersion       string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions           *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                  *lowasync.KafkaServerBinding
}

// NewKafkaServerBinding creates a new high-level KafkaServerBinding.
func NewKafkaServerBinding(b *lowasync.KafkaServerBinding) *KafkaServerBinding {
	k := new(KafkaServerBinding)
	k.low = b
	k.SchemaRegistryURL = b.SchemaRegistryURL.Value
	k.SchemaRegistryVendor = b.SchemaRegistryVendor.Value
	k.BindingVersion = b.BindingVersion.Value
	if orderedmap.Len(b.Extensions) > 0 {
		k.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return k
}

// GoLow returns the low-level KafkaServerBinding.
func (k *KafkaServerBinding) GoLow() *lowasync.KafkaServerBinding { return k.low }

// GoLowUntyped returns the low-level KafkaServerBinding with no type.
func (k *KafkaServerBinding) GoLowUntyped() any { return k.low }

// Kafka Channel Binding

// KafkaChannelBinding represents a high-level Kafka Channel Binding.
type KafkaChannelBinding struct {
	Topic              string                              `json:"topic,omitempty" yaml:"topic,omitempty"`
	Partitions         int                                 `json:"partitions,omitempty" yaml:"partitions,omitempty"`
	Replicas           int                                 `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	TopicConfiguration *KafkaTopicConfiguration            `json:"topicConfiguration,omitempty" yaml:"topicConfiguration,omitempty"`
	BindingVersion     string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions         *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                *lowasync.KafkaChannelBinding
}

// NewKafkaChannelBinding creates a new high-level KafkaChannelBinding.
func NewKafkaChannelBinding(b *lowasync.KafkaChannelBinding) *KafkaChannelBinding {
	k := new(KafkaChannelBinding)
	k.low = b
	k.Topic = b.Topic.Value
	k.Partitions = b.Partitions.Value
	k.Replicas = b.Replicas.Value
	k.BindingVersion = b.BindingVersion.Value
	if !b.TopicConfiguration.IsEmpty() {
		k.TopicConfiguration = NewKafkaTopicConfiguration(b.TopicConfiguration.Value)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		k.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return k
}

// KafkaTopicConfiguration represents high-level Kafka topic configuration.
type KafkaTopicConfiguration struct {
	CleanupPolicy                     []string                            `json:"cleanup.policy,omitempty" yaml:"cleanup.policy,omitempty"`
	RetentionMs                       int64                               `json:"retention.ms,omitempty" yaml:"retention.ms,omitempty"`
	RetentionBytes                    int64                               `json:"retention.bytes,omitempty" yaml:"retention.bytes,omitempty"`
	DeleteRetentionMs                 int64                               `json:"delete.retention.ms,omitempty" yaml:"delete.retention.ms,omitempty"`
	MaxMessageBytes                   int                                 `json:"max.message.bytes,omitempty" yaml:"max.message.bytes,omitempty"`
	ConfluentKeySchemaValidation      bool                                `json:"confluent.key.schema.validation,omitempty" yaml:"confluent.key.schema.validation,omitempty"`
	ConfluentKeySubjectNameStrategy   string                              `json:"confluent.key.subject.name.strategy,omitempty" yaml:"confluent.key.subject.name.strategy,omitempty"`
	ConfluentValueSchemaValidation    bool                                `json:"confluent.value.schema.validation,omitempty" yaml:"confluent.value.schema.validation,omitempty"`
	ConfluentValueSubjectNameStrategy string                              `json:"confluent.value.subject.name.strategy,omitempty" yaml:"confluent.value.subject.name.strategy,omitempty"`
	Extensions                        *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                               *lowasync.KafkaTopicConfiguration
}

// NewKafkaTopicConfiguration creates a new high-level KafkaTopicConfiguration.
func NewKafkaTopicConfiguration(tc *lowasync.KafkaTopicConfiguration) *KafkaTopicConfiguration {
	k := new(KafkaTopicConfiguration)
	k.low = tc
	if tc.CleanupPolicy.Value != nil {
		for _, v := range tc.CleanupPolicy.Value {
			k.CleanupPolicy = append(k.CleanupPolicy, v.Value)
		}
	}
	k.RetentionMs = tc.RetentionMs.Value
	k.RetentionBytes = tc.RetentionBytes.Value
	k.DeleteRetentionMs = tc.DeleteRetentionMs.Value
	k.MaxMessageBytes = tc.MaxMessageBytes.Value
	k.ConfluentKeySchemaValidation = tc.ConfluentKeySchemaValidation.Value
	k.ConfluentKeySubjectNameStrategy = tc.ConfluentKeySubjectNameStrategy.Value
	k.ConfluentValueSchemaValidation = tc.ConfluentValueSchemaValidation.Value
	k.ConfluentValueSubjectNameStrategy = tc.ConfluentValueSubjectNameStrategy.Value
	if orderedmap.Len(tc.Extensions) > 0 {
		k.Extensions = high.ExtractExtensions(tc.Extensions)
	}
	return k
}

// GoLow returns the low-level KafkaTopicConfiguration.
func (k *KafkaTopicConfiguration) GoLow() *lowasync.KafkaTopicConfiguration { return k.low }

// GoLowUntyped returns the low-level KafkaTopicConfiguration with no type.
func (k *KafkaTopicConfiguration) GoLowUntyped() any { return k.low }

// GoLow returns the low-level KafkaChannelBinding.
func (k *KafkaChannelBinding) GoLow() *lowasync.KafkaChannelBinding { return k.low }

// GoLowUntyped returns the low-level KafkaChannelBinding with no type.
func (k *KafkaChannelBinding) GoLowUntyped() any { return k.low }

// Kafka Operation Binding

// KafkaOperationBinding represents a high-level Kafka Operation Binding.
type KafkaOperationBinding struct {
	GroupID        *highbase.SchemaProxy               `json:"groupId,omitempty" yaml:"groupId,omitempty"`
	ClientID       *highbase.SchemaProxy               `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.KafkaOperationBinding
}

// NewKafkaOperationBinding creates a new high-level KafkaOperationBinding.
func NewKafkaOperationBinding(b *lowasync.KafkaOperationBinding) *KafkaOperationBinding {
	k := new(KafkaOperationBinding)
	k.low = b
	k.BindingVersion = b.BindingVersion.Value
	if !b.GroupID.IsEmpty() {
		k.GroupID = highbase.NewSchemaProxy(&b.GroupID)
	}
	if !b.ClientID.IsEmpty() {
		k.ClientID = highbase.NewSchemaProxy(&b.ClientID)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		k.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return k
}

// GoLow returns the low-level KafkaOperationBinding.
func (k *KafkaOperationBinding) GoLow() *lowasync.KafkaOperationBinding { return k.low }

// GoLowUntyped returns the low-level KafkaOperationBinding with no type.
func (k *KafkaOperationBinding) GoLowUntyped() any { return k.low }

// Kafka Message Binding

// KafkaMessageBinding represents a high-level Kafka Message Binding.
type KafkaMessageBinding struct {
	Key                     *highbase.SchemaProxy               `json:"key,omitempty" yaml:"key,omitempty"`
	SchemaIDLocation        string                              `json:"schemaIdLocation,omitempty" yaml:"schemaIdLocation,omitempty"`
	SchemaIDPayloadEncoding string                              `json:"schemaIdPayloadEncoding,omitempty" yaml:"schemaIdPayloadEncoding,omitempty"`
	SchemaLookupStrategy    string                              `json:"schemaLookupStrategy,omitempty" yaml:"schemaLookupStrategy,omitempty"`
	BindingVersion          string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions              *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                     *lowasync.KafkaMessageBinding
}

// NewKafkaMessageBinding creates a new high-level KafkaMessageBinding.
func NewKafkaMessageBinding(b *lowasync.KafkaMessageBinding) *KafkaMessageBinding {
	k := new(KafkaMessageBinding)
	k.low = b
	k.BindingVersion = b.BindingVersion.Value
	k.SchemaIDLocation = b.SchemaIDLocation.Value
	k.SchemaIDPayloadEncoding = b.SchemaIDPayloadEncoding.Value
	k.SchemaLookupStrategy = b.SchemaLookupStrategy.Value
	if !b.Key.IsEmpty() {
		k.Key = highbase.NewSchemaProxy(&b.Key)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		k.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return k
}

// GoLow returns the low-level KafkaMessageBinding.
func (k *KafkaMessageBinding) GoLow() *lowasync.KafkaMessageBinding { return k.low }

// GoLowUntyped returns the low-level KafkaMessageBinding with no type.
func (k *KafkaMessageBinding) GoLowUntyped() any { return k.low }

// WebSocket Channel Binding

// WebSocketChannelBinding represents a high-level WebSocket Channel Binding.
type WebSocketChannelBinding struct {
	Method         string                              `json:"method,omitempty" yaml:"method,omitempty"`
	Query          *highbase.SchemaProxy               `json:"query,omitempty" yaml:"query,omitempty"`
	Headers        *highbase.SchemaProxy               `json:"headers,omitempty" yaml:"headers,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.WebSocketChannelBinding
}

// NewWebSocketChannelBinding creates a new high-level WebSocketChannelBinding.
func NewWebSocketChannelBinding(b *lowasync.WebSocketChannelBinding) *WebSocketChannelBinding {
	w := new(WebSocketChannelBinding)
	w.low = b
	w.Method = b.Method.Value
	w.BindingVersion = b.BindingVersion.Value
	if !b.Query.IsEmpty() {
		w.Query = highbase.NewSchemaProxy(&b.Query)
	}
	if !b.Headers.IsEmpty() {
		w.Headers = highbase.NewSchemaProxy(&b.Headers)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		w.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return w
}

// GoLow returns the low-level WebSocketChannelBinding.
func (w *WebSocketChannelBinding) GoLow() *lowasync.WebSocketChannelBinding { return w.low }

// GoLowUntyped returns the low-level WebSocketChannelBinding with no type.
func (w *WebSocketChannelBinding) GoLowUntyped() any { return w.low }

// AMQP Channel Binding

// AMQPChannelBinding represents a high-level AMQP Channel Binding.
type AMQPChannelBinding struct {
	Is             string                              `json:"is,omitempty" yaml:"is,omitempty"`
	Exchange       *AMQPExchange                       `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	Queue          *AMQPQueue                          `json:"queue,omitempty" yaml:"queue,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.AMQPChannelBinding
}

// NewAMQPChannelBinding creates a new high-level AMQPChannelBinding.
func NewAMQPChannelBinding(b *lowasync.AMQPChannelBinding) *AMQPChannelBinding {
	a := new(AMQPChannelBinding)
	a.low = b
	a.Is = b.Is.Value
	a.BindingVersion = b.BindingVersion.Value
	if !b.Exchange.IsEmpty() {
		a.Exchange = NewAMQPExchange(b.Exchange.Value)
	}
	if !b.Queue.IsEmpty() {
		a.Queue = NewAMQPQueue(b.Queue.Value)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return a
}

// AMQPExchange represents high-level AMQP exchange configuration.
type AMQPExchange struct {
	Name       string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Type       string                              `json:"type,omitempty" yaml:"type,omitempty"`
	Durable    bool                                `json:"durable,omitempty" yaml:"durable,omitempty"`
	AutoDelete bool                                `json:"autoDelete,omitempty" yaml:"autoDelete,omitempty"`
	VHost      string                              `json:"vhost,omitempty" yaml:"vhost,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.AMQPExchange
}

// NewAMQPExchange creates a new high-level AMQPExchange.
func NewAMQPExchange(ex *lowasync.AMQPExchange) *AMQPExchange {
	a := new(AMQPExchange)
	a.low = ex
	a.Name = ex.Name.Value
	a.Type = ex.Type.Value
	a.Durable = ex.Durable.Value
	a.AutoDelete = ex.AutoDelete.Value
	a.VHost = ex.VHost.Value
	if orderedmap.Len(ex.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(ex.Extensions)
	}
	return a
}

// GoLow returns the low-level AMQPExchange.
func (a *AMQPExchange) GoLow() *lowasync.AMQPExchange { return a.low }

// GoLowUntyped returns the low-level AMQPExchange with no type.
func (a *AMQPExchange) GoLowUntyped() any { return a.low }

// AMQPQueue represents high-level AMQP queue configuration.
type AMQPQueue struct {
	Name       string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Durable    bool                                `json:"durable,omitempty" yaml:"durable,omitempty"`
	Exclusive  bool                                `json:"exclusive,omitempty" yaml:"exclusive,omitempty"`
	AutoDelete bool                                `json:"autoDelete,omitempty" yaml:"autoDelete,omitempty"`
	VHost      string                              `json:"vhost,omitempty" yaml:"vhost,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.AMQPQueue
}

// NewAMQPQueue creates a new high-level AMQPQueue.
func NewAMQPQueue(q *lowasync.AMQPQueue) *AMQPQueue {
	a := new(AMQPQueue)
	a.low = q
	a.Name = q.Name.Value
	a.Durable = q.Durable.Value
	a.Exclusive = q.Exclusive.Value
	a.AutoDelete = q.AutoDelete.Value
	a.VHost = q.VHost.Value
	if orderedmap.Len(q.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(q.Extensions)
	}
	return a
}

// GoLow returns the low-level AMQPQueue.
func (a *AMQPQueue) GoLow() *lowasync.AMQPQueue { return a.low }

// GoLowUntyped returns the low-level AMQPQueue with no type.
func (a *AMQPQueue) GoLowUntyped() any { return a.low }

// GoLow returns the low-level AMQPChannelBinding.
func (a *AMQPChannelBinding) GoLow() *lowasync.AMQPChannelBinding { return a.low }

// GoLowUntyped returns the low-level AMQPChannelBinding with no type.
func (a *AMQPChannelBinding) GoLowUntyped() any { return a.low }

// AMQP Operation Binding

// AMQPOperationBinding represents a high-level AMQP Operation Binding.
type AMQPOperationBinding struct {
	Expiration     int                                 `json:"expiration,omitempty" yaml:"expiration,omitempty"`
	UserID         string                              `json:"userId,omitempty" yaml:"userId,omitempty"`
	CC             []string                            `json:"cc,omitempty" yaml:"cc,omitempty"`
	Priority       int                                 `json:"priority,omitempty" yaml:"priority,omitempty"`
	DeliveryMode   int                                 `json:"deliveryMode,omitempty" yaml:"deliveryMode,omitempty"`
	Mandatory      bool                                `json:"mandatory,omitempty" yaml:"mandatory,omitempty"`
	BCC            []string                            `json:"bcc,omitempty" yaml:"bcc,omitempty"`
	Timestamp      bool                                `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Ack            bool                                `json:"ack,omitempty" yaml:"ack,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.AMQPOperationBinding
}

// NewAMQPOperationBinding creates a new high-level AMQPOperationBinding.
func NewAMQPOperationBinding(b *lowasync.AMQPOperationBinding) *AMQPOperationBinding {
	a := new(AMQPOperationBinding)
	a.low = b
	a.Expiration = b.Expiration.Value
	a.UserID = b.UserID.Value
	a.Priority = b.Priority.Value
	a.DeliveryMode = b.DeliveryMode.Value
	a.Mandatory = b.Mandatory.Value
	a.Timestamp = b.Timestamp.Value
	a.Ack = b.Ack.Value
	a.BindingVersion = b.BindingVersion.Value
	if b.CC.Value != nil {
		for _, v := range b.CC.Value {
			a.CC = append(a.CC, v.Value)
		}
	}
	if b.BCC.Value != nil {
		for _, v := range b.BCC.Value {
			a.BCC = append(a.BCC, v.Value)
		}
	}
	if orderedmap.Len(b.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return a
}

// GoLow returns the low-level AMQPOperationBinding.
func (a *AMQPOperationBinding) GoLow() *lowasync.AMQPOperationBinding { return a.low }

// GoLowUntyped returns the low-level AMQPOperationBinding with no type.
func (a *AMQPOperationBinding) GoLowUntyped() any { return a.low }

// AMQP Message Binding

// AMQPMessageBinding represents a high-level AMQP Message Binding.
type AMQPMessageBinding struct {
	ContentEncoding string                              `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`
	MessageType     string                              `json:"messageType,omitempty" yaml:"messageType,omitempty"`
	BindingVersion  string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions      *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low             *lowasync.AMQPMessageBinding
}

// NewAMQPMessageBinding creates a new high-level AMQPMessageBinding.
func NewAMQPMessageBinding(b *lowasync.AMQPMessageBinding) *AMQPMessageBinding {
	a := new(AMQPMessageBinding)
	a.low = b
	a.ContentEncoding = b.ContentEncoding.Value
	a.MessageType = b.MessageType.Value
	a.BindingVersion = b.BindingVersion.Value
	if orderedmap.Len(b.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return a
}

// GoLow returns the low-level AMQPMessageBinding.
func (a *AMQPMessageBinding) GoLow() *lowasync.AMQPMessageBinding { return a.low }

// GoLowUntyped returns the low-level AMQPMessageBinding with no type.
func (a *AMQPMessageBinding) GoLowUntyped() any { return a.low }

// MQTT Server Binding

// MQTTServerBinding represents a high-level MQTT Server Binding.
type MQTTServerBinding struct {
	ClientID              string                              `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	CleanSession          bool                                `json:"cleanSession,omitempty" yaml:"cleanSession,omitempty"`
	LastWill              *MQTTLastWill                       `json:"lastWill,omitempty" yaml:"lastWill,omitempty"`
	KeepAlive             int                                 `json:"keepAlive,omitempty" yaml:"keepAlive,omitempty"`
	SessionExpiryInterval int                                 `json:"sessionExpiryInterval,omitempty" yaml:"sessionExpiryInterval,omitempty"`
	MaximumPacketSize     int                                 `json:"maximumPacketSize,omitempty" yaml:"maximumPacketSize,omitempty"`
	BindingVersion        string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions            *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                   *lowasync.MQTTServerBinding
}

// NewMQTTServerBinding creates a new high-level MQTTServerBinding.
func NewMQTTServerBinding(b *lowasync.MQTTServerBinding) *MQTTServerBinding {
	m := new(MQTTServerBinding)
	m.low = b
	m.ClientID = b.ClientID.Value
	m.CleanSession = b.CleanSession.Value
	m.KeepAlive = b.KeepAlive.Value
	m.SessionExpiryInterval = b.SessionExpiryInterval.Value
	m.MaximumPacketSize = b.MaximumPacketSize.Value
	m.BindingVersion = b.BindingVersion.Value
	if !b.LastWill.IsEmpty() {
		m.LastWill = NewMQTTLastWill(b.LastWill.Value)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		m.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return m
}

// MQTTLastWill represents high-level MQTT Last Will configuration.
type MQTTLastWill struct {
	Topic      string                              `json:"topic,omitempty" yaml:"topic,omitempty"`
	QoS        int                                 `json:"qos,omitempty" yaml:"qos,omitempty"`
	Message    string                              `json:"message,omitempty" yaml:"message,omitempty"`
	Retain     bool                                `json:"retain,omitempty" yaml:"retain,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.MQTTLastWill
}

// NewMQTTLastWill creates a new high-level MQTTLastWill.
func NewMQTTLastWill(lw *lowasync.MQTTLastWill) *MQTTLastWill {
	m := new(MQTTLastWill)
	m.low = lw
	m.Topic = lw.Topic.Value
	m.QoS = lw.QoS.Value
	m.Message = lw.Message.Value
	m.Retain = lw.Retain.Value
	if orderedmap.Len(lw.Extensions) > 0 {
		m.Extensions = high.ExtractExtensions(lw.Extensions)
	}
	return m
}

// GoLow returns the low-level MQTTLastWill.
func (m *MQTTLastWill) GoLow() *lowasync.MQTTLastWill { return m.low }

// GoLowUntyped returns the low-level MQTTLastWill with no type.
func (m *MQTTLastWill) GoLowUntyped() any { return m.low }

// GoLow returns the low-level MQTTServerBinding.
func (m *MQTTServerBinding) GoLow() *lowasync.MQTTServerBinding { return m.low }

// GoLowUntyped returns the low-level MQTTServerBinding with no type.
func (m *MQTTServerBinding) GoLowUntyped() any { return m.low }

// MQTT Operation Binding

// MQTTOperationBinding represents a high-level MQTT Operation Binding.
type MQTTOperationBinding struct {
	QoS            int                                 `json:"qos,omitempty" yaml:"qos,omitempty"`
	Retain         bool                                `json:"retain,omitempty" yaml:"retain,omitempty"`
	BindingVersion string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.MQTTOperationBinding
}

// NewMQTTOperationBinding creates a new high-level MQTTOperationBinding.
func NewMQTTOperationBinding(b *lowasync.MQTTOperationBinding) *MQTTOperationBinding {
	m := new(MQTTOperationBinding)
	m.low = b
	m.QoS = b.QoS.Value
	m.Retain = b.Retain.Value
	m.BindingVersion = b.BindingVersion.Value
	if orderedmap.Len(b.Extensions) > 0 {
		m.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return m
}

// GoLow returns the low-level MQTTOperationBinding.
func (m *MQTTOperationBinding) GoLow() *lowasync.MQTTOperationBinding { return m.low }

// GoLowUntyped returns the low-level MQTTOperationBinding with no type.
func (m *MQTTOperationBinding) GoLowUntyped() any { return m.low }

// MQTT Message Binding

// MQTTMessageBinding represents a high-level MQTT Message Binding.
type MQTTMessageBinding struct {
	PayloadFormatIndicator int                                 `json:"payloadFormatIndicator,omitempty" yaml:"payloadFormatIndicator,omitempty"`
	CorrelationData        *highbase.SchemaProxy               `json:"correlationData,omitempty" yaml:"correlationData,omitempty"`
	ContentType            string                              `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	ResponseTopic          string                              `json:"responseTopic,omitempty" yaml:"responseTopic,omitempty"`
	BindingVersion         string                              `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
	Extensions             *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                    *lowasync.MQTTMessageBinding
}

// NewMQTTMessageBinding creates a new high-level MQTTMessageBinding.
func NewMQTTMessageBinding(b *lowasync.MQTTMessageBinding) *MQTTMessageBinding {
	m := new(MQTTMessageBinding)
	m.low = b
	m.PayloadFormatIndicator = b.PayloadFormatIndicator.Value
	m.ContentType = b.ContentType.Value
	m.ResponseTopic = b.ResponseTopic.Value
	m.BindingVersion = b.BindingVersion.Value
	if !b.CorrelationData.IsEmpty() {
		m.CorrelationData = highbase.NewSchemaProxy(&b.CorrelationData)
	}
	if orderedmap.Len(b.Extensions) > 0 {
		m.Extensions = high.ExtractExtensions(b.Extensions)
	}
	return m
}

// GoLow returns the low-level MQTTMessageBinding.
func (m *MQTTMessageBinding) GoLow() *lowasync.MQTTMessageBinding { return m.low }

// GoLowUntyped returns the low-level MQTTMessageBinding with no type.
func (m *MQTTMessageBinding) GoLowUntyped() any { return m.low }
