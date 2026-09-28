// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
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

func createLowDoc(t *testing.T, root *yaml.Node) *lowasync.AsyncAPI {
	t.Helper()
	doc, err := lowasync.CreateDocumentWithConfig(root, lowasync.NewDocumentConfiguration())
	require.NoError(t, err)
	require.NotNil(t, doc)
	return doc
}

func TestNewAsyncAPI_ComprehensiveBindings(t *testing.T) {
	lowDoc := createLowDoc(t, loadFixtureNode(t, "comprehensive-bindings.yaml"))
	doc := NewAsyncAPI(lowDoc)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.AsyncAPI)
	require.NotNil(t, doc.Info)
	assert.Equal(t, "Comprehensive Bindings Test", doc.Info.Title)
	assert.NotNil(t, doc.Info.GoLow())
	assert.NotNil(t, doc.Info.GoLowUntyped())

	rendered, err := doc.Render()
	require.NoError(t, err)
	assert.Contains(t, string(rendered), "asyncapi")

	server, ok := doc.Servers.Get("kafka-server")
	require.True(t, ok)
	require.NotNil(t, server.Bindings)
	require.NotNil(t, server.Bindings.Kafka)
	assert.Equal(t, "https://schema-registry.example.com", server.Bindings.Kafka.SchemaRegistryURL)
	assert.NotNil(t, server.Bindings.Kafka.GoLow())
	assert.NotNil(t, server.Bindings.Kafka.GoLowUntyped())

	channel, ok := doc.Channels.Get("amqp-channel")
	require.True(t, ok)
	require.NotNil(t, channel.Bindings)
	require.NotNil(t, channel.Bindings.AMQP)
	assert.Equal(t, "routingKey", channel.Bindings.AMQP.Is)
	require.NotNil(t, channel.Bindings.AMQP.Exchange)
	require.NotNil(t, channel.Bindings.AMQP.Queue)
	assert.NotNil(t, channel.Bindings.AMQP.GoLow())

	op, ok := doc.Operations.Get("http-operation")
	require.True(t, ok)
	require.NotNil(t, op.Bindings)
	require.NotNil(t, op.Bindings.HTTP)
	assert.Equal(t, "POST", op.Bindings.HTTP.Method)
	assert.NotNil(t, op.Bindings.HTTP.GoLow())
	assert.NotNil(t, op.Bindings.HTTP.GoLowUntyped())

	message, ok := doc.Components.Messages.Get("kafkaMessage")
	require.True(t, ok)
	require.NotNil(t, message.Bindings)
	require.NotNil(t, message.Bindings.Kafka)
	assert.Equal(t, "TopicIdStrategy", message.Bindings.Kafka.SchemaLookupStrategy)
	assert.NotNil(t, message.Bindings.Kafka.GoLow())
	assert.NotNil(t, message.Bindings.Kafka.GoLowUntyped())

	infoRendered, err := doc.Info.Render()
	require.NoError(t, err)
	assert.Contains(t, string(infoRendered), "Comprehensive Bindings Test")

	channelYAML, err := channel.Render()
	require.NoError(t, err)
	assert.Contains(t, string(channelYAML), "orders")
}

func TestNewAsyncAPI_MultiProtocolComponents(t *testing.T) {
	lowDoc := createLowDoc(t, loadFixtureNode(t, "multi-protocol.yaml"))
	doc := NewAsyncAPI(lowDoc)
	require.NotNil(t, doc)

	require.NotNil(t, doc.Components)
	assert.NotNil(t, doc.Components.GoLow())
	assert.NotNil(t, doc.Components.GoLowUntyped())

	oauth, ok := doc.Components.SecuritySchemes.Get("oauth2")
	require.True(t, ok)
	require.NotNil(t, oauth.Flows)
	require.NotNil(t, oauth.Flows.ClientCredentials)
	assert.NotNil(t, oauth.GoLow())
	assert.NotNil(t, oauth.GoLowUntyped())

	apiKey, ok := doc.Components.SecuritySchemes.Get("apiKey")
	require.True(t, ok)
	assert.Equal(t, "header", apiKey.In)

	param, ok := doc.Components.Parameters.Get("userId")
	require.True(t, ok)
	assert.Equal(t, "The unique user identifier", param.Description)
	assert.NotNil(t, param.GoLow())

	serverVar, ok := doc.Components.ServerVariables.Get("environment")
	require.True(t, ok)
	assert.Equal(t, "production", serverVar.Default)
	assert.NotNil(t, serverVar.GoLow())

	tag, ok := doc.Components.Tags.Get("user-events")
	require.True(t, ok)
	assert.Equal(t, "user-events", tag.Name)
	assert.NotNil(t, tag.GoLow())

	extDoc, ok := doc.Components.ExternalDocs.Get("apiDocs")
	require.True(t, ok)
	assert.Equal(t, "https://docs.example.com/api", extDoc.URL)
	assert.NotNil(t, extDoc.GoLow())

	secRendered, err := oauth.Render()
	require.NoError(t, err)
	assert.Contains(t, string(secRendered), "oauth2")
}

func TestNewAsyncAPI_OperationReply(t *testing.T) {
	spec := []byte(`asyncapi: 3.0.0
info:
  title: Reply API
  version: 1.0.0
channels:
  replyChannel:
    address: reply/channel
    messages:
      replyMessage:
        $ref: '#/components/messages/replyMessage'
operations:
  handleReply:
    action: send
    channel:
      $ref: '#/channels/replyChannel'
    reply:
      address:
        description: reply address
        location: $message.header#/replyTo
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

	lowDoc := createLowDoc(t, loadSpecNode(t, spec))
	doc := NewAsyncAPI(lowDoc)
	require.NotNil(t, doc)

	op, ok := doc.Operations.Get("handleReply")
	require.True(t, ok)
	require.NotNil(t, op.Reply)
	require.NotNil(t, op.Reply.Address)
	assert.NotNil(t, op.Reply.GoLow())
	assert.NotNil(t, op.Reply.GoLowUntyped())
	assert.NotNil(t, op.Reply.Address.GoLow())
	assert.NotNil(t, op.Reply.Address.GoLowUntyped())

	opYAML, err := op.Render()
	require.NoError(t, err)
	assert.Contains(t, string(opYAML), "action: send")
	assert.Contains(t, string(opYAML), "reply:")
}
