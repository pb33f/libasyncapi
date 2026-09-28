// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Info represents a high-level AsyncAPI 3.0 Info object, backed by a low-level one.
//
// Provides metadata about the API. The metadata can be used by the clients if needed.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#infoObject
type Info struct {
	Title          string                              `json:"title,omitempty" yaml:"title,omitempty"`
	Version        string                              `json:"version,omitempty" yaml:"version,omitempty"`
	Description    string                              `json:"description,omitempty" yaml:"description,omitempty"`
	TermsOfService string                              `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
	Contact        *highbase.Contact                   `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *highbase.License                   `json:"license,omitempty" yaml:"license,omitempty"`
	Tags           []*Tag                              `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs   *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Extensions     *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low            *lowasync.Info
}

// NewInfo creates a new high-level Info instance from a low-level one.
func NewInfo(info *lowasync.Info) *Info {
	i := new(Info)
	i.low = info
	i.Title = info.Title.Value
	i.Version = info.Version.Value
	i.Description = info.Description.Value
	i.TermsOfService = info.TermsOfService.Value

	if !info.Contact.IsEmpty() {
		i.Contact = highbase.NewContact(info.Contact.Value)
	}
	if !info.License.IsEmpty() {
		i.License = highbase.NewLicense(info.License.Value)
	}
	if info.Tags.Value != nil {
		for pair := range info.Tags.Value.ValuesFromOldest() {
			i.Tags = append(i.Tags, NewTag(pair.Value))
		}
	}
	if !info.ExternalDocs.IsEmpty() {
		i.ExternalDocs = NewExternalDoc(info.ExternalDocs.Value)
	}
	if orderedmap.Len(info.Extensions) > 0 {
		i.Extensions = high.ExtractExtensions(info.Extensions)
	}
	return i
}

// GoLow returns the low-level Info instance that was used to create the high-level one.
func (i *Info) GoLow() *lowasync.Info {
	return i.low
}

// GoLowUntyped returns the low-level Info instance with no type.
func (i *Info) GoLowUntyped() any {
	return i.low
}

// Render will return a YAML representation of the Info object as a byte slice.
func (i *Info) Render() ([]byte, error) {
	return yaml.Marshal(i)
}

// MarshalYAML will create a ready to render YAML representation of the Info object.
func (i *Info) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(i, i.low)
	return nb.Render(), nil
}

// Tag represents a high-level AsyncAPI 3.0 Tag object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#tagObject
type Tag struct {
	Name         string                              `json:"name,omitempty" yaml:"name,omitempty"`
	Description  string                              `json:"description,omitempty" yaml:"description,omitempty"`
	ExternalDocs *ExternalDoc                        `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Extensions   *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low          *lowasync.Tag
}

// NewTag creates a new high-level Tag instance from a low-level one.
func NewTag(tag *lowasync.Tag) *Tag {
	t := new(Tag)
	t.low = tag
	t.Name = tag.Name.Value
	t.Description = tag.Description.Value
	if !tag.ExternalDocs.IsEmpty() {
		t.ExternalDocs = NewExternalDoc(tag.ExternalDocs.Value)
	}
	if orderedmap.Len(tag.Extensions) > 0 {
		t.Extensions = high.ExtractExtensions(tag.Extensions)
	}
	return t
}

// GoLow returns the low-level Tag instance.
func (t *Tag) GoLow() *lowasync.Tag {
	return t.low
}

// GoLowUntyped returns the low-level Tag instance with no type.
func (t *Tag) GoLowUntyped() any {
	return t.low
}

// Render will return a YAML representation of the Tag object as a byte slice.
func (t *Tag) Render() ([]byte, error) {
	return yaml.Marshal(t)
}

// MarshalYAML will create a ready to render YAML representation of the Tag object.
func (t *Tag) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(t, t.low)
	return nb.Render(), nil
}

// ExternalDoc represents a high-level AsyncAPI 3.0 External Documentation object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#externalDocumentationObject
type ExternalDoc struct {
	Description string                              `json:"description,omitempty" yaml:"description,omitempty"`
	URL         string                              `json:"url,omitempty" yaml:"url,omitempty"`
	Extensions  *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low         *lowasync.ExternalDoc
}

// NewExternalDoc creates a new high-level ExternalDoc instance from a low-level one.
func NewExternalDoc(extDoc *lowasync.ExternalDoc) *ExternalDoc {
	e := new(ExternalDoc)
	e.low = extDoc
	e.Description = extDoc.Description.Value
	e.URL = extDoc.URL.Value
	if orderedmap.Len(extDoc.Extensions) > 0 {
		e.Extensions = high.ExtractExtensions(extDoc.Extensions)
	}
	return e
}

// GoLow returns the low-level ExternalDoc instance.
func (e *ExternalDoc) GoLow() *lowasync.ExternalDoc {
	return e.low
}

// GoLowUntyped returns the low-level ExternalDoc instance with no type.
func (e *ExternalDoc) GoLowUntyped() any {
	return e.low
}

// Render will return a YAML representation of the ExternalDoc object as a byte slice.
func (e *ExternalDoc) Render() ([]byte, error) {
	return yaml.Marshal(e)
}

// MarshalYAML will create a ready to render YAML representation of the ExternalDoc object.
func (e *ExternalDoc) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(e, e.low)
	return nb.Render(), nil
}
