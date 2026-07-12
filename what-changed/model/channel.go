// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// ParameterChanges represents changes made to a single AsyncAPI Parameter object.
type ParameterChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Parameter objects.
func (p *ParameterChanges) GetAllChanges() []*Change {
	if p == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, p.Changes...)
	if p.ExtensionChanges != nil {
		changes = append(changes, p.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (p *ParameterChanges) TotalChanges() int {
	if p == nil {
		return 0
	}
	c := p.PropertyChanges.TotalChanges()
	if p.ExtensionChanges != nil {
		c += p.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (p *ParameterChanges) TotalBreakingChanges() int {
	if p == nil {
		return 0
	}
	c := p.PropertyChanges.TotalBreakingChanges()
	if p.ExtensionChanges != nil {
		c += p.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareParameters compares two AsyncAPI Parameter objects and returns a pointer to
// ParameterChanges, or nil if nothing changed.
func CompareParameters(l, r *lowasync.Parameter, configs ...*BreakingRulesConfig) *ParameterChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompParameter, PropDefault,
			l.Default.ValueNode, r.Default.ValueNode, lowasync.DefaultLabel, &changes, l, r, config),
		NewPropertyCheck(CompParameter, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
		NewPropertyCheck(CompParameter, PropLocation,
			l.Location.ValueNode, r.Location.ValueNode, lowasync.LocationLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	if len(l.Enum.Value) > 0 || len(r.Enum.Value) > 0 {
		ExtractStringValueSliceChangesWithRules(l.Enum.Value, r.Enum.Value, &changes,
			lowasync.EnumLabel, CompParameter, PropEnum, config)
	}
	if len(l.Examples.Value) > 0 || len(r.Examples.Value) > 0 {
		ExtractStringValueSliceChangesWithRules(l.Examples.Value, r.Examples.Value, &changes,
			lowasync.ExamplesLabel, CompParameter, PropExamples, config)
	}

	pc := new(ParameterChanges)
	pc.ExtensionChanges = CheckExtensions(l, r)
	pc.PropertyChanges = NewPropertyChanges(changes)
	if pc.TotalChanges() <= 0 {
		return nil
	}
	return pc
}

// ChannelChanges represents changes made to a single AsyncAPI Channel object.
type ChannelChanges struct {
	*PropertyChanges
	MessageChanges     map[string]*MessageChanges   `json:"messages,omitempty" yaml:"messages,omitempty"`
	ParameterChanges   map[string]*ParameterChanges `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	TagChanges         map[string]*TagChanges       `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges *ExternalDocChanges          `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges    *ChannelBindingsChanges      `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ExtensionChanges   *ExtensionChanges            `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Channel objects.
func (c *ChannelChanges) GetAllChanges() []*Change {
	if c == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, c.Changes...)
	for k := range c.MessageChanges {
		changes = append(changes, c.MessageChanges[k].GetAllChanges()...)
	}
	for k := range c.ParameterChanges {
		changes = append(changes, c.ParameterChanges[k].GetAllChanges()...)
	}
	for k := range c.TagChanges {
		changes = append(changes, c.TagChanges[k].GetAllChanges()...)
	}
	if c.ExternalDocChanges != nil {
		changes = append(changes, c.ExternalDocChanges.GetAllChanges()...)
	}
	if c.BindingsChanges != nil {
		changes = append(changes, c.BindingsChanges.GetAllChanges()...)
	}
	if c.ExtensionChanges != nil {
		changes = append(changes, c.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (c *ChannelChanges) TotalChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalChanges()
	for k := range c.MessageChanges {
		t += c.MessageChanges[k].TotalChanges()
	}
	for k := range c.ParameterChanges {
		t += c.ParameterChanges[k].TotalChanges()
	}
	for k := range c.TagChanges {
		t += c.TagChanges[k].TotalChanges()
	}
	if c.ExternalDocChanges != nil {
		t += c.ExternalDocChanges.TotalChanges()
	}
	if c.BindingsChanges != nil {
		t += c.BindingsChanges.TotalChanges()
	}
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalChanges()
	}
	return t
}

// TotalBreakingChanges returns the number of breaking changes made.
func (c *ChannelChanges) TotalBreakingChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalBreakingChanges()
	for k := range c.MessageChanges {
		t += c.MessageChanges[k].TotalBreakingChanges()
	}
	for k := range c.ParameterChanges {
		t += c.ParameterChanges[k].TotalBreakingChanges()
	}
	for k := range c.TagChanges {
		t += c.TagChanges[k].TotalBreakingChanges()
	}
	if c.ExternalDocChanges != nil {
		t += c.ExternalDocChanges.TotalBreakingChanges()
	}
	if c.BindingsChanges != nil {
		t += c.BindingsChanges.TotalBreakingChanges()
	}
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// CompareChannels compares two AsyncAPI Channel objects and returns a pointer to
// ChannelChanges, or nil if nothing changed.
//
// The channel address is nullable in AsyncAPI 3.0, so addition, removal and modification
// of the address are checked by hand. Server entries are reference-only and compared as
// sets of reference strings, never recursed into.
func CompareChannels(l, r *lowasync.Channel, configs ...*BreakingRulesConfig) *ChannelChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change

	// address is nullable (*string) and has three states: absent (no key), an explicit
	// `address: null` (address unknown) and a concrete value. all state transitions are
	// reported; null<->value transitions are modifications because the property itself
	// was already present.
	lVal, rVal := l.Address.Value, r.Address.Value
	lNull := lVal == nil && l.Address.ValueNode != nil
	rNull := rVal == nil && r.Address.ValueNode != nil
	switch {
	case lVal != nil && rVal != nil:
		if *lVal != *rVal {
			CreateChange(&changes, Modified, lowasync.AddressLabel,
				l.Address.ValueNode, r.Address.ValueNode,
				BreakingModified(CompChannel, PropAddress, config), lVal, rVal)
		}
	case lVal == nil && rVal != nil:
		if lNull {
			CreateChange(&changes, Modified, lowasync.AddressLabel,
				l.Address.ValueNode, r.Address.ValueNode,
				BreakingModified(CompChannel, PropAddress, config), nil, rVal)
		} else {
			CreateChange(&changes, PropertyAdded, lowasync.AddressLabel,
				nil, r.Address.ValueNode, BreakingAdded(CompChannel, PropAddress, config), nil, rVal)
		}
	case lVal != nil && rVal == nil:
		if rNull {
			CreateChange(&changes, Modified, lowasync.AddressLabel,
				l.Address.ValueNode, r.Address.ValueNode,
				BreakingModified(CompChannel, PropAddress, config), lVal, nil)
		} else {
			CreateChange(&changes, PropertyRemoved, lowasync.AddressLabel,
				l.Address.ValueNode, nil, BreakingRemoved(CompChannel, PropAddress, config), lVal, nil)
		}
	case !lNull && rNull:
		CreateChange(&changes, PropertyAdded, lowasync.AddressLabel,
			nil, r.Address.ValueNode, BreakingAdded(CompChannel, PropAddress, config), nil, nil)
	case lNull && !rNull:
		CreateChange(&changes, PropertyRemoved, lowasync.AddressLabel,
			l.Address.ValueNode, nil, BreakingRemoved(CompChannel, PropAddress, config), nil, nil)
	}

	props := []*PropertyCheck{
		NewPropertyCheck(CompChannel, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompChannel, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompChannel, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	cc := new(ChannelChanges)

	if mc := CheckMapForChangesWithRules(l.Messages.Value, r.Messages.Value, &changes,
		lowasync.MessagesLabel, configuredCompare(config, CompareMessages), CompChannel, PropMessages, config); len(mc) > 0 {
		cc.MessageChanges = mc
	}

	if pc := CheckMapForChangesWithRules(l.Parameters.Value, r.Parameters.Value, &changes,
		lowasync.ParametersLabel, configuredCompare(config, CompareParameters), CompChannel, PropParameters, config); len(pc) > 0 {
		cc.ParameterChanges = pc
	}

	// servers are reference-only entries, compare them as reference strings.
	compareRefSlices(l.Servers.Value, r.Servers.Value,
		lowasync.ServersLabel, CompChannel, PropServers, &changes, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes, CompChannel, PropTags, config); len(tc) > 0 {
		cc.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompChannel, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &cc.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompChannel, PropBindings, &changes, configuredNestedCompare(config, CompareChannelBindings), &cc.BindingsChanges, config)

	cc.ExtensionChanges = CheckExtensions(l, r)
	cc.PropertyChanges = NewPropertyChanges(changes)
	if cc.TotalChanges() <= 0 {
		return nil
	}
	return cc
}
