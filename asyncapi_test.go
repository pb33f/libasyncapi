// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestNewDocument_ValidAsyncAPI3(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
channels:
  userSignups:
    address: user/signedup
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.GetVersion())
	assert.NotNil(t, doc.GetSpecInfo())
	assert.Equal(t, "asyncapi", doc.GetSpecInfo().SpecType)
	assert.NotNil(t, doc.Model())
	assert.NotNil(t, doc.Index())
	assert.NotNil(t, doc.Rolodex())
	assert.NotNil(t, doc.RootNode())
	assert.False(t, doc.IsPartial())
	assert.Empty(t, doc.Errors())
}

func TestNewDocument_AsyncAPI2NotSupported(t *testing.T) {
	spec := []byte(`asyncapi: "2.6.0"
info:
  title: Test API
  version: "1.0.0"
`)

	doc, err := NewDocument(spec)
	assert.ErrorIs(t, err, ErrAsyncAPI2NotSupported)
	assert.Nil(t, doc)
}

func TestNewDocument_MissingAsyncAPIVersion(t *testing.T) {
	spec := []byte(`openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	doc, err := NewDocument(spec)
	assert.ErrorIs(t, err, ErrNoAsyncAPIVersion)
	assert.Nil(t, doc)
}

func TestNewDocument_InvalidYAML(t *testing.T) {
	spec := []byte(`{{{invalid yaml`)

	doc, err := NewDocument(spec)
	assert.ErrorIs(t, err, ErrInvalidYAML)
	assert.Nil(t, doc)
}

func TestNewDocument_EmptySpec(t *testing.T) {
	spec := []byte(``)

	doc, err := NewDocument(spec)
	assert.ErrorIs(t, err, ErrInvalidYAML)
	assert.Nil(t, doc)
}

func TestNewDocumentWithConfiguration(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	config := NewDocumentConfiguration()
	config.AllowFileReferences = false
	config.AllowRemoteReferences = false

	doc, err := NewDocumentWithConfiguration(spec, config)
	require.NoError(t, err)
	require.NotNil(t, doc)

	assert.Equal(t, "3.0.0", doc.GetVersion())
}

func TestDocument_Serialize(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	serialized, err := doc.Serialize()
	require.NoError(t, err)
	assert.Contains(t, string(serialized), "asyncapi")
	assert.Contains(t, string(serialized), "3.0.0")
}

func TestNewDocumentConfiguration(t *testing.T) {
	config := NewDocumentConfiguration()
	require.NotNil(t, config)
	assert.False(t, config.AllowFileReferences)
	assert.False(t, config.AllowRemoteReferences)
	assert.NotNil(t, config.Logger)
}

func TestDocumentConfiguration_ToLibOpenAPIConfig(t *testing.T) {
	config := &DocumentConfiguration{
		BasePath:                   "/test/path",
		AllowFileReferences:        true,
		AllowRemoteReferences:      true,
		SkipCircularReferenceCheck: true,
	}
	loConfig := config.ToLibOpenAPIConfig()
	require.NotNil(t, loConfig)
	assert.Equal(t, "/test/path", loConfig.BasePath)
	assert.True(t, loConfig.AllowFileReferences)
	assert.True(t, loConfig.AllowRemoteReferences)
	assert.True(t, loConfig.SkipCircularReferenceCheck)
}

func TestDocumentConfiguration_ToLibOpenAPIConfig_Nil(t *testing.T) {
	var config *DocumentConfiguration
	loConfig := config.ToLibOpenAPIConfig()
	require.NotNil(t, loConfig)
}

func TestDocumentConfiguration_ToLibOpenAPIConfig_CustomHandlers(t *testing.T) {
	customHandler := func(url string) (*http.Response, error) {
		return nil, nil
	}
	customFS := os.DirFS(".")

	config := &DocumentConfiguration{
		RemoteURLHandler: customHandler,
		LocalFS:          customFS,
	}

	loConfig := config.ToLibOpenAPIConfig()
	require.NotNil(t, loConfig)
	assert.NotNil(t, loConfig.RemoteURLHandler, "RemoteURLHandler should be propagated")
	assert.NotNil(t, loConfig.LocalFS, "LocalFS should be propagated")
}

func TestNewDocumentWithConfiguration_CustomLocalFS(t *testing.T) {
	// Use the test_fixtures directory as a custom FS
	customFS := os.DirFS("test_fixtures")

	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	config := &DocumentConfiguration{
		AllowFileReferences: true,
		LocalFS:             customFS,
	}

	doc, err := NewDocumentWithConfiguration(spec, config)
	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.Equal(t, "3.0.0", doc.GetVersion())
}

func TestDocument_PartialParse_UnresolvedReference(t *testing.T) {
	// A spec with an unresolvable reference should still parse partially
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
channels:
  testChannel:
    address: test/channel
    messages:
      testMessage:
        $ref: '#/components/messages/NonExistent'
components:
  messages: {}
`)

	doc, err := NewDocument(spec)

	// Document should be created even with reference errors (partial parse)
	require.NoError(t, err, "Partial parse should not return error from NewDocument")
	require.NotNil(t, doc, "Document should be returned even with unresolved references")

	// Document should be marked as partial with errors
	assert.True(t, doc.IsPartial(), "Document with unresolved references should be partial")
	assert.NotEmpty(t, doc.Errors(), "Document with unresolved references should have errors")
	assert.Greater(t, len(doc.Errors()), 0, "Should have at least one error for unresolved reference")

	// Model should still be accessible
	assert.NotNil(t, doc.Model(), "Model should still be accessible on partial document")
	assert.Equal(t, "3.0.0", doc.GetVersion(), "Version should still be parsed")
}

func TestDocument_Errors_EmptyForValidSpec(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
channels:
  userChannel:
    address: user/events
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc)

	// Valid spec should have no errors
	assert.Empty(t, doc.Errors(), "Valid spec should have no errors")
	assert.False(t, doc.IsPartial(), "Valid spec should not be partial")
}

func TestDocument_InvalidReferenceOnlyFieldsProducePartialParseDiagnostics(t *testing.T) {
	doc, err := NewDocument([]byte(`asyncapi: 3.0.0
info:
  title: Invalid references
  version: 1.0.0
channels:
  events:
    servers:
      - '#/servers/scalar'
operations:
  receive:
    action: receive
    channel:
      $ref: '#/channels/events'
    reply:
      channel:
        address: inline
`))

	require.NoError(t, err)
	require.NotNil(t, doc)
	assert.True(t, doc.IsPartial())
	require.NotEmpty(t, doc.Errors())
	joined := errors.Join(doc.Errors()...).Error()
	assert.Contains(t, joined, "channel servers entry must be a Reference Object")
	assert.Contains(t, joined, "operation reply channel must be a Reference Object")
}

func TestDocument_Render_ReturnsModelBasedYAML(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	rendered, err := doc.Render()
	require.NoError(t, err)
	assert.NotEmpty(t, rendered)
	assert.Contains(t, string(rendered), "asyncapi")
}

func TestDocument_Serialize_PreservesOriginalStructure(t *testing.T) {
	// Serialize should preserve original YAML structure
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
  # This comment should be preserved in serialization
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	serialized, err := doc.Serialize()
	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
	// Should contain original content
	assert.Contains(t, string(serialized), "Test API")
}

func TestDocument_GoLow_ReturnsLowLevelModel(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err)

	low := doc.GoLow()
	require.NotNil(t, low)
	assert.Equal(t, "3.0.0", low.AsyncAPI.Value)
}

func TestDocument_FileRefResolution(t *testing.T) {
	// Test that file $ref resolution works with AllowFileReferences
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API with File Refs
  version: "1.0.0"
channels:
  events:
    address: events/stream
    messages:
      externalMessage:
        $ref: 'shared-message.yaml'
`)

	config := &DocumentConfiguration{
		BasePath:            "test_fixtures",
		AllowFileReferences: true,
	}

	doc, err := NewDocumentWithConfiguration(spec, config)
	require.NoError(t, err)
	require.NotNil(t, doc)

	// Verify the document parsed successfully
	assert.Equal(t, "3.0.0", doc.GetVersion())
	assert.NotNil(t, doc.Model())

	// Verify channels were parsed
	require.NotNil(t, doc.Model().Channels)
	channel, ok := doc.Model().Channels.Get("events")
	require.True(t, ok, "Channel 'events' should exist")
	require.NotNil(t, channel.Address, "Channel address should not be nil")
	assert.Equal(t, "events/stream", *channel.Address)

	// Verify messages exist (reference was resolved)
	require.NotNil(t, channel.Messages)
	assert.Greater(t, channel.Messages.Len(), 0, "Channel should have messages after ref resolution")
}

func TestDocument_MultipleErrorsAggregated(t *testing.T) {
	// Test that multiple parse errors are aggregated (not just the first one)
	// This exercises the errors.Join behavior in asyncapi.go and components.go
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
channels:
  channel1:
    address: test/channel1
    messages:
      msg1:
        $ref: '#/components/messages/NonExistent1'
  channel2:
    address: test/channel2
    messages:
      msg2:
        $ref: '#/components/messages/NonExistent2'
  channel3:
    address: test/channel3
    messages:
      msg3:
        $ref: '#/components/messages/NonExistent3'
components:
  messages: {}
`)

	doc, err := NewDocument(spec)
	require.NoError(t, err, "Document should parse despite reference errors")
	require.NotNil(t, doc)

	// Should be partial due to unresolved references
	assert.True(t, doc.IsPartial(), "Document should be partial")

	// Should have multiple errors aggregated (one for each unresolved ref)
	errors := doc.Errors()
	assert.NotEmpty(t, errors, "Should have errors for unresolved references")

	// Verify we got multiple errors, not just the first one
	// This confirms errors.Join is working to aggregate all errors
	assert.GreaterOrEqual(t, len(errors), 3,
		"Should have at least 3 errors (one per unresolved ref), got %d: %v", len(errors), errors)
}
