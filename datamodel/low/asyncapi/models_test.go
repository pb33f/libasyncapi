// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadFixtureNode(t *testing.T, name string) *yaml.Node {
	t.Helper()
	path := filepath.Join("..", "..", "..", "test_fixtures", name)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	return loadSpecNode(t, data)
}

func loadSpecNode(t *testing.T, spec []byte) *yaml.Node {
	t.Helper()
	var root yaml.Node
	require.NoError(t, yaml.Unmarshal(spec, &root))
	require.Equal(t, yaml.DocumentNode, root.Kind)
	require.NotEmpty(t, root.Content)
	return &root
}

func createDocument(t *testing.T, root *yaml.Node, config *DocumentConfiguration) *AsyncAPI {
	t.Helper()
	doc, err := CreateDocumentWithConfig(root, config)
	require.NoError(t, err)
	require.NotNil(t, doc)
	return doc
}

func assertHashNotZero(t *testing.T, hash uint64) {
	t.Helper()
	assert.NotEqual(t, uint64(0), hash)
}

func TestCreateDocument_ComprehensiveBindings(t *testing.T) {
	root := loadFixtureNode(t, "comprehensive-bindings.yaml")
	doc := createDocument(t, root, NewDocumentConfiguration())

	assert.Equal(t, "3.0.0", doc.AsyncAPI.Value)
	assert.NotNil(t, doc.GetRootNode())
	assertHashNotZero(t, doc.Hash())

	require.NotNil(t, doc.Info.Value)
	assertHashNotZero(t, doc.Info.Value.Hash())

	require.NotNil(t, doc.Servers.Value)
	kafkaServer := low.FindItemInOrderedMap("kafka-server", doc.Servers.Value)
	require.NotNil(t, kafkaServer)
	assertHashNotZero(t, kafkaServer.Value.Hash())

	kafkaBinding := kafkaServer.Value.Bindings.Value.Kafka.Value
	require.NotNil(t, kafkaBinding)
	assert.NotNil(t, kafkaBinding.GetRootNode())
	assert.NotNil(t, kafkaBinding.GetContext())
	assert.NotNil(t, kafkaBinding.GetIndex())
	assert.NotNil(t, low.FindItemInOrderedMap("x-custom", kafkaBinding.GetExtensions()))
	assertHashNotZero(t, kafkaBinding.Hash())

	httpServer := low.FindItemInOrderedMap("http-server", doc.Servers.Value)
	require.NotNil(t, httpServer)
	require.NotNil(t, httpServer.Value.Bindings.Value.HTTP.Value)
	assertHashNotZero(t, httpServer.Value.Bindings.Value.HTTP.Value.Hash())

	require.NotNil(t, doc.Channels.Value)
	kafkaChannel := low.FindItemInOrderedMap("kafka-channel", doc.Channels.Value)
	require.NotNil(t, kafkaChannel)
	assertHashNotZero(t, kafkaChannel.Value.Hash())

	require.NotNil(t, kafkaChannel.Value.Bindings.Value.Kafka.Value)
	assertHashNotZero(t, kafkaChannel.Value.Bindings.Value.Kafka.Value.Hash())
	require.NotNil(t, kafkaChannel.Value.Bindings.Value.Kafka.Value.TopicConfiguration.Value)
	assertHashNotZero(t, kafkaChannel.Value.Bindings.Value.Kafka.Value.TopicConfiguration.Value.Hash())

	amqpChannel := low.FindItemInOrderedMap("amqp-channel", doc.Channels.Value)
	require.NotNil(t, amqpChannel)
	require.NotNil(t, amqpChannel.Value.Bindings.Value.AMQP.Value)
	assertHashNotZero(t, amqpChannel.Value.Bindings.Value.AMQP.Value.Hash())
	require.NotNil(t, amqpChannel.Value.Bindings.Value.AMQP.Value.Exchange.Value)
	assertHashNotZero(t, amqpChannel.Value.Bindings.Value.AMQP.Value.Exchange.Value.Hash())
	require.NotNil(t, amqpChannel.Value.Bindings.Value.AMQP.Value.Queue.Value)
	assertHashNotZero(t, amqpChannel.Value.Bindings.Value.AMQP.Value.Queue.Value.Hash())

	wsChannel := low.FindItemInOrderedMap("ws-channel", doc.Channels.Value)
	require.NotNil(t, wsChannel)
	require.NotNil(t, wsChannel.Value.Bindings.Value.WebSocket.Value)
	assertHashNotZero(t, wsChannel.Value.Bindings.Value.WebSocket.Value.Hash())

	httpChannel := low.FindItemInOrderedMap("http-channel", doc.Channels.Value)
	require.NotNil(t, httpChannel)
	require.NotNil(t, httpChannel.Value.Bindings.Value.HTTP.Value)
	assertHashNotZero(t, httpChannel.Value.Bindings.Value.HTTP.Value.Hash())

	require.NotNil(t, doc.Operations.Value)
	kafkaOp := low.FindItemInOrderedMap("kafka-operation", doc.Operations.Value)
	require.NotNil(t, kafkaOp)
	assertHashNotZero(t, kafkaOp.Value.Hash())
	require.NotNil(t, kafkaOp.Value.Bindings.Value.Kafka.Value)
	assertHashNotZero(t, kafkaOp.Value.Bindings.Value.Kafka.Value.Hash())

	amqpOp := low.FindItemInOrderedMap("amqp-operation", doc.Operations.Value)
	require.NotNil(t, amqpOp)
	require.NotNil(t, amqpOp.Value.Bindings.Value.AMQP.Value)
	assertHashNotZero(t, amqpOp.Value.Bindings.Value.AMQP.Value.Hash())

	mqttOp := low.FindItemInOrderedMap("mqtt-operation", doc.Operations.Value)
	require.NotNil(t, mqttOp)
	require.NotNil(t, mqttOp.Value.Bindings.Value.MQTT.Value)
	assertHashNotZero(t, mqttOp.Value.Bindings.Value.MQTT.Value.Hash())

	httpOp := low.FindItemInOrderedMap("http-operation", doc.Operations.Value)
	require.NotNil(t, httpOp)
	require.NotNil(t, httpOp.Value.Bindings.Value.HTTP.Value)
	assertHashNotZero(t, httpOp.Value.Bindings.Value.HTTP.Value.Hash())

	require.NotNil(t, doc.Components.Value)
	components := doc.Components.Value
	assertHashNotZero(t, components.Hash())

	kafkaMessage := low.FindItemInOrderedMap("kafkaMessage", components.Messages.Value)
	require.NotNil(t, kafkaMessage)
	assertHashNotZero(t, kafkaMessage.Value.Hash())
	require.NotNil(t, kafkaMessage.Value.Bindings.Value.Kafka.Value)
	assertHashNotZero(t, kafkaMessage.Value.Bindings.Value.Kafka.Value.Hash())

	amqpMessage := low.FindItemInOrderedMap("amqpMessage", components.Messages.Value)
	require.NotNil(t, amqpMessage)
	require.NotNil(t, amqpMessage.Value.Bindings.Value.AMQP.Value)
	assertHashNotZero(t, amqpMessage.Value.Bindings.Value.AMQP.Value.Hash())

	mqttMessage := low.FindItemInOrderedMap("mqttMessage", components.Messages.Value)
	require.NotNil(t, mqttMessage)
	require.NotNil(t, mqttMessage.Value.Bindings.Value.MQTT.Value)
	assertHashNotZero(t, mqttMessage.Value.Bindings.Value.MQTT.Value.Hash())

	httpMessage := low.FindItemInOrderedMap("httpMessage", components.Messages.Value)
	require.NotNil(t, httpMessage)
	require.NotNil(t, httpMessage.Value.Bindings.Value.HTTP.Value)
	assertHashNotZero(t, httpMessage.Value.Bindings.Value.HTTP.Value.Hash())
}

func TestCreateDocument_MultiProtocolComponents(t *testing.T) {
	root := loadFixtureNode(t, "multi-protocol.yaml")
	doc := createDocument(t, root, NewDocumentConfiguration())

	require.NotNil(t, doc.Components.Value)
	comp := doc.Components.Value
	assertHashNotZero(t, comp.Hash())

	assert.NotNil(t, comp.FindSchema("User"))
	assert.NotNil(t, comp.FindMessage("userCreated"))
	assert.NotNil(t, comp.FindSecurityScheme("oauth2"))

	oauth := comp.FindSecurityScheme("oauth2")
	require.NotNil(t, oauth)
	assertHashNotZero(t, oauth.Value.Hash())
	require.NotNil(t, oauth.Value.Flows.Value)
	assertHashNotZero(t, oauth.Value.Flows.Value.Hash())
	require.NotNil(t, oauth.Value.Flows.Value.ClientCredentials.Value)
	assertHashNotZero(t, oauth.Value.Flows.Value.ClientCredentials.Value.Hash())

	apiKey := comp.FindSecurityScheme("apiKey")
	require.NotNil(t, apiKey)
	assertHashNotZero(t, apiKey.Value.Hash())

	serverVar := low.FindItemInOrderedMap("environment", comp.ServerVariables.Value)
	require.NotNil(t, serverVar)
	assertHashNotZero(t, serverVar.Value.Hash())

	param := low.FindItemInOrderedMap("userId", comp.Parameters.Value)
	require.NotNil(t, param)
	assertHashNotZero(t, param.Value.Hash())

	cid := low.FindItemInOrderedMap("orderCorrelation", comp.CorrelationIDs.Value)
	require.NotNil(t, cid)
	assertHashNotZero(t, cid.Value.Hash())

	opTrait := low.FindItemInOrderedMap("kafkaCommon", comp.OperationTraits.Value)
	require.NotNil(t, opTrait)
	assertHashNotZero(t, opTrait.Value.Hash())

	msgTrait := low.FindItemInOrderedMap("commonHeaders", comp.MessageTraits.Value)
	require.NotNil(t, msgTrait)
	assertHashNotZero(t, msgTrait.Value.Hash())

	serverBindings := low.FindItemInOrderedMap("kafkaBinding", comp.ServerBindings.Value)
	require.NotNil(t, serverBindings)
	assertHashNotZero(t, serverBindings.Value.Hash())
	require.NotNil(t, serverBindings.Value.Kafka.Value)
	assertHashNotZero(t, serverBindings.Value.Kafka.Value.Hash())

	channelBindings := low.FindItemInOrderedMap("kafkaChannelBinding", comp.ChannelBindings.Value)
	require.NotNil(t, channelBindings)
	assertHashNotZero(t, channelBindings.Value.Hash())
	require.NotNil(t, channelBindings.Value.Kafka.Value)
	assertHashNotZero(t, channelBindings.Value.Kafka.Value.Hash())

	operationBindings := low.FindItemInOrderedMap("kafkaOperationBinding", comp.OperationBindings.Value)
	require.NotNil(t, operationBindings)
	assertHashNotZero(t, operationBindings.Value.Hash())
	require.NotNil(t, operationBindings.Value.Kafka.Value)
	assertHashNotZero(t, operationBindings.Value.Kafka.Value.Hash())

	messageBindings := low.FindItemInOrderedMap("kafkaMessageBinding", comp.MessageBindings.Value)
	require.NotNil(t, messageBindings)
	assertHashNotZero(t, messageBindings.Value.Hash())
	require.NotNil(t, messageBindings.Value.Kafka.Value)
	assertHashNotZero(t, messageBindings.Value.Kafka.Value.Hash())

	tag := low.FindItemInOrderedMap("user-events", comp.Tags.Value)
	require.NotNil(t, tag)
	assertHashNotZero(t, tag.Value.Hash())

	extDoc := low.FindItemInOrderedMap("apiDocs", comp.ExternalDocs.Value)
	require.NotNil(t, extDoc)
	assertHashNotZero(t, extDoc.Value.Hash())
}

func TestCreateDocument_FileRefResolution(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: File Ref API
  version: "1.0.0"
channels:
  events:
    address: events/stream
    messages:
      externalMessage:
        $ref: 'shared-message.yaml'
`)

	root := loadSpecNode(t, spec)
	basePath := filepath.Join("..", "..", "..", "test_fixtures")
	config := &DocumentConfiguration{
		BasePath:            basePath,
		AllowFileReferences: true,
	}
	doc := createDocument(t, root, config)

	channel := low.FindItemInOrderedMap("events", doc.Channels.Value)
	require.NotNil(t, channel)
	require.NotNil(t, channel.Value.Messages.Value)
	assert.Greater(t, channel.Value.Messages.Value.Len(), 0)
}

func TestOperationReplyAndParameters(t *testing.T) {
	spec := []byte(`asyncapi: 3.0.0
info:
  title: Reply API
  version: 1.0.0
channels:
  replyChannel:
    address: reply/channel
    parameters:
      status:
        enum:
          - ok
          - error
        default: ok
        location: $message.payload#/status
        description: Response status
        examples:
          - ok
        x-param-ext: value
    messages:
      replyMessage:
        $ref: '#/components/messages/replyMessage'
operations:
  handleReply:
    action: send
    channel:
      $ref: '#/channels/replyChannel'
    reply:
      x-reply-ext: value
      address:
        description: reply address
        location: $message.header#/replyTo
        x-reply-address-ext: value
      channel:
        $ref: '#/channels/replyChannel'
      messages:
        - $ref: '#/channels/replyChannel/messages/replyMessage'
components:
  messages:
    replyMessage:
      name: replyMessage
      contentType: application/json
      payload:
        type: object
`)

	doc := createDocument(t, loadSpecNode(t, spec), NewDocumentConfiguration())

	channel := low.FindItemInOrderedMap("replyChannel", doc.Channels.Value)
	require.NotNil(t, channel)
	param := low.FindItemInOrderedMap("status", channel.Value.Parameters.Value)
	require.NotNil(t, param)
	assertHashNotZero(t, param.Value.Hash())
	assert.NotNil(t, low.FindItemInOrderedMap("x-param-ext", param.Value.Extensions))

	op := low.FindItemInOrderedMap("handleReply", doc.Operations.Value)
	require.NotNil(t, op)
	require.NotNil(t, op.Value.Reply.Value)
	assertHashNotZero(t, op.Value.Reply.Value.Hash())
	require.NotNil(t, op.Value.Reply.Value.Address.Value)
	assertHashNotZero(t, op.Value.Reply.Value.Address.Value.Hash())
	assert.NotNil(t, low.FindItemInOrderedMap("x-reply-address-ext", op.Value.Reply.Value.Address.Value.Extensions))
}
