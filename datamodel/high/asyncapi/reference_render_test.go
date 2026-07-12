// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi_test

import (
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReferenceOnlyFields_RenderAndReparseWithoutLoss(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: 3.0.0
info:
  title: Reference rendering
  version: 1.0.0
servers:
  production:
    host: example.com
    protocol: kafka
channels:
  events:
    servers:
      - $ref: '#/servers/production'
    messages:
      event:
        payload:
          type: string
      reply:
        payload:
          type: string
operations:
  receive:
    action: receive
    channel:
      $ref: '#/channels/events'
    messages:
      - $ref: '#/channels/events/messages/event'
    reply:
      channel:
        $ref: '#/channels/events'
      messages:
        - $ref: '#/channels/events/messages/reply'
`))
	require.NoError(t, err)
	require.False(t, doc.IsPartial())

	rendered, err := doc.Render()
	require.NoError(t, err)
	assert.NotContains(t, string(rendered), "channel: {}")
	assert.NotContains(t, string(rendered), "- {}")
	assert.Contains(t, string(rendered), "$ref: '#/channels/events'")

	roundTrip, err := libasyncapi.NewDocument(rendered)
	require.NoError(t, err)
	require.False(t, roundTrip.IsPartial(), "rendered references must remain valid")

	channel, ok := roundTrip.Model().Channels.Get("events")
	require.True(t, ok)
	require.Len(t, channel.Servers, 1)
	assert.Equal(t, "#/servers/production", channel.Servers[0].GetReference())

	operation, ok := roundTrip.Model().Operations.Get("receive")
	require.True(t, ok)
	require.NotNil(t, operation.Channel)
	assert.Equal(t, "#/channels/events", operation.Channel.GetReference())
	require.Len(t, operation.Messages, 1)
	require.NotNil(t, operation.Reply)
	require.NotNil(t, operation.Reply.Channel)
	require.Len(t, operation.Reply.Messages, 1)

	replyYAML, err := operation.Reply.Render()
	require.NoError(t, err)
	assert.Contains(t, string(replyYAML), "$ref: '#/channels/events'")
}
