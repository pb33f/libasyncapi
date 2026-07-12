// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package model contains all the comparison logic and change types used to determine
// what changed between two AsyncAPI documents. The package builds on top of
// libopenapi's what-changed model, reusing its core change types and comparison
// primitives, while providing AsyncAPI specific comparators and an AsyncAPI shaped
// breaking-change rules system.
//
// Breaking-rules boundary: every AsyncAPI object (including Contact and License) is
// classified against THIS package's active breaking-rules configuration. The single
// deliberate exception is JSON Schema comparison: schema-bearing properties (message
// payloads/headers, binding query/headers/key fields, components.schemas) delegate to
// libopenapi's CompareSchemas, which applies libopenapi's own schema rules — including
// any custom configuration installed with libopenapi's SetActiveBreakingRulesConfig.
// Each exported AsyncAPI comparator also accepts an optional BreakingRulesConfig. When
// supplied, that immutable configuration is propagated through the complete nested
// comparison without changing this package's process-global default.
package model

import (
	"github.com/pb33f/libopenapi/datamodel/low"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// Core change types are aliased from libopenapi's what-changed model so a single set
// of types flows through both the delegated (OpenAPI base object) comparisons and the
// AsyncAPI specific ones.
type (
	// Change represents a change between two AsyncAPI documents.
	Change = wcmodel.Change

	// PropertyChanges holds a slice of Change pointers.
	PropertyChanges = wcmodel.PropertyChanges

	// PropertyCheck is used by checking mechanisms to define a property to be checked.
	PropertyCheck = wcmodel.PropertyCheck

	// ChangeContext holds a reference to the line and column positions of original and new change.
	ChangeContext = wcmodel.ChangeContext

	// ExtensionChanges represents any changes to custom extensions defined for an AsyncAPI object.
	ExtensionChanges = wcmodel.ExtensionChanges

	// BreakingChangeRule holds the breaking status for a property's change types.
	// nil values mean "use default" - only set values override the defaults.
	BreakingChangeRule = wcmodel.BreakingChangeRule
)

// Change type constants, aliased from libopenapi's what-changed model.
const (
	// Modified means the value was changed.
	Modified = wcmodel.Modified

	// PropertyAdded means a new property to an object was added.
	PropertyAdded = wcmodel.PropertyAdded

	// ObjectAdded means a new object was added.
	ObjectAdded = wcmodel.ObjectAdded

	// ObjectRemoved means an object was removed.
	ObjectRemoved = wcmodel.ObjectRemoved

	// PropertyRemoved means a property of an object was removed.
	PropertyRemoved = wcmodel.PropertyRemoved
)

// ChangeType constants for breaking rules lookups.
const (
	ChangeTypeAdded    = wcmodel.ChangeTypeAdded
	ChangeTypeModified = wcmodel.ChangeTypeModified
	ChangeTypeRemoved  = wcmodel.ChangeTypeRemoved
)

// Non-generic comparison primitives re-exported from libopenapi's what-changed model.
// These are config free (they take an explicit breaking flag or none at all), so they
// behave identically regardless of which breaking-rules configuration is active.
var (
	// NewPropertyChanges creates a new PropertyChanges instance from a set of changes.
	NewPropertyChanges = wcmodel.NewPropertyChanges

	// CreateChange creates a new Change of type T with the provided types
	CreateChange = wcmodel.CreateChange

	// CreateChangeWithEncoding creates a new Change with YAML-encoded original/new values,
	// used for complex types like extensions.
	CreateChangeWithEncoding = wcmodel.CreateChangeWithEncoding

	// CompareExtensions compares a left and right map of extension nodes for any changes.
	CompareExtensions = wcmodel.CompareExtensions

	// CountBreakingChanges counts the number of changes in a slice that are breaking.
	CountBreakingChanges = wcmodel.CountBreakingChanges
)

// CheckExtensions unpacks a left and right model that contains extensions, compares
// them and returns a pointer to ExtensionChanges. If nothing changed, nil is returned.
func CheckExtensions[T low.HasExtensions[T]](l, r T) *ExtensionChanges {
	return wcmodel.CheckExtensions(l, r)
}
