// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"fmt"

	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// MessageExampleChanges represents changes made to an AsyncAPI Message Example object.
type MessageExampleChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Message Example objects.
func (m *MessageExampleChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MessageExampleChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MessageExampleChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMessageExampleSlices compares two slices of message example value references,
// keyed by example name. Named examples present on both sides are compared field by
// field with CompareMessageExamples; unnamed examples fall back to hash identity, so an
// edited unnamed example reports as one removal plus one addition. Returns a map of
// example name (or hash for unnamed entries) to MessageExampleChanges, or nil if no
// keyed examples changed.
func CompareMessageExampleSlices(l, r []low.ValueReference[*lowasync.MessageExample],
	changes *[]*Change, component, property string,
	configs ...*BreakingRulesConfig,
) map[string]*MessageExampleChanges {
	config := comparisonConfig(configs)
	if len(l) == 0 && len(r) == 0 {
		return nil
	}
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	key := func(v low.ValueReference[*lowasync.MessageExample]) string {
		if v.Value.Name.Value != "" {
			return v.Value.Name.Value
		}
		return low.GenerateHashString(v.Value)
	}

	lExamples, lOrder := groupMessageExamples(l, key)
	rExamples, rOrder := groupMessageExamples(r, key)

	exampleChanges := make(map[string]*MessageExampleChanges)
	for _, k := range lOrder {
		left, right := lExamples[k], rExamples[k]
		common := min(len(left), len(right))
		for i := 0; i < common; i++ {
			if ec := CompareMessageExamples(left[i].Value, right[i].Value, config); ec != nil {
				exampleChanges[messageExampleOccurrenceKey(k, i, len(left), len(right), exampleChanges)] = ec
			}
		}
		for i := common; i < len(left); i++ {
			lv := left[i]
			CreateChange(changes, ObjectRemoved, lowasync.ExamplesLabel,
				lv.ValueNode, nil, breakingRemoved, lv.Value, nil)
		}
		for i := common; i < len(right); i++ {
			rv := right[i]
			CreateChange(changes, ObjectAdded, lowasync.ExamplesLabel,
				nil, rv.ValueNode, breakingAdded, nil, rv.Value)
		}
	}
	for _, k := range rOrder {
		if _, exists := lExamples[k]; exists {
			continue
		}
		for _, rv := range rExamples[k] {
			CreateChange(changes, ObjectAdded, lowasync.ExamplesLabel,
				nil, rv.ValueNode, breakingAdded, nil, rv.Value)
		}
	}
	if len(exampleChanges) == 0 {
		return nil
	}
	return exampleChanges
}

func groupMessageExamples(examples []low.ValueReference[*lowasync.MessageExample],
	key func(low.ValueReference[*lowasync.MessageExample]) string,
) (map[string][]low.ValueReference[*lowasync.MessageExample], []string) {
	grouped := make(map[string][]low.ValueReference[*lowasync.MessageExample], len(examples))
	order := make([]string, 0, len(examples))
	for _, example := range examples {
		if example.Value == nil {
			continue
		}
		k := key(example)
		if _, exists := grouped[k]; !exists {
			order = append(order, k)
		}
		grouped[k] = append(grouped[k], example)
	}
	return grouped, order
}

func messageExampleOccurrenceKey(key string, occurrence, leftCount, rightCount int,
	existing map[string]*MessageExampleChanges,
) string {
	if leftCount == 1 && rightCount == 1 {
		if _, used := existing[key]; !used {
			return key
		}
	}
	for suffix := occurrence + 1; ; suffix++ {
		candidate := fmt.Sprintf("%s#%d", key, suffix)
		if _, used := existing[candidate]; !used {
			return candidate
		}
	}
}

// CompareMessageExamples compares two AsyncAPI Message Example objects and returns a
// pointer to MessageExampleChanges, or nil if nothing changed.
//
// Headers and payload are free-form yaml nodes, so they are compared structurally
// as raw nodes.
func CompareMessageExamples(l, r *lowasync.MessageExample, configs ...*BreakingRulesConfig) *MessageExampleChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMessageExample, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessageExample, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	compareRawNode(l.Headers, r.Headers, lowasync.HeadersLabel,
		CompMessageExample, PropHeaders, &changes, config)
	compareRawNode(l.Payload, r.Payload, lowasync.PayloadLabel,
		CompMessageExample, PropPayload, &changes, config)

	m := new(MessageExampleChanges)
	m.ExtensionChanges = CheckExtensions(l, r)
	m.PropertyChanges = NewPropertyChanges(changes)
	if m.TotalChanges() <= 0 {
		return nil
	}
	return m
}
