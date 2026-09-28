// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"github.com/pb33f/go-yaml"
	highasync "github.com/pb33f/libasyncapi/datamodel/high/asyncapi"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/index"
)

// Document represents a parsed AsyncAPI document that can be rendered into a model or serialized.
type Document interface {
	// GetVersion will return the exact version of the AsyncAPI specification set for the document.
	GetVersion() string

	// GetSpecInfo will return the SpecInfo instance that contains all specification information.
	GetSpecInfo() *SpecInfo

	// Model returns the high-level AsyncAPI model.
	// Returns nil if parsing completely failed.
	Model() *highasync.AsyncAPI

	// GoLow returns the low-level AsyncAPI model.
	// This provides access to line/column numbers and raw YAML nodes.
	GoLow() *lowasync.AsyncAPI

	// Index returns the reference index for this document.
	Index() *index.SpecIndex

	// Rolodex returns the file management system for multi-file specs.
	Rolodex() *index.Rolodex

	// RootNode returns the original YAML tree (not normalized).
	// This provides raw access for tooling that needs the original structure.
	RootNode() *yaml.Node

	// Serialize marshals the original YAML root node back to bytes.
	// This preserves the original structure and formatting.
	// For model-based rendering, use Render() instead.
	Serialize() ([]byte, error)

	// Render will return a YAML representation of the high-level model.
	Render() ([]byte, error)

	// Errors returns all errors accumulated during parsing.
	// The document may still contain a partial model even if errors occurred.
	Errors() []error

	// IsPartial returns true if errors occurred during parsing.
	// When true, the model may be incomplete but still usable.
	IsPartial() bool
}

type document struct {
	version   string
	specInfo  *SpecInfo
	highModel *highasync.AsyncAPI
	lowModel  *lowasync.AsyncAPI
	specIndex *index.SpecIndex
	rolodex   *index.Rolodex
	rootNode  *yaml.Node
	errors    []error
	config    *DocumentConfiguration
}

func (d *document) GetVersion() string {
	return d.version
}

func (d *document) GetSpecInfo() *SpecInfo {
	return d.specInfo
}

func (d *document) Model() *highasync.AsyncAPI {
	return d.highModel
}

func (d *document) GoLow() *lowasync.AsyncAPI {
	return d.lowModel
}

func (d *document) Index() *index.SpecIndex {
	return d.specIndex
}

func (d *document) Rolodex() *index.Rolodex {
	return d.rolodex
}

func (d *document) RootNode() *yaml.Node {
	return d.rootNode
}

func (d *document) Serialize() ([]byte, error) {
	if d.rootNode == nil {
		return nil, ErrInvalidYAML
	}
	return yaml.Marshal(d.rootNode)
}

func (d *document) Render() ([]byte, error) {
	if d.highModel == nil {
		return nil, ErrInvalidYAML
	}
	return d.highModel.Render()
}

func (d *document) Errors() []error {
	return d.errors
}

func (d *document) IsPartial() bool {
	return len(d.errors) > 0
}

func (d *document) addError(err error) {
	d.errors = append(d.errors, err)
}
