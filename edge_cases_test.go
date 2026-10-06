// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"os"
	"testing"

	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestEdgeCase_MultiProtocolDocument(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/multi-protocol.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)

	// Test ID field
	assert.Equal(t, "urn:com:example:multi-protocol-api", model.ID)

	// Test multiple servers
	require.NotNil(t, model.Servers)
	assert.Equal(t, 3, model.Servers.Len())

	kafkaServer, ok := model.Servers.Get("production-kafka")
	require.True(t, ok)
	assert.Equal(t, "kafka-secure", kafkaServer.Protocol)
	require.NotNil(t, kafkaServer.Bindings)
	require.NotNil(t, kafkaServer.Bindings.Kafka)

	mqttServer, ok := model.Servers.Get("staging-mqtt")
	require.True(t, ok)
	assert.Equal(t, "secure-mqtt", mqttServer.Protocol)
	require.NotNil(t, mqttServer.Bindings)
	require.NotNil(t, mqttServer.Bindings.MQTT)

	httpServer, ok := model.Servers.Get("development-http")
	require.True(t, ok)
	assert.Equal(t, "http", httpServer.Protocol)

	// Test channels with different bindings
	require.NotNil(t, model.Channels)
	assert.Equal(t, 3, model.Channels.Len())

	// Kafka channel
	userEvents, ok := model.Channels.Get("userEvents")
	require.True(t, ok)
	require.NotNil(t, userEvents.Bindings)
	require.NotNil(t, userEvents.Bindings.Kafka)
	assert.Equal(t, "user-events", userEvents.Bindings.Kafka.Topic)
	assert.Equal(t, 12, userEvents.Bindings.Kafka.Partitions)

	// AMQP channel
	orderQueue, ok := model.Channels.Get("orderQueue")
	require.True(t, ok)
	require.NotNil(t, orderQueue.Bindings)
	require.NotNil(t, orderQueue.Bindings.AMQP)
	assert.Equal(t, "queue", orderQueue.Bindings.AMQP.Is)

	// WebSocket channel
	notifications, ok := model.Channels.Get("notifications")
	require.True(t, ok)
	require.NotNil(t, notifications.Bindings)
	require.NotNil(t, notifications.Bindings.WebSocket)
	assert.Equal(t, "GET", notifications.Bindings.WebSocket.Method)

	// Test operations
	require.NotNil(t, model.Operations)
	assert.Equal(t, 4, model.Operations.Len())

	// Test security schemes
	require.NotNil(t, model.Components)
	require.NotNil(t, model.Components.SecuritySchemes)
	assert.Equal(t, 2, model.Components.SecuritySchemes.Len())

	oauth2, ok := model.Components.SecuritySchemes.Get("oauth2")
	require.True(t, ok)
	assert.Equal(t, "oauth2", oauth2.Type)

	apiKey, ok := model.Components.SecuritySchemes.Get("apiKey")
	require.True(t, ok)
	assert.Equal(t, "apiKey", apiKey.Type)
	assert.Equal(t, "header", apiKey.In)
}

func TestEdgeCase_EmptyChannels(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Empty Channels API
  version: 1.0.0
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	assert.Nil(t, model.Channels)
	assert.Nil(t, model.Operations)
}

func TestEdgeCase_ExtensionsPreserved(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Extensions API
  version: 1.0.0
  x-custom-field: custom value
  x-another-extension:
    nested: data
    array:
      - item1
      - item2
x-root-extension: root value
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Info)
	require.NotNil(t, model.Info.Extensions)
	assert.Greater(t, model.Info.Extensions.Len(), 0)

	require.NotNil(t, model.Extensions)
	assert.Greater(t, model.Extensions.Len(), 0)
}

func TestEdgeCase_ComplexSchemaReferences(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Schema Refs API
  version: 1.0.0
components:
  schemas:
    BaseEntity:
      type: object
      properties:
        id:
          type: string
    User:
      allOf:
        - $ref: '#/components/schemas/BaseEntity'
        - type: object
          properties:
            name:
              type: string
    Address:
      type: object
      properties:
        street:
          type: string
        city:
          type: string
    UserWithAddress:
      type: object
      properties:
        user:
          $ref: '#/components/schemas/User'
        addresses:
          type: array
          items:
            $ref: '#/components/schemas/Address'
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Components)
	require.NotNil(t, model.Components.Schemas)
	assert.Equal(t, 4, model.Components.Schemas.Len())
}

func TestEdgeCase_OperationWithMultipleTraits(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Traits API
  version: 1.0.0
channels:
  test:
    address: test/channel
    messages:
      testMessage:
        payload:
          type: string
operations:
  testOp:
    action: send
    channel:
      $ref: '#/channels/test'
    traits:
      - $ref: '#/components/operationTraits/trait1'
      - $ref: '#/components/operationTraits/trait2'
    messages:
      - $ref: '#/channels/test/messages/testMessage'
components:
  operationTraits:
    trait1:
      description: First trait
    trait2:
      description: Second trait
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Operations)

	testOp, ok := model.Operations.Get("testOp")
	require.True(t, ok)
	require.NotNil(t, testOp.Traits)
	assert.Len(t, testOp.Traits, 2)
}

func TestEdgeCase_MessageWithTraitsAndExamples(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Message Traits API
  version: 1.0.0
components:
  messages:
    userEvent:
      name: userEvent
      title: User Event
      summary: A user-related event
      contentType: application/json
      traits:
        - $ref: '#/components/messageTraits/commonHeaders'
      payload:
        type: object
        properties:
          userId:
            type: string
      examples:
        - name: example1
          summary: Example user event
          payload:
            userId: "12345"
  messageTraits:
    commonHeaders:
      headers:
        type: object
        properties:
          correlationId:
            type: string
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Components)
	require.NotNil(t, model.Components.Messages)

	userEvent, ok := model.Components.Messages.Get("userEvent")
	require.True(t, ok)
	assert.Equal(t, "userEvent", userEvent.Name)
	assert.Equal(t, "User Event", userEvent.Title)
	require.NotNil(t, userEvent.Traits)
	assert.Len(t, userEvent.Traits, 1)
	require.NotNil(t, userEvent.Examples)
	assert.Len(t, userEvent.Examples, 1)
}

func TestEdgeCase_ServerVariables(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Server Variables API
  version: 1.0.0
servers:
  production:
    host: "{environment}.example.com:{port}"
    protocol: kafka
    variables:
      environment:
        default: prod
        enum:
          - prod
          - staging
        description: Environment name
      port:
        default: "9092"
        description: Port number
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Servers)

	server, ok := model.Servers.Get("production")
	require.True(t, ok)
	require.NotNil(t, server.Variables)
	assert.Equal(t, 2, server.Variables.Len())

	envVar, ok := server.Variables.Get("environment")
	require.True(t, ok)
	assert.Equal(t, "prod", envVar.Default)
	assert.Len(t, envVar.Enum, 2)
}

func TestEdgeCase_ChannelParametersWithSchema(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Parameters API
  version: 1.0.0
channels:
  userChannel:
    address: users/{userId}/events/{eventType}
    parameters:
      userId:
        description: User identifier
        examples:
          - "user-123"
      eventType:
        description: Type of event
        enum:
          - created
          - updated
          - deleted
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Channels)

	channel, ok := model.Channels.Get("userChannel")
	require.True(t, ok)
	require.NotNil(t, channel.Parameters)
	assert.Equal(t, 2, channel.Parameters.Len())
}

func TestEdgeCase_ReplyAndReplyAddress(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Reply API
  version: 1.0.0
channels:
  request:
    address: request/channel
    messages:
      request:
        payload:
          type: string
  response:
    address: response/channel
    messages:
      response:
        payload:
          type: string
operations:
  sendRequest:
    action: send
    channel:
      $ref: '#/channels/request'
    reply:
      channel:
        $ref: '#/channels/response'
      address:
        location: $message.header#/replyTo
      messages:
        - $ref: '#/channels/response/messages/response'
    messages:
      - $ref: '#/channels/request/messages/request'
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Operations)

	op, ok := model.Operations.Get("sendRequest")
	require.True(t, ok)
	require.NotNil(t, op.Reply)
	require.NotNil(t, op.Reply.Address)
	assert.Equal(t, "$message.header#/replyTo", op.Reply.Address.Location)
}

func TestEdgeCase_CorrelationId(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Correlation API
  version: 1.0.0
components:
  correlationIds:
    defaultCorrelation:
      description: Default correlation ID
      location: $message.header#/correlationId
  messages:
    correlatedMessage:
      payload:
        type: string
      correlationId:
        $ref: '#/components/correlationIds/defaultCorrelation'
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Components)
	require.NotNil(t, model.Components.CorrelationIDs)
	assert.Equal(t, 1, model.Components.CorrelationIDs.Len())
}

func TestEdgeCase_SecurityAtMultipleLevels(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Security API
  version: 1.0.0
servers:
  production:
    host: api.example.com
    protocol: kafka
    security:
      - $ref: '#/components/securitySchemes/oauth2'
channels:
  secure:
    address: secure/channel
    messages:
      secureMessage:
        payload:
          type: string
operations:
  secureOp:
    action: send
    channel:
      $ref: '#/channels/secure'
    security:
      - $ref: '#/components/securitySchemes/apiKey'
    messages:
      - $ref: '#/channels/secure/messages/secureMessage'
components:
  securitySchemes:
    oauth2:
      type: oauth2
      flows:
        clientCredentials:
          tokenUrl: https://auth.example.com/token
          scopes:
            read: Read access
    apiKey:
      type: httpApiKey
      name: X-API-Key
      in: header
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)

	// Server security
	server, ok := model.Servers.Get("production")
	require.True(t, ok)
	require.NotNil(t, server.Security)
	assert.Len(t, server.Security, 1)

	// Operation security
	op, ok := model.Operations.Get("secureOp")
	require.True(t, ok)
	require.NotNil(t, op.Security)
	assert.Len(t, op.Security, 1)
}

func TestEdgeCase_AllBindingTypes(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/multi-protocol.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Server bindings
	kafkaServer, _ := model.Servers.Get("production-kafka")
	require.NotNil(t, kafkaServer.Bindings.Kafka)
	assert.Equal(t, "https://schema-registry.production.example.com", kafkaServer.Bindings.Kafka.SchemaRegistryURL)

	mqttServer, _ := model.Servers.Get("staging-mqtt")
	require.NotNil(t, mqttServer.Bindings.MQTT)
	assert.Equal(t, "multi-protocol-client", mqttServer.Bindings.MQTT.ClientID)

	// Channel bindings
	userEvents, _ := model.Channels.Get("userEvents")
	require.NotNil(t, userEvents.Bindings.Kafka)
	assert.Equal(t, 12, userEvents.Bindings.Kafka.Partitions)

	orderQueue, _ := model.Channels.Get("orderQueue")
	require.NotNil(t, orderQueue.Bindings.AMQP)
	assert.Equal(t, "queue", orderQueue.Bindings.AMQP.Is)

	notifications, _ := model.Channels.Get("notifications")
	require.NotNil(t, notifications.Bindings.WebSocket)
	assert.Equal(t, "GET", notifications.Bindings.WebSocket.Method)

	// Operation bindings
	receiveOp, _ := model.Operations.Get("receiveUserUpdated")
	require.NotNil(t, receiveOp.Bindings)
	require.NotNil(t, receiveOp.Bindings.Kafka)

	processOp, _ := model.Operations.Get("processOrder")
	require.NotNil(t, processOp.Bindings)
	require.NotNil(t, processOp.Bindings.AMQP)
	assert.Equal(t, 60000, processOp.Bindings.AMQP.Expiration)

	sendOp, _ := model.Operations.Get("sendNotification")
	require.NotNil(t, sendOp.Bindings)
	require.NotNil(t, sendOp.Bindings.HTTP)
	assert.Equal(t, "POST", sendOp.Bindings.HTTP.Method)

	// Message bindings
	userCreatedMsg, _ := model.Components.Messages.Get("userCreated")
	require.NotNil(t, userCreatedMsg.Bindings)
	require.NotNil(t, userCreatedMsg.Bindings.Kafka)

	orderPlacedMsg, _ := model.Components.Messages.Get("orderPlaced")
	require.NotNil(t, orderPlacedMsg.Bindings)
	require.NotNil(t, orderPlacedMsg.Bindings.AMQP)
	assert.Equal(t, "utf-8", orderPlacedMsg.Bindings.AMQP.ContentEncoding)

	notificationMsg, _ := model.Components.Messages.Get("notification")
	require.NotNil(t, notificationMsg.Bindings)
	require.NotNil(t, notificationMsg.Bindings.HTTP)
}

func TestEdgeCase_ComponentBindingsReusable(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/multi-protocol.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)
	require.NotNil(t, model.Components)

	// Test reusable bindings in components
	require.NotNil(t, model.Components.ServerBindings)
	assert.Equal(t, 1, model.Components.ServerBindings.Len())

	require.NotNil(t, model.Components.ChannelBindings)
	assert.Equal(t, 1, model.Components.ChannelBindings.Len())

	require.NotNil(t, model.Components.OperationBindings)
	assert.Equal(t, 1, model.Components.OperationBindings.Len())

	require.NotNil(t, model.Components.MessageBindings)
	assert.Equal(t, 1, model.Components.MessageBindings.Len())
}

func TestEdgeCase_RenderPreservesStructure(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/multi-protocol.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	rendered, err := doc.Render()
	require.NoError(t, err)
	require.NotEmpty(t, rendered)

	// Parse the rendered output
	doc2, err := NewDocument(rendered)
	require.NoError(t, err)

	// Verify structure is preserved
	model1 := doc.Model()
	model2 := doc2.Model()

	assert.Equal(t, model1.Info.Title, model2.Info.Title)
	assert.Equal(t, model1.Info.Version, model2.Info.Version)
	assert.Equal(t, model1.Servers.Len(), model2.Servers.Len())
	assert.Equal(t, model1.Channels.Len(), model2.Channels.Len())
	assert.Equal(t, model1.Operations.Len(), model2.Operations.Len())
}

func TestEdgeCase_GoLowPreservesLineNumbers(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	lowModel := doc.GoLow()
	require.NotNil(t, lowModel)

	// Verify line numbers are preserved
	assert.Greater(t, lowModel.Info.KeyNode.Line, 0)
	assert.Greater(t, lowModel.Info.ValueNode.Line, 0)

	// Info title should be on line 3
	assert.Equal(t, 2, lowModel.Info.KeyNode.Line)
}

func TestEdgeCase_ExternalDocs(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: External Docs API
  version: 1.0.0
  externalDocs:
    description: Full documentation
    url: https://docs.example.com
components:
  externalDocs:
    wiki:
      description: Wiki documentation
      url: https://wiki.example.com
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Info level external docs
	require.NotNil(t, model.Info.ExternalDocs)
	assert.Equal(t, "Full documentation", model.Info.ExternalDocs.Description)
	assert.Equal(t, "https://docs.example.com", model.Info.ExternalDocs.URL)

	// Component external docs
	require.NotNil(t, model.Components)
	require.NotNil(t, model.Components.ExternalDocs)
	assert.Equal(t, 1, model.Components.ExternalDocs.Len())
}

func TestEdgeCase_TagsAtMultipleLevels(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Tags API
  version: 1.0.0
  tags:
    - name: info-tag
      description: Tag at info level
servers:
  production:
    host: api.example.com
    protocol: kafka
    tags:
      - name: server-tag
        description: Tag at server level
channels:
  events:
    address: events
    tags:
      - name: channel-tag
        description: Tag at channel level
    messages:
      event:
        payload:
          type: string
operations:
  sendEvent:
    action: send
    channel:
      $ref: '#/channels/events'
    tags:
      - name: operation-tag
        description: Tag at operation level
    messages:
      - $ref: '#/channels/events/messages/event'
components:
  tags:
    reusable-tag:
      name: reusable-tag
      description: Reusable tag in components
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)

	// Info tags
	require.NotNil(t, model.Info.Tags)
	assert.Len(t, model.Info.Tags, 1)
	assert.Equal(t, "info-tag", model.Info.Tags[0].Name)

	// Server tags
	server, _ := model.Servers.Get("production")
	require.NotNil(t, server.Tags)
	assert.Len(t, server.Tags, 1)

	// Channel tags
	channel, _ := model.Channels.Get("events")
	require.NotNil(t, channel.Tags)
	assert.Len(t, channel.Tags, 1)

	// Operation tags
	op, _ := model.Operations.Get("sendEvent")
	require.NotNil(t, op.Tags)
	assert.Len(t, op.Tags, 1)

	// Component tags
	require.NotNil(t, model.Components.Tags)
	assert.Equal(t, 1, model.Components.Tags.Len())
}
