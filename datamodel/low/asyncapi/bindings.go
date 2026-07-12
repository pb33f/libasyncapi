// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"hash/maphash"
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

// hashExtensions appends extension hashes to the hasher.
// This is the common ending for all Hash() methods.
func (b *BaseBinding) hashExtensions(h *maphash.Hash) {
	for _, ext := range low.HashExtensions(b.Extensions) {
		h.WriteString(ext)
		h.WriteByte(low.HASH_PIPE)
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

func (h *HTTPServerBinding) Hash() uint64 {
	return low.WithHasher(func(f *maphash.Hash) uint64 {
		h.hashExtensions(f)
		return f.Sum64()
	})
}

// HTTPChannelBinding represents a low-level AsyncAPI 3.0 HTTP Channel Binding object.
type HTTPChannelBinding struct {
	BaseBinding
}

func (h *HTTPChannelBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	h.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (h *HTTPChannelBinding) Hash() uint64 {
	return low.WithHasher(func(f *maphash.Hash) uint64 {
		h.hashExtensions(f)
		return f.Sum64()
	})
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

func (h *HTTPOperationBinding) Hash() uint64 {
	return low.WithHasher(func(f *maphash.Hash) uint64 {
		if !h.Method.IsEmpty() {
			f.WriteString(h.Method.Value)
			f.WriteByte(low.HASH_PIPE)
		}
		if !h.Query.IsEmpty() {
			f.WriteString(low.GenerateHashString(h.Query.Value))
			f.WriteByte(low.HASH_PIPE)
		}
		if !h.BindingVersion.IsEmpty() {
			f.WriteString(h.BindingVersion.Value)
			f.WriteByte(low.HASH_PIPE)
		}
		h.hashExtensions(f)
		return f.Sum64()
	})
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

func (h *HTTPMessageBinding) Hash() uint64 {
	return low.WithHasher(func(f *maphash.Hash) uint64 {
		if !h.Headers.IsEmpty() {
			f.WriteString(low.GenerateHashString(h.Headers.Value))
			f.WriteByte(low.HASH_PIPE)
		}
		if !h.StatusCode.IsEmpty() {
			low.HashInt64(f, int64(h.StatusCode.Value))
			f.WriteByte(low.HASH_PIPE)
		}
		if !h.BindingVersion.IsEmpty() {
			f.WriteString(h.BindingVersion.Value)
			f.WriteByte(low.HASH_PIPE)
		}
		h.hashExtensions(f)
		return f.Sum64()
	})
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

func (s *SQSServerBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSChannelBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !s.Queue.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.Queue.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.DeadLetterQueue.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.DeadLetterQueue.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.BindingVersion.IsEmpty() {
			h.WriteString(s.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSOperationBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if s.Queues.Value != nil {
			for _, queue := range s.Queues.Value {
				h.WriteString(low.GenerateHashString(queue.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !s.BindingVersion.IsEmpty() {
			h.WriteString(s.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		s.hashExtensions(h)
		return h.Sum64()
	})
}

// SQSMessageBinding represents a low-level AsyncAPI SQS Message Binding object.
type SQSMessageBinding struct {
	BaseBinding
}

func (s *SQSMessageBinding) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	s.initBuild(ctx, keyNode, root, idx)
	return nil
}

func (s *SQSMessageBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSQueue) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		hashSQSIdentifierFields(h, s.Name, s.ARN, s.FifoQueue)
		if !s.DeduplicationScope.IsEmpty() {
			h.WriteString(s.DeduplicationScope.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.FifoThroughputLimit.IsEmpty() {
			h.WriteString(s.FifoThroughputLimit.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.DeliveryDelay.IsEmpty() {
			low.HashInt64(h, int64(s.DeliveryDelay.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.VisibilityTimeout.IsEmpty() {
			low.HashInt64(h, int64(s.VisibilityTimeout.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.ReceiveMessageWaitTime.IsEmpty() {
			low.HashInt64(h, int64(s.ReceiveMessageWaitTime.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.MessageRetentionPeriod.IsEmpty() {
			low.HashInt64(h, int64(s.MessageRetentionPeriod.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.RedrivePolicy.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.RedrivePolicy.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.Policy.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.Policy.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.Tags.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.Tags.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSIdentifier) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		hashSQSIdentifierFields(h, s.Name, s.ARN, s.FifoQueue)
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSRedrivePolicy) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !s.DeadLetterQueue.IsEmpty() {
			h.WriteString(low.GenerateHashString(s.DeadLetterQueue.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !s.MaxReceiveCount.IsEmpty() {
			low.HashInt64(h, int64(s.MaxReceiveCount.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSPolicy) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if s.Statements.Value != nil {
			for _, statement := range s.Statements.Value {
				h.WriteString(low.GenerateHashString(statement.Value))
				h.WriteByte(low.HASH_PIPE)
			}
		}
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func (s *SQSPolicyStatement) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !s.Effect.IsEmpty() {
			h.WriteString(s.Effect.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		hashRawNodeReference(h, s.Principal)
		hashRawNodeReference(h, s.Action)
		hashRawNodeReference(h, s.Resource)
		hashRawNodeReference(h, s.Condition)
		s.hashExtensions(h)
		return h.Sum64()
	})
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

func hashRawNodeReference(h *maphash.Hash, ref low.NodeReference[*yaml.Node]) {
	if !ref.IsEmpty() {
		h.WriteString(low.GenerateHashString(ref.Value))
		h.WriteByte(low.HASH_PIPE)
	}
}

func hashSQSIdentifierFields(h *maphash.Hash, name low.NodeReference[string], arn low.NodeReference[string], fifoQueue low.NodeReference[bool]) {
	if !name.IsEmpty() {
		h.WriteString(name.Value)
		h.WriteByte(low.HASH_PIPE)
	}
	if !arn.IsEmpty() {
		h.WriteString(arn.Value)
		h.WriteByte(low.HASH_PIPE)
	}
	if !fifoQueue.IsEmpty() {
		low.HashBool(h, fifoQueue.Value)
		h.WriteByte(low.HASH_PIPE)
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

func (k *KafkaServerBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !k.SchemaRegistryURL.IsEmpty() {
			h.WriteString(k.SchemaRegistryURL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.SchemaRegistryVendor.IsEmpty() {
			h.WriteString(k.SchemaRegistryVendor.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.BindingVersion.IsEmpty() {
			h.WriteString(k.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		k.hashExtensions(h)
		return h.Sum64()
	})
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

func (k *KafkaChannelBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !k.Topic.IsEmpty() {
			h.WriteString(k.Topic.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.Partitions.IsEmpty() {
			low.HashInt64(h, int64(k.Partitions.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.Replicas.IsEmpty() {
			low.HashInt64(h, int64(k.Replicas.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.TopicConfiguration.IsEmpty() {
			h.WriteString(low.GenerateHashString(k.TopicConfiguration.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.BindingVersion.IsEmpty() {
			h.WriteString(k.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		k.hashExtensions(h)
		return h.Sum64()
	})
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

func (k *KafkaTopicConfiguration) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if k.CleanupPolicy.Value != nil {
			for _, policy := range k.CleanupPolicy.Value {
				h.WriteString(policy.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !k.RetentionMs.IsEmpty() {
			low.HashInt64(h, k.RetentionMs.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.RetentionBytes.IsEmpty() {
			low.HashInt64(h, k.RetentionBytes.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.DeleteRetentionMs.IsEmpty() {
			low.HashInt64(h, k.DeleteRetentionMs.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.MaxMessageBytes.IsEmpty() {
			low.HashInt64(h, int64(k.MaxMessageBytes.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.ConfluentKeySchemaValidation.IsEmpty() {
			low.HashBool(h, k.ConfluentKeySchemaValidation.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.ConfluentKeySubjectNameStrategy.IsEmpty() {
			h.WriteString(k.ConfluentKeySubjectNameStrategy.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.ConfluentValueSchemaValidation.IsEmpty() {
			low.HashBool(h, k.ConfluentValueSchemaValidation.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.ConfluentValueSubjectNameStrategy.IsEmpty() {
			h.WriteString(k.ConfluentValueSubjectNameStrategy.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		k.hashExtensions(h)
		return h.Sum64()
	})
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

func (k *KafkaOperationBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !k.GroupID.IsEmpty() {
			h.WriteString(low.GenerateHashString(k.GroupID.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.ClientID.IsEmpty() {
			h.WriteString(low.GenerateHashString(k.ClientID.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.BindingVersion.IsEmpty() {
			h.WriteString(k.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		k.hashExtensions(h)
		return h.Sum64()
	})
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

func (k *KafkaMessageBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !k.Key.IsEmpty() {
			h.WriteString(low.GenerateHashString(k.Key.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.SchemaIDLocation.IsEmpty() {
			h.WriteString(k.SchemaIDLocation.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.SchemaIDPayloadEncoding.IsEmpty() {
			h.WriteString(k.SchemaIDPayloadEncoding.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.SchemaLookupStrategy.IsEmpty() {
			h.WriteString(k.SchemaLookupStrategy.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !k.BindingVersion.IsEmpty() {
			h.WriteString(k.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		k.hashExtensions(h)
		return h.Sum64()
	})
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

func (w *WebSocketChannelBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !w.Method.IsEmpty() {
			h.WriteString(w.Method.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !w.Query.IsEmpty() {
			h.WriteString(low.GenerateHashString(w.Query.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !w.Headers.IsEmpty() {
			h.WriteString(low.GenerateHashString(w.Headers.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !w.BindingVersion.IsEmpty() {
			h.WriteString(w.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		w.hashExtensions(h)
		return h.Sum64()
	})
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

func (a *AMQPChannelBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !a.Is.IsEmpty() {
			h.WriteString(a.Is.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Exchange.IsEmpty() {
			h.WriteString(low.GenerateHashString(a.Exchange.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Queue.IsEmpty() {
			h.WriteString(low.GenerateHashString(a.Queue.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.BindingVersion.IsEmpty() {
			h.WriteString(a.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		a.hashExtensions(h)
		return h.Sum64()
	})
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

func (a *AMQPExchange) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !a.Name.IsEmpty() {
			h.WriteString(a.Name.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Type.IsEmpty() {
			h.WriteString(a.Type.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Durable.IsEmpty() {
			low.HashBool(h, a.Durable.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.AutoDelete.IsEmpty() {
			low.HashBool(h, a.AutoDelete.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.VHost.IsEmpty() {
			h.WriteString(a.VHost.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		a.hashExtensions(h)
		return h.Sum64()
	})
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

func (a *AMQPQueue) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !a.Name.IsEmpty() {
			h.WriteString(a.Name.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Durable.IsEmpty() {
			low.HashBool(h, a.Durable.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Exclusive.IsEmpty() {
			low.HashBool(h, a.Exclusive.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.AutoDelete.IsEmpty() {
			low.HashBool(h, a.AutoDelete.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.VHost.IsEmpty() {
			h.WriteString(a.VHost.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		a.hashExtensions(h)
		return h.Sum64()
	})
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

func (a *AMQPOperationBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !a.Expiration.IsEmpty() {
			low.HashInt64(h, int64(a.Expiration.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.UserID.IsEmpty() {
			h.WriteString(a.UserID.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if a.CC.Value != nil {
			for _, cc := range a.CC.Value {
				h.WriteString(cc.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !a.Priority.IsEmpty() {
			low.HashInt64(h, int64(a.Priority.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.DeliveryMode.IsEmpty() {
			low.HashInt64(h, int64(a.DeliveryMode.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Mandatory.IsEmpty() {
			low.HashBool(h, a.Mandatory.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if a.BCC.Value != nil {
			for _, bcc := range a.BCC.Value {
				h.WriteString(bcc.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		if !a.Timestamp.IsEmpty() {
			low.HashBool(h, a.Timestamp.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.Ack.IsEmpty() {
			low.HashBool(h, a.Ack.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.BindingVersion.IsEmpty() {
			h.WriteString(a.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		a.hashExtensions(h)
		return h.Sum64()
	})
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

func (a *AMQPMessageBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !a.ContentEncoding.IsEmpty() {
			h.WriteString(a.ContentEncoding.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.MessageType.IsEmpty() {
			h.WriteString(a.MessageType.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !a.BindingVersion.IsEmpty() {
			h.WriteString(a.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		a.hashExtensions(h)
		return h.Sum64()
	})
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

func (m *MQTTServerBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !m.ClientID.IsEmpty() {
			h.WriteString(m.ClientID.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.CleanSession.IsEmpty() {
			low.HashBool(h, m.CleanSession.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.LastWill.IsEmpty() {
			h.WriteString(low.GenerateHashString(m.LastWill.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.KeepAlive.IsEmpty() {
			low.HashInt64(h, int64(m.KeepAlive.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.SessionExpiryInterval.IsEmpty() {
			low.HashInt64(h, int64(m.SessionExpiryInterval.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.MaximumPacketSize.IsEmpty() {
			low.HashInt64(h, int64(m.MaximumPacketSize.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.BindingVersion.IsEmpty() {
			h.WriteString(m.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		m.hashExtensions(h)
		return h.Sum64()
	})
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

func (m *MQTTLastWill) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !m.Topic.IsEmpty() {
			h.WriteString(m.Topic.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.QoS.IsEmpty() {
			low.HashInt64(h, int64(m.QoS.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.Message.IsEmpty() {
			h.WriteString(m.Message.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.Retain.IsEmpty() {
			low.HashBool(h, m.Retain.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		m.hashExtensions(h)
		return h.Sum64()
	})
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

func (m *MQTTOperationBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !m.QoS.IsEmpty() {
			low.HashInt64(h, int64(m.QoS.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.Retain.IsEmpty() {
			low.HashBool(h, m.Retain.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.BindingVersion.IsEmpty() {
			h.WriteString(m.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		m.hashExtensions(h)
		return h.Sum64()
	})
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

func (m *MQTTMessageBinding) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !m.PayloadFormatIndicator.IsEmpty() {
			low.HashInt64(h, int64(m.PayloadFormatIndicator.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.CorrelationData.IsEmpty() {
			h.WriteString(low.GenerateHashString(m.CorrelationData.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.ContentType.IsEmpty() {
			h.WriteString(m.ContentType.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.ResponseTopic.IsEmpty() {
			h.WriteString(m.ResponseTopic.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !m.BindingVersion.IsEmpty() {
			h.WriteString(m.BindingVersion.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		m.hashExtensions(h)
		return h.Sum64()
	})
}
