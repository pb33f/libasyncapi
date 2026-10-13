// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"unicode"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libasyncapi/internal/yamlparse"
	"github.com/pb33f/libopenapi/utils"
)

func field(n *yaml.Node, key string) *yaml.Node {
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			return n.Content[i+1]
		}
	}
	return nil
}

// referenceOnly separates a foreign schema import from schema metadata named $ref.
func referenceOnly(n *yaml.Node) bool {
	return n.Kind == yaml.MappingNode && len(n.Content) == 2 && field(n, "$ref") != nil
}
func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}
func set(n *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(n.Content); i += 2 {
		if n.Content[i].Value == key {
			n.Content[i+1] = value
			return
		}
	}
	n.Content = append(n.Content, scalar(key), value)
}
func ensureMap(n *yaml.Node, key string) *yaml.Node {
	if v := field(n, key); v != nil {
		return v
	}
	v := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	set(n, key, v)
	return v
}
func clone(n *yaml.Node) *yaml.Node {
	out := *n
	out.Content = make([]*yaml.Node, len(n.Content))
	for i, c := range n.Content {
		out.Content[i] = clone(c)
	}
	return &out
}
func reference(n *yaml.Node, ref string) *yaml.Node {
	out := *n
	out.Kind = yaml.MappingNode
	out.Tag = "!!map"
	out.Anchor = ""
	out.Alias = nil
	out.Content = []*yaml.Node{scalar("$ref"), scalar(ref)}
	if r := field(n, "$ref"); r != nil {
		c := *r
		c.Value = ref
		out.Content[1] = &c
	}
	return &out
}
func escape(s string) string { return strings.ReplaceAll(strings.ReplaceAll(s, "~", "~0"), "/", "~1") }
func localRef(pointer string) string {
	if pointer == "" {
		return "#"
	}
	return (&url.URL{Fragment: pointer}).String()
}
func fragment(ref string) (string, error) {
	u, err := url.Parse(ref)
	if err != nil || !strings.HasPrefix(ref, "#") {
		return "", fmt.Errorf("expected local reference, got %q", ref)
	}
	return u.Fragment, nil
}
func pointer(n *yaml.Node, p string) (*yaml.Node, error) {
	if p == "" {
		return n, nil
	}
	if !strings.HasPrefix(p, "/") {
		return nil, fmt.Errorf("invalid JSON pointer %q", p)
	}
	for _, part := range strings.Split(p[1:], "/") {
		for i := 0; i < len(part); i++ {
			if part[i] == '~' {
				if i+1 >= len(part) || (part[i+1] != '0' && part[i+1] != '1') {
					return nil, fmt.Errorf("invalid JSON pointer %q", p)
				}
				i++
			}
		}
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		switch n.Kind {
		case yaml.MappingNode:
			n = field(n, part)
		case yaml.SequenceNode:
			i, err := strconv.Atoi(part)
			if err != nil || i < 0 || i >= len(n.Content) || strconv.Itoa(i) != part {
				n = nil
			} else {
				n = n.Content[i]
			}
		default:
			n = nil
		}
		if n == nil {
			return nil, fmt.Errorf("JSON pointer #%s not found", p)
		}
	}
	return n, nil
}
func componentName(t target) string {
	name := path.Base(t.pointer)
	if t.pointer == "" {
		u, _ := url.Parse(t.doc.uri)
		name = strings.TrimSuffix(path.Base(u.Path), path.Ext(u.Path))
	}
	name = strings.ReplaceAll(strings.ReplaceAll(name, "~1", "/"), "~0", "~")
	name = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)
	if name == "" || name == "." {
		return "imported"
	}
	return name
}
func indentation(root *yaml.Node) int {
	indent := 0
	var visit func(*yaml.Node)
	visit = func(n *yaml.Node) {
		if n.Style&yaml.FlowStyle != 0 || n.Kind == yaml.ScalarNode {
			return
		}
		if n.Kind == yaml.MappingNode {
			for i := 0; i < len(n.Content); i += 2 {
				key, value := n.Content[i], n.Content[i+1]
				if (value.Kind == yaml.MappingNode || value.Kind == yaml.SequenceNode) && value.Style&yaml.FlowStyle == 0 && value.Line > key.Line {
					if delta := value.Column - key.Column; delta > 0 && (indent == 0 || delta < indent) {
						indent = delta
					}
				}
				visit(value)
			}
			return
		}
		for _, child := range n.Content {
			visit(child)
		}
	}
	visit(root)
	if indent == 0 {
		return 2
	}
	return indent
}

func componentKind(kind string) bool {
	switch kind {
	case "schemas", "servers", "channels", "operations", "messages", "securitySchemes", "serverVariables", "parameters", "correlationIds", "replies", "replyAddresses", "externalDocs", "tags", "operationTraits", "messageTraits", "serverBindings", "channelBindings", "operationBindings", "messageBindings":
		return true
	}
	return false
}

func childKind(parent, key string) string {
	if strings.HasPrefix(parent, "map:") {
		return strings.TrimPrefix(parent, "map:")
	}
	if parent == "extension" {
		return "extension"
	}
	if parent == "root" {
		switch key {
		case "servers", "channels", "operations":
			return "map:" + key
		case "components":
			return "components"
		case "info":
			return "info"
		}
	}
	if parent == "components" && componentKind(key) {
		return "map:" + key
	}
	if parent == "schemas" {
		switch key {
		case "properties", "patternProperties", "definitions", "$defs", "dependentSchemas", "dependencies":
			return "map:schemas"
		case "allOf", "oneOf", "anyOf", "prefixItems":
			return "list:schemas"
		case "items", "additionalItems", "additionalProperties", "unevaluatedProperties", "unevaluatedItems", "contains", "not", "if", "then", "else", "propertyNames", "schema":
			return "schemas"
		}
		return "opaque"
	}
	switch key {
	case "tags":
		return "list:tags"
	case "externalDocs":
		return "externalDocs"
	case "bindings":
		switch parent {
		case "servers":
			return "serverBindings"
		case "channels":
			return "channelBindings"
		case "operations", "operationTraits":
			return "operationBindings"
		case "messages", "messageTraits":
			return "messageBindings"
		}
	case "traits":
		if parent == "operations" {
			return "list:operationTraits"
		}
		if parent == "messages" {
			return "list:messageTraits"
		}
	case "security":
		return "list:securitySchemes"
	case "variables":
		if parent == "servers" {
			return "map:serverVariables"
		}
	case "parameters":
		if parent == "channels" {
			return "map:parameters"
		}
	case "servers":
		if parent == "channels" {
			return "list:ref:servers"
		}
	case "channel":
		if parent == "operations" || parent == "replies" {
			return "ref:channels"
		}
	case "messages":
		if parent == "channels" {
			return "map:messages"
		}
		if parent == "operations" || parent == "replies" {
			return "list:ref:messages"
		}
	case "reply":
		if parent == "operations" {
			return "replies"
		}
	case "address":
		if parent == "replies" {
			return "replyAddresses"
		}
	case "headers", "payload":
		if parent == "messages" || parent == "messageTraits" {
			return "schemas"
		}
	case "correlationId":
		if parent == "messages" || parent == "messageTraits" {
			return "correlationIds"
		}
	}
	if strings.HasSuffix(parent, "Bindings") {
		return "binding:" + key
	}
	if strings.HasPrefix(parent, "binding:") {
		switch key {
		case "schema", "headers", "query", "key":
			return "schemas"
		}
		if parent == "binding:kafka" && (key == "groupId" || key == "clientId") {
			return "schemas"
		}
		if parent == "binding:mqtt" {
			switch key {
			case "correlationData", "responseTopic", "sessionExpiryInterval", "maximumPacketSize", "messageExpiryInterval":
				return "schemas"
			}
		}
	}
	return "opaque"
}

func validDelimiter(s string) bool {
	return !strings.ContainsAny(s, "#/") && strings.IndexFunc(s, unicode.IsSpace) < 0
}

func jsonSchemaFormat(format string) bool {
	return schemaformat.IsJSON(format)
}

func parse(data []byte) (*yaml.Node, error) {
	root, err := yamlparse.Parse(data)
	if err != nil {
		return nil, err
	}
	budget := 1_000_000
	return expandAliases(root, make(map[*yaml.Node]bool), 0, &budget, false)
}

func expandAliases(n *yaml.Node, active map[*yaml.Node]bool, depth int, budget *int, aliased bool) (*yaml.Node, error) {
	if n.Kind == yaml.AliasNode {
		n = n.Alias
		aliased = true
	}
	if aliased {
		*budget--
	}
	if depth > 512 || *budget < 0 {
		return nil, fmt.Errorf("document expansion limit exceeded")
	}
	if active[n] {
		return nil, fmt.Errorf("circular YAML alias at line %d", n.Line)
	}
	active[n] = true
	defer delete(active, n)
	out := *n
	out.Anchor = ""
	out.Alias = nil
	out.Content = make([]*yaml.Node, len(n.Content))
	for i, child := range n.Content {
		value, err := expandAliases(child, active, depth+1, budget, aliased)
		if err != nil {
			return nil, err
		}
		out.Content[i] = value
	}
	utils.CheckForMergeNodes(&out)
	return &out, nil
}
