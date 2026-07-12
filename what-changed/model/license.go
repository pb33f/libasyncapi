// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
)

// LicenseChanges represents changes made to the License object of an AsyncAPI document.
type LicenseChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between License objects.
func (lc *LicenseChanges) GetAllChanges() []*Change {
	if lc == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, lc.Changes...)
	if lc.ExtensionChanges != nil {
		changes = append(changes, lc.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (lc *LicenseChanges) TotalChanges() int {
	if lc == nil {
		return 0
	}
	t := lc.PropertyChanges.TotalChanges()
	if lc.ExtensionChanges != nil {
		t += lc.ExtensionChanges.TotalChanges()
	}
	return t
}

// TotalBreakingChanges returns the number of breaking changes made.
func (lc *LicenseChanges) TotalBreakingChanges() int {
	if lc == nil {
		return 0
	}
	t := lc.PropertyChanges.TotalBreakingChanges()
	if lc.ExtensionChanges != nil {
		t += lc.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// CompareLicense compares two License objects and returns a pointer to LicenseChanges,
// or nil if nothing changed.
//
// License is a libopenapi base object, but the comparison is implemented locally so
// breaking classification resolves against this package's AsyncAPI breaking-rules
// configuration rather than libopenapi's global OpenAPI configuration.
func CompareLicense(l, r *base.License, configs ...*BreakingRulesConfig) *LicenseChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompLicense, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompLicense, PropIdentifier,
			l.Identifier.ValueNode, r.Identifier.ValueNode, PropIdentifier, &changes, l, r, config),
		NewPropertyCheck(CompLicense, PropURL,
			l.URL.ValueNode, r.URL.ValueNode, lowasync.URLLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	lc := new(LicenseChanges)
	lc.ExtensionChanges = CheckExtensions(l, r)
	lc.PropertyChanges = NewPropertyChanges(changes)
	if lc.TotalChanges() <= 0 {
		return nil
	}
	return lc
}
