// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/high"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
	"go.yaml.in/yaml/v4"
)

// Server represents a high-level AsyncAPI 3.0 Server object, backed by a low-level one.
//
// An object representing a message broker, a server or any other kind of computer program capable of
// sending and/or receiving data.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverObject
type Server struct {
	Host            string                                   `json:"host,omitempty" yaml:"host,omitempty"`
	Protocol        string                                   `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	ProtocolVersion string                                   `json:"protocolVersion,omitempty" yaml:"protocolVersion,omitempty"`
	Pathname        string                                   `json:"pathname,omitempty" yaml:"pathname,omitempty"`
	Description     string                                   `json:"description,omitempty" yaml:"description,omitempty"`
	Title           string                                   `json:"title,omitempty" yaml:"title,omitempty"`
	Summary         string                                   `json:"summary,omitempty" yaml:"summary,omitempty"`
	Variables       *orderedmap.Map[string, *ServerVariable] `json:"variables,omitempty" yaml:"variables,omitempty"`
	Security        []*SecurityScheme                        `json:"security,omitempty" yaml:"security,omitempty"`
	Tags            []*Tag                                   `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocs    *ExternalDoc                             `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	Bindings        *ServerBindings                          `json:"bindings,omitempty" yaml:"bindings,omitempty"`
	Extensions      *orderedmap.Map[string, *yaml.Node]      `json:"-" yaml:"-"`
	low             *lowasync.Server
}

// NewServer creates a new high-level Server instance from a low-level one.
func NewServer(server *lowasync.Server) *Server {
	s := new(Server)
	s.low = server
	s.Host = server.Host.Value
	s.Protocol = server.Protocol.Value
	s.ProtocolVersion = server.ProtocolVersion.Value
	s.Pathname = server.Pathname.Value
	s.Description = server.Description.Value
	s.Title = server.Title.Value
	s.Summary = server.Summary.Value

	if server.Variables.Value != nil {
		s.Variables = low.FromReferenceMapWithFunc(server.Variables.Value, NewServerVariable)
	}
	if server.Security.Value != nil {
		for _, sec := range server.Security.Value {
			s.Security = append(s.Security, NewSecurityScheme(sec.Value))
		}
	}
	if server.Tags.Value != nil {
		for _, tag := range server.Tags.Value {
			s.Tags = append(s.Tags, NewTag(tag.Value))
		}
	}
	if !server.ExternalDocs.IsEmpty() {
		s.ExternalDocs = NewExternalDoc(server.ExternalDocs.Value)
	}
	if !server.Bindings.IsEmpty() {
		s.Bindings = NewServerBindings(server.Bindings.Value)
	}
	if orderedmap.Len(server.Extensions) > 0 {
		s.Extensions = high.ExtractExtensions(server.Extensions)
	}
	return s
}

// GoLow returns the low-level Server instance.
func (s *Server) GoLow() *lowasync.Server {
	return s.low
}

// GoLowUntyped returns the low-level Server instance with no type.
func (s *Server) GoLowUntyped() any {
	return s.low
}

// Render will return a YAML representation of the Server object as a byte slice.
func (s *Server) Render() ([]byte, error) {
	return yaml.Marshal(s)
}

// MarshalYAML will create a ready to render YAML representation of the Server object.
func (s *Server) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(s, s.low)
	return nb.Render(), nil
}

// ServerVariable represents a high-level AsyncAPI 3.0 Server Variable object.
//
//	https://www.asyncapi.com/docs/reference/specification/v3.0.0#serverVariableObject
type ServerVariable struct {
	Enum        []string                            `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string                              `json:"default,omitempty" yaml:"default,omitempty"`
	Description string                              `json:"description,omitempty" yaml:"description,omitempty"`
	Examples    []string                            `json:"examples,omitempty" yaml:"examples,omitempty"`
	Extensions  *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low         *lowasync.ServerVariable
}

// NewServerVariable creates a new high-level ServerVariable instance from a low-level one.
func NewServerVariable(sv *lowasync.ServerVariable) *ServerVariable {
	v := new(ServerVariable)
	v.low = sv
	v.Default = sv.Default.Value
	v.Description = sv.Description.Value
	if sv.Enum.Value != nil {
		for _, e := range sv.Enum.Value {
			v.Enum = append(v.Enum, e.Value)
		}
	}
	if sv.Examples.Value != nil {
		for _, ex := range sv.Examples.Value {
			v.Examples = append(v.Examples, ex.Value)
		}
	}
	if orderedmap.Len(sv.Extensions) > 0 {
		v.Extensions = high.ExtractExtensions(sv.Extensions)
	}
	return v
}

// GoLow returns the low-level ServerVariable instance.
func (v *ServerVariable) GoLow() *lowasync.ServerVariable {
	return v.low
}

// GoLowUntyped returns the low-level ServerVariable instance with no type.
func (v *ServerVariable) GoLowUntyped() any {
	return v.low
}

// Render will return a YAML representation of the ServerVariable object as a byte slice.
func (v *ServerVariable) Render() ([]byte, error) {
	return yaml.Marshal(v)
}

// MarshalYAML will create a ready to render YAML representation of the ServerVariable object.
func (v *ServerVariable) MarshalYAML() (interface{}, error) {
	nb := high.NewNodeBuilder(v, v.low)
	return nb.Render(), nil
}

// ServerBindings represents a high-level AsyncAPI 3.0 Server Bindings object.
type ServerBindings struct {
	HTTP       *HTTPServerBinding                  `json:"http,omitempty" yaml:"http,omitempty"`
	Kafka      *KafkaServerBinding                 `json:"kafka,omitempty" yaml:"kafka,omitempty"`
	MQTT       *MQTTServerBinding                  `json:"mqtt,omitempty" yaml:"mqtt,omitempty"`
	SQS        *SQSServerBinding                   `json:"sqs,omitempty" yaml:"sqs,omitempty"`
	Extensions *orderedmap.Map[string, *yaml.Node] `json:"-" yaml:"-"`
	low        *lowasync.ServerBindings
}

// NewServerBindings creates a new high-level ServerBindings instance from a low-level one.
func NewServerBindings(sb *lowasync.ServerBindings) *ServerBindings {
	b := new(ServerBindings)
	b.low = sb
	if !sb.HTTP.IsEmpty() {
		b.HTTP = NewHTTPServerBinding(sb.HTTP.Value)
	}
	if !sb.Kafka.IsEmpty() {
		b.Kafka = NewKafkaServerBinding(sb.Kafka.Value)
	}
	if !sb.MQTT.IsEmpty() {
		b.MQTT = NewMQTTServerBinding(sb.MQTT.Value)
	}
	if !sb.SQS.IsEmpty() {
		b.SQS = NewSQSServerBinding(sb.SQS.Value)
	}
	if orderedmap.Len(sb.Extensions) > 0 {
		b.Extensions = high.ExtractExtensions(sb.Extensions)
	}
	return b
}

// GoLow returns the low-level ServerBindings instance.
func (b *ServerBindings) GoLow() *lowasync.ServerBindings {
	return b.low
}

// GoLowUntyped returns the low-level ServerBindings instance with no type.
func (b *ServerBindings) GoLowUntyped() any {
	return b.low
}
