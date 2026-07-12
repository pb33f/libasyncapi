// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// OAuthFlowChanges represents changes made to an AsyncAPI OAuth Flow object.
type OAuthFlowChanges struct {
	*PropertyChanges
	ExtensionChanges *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OAuth Flow objects.
func (o *OAuthFlowChanges) GetAllChanges() []*Change {
	if o == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, o.Changes...)
	if o.ExtensionChanges != nil {
		changes = append(changes, o.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (o *OAuthFlowChanges) TotalChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalChanges()
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (o *OAuthFlowChanges) TotalBreakingChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalBreakingChanges()
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOAuthFlow compares two AsyncAPI OAuth Flow objects and returns a pointer to
// OAuthFlowChanges, or nil if nothing changed.
//
// Available scopes are compared key-by-key: removed scopes, added scopes and modified
// scope descriptions are each recorded as individual changes.
func CompareOAuthFlow(l, r *lowasync.OAuthFlow, configs ...*BreakingRulesConfig) *OAuthFlowChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompOAuthFlow, PropAuthorizationURL,
			l.AuthorizationURL.ValueNode, r.AuthorizationURL.ValueNode, lowasync.AuthorizationURLLabel, &changes, l, r, config),
		NewPropertyCheck(CompOAuthFlow, PropTokenURL,
			l.TokenURL.ValueNode, r.TokenURL.ValueNode, lowasync.TokenURLLabel, &changes, l, r, config),
		NewPropertyCheck(CompOAuthFlow, PropRefreshURL,
			l.RefreshURL.ValueNode, r.RefreshURL.ValueNode, lowasync.RefreshURLLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	// compare available scopes key-by-key.
	lScopes := make(map[string]low.ValueReference[string])
	rScopes := make(map[string]low.ValueReference[string])
	var lOrder, rOrder []string
	if l.AvailableScopes.Value != nil {
		for k, v := range l.AvailableScopes.Value.FromOldest() {
			lOrder = append(lOrder, k.Value)
			lScopes[k.Value] = v
		}
	}
	if r.AvailableScopes.Value != nil {
		for k, v := range r.AvailableScopes.Value.FromOldest() {
			rOrder = append(rOrder, k.Value)
			rScopes[k.Value] = v
		}
	}
	for _, k := range lOrder {
		lv := lScopes[k]
		if rv, ok := rScopes[k]; ok {
			if lv.Value != rv.Value {
				CreateChange(&changes, Modified, lowasync.AvailableScopesLabel,
					lv.ValueNode, rv.ValueNode,
					BreakingModified(CompOAuthFlow, PropAvailableScopes, config), lv.Value, rv.Value)
			}
			continue
		}
		CreateChange(&changes, PropertyRemoved, lowasync.AvailableScopesLabel,
			lv.ValueNode, nil,
			BreakingRemoved(CompOAuthFlow, PropAvailableScopes, config), lv.Value, nil)
	}
	for _, k := range rOrder {
		rv := rScopes[k]
		if _, ok := lScopes[k]; !ok {
			CreateChange(&changes, PropertyAdded, lowasync.AvailableScopesLabel,
				nil, rv.ValueNode,
				BreakingAdded(CompOAuthFlow, PropAvailableScopes, config), nil, rv.Value)
		}
	}

	o := new(OAuthFlowChanges)
	o.ExtensionChanges = CheckExtensions(l, r)
	o.PropertyChanges = NewPropertyChanges(changes)
	if o.TotalChanges() <= 0 {
		return nil
	}
	return o
}

// OAuthFlowsChanges represents changes made to an AsyncAPI OAuth Flows object.
type OAuthFlowsChanges struct {
	*PropertyChanges
	ImplicitChanges          *OAuthFlowChanges `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	PasswordChanges          *OAuthFlowChanges `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentialsChanges *OAuthFlowChanges `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCodeChanges *OAuthFlowChanges `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
	ExtensionChanges         *ExtensionChanges `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between OAuth Flows objects.
func (o *OAuthFlowsChanges) GetAllChanges() []*Change {
	if o == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, o.Changes...)
	if o.ImplicitChanges != nil {
		changes = append(changes, o.ImplicitChanges.GetAllChanges()...)
	}
	if o.PasswordChanges != nil {
		changes = append(changes, o.PasswordChanges.GetAllChanges()...)
	}
	if o.ClientCredentialsChanges != nil {
		changes = append(changes, o.ClientCredentialsChanges.GetAllChanges()...)
	}
	if o.AuthorizationCodeChanges != nil {
		changes = append(changes, o.AuthorizationCodeChanges.GetAllChanges()...)
	}
	if o.ExtensionChanges != nil {
		changes = append(changes, o.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (o *OAuthFlowsChanges) TotalChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalChanges()
	if o.ImplicitChanges != nil {
		c += o.ImplicitChanges.TotalChanges()
	}
	if o.PasswordChanges != nil {
		c += o.PasswordChanges.TotalChanges()
	}
	if o.ClientCredentialsChanges != nil {
		c += o.ClientCredentialsChanges.TotalChanges()
	}
	if o.AuthorizationCodeChanges != nil {
		c += o.AuthorizationCodeChanges.TotalChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (o *OAuthFlowsChanges) TotalBreakingChanges() int {
	if o == nil {
		return 0
	}
	c := o.PropertyChanges.TotalBreakingChanges()
	if o.ImplicitChanges != nil {
		c += o.ImplicitChanges.TotalBreakingChanges()
	}
	if o.PasswordChanges != nil {
		c += o.PasswordChanges.TotalBreakingChanges()
	}
	if o.ClientCredentialsChanges != nil {
		c += o.ClientCredentialsChanges.TotalBreakingChanges()
	}
	if o.AuthorizationCodeChanges != nil {
		c += o.AuthorizationCodeChanges.TotalBreakingChanges()
	}
	if o.ExtensionChanges != nil {
		c += o.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareOAuthFlows compares two AsyncAPI OAuth Flows objects and returns a pointer to
// OAuthFlowsChanges, or nil if nothing changed.
func CompareOAuthFlows(l, r *lowasync.OAuthFlows, configs ...*BreakingRulesConfig) *OAuthFlowsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	o := new(OAuthFlowsChanges)

	compareNestedObject(l.Implicit, r.Implicit, lowasync.ImplicitLabel,
		CompOAuthFlows, PropImplicit, &changes, configuredNestedCompare(config, CompareOAuthFlow), &o.ImplicitChanges, config)

	compareNestedObject(l.Password, r.Password, lowasync.PasswordLabel,
		CompOAuthFlows, PropPassword, &changes, configuredNestedCompare(config, CompareOAuthFlow), &o.PasswordChanges, config)

	compareNestedObject(l.ClientCredentials, r.ClientCredentials, lowasync.ClientCredentialsLabel,
		CompOAuthFlows, PropClientCredentials, &changes, configuredNestedCompare(config, CompareOAuthFlow), &o.ClientCredentialsChanges, config)

	compareNestedObject(l.AuthorizationCode, r.AuthorizationCode, lowasync.AuthorizationCodeLabel,
		CompOAuthFlows, PropAuthorizationCode, &changes, configuredNestedCompare(config, CompareOAuthFlow), &o.AuthorizationCodeChanges, config)

	o.ExtensionChanges = CheckExtensions(l, r)
	o.PropertyChanges = NewPropertyChanges(changes)
	if o.TotalChanges() <= 0 {
		return nil
	}
	return o
}

// SecuritySchemeChanges represents changes made to an AsyncAPI Security Scheme object.
type SecuritySchemeChanges struct {
	*PropertyChanges
	FlowsChanges     *OAuthFlowsChanges `json:"flows,omitempty" yaml:"flows,omitempty"`
	ExtensionChanges *ExtensionChanges  `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Security Scheme objects.
func (s *SecuritySchemeChanges) GetAllChanges() []*Change {
	if s == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, s.Changes...)
	if s.FlowsChanges != nil {
		changes = append(changes, s.FlowsChanges.GetAllChanges()...)
	}
	if s.ExtensionChanges != nil {
		changes = append(changes, s.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (s *SecuritySchemeChanges) TotalChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalChanges()
	if s.FlowsChanges != nil {
		c += s.FlowsChanges.TotalChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (s *SecuritySchemeChanges) TotalBreakingChanges() int {
	if s == nil {
		return 0
	}
	c := s.PropertyChanges.TotalBreakingChanges()
	if s.FlowsChanges != nil {
		c += s.FlowsChanges.TotalBreakingChanges()
	}
	if s.ExtensionChanges != nil {
		c += s.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareSecuritySchemes compares two AsyncAPI Security Scheme objects and returns a
// pointer to SecuritySchemeChanges, or nil if nothing changed.
//
// Scopes is the AsyncAPI specific array-of-strings field, compared as a set; the
// OAuth flow scope maps are handled by CompareOAuthFlow.
func CompareSecuritySchemes(l, r *lowasync.SecurityScheme, configs ...*BreakingRulesConfig) *SecuritySchemeChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompSecurityScheme, PropType,
			l.Type.ValueNode, r.Type.ValueNode, lowasync.TypeLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropIn,
			l.In.ValueNode, r.In.ValueNode, lowasync.InLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropScheme,
			l.Scheme.ValueNode, r.Scheme.ValueNode, lowasync.SchemeLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropBearerFormat,
			l.BearerFormat.ValueNode, r.BearerFormat.ValueNode, lowasync.BearerFormatLabel, &changes, l, r, config),
		NewPropertyCheck(CompSecurityScheme, PropOpenIDConnectURL,
			l.OpenIDConnectURL.ValueNode, r.OpenIDConnectURL.ValueNode, lowasync.OpenIDConnectURLLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	s := new(SecuritySchemeChanges)

	compareNestedObject(l.Flows, r.Flows, lowasync.FlowsLabel,
		CompSecurityScheme, PropFlows, &changes, configuredNestedCompare(config, CompareOAuthFlows), &s.FlowsChanges, config)

	if len(l.Scopes.Value) > 0 || len(r.Scopes.Value) > 0 {
		ExtractStringValueSliceChangesWithRules(l.Scopes.Value, r.Scopes.Value, &changes,
			lowasync.ScopesLabel, CompSecurityScheme, PropScopes, config)
	}

	s.ExtensionChanges = CheckExtensions(l, r)
	s.PropertyChanges = NewPropertyChanges(changes)
	if s.TotalChanges() <= 0 {
		return nil
	}
	return s
}
