// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"crypto/sha256"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
	"go.yaml.in/yaml/v4"
)

// BaseBinding provides common fields and methods for all binding types.
// This reduces code duplication across the 15+ binding structs.
type BaseBinding struct {
	Extensions *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode    *yaml.Node
	RootNode   *yaml.Node
	idx        *index.SpecIndex
	ctx        context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the binding.
func (b *BaseBinding) GetRootNode() *yaml.Node { return b.RootNode }

// GetKeyNode returns the key yaml node of the binding.
func (b *BaseBinding) GetKeyNode() *yaml.Node { return b.KeyNode }

// GetIndex returns the spec index.
func (b *BaseBinding) GetIndex() *index.SpecIndex { return b.idx }

// GetContext returns the context.
func (b *BaseBinding) GetContext() context.Context { return b.ctx }

// GetExtensions returns all extensions for the binding.
func (b *BaseBinding) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return b.Extensions
}

// initBuild performs common initialization for all binding Build() methods.
// Returns the processed root node for further type-specific extraction.
func (b *BaseBinding) initBuild(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) *yaml.Node {
	b.KeyNode = keyNode
	root = utils.NodeAlias(root)
	b.RootNode = root
	utils.CheckForMergeNodes(root)
	b.Reference = new(low.Reference)
	b.Nodes = low.ExtractNodes(ctx, root)
	b.Extensions = low.ExtractExtensions(root)
	b.idx = idx
	b.ctx = ctx
	return root
}

// hashExtensions appends extension hashes to the string builder.
// This is the common ending for all Hash() methods.
func (b *BaseBinding) hashExtensions(sb *strings.Builder) {
	for _, ext := range low.HashExtensions(b.Extensions) {
		sb.WriteString(ext)
		sb.WriteByte('|')
	}
}

// HTTP Bindings

// HTTPServerBinding represents a low-level AsyncAPI 3.0 HTTP Server Binding object.
type HTTPServerBinding struct {
	BaseBinding
}

func (h *HTTPServerBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	h.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (h *HTTPServerBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	h.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// HTTPChannelBinding represents a low-level AsyncAPI 3.0 HTTP Channel Binding object.
type HTTPChannelBinding struct {
	BaseBinding
}

func (h *HTTPChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	h.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (h *HTTPChannelBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	h.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// HTTPOperationBinding represents a low-level AsyncAPI 3.0 HTTP Operation Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/http/README.md#operation
type HTTPOperationBinding struct {
	BaseBinding
	Method         low.NodeReference[string]
	Query          low.NodeReference[*base.SchemaProxy]
	BindingVersion low.NodeReference[string]
}

func (h *HTTPOperationBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = h.initBuild(ctx, keyNode, root, idx)

	query, _ := low.ExtractObject[*base.SchemaProxy](ctx, QueryLabel, root, idx)
	h.Query = query

	return nil
}

func (h *HTTPOperationBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !h.Method.IsEmpty() {
		sb.WriteString(h.Method.Value)
		sb.WriteByte('|')
	}
	if !h.Query.IsEmpty() {
		sb.WriteString(low.GenerateHashString(h.Query.Value))
		sb.WriteByte('|')
	}
	if !h.BindingVersion.IsEmpty() {
		sb.WriteString(h.BindingVersion.Value)
		sb.WriteByte('|')
	}
	h.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// HTTPMessageBinding represents a low-level AsyncAPI 3.0 HTTP Message Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/http/README.md#message
type HTTPMessageBinding struct {
	BaseBinding
	Headers        low.NodeReference[*base.SchemaProxy]
	StatusCode     low.NodeReference[int]
	BindingVersion low.NodeReference[string]
}

func (h *HTTPMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = h.initBuild(ctx, keyNode, root, idx)

	headers, _ := low.ExtractObject[*base.SchemaProxy](ctx, HeadersLabel, root, idx)
	h.Headers = headers

	return nil
}

func (h *HTTPMessageBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !h.Headers.IsEmpty() {
		sb.WriteString(low.GenerateHashString(h.Headers.Value))
		sb.WriteByte('|')
	}
	if !h.BindingVersion.IsEmpty() {
		sb.WriteString(h.BindingVersion.Value)
		sb.WriteByte('|')
	}
	h.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQS Bindings

// SQSServerBinding represents a low-level AsyncAPI SQS Server Binding object.
type SQSServerBinding struct {
	BaseBinding
}

func (s *SQSServerBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	s.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (s *SQSServerBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSChannelBinding represents a low-level AsyncAPI SQS Channel Binding object.
type SQSChannelBinding struct {
	BaseBinding
	Queue           low.NodeReference[*SQSQueue]
	DeadLetterQueue low.NodeReference[*SQSQueue]
	BindingVersion  low.NodeReference[string]
}

func (s *SQSChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)

	queue, _ := low.ExtractObject[*SQSQueue](ctx, QueueLabel, root, idx)
	s.Queue = queue

	deadLetterQueue, _ := low.ExtractObject[*SQSQueue](ctx, DeadLetterQueueLabel, root, idx)
	s.DeadLetterQueue = deadLetterQueue

	return nil
}

func (s *SQSChannelBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !s.Queue.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.Queue.Value))
		sb.WriteByte('|')
	}
	if !s.DeadLetterQueue.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.DeadLetterQueue.Value))
		sb.WriteByte('|')
	}
	if !s.BindingVersion.IsEmpty() {
		sb.WriteString(s.BindingVersion.Value)
		sb.WriteByte('|')
	}
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSOperationBinding represents a low-level AsyncAPI SQS Operation Binding object.
type SQSOperationBinding struct {
	BaseBinding
	Queues         low.NodeReference[[]low.ValueReference[*SQSQueue]]
	BindingVersion low.NodeReference[string]
}

func (s *SQSOperationBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)

	queues, qLabel, qValue, err := low.ExtractArray[*SQSQueue](ctx, QueuesLabel, root, idx)
	if err != nil {
		return err
	}
	if queues != nil {
		s.Queues = low.NodeReference[[]low.ValueReference[*SQSQueue]]{
			Value:     queues,
			KeyNode:   qLabel,
			ValueNode: qValue,
		}
	}

	return nil
}

func (s *SQSOperationBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if s.Queues.Value != nil {
		for _, queue := range s.Queues.Value {
			sb.WriteString(low.GenerateHashString(queue.Value))
			sb.WriteByte('|')
		}
	}
	if !s.BindingVersion.IsEmpty() {
		sb.WriteString(s.BindingVersion.Value)
		sb.WriteByte('|')
	}
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSMessageBinding represents a low-level AsyncAPI SQS Message Binding object.
type SQSMessageBinding struct {
	BaseBinding
}

func (s *SQSMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	s.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (s *SQSMessageBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSQueue represents SQS queue configuration.
type SQSQueue struct {
	BaseBinding
	Name                   low.NodeReference[string]
	ARN                    low.NodeReference[string]
	FifoQueue              low.NodeReference[bool]
	DeduplicationScope     low.NodeReference[string]
	FifoThroughputLimit    low.NodeReference[string]
	DeliveryDelay          low.NodeReference[int]
	VisibilityTimeout      low.NodeReference[int]
	ReceiveMessageWaitTime low.NodeReference[int]
	MessageRetentionPeriod low.NodeReference[int]
	RedrivePolicy          low.NodeReference[*SQSRedrivePolicy]
	Policy                 low.NodeReference[*SQSPolicy]
	Tags                   low.NodeReference[*yaml.Node]
}

func (s *SQSQueue) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)

	redrivePolicy, _ := low.ExtractObject[*SQSRedrivePolicy](ctx, RedrivePolicyLabel, root, idx)
	s.RedrivePolicy = redrivePolicy

	policy, _ := low.ExtractObject[*SQSPolicy](ctx, PolicyLabel, root, idx)
	s.Policy = policy

	_, tagsLabel, tagsValue := utils.FindKeyNodeFullTop(TagsLabel, root.Content)
	if tagsValue != nil {
		s.Tags = low.NodeReference[*yaml.Node]{
			Value:     tagsValue,
			KeyNode:   tagsLabel,
			ValueNode: tagsValue,
		}
	}

	return nil
}

func (s *SQSQueue) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	hashSQSIdentifierFields(sb, s.Name, s.ARN, s.FifoQueue)
	if !s.DeduplicationScope.IsEmpty() {
		sb.WriteString(s.DeduplicationScope.Value)
		sb.WriteByte('|')
	}
	if !s.FifoThroughputLimit.IsEmpty() {
		sb.WriteString(s.FifoThroughputLimit.Value)
		sb.WriteByte('|')
	}
	if !s.DeliveryDelay.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.DeliveryDelay.Value))
		sb.WriteByte('|')
	}
	if !s.VisibilityTimeout.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.VisibilityTimeout.Value))
		sb.WriteByte('|')
	}
	if !s.ReceiveMessageWaitTime.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.ReceiveMessageWaitTime.Value))
		sb.WriteByte('|')
	}
	if !s.MessageRetentionPeriod.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.MessageRetentionPeriod.Value))
		sb.WriteByte('|')
	}
	if !s.RedrivePolicy.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.RedrivePolicy.Value))
		sb.WriteByte('|')
	}
	if !s.Policy.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.Policy.Value))
		sb.WriteByte('|')
	}
	if !s.Tags.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.Tags.Value))
		sb.WriteByte('|')
	}
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSIdentifier represents a named SQS queue reference.
type SQSIdentifier struct {
	BaseBinding
	Name      low.NodeReference[string]
	ARN       low.NodeReference[string]
	FifoQueue low.NodeReference[bool]
}

func (s *SQSIdentifier) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	s.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (s *SQSIdentifier) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	hashSQSIdentifierFields(sb, s.Name, s.ARN, s.FifoQueue)
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSRedrivePolicy represents SQS redrive policy configuration.
type SQSRedrivePolicy struct {
	BaseBinding
	DeadLetterQueue low.NodeReference[*SQSIdentifier]
	MaxReceiveCount low.NodeReference[int]
}

func (s *SQSRedrivePolicy) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)

	deadLetterQueue, _ := low.ExtractObject[*SQSIdentifier](ctx, DeadLetterQueueLabel, root, idx)
	s.DeadLetterQueue = deadLetterQueue

	return nil
}

func (s *SQSRedrivePolicy) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !s.DeadLetterQueue.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.DeadLetterQueue.Value))
		sb.WriteByte('|')
	}
	if !s.MaxReceiveCount.IsEmpty() {
		sb.WriteString(low.GenerateHashString(s.MaxReceiveCount.Value))
		sb.WriteByte('|')
	}
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSPolicy represents an SQS queue policy.
type SQSPolicy struct {
	BaseBinding
	Statements low.NodeReference[[]low.ValueReference[*SQSPolicyStatement]]
}

func (s *SQSPolicy) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)

	statements, stLabel, stValue, err := low.ExtractArray[*SQSPolicyStatement](ctx, StatementsLabel, root, idx)
	if err != nil {
		return err
	}
	if statements != nil {
		s.Statements = low.NodeReference[[]low.ValueReference[*SQSPolicyStatement]]{
			Value:     statements,
			KeyNode:   stLabel,
			ValueNode: stValue,
		}
	}

	return nil
}

func (s *SQSPolicy) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if s.Statements.Value != nil {
		for _, statement := range s.Statements.Value {
			sb.WriteString(low.GenerateHashString(statement.Value))
			sb.WriteByte('|')
		}
	}
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// SQSPolicyStatement represents an SQS queue policy statement.
type SQSPolicyStatement struct {
	BaseBinding
	Effect    low.NodeReference[string]
	Principal low.NodeReference[*yaml.Node]
	Action    low.NodeReference[*yaml.Node]
	Resource  low.NodeReference[*yaml.Node]
	Condition low.NodeReference[*yaml.Node]
}

func (s *SQSPolicyStatement) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = s.initBuild(ctx, keyNode, root, idx)
	s.Principal = extractRawNodeReference(PrincipalLabel, root)
	s.Action = extractRawNodeReference(ActionLabel, root)
	s.Resource = extractRawNodeReference(ResourceLabel, root)
	s.Condition = extractRawNodeReference(ConditionLabel, root)
	return nil
}

func (s *SQSPolicyStatement) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !s.Effect.IsEmpty() {
		sb.WriteString(s.Effect.Value)
		sb.WriteByte('|')
	}
	hashRawNodeReference(sb, s.Principal)
	hashRawNodeReference(sb, s.Action)
	hashRawNodeReference(sb, s.Resource)
	hashRawNodeReference(sb, s.Condition)
	s.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

func extractRawNodeReference(label string, root *yaml.Node) low.NodeReference[*yaml.Node] {
	var keyNode, valueNode *yaml.Node
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i]
		value := root.Content[i+1]
		if strings.EqualFold(key.Value, label) {
			keyNode = key
			valueNode = value
			break
		}
	}
	if valueNode == nil {
		return low.NodeReference[*yaml.Node]{}
	}
	return low.NodeReference[*yaml.Node]{
		Value:     valueNode,
		KeyNode:   keyNode,
		ValueNode: valueNode,
	}
}

func hashRawNodeReference(sb *strings.Builder, ref low.NodeReference[*yaml.Node]) {
	if !ref.IsEmpty() {
		sb.WriteString(low.GenerateHashString(ref.Value))
		sb.WriteByte('|')
	}
}

func hashSQSIdentifierFields(sb *strings.Builder, name low.NodeReference[string], arn low.NodeReference[string], fifoQueue low.NodeReference[bool]) {
	if !name.IsEmpty() {
		sb.WriteString(name.Value)
		sb.WriteByte('|')
	}
	if !arn.IsEmpty() {
		sb.WriteString(arn.Value)
		sb.WriteByte('|')
	}
	if !fifoQueue.IsEmpty() {
		sb.WriteString(low.GenerateHashString(fifoQueue.Value))
		sb.WriteByte('|')
	}
}

// Kafka Bindings

// KafkaServerBinding represents a low-level AsyncAPI 3.0 Kafka Server Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/kafka/README.md#server
type KafkaServerBinding struct {
	BaseBinding
	SchemaRegistryURL    low.NodeReference[string]
	SchemaRegistryVendor low.NodeReference[string]
	BindingVersion       low.NodeReference[string]
}

func (k *KafkaServerBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	k.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (k *KafkaServerBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !k.SchemaRegistryURL.IsEmpty() {
		sb.WriteString(k.SchemaRegistryURL.Value)
		sb.WriteByte('|')
	}
	if !k.SchemaRegistryVendor.IsEmpty() {
		sb.WriteString(k.SchemaRegistryVendor.Value)
		sb.WriteByte('|')
	}
	if !k.BindingVersion.IsEmpty() {
		sb.WriteString(k.BindingVersion.Value)
		sb.WriteByte('|')
	}
	k.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// KafkaChannelBinding represents a low-level AsyncAPI 3.0 Kafka Channel Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/kafka/README.md#channel
type KafkaChannelBinding struct {
	BaseBinding
	Topic              low.NodeReference[string]
	Partitions         low.NodeReference[int]
	Replicas           low.NodeReference[int]
	TopicConfiguration low.NodeReference[*KafkaTopicConfiguration]
	BindingVersion     low.NodeReference[string]
}

func (k *KafkaChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = k.initBuild(ctx, keyNode, root, idx)

	topicConfig, _ := low.ExtractObject[*KafkaTopicConfiguration](ctx, TopicConfigurationLabel, root, idx)
	k.TopicConfiguration = topicConfig

	return nil
}

func (k *KafkaChannelBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !k.Topic.IsEmpty() {
		sb.WriteString(k.Topic.Value)
		sb.WriteByte('|')
	}
	if !k.TopicConfiguration.IsEmpty() {
		sb.WriteString(low.GenerateHashString(k.TopicConfiguration.Value))
		sb.WriteByte('|')
	}
	if !k.BindingVersion.IsEmpty() {
		sb.WriteString(k.BindingVersion.Value)
		sb.WriteByte('|')
	}
	k.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// KafkaTopicConfiguration represents Kafka topic configuration.
type KafkaTopicConfiguration struct {
	BaseBinding
	CleanupPolicy                     low.NodeReference[[]low.ValueReference[string]]
	RetentionMs                       low.NodeReference[int64]
	RetentionBytes                    low.NodeReference[int64]
	DeleteRetentionMs                 low.NodeReference[int64]
	MaxMessageBytes                   low.NodeReference[int]
	ConfluentKeySchemaValidation      low.NodeReference[bool]
	ConfluentKeySubjectNameStrategy   low.NodeReference[string]
	ConfluentValueSchemaValidation    low.NodeReference[bool]
	ConfluentValueSubjectNameStrategy low.NodeReference[string]
}

func (k *KafkaTopicConfiguration) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	k.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (k *KafkaTopicConfiguration) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	k.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// KafkaOperationBinding represents a low-level AsyncAPI 3.0 Kafka Operation Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/kafka/README.md#operation
type KafkaOperationBinding struct {
	BaseBinding
	GroupID        low.NodeReference[*base.SchemaProxy]
	ClientID       low.NodeReference[*base.SchemaProxy]
	BindingVersion low.NodeReference[string]
}

func (k *KafkaOperationBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = k.initBuild(ctx, keyNode, root, idx)

	groupID, _ := low.ExtractObject[*base.SchemaProxy](ctx, GroupIDLabel, root, idx)
	k.GroupID = groupID

	clientID, _ := low.ExtractObject[*base.SchemaProxy](ctx, ClientIDLabel, root, idx)
	k.ClientID = clientID

	return nil
}

func (k *KafkaOperationBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !k.GroupID.IsEmpty() {
		sb.WriteString(low.GenerateHashString(k.GroupID.Value))
		sb.WriteByte('|')
	}
	if !k.ClientID.IsEmpty() {
		sb.WriteString(low.GenerateHashString(k.ClientID.Value))
		sb.WriteByte('|')
	}
	if !k.BindingVersion.IsEmpty() {
		sb.WriteString(k.BindingVersion.Value)
		sb.WriteByte('|')
	}
	k.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// KafkaMessageBinding represents a low-level AsyncAPI 3.0 Kafka Message Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/kafka/README.md#message
type KafkaMessageBinding struct {
	BaseBinding
	Key                     low.NodeReference[*base.SchemaProxy]
	SchemaIDLocation        low.NodeReference[string]
	SchemaIDPayloadEncoding low.NodeReference[string]
	SchemaLookupStrategy    low.NodeReference[string]
	BindingVersion          low.NodeReference[string]
}

func (k *KafkaMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = k.initBuild(ctx, keyNode, root, idx)

	key, _ := low.ExtractObject[*base.SchemaProxy](ctx, KeyLabel, root, idx)
	k.Key = key

	return nil
}

func (k *KafkaMessageBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !k.Key.IsEmpty() {
		sb.WriteString(low.GenerateHashString(k.Key.Value))
		sb.WriteByte('|')
	}
	if !k.SchemaIDLocation.IsEmpty() {
		sb.WriteString(k.SchemaIDLocation.Value)
		sb.WriteByte('|')
	}
	if !k.SchemaIDPayloadEncoding.IsEmpty() {
		sb.WriteString(k.SchemaIDPayloadEncoding.Value)
		sb.WriteByte('|')
	}
	if !k.SchemaLookupStrategy.IsEmpty() {
		sb.WriteString(k.SchemaLookupStrategy.Value)
		sb.WriteByte('|')
	}
	if !k.BindingVersion.IsEmpty() {
		sb.WriteString(k.BindingVersion.Value)
		sb.WriteByte('|')
	}
	k.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// WebSocket Bindings

// WebSocketChannelBinding represents a low-level AsyncAPI 3.0 WebSocket Channel Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/websockets/README.md#channel
type WebSocketChannelBinding struct {
	BaseBinding
	Method         low.NodeReference[string]
	Query          low.NodeReference[*base.SchemaProxy]
	Headers        low.NodeReference[*base.SchemaProxy]
	BindingVersion low.NodeReference[string]
}

func (w *WebSocketChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = w.initBuild(ctx, keyNode, root, idx)

	query, _ := low.ExtractObject[*base.SchemaProxy](ctx, QueryLabel, root, idx)
	w.Query = query

	headers, _ := low.ExtractObject[*base.SchemaProxy](ctx, HeadersLabel, root, idx)
	w.Headers = headers

	return nil
}

func (w *WebSocketChannelBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !w.Method.IsEmpty() {
		sb.WriteString(w.Method.Value)
		sb.WriteByte('|')
	}
	if !w.Query.IsEmpty() {
		sb.WriteString(low.GenerateHashString(w.Query.Value))
		sb.WriteByte('|')
	}
	if !w.Headers.IsEmpty() {
		sb.WriteString(low.GenerateHashString(w.Headers.Value))
		sb.WriteByte('|')
	}
	if !w.BindingVersion.IsEmpty() {
		sb.WriteString(w.BindingVersion.Value)
		sb.WriteByte('|')
	}
	w.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// AMQP Bindings

// AMQPChannelBinding represents a low-level AsyncAPI 3.0 AMQP Channel Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/amqp/README.md#channel
type AMQPChannelBinding struct {
	BaseBinding
	Is             low.NodeReference[string]
	Exchange       low.NodeReference[*AMQPExchange]
	Queue          low.NodeReference[*AMQPQueue]
	BindingVersion low.NodeReference[string]
}

func (a *AMQPChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = a.initBuild(ctx, keyNode, root, idx)

	exchange, _ := low.ExtractObject[*AMQPExchange](ctx, ExchangeLabel, root, idx)
	a.Exchange = exchange

	queue, _ := low.ExtractObject[*AMQPQueue](ctx, QueueLabel, root, idx)
	a.Queue = queue

	return nil
}

func (a *AMQPChannelBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !a.Is.IsEmpty() {
		sb.WriteString(a.Is.Value)
		sb.WriteByte('|')
	}
	if !a.Exchange.IsEmpty() {
		sb.WriteString(low.GenerateHashString(a.Exchange.Value))
		sb.WriteByte('|')
	}
	if !a.Queue.IsEmpty() {
		sb.WriteString(low.GenerateHashString(a.Queue.Value))
		sb.WriteByte('|')
	}
	if !a.BindingVersion.IsEmpty() {
		sb.WriteString(a.BindingVersion.Value)
		sb.WriteByte('|')
	}
	a.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// AMQPExchange represents AMQP exchange configuration.
type AMQPExchange struct {
	BaseBinding
	Name       low.NodeReference[string]
	Type       low.NodeReference[string]
	Durable    low.NodeReference[bool]
	AutoDelete low.NodeReference[bool]
	VHost      low.NodeReference[string]
}

func (a *AMQPExchange) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	a.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (a *AMQPExchange) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !a.Name.IsEmpty() {
		sb.WriteString(a.Name.Value)
		sb.WriteByte('|')
	}
	if !a.Type.IsEmpty() {
		sb.WriteString(a.Type.Value)
		sb.WriteByte('|')
	}
	if !a.VHost.IsEmpty() {
		sb.WriteString(a.VHost.Value)
		sb.WriteByte('|')
	}
	a.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// AMQPQueue represents AMQP queue configuration.
type AMQPQueue struct {
	BaseBinding
	Name       low.NodeReference[string]
	Durable    low.NodeReference[bool]
	Exclusive  low.NodeReference[bool]
	AutoDelete low.NodeReference[bool]
	VHost      low.NodeReference[string]
}

func (a *AMQPQueue) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	a.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (a *AMQPQueue) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !a.Name.IsEmpty() {
		sb.WriteString(a.Name.Value)
		sb.WriteByte('|')
	}
	if !a.VHost.IsEmpty() {
		sb.WriteString(a.VHost.Value)
		sb.WriteByte('|')
	}
	a.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// AMQPOperationBinding represents a low-level AsyncAPI 3.0 AMQP Operation Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/amqp/README.md#operation
type AMQPOperationBinding struct {
	BaseBinding
	Expiration     low.NodeReference[int]
	UserID         low.NodeReference[string]
	CC             low.NodeReference[[]low.ValueReference[string]]
	Priority       low.NodeReference[int]
	DeliveryMode   low.NodeReference[int]
	Mandatory      low.NodeReference[bool]
	BCC            low.NodeReference[[]low.ValueReference[string]]
	Timestamp      low.NodeReference[bool]
	Ack            low.NodeReference[bool]
	BindingVersion low.NodeReference[string]
}

func (a *AMQPOperationBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	a.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (a *AMQPOperationBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !a.UserID.IsEmpty() {
		sb.WriteString(a.UserID.Value)
		sb.WriteByte('|')
	}
	if !a.BindingVersion.IsEmpty() {
		sb.WriteString(a.BindingVersion.Value)
		sb.WriteByte('|')
	}
	a.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// AMQPMessageBinding represents a low-level AsyncAPI 3.0 AMQP Message Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/amqp/README.md#message
type AMQPMessageBinding struct {
	BaseBinding
	ContentEncoding low.NodeReference[string]
	MessageType     low.NodeReference[string]
	BindingVersion  low.NodeReference[string]
}

func (a *AMQPMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	a.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (a *AMQPMessageBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !a.ContentEncoding.IsEmpty() {
		sb.WriteString(a.ContentEncoding.Value)
		sb.WriteByte('|')
	}
	if !a.MessageType.IsEmpty() {
		sb.WriteString(a.MessageType.Value)
		sb.WriteByte('|')
	}
	if !a.BindingVersion.IsEmpty() {
		sb.WriteString(a.BindingVersion.Value)
		sb.WriteByte('|')
	}
	a.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// MQTT Bindings

// MQTTServerBinding represents a low-level AsyncAPI 3.0 MQTT Server Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/mqtt/README.md#server
type MQTTServerBinding struct {
	BaseBinding
	ClientID              low.NodeReference[string]
	CleanSession          low.NodeReference[bool]
	LastWill              low.NodeReference[*MQTTLastWill]
	KeepAlive             low.NodeReference[int]
	SessionExpiryInterval low.NodeReference[int]
	MaximumPacketSize     low.NodeReference[int]
	BindingVersion        low.NodeReference[string]
}

func (m *MQTTServerBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = m.initBuild(ctx, keyNode, root, idx)

	lastWill, _ := low.ExtractObject[*MQTTLastWill](ctx, LastWillLabel, root, idx)
	m.LastWill = lastWill

	return nil
}

func (m *MQTTServerBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !m.ClientID.IsEmpty() {
		sb.WriteString(m.ClientID.Value)
		sb.WriteByte('|')
	}
	if !m.LastWill.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.LastWill.Value))
		sb.WriteByte('|')
	}
	if !m.BindingVersion.IsEmpty() {
		sb.WriteString(m.BindingVersion.Value)
		sb.WriteByte('|')
	}
	m.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// MQTTLastWill represents MQTT Last Will configuration.
type MQTTLastWill struct {
	BaseBinding
	Topic   low.NodeReference[string]
	QoS     low.NodeReference[int]
	Message low.NodeReference[string]
	Retain  low.NodeReference[bool]
}

func (m *MQTTLastWill) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	m.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (m *MQTTLastWill) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !m.Topic.IsEmpty() {
		sb.WriteString(m.Topic.Value)
		sb.WriteByte('|')
	}
	if !m.Message.IsEmpty() {
		sb.WriteString(m.Message.Value)
		sb.WriteByte('|')
	}
	m.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// MQTTOperationBinding represents a low-level AsyncAPI 3.0 MQTT Operation Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/mqtt/README.md#operation
type MQTTOperationBinding struct {
	BaseBinding
	QoS            low.NodeReference[int]
	Retain         low.NodeReference[bool]
	BindingVersion low.NodeReference[string]
}

func (m *MQTTOperationBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	m.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (m *MQTTOperationBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !m.BindingVersion.IsEmpty() {
		sb.WriteString(m.BindingVersion.Value)
		sb.WriteByte('|')
	}
	m.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}

// MQTTMessageBinding represents a low-level AsyncAPI 3.0 MQTT Message Binding object.
//
//	https://github.com/asyncapi/bindings/blob/master/mqtt/README.md#message
type MQTTMessageBinding struct {
	BaseBinding
	PayloadFormatIndicator low.NodeReference[int]
	CorrelationData        low.NodeReference[*base.SchemaProxy]
	ContentType            low.NodeReference[string]
	ResponseTopic          low.NodeReference[string]
	BindingVersion         low.NodeReference[string]
}

func (m *MQTTMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	root = m.initBuild(ctx, keyNode, root, idx)

	correlationData, _ := low.ExtractObject[*base.SchemaProxy](ctx, CorrelationDataLabel, root, idx)
	m.CorrelationData = correlationData

	return nil
}

func (m *MQTTMessageBinding) Hash() [32]byte {
	sb := low.GetStringBuilder()
	defer low.PutStringBuilder(sb)
	if !m.CorrelationData.IsEmpty() {
		sb.WriteString(low.GenerateHashString(m.CorrelationData.Value))
		sb.WriteByte('|')
	}
	if !m.ContentType.IsEmpty() {
		sb.WriteString(m.ContentType.Value)
		sb.WriteByte('|')
	}
	if !m.ResponseTopic.IsEmpty() {
		sb.WriteString(m.ResponseTopic.Value)
		sb.WriteByte('|')
	}
	if !m.BindingVersion.IsEmpty() {
		sb.WriteString(m.BindingVersion.Value)
		sb.WriteByte('|')
	}
	m.hashExtensions(sb)
	return sha256.Sum256([]byte(sb.String()))
}
