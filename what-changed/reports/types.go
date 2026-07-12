// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package reports provides summarized, flattened views over the hierarchical change
// trees produced by the what-changed model.
package reports

// HasChanges represents a change model that provides a total change count and a breaking change count.
type HasChanges interface {
	// TotalChanges represents number of all changes found.
	TotalChanges() int

	// TotalBreakingChanges represents the number of contract breaking changes only.
	TotalBreakingChanges() int
}
