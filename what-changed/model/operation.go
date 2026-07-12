// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// OperationReplyAddressChanges represents changes made to a single AsyncAPI
// OperationReplyAddress object.
type OperationReplyAddressChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OperationReplyAddress objects.
func (ora *OperationReplyAddressChanges) GetAllChanges() []*Change {
	if ora == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, ora.Changes...)
	if ora.ExtensionChanges != nil {
		changes = append(changes, ora.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (ora *OperationReplyAddressChanges) TotalChanges() int {
	if ora == nil {
		return 0
	}
	c := ora.PropertyChanges.TotalChanges()
	if ora.ExtensionChanges != nil {
		c += ora.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (ora *OperationReplyAddressChanges) TotalBreakingChanges() int {
	if ora == nil {
		return 0
	}
	c := ora.PropertyChanges.TotalBreakingChanges()
	if ora.ExtensionChanges != nil {
		c += ora.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOperationReplyAddress compares two AsyncAPI OperationReplyAddress objects and
// returns a pointer to OperationReplyAddressChanges, or nil if nothing changed.
func CompareOperationReplyAddress(l, r *lowasync.OperationReplyAddress, configs ...*BreakingRulesConfig) *OperationReplyAddressChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompOperationReplyAddress, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperationReplyAddress, PropLocation,
			l.Location.ValueNode, r.Location.ValueNode, lowasync.LocationLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ora := new(OperationReplyAddressChanges)
	ora.ExtensionChanges = CheckExtensions(l, r)
	ora.PropertyChanges = NewPropertyChanges(changes)
	if ora.TotalChanges() <= 0 {
		return nil
	}
	return ora
}

// OperationReplyChanges represents changes made to a single AsyncAPI OperationReply object.
type OperationReplyChanges struct {
	*PropertyChanges
	AddressChanges   *OperationReplyAddressChanges `json:"address,omitempty" yaml:"address,omitempty"`
	ExtensionChanges *ExtensionChanges             `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OperationReply objects.
func (or *OperationReplyChanges) GetAllChanges() []*Change {
	if or == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, or.Changes...)
	if or.AddressChanges != nil {
		changes = append(changes, or.AddressChanges.GetAllChanges()...)
	}
	if or.ExtensionChanges != nil {
		changes = append(changes, or.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (or *OperationReplyChanges) TotalChanges() int {
	if or == nil {
		return 0
	}
	c := or.PropertyChanges.TotalChanges()
	if or.AddressChanges != nil {
		c += or.AddressChanges.TotalChanges()
	}
	if or.ExtensionChanges != nil {
		c += or.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (or *OperationReplyChanges) TotalBreakingChanges() int {
	if or == nil {
		return 0
	}
	c := or.PropertyChanges.TotalBreakingChanges()
	if or.AddressChanges != nil {
		c += or.AddressChanges.TotalBreakingChanges()
	}
	if or.ExtensionChanges != nil {
		c += or.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOperationReply compares two AsyncAPI OperationReply objects and returns a pointer
// to OperationReplyChanges, or nil if nothing changed.
//
// The reply channel and messages are reference-only fields and are compared as reference
// strings, never recursed into.
func CompareOperationReply(l, r *lowasync.OperationReply, configs ...*BreakingRulesConfig) *OperationReplyChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	or := new(OperationReplyChanges)

	compareNestedObject(l.Address, r.Address, lowasync.AddressLabel,
		CompOperationReply, PropAddress, &changes, configuredNestedCompare(config, CompareOperationReplyAddress), &or.AddressChanges, config)

	compareSingleRef(l.Channel, r.Channel, lowasync.ChannelLabel,
		CompOperationReply, PropChannel, &changes, config)

	compareRefSlices(l.Messages.Value, r.Messages.Value,
		lowasync.MessagesLabel, CompOperationReply, PropMessages, &changes, config)

	or.ExtensionChanges = CheckExtensions(l, r)
	or.PropertyChanges = NewPropertyChanges(changes)
	if or.TotalChanges() <= 0 {
		return nil
	}
	return or
}

// OperationTraitChanges represents changes made to a single AsyncAPI OperationTrait object.
type OperationTraitChanges struct {
	*PropertyChanges
	TagChanges         map[string]*TagChanges    `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges *ExternalDocChanges       `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges    *OperationBindingsChanges `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ExtensionChanges   *ExtensionChanges         `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OperationTrait objects.
func (ot *OperationTraitChanges) GetAllChanges() []*Change {
	if ot == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, ot.Changes...)
	for k := range ot.TagChanges {
		changes = append(changes, ot.TagChanges[k].GetAllChanges()...)
	}
	if ot.ExternalDocChanges != nil {
		changes = append(changes, ot.ExternalDocChanges.GetAllChanges()...)
	}
	if ot.BindingsChanges != nil {
		changes = append(changes, ot.BindingsChanges.GetAllChanges()...)
	}
	if ot.ExtensionChanges != nil {
		changes = append(changes, ot.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (ot *OperationTraitChanges) TotalChanges() int {
	if ot == nil {
		return 0
	}
	c := ot.PropertyChanges.TotalChanges()
	for k := range ot.TagChanges {
		c += ot.TagChanges[k].TotalChanges()
	}
	if ot.ExternalDocChanges != nil {
		c += ot.ExternalDocChanges.TotalChanges()
	}
	if ot.BindingsChanges != nil {
		c += ot.BindingsChanges.TotalChanges()
	}
	if ot.ExtensionChanges != nil {
		c += ot.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (ot *OperationTraitChanges) TotalBreakingChanges() int {
	if ot == nil {
		return 0
	}
	c := ot.PropertyChanges.TotalBreakingChanges()
	for k := range ot.TagChanges {
		c += ot.TagChanges[k].TotalBreakingChanges()
	}
	if ot.ExternalDocChanges != nil {
		c += ot.ExternalDocChanges.TotalBreakingChanges()
	}
	if ot.BindingsChanges != nil {
		c += ot.BindingsChanges.TotalBreakingChanges()
	}
	if ot.ExtensionChanges != nil {
		c += ot.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOperationTraits compares two AsyncAPI OperationTrait objects and returns a pointer
// to OperationTraitChanges, or nil if nothing changed.
func CompareOperationTraits(l, r *lowasync.OperationTrait, configs ...*BreakingRulesConfig) *OperationTraitChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompOperationTrait, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperationTrait, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperationTrait, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	ot := new(OperationTraitChanges)

	compareUnkeyedSlices(l.Security.Value, r.Security.Value,
		lowasync.SecurityLabel, CompOperationTrait, PropSecurity, &changes, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes, CompOperationTrait, PropTags, config); len(tc) > 0 {
		ot.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompOperationTrait, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &ot.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompOperationTrait, PropBindings, &changes, configuredNestedCompare(config, CompareOperationBindings), &ot.BindingsChanges, config)

	ot.ExtensionChanges = CheckExtensions(l, r)
	ot.PropertyChanges = NewPropertyChanges(changes)
	if ot.TotalChanges() <= 0 {
		return nil
	}
	return ot
}

// OperationChanges represents changes made to a single AsyncAPI Operation object.
type OperationChanges struct {
	*PropertyChanges
	TagChanges         map[string]*TagChanges    `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges *ExternalDocChanges       `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges    *OperationBindingsChanges `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ReplyChanges       *OperationReplyChanges    `json:"reply,omitempty" yaml:"reply,omitempty"`
	ExtensionChanges   *ExtensionChanges         `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Operation objects.
func (o *OperationChanges) GetAllChanges() []*Change {
	if o == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, o.Changes...)
	for k := range o.TagChanges {
		changes = append(changes, o.TagChanges[k].GetAllChanges()...)
	}
	if o.ExternalDocChanges != nil {
		changes = append(changes, o.ExternalDocChanges.GetAllChanges()...)
	}
	if o.BindingsChanges != nil {
		changes = append(changes, o.BindingsChanges.GetAllChanges()...)
	}
	if o.ReplyChanges != nil {
		changes = append(changes, o.ReplyChanges.GetAllChanges()...)
	}
	if o.ExtensionChanges != nil {
		changes = append(changes, o.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (o *OperationChanges) TotalChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalChanges()
	for k := range o.TagChanges {
		c += o.TagChanges[k].TotalChanges()
	}
	if o.ExternalDocChanges != nil {
		c += o.ExternalDocChanges.TotalChanges()
	}
	if o.BindingsChanges != nil {
		c += o.BindingsChanges.TotalChanges()
	}
	if o.ReplyChanges != nil {
		c += o.ReplyChanges.TotalChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (o *OperationChanges) TotalBreakingChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalBreakingChanges()
	for k := range o.TagChanges {
		c += o.TagChanges[k].TotalBreakingChanges()
	}
	if o.ExternalDocChanges != nil {
		c += o.ExternalDocChanges.TotalBreakingChanges()
	}
	if o.BindingsChanges != nil {
		c += o.BindingsChanges.TotalBreakingChanges()
	}
	if o.ReplyChanges != nil {
		c += o.ReplyChanges.TotalBreakingChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOperations compares two AsyncAPI Operation objects and returns a pointer to
// OperationChanges, or nil if nothing changed.
//
// The operation channel and messages are reference-only fields and are compared as
// reference strings, never recursed into. Traits have no natural identity key so they are
// compared as hashed sets; an in-place trait edit reports as one removal and one addition.
func CompareOperations(l, r *lowasync.Operation, configs ...*BreakingRulesConfig) *OperationChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompOperation, PropAction,
			l.Action.ValueNode, r.Action.ValueNode, lowasync.ActionLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperation, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperation, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompOperation, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	oc := new(OperationChanges)

	compareSingleRef(l.Channel, r.Channel, lowasync.ChannelLabel,
		CompOperation, PropChannel, &changes, config)

	compareUnkeyedSlices(l.Security.Value, r.Security.Value,
		lowasync.SecurityLabel, CompOperation, PropSecurity, &changes, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes, CompOperation, PropTags, config); len(tc) > 0 {
		oc.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompOperation, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &oc.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompOperation, PropBindings, &changes, configuredNestedCompare(config, CompareOperationBindings), &oc.BindingsChanges, config)

	compareUnkeyedSlices(l.Traits.Value, r.Traits.Value,
		lowasync.TraitsLabel, CompOperation, PropTraits, &changes, config)

	compareRefSlices(l.Messages.Value, r.Messages.Value,
		lowasync.MessagesLabel, CompOperation, PropMessages, &changes, config)

	compareNestedObject(l.Reply, r.Reply, lowasync.ReplyLabel,
		CompOperation, PropReply, &changes, configuredNestedCompare(config, CompareOperationReply), &oc.ReplyChanges, config)

	oc.ExtensionChanges = CheckExtensions(l, r)
	oc.PropertyChanges = NewPropertyChanges(changes)
	if oc.TotalChanges() <= 0 {
		return nil
	}
	return oc
}
