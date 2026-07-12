// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"sync"
	"testing"

	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"
)

type comparisonMapValue struct {
	hash uint64
}

func (v *comparisonMapValue) Hash() uint64 { return v.hash }

func comparisonMap(entries ...struct {
	key   string
	value *comparisonMapValue
}) *orderedmap.Map[low.KeyReference[string], low.ValueReference[*comparisonMapValue]] {
	m := orderedmap.New[low.KeyReference[string], low.ValueReference[*comparisonMapValue]]()
	for _, entry := range entries {
		keyNode := &yaml.Node{Kind: yaml.ScalarNode, Value: entry.key}
		valueNode := &yaml.Node{Kind: yaml.MappingNode}
		m.Set(low.KeyReference[string]{Value: entry.key, KeyNode: keyNode},
			low.ValueReference[*comparisonMapValue]{Value: entry.value, ValueNode: valueNode})
	}
	return m
}

func mapEntry(key string, hash uint64) struct {
	key   string
	value *comparisonMapValue
} {
	return struct {
		key   string
		value *comparisonMapValue
	}{key: key, value: &comparisonMapValue{hash: hash}}
}

func nilMapEntry(key string) struct {
	key   string
	value *comparisonMapValue
} {
	return struct {
		key   string
		value *comparisonMapValue
	}{key: key}
}

func compareMapValue(l, r *comparisonMapValue) *struct{} {
	if l == nil || r == nil || l.hash == r.hash {
		return nil
	}
	return &struct{}{}
}

func TestCheckMapForChangesInternal_IsReadOnlyAndSourceOrdered(t *testing.T) {
	left := comparisonMap(mapEntry("first", 1), mapEntry("second", 2))
	right := comparisonMap(mapEntry("third", 3))

	var changes []*Change
	result := checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
		true, false, true)

	assert.Empty(t, result)
	require.Len(t, changes, 3)
	assert.Equal(t, []int{ObjectRemoved, ObjectRemoved, ObjectAdded}, []int{
		changes[0].ChangeType, changes[1].ChangeType, changes[2].ChangeType,
	})
	assert.Equal(t, []any{"first", "second", "third"}, []any{
		changes[0].Original, changes[1].Original, changes[2].New,
	})
	for _, value := range left.FromOldest() {
		assert.Empty(t, value.ValueNode.Value, "comparison must not mutate mapping nodes")
	}
	for _, value := range right.FromOldest() {
		assert.Empty(t, value.ValueNode.Value, "comparison must not mutate mapping nodes")
	}
}

func TestCheckMapForChangesInternal_TypedNilIsNotMissingOrPanicking(t *testing.T) {
	left := comparisonMap(nilMapEntry("nil"))
	right := comparisonMap(nilMapEntry("nil"))

	var changes []*Change
	assert.NotPanics(t, func() {
		result := checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
			true, false, true)
		assert.Empty(t, result)
	})
	assert.Empty(t, changes)

	right = comparisonMap(mapEntry("nil", 1))
	changes = nil
	checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
		true, false, true)
	require.Len(t, changes, 1)
	assert.Equal(t, ObjectAdded, changes[0].ChangeType)

	left = comparisonMap(mapEntry("nil", 1))
	right = comparisonMap(nilMapEntry("nil"))
	changes = nil
	checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
		true, false, true)
	require.Len(t, changes, 1)
	assert.Equal(t, ObjectRemoved, changes[0].ChangeType)
}

func TestCheckMapForChangesInternal_ModifiedAndStructureOnlyModes(t *testing.T) {
	left := comparisonMap(mapEntry("same-key", 1))
	right := comparisonMap(mapEntry("same-key", 2))

	var changes []*Change
	result := checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
		true, false, true)
	require.Contains(t, result, "same-key")
	assert.Empty(t, changes)

	result = checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
		false, false, true)
	assert.Empty(t, result, "structure-only comparisons ignore value modifications")
}

func TestChangeNode_PrefersUsefulValueThenKeyThenOriginalNode(t *testing.T) {
	value := &yaml.Node{Kind: yaml.ScalarNode, Value: "value"}
	key := &yaml.Node{Kind: yaml.ScalarNode, Value: "key"}
	mapping := &yaml.Node{Kind: yaml.MappingNode}

	assert.Same(t, value, changeNode(value, key))
	assert.Same(t, key, changeNode(mapping, key))
	assert.Same(t, mapping, changeNode(mapping, nil))
	assert.Nil(t, changeNode(nil, nil))
}

func TestCheckMapForChangesInternal_SharedInputsAreConcurrentSafe(t *testing.T) {
	left := comparisonMap(mapEntry("left", 1))
	right := comparisonMap(mapEntry("right", 2))

	const comparisons = 64
	var wg sync.WaitGroup
	errs := make(chan string, comparisons)
	for range comparisons {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var changes []*Change
			checkMapForChangesInternal(left, right, &changes, "items", compareMapValue,
				true, false, true)
			if len(changes) != 2 || changes[0].Original != "left" || changes[1].New != "right" {
				errs <- "unexpected comparison result"
			}
		}()
	}
	wg.Wait()
	close(errs)
	assert.Empty(t, errs)
}

func TestExtractStringValueSliceChangesWithRules_PreservesSourceOrder(t *testing.T) {
	values := func(items ...string) []low.ValueReference[string] {
		refs := make([]low.ValueReference[string], 0, len(items))
		for _, item := range items {
			refs = append(refs, low.ValueReference[string]{
				Value: item, ValueNode: &yaml.Node{Kind: yaml.ScalarNode, Value: item},
			})
		}
		return refs
	}

	var changes []*Change
	ExtractStringValueSliceChangesWithRules(values("left-b", "left-a"), values("right-d", "right-c"),
		&changes, "values", CompParameter, PropEnum)
	require.Len(t, changes, 4)
	assert.Equal(t, []any{"left-b", "left-a", "right-d", "right-c"}, []any{
		changes[0].Original, changes[1].Original, changes[2].New, changes[3].New,
	})
}
