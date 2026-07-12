// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
)

// DocumentChanges represents all the changes made between two AsyncAPI documents.
type DocumentChanges struct {
	*PropertyChanges
	InfoChanges       *InfoChanges                 `json:"info,omitempty" yaml:"info,omitempty"`
	ServerChanges     map[string]*ServerChanges    `json:"servers,omitempty" yaml:"servers,omitempty"`
	ChannelChanges    map[string]*ChannelChanges   `json:"channels,omitempty" yaml:"channels,omitempty"`
	OperationChanges  map[string]*OperationChanges `json:"operations,omitempty" yaml:"operations,omitempty"`
	ComponentsChanges *ComponentsChanges           `json:"components,omitempty" yaml:"components,omitempty"`
	ExtensionChanges  *ExtensionChanges            `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between AsyncAPI documents.
func (d *DocumentChanges) GetAllChanges() []*Change {
	if d == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, d.Changes...)
	if d.InfoChanges != nil {
		changes = append(changes, d.InfoChanges.GetAllChanges()...)
	}
	changes = appendMapChanges(changes, d.ServerChanges)
	changes = appendMapChanges(changes, d.ChannelChanges)
	changes = appendMapChanges(changes, d.OperationChanges)
	if d.ComponentsChanges != nil {
		changes = append(changes, d.ComponentsChanges.GetAllChanges()...)
	}
	if d.ExtensionChanges != nil {
		changes = append(changes, d.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (d *DocumentChanges) TotalChanges() int {
	if d == nil {
		return 0
	}
	c := d.PropertyChanges.TotalChanges()
	if d.InfoChanges != nil {
		c += d.InfoChanges.TotalChanges()
	}
	c += countMapChanges(d.ServerChanges)
	c += countMapChanges(d.ChannelChanges)
	c += countMapChanges(d.OperationChanges)
	if d.ComponentsChanges != nil {
		c += d.ComponentsChanges.TotalChanges()
	}
	if d.ExtensionChanges != nil {
		c += d.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (d *DocumentChanges) TotalBreakingChanges() int {
	if d == nil {
		return 0
	}
	c := d.PropertyChanges.TotalBreakingChanges()
	if d.InfoChanges != nil {
		c += d.InfoChanges.TotalBreakingChanges()
	}
	c += countMapBreakingChanges(d.ServerChanges)
	c += countMapBreakingChanges(d.ChannelChanges)
	c += countMapBreakingChanges(d.OperationChanges)
	if d.ComponentsChanges != nil {
		c += d.ComponentsChanges.TotalBreakingChanges()
	}
	if d.ExtensionChanges != nil {
		c += d.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareDocuments compares two AsyncAPI documents and returns a pointer to
// DocumentChanges describing everything that changed, or nil if nothing changed.
func CompareDocuments(l, r *lowasync.AsyncAPI) *DocumentChanges {
	return compareDocuments(l, r, GetActiveBreakingRulesConfig())
}

// CompareDocumentsWithConfig compares two AsyncAPI documents using config for this call
// only. A nil config selects the default rules, independent of any process-global active
// configuration. Custom-config comparisons are serialized so their rules cannot cross-talk
// with concurrent comparisons or global configuration changes.
func CompareDocumentsWithConfig(l, r *lowasync.AsyncAPI, config *BreakingRulesConfig) *DocumentChanges {
	if config == nil {
		config = GenerateDefaultBreakingRules()
	}
	return compareDocuments(l, r, config)
}

func compareDocuments(l, r *lowasync.AsyncAPI, config *BreakingRulesConfig) *DocumentChanges {
	if l == nil || r == nil {
		return nil
	}

	// reset the schema quick-hash map and the hash cache so stale entries from a
	// previous comparison cannot collide with this one.
	base.SchemaQuickHashMap.Clear()
	low.ClearHashCache()

	var changes []*Change
	dc := new(DocumentChanges)

	props := []*PropertyCheck{
		NewPropertyCheck(CompAsyncAPI, "",
			l.AsyncAPI.ValueNode, r.AsyncAPI.ValueNode, lowasync.AsyncAPILabel, &changes, l, r, config),
		NewPropertyCheck(CompID, "",
			l.ID.ValueNode, r.ID.ValueNode, lowasync.IDLabel, &changes, l, r, config),
		NewPropertyCheck(CompDefaultContentType, "",
			l.DefaultContentType.ValueNode, r.DefaultContentType.ValueNode,
			lowasync.DefaultContentTypeLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	// info. adding or removing the whole info or components container sits outside the
	// configurable rules system (nested rule structs only register component.property
	// keys), so the flags are fixed: additions are never breaking; removals are breaking
	// because info is required by the specification and components removal severs refs.
	if l.Info.Value != nil && r.Info.Value != nil {
		dc.InfoChanges = CompareInfo(l.Info.Value, r.Info.Value, config)
	} else {
		if l.Info.Value == nil && r.Info.Value != nil {
			CreateChange(&changes, ObjectAdded, lowasync.InfoLabel,
				nil, r.Info.ValueNode, false, nil, r.Info.Value)
		}
		if l.Info.Value != nil && r.Info.Value == nil {
			CreateChange(&changes, ObjectRemoved, lowasync.InfoLabel,
				l.Info.ValueNode, nil, true, l.Info.Value, nil)
		}
	}

	// servers, channels, operations
	assignIfChanged(&dc.ServerChanges, CheckMapForChangesWithRules(l.Servers.Value, r.Servers.Value,
		&changes, lowasync.ServersLabel, configuredCompare(config, CompareServers), CompServers, "", config))
	assignIfChanged(&dc.ChannelChanges, CheckMapForChangesWithRules(l.Channels.Value, r.Channels.Value,
		&changes, lowasync.ChannelsLabel, configuredCompare(config, CompareChannels), CompChannels, "", config))
	assignIfChanged(&dc.OperationChanges, CheckMapForChangesWithRules(l.Operations.Value, r.Operations.Value,
		&changes, lowasync.OperationsLabel, configuredCompare(config, CompareOperations), CompOperations, "", config))

	// components
	if l.Components.Value != nil && r.Components.Value != nil {
		dc.ComponentsChanges = CompareComponents(l.Components.Value, r.Components.Value, config)
	} else {
		if l.Components.Value == nil && r.Components.Value != nil {
			CreateChange(&changes, ObjectAdded, lowasync.ComponentsLabel,
				nil, r.Components.ValueNode, false, nil, r.Components.Value)
		}
		if l.Components.Value != nil && r.Components.Value == nil {
			CreateChange(&changes, ObjectRemoved, lowasync.ComponentsLabel,
				l.Components.ValueNode, nil, true, l.Components.Value, nil)
		}
	}

	dc.ExtensionChanges = CheckExtensions(l, r)
	dc.PropertyChanges = NewPropertyChanges(changes)
	if dc.TotalChanges() <= 0 {
		return nil
	}
	return dc
}
