// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	"github.com/pb33f/libopenapi/orderedmap"
)

// AsyncAPI represents a high-level AsyncAPI 3.0 document.
//
// This is the root object of the AsyncAPI document.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#A2SObject
type AsyncAPI struct {
	AsyncAPI           string                              `json:"asyncapi,omitempty" yaml:"asyncapi,omitempty"`
	ID                 string                              `json:"id,omitempty" yaml:"id,omitempty"`
	Info               *Info                               `json:"info,omitempty" yaml:"info,omitempty"`
	Servers            *orderedmap.Map[string, *Server]    `json:"servers,omitempty" yaml:"servers,omitempty"`
	DefaultContentType string                              `json:"defaultContentType,omitempty" yaml:"defaultContentType,omitempty"`
	Channels           *orderedmap.Map[string, *Channel]   `json:"channels,omitempty" yaml:"channels,omitempty"`
	Operations         *orderedmap.Map[string, *Operation] `json:"operations,omitempty" yaml:"operations,omitempty"`
	Components         *Components                         `json:"components,omitempty" yaml:"components,omitempty"`
	Extensions         *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low                *lowasync.AsyncAPI
}

// NewAsyncAPI creates a new high-level AsyncAPI instance from a low-level one.
func NewAsyncAPI(doc *lowasync.AsyncAPI) *AsyncAPI {
	a := new(AsyncAPI)
	a.low = doc
	a.AsyncAPI = doc.AsyncAPI.Value
	a.ID = doc.ID.Value
	a.DefaultContentType = doc.DefaultContentType.Value

	if !doc.Info.IsEmpty() {
		a.Info = NewInfo(doc.Info.Value)
	}
	if doc.Servers.Value != nil {
		a.Servers = orderedmap.New[string, *Server]()
		for k, v := range doc.Servers.Value.FromOldest() {
			a.Servers.Set(k.Value, NewServer(v.Value))
		}
	}
	if doc.Channels.Value != nil {
		a.Channels = orderedmap.New[string, *Channel]()
		for k, v := range doc.Channels.Value.FromOldest() {
			a.Channels.Set(k.Value, NewChannel(v.Value))
		}
	}
	if doc.Operations.Value != nil {
		a.Operations = orderedmap.New[string, *Operation]()
		for k, v := range doc.Operations.Value.FromOldest() {
			a.Operations.Set(k.Value, NewOperation(v.Value))
		}
	}
	if !doc.Components.IsEmpty() {
		a.Components = NewComponents(doc.Components.Value)
	}
	if orderedmap.Len(doc.Extensions) > 0 {
		a.Extensions = high.ExtractExtensions(doc.Extensions)
	}
	return a
}

// GoLow returns the low-level AsyncAPI instance.
func (a *AsyncAPI) GoLow() *lowasync.AsyncAPI {
	return a.low
}

// GoLowUntyped returns the low-level AsyncAPI instance with no type.
func (a *AsyncAPI) GoLowUntyped() any {
	return a.low
}

// Render will return a YAML representation of the AsyncAPI document as a byte slice.
func (a *AsyncAPI) Render() ([]byte, error) {
	return yaml.Marshal(a)
}

// MarshalYAML will create a ready to render YAML representation of the AsyncAPI document.
func (a *AsyncAPI) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(a, a.low)
	return nb.Render(), nil
}
