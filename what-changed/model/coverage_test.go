// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"context"
	"testing"

	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func buildParameter(t *testing.T, y string) *lowasync.Parameter {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	p := &lowasync.Parameter{}
	require.NoError(t, low.BuildModel(node.Content[0], p))
	require.NoError(t, p.Build(context.Background(), nil, node.Content[0], idx))
	return p
}

func buildSecurityScheme(t *testing.T, y string) *lowasync.SecurityScheme {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	ss := &lowasync.SecurityScheme{}
	require.NoError(t, low.BuildModel(node.Content[0], ss))
	require.NoError(t, ss.Build(context.Background(), nil, node.Content[0], idx))
	return ss
}

func buildOAuthFlow(t *testing.T, y string) *lowasync.OAuthFlow {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	f := &lowasync.OAuthFlow{}
	require.NoError(t, low.BuildModel(node.Content[0], f))
	require.NoError(t, f.Build(context.Background(), nil, node.Content[0], idx))
	return f
}

func buildOAuthFlows(t *testing.T, y string) *lowasync.OAuthFlows {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	f := &lowasync.OAuthFlows{}
	require.NoError(t, low.BuildModel(node.Content[0], f))
	require.NoError(t, f.Build(context.Background(), nil, node.Content[0], idx))
	return f
}

func buildServerVariable(t *testing.T, y string) *lowasync.ServerVariable {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	sv := &lowasync.ServerVariable{}
	require.NoError(t, low.BuildModel(node.Content[0], sv))
	require.NoError(t, sv.Build(context.Background(), nil, node.Content[0], idx))
	return sv
}

func buildCorrelationID(t *testing.T, y string) *lowasync.CorrelationID {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	c := &lowasync.CorrelationID{}
	require.NoError(t, low.BuildModel(node.Content[0], c))
	require.NoError(t, c.Build(context.Background(), nil, node.Content[0], idx))
	return c
}

func buildOperationReply(t *testing.T, y string) *lowasync.OperationReply {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	or := &lowasync.OperationReply{}
	require.NoError(t, low.BuildModel(node.Content[0], or))
	require.NoError(t, or.Build(context.Background(), nil, node.Content[0], idx))
	return or
}

func buildOperationReplyAddress(t *testing.T, y string) *lowasync.OperationReplyAddress {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	ora := &lowasync.OperationReplyAddress{}
	require.NoError(t, low.BuildModel(node.Content[0], ora))
	require.NoError(t, ora.Build(context.Background(), nil, node.Content[0], idx))
	return ora
}

func buildOperationTrait(t *testing.T, y string) *lowasync.OperationTrait {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	ot := &lowasync.OperationTrait{}
	require.NoError(t, low.BuildModel(node.Content[0], ot))
	require.NoError(t, ot.Build(context.Background(), nil, node.Content[0], idx))
	return ot
}

func buildMessageTrait(t *testing.T, y string) *lowasync.MessageTrait {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	mt := &lowasync.MessageTrait{}
	require.NoError(t, low.BuildModel(node.Content[0], mt))
	require.NoError(t, mt.Build(context.Background(), nil, node.Content[0], idx))
	return mt
}

func buildServer(t *testing.T, y string) *lowasync.Server {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	s := &lowasync.Server{}
	require.NoError(t, low.BuildModel(node.Content[0], s))
	require.NoError(t, s.Build(context.Background(), nil, node.Content[0], idx))
	return s
}

func buildExternalDoc(t *testing.T, y string) *lowasync.ExternalDoc {
	t.Helper()
	var node yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte(y), &node))
	idx := index.NewSpecIndexWithConfig(&node, index.CreateOpenAPIIndexConfig())
	ed := &lowasync.ExternalDoc{}
	require.NoError(t, low.BuildModel(node.Content[0], ed))
	require.NoError(t, ed.Build(context.Background(), nil, node.Content[0], idx))
	return ed
}

// --- CompareParameters ---------------------------------------------------------------

func TestCompareParameters_LocationModified(t *testing.T) {
	left := buildParameter(t, `location: $message.payload#/region
description: a region`)
	right := buildParameter(t, `location: $message.payload#/zone
description: a region`)

	changes := CompareParameters(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "parameter.location modification is breaking")
	assert.Equal(t, "$message.payload#/region", ch.Original)
	assert.Equal(t, "$message.payload#/zone", ch.New)
}

func TestCompareParameters_EnumAddedAndRemoved(t *testing.T) {
	left := buildParameter(t, `enum:
  - us-east
  - us-west`)
	right := buildParameter(t, `enum:
  - us-east
  - eu-central`)

	changes := CompareParameters(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	var added, removed *Change
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case PropertyAdded:
			added = ch
		case PropertyRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
	assert.False(t, added.Breaking, "parameter.enum addition is not breaking")
	assert.True(t, removed.Breaking, "parameter.enum removal is breaking")
	assert.Equal(t, "us-west", removed.Original)
	assert.Equal(t, "eu-central", added.New)
}

func TestCompareParameters_DescriptionNotBreaking(t *testing.T) {
	left := buildParameter(t, `description: old words`)
	right := buildParameter(t, `description: new words`)

	changes := CompareParameters(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	assert.Equal(t, Modified, changes.Changes[0].ChangeType)
}

func TestCompareParameters_IdenticalAndNil(t *testing.T) {
	p := buildParameter(t, `location: $message.payload#/id`)
	assert.Nil(t, CompareParameters(p, buildParameter(t, `location: $message.payload#/id`)))
	assert.Nil(t, CompareParameters(nil, p))
	assert.Nil(t, CompareParameters(p, nil))
}

// --- CompareSecuritySchemes ----------------------------------------------------------

// securityScheme.type collides with OpenAPI rule names — the breaking result must come
// from THIS package's config, where modification is breaking.
func TestCompareSecuritySchemes_TypeModifiedBreaking(t *testing.T) {
	left := buildSecurityScheme(t, `type: apiKey
in: user`)
	right := buildSecurityScheme(t, `type: http
in: user`)

	changes := CompareSecuritySchemes(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "securityScheme.type modification is breaking per OUR config")
	assert.Equal(t, "apiKey", ch.Original)
	assert.Equal(t, "http", ch.New)
}

func TestCompareSecuritySchemes_ScopesAddedAndRemoved(t *testing.T) {
	left := buildSecurityScheme(t, `type: oauth2
scopes:
  - read
  - write`)
	right := buildSecurityScheme(t, `type: oauth2
scopes:
  - read
  - admin`)

	changes := CompareSecuritySchemes(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	var added, removed *Change
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case PropertyAdded:
			added = ch
		case PropertyRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
	assert.False(t, added.Breaking, "securityScheme.scopes addition is not breaking")
	assert.True(t, removed.Breaking, "securityScheme.scopes removal is breaking")
	assert.Equal(t, "write", removed.Original)
	assert.Equal(t, "admin", added.New)
}

func TestCompareSecuritySchemes_DescriptionNotBreaking(t *testing.T) {
	left := buildSecurityScheme(t, `type: apiKey
description: old`)
	right := buildSecurityScheme(t, `type: apiKey
description: new`)

	changes := CompareSecuritySchemes(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

// flipping securityScheme.type to non-breaking via a custom config must be honored,
// proving the lookup resolves against this package's active configuration.
func TestCompareSecuritySchemes_CustomTypeRuleHonored(t *testing.T) {
	left := buildSecurityScheme(t, `type: apiKey`)
	right := buildSecurityScheme(t, `type: http`)

	// default: breaking.
	changes := CompareSecuritySchemes(left, right)
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalBreakingChanges())

	custom := NewDefaultBreakingRulesConfig()
	custom.Merge(&BreakingRulesConfig{
		SecurityScheme: &SecuritySchemeRules{Type: rule(false, false, false)},
	})
	SetActiveBreakingRulesConfig(custom)
	defer ResetActiveBreakingRulesConfig()

	changes = CompareSecuritySchemes(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	assert.False(t, changes.Changes[0].Breaking)
}

func TestCompareSecuritySchemes_IdenticalAndNil(t *testing.T) {
	s := buildSecurityScheme(t, `type: apiKey`)
	assert.Nil(t, CompareSecuritySchemes(s, buildSecurityScheme(t, `type: apiKey`)))
	assert.Nil(t, CompareSecuritySchemes(nil, s))
	assert.Nil(t, CompareSecuritySchemes(s, nil))
}

// --- CompareOAuthFlow / CompareOAuthFlows --------------------------------------------

func TestCompareOAuthFlow_TokenURLModifiedBreaking(t *testing.T) {
	left := buildOAuthFlow(t, `tokenUrl: https://old.example.com/token`)
	right := buildOAuthFlow(t, `tokenUrl: https://new.example.com/token`)

	changes := CompareOAuthFlow(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "oauthFlow.tokenUrl modification is breaking")
	assert.Equal(t, "https://old.example.com/token", ch.Original)
	assert.Equal(t, "https://new.example.com/token", ch.New)
}

func TestCompareOAuthFlow_AvailableScopesChanges(t *testing.T) {
	base := `authorizationUrl: https://auth.example.com
availableScopes:
  read: read stuff
  write: write stuff`

	// scope added: not breaking per oauthFlow.availableScopes (f,t,t).
	added := buildOAuthFlow(t, base+`
  admin: admin stuff`)
	changes := CompareOAuthFlow(buildOAuthFlow(t, base), added)
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, PropertyAdded, changes.Changes[0].ChangeType)
	assert.False(t, changes.Changes[0].Breaking)
	assert.Equal(t, "admin stuff", changes.Changes[0].New)

	// scope removed: breaking.
	changes = CompareOAuthFlow(added, buildOAuthFlow(t, base))
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, PropertyRemoved, changes.Changes[0].ChangeType)
	assert.True(t, changes.Changes[0].Breaking)

	// scope description modified: breaking.
	modified := buildOAuthFlow(t, `authorizationUrl: https://auth.example.com
availableScopes:
  read: read more stuff
  write: write stuff`)
	changes = CompareOAuthFlow(buildOAuthFlow(t, base), modified)
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, Modified, changes.Changes[0].ChangeType)
	assert.True(t, changes.Changes[0].Breaking)
	assert.Equal(t, "read stuff", changes.Changes[0].Original)
	assert.Equal(t, "read more stuff", changes.Changes[0].New)
}

func TestCompareOAuthFlow_IdenticalAndNil(t *testing.T) {
	f := buildOAuthFlow(t, `tokenUrl: https://t.example.com`)
	assert.Nil(t, CompareOAuthFlow(f, buildOAuthFlow(t, `tokenUrl: https://t.example.com`)))
	assert.Nil(t, CompareOAuthFlow(nil, f))
	assert.Nil(t, CompareOAuthFlow(f, nil))
}

func TestCompareOAuthFlows_FlowRemovedAndAdded(t *testing.T) {
	left := buildOAuthFlows(t, `implicit:
  authorizationUrl: https://auth.example.com
  availableScopes:
    read: read stuff`)
	right := buildOAuthFlows(t, `password:
  tokenUrl: https://token.example.com
  availableScopes:
    read: read stuff`)

	changes := CompareOAuthFlows(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	var added, removed *Change
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case ObjectAdded:
			added = ch
		case ObjectRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
	assert.True(t, removed.Breaking, "removing the implicit flow is breaking")
	assert.False(t, added.Breaking, "adding a password flow is not breaking")
}

func TestCompareOAuthFlows_NestedFlowChangeIsGranular(t *testing.T) {
	left := buildOAuthFlows(t, `implicit:
  authorizationUrl: https://auth.example.com
  availableScopes:
    read: read stuff`)
	right := buildOAuthFlows(t, `implicit:
  authorizationUrl: https://auth2.example.com
  availableScopes:
    read: read stuff`)

	changes := CompareOAuthFlows(left, right)
	require.NotNil(t, changes)
	require.NotNil(t, changes.ImplicitChanges)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.Equal(t, Modified, changes.ImplicitChanges.Changes[0].ChangeType)
}

func TestCompareOAuthFlows_IdenticalAndNil(t *testing.T) {
	f := buildOAuthFlows(t, `implicit:
  authorizationUrl: https://auth.example.com`)
	assert.Nil(t, CompareOAuthFlows(f, buildOAuthFlows(t, `implicit:
  authorizationUrl: https://auth.example.com`)))
	assert.Nil(t, CompareOAuthFlows(nil, f))
	assert.Nil(t, CompareOAuthFlows(f, nil))
}

// --- CompareServerVariables ----------------------------------------------------------

func TestCompareServerVariables_DefaultModifiedBreaking(t *testing.T) {
	left := buildServerVariable(t, `default: "8080"`)
	right := buildServerVariable(t, `default: "9090"`)

	changes := CompareServerVariables(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "serverVariable.default modification is breaking")
	assert.Equal(t, "8080", ch.Original)
	assert.Equal(t, "9090", ch.New)
}

func TestCompareServerVariables_EnumEntryRemovedBreaking(t *testing.T) {
	left := buildServerVariable(t, `default: "8080"
enum:
  - "8080"
  - "9090"`)
	right := buildServerVariable(t, `default: "8080"
enum:
  - "8080"`)

	changes := CompareServerVariables(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.Equal(t, PropertyRemoved, changes.Changes[0].ChangeType)
	assert.Equal(t, "9090", changes.Changes[0].Original)
}

func TestCompareServerVariables_ExamplesNotBreaking(t *testing.T) {
	left := buildServerVariable(t, `default: "8080"
examples:
  - "8080"`)
	right := buildServerVariable(t, `default: "8080"
examples:
  - "9090"`)

	changes := CompareServerVariables(left, right)
	require.NotNil(t, changes)
	// example swap is one removal plus one addition, neither breaking.
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

func TestCompareServerVariables_IdenticalAndNil(t *testing.T) {
	sv := buildServerVariable(t, `default: "8080"`)
	assert.Nil(t, CompareServerVariables(sv, buildServerVariable(t, `default: "8080"`)))
	assert.Nil(t, CompareServerVariables(nil, sv))
	assert.Nil(t, CompareServerVariables(sv, nil))
}

// --- CompareCorrelationID ------------------------------------------------------------

func TestCompareCorrelationID_LocationModifiedBreaking(t *testing.T) {
	left := buildCorrelationID(t, `location: $message.header#/correlationId
description: correlates`)
	right := buildCorrelationID(t, `location: $message.header#/traceId
description: correlates`)

	changes := CompareCorrelationID(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "correlationId.location modification is breaking")
	assert.Equal(t, "$message.header#/correlationId", ch.Original)
	assert.Equal(t, "$message.header#/traceId", ch.New)
}

func TestCompareCorrelationID_DescriptionNotBreaking(t *testing.T) {
	left := buildCorrelationID(t, `location: $message.header#/correlationId
description: old`)
	right := buildCorrelationID(t, `location: $message.header#/correlationId
description: new`)

	changes := CompareCorrelationID(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

func TestCompareCorrelationID_IdenticalAndNil(t *testing.T) {
	c := buildCorrelationID(t, `location: $message.header#/correlationId`)
	assert.Nil(t, CompareCorrelationID(c, buildCorrelationID(t, `location: $message.header#/correlationId`)))
	assert.Nil(t, CompareCorrelationID(nil, c))
	assert.Nil(t, CompareCorrelationID(c, nil))
}

// --- CompareOperationReply / CompareOperationReplyAddress ----------------------------

func TestCompareOperationReplyAddress_LocationModifiedBreaking(t *testing.T) {
	left := buildOperationReplyAddress(t, `location: $message.header#/replyTo
description: where to reply`)
	right := buildOperationReplyAddress(t, `location: $message.header#/replyQueue
description: where to reply`)

	changes := CompareOperationReplyAddress(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "operationReplyAddress.location modification is breaking")
	assert.Equal(t, "$message.header#/replyTo", ch.Original)
	assert.Equal(t, "$message.header#/replyQueue", ch.New)
}

func TestCompareOperationReplyAddress_IdenticalAndNil(t *testing.T) {
	a := buildOperationReplyAddress(t, `location: $message.header#/replyTo`)
	assert.Nil(t, CompareOperationReplyAddress(a, buildOperationReplyAddress(t, `location: $message.header#/replyTo`)))
	assert.Nil(t, CompareOperationReplyAddress(nil, a))
	assert.Nil(t, CompareOperationReplyAddress(a, nil))
}

func TestCompareOperationReply_AddressLocationChangeIsGranular(t *testing.T) {
	left := buildOperationReply(t, `address:
  location: $message.header#/replyTo
channel:
  $ref: '#/channels/x'`)
	right := buildOperationReply(t, `address:
  location: $message.header#/replyQueue
channel:
  $ref: '#/channels/x'`)

	changes := CompareOperationReply(left, right)
	require.NotNil(t, changes)
	require.NotNil(t, changes.AddressChanges)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.Equal(t, Modified, changes.AddressChanges.Changes[0].ChangeType)
}

func TestCompareOperationReply_ChannelRefSwapped(t *testing.T) {
	left := buildOperationReply(t, `channel:
  $ref: '#/channels/x'`)
	right := buildOperationReply(t, `channel:
  $ref: '#/channels/y'`)

	changes := CompareOperationReply(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "operationReply.channel modification is breaking")
	assert.Equal(t, "#/channels/x", ch.Original)
	assert.Equal(t, "#/channels/y", ch.New)
}

func TestCompareOperationReply_MessagesRefAddedAndRemoved(t *testing.T) {
	one := buildOperationReply(t, `channel:
  $ref: '#/channels/x'
messages:
  - $ref: '#/messages/a'`)
	two := buildOperationReply(t, `channel:
  $ref: '#/channels/x'
messages:
  - $ref: '#/messages/a'
  - $ref: '#/messages/b'`)

	// message ref added: breaking per operationReply.messages (t,f,t).
	changes := CompareOperationReply(one, two)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.Equal(t, ObjectAdded, changes.Changes[0].ChangeType)
	assert.Equal(t, "#/messages/b", changes.Changes[0].New)

	// message ref removed: breaking.
	changes = CompareOperationReply(two, one)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())
	assert.Equal(t, ObjectRemoved, changes.Changes[0].ChangeType)
	assert.Equal(t, "#/messages/b", changes.Changes[0].Original)
}

func TestCompareOperationReply_IdenticalAndNil(t *testing.T) {
	r := buildOperationReply(t, `channel:
  $ref: '#/channels/x'`)
	assert.Nil(t, CompareOperationReply(r, buildOperationReply(t, `channel:
  $ref: '#/channels/x'`)))
	assert.Nil(t, CompareOperationReply(nil, r))
	assert.Nil(t, CompareOperationReply(r, nil))
}

// --- CompareOperationTraits ----------------------------------------------------------

func TestCompareOperationTraits_InformationalChangesNotBreaking(t *testing.T) {
	left := buildOperationTrait(t, `title: old title
description: old words`)
	right := buildOperationTrait(t, `title: new title
description: new words`)

	changes := CompareOperationTraits(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	for _, ch := range changes.Changes {
		assert.Equal(t, Modified, ch.ChangeType)
		assert.False(t, ch.Breaking)
	}
}

func TestCompareOperationTraits_IdenticalAndNil(t *testing.T) {
	ot := buildOperationTrait(t, `title: a trait`)
	assert.Nil(t, CompareOperationTraits(ot, buildOperationTrait(t, `title: a trait`)))
	assert.Nil(t, CompareOperationTraits(nil, ot))
	assert.Nil(t, CompareOperationTraits(ot, nil))
}

// --- CompareMessageTraits ------------------------------------------------------------

func TestCompareMessageTraits_ContentTypeModifiedBreaking(t *testing.T) {
	left := buildMessageTrait(t, `contentType: application/json
title: order trait`)
	right := buildMessageTrait(t, `contentType: application/xml
title: order trait`)

	changes := CompareMessageTraits(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "messageTrait.contentType modification is breaking")
	assert.Equal(t, "application/json", ch.Original)
	assert.Equal(t, "application/xml", ch.New)
}

func TestCompareMessageTraits_TitleNotBreaking(t *testing.T) {
	left := buildMessageTrait(t, `contentType: application/json
title: old title`)
	right := buildMessageTrait(t, `contentType: application/json
title: new title`)

	changes := CompareMessageTraits(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
}

func TestCompareMessageTraits_NamedExampleEditIsGranular(t *testing.T) {
	left := buildMessageTrait(t, `title: order trait
examples:
  - name: minimal
    summary: a minimal example`)
	right := buildMessageTrait(t, `title: order trait
examples:
  - name: minimal
    summary: a minimal example, updated`)

	changes := CompareMessageTraits(left, right)
	require.NotNil(t, changes)
	require.Contains(t, changes.ExampleChanges, "minimal")
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	ch := changes.ExampleChanges["minimal"].Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.Equal(t, "a minimal example", ch.Original)
	assert.Equal(t, "a minimal example, updated", ch.New)
}

func TestCompareMessageTraits_IdenticalAndNil(t *testing.T) {
	mt := buildMessageTrait(t, `title: a trait`)
	assert.Nil(t, CompareMessageTraits(mt, buildMessageTrait(t, `title: a trait`)))
	assert.Nil(t, CompareMessageTraits(nil, mt))
	assert.Nil(t, CompareMessageTraits(mt, nil))
}

// --- CompareTagSlices (via CompareServers) -------------------------------------------

func TestCompareServers_TagAddedAndRemoved(t *testing.T) {
	one := buildServer(t, `host: example.com
protocol: kafka
tags:
  - name: env
    description: environment tag`)
	two := buildServer(t, `host: example.com
protocol: kafka
tags:
  - name: env
    description: environment tag
  - name: region
    description: region tag`)

	// tag added: not breaking per server.tags (f,f,f).
	changes := CompareServers(one, two)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	assert.Equal(t, ObjectAdded, changes.Changes[0].ChangeType)

	// tag removed: not breaking per server.tags.
	changes = CompareServers(two, one)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	assert.Equal(t, ObjectRemoved, changes.Changes[0].ChangeType)
}

func TestCompareServers_TagDescriptionModifiedIsGranular(t *testing.T) {
	left := buildServer(t, `host: example.com
protocol: kafka
tags:
  - name: env
    description: production`)
	right := buildServer(t, `host: example.com
protocol: kafka
tags:
  - name: env
    description: staging`)

	changes := CompareServers(left, right)
	require.NotNil(t, changes)
	require.Contains(t, changes.TagChanges, "env")
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())

	ch := changes.TagChanges["env"].Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.Equal(t, "production", ch.Original)
	assert.Equal(t, "staging", ch.New)
}

// --- CompareServers scalars / CompareExternalDocs ------------------------------------

func TestCompareServers_ProtocolModifiedBreaking(t *testing.T) {
	left := buildServer(t, `host: example.com
protocol: kafka`)
	right := buildServer(t, `host: example.com
protocol: mqtt`)

	changes := CompareServers(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Equal(t, 1, changes.TotalBreakingChanges())

	ch := changes.Changes[0]
	assert.Equal(t, Modified, ch.ChangeType)
	assert.True(t, ch.Breaking, "server.protocol modification is breaking")
	assert.Equal(t, "kafka", ch.Original)
	assert.Equal(t, "mqtt", ch.New)
}

func TestCompareServers_IdenticalAndNil(t *testing.T) {
	s := buildServer(t, `host: example.com
protocol: kafka`)
	assert.Nil(t, CompareServers(s, buildServer(t, `host: example.com
protocol: kafka`)))
	assert.Nil(t, CompareServers(nil, s))
	assert.Nil(t, CompareServers(s, nil))
}

func TestCompareExternalDocs_NotBreaking(t *testing.T) {
	left := buildExternalDoc(t, `url: https://docs.example.com
description: the docs`)
	right := buildExternalDoc(t, `url: https://docs2.example.com
description: the better docs`)

	changes := CompareExternalDocs(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 2, changes.TotalChanges())
	assert.Equal(t, 0, changes.TotalBreakingChanges())
	for _, ch := range changes.Changes {
		assert.Equal(t, Modified, ch.ChangeType)
		assert.False(t, ch.Breaking)
	}
}

func TestCompareExternalDocs_IdenticalAndNil(t *testing.T) {
	ed := buildExternalDoc(t, `url: https://docs.example.com`)
	assert.Nil(t, CompareExternalDocs(ed, buildExternalDoc(t, `url: https://docs.example.com`)))
	assert.Nil(t, CompareExternalDocs(nil, ed))
	assert.Nil(t, CompareExternalDocs(ed, nil))
}
