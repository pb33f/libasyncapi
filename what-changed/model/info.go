// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// InfoChanges represents changes made to the Info object of an AsyncAPI document.
type InfoChanges struct {
	*PropertyChanges
	ContactChanges     *ContactChanges        `json:"contact,omitempty" yaml:"contact,omitempty"`
	LicenseChanges     *LicenseChanges        `json:"license,omitempty" yaml:"license,omitempty"`
	TagChanges         map[string]*TagChanges `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges *ExternalDocChanges    `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	ExtensionChanges   *ExtensionChanges      `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Info objects.
func (i *InfoChanges) GetAllChanges() []*Change {
	if i == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, i.Changes...)
	if i.ContactChanges != nil {
		changes = append(changes, i.ContactChanges.GetAllChanges()...)
	}
	if i.LicenseChanges != nil {
		changes = append(changes, i.LicenseChanges.GetAllChanges()...)
	}
	for k := range i.TagChanges {
		changes = append(changes, i.TagChanges[k].GetAllChanges()...)
	}
	if i.ExternalDocChanges != nil {
		changes = append(changes, i.ExternalDocChanges.GetAllChanges()...)
	}
	if i.ExtensionChanges != nil {
		changes = append(changes, i.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (i *InfoChanges) TotalChanges() int {
	if i == nil {
		return 0
	}
	c := i.PropertyChanges.TotalChanges()
	if i.ContactChanges != nil {
		c += i.ContactChanges.TotalChanges()
	}
	if i.LicenseChanges != nil {
		c += i.LicenseChanges.TotalChanges()
	}
	for k := range i.TagChanges {
		c += i.TagChanges[k].TotalChanges()
	}
	if i.ExternalDocChanges != nil {
		c += i.ExternalDocChanges.TotalChanges()
	}
	if i.ExtensionChanges != nil {
		c += i.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (i *InfoChanges) TotalBreakingChanges() int {
	if i == nil {
		return 0
	}
	c := i.PropertyChanges.TotalBreakingChanges()
	if i.ContactChanges != nil {
		c += i.ContactChanges.TotalBreakingChanges()
	}
	if i.LicenseChanges != nil {
		c += i.LicenseChanges.TotalBreakingChanges()
	}
	for k := range i.TagChanges {
		c += i.TagChanges[k].TotalBreakingChanges()
	}
	if i.ExternalDocChanges != nil {
		c += i.ExternalDocChanges.TotalBreakingChanges()
	}
	if i.ExtensionChanges != nil {
		c += i.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareInfo compares two AsyncAPI Info objects and returns a pointer to InfoChanges,
// or nil if nothing changed.
//
// Contact and License are libopenapi base objects, but they are compared with this
// package's local comparators so breaking classification follows the AsyncAPI rules.
func CompareInfo(l, r *lowasync.Info, configs ...*BreakingRulesConfig) *InfoChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompInfo, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompInfo, PropVersion,
			l.Version.ValueNode, r.Version.ValueNode, lowasync.VersionLabel, &changes, l, r, config),
		NewPropertyCheck(CompInfo, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
		NewPropertyCheck(CompInfo, PropTermsOfService,
			l.TermsOfService.ValueNode, r.TermsOfService.ValueNode, lowasync.TermsOfServiceLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	i := new(InfoChanges)

	// contact + license are libopenapi base objects, but they are compared locally so
	// breaking classification follows this package's AsyncAPI rules (see CompareContact).
	if l.Contact.Value != nil && r.Contact.Value != nil {
		i.ContactChanges = CompareContact(l.Contact.Value, r.Contact.Value, config)
	} else {
		if l.Contact.Value == nil && r.Contact.Value != nil {
			CreateChange(&changes, ObjectAdded, lowasync.ContactLabel,
				nil, r.Contact.ValueNode, BreakingAdded(CompInfo, PropContact, config), nil, r.Contact.Value)
		}
		if l.Contact.Value != nil && r.Contact.Value == nil {
			CreateChange(&changes, ObjectRemoved, lowasync.ContactLabel,
				l.Contact.ValueNode, nil, BreakingRemoved(CompInfo, PropContact, config), l.Contact.Value, nil)
		}
	}

	if l.License.Value != nil && r.License.Value != nil {
		i.LicenseChanges = CompareLicense(l.License.Value, r.License.Value, config)
	} else {
		if l.License.Value == nil && r.License.Value != nil {
			CreateChange(&changes, ObjectAdded, lowasync.LicenseLabel,
				nil, r.License.ValueNode, BreakingAdded(CompInfo, PropLicense, config), nil, r.License.Value)
		}
		if l.License.Value != nil && r.License.Value == nil {
			CreateChange(&changes, ObjectRemoved, lowasync.LicenseLabel,
				l.License.ValueNode, nil, BreakingRemoved(CompInfo, PropLicense, config), l.License.Value, nil)
		}
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompInfo, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &i.ExternalDocChanges, config)

	if tc := CheckMapForChangesWithRules(l.Tags.Value, r.Tags.Value, &changes,
		lowasync.TagsLabel, configuredCompare(config, CompareTags), CompInfo, PropTags, config); len(tc) > 0 {
		i.TagChanges = tc
	}

	i.ExtensionChanges = CheckExtensions(l, r)
	i.PropertyChanges = NewPropertyChanges(changes)
	if i.TotalChanges() <= 0 {
		return nil
	}
	return i
}
