// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

// InfoRules defines breaking rules for the Info object properties.
type InfoRules struct {
	Title          *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Version        *BreakingChangeRule `json:"version,omitempty" yaml:"version,omitempty"`
	Description    *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService *BreakingChangeRule `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *BreakingChangeRule `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *BreakingChangeRule `json:"license,omitempty" yaml:"license,omitempty"`
	Tags           *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs   *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ContactRules defines breaking rules for the Contact object properties.
type ContactRules struct {
	Name  *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	URL   *BreakingChangeRule `json:"url,omitempty" yaml:"url,omitempty"`
	Email *BreakingChangeRule `json:"email,omitempty" yaml:"email,omitempty"`
}

// LicenseRules defines breaking rules for the License object properties.
type LicenseRules struct {
	Name       *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Identifier *BreakingChangeRule `json:"identifier,omitempty" yaml:"identifier,omitempty"`
	URL        *BreakingChangeRule `json:"url,omitempty" yaml:"url,omitempty"`
}

// ServerRules defines breaking rules for the Server object properties.
type ServerRules struct {
	Host            *BreakingChangeRule `json:"host,omitempty" yaml:"host,omitempty"`
	Protocol        *BreakingChangeRule `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	ProtocolVersion *BreakingChangeRule `json:"protocolVersion,omitempty" yaml:"protocolVersion,omitempty"`
	Pathname        *BreakingChangeRule `json:"pathname,omitempty" yaml:"pathname,omitempty"`
	Title           *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary         *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description     *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Variables       *BreakingChangeRule `json:"variables,omitempty" yaml:"variables,omitempty"`
	Security        *BreakingChangeRule `json:"security,omitempty" yaml:"security,omitempty"`
	Tags            *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs    *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings        *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
}

// ServerVariableRules defines breaking rules for the Server Variable object properties.
type ServerVariableRules struct {
	Enum        *BreakingChangeRule `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     *BreakingChangeRule `json:"default,omitempty" yaml:"default,omitempty"`
	Description *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Examples    *BreakingChangeRule `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// ChannelRules defines breaking rules for the Channel object properties.
type ChannelRules struct {
	Address      *BreakingChangeRule `json:"address,omitempty" yaml:"address,omitempty"`
	Title        *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Messages     *BreakingChangeRule `json:"messages,omitempty" yaml:"messages,omitempty"`
	Parameters   *BreakingChangeRule `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Servers      *BreakingChangeRule `json:"servers,omitempty" yaml:"servers,omitempty"`
	Tags         *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
}

// ParameterRules defines breaking rules for the channel Parameter object properties.
type ParameterRules struct {
	Enum        *BreakingChangeRule `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     *BreakingChangeRule `json:"default,omitempty" yaml:"default,omitempty"`
	Description *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Examples    *BreakingChangeRule `json:"examples,omitempty" yaml:"examples,omitempty"`
	Location    *BreakingChangeRule `json:"location,omitempty" yaml:"location,omitempty"`
}

// OperationRules defines breaking rules for the Operation object properties.
type OperationRules struct {
	Action       *BreakingChangeRule `json:"action,omitempty" yaml:"action,omitempty"`
	Channel      *BreakingChangeRule `json:"channel,omitempty" yaml:"channel,omitempty"`
	Title        *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Security     *BreakingChangeRule `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Traits       *BreakingChangeRule `json:"traits,omitempty" yaml:"traits,omitempty"`
	Messages     *BreakingChangeRule `json:"messages,omitempty" yaml:"messages,omitempty"`
	Reply        *BreakingChangeRule `json:"reply,omitempty" yaml:"reply,omitempty"`
}

// OperationTraitRules defines breaking rules for the Operation Trait object properties.
type OperationTraitRules struct {
	Title        *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary      *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description  *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Security     *BreakingChangeRule `json:"security,omitempty" yaml:"security,omitempty"`
	Tags         *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings     *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
}

// OperationReplyRules defines breaking rules for the Operation Reply object properties.
type OperationReplyRules struct {
	Address  *BreakingChangeRule `json:"address,omitempty" yaml:"address,omitempty"`
	Channel  *BreakingChangeRule `json:"channel,omitempty" yaml:"channel,omitempty"`
	Messages *BreakingChangeRule `json:"messages,omitempty" yaml:"messages,omitempty"`
}

// OperationReplyAddressRules defines breaking rules for the Operation Reply Address object properties.
type OperationReplyAddressRules struct {
	Description *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Location    *BreakingChangeRule `json:"location,omitempty" yaml:"location,omitempty"`
}

// MessageRules defines breaking rules for the Message object properties.
type MessageRules struct {
	Headers       *BreakingChangeRule `json:"headers,omitempty" yaml:"headers,omitempty"`
	Payload       *BreakingChangeRule `json:"payload,omitempty" yaml:"payload,omitempty"`
	CorrelationID *BreakingChangeRule `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	ContentType   *BreakingChangeRule `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Name          *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Title         *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary       *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Tags          *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs  *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings      *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Examples      *BreakingChangeRule `json:"examples,omitempty" yaml:"examples,omitempty"`
	Traits        *BreakingChangeRule `json:"traits,omitempty" yaml:"traits,omitempty"`
}

// MessageTraitRules defines breaking rules for the Message Trait object properties.
type MessageTraitRules struct {
	Headers       *BreakingChangeRule `json:"headers,omitempty" yaml:"headers,omitempty"`
	CorrelationID *BreakingChangeRule `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	ContentType   *BreakingChangeRule `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Name          *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Title         *BreakingChangeRule `json:"title,omitempty" yaml:"title,omitempty"`
	Summary       *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Tags          *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs  *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings      *BreakingChangeRule `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Examples      *BreakingChangeRule `json:"examples,omitempty" yaml:"examples,omitempty"`
}

// MessageExampleRules defines breaking rules for the Message Example object properties.
type MessageExampleRules struct {
	Headers *BreakingChangeRule `json:"headers,omitempty" yaml:"headers,omitempty"`
	Payload *BreakingChangeRule `json:"payload,omitempty" yaml:"payload,omitempty"`
	Name    *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Summary *BreakingChangeRule `json:"summary,omitempty" yaml:"summary,omitempty"`
}

// CorrelationIDRules defines breaking rules for the Correlation ID object properties.
type CorrelationIDRules struct {
	Description *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Location    *BreakingChangeRule `json:"location,omitempty" yaml:"location,omitempty"`
}

// ComponentsRules defines breaking rules for the Components object maps.
type ComponentsRules struct {
	Schemas           *BreakingChangeRule `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Servers           *BreakingChangeRule `json:"servers,omitempty" yaml:"servers,omitempty"`
	Channels          *BreakingChangeRule `json:"channels,omitempty" yaml:"channels,omitempty"`
	Operations        *BreakingChangeRule `json:"operations,omitempty" yaml:"operations,omitempty"`
	Messages          *BreakingChangeRule `json:"messages,omitempty" yaml:"messages,omitempty"`
	SecuritySchemes   *BreakingChangeRule `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	ServerVariables   *BreakingChangeRule `json:"serverVariables,omitempty" yaml:"serverVariables,omitempty"`
	Parameters        *BreakingChangeRule `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	CorrelationIDs    *BreakingChangeRule `json:"correlationIds,omitempty" yaml:"correlationIds,omitempty"`
	Replies           *BreakingChangeRule `json:"replies,omitempty" yaml:"replies,omitempty"`
	ReplyAddresses    *BreakingChangeRule `json:"replyAddresses,omitempty" yaml:"replyAddresses,omitempty"`
	ExternalDocs      *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Tags              *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
	OperationTraits   *BreakingChangeRule `json:"operationTraits,omitempty" yaml:"operationTraits,omitempty"`
	MessageTraits     *BreakingChangeRule `json:"messageTraits,omitempty" yaml:"messageTraits,omitempty"`
	ServerBindings    *BreakingChangeRule `json:"serverBindings,omitempty" yaml:"serverBindings,omitempty"`
	ChannelBindings   *BreakingChangeRule `json:"channelBindings,omitempty" yaml:"channelBindings,omitempty"`
	OperationBindings *BreakingChangeRule `json:"operationBindings,omitempty" yaml:"operationBindings,omitempty"`
	MessageBindings   *BreakingChangeRule `json:"messageBindings,omitempty" yaml:"messageBindings,omitempty"`
}

// SecuritySchemeRules defines breaking rules for the Security Scheme object properties.
type SecuritySchemeRules struct {
	Type             *BreakingChangeRule `json:"type,omitempty" yaml:"type,omitempty"`
	Description      *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	Name             *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	In               *BreakingChangeRule `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           *BreakingChangeRule `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     *BreakingChangeRule `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *BreakingChangeRule `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL *BreakingChangeRule `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
	Scopes           *BreakingChangeRule `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

// OAuthFlowsRules defines breaking rules for the OAuth Flows object properties.
type OAuthFlowsRules struct {
	Implicit          *BreakingChangeRule `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *BreakingChangeRule `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *BreakingChangeRule `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *BreakingChangeRule `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlowRules defines breaking rules for the OAuth Flow object properties.
type OAuthFlowRules struct {
	AuthorizationURL *BreakingChangeRule `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         *BreakingChangeRule `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       *BreakingChangeRule `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	AvailableScopes  *BreakingChangeRule `json:"availableScopes,omitempty" yaml:"availableScopes,omitempty"`
}

// TagRules defines breaking rules for the Tag object properties.
type TagRules struct {
	Name         *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Description  *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *BreakingChangeRule `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
}

// ExternalDocsRules defines breaking rules for the External Documentation object properties.
type ExternalDocsRules struct {
	URL         *BreakingChangeRule `json:"url,omitempty" yaml:"url,omitempty"`
	Description *BreakingChangeRule `json:"description,omitempty" yaml:"description,omitempty"`
}

// ServerBindingsRules defines breaking rules for the Server Bindings container properties.
type ServerBindingsRules struct {
	HTTP  *BreakingChangeRule `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka *BreakingChangeRule `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	MQTT  *BreakingChangeRule `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS   *BreakingChangeRule `json:"sqs,omitempty" yaml:"sqs,omitempty"`
}

// ChannelBindingsRules defines breaking rules for the Channel Bindings container properties.
type ChannelBindingsRules struct {
	HTTP      *BreakingChangeRule `json:"http,omitempty" yaml:"http,omitempty"`
	WebSocket *BreakingChangeRule `json:"ws,omitempty" yaml:"ws,omitempty"`
	Kafka     *BreakingChangeRule `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP      *BreakingChangeRule `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	SQS       *BreakingChangeRule `json:"sqs,omitempty" yaml:"sqs,omitempty"`
}

// OperationBindingsRules defines breaking rules for the Operation Bindings container properties.
type OperationBindingsRules struct {
	HTTP  *BreakingChangeRule `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka *BreakingChangeRule `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP  *BreakingChangeRule `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTT  *BreakingChangeRule `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS   *BreakingChangeRule `json:"sqs,omitempty" yaml:"sqs,omitempty"`
}

// MessageBindingsRules defines breaking rules for the Message Bindings container properties.
type MessageBindingsRules struct {
	HTTP  *BreakingChangeRule `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka *BreakingChangeRule `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	AMQP  *BreakingChangeRule `json:"amqp,omitempty" yaml:"amqp,omitempty"`
	MQTT  *BreakingChangeRule `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS   *BreakingChangeRule `json:"sqs,omitempty" yaml:"sqs,omitempty"`
}

// HTTPOperationBindingRules defines breaking rules for the HTTP Operation Binding properties.
type HTTPOperationBindingRules struct {
	Method         *BreakingChangeRule `json:"method,omitempty" yaml:"method,omitempty"`
	Query          *BreakingChangeRule `json:"query,omitempty" yaml:"query,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// HTTPMessageBindingRules defines breaking rules for the HTTP Message Binding properties.
type HTTPMessageBindingRules struct {
	Headers        *BreakingChangeRule `json:"headers,omitempty" yaml:"headers,omitempty"`
	StatusCode     *BreakingChangeRule `json:"statusCode,omitempty" yaml:"statusCode,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// KafkaServerBindingRules defines breaking rules for the Kafka Server Binding properties.
type KafkaServerBindingRules struct {
	SchemaRegistryURL    *BreakingChangeRule `json:"schemaRegistryUrl,omitempty" yaml:"schemaRegistryUrl,omitempty"`
	SchemaRegistryVendor *BreakingChangeRule `json:"schemaRegistryVendor,omitempty" yaml:"schemaRegistryVendor,omitempty"`
	BindingVersion       *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// KafkaChannelBindingRules defines breaking rules for the Kafka Channel Binding properties.
type KafkaChannelBindingRules struct {
	Topic              *BreakingChangeRule `json:"topic,omitempty" yaml:"topic,omitempty"`
	Partitions         *BreakingChangeRule `json:"partitions,omitempty" yaml:"partitions,omitempty"`
	Replicas           *BreakingChangeRule `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	TopicConfiguration *BreakingChangeRule `json:"topicConfiguration,omitempty" yaml:"topicConfiguration,omitempty"`
	BindingVersion     *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// KafkaTopicConfigurationRules defines breaking rules for the Kafka Topic Configuration properties.
type KafkaTopicConfigurationRules struct {
	CleanupPolicy                     *BreakingChangeRule `json:"cleanup.policy,omitempty" yaml:"cleanup.policy,omitempty"`
	RetentionMs                       *BreakingChangeRule `json:"retention.ms,omitempty" yaml:"retention.ms,omitempty"`
	RetentionBytes                    *BreakingChangeRule `json:"retention.bytes,omitempty" yaml:"retention.bytes,omitempty"`
	DeleteRetentionMs                 *BreakingChangeRule `json:"delete.retention.ms,omitempty" yaml:"delete.retention.ms,omitempty"`
	MaxMessageBytes                   *BreakingChangeRule `json:"max.message.bytes,omitempty" yaml:"max.message.bytes,omitempty"`
	ConfluentKeySchemaValidation      *BreakingChangeRule `json:"confluent.key.schema.validation,omitempty" yaml:"confluent.key.schema.validation,omitempty"`
	ConfluentKeySubjectNameStrategy   *BreakingChangeRule `json:"confluent.key.subject.name.strategy,omitempty" yaml:"confluent.key.subject.name.strategy,omitempty"`
	ConfluentValueSchemaValidation    *BreakingChangeRule `json:"confluent.value.schema.validation,omitempty" yaml:"confluent.value.schema.validation,omitempty"`
	ConfluentValueSubjectNameStrategy *BreakingChangeRule `json:"confluent.value.subject.name.strategy,omitempty" yaml:"confluent.value.subject.name.strategy,omitempty"`
}

// KafkaOperationBindingRules defines breaking rules for the Kafka Operation Binding properties.
type KafkaOperationBindingRules struct {
	GroupID        *BreakingChangeRule `json:"groupId,omitempty" yaml:"groupId,omitempty"`
	ClientID       *BreakingChangeRule `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// KafkaMessageBindingRules defines breaking rules for the Kafka Message Binding properties.
type KafkaMessageBindingRules struct {
	Key                     *BreakingChangeRule `json:"key,omitempty" yaml:"key,omitempty"`
	SchemaIDLocation        *BreakingChangeRule `json:"schemaIdLocation,omitempty" yaml:"schemaIdLocation,omitempty"`
	SchemaIDPayloadEncoding *BreakingChangeRule `json:"schemaIdPayloadEncoding,omitempty" yaml:"schemaIdPayloadEncoding,omitempty"`
	SchemaLookupStrategy    *BreakingChangeRule `json:"schemaLookupStrategy,omitempty" yaml:"schemaLookupStrategy,omitempty"`
	BindingVersion          *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// WebSocketChannelBindingRules defines breaking rules for the WebSocket Channel Binding properties.
type WebSocketChannelBindingRules struct {
	Method         *BreakingChangeRule `json:"method,omitempty" yaml:"method,omitempty"`
	Query          *BreakingChangeRule `json:"query,omitempty" yaml:"query,omitempty"`
	Headers        *BreakingChangeRule `json:"headers,omitempty" yaml:"headers,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// AMQPChannelBindingRules defines breaking rules for the AMQP Channel Binding properties.
type AMQPChannelBindingRules struct {
	Is             *BreakingChangeRule `json:"is,omitempty" yaml:"is,omitempty"`
	Exchange       *BreakingChangeRule `json:"exchange,omitempty" yaml:"exchange,omitempty"`
	Queue          *BreakingChangeRule `json:"queue,omitempty" yaml:"queue,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// AMQPExchangeRules defines breaking rules for the AMQP Exchange properties.
type AMQPExchangeRules struct {
	Name       *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Type       *BreakingChangeRule `json:"type,omitempty" yaml:"type,omitempty"`
	Durable    *BreakingChangeRule `json:"durable,omitempty" yaml:"durable,omitempty"`
	AutoDelete *BreakingChangeRule `json:"autoDelete,omitempty" yaml:"autoDelete,omitempty"`
	VHost      *BreakingChangeRule `json:"vhost,omitempty" yaml:"vhost,omitempty"`
}

// AMQPQueueRules defines breaking rules for the AMQP Queue properties.
type AMQPQueueRules struct {
	Name       *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	Durable    *BreakingChangeRule `json:"durable,omitempty" yaml:"durable,omitempty"`
	Exclusive  *BreakingChangeRule `json:"exclusive,omitempty" yaml:"exclusive,omitempty"`
	AutoDelete *BreakingChangeRule `json:"autoDelete,omitempty" yaml:"autoDelete,omitempty"`
	VHost      *BreakingChangeRule `json:"vhost,omitempty" yaml:"vhost,omitempty"`
}

// AMQPOperationBindingRules defines breaking rules for the AMQP Operation Binding properties.
type AMQPOperationBindingRules struct {
	Expiration     *BreakingChangeRule `json:"expiration,omitempty" yaml:"expiration,omitempty"`
	UserID         *BreakingChangeRule `json:"userId,omitempty" yaml:"userId,omitempty"`
	CC             *BreakingChangeRule `json:"cc,omitempty" yaml:"cc,omitempty"`
	Priority       *BreakingChangeRule `json:"priority,omitempty" yaml:"priority,omitempty"`
	DeliveryMode   *BreakingChangeRule `json:"deliveryMode,omitempty" yaml:"deliveryMode,omitempty"`
	Mandatory      *BreakingChangeRule `json:"mandatory,omitempty" yaml:"mandatory,omitempty"`
	BCC            *BreakingChangeRule `json:"bcc,omitempty" yaml:"bcc,omitempty"`
	Timestamp      *BreakingChangeRule `json:"timestamp,omitempty" yaml:"timestamp,omitempty"`
	Ack            *BreakingChangeRule `json:"ack,omitempty" yaml:"ack,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// AMQPMessageBindingRules defines breaking rules for the AMQP Message Binding properties.
type AMQPMessageBindingRules struct {
	ContentEncoding *BreakingChangeRule `json:"contentEncoding,omitempty" yaml:"contentEncoding,omitempty"`
	MessageType     *BreakingChangeRule `json:"messageType,omitempty" yaml:"messageType,omitempty"`
	BindingVersion  *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// MQTTServerBindingRules defines breaking rules for the MQTT Server Binding properties.
type MQTTServerBindingRules struct {
	ClientID              *BreakingChangeRule `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	CleanSession          *BreakingChangeRule `json:"cleanSession,omitempty" yaml:"cleanSession,omitempty"`
	LastWill              *BreakingChangeRule `json:"lastWill,omitempty" yaml:"lastWill,omitempty"`
	KeepAlive             *BreakingChangeRule `json:"keepAlive,omitempty" yaml:"keepAlive,omitempty"`
	SessionExpiryInterval *BreakingChangeRule `json:"sessionExpiryInterval,omitempty" yaml:"sessionExpiryInterval,omitempty"`
	MaximumPacketSize     *BreakingChangeRule `json:"maximumPacketSize,omitempty" yaml:"maximumPacketSize,omitempty"`
	BindingVersion        *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// MQTTLastWillRules defines breaking rules for the MQTT Last Will properties.
type MQTTLastWillRules struct {
	Topic   *BreakingChangeRule `json:"topic,omitempty" yaml:"topic,omitempty"`
	QoS     *BreakingChangeRule `json:"qos,omitempty" yaml:"qos,omitempty"`
	Message *BreakingChangeRule `json:"message,omitempty" yaml:"message,omitempty"`
	Retain  *BreakingChangeRule `json:"retain,omitempty" yaml:"retain,omitempty"`
}

// MQTTOperationBindingRules defines breaking rules for the MQTT Operation Binding properties.
type MQTTOperationBindingRules struct {
	QoS            *BreakingChangeRule `json:"qos,omitempty" yaml:"qos,omitempty"`
	Retain         *BreakingChangeRule `json:"retain,omitempty" yaml:"retain,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// MQTTMessageBindingRules defines breaking rules for the MQTT Message Binding properties.
type MQTTMessageBindingRules struct {
	PayloadFormatIndicator *BreakingChangeRule `json:"payloadFormatIndicator,omitempty" yaml:"payloadFormatIndicator,omitempty"`
	CorrelationData        *BreakingChangeRule `json:"correlationData,omitempty" yaml:"correlationData,omitempty"`
	ContentType            *BreakingChangeRule `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	ResponseTopic          *BreakingChangeRule `json:"responseTopic,omitempty" yaml:"responseTopic,omitempty"`
	BindingVersion         *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// SQSChannelBindingRules defines breaking rules for the SQS Channel Binding properties.
type SQSChannelBindingRules struct {
	Queue           *BreakingChangeRule `json:"queue,omitempty" yaml:"queue,omitempty"`
	DeadLetterQueue *BreakingChangeRule `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	BindingVersion  *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// SQSOperationBindingRules defines breaking rules for the SQS Operation Binding properties.
type SQSOperationBindingRules struct {
	Queues         *BreakingChangeRule `json:"queues,omitempty" yaml:"queues,omitempty"`
	BindingVersion *BreakingChangeRule `json:"bindingVersion,omitempty" yaml:"bindingVersion,omitempty"`
}

// SQSQueueRules defines breaking rules for the SQS Queue properties.
type SQSQueueRules struct {
	Name                   *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	ARN                    *BreakingChangeRule `json:"arn,omitempty" yaml:"arn,omitempty"`
	FifoQueue              *BreakingChangeRule `json:"fifoQueue,omitempty" yaml:"fifoQueue,omitempty"`
	DeduplicationScope     *BreakingChangeRule `json:"deduplicationScope,omitempty" yaml:"deduplicationScope,omitempty"`
	FifoThroughputLimit    *BreakingChangeRule `json:"fifoThroughputLimit,omitempty" yaml:"fifoThroughputLimit,omitempty"`
	DeliveryDelay          *BreakingChangeRule `json:"deliveryDelay,omitempty" yaml:"deliveryDelay,omitempty"`
	VisibilityTimeout      *BreakingChangeRule `json:"visibilityTimeout,omitempty" yaml:"visibilityTimeout,omitempty"`
	ReceiveMessageWaitTime *BreakingChangeRule `json:"receiveMessageWaitTime,omitempty" yaml:"receiveMessageWaitTime,omitempty"`
	MessageRetentionPeriod *BreakingChangeRule `json:"messageRetentionPeriod,omitempty" yaml:"messageRetentionPeriod,omitempty"`
	RedrivePolicy          *BreakingChangeRule `json:"redrivePolicy,omitempty" yaml:"redrivePolicy,omitempty"`
	Policy                 *BreakingChangeRule `json:"policy,omitempty" yaml:"policy,omitempty"`
	Tags                   *BreakingChangeRule `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// SQSIdentifierRules defines breaking rules for the SQS Identifier properties.
type SQSIdentifierRules struct {
	Name      *BreakingChangeRule `json:"name,omitempty" yaml:"name,omitempty"`
	ARN       *BreakingChangeRule `json:"arn,omitempty" yaml:"arn,omitempty"`
	FifoQueue *BreakingChangeRule `json:"fifoQueue,omitempty" yaml:"fifoQueue,omitempty"`
}

// SQSRedrivePolicyRules defines breaking rules for the SQS Redrive Policy properties.
type SQSRedrivePolicyRules struct {
	DeadLetterQueue *BreakingChangeRule `json:"deadLetterQueue,omitempty" yaml:"deadLetterQueue,omitempty"`
	MaxReceiveCount *BreakingChangeRule `json:"maxReceiveCount,omitempty" yaml:"maxReceiveCount,omitempty"`
}

// SQSPolicyRules defines breaking rules for the SQS Policy properties.
type SQSPolicyRules struct {
	Statements *BreakingChangeRule `json:"statements,omitempty" yaml:"statements,omitempty"`
}

// SQSPolicyStatementRules defines breaking rules for the SQS Policy Statement properties.
type SQSPolicyStatementRules struct {
	Effect    *BreakingChangeRule `json:"effect,omitempty" yaml:"effect,omitempty"`
	Principal *BreakingChangeRule `json:"principal,omitempty" yaml:"principal,omitempty"`
	Action    *BreakingChangeRule `json:"action,omitempty" yaml:"action,omitempty"`
	Resource  *BreakingChangeRule `json:"resource,omitempty" yaml:"resource,omitempty"`
	Condition *BreakingChangeRule `json:"condition,omitempty" yaml:"condition,omitempty"`
}
