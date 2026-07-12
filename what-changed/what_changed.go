// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package what_changed determines what changed between two AsyncAPI documents.
//
// The package compares the low-level models of two parsed documents and produces a
// hierarchical tree of changes — additions, removals and modifications — each classified
// as breaking or non-breaking with full YAML line and column context. Breaking
// classification is driven by a configurable, AsyncAPI shaped rules system in the model
// sub-package.
package what_changed

import (
	libasyncapi "github.com/pb33f/libasyncapi"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/what-changed/model"
)

// CompareAsyncAPIDocuments compares two low-level AsyncAPI documents and returns a
// pointer to model.DocumentChanges describing everything that changed, or nil if
// nothing changed.
func CompareAsyncAPIDocuments(original, updated *lowasync.AsyncAPI) *model.DocumentChanges {
	return model.CompareDocuments(original, updated)
}

// CompareAsyncAPIDocumentsWithConfig compares two low-level AsyncAPI documents using
// config for this call only. Pass nil to use the default breaking rules.
func CompareAsyncAPIDocumentsWithConfig(original, updated *lowasync.AsyncAPI,
	config *model.BreakingRulesConfig,
) *model.DocumentChanges {
	return model.CompareDocumentsWithConfig(original, updated, config)
}

// CompareDocuments compares two parsed AsyncAPI documents and returns a pointer to
// model.DocumentChanges describing everything that changed, or nil if nothing changed.
func CompareDocuments(original, updated libasyncapi.Document) *model.DocumentChanges {
	if original == nil || updated == nil {
		return nil
	}
	return model.CompareDocuments(original.GoLow(), updated.GoLow())
}

// CompareDocumentsWithConfig compares two parsed AsyncAPI documents using config for this
// call only. Pass nil to use the default breaking rules.
func CompareDocumentsWithConfig(original, updated libasyncapi.Document,
	config *model.BreakingRulesConfig,
) *model.DocumentChanges {
	if original == nil || updated == nil {
		return nil
	}
	return model.CompareDocumentsWithConfig(original.GoLow(), updated.GoLow(), config)
}
