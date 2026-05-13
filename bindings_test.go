// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

func TestComprehensiveBindings_KafkaServerBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	server, ok := model.Servers.Get("kafka-server")
	require.True(t, ok)
	require.NotNil(t, server.Bindings)
	require.NotNil(t, server.Bindings.Kafka)

	kafka := server.Bindings.Kafka
	assert.Equal(t, "https://schema-registry.example.com", kafka.SchemaRegistryURL)
	assert.Equal(t, "confluent", kafka.SchemaRegistryVendor)
	assert.Equal(t, "0.5.0", kafka.BindingVersion)
	require.NotNil(t, kafka.Extensions)
	ext, ok := kafka.Extensions.Get("x-custom")
	require.True(t, ok)
	assert.NotNil(t, ext)
}

func TestComprehensiveBindings_MQTTServerBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	server, ok := model.Servers.Get("mqtt-server")
	require.True(t, ok)
	require.NotNil(t, server.Bindings)
	require.NotNil(t, server.Bindings.MQTT)

	mqtt := server.Bindings.MQTT
	assert.Equal(t, "test-client-id", mqtt.ClientID)
	assert.True(t, mqtt.CleanSession)
	assert.Equal(t, 120, mqtt.KeepAlive)
	assert.Equal(t, 3600, mqtt.SessionExpiryInterval)
	assert.Equal(t, 65535, mqtt.MaximumPacketSize)
	assert.Equal(t, "0.2.0", mqtt.BindingVersion)

	// Test LastWill
	require.NotNil(t, mqtt.LastWill)
	assert.Equal(t, "/will/topic", mqtt.LastWill.Topic)
	assert.Equal(t, 1, mqtt.LastWill.QoS)
	assert.Equal(t, "Client disconnected unexpectedly", mqtt.LastWill.Message)
	assert.True(t, mqtt.LastWill.Retain)
}

func TestComprehensiveBindings_KafkaChannelBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	channel, ok := model.Channels.Get("kafka-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.Kafka)

	kafka := channel.Bindings.Kafka
	assert.Equal(t, "test-topic", kafka.Topic)
	assert.Equal(t, 24, kafka.Partitions)
	assert.Equal(t, 3, kafka.Replicas)
	assert.Equal(t, "0.5.0", kafka.BindingVersion)

	// Test TopicConfiguration exists (note: individual fields not yet extracted by low-level Build)
	require.NotNil(t, kafka.TopicConfiguration)
	require.NotNil(t, kafka.TopicConfiguration.GoLow())
}

func TestComprehensiveBindings_AMQPChannelBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	channel, ok := model.Channels.Get("amqp-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.AMQP)

	amqp := channel.Bindings.AMQP
	assert.Equal(t, "routingKey", amqp.Is)
	assert.Equal(t, "0.3.0", amqp.BindingVersion)

	// Test Exchange
	require.NotNil(t, amqp.Exchange)
	assert.Equal(t, "orders-exchange", amqp.Exchange.Name)
	assert.Equal(t, "topic", amqp.Exchange.Type)
	assert.True(t, amqp.Exchange.Durable)
	assert.False(t, amqp.Exchange.AutoDelete)
	assert.Equal(t, "/test", amqp.Exchange.VHost)

	// Test Queue
	require.NotNil(t, amqp.Queue)
	assert.Equal(t, "orders-queue", amqp.Queue.Name)
	assert.True(t, amqp.Queue.Durable)
	assert.False(t, amqp.Queue.Exclusive)
	assert.False(t, amqp.Queue.AutoDelete)
	assert.Equal(t, "/test", amqp.Queue.VHost)
}

func TestComprehensiveBindings_WebSocketChannelBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	channel, ok := model.Channels.Get("ws-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.WebSocket)

	ws := channel.Bindings.WebSocket
	assert.Equal(t, "GET", ws.Method)
	assert.Equal(t, "0.1.0", ws.BindingVersion)
	require.NotNil(t, ws.Query)
	require.NotNil(t, ws.Headers)
}

func TestComprehensiveBindings_KafkaOperationBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	op, ok := model.Operations.Get("kafka-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.Kafka)

	kafka := op.Bindings.Kafka
	assert.Equal(t, "0.5.0", kafka.BindingVersion)
	require.NotNil(t, kafka.GroupID)
	require.NotNil(t, kafka.ClientID)
}

func TestComprehensiveBindings_AMQPOperationBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	op, ok := model.Operations.Get("amqp-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.AMQP)

	amqp := op.Bindings.AMQP
	assert.Equal(t, 30000, amqp.Expiration)
	assert.Equal(t, "guest", amqp.UserID)
	assert.Equal(t, 10, amqp.Priority)
	assert.Equal(t, 2, amqp.DeliveryMode)
	assert.True(t, amqp.Mandatory)
	assert.True(t, amqp.Timestamp)
	assert.True(t, amqp.Ack)
	assert.Equal(t, "0.3.0", amqp.BindingVersion)

	// Test CC and BCC arrays
	require.Len(t, amqp.CC, 2)
	assert.Contains(t, amqp.CC, "routing.key.1")
	assert.Contains(t, amqp.CC, "routing.key.2")

	require.Len(t, amqp.BCC, 2)
	assert.Contains(t, amqp.BCC, "hidden.key.1")
	assert.Contains(t, amqp.BCC, "hidden.key.2")
}

func TestComprehensiveBindings_MQTTOperationBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	op, ok := model.Operations.Get("mqtt-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.MQTT)

	mqtt := op.Bindings.MQTT
	assert.Equal(t, 2, mqtt.QoS)
	assert.True(t, mqtt.Retain)
	assert.Equal(t, "0.2.0", mqtt.BindingVersion)
}

func TestComprehensiveBindings_HTTPOperationBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	op, ok := model.Operations.Get("http-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.HTTP)

	http := op.Bindings.HTTP
	assert.Equal(t, "POST", http.Method)
	assert.Equal(t, "0.3.0", http.BindingVersion)
	require.NotNil(t, http.Query)
}

func TestComprehensiveBindings_KafkaMessageBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	msg, ok := model.Components.Messages.Get("kafkaMessage")
	require.True(t, ok)
	require.NotNil(t, msg.Bindings)
	require.NotNil(t, msg.Bindings.Kafka)

	kafka := msg.Bindings.Kafka
	require.NotNil(t, kafka.Key)
	assert.Equal(t, "header", kafka.SchemaIDLocation)
	assert.Equal(t, "confluent", kafka.SchemaIDPayloadEncoding)
	assert.Equal(t, "TopicIdStrategy", kafka.SchemaLookupStrategy)
	assert.Equal(t, "0.5.0", kafka.BindingVersion)
}

func TestComprehensiveBindings_AMQPMessageBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	msg, ok := model.Components.Messages.Get("amqpMessage")
	require.True(t, ok)
	require.NotNil(t, msg.Bindings)
	require.NotNil(t, msg.Bindings.AMQP)

	amqp := msg.Bindings.AMQP
	assert.Equal(t, "utf-8", amqp.ContentEncoding)
	assert.Equal(t, "order.created", amqp.MessageType)
	assert.Equal(t, "0.3.0", amqp.BindingVersion)
}

func TestComprehensiveBindings_MQTTMessageBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	msg, ok := model.Components.Messages.Get("mqttMessage")
	require.True(t, ok)
	require.NotNil(t, msg.Bindings)
	require.NotNil(t, msg.Bindings.MQTT)

	mqtt := msg.Bindings.MQTT
	assert.Equal(t, 1, mqtt.PayloadFormatIndicator)
	assert.Equal(t, "application/json", mqtt.ContentType)
	assert.Equal(t, "/response/topic", mqtt.ResponseTopic)
	assert.Equal(t, "0.2.0", mqtt.BindingVersion)
	require.NotNil(t, mqtt.CorrelationData)
}

func TestComprehensiveBindings_SQSServerBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	server, ok := model.Servers.Get("sqs-server")
	require.True(t, ok)
	require.NotNil(t, server.Bindings)
	require.NotNil(t, server.Bindings.SQS)
	require.NotNil(t, server.Bindings.SQS.Extensions)

	ext, ok := server.Bindings.SQS.Extensions.Get("x-server-binding-ext")
	require.True(t, ok)
	assert.NotNil(t, ext)
}

func TestComprehensiveBindings_SQSChannelBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	channel, ok := model.Channels.Get("sqs-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.SQS)

	sqs := channel.Bindings.SQS
	assert.Equal(t, "0.3.0", sqs.BindingVersion)
	require.NotNil(t, sqs.Queue)
	assert.Equal(t, "findings-worker", sqs.Queue.Name)
	assert.False(t, sqs.Queue.FifoQueue)
	assert.Equal(t, 60, sqs.Queue.VisibilityTimeout)
	assert.Equal(t, 20, sqs.Queue.ReceiveMessageWaitTime)
	assert.Equal(t, 345600, sqs.Queue.MessageRetentionPeriod)

	require.NotNil(t, sqs.Queue.RedrivePolicy)
	require.NotNil(t, sqs.Queue.RedrivePolicy.DeadLetterQueue)
	assert.Equal(t, "findings-worker-dlq", sqs.Queue.RedrivePolicy.DeadLetterQueue.Name)
	assert.Equal(t, 5, sqs.Queue.RedrivePolicy.MaxReceiveCount)

	require.NotNil(t, sqs.Queue.Policy)
	require.Len(t, sqs.Queue.Policy.Statements, 1)
	statement := sqs.Queue.Policy.Statements[0]
	assert.Equal(t, "Allow", statement.Effect)
	require.NotNil(t, statement.Principal)
	assert.Equal(t, yaml.MappingNode, statement.Principal.Kind)
	require.NotNil(t, statement.Action)
	assert.Equal(t, yaml.SequenceNode, statement.Action.Kind)
	require.Len(t, statement.Action.Content, 1)
	assert.Equal(t, "sqs:SendMessage", statement.Action.Content[0].Value)
	require.NotNil(t, statement.Resource)
	assert.Equal(t, yaml.ScalarNode, statement.Resource.Kind)
	assert.Equal(t, "arn:aws:sqs:us-east-1:123456789012:findings-worker", statement.Resource.Value)
	require.NotNil(t, statement.Condition)
	assert.Equal(t, yaml.MappingNode, statement.Condition.Kind)

	require.NotNil(t, sqs.Queue.Tags)
	assert.Equal(t, "linus", sqs.Queue.Tags["owner"])
	assert.Equal(t, "findings", sqs.Queue.Tags["service"])

	require.NotNil(t, sqs.DeadLetterQueue)
	assert.Equal(t, "findings-worker-dlq", sqs.DeadLetterQueue.Name)
	assert.False(t, sqs.DeadLetterQueue.FifoQueue)
	assert.Equal(t, 1209600, sqs.DeadLetterQueue.MessageRetentionPeriod)

	require.NotNil(t, sqs.Extensions)
	ext, ok := sqs.Extensions.Get("x-sqs-channel-ext")
	require.True(t, ok)
	assert.NotNil(t, ext)
}

func TestComprehensiveBindings_SQSFIFOChannelBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	channel, ok := model.Channels.Get("sqs-fifo-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.SQS)
	require.NotNil(t, channel.Bindings.SQS.Queue)

	queue := channel.Bindings.SQS.Queue
	assert.Equal(t, "audit-events.fifo", queue.Name)
	assert.True(t, queue.FifoQueue)
	assert.Equal(t, "messageGroup", queue.DeduplicationScope)
	assert.Equal(t, "perMessageGroupId", queue.FifoThroughputLimit)
	assert.Equal(t, 15, queue.DeliveryDelay)
	assert.Equal(t, 45, queue.VisibilityTimeout)
	assert.Equal(t, 10, queue.ReceiveMessageWaitTime)
	assert.Equal(t, 86400, queue.MessageRetentionPeriod)
	assert.Equal(t, "fifo", queue.Tags["queue_kind"])
}

func TestComprehensiveBindings_SQSOperationBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	op, ok := model.Operations.Get("sqs-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.SQS)

	sqs := op.Bindings.SQS
	assert.Equal(t, "0.3.0", sqs.BindingVersion)
	require.Len(t, sqs.Queues, 2)
	assert.Equal(t, "findings-worker", sqs.Queues[0].Name)
	assert.False(t, sqs.Queues[0].FifoQueue)
	assert.Equal(t, 75, sqs.Queues[0].VisibilityTimeout)
	assert.Equal(t, 12, sqs.Queues[0].ReceiveMessageWaitTime)
	assert.Equal(t, 172800, sqs.Queues[0].MessageRetentionPeriod)
	require.NotNil(t, sqs.Queues[0].RedrivePolicy)
	require.NotNil(t, sqs.Queues[0].RedrivePolicy.DeadLetterQueue)
	assert.Equal(t, "findings-worker-dlq", sqs.Queues[0].RedrivePolicy.DeadLetterQueue.Name)
	assert.Equal(t, 7, sqs.Queues[0].RedrivePolicy.MaxReceiveCount)
	require.NotNil(t, sqs.Queues[0].Policy)
	require.Len(t, sqs.Queues[0].Policy.Statements, 1)
	assert.Equal(t, "Allow", sqs.Queues[0].Policy.Statements[0].Effect)
	require.NotNil(t, sqs.Queues[0].Tags)
	assert.Equal(t, "receive", sqs.Queues[0].Tags["operation_scope"])
	assert.Equal(t, "findings-worker-dlq", sqs.Queues[1].Name)
	assert.False(t, sqs.Queues[1].FifoQueue)
	assert.Equal(t, 1209600, sqs.Queues[1].MessageRetentionPeriod)
}

func TestComprehensiveBindings_SQSMessageBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	msg, ok := model.Components.Messages.Get("sqsMessage")
	require.True(t, ok)
	require.NotNil(t, msg.Bindings)
	require.NotNil(t, msg.Bindings.SQS)
	require.NotNil(t, msg.Bindings.SQS.Extensions)

	ext, ok := msg.Bindings.SQS.Extensions.Get("x-sqs-message-ext")
	require.True(t, ok)
	assert.NotNil(t, ext)
}

func TestComprehensiveBindings_SQSGoLowHashAndLineNumbers(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	server, _ := model.Servers.Get("sqs-server")
	require.NotNil(t, server.Bindings.SQS.GoLow())
	require.NotNil(t, server.Bindings.SQS.GoLowUntyped())
	require.NotNil(t, server.Bindings.SQS.GoLow().GetRootNode())
	assert.Greater(t, server.Bindings.SQS.GoLow().GetRootNode().Line, 0)

	channel, _ := model.Channels.Get("sqs-channel")
	require.NotNil(t, channel.Bindings.SQS.GoLow())
	require.NotNil(t, channel.Bindings.SQS.Queue.GoLow())
	require.NotNil(t, channel.Bindings.SQS.Queue.RedrivePolicy.GoLow())
	require.NotNil(t, channel.Bindings.SQS.Queue.Policy.GoLow())
	require.NotNil(t, channel.Bindings.SQS.DeadLetterQueue.GoLow())

	hash1 := channel.Bindings.SQS.GoLow().Hash()
	hash2 := channel.Bindings.SQS.GoLow().Hash()
	assert.Equal(t, hash1, hash2)

	op, _ := model.Operations.Get("sqs-operation")
	require.NotNil(t, op.Bindings.SQS.GoLow())

	msg, _ := model.Components.Messages.Get("sqsMessage")
	require.NotNil(t, msg.Bindings.SQS.GoLow())
}

func TestComprehensiveBindings_HTTPMessageBinding(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	msg, ok := model.Components.Messages.Get("httpMessage")
	require.True(t, ok)
	require.NotNil(t, msg.Bindings)
	require.NotNil(t, msg.Bindings.HTTP)

	http := msg.Bindings.HTTP
	assert.Equal(t, 200, http.StatusCode)
	assert.Equal(t, "0.3.0", http.BindingVersion)
	require.NotNil(t, http.Headers)
}

func TestComprehensiveBindings_GoLowMethods(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Test GoLow returns non-nil for all binding types
	kafkaServer, _ := model.Servers.Get("kafka-server")
	require.NotNil(t, kafkaServer.Bindings.Kafka.GoLow())
	require.NotNil(t, kafkaServer.Bindings.Kafka.GoLowUntyped())

	mqttServer, _ := model.Servers.Get("mqtt-server")
	require.NotNil(t, mqttServer.Bindings.MQTT.GoLow())
	require.NotNil(t, mqttServer.Bindings.MQTT.GoLowUntyped())
	require.NotNil(t, mqttServer.Bindings.MQTT.LastWill.GoLow())
	require.NotNil(t, mqttServer.Bindings.MQTT.LastWill.GoLowUntyped())

	kafkaChannel, _ := model.Channels.Get("kafka-channel")
	require.NotNil(t, kafkaChannel.Bindings.Kafka.GoLow())
	require.NotNil(t, kafkaChannel.Bindings.Kafka.TopicConfiguration.GoLow())

	amqpChannel, _ := model.Channels.Get("amqp-channel")
	require.NotNil(t, amqpChannel.Bindings.AMQP.GoLow())
	require.NotNil(t, amqpChannel.Bindings.AMQP.Exchange.GoLow())
	require.NotNil(t, amqpChannel.Bindings.AMQP.Queue.GoLow())

	wsChannel, _ := model.Channels.Get("ws-channel")
	require.NotNil(t, wsChannel.Bindings.WebSocket.GoLow())
}

func TestComprehensiveBindings_HashMethods(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Test Hash methods return consistent results
	kafkaServer, _ := model.Servers.Get("kafka-server")
	hash1 := kafkaServer.Bindings.Kafka.GoLow().Hash()
	hash2 := kafkaServer.Bindings.Kafka.GoLow().Hash()
	assert.Equal(t, hash1, hash2, "Hash should be consistent")

	mqttServer, _ := model.Servers.Get("mqtt-server")
	hash3 := mqttServer.Bindings.MQTT.GoLow().Hash()
	hash4 := mqttServer.Bindings.MQTT.GoLow().Hash()
	assert.Equal(t, hash3, hash4, "Hash should be consistent")

	// Different bindings should have different hashes
	assert.NotEqual(t, hash1, hash3, "Different bindings should have different hashes")
}

func TestComprehensiveBindings_LineNumbers(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Test that low-level objects have line number information
	kafkaServer, _ := model.Servers.Get("kafka-server")
	lowKafka := kafkaServer.Bindings.Kafka.GoLow()
	require.NotNil(t, lowKafka.GetRootNode())
	assert.Greater(t, lowKafka.GetRootNode().Line, 0, "Line number should be set")

	mqttServer, _ := model.Servers.Get("mqtt-server")
	lowMQTT := mqttServer.Bindings.MQTT.GoLow()
	require.NotNil(t, lowMQTT.GetRootNode())
	assert.Greater(t, lowMQTT.GetRootNode().Line, 0, "Line number should be set")
}

func TestComprehensiveBindings_SQSRenderRoundTrip(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	rendered, err := doc.Render()
	require.NoError(t, err)

	output := string(rendered)
	assert.Contains(t, output, "sqs:")
	assert.Contains(t, output, "findings-worker")
	assert.Contains(t, output, "deadLetterQueue:")
	assert.GreaterOrEqual(t, strings.Count(output, "fifoQueue: false"), 5)
	assert.Contains(t, output, "bindingVersion: 0.3.0")
	assert.Contains(t, output, "audit-events.fifo")
	assert.Contains(t, output, "x-server-binding-ext: value")
	assert.Contains(t, output, "x-sqs-channel-ext: value")
	assert.Contains(t, output, "x-sqs-message-ext: value")
}
