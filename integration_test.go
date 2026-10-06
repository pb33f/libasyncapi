// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"os"
	"testing"

	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestIntegration_StreetlightsKafka(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.GetVersion())
	assert.False(t, doc.IsPartial())

	model := doc.Model()
	require.NotNil(t, model)

	// Test Info
	require.NotNil(t, model.Info)
	assert.Equal(t, "Streetlights Kafka API", model.Info.Title)
	assert.Equal(t, "1.0.0", model.Info.Version)
	assert.Contains(t, model.Info.Description, "Smartylighting Streetlights API")
	assert.Equal(t, "https://asyncapi.org/terms/", model.Info.TermsOfService)

	// Test Info Contact
	require.NotNil(t, model.Info.Contact)
	assert.Equal(t, "API Support", model.Info.Contact.Name)
	assert.Equal(t, "https://www.asyncapi.org/support", model.Info.Contact.URL)
	assert.Equal(t, "support@asyncapi.org", model.Info.Contact.Email)

	// Test Info License
	require.NotNil(t, model.Info.License)
	assert.Equal(t, "Apache 2.0", model.Info.License.Name)
	assert.Equal(t, "https://www.apache.org/licenses/LICENSE-2.0.html", model.Info.License.URL)

	// Test Info Tags
	require.NotNil(t, model.Info.Tags)
	assert.Len(t, model.Info.Tags, 2)
	assert.Equal(t, "streetlights", model.Info.Tags[0].Name)
	assert.Equal(t, "kafka", model.Info.Tags[1].Name)

	// Test DefaultContentType
	assert.Equal(t, "application/json", model.DefaultContentType)

	// Test Servers
	require.NotNil(t, model.Servers)
	assert.Equal(t, 1, model.Servers.Len())

	server, ok := model.Servers.Get("scram-connections")
	require.True(t, ok)
	assert.Equal(t, "test.mykafkacluster.org:18092", server.Host)
	assert.Equal(t, "kafka-secure", server.Protocol)
	assert.Contains(t, server.Description, "Test broker secured with SASL/SCRAM")

	// Test Server Bindings
	require.NotNil(t, server.Bindings)
	require.NotNil(t, server.Bindings.Kafka)
	assert.Equal(t, "https://my-schema-registry.com", server.Bindings.Kafka.SchemaRegistryURL)
	assert.Equal(t, "confluent", server.Bindings.Kafka.SchemaRegistryVendor)
	assert.Equal(t, "0.4.0", server.Bindings.Kafka.BindingVersion)

	// Test Channels
	require.NotNil(t, model.Channels)
	assert.Equal(t, 3, model.Channels.Len())

	lightingMeasured, ok := model.Channels.Get("lightingMeasured")
	require.True(t, ok)
	require.NotNil(t, lightingMeasured.Address)
	assert.Contains(t, *lightingMeasured.Address, "lighting.measured")
	assert.Contains(t, lightingMeasured.Description, "measured values")

	// Test Channel Bindings
	require.NotNil(t, lightingMeasured.Bindings)
	require.NotNil(t, lightingMeasured.Bindings.Kafka)
	assert.Equal(t, "streetlights-lighting", lightingMeasured.Bindings.Kafka.Topic)
	assert.Equal(t, 3, lightingMeasured.Bindings.Kafka.Partitions)
	assert.Equal(t, 2, lightingMeasured.Bindings.Kafka.Replicas)

	// Test Channel Parameters
	require.NotNil(t, lightingMeasured.Parameters)
	assert.Equal(t, 1, lightingMeasured.Parameters.Len())

	// Test Operations
	require.NotNil(t, model.Operations)
	assert.Equal(t, 3, model.Operations.Len())

	receiveOp, ok := model.Operations.Get("receiveLightMeasurement")
	require.True(t, ok)
	assert.Equal(t, "receive", receiveOp.Action)
	assert.Contains(t, receiveOp.Summary, "environmental lighting conditions")
	// Test Operation Channel reference
	require.NotNil(t, receiveOp.Channel)
	assert.Contains(t, receiveOp.Channel.GetReference(), "#/channels/lightingMeasured")
	// Test Operation Messages references
	require.NotNil(t, receiveOp.Messages)
	assert.Len(t, receiveOp.Messages, 1)

	turnOnOp, ok := model.Operations.Get("turnOn")
	require.True(t, ok)
	assert.Equal(t, "send", turnOnOp.Action)
	require.NotNil(t, turnOnOp.Channel)
	require.NotNil(t, turnOnOp.Messages)

	// Test Components
	require.NotNil(t, model.Components)

	// Test Components Messages
	require.NotNil(t, model.Components.Messages)
	assert.Equal(t, 2, model.Components.Messages.Len())

	lightMeasuredMsg, ok := model.Components.Messages.Get("lightMeasured")
	require.True(t, ok)
	assert.Equal(t, "lightMeasured", lightMeasuredMsg.Name)
	assert.Equal(t, "Light measured", lightMeasuredMsg.Title)
	assert.Equal(t, "application/json", lightMeasuredMsg.ContentType)

	// Test Components Schemas
	require.NotNil(t, model.Components.Schemas)
	assert.Equal(t, 3, model.Components.Schemas.Len())

	// Test Components SecuritySchemes
	require.NotNil(t, model.Components.SecuritySchemes)
	assert.Equal(t, 1, model.Components.SecuritySchemes.Len())

	saslScram, ok := model.Components.SecuritySchemes.Get("saslScram")
	require.True(t, ok)
	assert.Equal(t, "scramSha256", saslScram.Type)

	// Test Components Parameters
	require.NotNil(t, model.Components.Parameters)
	assert.Equal(t, 1, model.Components.Parameters.Len())

	// Test Components OperationTraits
	require.NotNil(t, model.Components.OperationTraits)
	assert.Equal(t, 1, model.Components.OperationTraits.Len())

	kafkaTrait, ok := model.Components.OperationTraits.Get("kafka")
	require.True(t, ok)
	require.NotNil(t, kafkaTrait.Bindings)
	require.NotNil(t, kafkaTrait.Bindings.Kafka)

	// Test Components MessageTraits
	require.NotNil(t, model.Components.MessageTraits)
	assert.Equal(t, 1, model.Components.MessageTraits.Len())

	commonHeaders, ok := model.Components.MessageTraits.Get("commonHeaders")
	require.True(t, ok)
	require.NotNil(t, commonHeaders.Headers)
}

func TestIntegration_GoLowAccess(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test GoLow access
	lowModel := doc.GoLow()
	require.NotNil(t, lowModel)

	// Verify line/column number access on low-level model
	assert.Greater(t, lowModel.Info.KeyNode.Line, 0)
	assert.Greater(t, lowModel.Info.ValueNode.Line, 0)

	// Test that high model GoLow() returns to the same low model
	highModel := doc.Model()
	require.NotNil(t, highModel)

	lowFromHigh := highModel.GoLow()
	assert.Equal(t, lowModel, lowFromHigh)
}

func TestIntegration_Render(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test Render
	rendered, err := doc.Render()
	require.NoError(t, err)
	require.NotEmpty(t, rendered)

	// Rendered output should contain key elements
	assert.Contains(t, string(rendered), "asyncapi: 3.0.0")
	assert.Contains(t, string(rendered), "Streetlights Kafka API")
	assert.Contains(t, string(rendered), "scram-connections")
	assert.Contains(t, string(rendered), "lightingMeasured")
}

func TestIntegration_Serialize(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test Serialize (from original root node)
	serialized, err := doc.Serialize()
	require.NoError(t, err)
	require.NotEmpty(t, serialized)

	// Serialized output should contain key elements
	assert.Contains(t, string(serialized), "asyncapi: 3.0.0")
	assert.Contains(t, string(serialized), "Streetlights Kafka API")
}

func TestIntegration_PartialParse(t *testing.T) {
	// Test a document with some invalid content but still parseable
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Test API
  version: 1.0.0
channels:
  test:
    address: test/channel
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	model := doc.Model()
	require.NotNil(t, model)
	assert.Equal(t, "Test API", model.Info.Title)
	assert.Equal(t, "1.0.0", model.Info.Version)
}

func TestIntegration_MinimalDocument(t *testing.T) {
	spec := []byte(`
asyncapi: 3.0.0
info:
  title: Minimal API
  version: 0.1.0
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.GetVersion())
	assert.False(t, doc.IsPartial())

	model := doc.Model()
	require.NotNil(t, model)
	assert.Equal(t, "Minimal API", model.Info.Title)
	assert.Equal(t, "0.1.0", model.Info.Version)

	// Optional sections should be nil or empty
	assert.Nil(t, model.Servers)
	assert.Nil(t, model.Channels)
	assert.Nil(t, model.Operations)
	assert.Nil(t, model.Components)
}

func TestIntegration_IndexAccess(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test Index access
	idx := doc.Index()
	require.NotNil(t, idx)

	// Index should have reference information
	refs := idx.GetAllReferences()
	assert.NotEmpty(t, refs)
}

func TestIntegration_RolodexAccess(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test Rolodex access
	rolodex := doc.Rolodex()
	require.NotNil(t, rolodex)
}

func TestIntegration_RootNodeAccess(t *testing.T) {
	spec, err := os.ReadFile("test_fixtures/streetlights-kafka.yaml")
	require.NoError(t, err)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	// Test RootNode access
	rootNode := doc.RootNode()
	require.NotNil(t, rootNode)
	assert.NotEmpty(t, rootNode.Content)
}
