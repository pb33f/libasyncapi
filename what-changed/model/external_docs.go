// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// ExternalDocChanges represents changes made to a single AsyncAPI ExternalDoc object.
type ExternalDocChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between ExternalDoc objects.
func (e *ExternalDocChanges) GetAllChanges() []*Change {
	if e == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, e.Changes...)
	if e.ExtensionChanges != nil {
		changes = append(changes, e.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (e *ExternalDocChanges) TotalChanges() int {
	if e == nil {
		return 0
	}
	c := e.PropertyChanges.TotalChanges()
	if e.ExtensionChanges != nil {
		c += e.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (e *ExternalDocChanges) TotalBreakingChanges() int {
	if e == nil {
		return 0
	}
	c := e.PropertyChanges.TotalBreakingChanges()
	if e.ExtensionChanges != nil {
		c += e.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareExternalDocs compares two AsyncAPI ExternalDoc objects and returns a pointer to
// ExternalDocChanges, or nil if nothing changed.
func CompareExternalDocs(l, r *lowasync.ExternalDoc, configs ...*BreakingRulesConfig) *ExternalDocChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompExternalDocs, PropURL,
			l.URL.ValueNode, r.URL.ValueNode, lowasync.URLLabel, &changes, l, r, config),
		NewPropertyCheck(CompExternalDocs, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	dc := new(ExternalDocChanges)
	dc.ExtensionChanges = CheckExtensions(l, r)
	dc.PropertyChanges = NewPropertyChanges(changes)
	if dc.TotalChanges() <= 0 {
		return nil
	}
	return dc
}
