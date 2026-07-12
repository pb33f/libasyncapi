// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

func buildChannel(t *testing.T, y string) *lowasync.Channel {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	c := &lowasync.Channel{}
	require.NoError(t, low.BuildModel(node.Content[0], c))
	require.NoError(t, c.Build(context.Background(), nil, node.Content[0], idx))
	return c
}

// regression for the review repro: swapping a channel's server reference previously
// produced no changes at all because the extracted refs were empty strings.
func TestCompareChannels_ServerRefSwapDetected(t *testing.T) {
	left := buildChannel(t, `address: /events
servers:
  - $ref: '#/servers/production'`)
	right := buildChannel(t, `address: /events
servers:
  - $ref: '#/servers/staging'`)

	changes := CompareChannels(left, right)
	require.NotNil(t, changes)
	// a swap is one removal plus one addition; channel.servers removal is breaking.
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	var added, removed *Change
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case ObjectAdded:
			added = ch
		case ObjectRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
}

func TestCompareChannels_ServerRefAddedAndRemoved(t *testing.T) {
	one := buildChannel(t, `address: /events
servers:
  - $ref: '#/servers/production'`)
	two := buildChannel(t, `address: /events
servers:
  - $ref: '#/servers/production'
  - $ref: '#/servers/staging'`)

	// adding a server: not breaking.
	changes := CompareChannels(one, two)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	// removing a server: breaking.
	changes = CompareChannels(two, one)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
}

// the nullable address has three states: absent, explicit null, and a value. every
// transition must be reported.
func TestCompareChannels_AddressStateTransitions(t *testing.T) {
	withValue := `address: /orders
title: orders`
	withNull := `address: null
title: orders`
	absent := `title: orders`

	cases := []struct {
		name       string
		left       string
		right      string
		changeType int
	}{
		{"value modified", withValue, "address: /payments\ntitle: orders", Modified},
		{"absent to value", absent, withValue, PropertyAdded},
		{"value to absent", withValue, absent, PropertyRemoved},
		{"null to value", withNull, withValue, Modified},
		{"value to null", withValue, withNull, Modified},
		{"absent to null", absent, withNull, PropertyAdded},
		{"null to absent", withNull, absent, PropertyRemoved},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			changes := CompareChannels(buildChannel(t, tc.left), buildChannel(t, tc.right))
			require.NotNil(t, changes, "transition must be detected")
			require.Len(t, changes.Changes, 1)
			assert.Equal(t, tc.changeType, changes.Changes[0].ChangeType)
			assert.True(t, changes.Changes[0].Breaking, "address transitions are breaking")
		})
	}

	// no-change cases: identical states produce nil.
	assert.Nil(t, CompareChannels(buildChannel(t, withNull), buildChannel(t, withNull)))
	assert.Nil(t, CompareChannels(buildChannel(t, absent), buildChannel(t, absent)))
}
