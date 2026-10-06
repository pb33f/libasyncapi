// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi_test

import (
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBindingNamesIncludesUntypedProtocols(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: 3.0.0
info:
  title: Binding names
  version: 1.0.0
servers:
  events:
    host: events.example.com
    protocol: nats
    bindings:
      nats: {}
      x-provider: ignored
channels:
  events:
    address: events
    bindings:
      googlepubsub: {}
      kafka: {}
    messages:
      event:
        $ref: '#/components/messages/event'
operations:
  receiveEvent:
    action: receive
    channel:
      $ref: '#/channels/events'
    bindings:
      pulsar: {}
      redis: {}
    messages:
      - $ref: '#/components/messages/event'
components:
  messages:
    event:
      bindings:
        sns: {}
        sqs: {}
      payload:
        type: object
`))
	require.NoError(t, err)

	model := doc.Model()
	require.NotNil(t, model)
	server, ok := model.Servers.Get("events")
	require.True(t, ok)
	assert.Equal(t, []string{"nats"}, server.Bindings.BindingNames())
	channel, ok := model.Channels.Get("events")
	require.True(t, ok)
	assert.Equal(t, []string{"googlepubsub", "kafka"}, channel.Bindings.BindingNames())
	operation, ok := model.Operations.Get("receiveEvent")
	require.True(t, ok)
	assert.Equal(t, []string{"pulsar", "redis"}, operation.Bindings.BindingNames())
	message, ok := model.Components.Messages.Get("event")
	require.True(t, ok)
	assert.Equal(t, []string{"sns", "sqs"}, message.Bindings.BindingNames())
}
