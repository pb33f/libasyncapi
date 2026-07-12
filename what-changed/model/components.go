// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// hasChanges is the local contract every *XChanges type satisfies, used to fold
// map-valued children without per-field boilerplate.
type hasChanges interface {
	GetAllChanges() []*Change
	TotalChanges() int
	TotalBreakingChanges() int
}

// appendMapChanges collects all changes from a map of child change objects.
func appendMapChanges[T hasChanges](changes []*Change, m map[string]T) []*Change {
	for k := range m {
		changes = append(changes, m[k].GetAllChanges()...)
	}
	return changes
}

// countMapChanges sums TotalChanges across a map of child change objects.
func countMapChanges[T hasChanges](m map[string]T) int {
	c := 0
	for k := range m {
		c += m[k].TotalChanges()
	}
	return c
}

// countMapBreakingChanges sums TotalBreakingChanges across a map of child change objects.
func countMapBreakingChanges[T hasChanges](m map[string]T) int {
	c := 0
	for k := range m {
		c += m[k].TotalBreakingChanges()
	}
	return c
}

// ComponentsChanges represents changes made to the Components object of an AsyncAPI document.
type ComponentsChanges struct {
	*PropertyChanges
	SchemaChanges            map[string]*wcmodel.SchemaChanges        `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	ServerChanges            map[string]*ServerChanges                `json:"servers,omitempty" yaml:"servers,omitempty"`
	ChannelChanges           map[string]*ChannelChanges               `json:"channels,omitempty" yaml:"channels,omitempty"`
	OperationChanges         map[string]*OperationChanges             `json:"operations,omitempty" yaml:"operations,omitempty"`
	MessageChanges           map[string]*MessageChanges               `json:"messages,omitempty" yaml:"messages,omitempty"`
	SecuritySchemeChanges    map[string]*SecuritySchemeChanges        `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
	ServerVariableChanges    map[string]*ServerVariableChanges        `json:"serverVariables,omitempty" yaml:"serverVariables,omitempty"`
	ParameterChanges         map[string]*ParameterChanges             `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	CorrelationIDChanges     map[string]*CorrelationIDChanges         `json:"correlationIds,omitempty" yaml:"correlationIds,omitempty"`
	ReplyChanges             map[string]*OperationReplyChanges        `json:"replies,omitempty" yaml:"replies,omitempty"`
	ReplyAddressChanges      map[string]*OperationReplyAddressChanges `json:"replyAddresses,omitempty" yaml:"replyAddresses,omitempty"`
	ExternalDocChanges       map[string]*ExternalDocChanges           `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	TagChanges               map[string]*TagChanges                   `json:"tags,omitempty" yaml:"tags,omitempty"`
	OperationTraitChanges    map[string]*OperationTraitChanges        `json:"operationTraits,omitempty" yaml:"operationTraits,omitempty"`
	MessageTraitChanges      map[string]*MessageTraitChanges          `json:"messageTraits,omitempty" yaml:"messageTraits,omitempty"`
	ServerBindingsChanges    map[string]*ServerBindingsChanges        `json:"serverBindings,omitempty" yaml:"serverBindings,omitempty"`
	ChannelBindingsChanges   map[string]*ChannelBindingsChanges       `json:"channelBindings,omitempty" yaml:"channelBindings,omitempty"`
	OperationBindingsChanges map[string]*OperationBindingsChanges     `json:"operationBindings,omitempty" yaml:"operationBindings,omitempty"`
	MessageBindingsChanges   map[string]*MessageBindingsChanges       `json:"messageBindings,omitempty" yaml:"messageBindings,omitempty"`
	ExtensionChanges         *ExtensionChanges                        `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Components objects.
func (c *ComponentsChanges) GetAllChanges() []*Change {
	if c == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, c.Changes...)
	changes = appendMapChanges(changes, c.SchemaChanges)
	changes = appendMapChanges(changes, c.ServerChanges)
	changes = appendMapChanges(changes, c.ChannelChanges)
	changes = appendMapChanges(changes, c.OperationChanges)
	changes = appendMapChanges(changes, c.MessageChanges)
	changes = appendMapChanges(changes, c.SecuritySchemeChanges)
	changes = appendMapChanges(changes, c.ServerVariableChanges)
	changes = appendMapChanges(changes, c.ParameterChanges)
	changes = appendMapChanges(changes, c.CorrelationIDChanges)
	changes = appendMapChanges(changes, c.ReplyChanges)
	changes = appendMapChanges(changes, c.ReplyAddressChanges)
	changes = appendMapChanges(changes, c.ExternalDocChanges)
	changes = appendMapChanges(changes, c.TagChanges)
	changes = appendMapChanges(changes, c.OperationTraitChanges)
	changes = appendMapChanges(changes, c.MessageTraitChanges)
	changes = appendMapChanges(changes, c.ServerBindingsChanges)
	changes = appendMapChanges(changes, c.ChannelBindingsChanges)
	changes = appendMapChanges(changes, c.OperationBindingsChanges)
	changes = appendMapChanges(changes, c.MessageBindingsChanges)
	if c.ExtensionChanges != nil {
		changes = append(changes, c.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (c *ComponentsChanges) TotalChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalChanges()
	t += countMapChanges(c.SchemaChanges)
	t += countMapChanges(c.ServerChanges)
	t += countMapChanges(c.ChannelChanges)
	t += countMapChanges(c.OperationChanges)
	t += countMapChanges(c.MessageChanges)
	t += countMapChanges(c.SecuritySchemeChanges)
	t += countMapChanges(c.ServerVariableChanges)
	t += countMapChanges(c.ParameterChanges)
	t += countMapChanges(c.CorrelationIDChanges)
	t += countMapChanges(c.ReplyChanges)
	t += countMapChanges(c.ReplyAddressChanges)
	t += countMapChanges(c.ExternalDocChanges)
	t += countMapChanges(c.TagChanges)
	t += countMapChanges(c.OperationTraitChanges)
	t += countMapChanges(c.MessageTraitChanges)
	t += countMapChanges(c.ServerBindingsChanges)
	t += countMapChanges(c.ChannelBindingsChanges)
	t += countMapChanges(c.OperationBindingsChanges)
	t += countMapChanges(c.MessageBindingsChanges)
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalChanges()
	}
	return t
}

// TotalBreakingChanges returns the number of breaking changes made.
func (c *ComponentsChanges) TotalBreakingChanges() int {
	if c == nil {
		return 0
	}
	t := c.PropertyChanges.TotalBreakingChanges()
	t += countMapBreakingChanges(c.SchemaChanges)
	t += countMapBreakingChanges(c.ServerChanges)
	t += countMapBreakingChanges(c.ChannelChanges)
	t += countMapBreakingChanges(c.OperationChanges)
	t += countMapBreakingChanges(c.MessageChanges)
	t += countMapBreakingChanges(c.SecuritySchemeChanges)
	t += countMapBreakingChanges(c.ServerVariableChanges)
	t += countMapBreakingChanges(c.ParameterChanges)
	t += countMapBreakingChanges(c.CorrelationIDChanges)
	t += countMapBreakingChanges(c.ReplyChanges)
	t += countMapBreakingChanges(c.ReplyAddressChanges)
	t += countMapBreakingChanges(c.ExternalDocChanges)
	t += countMapBreakingChanges(c.TagChanges)
	t += countMapBreakingChanges(c.OperationTraitChanges)
	t += countMapBreakingChanges(c.MessageTraitChanges)
	t += countMapBreakingChanges(c.ServerBindingsChanges)
	t += countMapBreakingChanges(c.ChannelBindingsChanges)
	t += countMapBreakingChanges(c.OperationBindingsChanges)
	t += countMapBreakingChanges(c.MessageBindingsChanges)
	if c.ExtensionChanges != nil {
		t += c.ExtensionChanges.TotalBreakingChanges()
	}
	return t
}

// assignIfChanged assigns a comparison result map to target only when it has entries.
func assignIfChanged[T any](target *map[string]T, result map[string]T) {
	if len(result) > 0 {
		*target = result
	}
}

// CompareComponents compares two AsyncAPI Components objects and returns a pointer to
// ComponentsChanges, or nil if nothing changed.
//
// Schemas are libopenapi base objects, so their comparison is delegated to libopenapi's
// what-changed schema comparator.
func CompareComponents(l, r *lowasync.Components, configs ...*BreakingRulesConfig) *ComponentsChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	cc := new(ComponentsChanges)

	assignIfChanged(&cc.SchemaChanges, CheckMapForChangesWithRules(l.Schemas.Value, r.Schemas.Value,
		&changes, lowasync.SchemasLabel, wcmodel.CompareSchemas, CompComponents, PropSchemas, config))
	assignIfChanged(&cc.ServerChanges, CheckMapForChangesWithRules(l.Servers.Value, r.Servers.Value,
		&changes, lowasync.ServersLabel, configuredCompare(config, CompareServers), CompComponents, PropServers, config))
	assignIfChanged(&cc.ChannelChanges, CheckMapForChangesWithRules(l.Channels.Value, r.Channels.Value,
		&changes, lowasync.ChannelsLabel, configuredCompare(config, CompareChannels), CompComponents, PropChannels, config))
	assignIfChanged(&cc.OperationChanges, CheckMapForChangesWithRules(l.Operations.Value, r.Operations.Value,
		&changes, lowasync.OperationsLabel, configuredCompare(config, CompareOperations), CompComponents, PropOperations, config))
	assignIfChanged(&cc.MessageChanges, CheckMapForChangesWithRules(l.Messages.Value, r.Messages.Value,
		&changes, lowasync.MessagesLabel, configuredCompare(config, CompareMessages), CompComponents, PropMessages, config))
	assignIfChanged(&cc.SecuritySchemeChanges, CheckMapForChangesWithRules(l.SecuritySchemes.Value, r.SecuritySchemes.Value,
		&changes, lowasync.SecuritySchemesLabel, configuredCompare(config, CompareSecuritySchemes), CompComponents, PropSecuritySchemes, config))
	assignIfChanged(&cc.ServerVariableChanges, CheckMapForChangesWithRules(l.ServerVariables.Value, r.ServerVariables.Value,
		&changes, lowasync.ServerVariablesLabel, configuredCompare(config, CompareServerVariables), CompComponents, PropServerVariables, config))
	assignIfChanged(&cc.ParameterChanges, CheckMapForChangesWithRules(l.Parameters.Value, r.Parameters.Value,
		&changes, lowasync.ParametersLabel, configuredCompare(config, CompareParameters), CompComponents, PropParameters, config))
	assignIfChanged(&cc.CorrelationIDChanges, CheckMapForChangesWithRules(l.CorrelationIDs.Value, r.CorrelationIDs.Value,
		&changes, lowasync.CorrelationIDsLabel, configuredCompare(config, CompareCorrelationID), CompComponents, PropCorrelationIDs, config))
	assignIfChanged(&cc.ReplyChanges, CheckMapForChangesWithRules(l.Replies.Value, r.Replies.Value,
		&changes, lowasync.RepliesLabel, configuredCompare(config, CompareOperationReply), CompComponents, PropReplies, config))
	assignIfChanged(&cc.ReplyAddressChanges, CheckMapForChangesWithRules(l.ReplyAddresses.Value, r.ReplyAddresses.Value,
		&changes, lowasync.ReplyAddressesLabel, configuredCompare(config, CompareOperationReplyAddress), CompComponents, PropReplyAddresses, config))
	assignIfChanged(&cc.ExternalDocChanges, CheckMapForChangesWithRules(l.ExternalDocs.Value, r.ExternalDocs.Value,
		&changes, lowasync.ExternalDocsLabel, configuredCompare(config, CompareExternalDocs), CompComponents, PropExternalDocs, config))
	assignIfChanged(&cc.TagChanges, CheckMapForChangesWithRules(l.Tags.Value, r.Tags.Value,
		&changes, lowasync.TagsLabel, configuredCompare(config, CompareTags), CompComponents, PropTags, config))
	assignIfChanged(&cc.OperationTraitChanges, CheckMapForChangesWithRules(l.OperationTraits.Value, r.OperationTraits.Value,
		&changes, lowasync.OperationTraitsLabel, configuredCompare(config, CompareOperationTraits), CompComponents, PropOperationTraits, config))
	assignIfChanged(&cc.MessageTraitChanges, CheckMapForChangesWithRules(l.MessageTraits.Value, r.MessageTraits.Value,
		&changes, lowasync.MessageTraitsLabel, configuredCompare(config, CompareMessageTraits), CompComponents, PropMessageTraits, config))
	assignIfChanged(&cc.ServerBindingsChanges, CheckMapForChangesWithRules(l.ServerBindings.Value, r.ServerBindings.Value,
		&changes, lowasync.ServerBindingsLabel, configuredCompare(config, CompareServerBindings), CompComponents, PropServerBindings, config))
	assignIfChanged(&cc.ChannelBindingsChanges, CheckMapForChangesWithRules(l.ChannelBindings.Value, r.ChannelBindings.Value,
		&changes, lowasync.ChannelBindingsLabel, configuredCompare(config, CompareChannelBindings), CompComponents, PropChannelBindings, config))
	assignIfChanged(&cc.OperationBindingsChanges, CheckMapForChangesWithRules(l.OperationBindings.Value, r.OperationBindings.Value,
		&changes, lowasync.OperationBindingsLabel, configuredCompare(config, CompareOperationBindings), CompComponents, PropOperationBindings, config))
	assignIfChanged(&cc.MessageBindingsChanges, CheckMapForChangesWithRules(l.MessageBindings.Value, r.MessageBindings.Value,
		&changes, lowasync.MessageBindingsLabel, configuredCompare(config, CompareMessageBindings), CompComponents, PropMessageBindings, config))

	cc.ExtensionChanges = CheckExtensions(l, r)
	cc.PropertyChanges = NewPropertyChanges(changes)
	if cc.TotalChanges() <= 0 {
		return nil
	}
	return cc
}
