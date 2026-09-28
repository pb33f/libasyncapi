// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	"github.com/pb33f/libopenapi/orderedmap"
)

// SecurityScheme represents a high-level AsyncAPI 3.0 Security Scheme object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
type SecurityScheme struct {
	Type             string                              `json:"type,omitempty" yaml:"type,omitempty"`
	Description      string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string                              `json:"name,omitempty" yaml:"name,omitempty"`
	In               string                              `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string                              `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string                              `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows                         `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIDConnectURL string                              `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
	Scopes           []string                            `json:"scopes,omitempty" yaml:"scopes,omitempty"`
	Extensions       *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low              *lowasync.SecurityScheme
}

// NewSecurityScheme creates a new high-level SecurityScheme instance from a low-level one.
func NewSecurityScheme(ss *lowasync.SecurityScheme) *SecurityScheme {
	s := new(SecurityScheme)
	s.low = ss
	s.Type = ss.Type.Value
	s.Description = ss.Description.Value
	s.Name = ss.Name.Value
	s.In = ss.In.Value
	s.Scheme = ss.Scheme.Value
	s.BearerFormat = ss.BearerFormat.Value
	s.OpenIDConnectURL = ss.OpenIDConnectURL.Value
	if !ss.Flows.IsEmpty() {
		s.Flows = NewOAuthFlows(ss.Flows.Value)
	}
	if ss.Scopes.Value != nil {
		for _, sc := range ss.Scopes.Value {
			s.Scopes = append(s.Scopes, sc.Value)
		}
	}
	if orderedmap.Len(ss.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(ss.Extensions)
	}
	return s
}

// GoLow returns the low-level SecurityScheme instance.
func (s *SecurityScheme) GoLow() *lowasync.SecurityScheme {
	return s.low
}

// GoLowUntyped returns the low-level SecurityScheme instance with no type.
func (s *SecurityScheme) GoLowUntyped() any {
	return s.low
}

// Render will return a YAML representation of the SecurityScheme object as a byte slice.
func (s *SecurityScheme) Render() ([]byte, error) {
	return yaml.Marshal(s)
}

// MarshalYAML will create a ready to render YAML representation of the SecurityScheme object.
func (s *SecurityScheme) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(s, s.low)
	return nb.Render(), nil
}

// OAuthFlows represents a high-level AsyncAPI 3.0 OAuth Flows object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowsObject
type OAuthFlows struct {
	Implicit          *OAuthFlow                          `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow                          `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow                          `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow                          `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
	Extensions        *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low               *lowasync.OAuthFlows
}

// NewOAuthFlows creates a new high-level OAuthFlows instance from a low-level one.
func NewOAuthFlows(flows *lowasync.OAuthFlows) *OAuthFlows {
	f := new(OAuthFlows)
	f.low = flows
	if !flows.Implicit.IsEmpty() {
		f.Implicit = NewOAuthFlow(flows.Implicit.Value)
	}
	if !flows.Password.IsEmpty() {
		f.Password = NewOAuthFlow(flows.Password.Value)
	}
	if !flows.ClientCredentials.IsEmpty() {
		f.ClientCredentials = NewOAuthFlow(flows.ClientCredentials.Value)
	}
	if !flows.AuthorizationCode.IsEmpty() {
		f.AuthorizationCode = NewOAuthFlow(flows.AuthorizationCode.Value)
	}
	if orderedmap.Len(flows.Extensions) > 0 {
		f.Extensions = high.ExtractExtensions(flows.Extensions)
	}
	return f
}

// GoLow returns the low-level OAuthFlows instance.
func (f *OAuthFlows) GoLow() *lowasync.OAuthFlows {
	return f.low
}

// GoLowUntyped returns the low-level OAuthFlows instance with no type.
func (f *OAuthFlows) GoLowUntyped() any {
	return f.low
}

// OAuthFlow represents a high-level AsyncAPI 3.0 OAuth Flow object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowObject
type OAuthFlow struct {
	AuthorizationURL string                              `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenURL         string                              `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshURL       string                              `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	AvailableScopes  *orderedmap.Map[string, string]     `json:"availableScopes,omitempty" yaml:"availableScopes,omitempty"`
	Extensions       *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low              *lowasync.OAuthFlow
}

// NewOAuthFlow creates a new high-level OAuthFlow instance from a low-level one.
func NewOAuthFlow(flow *lowasync.OAuthFlow) *OAuthFlow {
	f := new(OAuthFlow)
	f.low = flow
	f.AuthorizationURL = flow.AuthorizationURL.Value
	f.TokenURL = flow.TokenURL.Value
	f.RefreshURL = flow.RefreshURL.Value
	if flow.AvailableScopes.Value != nil {
		f.AvailableScopes = orderedmap.New[string, string]()
		for k, v := range flow.AvailableScopes.Value.FromOldest() {
			f.AvailableScopes.Set(k.Value, v.Value)
		}
	}
	if orderedmap.Len(flow.Extensions) > 0 {
		f.Extensions = high.ExtractExtensions(flow.Extensions)
	}
	return f
}

// GoLow returns the low-level OAuthFlow instance.
func (f *OAuthFlow) GoLow() *lowasync.OAuthFlow {
	return f.low
}

// GoLowUntyped returns the low-level OAuthFlow instance with no type.
func (f *OAuthFlow) GoLowUntyped() any {
	return f.low
}
