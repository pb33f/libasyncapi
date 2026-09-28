// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"hash/maphash"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/orderedmap"
	"github.com/pb33f/libopenapi/utils"
)

// SecurityScheme represents a low-level AsyncAPI 3.0 Security Scheme object.
//
// Defines a security scheme that can be used by the operations. Supported schemes are:
// User/Password, API key, X.509, Symmetric/Asymmetric Encryption, HTTP (including Bearer),
// OAuth2, OpenID Connect, SASL (plain, scramSha256, scramSha512, gssapi).
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#securitySchemeObject
type SecurityScheme struct {
	Type             low.NodeReference[string]
	Description      low.NodeReference[string]
	Name             low.NodeReference[string]
	In               low.NodeReference[string]
	Scheme           low.NodeReference[string]
	BearerFormat     low.NodeReference[string]
	Flows            low.NodeReference[*OAuthFlows]
	OpenIDConnectURL low.NodeReference[string]
	Scopes           low.NodeReference[[]low.ValueReference[string]]
	Extensions       *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode          *yaml.Node
	RootNode         *yaml.Node
	idx              *index.SpecIndex
	ctx              context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the SecurityScheme object.
func (ss *SecurityScheme) GetRootNode() *yaml.Node {
	return ss.RootNode
}

// GetKeyNode returns the key yaml node of the SecurityScheme object.
func (ss *SecurityScheme) GetKeyNode() *yaml.Node {
	return ss.KeyNode
}

// GetExtensions returns all extensions for SecurityScheme.
func (ss *SecurityScheme) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return ss.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (ss *SecurityScheme) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, ss.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the SecurityScheme object.
func (ss *SecurityScheme) GetIndex() *index.SpecIndex {
	return ss.idx
}

// GetContext returns the context.Context instance used when building the SecurityScheme object.
func (ss *SecurityScheme) GetContext() context.Context {
	return ss.ctx
}

// Build extracts the SecurityScheme object from the supplied root node.
func (ss *SecurityScheme) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	ss.KeyNode = keyNode
	ss.Reference = new(low.Reference)
	if ok, _, ref := utils.IsNodeRefValue(root); ok {
		ss.SetReference(ref, root)
	}
	root = utils.NodeAlias(root)
	ss.RootNode = root
	utils.CheckForMergeNodes(root)
	ss.Nodes = low.ExtractNodes(ctx, root)
	ss.Extensions = low.ExtractExtensions(root)
	ss.idx = idx
	ss.ctx = ctx

	// extract flows
	flows, err := low.ExtractObject[*OAuthFlows](ctx, FlowsLabel, root, idx)
	if err != nil {
		return err
	}
	if flows.Value != nil {
		ss.Flows = flows
	}

	// extract scopes (AsyncAPI-specific: array of strings)
	_, scopesLabel, scopesValue := utils.FindKeyNodeFullTop(ScopesLabel, root.Content)
	if scopesValue != nil && scopesValue.Kind == yaml.SequenceNode {
		ss.Nodes.Store(scopesLabel.Line, scopesLabel)
		var scopeRefs []low.ValueReference[string]
		for _, item := range scopesValue.Content {
			scopeRefs = append(scopeRefs, low.ValueReference[string]{
				Value:     item.Value,
				ValueNode: item,
			})
		}
		ss.Scopes = low.NodeReference[[]low.ValueReference[string]]{
			Value:     scopeRefs,
			KeyNode:   scopesLabel,
			ValueNode: scopesValue,
		}
	}

	return nil
}

// Hash returns a process-local content hash of the SecurityScheme object.
func (ss *SecurityScheme) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !ss.Type.IsEmpty() {
			h.WriteString(ss.Type.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.Description.IsEmpty() {
			h.WriteString(ss.Description.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.Name.IsEmpty() {
			h.WriteString(ss.Name.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.In.IsEmpty() {
			h.WriteString(ss.In.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.Scheme.IsEmpty() {
			h.WriteString(ss.Scheme.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.BearerFormat.IsEmpty() {
			h.WriteString(ss.BearerFormat.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.Flows.IsEmpty() {
			h.WriteString(low.GenerateHashString(ss.Flows.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !ss.OpenIDConnectURL.IsEmpty() {
			h.WriteString(ss.OpenIDConnectURL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if ss.Scopes.Value != nil {
			for _, s := range ss.Scopes.Value {
				h.WriteString(s.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		for _, ext := range low.HashExtensions(ss.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OAuthFlows represents a low-level AsyncAPI 3.0 OAuth Flows object.
//
// Allows configuration of the supported OAuth Flows.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowsObject
type OAuthFlows struct {
	Implicit          low.NodeReference[*OAuthFlow]
	Password          low.NodeReference[*OAuthFlow]
	ClientCredentials low.NodeReference[*OAuthFlow]
	AuthorizationCode low.NodeReference[*OAuthFlow]
	Extensions        *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode           *yaml.Node
	RootNode          *yaml.Node
	idx               *index.SpecIndex
	ctx               context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the OAuthFlows object.
func (of *OAuthFlows) GetRootNode() *yaml.Node {
	return of.RootNode
}

// GetKeyNode returns the key yaml node of the OAuthFlows object.
func (of *OAuthFlows) GetKeyNode() *yaml.Node {
	return of.KeyNode
}

// GetExtensions returns all extensions for OAuthFlows.
func (of *OAuthFlows) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return of.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (of *OAuthFlows) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, of.Extensions)
}

// GetIndex returns the index.SpecIndex instance attached to the OAuthFlows object.
func (of *OAuthFlows) GetIndex() *index.SpecIndex {
	return of.idx
}

// GetContext returns the context.Context instance used when building the OAuthFlows object.
func (of *OAuthFlows) GetContext() context.Context {
	return of.ctx
}

// Build extracts the OAuthFlows object from the supplied root node.
func (of *OAuthFlows) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	of.KeyNode = keyNode
	root = utils.NodeAlias(root)
	of.RootNode = root
	utils.CheckForMergeNodes(root)
	of.Reference = new(low.Reference)
	of.Nodes = low.ExtractNodes(ctx, root)
	of.Extensions = low.ExtractExtensions(root)
	of.idx = idx
	of.ctx = ctx

	v, err := low.ExtractObject[*OAuthFlow](ctx, ImplicitLabel, root, idx)
	if err != nil {
		return err
	}
	of.Implicit = v

	v, err = low.ExtractObject[*OAuthFlow](ctx, PasswordLabel, root, idx)
	if err != nil {
		return err
	}
	of.Password = v

	v, err = low.ExtractObject[*OAuthFlow](ctx, ClientCredentialsLabel, root, idx)
	if err != nil {
		return err
	}
	of.ClientCredentials = v

	v, err = low.ExtractObject[*OAuthFlow](ctx, AuthorizationCodeLabel, root, idx)
	if err != nil {
		return err
	}
	of.AuthorizationCode = v

	return nil
}

// Hash returns a process-local content hash of the OAuthFlows object.
func (of *OAuthFlows) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !of.Implicit.IsEmpty() {
			h.WriteString(low.GenerateHashString(of.Implicit.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !of.Password.IsEmpty() {
			h.WriteString(low.GenerateHashString(of.Password.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !of.ClientCredentials.IsEmpty() {
			h.WriteString(low.GenerateHashString(of.ClientCredentials.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		if !of.AuthorizationCode.IsEmpty() {
			h.WriteString(low.GenerateHashString(of.AuthorizationCode.Value))
			h.WriteByte(low.HASH_PIPE)
		}
		for _, ext := range low.HashExtensions(of.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}

// OAuthFlow represents a low-level AsyncAPI 3.0 OAuth Flow object.
//
// Configuration details for a supported OAuth Flow.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#oauthFlowObject
type OAuthFlow struct {
	AuthorizationURL low.NodeReference[string]
	TokenURL         low.NodeReference[string]
	RefreshURL       low.NodeReference[string]
	AvailableScopes  low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[string]]]
	Extensions       *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]]
	KeyNode          *yaml.Node
	RootNode         *yaml.Node
	idx              *index.SpecIndex
	ctx              context.Context
	*low.Reference
	low.NodeMap
}

// GetRootNode returns the root yaml node of the OAuthFlow object.
func (of *OAuthFlow) GetRootNode() *yaml.Node {
	return of.RootNode
}

// GetKeyNode returns the key yaml node of the OAuthFlow object.
func (of *OAuthFlow) GetKeyNode() *yaml.Node {
	return of.KeyNode
}

// GetExtensions returns all extensions for OAuthFlow.
func (of *OAuthFlow) GetExtensions() *orderedmap.Map[low.KeyReference[string], low.ValueReference[*yaml.Node]] {
	return of.Extensions
}

// FindExtension attempts to locate an extension with the supplied key.
func (of *OAuthFlow) FindExtension(ext string) *low.ValueReference[*yaml.Node] {
	return low.FindItemInOrderedMap(ext, of.Extensions)
}

// FindScope attempts to locate a scope using a specified name.
func (of *OAuthFlow) FindScope(scope string) *low.ValueReference[string] {
	return low.FindItemInOrderedMap[string](scope, of.AvailableScopes.Value)
}

// GetIndex returns the index.SpecIndex instance attached to the OAuthFlow object.
func (of *OAuthFlow) GetIndex() *index.SpecIndex {
	return of.idx
}

// GetContext returns the context.Context instance used when building the OAuthFlow object.
func (of *OAuthFlow) GetContext() context.Context {
	return of.ctx
}

// Build extracts the OAuthFlow object from the supplied root node.
func (of *OAuthFlow) Build(ctx context.Context, keyNode, root *yaml.Node, idx *index.SpecIndex) error {
	of.KeyNode = keyNode
	root = utils.NodeAlias(root)
	of.RootNode = root
	utils.CheckForMergeNodes(root)
	of.Reference = new(low.Reference)
	of.Nodes = low.ExtractNodes(ctx, root)
	of.Extensions = low.ExtractExtensions(root)
	of.idx = idx
	of.ctx = ctx

	// extract availableScopes map (map of string to string)
	_, scopesLabel, scopesValue := utils.FindKeyNodeFullTop(AvailableScopesLabel, root.Content)
	if scopesValue != nil && scopesValue.Kind == yaml.MappingNode {
		of.Nodes.Store(scopesLabel.Line, scopesLabel)
		scopesMap := orderedmap.New[low.KeyReference[string], low.ValueReference[string]]()
		for i := 0; i < len(scopesValue.Content); i += 2 {
			keyNode := scopesValue.Content[i]
			valueNode := scopesValue.Content[i+1]
			scopesMap.Set(
				low.KeyReference[string]{
					Value:   keyNode.Value,
					KeyNode: keyNode,
				},
				low.ValueReference[string]{
					Value:     valueNode.Value,
					ValueNode: valueNode,
				},
			)
		}
		of.AvailableScopes = low.NodeReference[*orderedmap.Map[low.KeyReference[string], low.ValueReference[string]]]{
			Value:     scopesMap,
			KeyNode:   scopesLabel,
			ValueNode: scopesValue,
		}
	}

	return nil
}

// Hash returns a process-local content hash of the OAuthFlow object.
func (of *OAuthFlow) Hash() uint64 {
	return low.WithHasher(func(h *maphash.Hash) uint64 {
		if !of.AuthorizationURL.IsEmpty() {
			h.WriteString(of.AuthorizationURL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !of.TokenURL.IsEmpty() {
			h.WriteString(of.TokenURL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if !of.RefreshURL.IsEmpty() {
			h.WriteString(of.RefreshURL.Value)
			h.WriteByte(low.HASH_PIPE)
		}
		if of.AvailableScopes.Value != nil {
			for k, v := range orderedmap.SortAlpha(of.AvailableScopes.Value).FromOldest() {
				h.WriteString(k.Value)
				h.WriteByte(low.HASH_PIPE)
				h.WriteString(v.Value)
				h.WriteByte(low.HASH_PIPE)
			}
		}
		for _, ext := range low.HashExtensions(of.Extensions) {
			h.WriteString(ext)
			h.WriteByte(low.HASH_PIPE)
		}
		return h.Sum64()
	})
}
