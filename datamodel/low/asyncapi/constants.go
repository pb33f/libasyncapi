// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

// Constants for labels used to look up values within AsyncAPI 3.0 specifications.
const (
	// Root-level fields
	AsyncAPILabel           = "asyncapi"
	IDLabel                 = "id"
	DefaultContentTypeLabel = "defaultContentType"

	// Channel fields
	AddressLabel  = "address"
	MessagesLabel = "messages"
	ChannelsLabel = "channels"

	// Operation fields
	ActionLabel     = "action"
	ChannelLabel    = "channel"
	OperationsLabel = "operations"
	TraitsLabel     = "traits"
	ReplyLabel      = "reply"

	// Message fields
	HeadersLabel       = "headers"
	PayloadLabel       = "payload"
	CorrelationIDLabel = "correlationId"
	ContentTypeLabel   = "contentType"

	// Server fields
	HostLabel            = "host"
	ProtocolLabel        = "protocol"
	ProtocolVersionLabel = "protocolVersion"
	PathnameLabel        = "pathname"
	VariablesLabel       = "variables"

	// Parameter fields
	EnumLabel     = "enum"
	DefaultLabel  = "default"
	LocationLabel = "location"

	// Multi-format schema fields
	SchemaFormatLabel = "schemaFormat"
	SchemaLabel       = "schema"

	// Security fields
	TypeLabel             = "type"
	SchemeLabel           = "scheme"
	BearerFormatLabel     = "bearerFormat"
	FlowsLabel            = "flows"
	OpenIDConnectURLLabel = "openIdConnectUrl"
	ScopesLabel           = "scopes"
	InLabel               = "in"
	AvailableScopesLabel  = "availableScopes"

	// OAuth flow labels
	ImplicitLabel          = "implicit"
	PasswordLabel          = "password"
	ClientCredentialsLabel = "clientCredentials"
	AuthorizationCodeLabel = "authorizationCode"
	AuthorizationURLLabel  = "authorizationUrl"
	TokenURLLabel          = "tokenUrl"
	RefreshURLLabel        = "refreshUrl"

	// Common fields
	TitleLabel        = "title"
	SummaryLabel      = "summary"
	DescriptionLabel  = "description"
	TagsLabel         = "tags"
	ExternalDocsLabel = "externalDocs"
	ServersLabel      = "servers"
	ParametersLabel   = "parameters"
	SecurityLabel     = "security"
	BindingsLabel     = "bindings"
	ComponentsLabel   = "components"
	InfoLabel         = "info"
	NameLabel         = "name"
	ExamplesLabel     = "examples"
	VersionLabel      = "version"

	// Binding labels
	HTTPLabel           = "http"
	KafkaLabel          = "kafka"
	WebSocketLabel      = "ws"
	AMQPLabel           = "amqp"
	MQTTLabel           = "mqtt"
	SQSLabel            = "sqs"
	BindingVersionLabel = "bindingVersion"

	// HTTP binding fields
	QueryLabel      = "query"
	MethodLabel     = "method"
	StatusCodeLabel = "statusCode"

	// Kafka binding fields
	TopicLabel                   = "topic"
	PartitionsLabel              = "partitions"
	ReplicasLabel                = "replicas"
	TopicConfigurationLabel      = "topicConfiguration"
	SchemaRegistryURLLabel       = "schemaRegistryUrl"
	SchemaRegistryVendorLabel    = "schemaRegistryVendor"
	GroupIDLabel                 = "groupId"
	ClientIDLabel                = "clientId"
	KeyLabel                     = "key"
	SchemaIDLocationLabel        = "schemaIdLocation"
	SchemaIDPayloadEncodingLabel = "schemaIdPayloadEncoding"
	SchemaLookupStrategyLabel    = "schemaLookupStrategy"

	// Kafka topic configuration fields
	CleanupPolicyLabel                     = "cleanup.policy"
	RetentionMsLabel                       = "retention.ms"
	RetentionBytesLabel                    = "retention.bytes"
	DeleteRetentionMsLabel                 = "delete.retention.ms"
	MaxMessageBytesLabel                   = "max.message.bytes"
	ConfluentKeySchemaValidationLabel      = "confluent.key.schema.validation"
	ConfluentKeySubjectNameStrategyLabel   = "confluent.key.subject.name.strategy"
	ConfluentValueSchemaValidationLabel    = "confluent.value.schema.validation"
	ConfluentValueSubjectNameStrategyLabel = "confluent.value.subject.name.strategy"

	// AMQP binding fields
	IsLabel              = "is"
	ExchangeLabel        = "exchange"
	QueueLabel           = "queue"
	DurableLabel         = "durable"
	AutoDeleteLabel      = "autoDelete"
	ExclusiveLabel       = "exclusive"
	VHostLabel           = "vhost"
	ExpirationLabel      = "expiration"
	UserIDLabel          = "userId"
	CCLabel              = "cc"
	PriorityLabel        = "priority"
	DeliveryModeLabel    = "deliveryMode"
	MandatoryLabel       = "mandatory"
	BCCLabel             = "bcc"
	TimestampLabel       = "timestamp"
	AckLabel             = "ack"
	ContentEncodingLabel = "contentEncoding"
	MessageTypeLabel     = "messageType"

	// MQTT binding fields
	CleanSessionLabel           = "cleanSession"
	LastWillLabel               = "lastWill"
	KeepAliveLabel              = "keepAlive"
	SessionExpiryIntervalLabel  = "sessionExpiryInterval"
	MaximumPacketSizeLabel      = "maximumPacketSize"
	QoSLabel                    = "qos"
	RetainLabel                 = "retain"
	MessageExpiryIntervalLabel  = "messageExpiryInterval"
	PayloadFormatIndicatorLabel = "payloadFormatIndicator"
	CorrelationDataLabel        = "correlationData"
	ResponseTopicLabel          = "responseTopic"

	// SQS binding fields
	QueuesLabel                 = "queues"
	DeadLetterQueueLabel        = "deadLetterQueue"
	FifoQueueLabel              = "fifoQueue"
	DeduplicationScopeLabel     = "deduplicationScope"
	FifoThroughputLimitLabel    = "fifoThroughputLimit"
	DeliveryDelayLabel          = "deliveryDelay"
	VisibilityTimeoutLabel      = "visibilityTimeout"
	ReceiveMessageWaitTimeLabel = "receiveMessageWaitTime"
	MessageRetentionPeriodLabel = "messageRetentionPeriod"
	RedrivePolicyLabel          = "redrivePolicy"
	PolicyLabel                 = "policy"
	StatementsLabel             = "statements"
	PrincipalLabel              = "principal"
	ResourceLabel               = "resource"
	ConditionLabel              = "condition"
	ARNLabel                    = "arn"

	// Reference field
	RefLabel = "$ref"

	// Operation action values
	ActionSend    = "send"
	ActionReceive = "receive"

	// Components section labels
	SchemasLabel           = "schemas"
	SecuritySchemesLabel   = "securitySchemes"
	ServerVariablesLabel   = "serverVariables"
	CorrelationIDsLabel    = "correlationIds"
	RepliesLabel           = "replies"
	ReplyAddressesLabel    = "replyAddresses"
	OperationTraitsLabel   = "operationTraits"
	MessageTraitsLabel     = "messageTraits"
	ServerBindingsLabel    = "serverBindings"
	ChannelBindingsLabel   = "channelBindings"
	OperationBindingsLabel = "operationBindings"
	MessageBindingsLabel   = "messageBindings"

	// Contact/License fields (reused from base)
	ContactLabel = "contact"
	LicenseLabel = "license"
	EmailLabel   = "email"
	URLLabel     = "url"

	// Terms of service
	TermsOfServiceLabel = "termsOfService"
)
