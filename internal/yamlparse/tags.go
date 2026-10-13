// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package yamlparse

const (
	// MappingTag identifies a YAML mapping node.
	MappingTag = "!!map"
	// SequenceTag identifies a YAML sequence node.
	SequenceTag = "!!seq"
	// StringTag identifies a YAML string scalar.
	StringTag = "!!str"
	// IntegerTag identifies a YAML integer scalar.
	IntegerTag = "!!int"
	// FloatTag identifies a YAML floating-point scalar.
	FloatTag = "!!float"
	// BooleanTag identifies a YAML boolean scalar.
	BooleanTag = "!!bool"
	// NullTag identifies a YAML null scalar.
	NullTag = "!!null"
)
