// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// ServerChanges represents changes made to a single AsyncAPI Server object.
type ServerChanges struct {
	*PropertyChanges
	ServerVariableChanges map[string]*ServerVariableChanges `json:"variables,omitempty" yaml:"variables,omitempty"`
	TagChanges            map[string]*TagChanges            `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocChanges    *ExternalDocChanges               `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	BindingsChanges       *ServerBindingsChanges            `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	ExtensionChanges      *ExtensionChanges                 `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Server objects.
func (s *ServerChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	for k := range s.ServerVariableChanges {
		changes = append(changes, s.ServerVariableChanges[k].GetAllChanges()...)
	}
	for k := range s.TagChanges {
		changes = append(changes, s.TagChanges[k].GetAllChanges()...)
	}
	if s.ExternalDocChanges != nil {
		changes = append(changes, s.ExternalDocChanges.GetAllChanges()...)
	}
	if s.BindingsChanges != nil {
		changes = append(changes, s.BindingsChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *ServerChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	for k := range s.ServerVariableChanges {
		c += s.ServerVariableChanges[k].TotalChanges()
	}
	for k := range s.TagChanges {
		c += s.TagChanges[k].TotalChanges()
	}
	if s.ExternalDocChanges != nil {
		c += s.ExternalDocChanges.TotalChanges()
	}
	if s.BindingsChanges != nil {
		c += s.BindingsChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *ServerChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	for k := range s.ServerVariableChanges {
		c += s.ServerVariableChanges[k].TotalBreakingChanges()
	}
	for k := range s.TagChanges {
		c += s.TagChanges[k].TotalBreakingChanges()
	}
	if s.ExternalDocChanges != nil {
		c += s.ExternalDocChanges.TotalBreakingChanges()
	}
	if s.BindingsChanges != nil {
		c += s.BindingsChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareServers compares two AsyncAPI Server objects and returns a pointer to ServerChanges,
// or nil if nothing changed.
func CompareServers(l, r *lowasync.Server, configs ...*BreakingRulesConfig) *ServerChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompServer, PropHost,
			l.Host.ValueNode, r.Host.ValueNode, lowasync.HostLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropProtocol,
			l.Protocol.ValueNode, r.Protocol.ValueNode, lowasync.ProtocolLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropProtocolVersion,
			l.ProtocolVersion.ValueNode, r.ProtocolVersion.ValueNode, lowasync.ProtocolVersionLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropPathname,
			l.Pathname.ValueNode, r.Pathname.ValueNode, lowasync.PathnameLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropTitle,
			l.Title.ValueNode, r.Title.ValueNode, lowasync.TitleLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropSummary,
			l.Summary.ValueNode, r.Summary.ValueNode, lowasync.SummaryLabel, &changes, l, r, config),
		NewPropertyCheck(CompServer, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	sc := new(ServerChanges)

	if vc := CheckMapForChangesWithRules(l.Variables.Value, r.Variables.Value, &changes,
		lowasync.VariablesLabel, configuredCompare(config, CompareServerVariables), CompServer, PropVariables, config); len(vc) > 0 {
		sc.ServerVariableChanges = vc
	}

	compareUnkeyedSlices(l.Security.Value, r.Security.Value,
		lowasync.SecurityLabel, CompServer, PropSecurity, &changes, config)

	if tc := CompareTagSlices(l.Tags.Value, r.Tags.Value, &changes, CompServer, PropTags, config); len(tc) > 0 {
		sc.TagChanges = tc
	}

	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompServer, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &sc.ExternalDocChanges, config)

	compareNestedObject(l.Bindings, r.Bindings, lowasync.BindingsLabel,
		CompServer, PropBindings, &changes, configuredNestedCompare(config, CompareServerBindings), &sc.BindingsChanges, config)

	sc.ExtensionChanges = CheckExtensions(l, r)
	sc.PropertyChanges = NewPropertyChanges(changes)
	if sc.TotalChanges() <= 0 {
		return nil
	}
	return sc
}
