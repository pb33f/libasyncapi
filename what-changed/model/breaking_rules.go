// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"sync"
)

var (
	defaultRulesOnce  sync.Once
	defaultRulesCache *BreakingRulesConfig

	activeConfigMu sync.RWMutex
	activeConfig   *BreakingRulesConfig
)

// ResetDefaultBreakingRules resets the cached default rules. This is primarily
// intended for testing scenarios where the cache needs to be cleared.
func ResetDefaultBreakingRules() {
	defaultRulesOnce = sync.Once{}
	defaultRulesCache = nil
}

// SetActiveBreakingRulesConfig sets the active breaking rules configuration used
// by comparison functions. Pass nil to reset to defaults.
//
// The active configuration is process-global and is captured when a comparison starts.
// Prefer passing a config directly to a comparator when callers need independent rules.
func SetActiveBreakingRulesConfig(config *BreakingRulesConfig) {
	activeConfigMu.Lock()
	defer activeConfigMu.Unlock()
	activeConfig = config
}

// GetActiveBreakingRulesConfig returns the currently active breaking rules config.
// If no custom config has been set, returns the default rules.
func GetActiveBreakingRulesConfig() *BreakingRulesConfig {
	return activeBreakingRulesConfig()
}

// ResetActiveBreakingRulesConfig clears any custom config and reverts to defaults.
func ResetActiveBreakingRulesConfig() {
	activeConfigMu.Lock()
	defer activeConfigMu.Unlock()
	activeConfig = nil
}

func activeBreakingRulesConfig() *BreakingRulesConfig {
	activeConfigMu.RLock()
	defer activeConfigMu.RUnlock()
	if activeConfig != nil {
		return activeConfig
	}
	return GenerateDefaultBreakingRules()
}

func comparisonConfig(configs []*BreakingRulesConfig) *BreakingRulesConfig {
	if len(configs) > 0 && configs[0] != nil {
		return configs[0]
	}
	return activeBreakingRulesConfig()
}

// GenerateDefaultBreakingRules returns the default breaking change rules for AsyncAPI 3.0.
// The returned config is cached and reused for performance — do not mutate it. To
// customize rules, use NewDefaultBreakingRulesConfig and Merge.
func GenerateDefaultBreakingRules() *BreakingRulesConfig {
	defaultRulesOnce.Do(func() {
		defaultRulesCache = buildDefaultRules()
	})
	return defaultRulesCache
}

// NewDefaultBreakingRulesConfig returns a fresh copy of the default breaking change rules
// that is safe to mutate or Merge custom overrides into. Pass it directly to a comparator,
// or activate it process-wide with SetActiveBreakingRulesConfig.
func NewDefaultBreakingRulesConfig() *BreakingRulesConfig {
	return buildDefaultRules()
}

// IsBreakingChange is a package-level helper that looks up whether a change is breaking
// using the currently active configuration.
func IsBreakingChange(component, property, changeType string, configs ...*BreakingRulesConfig) bool {
	return comparisonConfig(configs).IsBreaking(component, property, changeType)
}

// BreakingAdded returns whether adding the specified property is a breaking change.
func BreakingAdded(component, property string, configs ...*BreakingRulesConfig) bool {
	return IsBreakingChange(component, property, ChangeTypeAdded, configs...)
}

// BreakingModified returns whether modifying the specified property is a breaking change.
func BreakingModified(component, property string, configs ...*BreakingRulesConfig) bool {
	return IsBreakingChange(component, property, ChangeTypeModified, configs...)
}

// BreakingRemoved returns whether removing the specified property is a breaking change.
func BreakingRemoved(component, property string, configs ...*BreakingRulesConfig) bool {
	return IsBreakingChange(component, property, ChangeTypeRemoved, configs...)
}

func boolPtr(b bool) *bool {
	return &b
}

func rule(added, modified, removed bool) *BreakingChangeRule {
	return &BreakingChangeRule{
		Added:    boolPtr(added),
		Modified: boolPtr(modified),
		Removed:  boolPtr(removed),
	}
}

// buildDefaultRules creates the actual default rules configuration for AsyncAPI 3.0.
//
// The general principles are:
//   - identity and routing values (addresses, topics, queue names/ARNs, hosts,
//     protocols, channel refs) break when added, modified or removed.
//   - schema-bearing properties break when added or removed; modifications recurse
//     into libopenapi's schema comparison, which applies libopenapi's own schema rules
//     (configured via libopenapi's SetActiveBreakingRulesConfig, not this package's).
//   - delivery semantics (QoS, durability, ack, fifo, delivery mode) break on
//     modification and removal.
//   - informational values (titles, summaries, descriptions, docs, examples) never break.
//   - collection entries (channels, operations, messages, components maps) break on
//     removal, never on addition.
func buildDefaultRules() *BreakingRulesConfig {
	return &BreakingRulesConfig{
		AsyncAPI:           rule(true, true, true),
		ID:                 rule(false, true, true),
		DefaultContentType: rule(false, true, true),
		Servers:            rule(false, false, true),
		Channels:           rule(false, false, true),
		Operations:         rule(false, false, true),

		Info: &InfoRules{
			Title:          rule(false, false, false),
			Version:        rule(false, false, false),
			Description:    rule(false, false, false),
			TermsOfService: rule(false, false, false),
			Contact:        rule(false, false, false),
			License:        rule(false, false, false),
			Tags:           rule(false, false, false),
			ExternalDocs:   rule(false, false, false),
		},

		Contact: &ContactRules{
			Name:  rule(false, false, false),
			URL:   rule(false, false, false),
			Email: rule(false, false, false),
		},

		License: &LicenseRules{
			Name:       rule(false, false, false),
			Identifier: rule(false, false, false),
			URL:        rule(false, false, false),
		},

		Server: &ServerRules{
			Host:            rule(true, true, true),
			Protocol:        rule(true, true, true),
			ProtocolVersion: rule(false, true, false),
			Pathname:        rule(true, true, true),
			Title:           rule(false, false, false),
			Summary:         rule(false, false, false),
			Description:     rule(false, false, false),
			Variables:       rule(false, false, true),
			Security:        rule(false, false, true),
			Tags:            rule(false, false, false),
			ExternalDocs:    rule(false, false, false),
			Bindings:        rule(false, false, true),
		},

		ServerVariable: &ServerVariableRules{
			Enum:        rule(false, false, true),
			Default:     rule(true, true, true),
			Description: rule(false, false, false),
			Examples:    rule(false, false, false),
		},

		Channel: &ChannelRules{
			Address:      rule(true, true, true),
			Title:        rule(false, false, false),
			Summary:      rule(false, false, false),
			Description:  rule(false, false, false),
			Messages:     rule(false, false, true),
			Parameters:   rule(false, false, true),
			Servers:      rule(false, false, true),
			Tags:         rule(false, false, false),
			ExternalDocs: rule(false, false, false),
			Bindings:     rule(false, false, true),
		},

		Parameter: &ParameterRules{
			Enum:        rule(false, false, true),
			Default:     rule(false, true, false),
			Description: rule(false, false, false),
			Examples:    rule(false, false, false),
			Location:    rule(true, true, true),
		},

		Operation: &OperationRules{
			Action:       rule(true, true, true),
			Channel:      rule(true, true, true),
			Title:        rule(false, false, false),
			Summary:      rule(false, false, false),
			Description:  rule(false, false, false),
			Security:     rule(false, false, true),
			Tags:         rule(false, false, false),
			ExternalDocs: rule(false, false, false),
			Bindings:     rule(false, false, true),
			Traits:       rule(false, false, true),
			Messages:     rule(true, false, true),
			Reply:        rule(true, false, true),
		},

		OperationTrait: &OperationTraitRules{
			Title:        rule(false, false, false),
			Summary:      rule(false, false, false),
			Description:  rule(false, false, false),
			Security:     rule(false, false, true),
			Tags:         rule(false, false, false),
			ExternalDocs: rule(false, false, false),
			Bindings:     rule(false, false, true),
		},

		OperationReply: &OperationReplyRules{
			Address:  rule(false, true, true),
			Channel:  rule(true, true, true),
			Messages: rule(true, false, true),
		},

		OperationReplyAddr: &OperationReplyAddressRules{
			Description: rule(false, false, false),
			Location:    rule(true, true, true),
		},

		Message: &MessageRules{
			Headers:       rule(true, false, true),
			Payload:       rule(true, false, true),
			CorrelationID: rule(false, true, true),
			ContentType:   rule(false, true, true),
			Name:          rule(false, true, true),
			Title:         rule(false, false, false),
			Summary:       rule(false, false, false),
			Description:   rule(false, false, false),
			Tags:          rule(false, false, false),
			ExternalDocs:  rule(false, false, false),
			Bindings:      rule(false, false, true),
			Examples:      rule(false, false, false),
			Traits:        rule(false, false, true),
		},

		MessageTrait: &MessageTraitRules{
			Headers:       rule(true, false, true),
			CorrelationID: rule(false, true, true),
			ContentType:   rule(false, true, true),
			Name:          rule(false, true, true),
			Title:         rule(false, false, false),
			Summary:       rule(false, false, false),
			Description:   rule(false, false, false),
			Tags:          rule(false, false, false),
			ExternalDocs:  rule(false, false, false),
			Bindings:      rule(false, false, true),
			Examples:      rule(false, false, false),
		},

		MessageExample: &MessageExampleRules{
			Headers: rule(false, false, false),
			Payload: rule(false, false, false),
			Name:    rule(false, false, false),
			Summary: rule(false, false, false),
		},

		CorrelationID: &CorrelationIDRules{
			Description: rule(false, false, false),
			Location:    rule(false, true, true),
		},

		Components: &ComponentsRules{
			Schemas:           rule(false, false, true),
			Servers:           rule(false, false, true),
			Channels:          rule(false, false, true),
			Operations:        rule(false, false, true),
			Messages:          rule(false, false, true),
			SecuritySchemes:   rule(false, false, true),
			ServerVariables:   rule(false, false, true),
			Parameters:        rule(false, false, true),
			CorrelationIDs:    rule(false, false, true),
			Replies:           rule(false, false, true),
			ReplyAddresses:    rule(false, false, true),
			ExternalDocs:      rule(false, false, false),
			Tags:              rule(false, false, false),
			OperationTraits:   rule(false, false, true),
			MessageTraits:     rule(false, false, true),
			ServerBindings:    rule(false, false, true),
			ChannelBindings:   rule(false, false, true),
			OperationBindings: rule(false, false, true),
			MessageBindings:   rule(false, false, true),
		},

		SecurityScheme: &SecuritySchemeRules{
			Type:             rule(true, true, true),
			Description:      rule(false, false, false),
			Name:             rule(true, true, true),
			In:               rule(true, true, true),
			Scheme:           rule(true, true, true),
			BearerFormat:     rule(false, false, false),
			Flows:            rule(false, false, true),
			OpenIDConnectURL: rule(false, false, false),
			Scopes:           rule(false, false, true),
		},

		OAuthFlows: &OAuthFlowsRules{
			Implicit:          rule(false, false, true),
			Password:          rule(false, false, true),
			ClientCredentials: rule(false, false, true),
			AuthorizationCode: rule(false, false, true),
		},

		OAuthFlow: &OAuthFlowRules{
			AuthorizationURL: rule(true, true, true),
			TokenURL:         rule(true, true, true),
			RefreshURL:       rule(false, true, true),
			AvailableScopes:  rule(false, true, true),
		},

		Tag: &TagRules{
			Name:         rule(false, true, true),
			Description:  rule(false, false, false),
			ExternalDocs: rule(false, false, false),
		},

		ExternalDocs: &ExternalDocsRules{
			URL:         rule(false, false, false),
			Description: rule(false, false, false),
		},

		ServerBindings: &ServerBindingsRules{
			HTTP:  rule(false, false, true),
			Kafka: rule(false, false, true),
			MQTT:  rule(false, false, true),
			SQS:   rule(false, false, true),
		},

		ChannelBindings: &ChannelBindingsRules{
			HTTP:      rule(false, false, true),
			WebSocket: rule(false, false, true),
			Kafka:     rule(false, false, true),
			AMQP:      rule(false, false, true),
			SQS:       rule(false, false, true),
		},

		OperationBindings: &OperationBindingsRules{
			HTTP:  rule(false, false, true),
			Kafka: rule(false, false, true),
			AMQP:  rule(false, false, true),
			MQTT:  rule(false, false, true),
			SQS:   rule(false, false, true),
		},

		MessageBindings: &MessageBindingsRules{
			HTTP:  rule(false, false, true),
			Kafka: rule(false, false, true),
			AMQP:  rule(false, false, true),
			MQTT:  rule(false, false, true),
			SQS:   rule(false, false, true),
		},

		HTTPOperation: &HTTPOperationBindingRules{
			Method:         rule(true, true, true),
			Query:          rule(true, false, true),
			BindingVersion: rule(false, false, false),
		},

		HTTPMessage: &HTTPMessageBindingRules{
			Headers:        rule(true, false, true),
			StatusCode:     rule(false, true, true),
			BindingVersion: rule(false, false, false),
		},

		KafkaServer: &KafkaServerBindingRules{
			SchemaRegistryURL:    rule(false, true, true),
			SchemaRegistryVendor: rule(false, true, false),
			BindingVersion:       rule(false, false, false),
		},

		KafkaChannel: &KafkaChannelBindingRules{
			Topic:              rule(true, true, true),
			Partitions:         rule(false, true, false),
			Replicas:           rule(false, false, false),
			TopicConfiguration: rule(false, false, true),
			BindingVersion:     rule(false, false, false),
		},

		KafkaTopicConfig: &KafkaTopicConfigurationRules{
			CleanupPolicy:                     rule(false, true, false),
			RetentionMs:                       rule(false, true, false),
			RetentionBytes:                    rule(false, true, false),
			DeleteRetentionMs:                 rule(false, true, false),
			MaxMessageBytes:                   rule(false, true, true),
			ConfluentKeySchemaValidation:      rule(false, true, false),
			ConfluentKeySubjectNameStrategy:   rule(false, true, false),
			ConfluentValueSchemaValidation:    rule(false, true, false),
			ConfluentValueSubjectNameStrategy: rule(false, true, false),
		},

		KafkaOperation: &KafkaOperationBindingRules{
			GroupID:        rule(true, false, true),
			ClientID:       rule(true, false, true),
			BindingVersion: rule(false, false, false),
		},

		KafkaMessage: &KafkaMessageBindingRules{
			Key:                     rule(true, false, true),
			SchemaIDLocation:        rule(false, true, true),
			SchemaIDPayloadEncoding: rule(false, true, true),
			SchemaLookupStrategy:    rule(false, true, true),
			BindingVersion:          rule(false, false, false),
		},

		WebSocketChannel: &WebSocketChannelBindingRules{
			Method:         rule(true, true, true),
			Query:          rule(true, false, true),
			Headers:        rule(true, false, true),
			BindingVersion: rule(false, false, false),
		},

		AMQPChannel: &AMQPChannelBindingRules{
			Is:             rule(true, true, true),
			Exchange:       rule(true, false, true),
			Queue:          rule(true, false, true),
			BindingVersion: rule(false, false, false),
		},

		AMQPExchange: &AMQPExchangeRules{
			Name:       rule(true, true, true),
			Type:       rule(true, true, true),
			Durable:    rule(false, true, true),
			AutoDelete: rule(false, true, true),
			VHost:      rule(false, true, true),
		},

		AMQPQueue: &AMQPQueueRules{
			Name:       rule(true, true, true),
			Durable:    rule(false, true, true),
			Exclusive:  rule(false, true, true),
			AutoDelete: rule(false, true, true),
			VHost:      rule(false, true, true),
		},

		AMQPOperation: &AMQPOperationBindingRules{
			Expiration:     rule(false, true, false),
			UserID:         rule(false, true, true),
			CC:             rule(false, true, true),
			Priority:       rule(false, true, false),
			DeliveryMode:   rule(false, true, true),
			Mandatory:      rule(false, true, true),
			BCC:            rule(false, true, true),
			Timestamp:      rule(false, true, false),
			Ack:            rule(false, true, true),
			BindingVersion: rule(false, false, false),
		},

		AMQPMessage: &AMQPMessageBindingRules{
			ContentEncoding: rule(false, true, true),
			MessageType:     rule(false, true, false),
			BindingVersion:  rule(false, false, false),
		},

		MQTTServer: &MQTTServerBindingRules{
			ClientID:              rule(false, true, true),
			CleanSession:          rule(false, true, true),
			LastWill:              rule(false, true, true),
			KeepAlive:             rule(false, true, false),
			SessionExpiryInterval: rule(false, true, false),
			MaximumPacketSize:     rule(false, true, true),
			BindingVersion:        rule(false, false, false),
		},

		MQTTLastWill: &MQTTLastWillRules{
			Topic:   rule(true, true, true),
			QoS:     rule(false, true, false),
			Message: rule(false, false, false),
			Retain:  rule(false, true, false),
		},

		MQTTOperation: &MQTTOperationBindingRules{
			QoS:            rule(false, true, false),
			Retain:         rule(false, true, false),
			BindingVersion: rule(false, false, false),
		},

		MQTTMessage: &MQTTMessageBindingRules{
			PayloadFormatIndicator: rule(false, true, true),
			CorrelationData:        rule(true, false, true),
			ContentType:            rule(false, true, true),
			ResponseTopic:          rule(true, true, true),
			BindingVersion:         rule(false, false, false),
		},

		SQSChannel: &SQSChannelBindingRules{
			Queue:           rule(true, false, true),
			DeadLetterQueue: rule(false, false, true),
			BindingVersion:  rule(false, false, false),
		},

		SQSOperation: &SQSOperationBindingRules{
			Queues:         rule(true, false, true),
			BindingVersion: rule(false, false, false),
		},

		SQSQueue: &SQSQueueRules{
			Name:                   rule(true, true, true),
			ARN:                    rule(false, true, true),
			FifoQueue:              rule(false, true, true),
			DeduplicationScope:     rule(false, true, false),
			FifoThroughputLimit:    rule(false, true, false),
			DeliveryDelay:          rule(false, true, false),
			VisibilityTimeout:      rule(false, true, false),
			ReceiveMessageWaitTime: rule(false, true, false),
			MessageRetentionPeriod: rule(false, true, false),
			RedrivePolicy:          rule(false, false, true),
			Policy:                 rule(false, true, true),
			Tags:                   rule(false, false, false),
		},

		SQSIdentifier: &SQSIdentifierRules{
			Name:      rule(true, true, true),
			ARN:       rule(false, true, true),
			FifoQueue: rule(false, true, true),
		},

		SQSRedrivePolicy: &SQSRedrivePolicyRules{
			DeadLetterQueue: rule(false, true, true),
			MaxReceiveCount: rule(false, true, false),
		},

		SQSPolicy: &SQSPolicyRules{
			Statements: rule(false, true, true),
		},

		SQSPolicyStatement: &SQSPolicyStatementRules{
			Effect:    rule(false, true, true),
			Principal: rule(false, true, true),
			Action:    rule(false, true, true),
			Resource:  rule(false, true, true),
			Condition: rule(false, true, true),
		},
	}
}
