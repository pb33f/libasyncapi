// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func buildContact(t *testing.T, y string) *base.Contact {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	c := &base.Contact{}
	require.NoError(t, low.BuildModel(node.Content[0], c))
	require.NoError(t, c.Build(context.Background(), nil, node.Content[0], idx))
	return c
}

func buildLicense(t *testing.T, y string) *base.License {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	l := &base.License{}
	require.NoError(t, low.BuildModel(node.Content[0], l))
	require.NoError(t, l.Build(context.Background(), nil, node.Content[0], idx))
	return l
}

func TestCompareContact(t *testing.T) {
	left := buildContact(t, `name: dave
url: https://pb33f.io
email: dave@pb33f.io`)
	right := buildContact(t, `name: dave
url: https://pb33f.io
email: dave@quobix.com`)

	changes := CompareContact(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	assert.Equal(t, "dave@pb33f.io", changes.Changes[0].Original)
	assert.Equal(t, "dave@quobix.com", changes.Changes[0].New)

	assert.Nil(t, CompareContact(left, left))
	assert.Nil(t, CompareContact(nil, right))
}

func TestCompareLicense(t *testing.T) {
	left := buildLicense(t, `name: MIT
identifier: MIT`)
	right := buildLicense(t, `name: Apache 2.0
identifier: Apache-2.0`)

	changes := CompareLicense(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	assert.Nil(t, CompareLicense(left, left))
	assert.Nil(t, CompareLicense(left, nil))
}

// regression for the rule bleed-through finding: a custom rule set on THIS package's
// config must govern contact/license breaking classification. before localization the
// modified-rules resolved against libopenapi's global OpenAPI config and the local
// knobs silently did nothing.
func TestCompareContact_CustomRuleHonored(t *testing.T) {
	left := buildContact(t, `email: dave@pb33f.io`)
	right := buildContact(t, `email: dave@quobix.com`)

	changes := CompareContact(left, right)
	require.NotNil(t, changes)
	require.Equal(t, 0, changes.TotalBreakingChanges())

	custom := NewDefaultBreakingRulesConfig()
	custom.Merge(&BreakingRulesConfig{
		Contact: &ContactRules{Email: rule(false, true, false)},
	})
	SetActiveBreakingRulesConfig(custom)
	defer ResetActiveBreakingRulesConfig()

	changes = CompareContact(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.True(t, changes.Changes[0].Breaking)
}
