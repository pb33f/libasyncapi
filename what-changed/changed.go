// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package what_changed

import "github.com/pb33f/libasyncapi/what-changed/model"

// Changed represents an object that was changed.
type Changed interface {
	// GetAllChanges returns all top level changes made to properties in this object,
	// including all children.
	GetAllChanges() []*model.Change

	// TotalChanges returns a count of all changes made on the object, including all children.
	TotalChanges() int

	// TotalBreakingChanges returns a count of all breaking changes on this object,
	// including all children.
	TotalBreakingChanges() int
}
