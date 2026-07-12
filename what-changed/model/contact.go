// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
)

// ContactChanges represents changes made to the Contact object of an AsyncAPI document.
type ContactChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Contact objects.
func (c *ContactChanges) GetAllChanges() []*Change {
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
func (c *ContactChanges) TotalChanges() int {
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
func (c *ContactChanges) TotalBreakingChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalBreakingChanges()
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// CompareContact compares two Contact objects and returns a pointer to ContactChanges,
// or nil if nothing changed.
//
// Contact is a libopenapi base object, but the comparison is implemented locally so
// breaking classification resolves against this package's AsyncAPI breaking-rules
// configuration rather than libopenapi's global OpenAPI configuration.
func CompareContact(l, r *base.Contact, configs ...*BreakingRulesConfig) *ContactChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompContact, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompContact, PropURL,
			l.URL.ValueNode, r.URL.ValueNode, lowasync.URLLabel, &changes, l, r, config),
		NewPropertyCheck(CompContact, PropEmail,
			l.Email.ValueNode, r.Email.ValueNode, lowasync.EmailLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	cc := new(ContactChanges)
	cc.ExtensionChanges = CheckExtensions(l, r)
	cc.PropertyChanges = NewPropertyChanges(changes)
	if cc.TotalChanges() <= 0 {
		return nil
	}
	return cc
}
