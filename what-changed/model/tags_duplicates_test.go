// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestCompareTagSlices_DuplicateModificationIsNotCollapsed(t *testing.T) {
	left := buildChannel(t, `tags:
  - name: env
    description: first
  - name: env
    description: second`)
	right := buildChannel(t, `tags:
  - name: env
    description: first, changed
  - name: env
    description: second`)

	changes := CompareChannels(left, right)
	require.NotNil(t, changes)
	require.Contains(t, changes.TagChanges, "env#1")
	assert.Equal(t, 1, changes.TagChanges["env#1"].TotalChanges())
	assert.Equal(t, 1, changes.TotalChanges())
}

func TestCompareTagSlices_DuplicateAndUnnamedRemovalIsReported(t *testing.T) {
	left := buildChannel(t, `tags:
  - name: env
  - name: env
  - description: unnamed one
  - description: unnamed two`)
	right := buildChannel(t, `tags:
  - name: env
  - description: unnamed one`)

	changes := CompareChannels(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	require.Len(t, changes.Changes, 2)
	assert.Equal(t, ObjectRemoved, changes.Changes[0].ChangeType)
	assert.Equal(t, ObjectRemoved, changes.Changes[1].ChangeType)
}

func TestCompareTagSlices_HandlesEmptyAndNilEntries(t *testing.T) {
	var changes []*Change
	assert.Nil(t, CompareTagSlices(nil, nil, &changes, CompChannel, PropTags))

	actual := buildChannel(t, `tags:
  - description: unnamed`).Tags.Value[0]
	nilTag := low.ValueReference[*lowasync.Tag]{}

	assert.Nil(t, CompareTagSlices([]low.ValueReference[*lowasync.Tag]{nilTag},
		[]low.ValueReference[*lowasync.Tag]{nilTag}, &changes, CompChannel, PropTags))
	assert.Empty(t, changes)

	CompareTagSlices([]low.ValueReference[*lowasync.Tag]{nilTag},
		[]low.ValueReference[*lowasync.Tag]{actual}, &changes, CompChannel, PropTags)
	require.Len(t, changes, 1)
	assert.Equal(t, ObjectAdded, changes[0].ChangeType)

	changes = nil
	CompareTagSlices([]low.ValueReference[*lowasync.Tag]{actual},
		[]low.ValueReference[*lowasync.Tag]{nilTag}, &changes, CompChannel, PropTags)
	require.Len(t, changes, 1)
	assert.Equal(t, ObjectRemoved, changes[0].ChangeType)
}

func TestTagOccurrenceKey_AvoidsExistingKeys(t *testing.T) {
	existing := map[string]*TagChanges{"env#1": {}}
	assert.Equal(t, "env#2", tagOccurrenceKey("env", 0, 2, 2, existing))
	assert.Equal(t, "<unnamed>#1", tagOccurrenceKey("", 0, 2, 2, nil))
}

func TestCompareTagSlices_DuplicateAdditionIsReported(t *testing.T) {
	left := buildChannel(t, `tags:
  - name: env`)
	right := buildChannel(t, `tags:
  - name: env
  - name: env`)

	changes := CompareChannels(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	require.Len(t, changes.Changes, 1)
	assert.Equal(t, ObjectAdded, changes.Changes[0].ChangeType)
}
