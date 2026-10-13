// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package yamlparse reads YAML and JSON into source nodes.
package yamlparse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/pb33f/go-yaml"
)

// Parse reads one document. Valid JSON uses the standard JSON decoder so JSON
// escapes and whitespace retain their meaning. Nodes keep source positions.
func Parse(data []byte) (*yaml.Node, error) {
	if json.Valid(data) {
		p := jsonNodes{data: data, decoder: json.NewDecoder(bytes.NewReader(data)), line: 1, column: 1}
		p.decoder.UseNumber()
		root, err := p.value(0)
		if err != nil {
			return nil, err
		}
		return &yaml.Node{Kind: yaml.DocumentNode, Line: root.Line, Column: root.Column, Content: []*yaml.Node{root}}, nil
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var root yaml.Node
	if err := decoder.Decode(&root); err != nil {
		return nil, err
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("expected one YAML document")
	}
	return &root, nil
}

type jsonNodes struct {
	data                 []byte
	decoder              *json.Decoder
	offset, line, column int
}

func (p *jsonNodes) value(depth int) (*yaml.Node, error) {
	if depth > 512 {
		return nil, fmt.Errorf("JSON nesting limit exceeded")
	}
	start := int(p.decoder.InputOffset())
	for start < len(p.data) && strings.ContainsRune(" \t\r\n:,", rune(p.data[start])) {
		start++
	}
	for p.offset < start {
		c, size := utf8.DecodeRune(p.data[p.offset:])
		if c == '\n' || (c == '\r' && (p.offset+1 == len(p.data) || p.data[p.offset+1] != '\n')) {
			p.line, p.column = p.line+1, 1
		} else if c != '\r' {
			p.column++
		}
		p.offset += size
	}
	token, err := p.decoder.Token()
	if err != nil {
		return nil, err
	}
	node := &yaml.Node{Kind: yaml.ScalarNode, Line: p.line, Column: p.column}
	switch value := token.(type) {
	case json.Delim:
		node.Style = yaml.FlowStyle
		if value == '{' {
			node.Kind, node.Tag = yaml.MappingNode, "!!map"
		} else {
			node.Kind, node.Tag = yaml.SequenceNode, "!!seq"
		}
		for p.decoder.More() {
			child, err := p.value(depth + 1)
			if err != nil {
				return nil, err
			}
			node.Content = append(node.Content, child)
			if value == '{' {
				child, err = p.value(depth + 1)
				if err != nil {
					return nil, err
				}
				node.Content = append(node.Content, child)
			}
		}
		if _, err := p.decoder.Token(); err != nil {
			return nil, err
		}
	case string:
		node.Tag, node.Style, node.Value = "!!str", yaml.DoubleQuotedStyle, value
	case json.Number:
		node.Tag, node.Value = "!!int", string(value)
		if strings.ContainsAny(string(value), ".eE") {
			node.Tag = "!!float"
		}
	case bool:
		node.Tag, node.Value = "!!bool", "false"
		if value {
			node.Value = "true"
		}
	case nil:
		node.Tag, node.Value = "!!null", "null"
	}
	return node, nil
}
