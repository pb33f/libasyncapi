// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// CorrelationIDChanges represents changes made to an AsyncAPI Correlation ID object.
type CorrelationIDChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Correlation ID objects.
func (c *CorrelationIDChanges) GetAllChanges() []*Change {
	if c == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, c.Changes...)
	if c.ExtensionChanges != nil {
		changes = append(changes, c.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (c *CorrelationIDChanges) TotalChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalChanges()
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalChanges()
	}
	return t
}

// TotalBreakingChanges returns the number of breaking changes made.
func (c *CorrelationIDChanges) TotalBreakingChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalBreakingChanges()
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// CompareCorrelationID compares two AsyncAPI Correlation ID objects and returns a pointer
// to CorrelationIDChanges, or nil if nothing changed.
func CompareCorrelationID(l, r *lowasync.CorrelationID, configs ...*BreakingRulesConfig) *CorrelationIDChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompCorrelationID, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
		NewPropertyCheck(CompCorrelationID, PropLocation,
			l.Location.ValueNode, r.Location.ValueNode, lowasync.LocationLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	c := new(CorrelationIDChanges)
	c.ExtensionChanges = CheckExtensions(l, r)
	c.PropertyChanges = NewPropertyChanges(changes)
	if c.TotalChanges() <= 0 {
		return nil
	}
	return c
}
