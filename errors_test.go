// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"testing"

	"github.com/pb33f/testify/assert"
)

func TestParseError_Error_WithPath(t *testing.T) {
	err := &ParseError{
		Message: "invalid field",
		Line:    10,
		Column:  5,
		Path:    "#/channels/myChannel",
	}
	assert.Equal(t, "invalid field at line 10, column 5 (#/channels/myChannel)", err.Error())
}

func TestParseError_Error_WithoutPath(t *testing.T) {
	err := &ParseError{
		Message: "syntax error",
		Line:    1,
		Column:  1,
		Path:    "",
	}
	assert.Equal(t, "syntax error at line 1, column 1", err.Error())
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("test message", 5, 10, "#/info/title")
	assert.Equal(t, "test message", err.Message)
	assert.Equal(t, 5, err.Line)
	assert.Equal(t, 10, err.Column)
	assert.Equal(t, "#/info/title", err.Path)
}

func TestValidationError_Error_WithPath(t *testing.T) {
	err := &ValidationError{
		Message: "required field missing",
		Path:    "#/info/version",
	}
	assert.Equal(t, "required field missing at #/info/version", err.Error())
}

func TestValidationError_Error_WithoutPath(t *testing.T) {
	err := &ValidationError{
		Message: "general error",
		Path:    "",
	}
	assert.Equal(t, "general error", err.Error())
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("missing required field", "#/operations/sendUser")
	assert.Equal(t, "missing required field", err.Message)
	assert.Equal(t, "#/operations/sendUser", err.Path)
}

func TestReferenceError_Error(t *testing.T) {
	err := &ReferenceError{
		Reference: "#/components/schemas/User",
		Message:   "reference not found",
		Line:      25,
		Column:    12,
	}
	assert.Equal(t, "reference error for '#/components/schemas/User': reference not found at line 25, column 12", err.Error())
}

func TestNewReferenceError(t *testing.T) {
	err := NewReferenceError("#/components/messages/UserCreated", "target not found", 30, 8)
	assert.Equal(t, "#/components/messages/UserCreated", err.Reference)
	assert.Equal(t, "target not found", err.Message)
	assert.Equal(t, 30, err.Line)
	assert.Equal(t, 8, err.Column)
}

func TestCircularReferenceError_Error(t *testing.T) {
	err := &CircularReferenceError{
		Path:  "#/components/schemas/Node",
		Chain: []string{"#/components/schemas/Node", "#/components/schemas/Child", "#/components/schemas/Node"},
	}
	assert.Equal(t, "circular reference detected at #/components/schemas/Node", err.Error())
}

func TestNewCircularReferenceError(t *testing.T) {
	chain := []string{"A", "B", "C", "A"}
	err := NewCircularReferenceError("#/schemas/A", chain)
	assert.Equal(t, "#/schemas/A", err.Path)
	assert.Equal(t, chain, err.Chain)
}
