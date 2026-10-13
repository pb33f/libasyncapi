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
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libasyncapi/internal/yamlparse"
	"github.com/pb33f/libopenapi/datamodel/low/base"
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
	return n.Kind == yaml.MappingNode && len(n.Content) == 2 && field(n, lowasync.RefLabel) != nil
}
func scalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: yamlparse.StringTag, Value: value}
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
	v := &yaml.Node{Kind: yaml.MappingNode, Tag: yamlparse.MappingTag}
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
	out.Tag = yamlparse.MappingTag
	out.Anchor = ""
	out.Alias = nil
	out.Content = []*yaml.Node{scalar(lowasync.RefLabel), scalar(ref)}
	if r := field(n, lowasync.RefLabel); r != nil {
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
		return defaultComponentName
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
	case lowasync.SchemasLabel, lowasync.ServersLabel, lowasync.ChannelsLabel,
		lowasync.OperationsLabel, lowasync.MessagesLabel, lowasync.SecuritySchemesLabel,
		lowasync.ServerVariablesLabel, lowasync.ParametersLabel, lowasync.CorrelationIDsLabel,
		lowasync.RepliesLabel, lowasync.ReplyAddressesLabel, lowasync.ExternalDocsLabel,
		lowasync.TagsLabel, lowasync.OperationTraitsLabel, lowasync.MessageTraitsLabel,
		lowasync.ServerBindingsLabel, lowasync.ChannelBindingsLabel,
		lowasync.OperationBindingsLabel, lowasync.MessageBindingsLabel:
		return true
	}
	return false
}

func childKind(parent, key string) string {
	if strings.HasPrefix(parent, mapKindPrefix) {
		return strings.TrimPrefix(parent, mapKindPrefix)
	}
	if parent == kindExtension {
		return kindExtension
	}
	if parent == kindRoot {
		switch key {
		case lowasync.ServersLabel, lowasync.ChannelsLabel, lowasync.OperationsLabel:
			return mapKindPrefix + key
		case lowasync.ComponentsLabel:
			return lowasync.ComponentsLabel
		case lowasync.InfoLabel:
			return lowasync.InfoLabel
		}
	}
	if parent == lowasync.ComponentsLabel && componentKind(key) {
		return mapKindPrefix + key
	}
	if parent == lowasync.SchemasLabel {
		switch key {
		case base.PropertiesLabel, base.PatternPropertiesLabel, definitionsLabel,
			base.DefsLabel, base.DependentSchemasLabel, dependenciesLabel:
			return mapKindPrefix + lowasync.SchemasLabel
		case base.AllOfLabel, base.OneOfLabel, base.AnyOfLabel, base.PrefixItemsLabel:
			return listKindPrefix + lowasync.SchemasLabel
		case base.ItemsLabel, additionalItemsLabel, base.AdditionalPropertiesLabel,
			base.UnevaluatedPropertiesLabel, base.UnevaluatedItemsLabel, base.ContainsLabel,
			base.NotLabel, base.IfLabel, base.ThenLabel, base.ElseLabel,
			base.PropertyNamesLabel, lowasync.SchemaLabel:
			return lowasync.SchemasLabel
		}
		return kindOpaque
	}
	switch key {
	case lowasync.TagsLabel:
		return listKindPrefix + lowasync.TagsLabel
	case lowasync.ExternalDocsLabel:
		return lowasync.ExternalDocsLabel
	case lowasync.BindingsLabel:
		switch parent {
		case lowasync.ServersLabel:
			return lowasync.ServerBindingsLabel
		case lowasync.ChannelsLabel:
			return lowasync.ChannelBindingsLabel
		case lowasync.OperationsLabel, lowasync.OperationTraitsLabel:
			return lowasync.OperationBindingsLabel
		case lowasync.MessagesLabel, lowasync.MessageTraitsLabel:
			return lowasync.MessageBindingsLabel
		}
	case lowasync.TraitsLabel:
		if parent == lowasync.OperationsLabel {
			return listKindPrefix + lowasync.OperationTraitsLabel
		}
		if parent == lowasync.MessagesLabel {
			return listKindPrefix + lowasync.MessageTraitsLabel
		}
	case lowasync.SecurityLabel:
		return listKindPrefix + lowasync.SecuritySchemesLabel
	case lowasync.VariablesLabel:
		if parent == lowasync.ServersLabel {
			return mapKindPrefix + lowasync.ServerVariablesLabel
		}
	case lowasync.ParametersLabel:
		if parent == lowasync.ChannelsLabel {
			return mapKindPrefix + lowasync.ParametersLabel
		}
	case lowasync.ServersLabel:
		if parent == lowasync.ChannelsLabel {
			return listKindPrefix + refKindPrefix + lowasync.ServersLabel
		}
	case lowasync.ChannelLabel:
		if parent == lowasync.OperationsLabel || parent == lowasync.RepliesLabel {
			return refKindPrefix + lowasync.ChannelsLabel
		}
	case lowasync.MessagesLabel:
		if parent == lowasync.ChannelsLabel {
			return mapKindPrefix + lowasync.MessagesLabel
		}
		if parent == lowasync.OperationsLabel || parent == lowasync.RepliesLabel {
			return listKindPrefix + refKindPrefix + lowasync.MessagesLabel
		}
	case lowasync.ReplyLabel:
		if parent == lowasync.OperationsLabel {
			return lowasync.RepliesLabel
		}
	case lowasync.AddressLabel:
		if parent == lowasync.RepliesLabel {
			return lowasync.ReplyAddressesLabel
		}
	case lowasync.HeadersLabel, lowasync.PayloadLabel:
		if parent == lowasync.MessagesLabel || parent == lowasync.MessageTraitsLabel {
			return lowasync.SchemasLabel
		}
	case lowasync.CorrelationIDLabel:
		if parent == lowasync.MessagesLabel || parent == lowasync.MessageTraitsLabel {
			return lowasync.CorrelationIDsLabel
		}
	}
	if strings.HasSuffix(parent, bindingsSuffix) {
		return bindingKindPrefix + key
	}
	if strings.HasPrefix(parent, bindingKindPrefix) {
		switch key {
		case lowasync.SchemaLabel, lowasync.HeadersLabel, lowasync.QueryLabel, lowasync.KeyLabel:
			return lowasync.SchemasLabel
		}
		if parent == bindingKindPrefix+lowasync.KafkaLabel && (key == lowasync.GroupIDLabel || key == lowasync.ClientIDLabel) {
			return lowasync.SchemasLabel
		}
		if parent == bindingKindPrefix+lowasync.MQTTLabel {
			switch key {
			case lowasync.CorrelationDataLabel, lowasync.ResponseTopicLabel,
				lowasync.SessionExpiryIntervalLabel, lowasync.MaximumPacketSizeLabel,
				lowasync.MessageExpiryIntervalLabel:
				return lowasync.SchemasLabel
			}
		}
	}
	return kindOpaque
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
