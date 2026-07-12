// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package reports

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/what-changed/model"
)

// Changed provides a simple wrapper for changed counts.
type Changed struct {
	Total    int `json:"totalChanges"`
	Breaking int `json:"breakingChanges"`
}

// OverallReport provides a document level overview of all changes to an AsyncAPI doc.
type OverallReport struct {
	ChangeReport map[string]*Changed `json:"overallSummaryReport"`
}

// CreateOverallReport creates a high level report for all top level changes (with deep counts).
func CreateOverallReport(changes *model.DocumentChanges) *OverallReport {
	changedReport := make(map[string]*Changed)
	if changes == nil {
		return &OverallReport{
			ChangeReport: changedReport,
		}
	}

	mergeRootPropertyChanges(changedReport, changes.PropertyChanges)
	if changes.InfoChanges != nil {
		mergeChangedModel(changedReport, lowasync.InfoLabel, createChangedModel(changes.InfoChanges))
	}
	if changes.ServerChanges != nil {
		mergeChangedModel(changedReport, lowasync.ServersLabel,
			createChangedModelFromMap(changes.ServerChanges))
	}
	if changes.ChannelChanges != nil {
		mergeChangedModel(changedReport, lowasync.ChannelsLabel,
			createChangedModelFromMap(changes.ChannelChanges))
	}
	if changes.OperationChanges != nil {
		mergeChangedModel(changedReport, lowasync.OperationsLabel,
			createChangedModelFromMap(changes.OperationChanges))
	}
	if changes.ComponentsChanges != nil {
		mergeChangedModel(changedReport, lowasync.ComponentsLabel, createChangedModel(changes.ComponentsChanges))
	}
	if changes.ExtensionChanges != nil {
		mergeChangedModel(changedReport, extensionsLabel, createChangedModel(changes.ExtensionChanges))
	}
	return &OverallReport{
		ChangeReport: changedReport,
	}
}

// extensionsLabel buckets root-level specification extension changes in the report,
// keeping the report totals equal to DocumentChanges.TotalChanges().
const extensionsLabel = "extensions"

func mergeRootPropertyChanges(changedReport map[string]*Changed, propertyChanges *model.PropertyChanges) {
	if propertyChanges == nil {
		return
	}
	for _, change := range propertyChanges.GetPropertyChanges() {
		if change == nil || change.Property == "" {
			continue
		}
		changed := getOrCreateChanged(changedReport, change.Property)
		changed.Total++
		if change.Breaking {
			changed.Breaking++
		}
	}
}

func mergeChangedModel(changedReport map[string]*Changed, label string, next *Changed) {
	if next == nil {
		return
	}
	changed := getOrCreateChanged(changedReport, label)
	changed.Total += next.Total
	changed.Breaking += next.Breaking
}

func getOrCreateChanged(changedReport map[string]*Changed, label string) *Changed {
	if changed, ok := changedReport[label]; ok {
		return changed
	}
	changed := &Changed{}
	changedReport[label] = changed
	return changed
}

func createChangedModel(ch HasChanges) *Changed {
	return &Changed{ch.TotalChanges(), ch.TotalBreakingChanges()}
}

func createChangedModelFromMap[T HasChanges](m map[string]T) *Changed {
	t := 0
	b := 0
	for k := range m {
		t += m[k].TotalChanges()
		b += m[k].TotalBreakingChanges()
	}
	return &Changed{t, b}
}
