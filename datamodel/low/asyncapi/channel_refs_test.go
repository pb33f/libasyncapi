// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func buildLowChannel(t *testing.T, y string) *Channel {
	t.Helper()
	c, err := buildLowChannelWithError(t, y)
	require.NoError(t, err)
	return c
}

func buildLowChannelWithError(t *testing.T, y string) (*Channel, error) {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	c := &Channel{}
	require.NoError(t, low.BuildModel(node.Content[0], c))
	return c, c.Build(context.Background(), nil, node.Content[0], idx)
}

func buildLowOperationReply(t *testing.T, y string) *OperationReply {
	t.Helper()
	or, err := buildLowOperationReplyWithError(t, y)
	require.NoError(t, err)
	return or
}

func buildLowOperationReplyWithError(t *testing.T, y string) (*OperationReply, error) {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	or := &OperationReply{}
	require.NoError(t, low.BuildModel(node.Content[0], or))
	return or, or.Build(context.Background(), nil, node.Content[0], idx)
}

func buildLowOperationWithError(t *testing.T, y string) (*Operation, error) {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	o := &Operation{}
	require.NoError(t, low.BuildModel(node.Content[0], o))
	return o, o.Build(context.Background(), nil, node.Content[0], idx)
}

// channel server entries are mappings containing $ref; the extracted references must
// carry the actual reference strings, not empty values.
func TestChannelBuild_ServerRefsExtracted(t *testing.T) {
	c := buildLowChannel(t, `address: /events
servers:
  - $ref: '#/servers/production'
  - $ref: '#/servers/staging'`)

	require.Len(t, c.Servers.Value, 2)
	assert.Equal(t, "#/servers/production", c.Servers.Value[0].Value.GetReference())
	assert.Equal(t, "#/servers/staging", c.Servers.Value[1].Value.GetReference())
}

// channels referencing different servers must hash differently, otherwise server list
// changes are invisible to hash-based comparison.
func TestChannelHash_ServerRefsAffectHash(t *testing.T) {
	a := buildLowChannel(t, `address: /events
servers:
  - $ref: '#/servers/production'`)
	b := buildLowChannel(t, `address: /events
servers:
  - $ref: '#/servers/staging'`)

	assert.NotEqual(t, a.Hash(), b.Hash())
}

// an explicit `address: null` is semantically distinct from an absent address and must
// hash differently.
func TestChannelHash_NullAddressDistinctFromAbsent(t *testing.T) {
	withNull := buildLowChannel(t, `address: null
title: a channel`)
	absent := buildLowChannel(t, `title: a channel`)

	require.Nil(t, withNull.Address.Value)
	require.NotNil(t, withNull.Address.ValueNode)
	require.Nil(t, absent.Address.ValueNode)
	assert.NotEqual(t, withNull.Hash(), absent.Hash())
}

// an operation reply without a $ref channel must leave the channel reference unset,
// matching Operation.Build behavior.
func TestOperationReplyBuild_NoRefLeavesChannelNil(t *testing.T) {
	or := buildLowOperationReply(t, `address:
  location: $message.header#/replyTo`)
	assert.Nil(t, or.Channel.Value)

	withRef := buildLowOperationReply(t, `channel:
  $ref: '#/channels/replies'
messages:
  - $ref: '#/channels/replies/messages/pong'`)
	require.NotNil(t, withRef.Channel.Value)
	assert.Equal(t, "#/channels/replies", withRef.Channel.Value.GetReference())
	require.Len(t, withRef.Messages.Value, 1)
	assert.Equal(t, "#/channels/replies/messages/pong", withRef.Messages.Value[0].Value.GetReference())
}

func TestChannelBuild_InvalidServerRefsReturnUsefulErrorsAndKeepValidRefs(t *testing.T) {
	c, err := buildLowChannelWithError(t, `servers:
  - '#/servers/scalar'
  - address: inline
  - $ref: ''
  - $ref: '#/servers/valid'`)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "channel servers entry must be a Reference Object")
	assert.Contains(t, err.Error(), "line 2, column 5")
	require.Len(t, c.Servers.Value, 1)
	assert.Equal(t, "#/servers/valid", c.Servers.Value[0].Value.GetReference())

	_, err = buildLowChannelWithError(t, `servers: '#/servers/not-an-array'`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "channel servers must be an array")
}

func TestChannelBuild_AggregatesServerAndTagErrors(t *testing.T) {
	_, err := buildLowChannelWithError(t, `servers: not-an-array
tags: not-an-array`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "channel servers must be an array")
	assert.Contains(t, err.Error(), "array build failed, input is not an array")
}

func TestOperationReplyBuild_InvalidRefsReturnUsefulErrorsAndKeepValidRefs(t *testing.T) {
	or, err := buildLowOperationReplyWithError(t, `channel:
  address: inline
messages:
  - '#/messages/scalar'
  - $ref: '#/messages/valid'
  - $ref: ''`)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "operation reply channel must be a Reference Object")
	assert.Contains(t, err.Error(), "operation reply messages entry must be a Reference Object")
	assert.Nil(t, or.Channel.Value)
	require.Len(t, or.Messages.Value, 1)
	assert.Equal(t, "#/messages/valid", or.Messages.Value[0].Value.GetReference())

	_, err = buildLowOperationReplyWithError(t, `messages: '#/messages/not-an-array'`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "operation reply messages must be an array")
}

func TestOperationBuild_InvalidRefsReturnUsefulErrorsAndKeepValidRefs(t *testing.T) {
	o, err := buildLowOperationWithError(t, `channel:
  address: inline
messages:
  - '#/messages/scalar'
  - $ref: '#/messages/valid'`)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "operation channel must be a Reference Object")
	assert.Contains(t, err.Error(), "operation messages entry must be a Reference Object")
	assert.Nil(t, o.Channel.Value)
	require.Len(t, o.Messages.Value, 1)
	assert.Equal(t, "#/messages/valid", o.Messages.Value[0].Value.GetReference())

	_, err = buildLowOperationWithError(t, `messages: '#/messages/not-an-array'`)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "operation messages must be an array")
}
