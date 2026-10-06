// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func buildSQSPolicy(t *testing.T, y string) *lowasync.SQSPolicy {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	p := &lowasync.SQSPolicy{}
	require.NoError(t, low.BuildModel(node.Content[0], p))
	require.NoError(t, p.Build(context.Background(), nil, node.Content[0], idx))
	return p
}

func buildMessage(t *testing.T, y string) *lowasync.Message {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	m := &lowasync.Message{}
	require.NoError(t, low.BuildModel(node.Content[0], m))
	require.NoError(t, m.Build(context.Background(), nil, node.Content[0], idx))
	return m
}

// regression for the dead-comparator finding: editing one field of one statement must
// produce a granular modification under sqsPolicyStatement rules, not a whole-object
// removal plus addition.
func TestCompareSQSPolicy_StatementEditIsGranular(t *testing.T) {
	left := buildSQSPolicy(t, `statements:
  - effect: Allow
    principal: events.amazonaws.com
    action: sqs:SendMessage
  - effect: Deny
    principal: '*'
    action: sqs:*`)
	right := buildSQSPolicy(t, `statements:
  - effect: Deny
    principal: events.amazonaws.com
    action: sqs:SendMessage
  - effect: Deny
    principal: '*'
    action: sqs:*`)

	changes := CompareSQSPolicy(left, right)
	require.NotNil(t, changes)
	require.Len(t, changes.StatementChanges, 1)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.StatementChanges[0].Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.Equal(t, "Allow", ch.Original)
	assert.Equal(t, "Deny", ch.New)
}

func TestCompareSQSPolicy_StatementAddedAndRemoved(t *testing.T) {
	one := buildSQSPolicy(t, `statements:
  - effect: Allow
    principal: events.amazonaws.com`)
	two := buildSQSPolicy(t, `statements:
  - effect: Allow
    principal: events.amazonaws.com
  - effect: Deny
    principal: '*'`)

	// statement added: breaking per sqsPolicy.statements rules (additions tighten policy).
	changes := CompareSQSPolicy(one, two)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, ObjectAdded, changes.Changes[0].ChangeType)

	// statement removed: breaking.
	changes = CompareSQSPolicy(two, one)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, ObjectRemoved, changes.Changes[0].ChangeType)
	assert.True(t, changes.Changes[0].Breaking)
}

// custom sqsPolicyStatement rules must now be honored — the rules group was dead before
// the comparator was wired in.
func TestCompareSQSPolicy_CustomStatementRuleHonored(t *testing.T) {
	left := buildSQSPolicy(t, `statements:
  - effect: Allow`)
	right := buildSQSPolicy(t, `statements:
  - effect: Deny`)

	// default: effect modification is breaking.
	changes := CompareSQSPolicy(left, right)
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalBreakingChanges())

	custom := NewDefaultBreakingRulesConfig()
	custom.Merge(&BreakingRulesConfig{
		SQSPolicyStatement: &SQSPolicyStatementRules{Effect: rule(false, false, false)},
	})
	SetActiveBreakingRulesConfig(custom)
	defer ResetActiveBreakingRulesConfig()

	changes = CompareSQSPolicy(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

// regression for the dead-comparator finding: editing a named example must produce a
// granular modification under messageExample rules, not a remove+add pair.
func TestCompareMessages_NamedExampleEditIsGranular(t *testing.T) {
	left := buildMessage(t, `name: order
examples:
  - name: minimal
    summary: a minimal order
    payload:
      id: 1`)
	right := buildMessage(t, `name: order
examples:
  - name: minimal
    summary: a minimal order, updated
    payload:
      id: 1`)

	changes := CompareMessages(left, right)
	require.NotNil(t, changes)
	require.Contains(t, changes.ExampleChanges, "minimal")
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	ch := changes.ExampleChanges["minimal"].Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.Equal(t, "a minimal order", ch.Original)
}

func TestCompareMessages_ExampleAddedRemovedAndUnnamed(t *testing.T) {
	named := buildMessage(t, `name: order
examples:
  - name: minimal
    payload:
      id: 1`)
	namedPlusUnnamed := buildMessage(t, `name: order
examples:
  - name: minimal
    payload:
      id: 1
  - payload:
      id: 2`)

	// unnamed example added.
	changes := CompareMessages(named, namedPlusUnnamed)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, ObjectAdded, changes.Changes[0].ChangeType)

	// unnamed example removed.
	changes = CompareMessages(namedPlusUnnamed, named)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, ObjectRemoved, changes.Changes[0].ChangeType)
}

func TestCompareMessages_DuplicateNamedExamplesRetainEveryOccurrence(t *testing.T) {
	left := buildMessage(t, `name: order
examples:
  - name: duplicate
    summary: before
  - name: duplicate
    summary: unchanged`)
	right := buildMessage(t, `name: order
examples:
  - name: duplicate
    summary: after
  - name: duplicate
    summary: unchanged`)

	changes := CompareMessages(left, right)
	require.NotNil(t, changes)
	require.Contains(t, changes.ExampleChanges, "duplicate#1")
	assert.Equal(t, 1, changes.ExampleChanges["duplicate#1"].TotalChanges())
	assert.Equal(t, 1, changes.TotalChanges())
}

func TestCompareMessages_DuplicateNamedExampleAdditionAndRemoval(t *testing.T) {
	one := buildMessage(t, `name: order
examples:
  - name: duplicate
    summary: first`)
	two := buildMessage(t, `name: order
examples:
  - name: duplicate
    summary: first
  - name: duplicate
    summary: second`)

	added := CompareMessages(one, two)
	require.NotNil(t, added)
	require.Len(t, added.Changes, 1)
	assert.Equal(t, ObjectAdded, added.Changes[0].ChangeType)

	removed := CompareMessages(two, one)
	require.NotNil(t, removed)
	require.Len(t, removed.Changes, 1)
	assert.Equal(t, ObjectRemoved, removed.Changes[0].ChangeType)
}

func TestGroupMessageExamples_IgnoresNilEntries(t *testing.T) {
	grouped, order := groupMessageExamples(
		[]low.ValueReference[*lowasync.MessageExample]{{}},
		func(low.ValueReference[*lowasync.MessageExample]) string {
			t.Fatal("key function must not run for nil examples")
			return ""
		},
	)
	assert.Empty(t, grouped)
	assert.Empty(t, order)
}
