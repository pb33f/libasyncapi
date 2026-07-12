// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// MessageTraitChanges represents changes made to an AsyncAPI Message Trait object.
type MessageTraitChanges struct {
	*PropertyChanges
	HeadersChanges       *wcmodel.SchemaChanges            `json:"headers,omitempty" yaml:"headers,omitempty"`
	CorrelationIDChanges *CorrelationIDChanges             `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	TagChanges           map[string]*TagChanges            `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges   *ExternalDocChanges               `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges      *MessageBindingsChanges           `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ExampleChanges       map[string]*MessageExampleChanges `json:"examples,omitempty" yaml:"examples,omitempty"`
	ExtensionChanges     *ExtensionChanges                 `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Message Trait objects.
func (m *MessageTraitChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.HeadersChanges != nil {
		changes = append(changes, m.HeadersChanges.GetAllChanges()...)
	}
	if m.CorrelationIDChanges != nil {
		changes = append(changes, m.CorrelationIDChanges.GetAllChanges()...)
	}
	for k := range m.TagChanges {
		changes = append(changes, m.TagChanges[k].GetAllChanges()...)
	}
	if m.ExternalDocChanges != nil {
		changes = append(changes, m.ExternalDocChanges.GetAllChanges()...)
	}
	if m.BindingsChanges != nil {
		changes = append(changes, m.BindingsChanges.GetAllChanges()...)
	}
	for k := range m.ExampleChanges {
		changes = append(changes, m.ExampleChanges[k].GetAllChanges()...)
	}
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MessageTraitChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.HeadersChanges != nil {
		c += m.HeadersChanges.TotalChanges()
	}
	if m.CorrelationIDChanges != nil {
		c += m.CorrelationIDChanges.TotalChanges()
	}
	for k := range m.TagChanges {
		c += m.TagChanges[k].TotalChanges()
	}
	if m.ExternalDocChanges != nil {
		c += m.ExternalDocChanges.TotalChanges()
	}
	if m.BindingsChanges != nil {
		c += m.BindingsChanges.TotalChanges()
	}
	for k := range m.ExampleChanges {
		c += m.ExampleChanges[k].TotalChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MessageTraitChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.HeadersChanges != nil {
		c += m.HeadersChanges.TotalBreakingChanges()
	}
	if m.CorrelationIDChanges != nil {
		c += m.CorrelationIDChanges.TotalBreakingChanges()
	}
	for k := range m.TagChanges {
		c += m.TagChanges[k].TotalBreakingChanges()
	}
	if m.ExternalDocChanges != nil {
		c += m.ExternalDocChanges.TotalBreakingChanges()
	}
	if m.BindingsChanges != nil {
		c += m.BindingsChanges.TotalBreakingChanges()
	}
	for k := range m.ExampleChanges {
		c += m.ExampleChanges[k].TotalBreakingChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMessageTraits compares two AsyncAPI Message Trait objects and returns a pointer
// to MessageTraitChanges, or nil if nothing changed.
//
// The headers field is a schema, so its comparison is delegated to libopenapi's
// what-changed schema comparator. Examples are keyed by name; unnamed examples fall back
// to hash identity, so an edited unnamed example reports as one removal plus one addition.
func CompareMessageTraits(l, r *lowasync.MessageTrait, configs ...*BreakingRulesConfig) *MessageTraitChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMessageTrait, PropContentType,
			l.ContentType.ValueNode, r.ContentType.ValueNode, lowasync.ContentTypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessageTrait, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessageTrait, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessageTrait, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessageTrait, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	m := new(MessageTraitChanges)

	compareNestedObject(l.Headers, r.Headers, lowasync.HeadersLabel,
		CompMessageTrait, PropHeaders, &changes, wcmodel.CompareSchemas, &m.HeadersChanges, config)

	compareNestedObject(l.CorrelationID, r.CorrelationID, lowasync.CorrelationIDLabel,
		CompMessageTrait, PropCorrelationID, &changes, configuredNestedCompare(config, CompareCorrelationID), &m.CorrelationIDChanges, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes,
		CompMessageTrait, PropTags, config); tc != nil {
		m.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompMessageTrait, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &m.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompMessageTrait, PropBindings, &changes, configuredNestedCompare(config, CompareMessageBindings), &m.BindingsChanges, config)

	if ec := CompareMessageExampleSlices(l.Examples.Value, r.Examples.Value, &changes,
		CompMessageTrait, PropExamples, config); ec != nil {
		m.ExampleChanges = ec
	}

	m.ExtensionChanges = CheckExtensions(l, r)
	m.PropertyChanges = NewPropertyChanges(changes)
	if m.TotalChanges() <= 0 {
		return nil
	}
	return m
}

// MessageChanges represents changes made to an AsyncAPI Message object.
type MessageChanges struct {
	*PropertyChanges
	HeadersChanges       *wcmodel.SchemaChanges            `json:"headers,omitempty" yaml:"headers,omitempty"`
	PayloadChanges       *wcmodel.SchemaChanges            `json:"payload,omitempty" yaml:"payload,omitempty"`
	CorrelationIDChanges *CorrelationIDChanges             `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	TagChanges           map[string]*TagChanges            `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges   *ExternalDocChanges               `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges      *MessageBindingsChanges           `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ExampleChanges       map[string]*MessageExampleChanges `json:"examples,omitempty" yaml:"examples,omitempty"`
	ExtensionChanges     *ExtensionChanges                 `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Message objects.
func (m *MessageChanges) GetAllChanges() []*Change {
	if m == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, m.Changes...)
	if m.HeadersChanges != nil {
		changes = append(changes, m.HeadersChanges.GetAllChanges()...)
	}
	if m.PayloadChanges != nil {
		changes = append(changes, m.PayloadChanges.GetAllChanges()...)
	}
	if m.CorrelationIDChanges != nil {
		changes = append(changes, m.CorrelationIDChanges.GetAllChanges()...)
	}
	for k := range m.TagChanges {
		changes = append(changes, m.TagChanges[k].GetAllChanges()...)
	}
	if m.ExternalDocChanges != nil {
		changes = append(changes, m.ExternalDocChanges.GetAllChanges()...)
	}
	if m.BindingsChanges != nil {
		changes = append(changes, m.BindingsChanges.GetAllChanges()...)
	}
	for k := range m.ExampleChanges {
		changes = append(changes, m.ExampleChanges[k].GetAllChanges()...)
	}
	if m.ExtensionChanges != nil {
		changes = append(changes, m.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (m *MessageChanges) TotalChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalChanges()
	if m.HeadersChanges != nil {
		c += m.HeadersChanges.TotalChanges()
	}
	if m.PayloadChanges != nil {
		c += m.PayloadChanges.TotalChanges()
	}
	if m.CorrelationIDChanges != nil {
		c += m.CorrelationIDChanges.TotalChanges()
	}
	for k := range m.TagChanges {
		c += m.TagChanges[k].TotalChanges()
	}
	if m.ExternalDocChanges != nil {
		c += m.ExternalDocChanges.TotalChanges()
	}
	if m.BindingsChanges != nil {
		c += m.BindingsChanges.TotalChanges()
	}
	for k := range m.ExampleChanges {
		c += m.ExampleChanges[k].TotalChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (m *MessageChanges) TotalBreakingChanges() int {
	if m == nil {
		return 0
	}
	c := m.PropertyChanges.TotalBreakingChanges()
	if m.HeadersChanges != nil {
		c += m.HeadersChanges.TotalBreakingChanges()
	}
	if m.PayloadChanges != nil {
		c += m.PayloadChanges.TotalBreakingChanges()
	}
	if m.CorrelationIDChanges != nil {
		c += m.CorrelationIDChanges.TotalBreakingChanges()
	}
	for k := range m.TagChanges {
		c += m.TagChanges[k].TotalBreakingChanges()
	}
	if m.ExternalDocChanges != nil {
		c += m.ExternalDocChanges.TotalBreakingChanges()
	}
	if m.BindingsChanges != nil {
		c += m.BindingsChanges.TotalBreakingChanges()
	}
	for k := range m.ExampleChanges {
		c += m.ExampleChanges[k].TotalBreakingChanges()
	}
	if m.ExtensionChanges != nil {
		c += m.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareMessages compares two AsyncAPI Message objects and returns a pointer to
// MessageChanges, or nil if nothing changed.
//
// Headers and payload are schemas, so their comparison is delegated to libopenapi's
// what-changed schema comparator. Examples are keyed by name (hash identity for unnamed
// entries). Traits have no natural identity key, so an edited trait reports as one
// removal plus one addition.
func CompareMessages(l, r *lowasync.Message, configs ...*BreakingRulesConfig) *MessageChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompMessage, PropContentType,
			l.ContentType.ValueNode, r.ContentType.ValueNode, lowasync.ContentTypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessage, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessage, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessage, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompMessage, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	m := new(MessageChanges)

	compareNestedObject(l.Headers, r.Headers, lowasync.HeadersLabel,
		CompMessage, PropHeaders, &changes, wcmodel.CompareSchemas, &m.HeadersChanges, config)

	compareNestedObject(l.Payload, r.Payload, lowasync.PayloadLabel,
		CompMessage, PropPayload, &changes, wcmodel.CompareSchemas, &m.PayloadChanges, config)

	compareNestedObject(l.CorrelationID, r.CorrelationID, lowasync.CorrelationIDLabel,
		CompMessage, PropCorrelationID, &changes, configuredNestedCompare(config, CompareCorrelationID), &m.CorrelationIDChanges, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes,
		CompMessage, PropTags, config); tc != nil {
		m.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompMessage, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &m.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompMessage, PropBindings, &changes, configuredNestedCompare(config, CompareMessageBindings), &m.BindingsChanges, config)

	if ec := CompareMessageExampleSlices(l.Examples.Value, r.Examples.Value, &changes,
		CompMessage, PropExamples, config); ec != nil {
		m.ExampleChanges = ec
	}

	compareUnkeyedSlices(l.Traits.Value, r.Traits.Value, lowasync.TraitsLabel,
		CompMessage, PropTraits, &changes, config)

	m.ExtensionChanges = CheckExtensions(l, r)
	m.PropertyChanges = NewPropertyChanges(changes)
	if m.TotalChanges() <= 0 {
		return nil
	}
	return m
}
