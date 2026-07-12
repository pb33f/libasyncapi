// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// ServerVariableChanges represents changes made to a single AsyncAPI ServerVariable object.
type ServerVariableChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between ServerVariable objects.
func (sv *ServerVariableChanges) GetAllChanges() []*Change {
	if sv == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, sv.Changes...)
	if sv.ExtensionChanges != nil {
		changes = append(changes, sv.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (sv *ServerVariableChanges) TotalChanges() int {
	if sv == nil {
		return 0
	}
	c := sv.PropertyChanges.TotalChanges()
	if sv.ExtensionChanges != nil {
		c += sv.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (sv *ServerVariableChanges) TotalBreakingChanges() int {
	if sv == nil {
		return 0
	}
	c := sv.PropertyChanges.TotalBreakingChanges()
	if sv.ExtensionChanges != nil {
		c += sv.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareServerVariables compares two AsyncAPI ServerVariable objects and returns a pointer
// to ServerVariableChanges, or nil if nothing changed.
func CompareServerVariables(l, r *lowasync.ServerVariable, configs ...*BreakingRulesConfig) *ServerVariableChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompServerVariable, PropDefault,
			l.Default.ValueNode, r.Default.ValueNode, lowasync.DefaultLabel, &changes, l, r, config),
		NewPropertyCheck(CompServerVariable, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	if len(l.Enum.Value) > 0 || len(r.Enum.Value) > 0 {
		ExtractStringValueSliceChangesWithRules(l.Enum.Value, r.Enum.Value, &changes,
			lowasync.EnumLabel, CompServerVariable, PropEnum, config)
	}
	if len(l.Examples.Value) > 0 || len(r.Examples.Value) > 0 {
		ExtractStringValueSliceChangesWithRules(l.Examples.Value, r.Examples.Value, &changes,
			lowasync.ExamplesLabel, CompServerVariable, PropExamples, config)
	}

	sv := new(ServerVariableChanges)
	sv.ExtensionChanges = CheckExtensions(l, r)
	sv.PropertyChanges = NewPropertyChanges(changes)
	if sv.TotalChanges() <= 0 {
		return nil
	}
	return sv
}
