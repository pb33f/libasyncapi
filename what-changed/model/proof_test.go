// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func buildInfo(t *testing.T, y string) *lowasync.Info {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	info := &lowasync.Info{}
	require.NoError(t, low.BuildModel(node.Content[0], info))
	require.NoError(t, info.Build(context.Background(), nil, node.Content[0], idx))
	return info
}

func buildSchemaProxy(t *testing.T, y string) *base.SchemaProxy {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	sp := &base.SchemaProxy{}
	require.NoError(t, sp.Build(context.Background(), nil, node.Content[0], idx))
	return sp
}

// proof point 1: a scalar property change is detected with line context.
func TestProof_ScalarChangeDetected(t *testing.T) {
	left := buildInfo(t, `title: Old API
version: 1.0.0
description: an api`)
	right := buildInfo(t, `title: New API
version: 1.0.0
description: an api`)

	changes := CompareInfo(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.Equal(t, "Old API", ch.Original)
	assert.Equal(t, "New API", ch.New)
	require.NotNil(t, ch.Context.OriginalLine)
	assert.Equal(t, 1, *ch.Context.OriginalLine)
}

// proof point 2: hash-based map diffing works on AsyncAPI types (uint64 Hash contract).
// Info tags are an ordered map of AsyncAPI Tag objects compared via CheckMapForChanges.
func TestProof_MapChangeDetected(t *testing.T) {
	left := buildInfo(t, `title: api
version: 1.0.0
tags:
  - name: one
    description: first tag
  - name: two
    description: second tag`)
	right := buildInfo(t, `title: api
version: 1.0.0
tags:
  - name: one
    description: first tag, changed
  - name: three
    description: third tag`)

	changes := CompareInfo(left, right)
	require.NotNil(t, changes)

	// tag "one" was modified, tag "two" removed, tag "three" added.
	require.Contains(t, changes.TagChanges, "one")
	assert.Equal(t, 1, changes.TagChanges["one"].TotalChanges())

	var removed, added int
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case ObjectRemoved:
			removed++
		case ObjectAdded:
			added++
		}
	}
	assert.Equal(t, 1, removed)
	assert.Equal(t, 1, added)
	assert.Equal(t, 3, changes.TotalChanges())
}

// proof point 3: a custom AsyncAPI breaking rule set on THIS package's config is honored
// by the localized helpers (and not swallowed by libopenapi's global config).
func TestProof_CustomBreakingRuleHonored(t *testing.T) {
	left := buildInfo(t, `title: Old API
version: 1.0.0`)
	right := buildInfo(t, `title: New API
version: 1.0.0`)

	// default: info.title modifications are not breaking.
	changes := CompareInfo(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	// flip the rule: info.title modifications are now breaking.
	custom := NewDefaultBreakingRulesConfig()
	custom.Merge(&BreakingRulesConfig{
		Info: &InfoRules{Title: rule(false, true, false)},
	})
	SetActiveBreakingRulesConfig(custom)
	defer ResetActiveBreakingRulesConfig()

	changes = CompareInfo(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.True(t, changes.Changes[0].Breaking)
}

// proof point 4: schema comparison is delegated to libopenapi's battle-tested comparator
// and reports changes on base.SchemaProxy objects as used by AsyncAPI message payloads.
func TestProof_DelegatedSchemaDiff(t *testing.T) {
	left := buildSchemaProxy(t, `type: object
properties:
  id:
    type: string
  total:
    type: integer`)
	right := buildSchemaProxy(t, `type: object
properties:
  id:
    type: string
  total:
    type: number`)

	sc := wcmodel.CompareSchemas(left, right)
	require.NotNil(t, sc)
	assert.Positive(t, sc.TotalChanges())
	assert.Positive(t, sc.TotalBreakingChanges())
}

// the converted uint64 Hash methods must satisfy low.Hashable, otherwise none of the
// hash-based map diffing above can work. compile-time assertions:
var (
	_ low.Hashable = (*lowasync.AsyncAPI)(nil)
	_ low.Hashable = (*lowasync.Info)(nil)
	_ low.Hashable = (*lowasync.Tag)(nil)
	_ low.Hashable = (*lowasync.Channel)(nil)
	_ low.Hashable = (*lowasync.Operation)(nil)
	_ low.Hashable = (*lowasync.Message)(nil)
	_ low.Hashable = (*lowasync.Components)(nil)
	_ low.Hashable = (*lowasync.Server)(nil)
)
