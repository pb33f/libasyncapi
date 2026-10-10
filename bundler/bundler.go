// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi"
)

// BundleCompositionConfig controls component names in composed output.
type BundleCompositionConfig struct {
	// Delimiter separates colliding names. The default is "__".
	Delimiter string
}

// BundleBytes inlines external references where AsyncAPI permits an object.
// Root-local references are preserved. Recursive and reference-only fields use
// local references. On failure no output is returned. Input and configuration
// are not modified. Output is YAML, including when the input is JSON.
func BundleBytes(spec []byte, config *libasyncapi.DocumentConfiguration) ([]byte, error) {
	return bundle(spec, config, false, "__")
}

// BundleBytesComposed places external objects in their AsyncAPI component
// sections and rewrites references to local targets. Reference-only fields keep
// the channel, message and server locations required by AsyncAPI.
func BundleBytesComposed(spec []byte, config *libasyncapi.DocumentConfiguration, composition *BundleCompositionConfig) ([]byte, error) {
	delimiter := "__"
	if composition != nil && composition.Delimiter != "" {
		delimiter = composition.Delimiter
	}
	if !validDelimiter(delimiter) {
		return nil, fmt.Errorf("composition delimiter cannot contain whitespace, '#' or '/'")
	}
	return bundle(spec, config, true, delimiter)
}

type componentImport struct {
	kind, name string
	node       *yaml.Node
}

type pendingReference struct {
	node      *yaml.Node
	target    target
	path      string
	role      string
	preferred string
}

type engine struct {
	loader            *loader
	root              *yaml.Node
	composed          bool
	delimiter         string
	excludeExtensions bool
	placements        map[target][]string
	active            map[target]string
	components        map[componentKey]string
	imports           []componentImport
	pending           []pendingReference
	aliases           map[string][]string
}

type componentKey struct {
	target target
	kind   string
}

func bundle(spec []byte, config *libasyncapi.DocumentConfiguration, composed bool, delimiter string) ([]byte, error) {
	input, err := parse(spec)
	if err != nil {
		return nil, fmt.Errorf("parse AsyncAPI: %w", err)
	}
	version, err := libasyncapi.ExtractVersionFromNode(input)
	if err != nil {
		return nil, err
	}
	if _, err = libasyncapi.ParseAsyncAPIVersion(version); err != nil {
		return nil, err
	}
	if !libasyncapi.IsAsyncAPI3(version) {
		return nil, libasyncapi.ErrInvalidAsyncAPIVersion
	}
	cfg := libasyncapi.DocumentConfiguration{}
	if config != nil {
		cfg = *config
	}
	l, err := newLoader(input, &cfg)
	if err != nil {
		return nil, err
	}
	e := &engine{loader: l, composed: composed, delimiter: delimiter,
		excludeExtensions: cfg.ExcludeExtensionRefs, active: make(map[target]string), components: make(map[componentKey]string)}
	e.aliases = make(map[string][]string)
	e.placements = make(map[target][]string)
	// Reserve all existing component names before imports are allocated.
	e.root = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	if components := field(input.Content[0], "components"); components != nil {
		set(e.root, "components", clone(components))
	}
	result, err := e.walk(target{l.root, ""}, "", "root", 0)
	if err != nil {
		return nil, err
	}
	// Reusable operations and channels may point to external reusable objects.
	// These fields require references even in inline mode, so retain a component.
	for i := 0; i < len(e.pending); i++ {
		p := e.pending[i]
		rootUse, err := e.rootUse(p.path)
		if err != nil {
			return nil, err
		}
		if p.role == "messages" || rootUse || !strings.HasPrefix(p.path, "/components/") || len(e.placements[p.target]) > 0 {
			continue
		}
		if _, err := e.compose(p.target, p.role, 0); err != nil {
			return nil, err
		}
	}
	// Imported components were collected separately so traversal never sees its
	// own output as input and root key order remains stable.
	for _, item := range e.imports {
		components := ensureMap(result, "components")
		set(ensureMap(components, item.kind), item.name, item.node)
	}
	e.root = result
	if err := e.finishReferences(); err != nil {
		return nil, err
	}
	output := *input
	output.Content = []*yaml.Node{result}
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(indentation(spec))
	if err := encoder.Encode(&output); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (e *engine) walk(src target, dest, kind string, depth int) (*yaml.Node, error) {
	if depth > 512 {
		return nil, fmt.Errorf("bundle depth exceeded at #%s", dest)
	}
	n, err := src.node()
	if err != nil {
		return nil, err
	}
	if n.Kind == yaml.AliasNode {
		return nil, fmt.Errorf("YAML aliases are not supported at #%s", dest)
	}
	if kind == "opaque" {
		return clone(n), nil
	}
	if (strings.HasPrefix(kind, "map:") || kind == "components" || (componentKind(kind) && kind != "schemas")) && n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("#%s: %s requires an object", dest, kind)
	}
	if strings.HasPrefix(kind, "list:") && n.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("#%s requires an array", dest)
	}
	if kind == "schemas" {
		if err := e.loader.scan(src, kind, src.doc.uri, 0); err != nil {
			return nil, err
		}
	}
	if ref := field(n, "$ref"); ref != nil && !strings.HasPrefix(kind, "map:") && kind != "root" && kind != "components" {
		if ref.Kind != yaml.ScalarNode || ref.Tag != "!!str" || ref.Value == "" {
			return nil, fmt.Errorf("invalid $ref at %s#%s", src.doc.uri, src.pointer)
		}
		var resolved target
		if kind == "schemas" {
			resolved, err = e.loader.schemaReference(src, ref.Value)
		} else {
			resolved, err = e.loader.resolve(src, ref.Value)
		}
		if err != nil {
			return nil, fmt.Errorf("%s#%s: %w", src.doc.uri, src.pointer, err)
		}
		// Reject reference-only loops before treating any local reference as safe.
		if kind != "schemas" {
			resolved, err = e.loader.terminal(resolved)
		} else {
			resolved, err = e.schemaTerminal(resolved)
		}
		if err != nil {
			return nil, err
		}
		if strings.HasPrefix(kind, "ref:") {
			out := reference(n, ref.Value)
			preferred := ""
			if src.doc == e.loader.root && strings.HasPrefix(ref.Value, "#") {
				preferred, _ = fragment(ref.Value)
			}
			e.pending = append(e.pending, pendingReference{out, resolved, dest, strings.TrimPrefix(kind, "ref:"), preferred})
			return out, nil
		}
		if kind == "schemas" && src.doc == e.loader.root && e.loader.scopes[src] == src.doc.uri && strings.HasPrefix(ref.Value, "#/") {
			return reference(n, ref.Value), nil
		}
		if kind == "schemas" && resolved.doc == e.loader.root {
			return reference(n, localRef(resolved.pointer)), nil
		}
		if src.doc == e.loader.root && strings.HasPrefix(ref.Value, "#") && (kind != "schemas" || e.loader.scopes[src] == src.doc.uri) {
			e.placements[resolved] = append(e.placements[resolved], dest)
			if kind != "schemas" && resolved.doc == e.loader.root && resolved.pointer != dest {
				e.aliases[resolved.pointer] = append(e.aliases[resolved.pointer], dest)
			}
			return reference(n, ref.Value), nil
		}
		e.placements[resolved] = append(e.placements[resolved], dest)
		if path, ok := e.active[resolved]; ok {
			return reference(n, localRef(path)), nil
		}
		if e.composed && componentKind(kind) {
			path, err := e.compose(resolved, kind, depth+1)
			if err != nil {
				return nil, err
			}
			e.aliases[path] = append(e.aliases[path], dest)
			return reference(n, localRef(path)), nil
		}
		e.active[resolved] = dest
		out, err := e.walk(resolved, dest, kind, depth+1)
		delete(e.active, resolved)
		return out, err
	}
	if strings.HasPrefix(kind, "ref:") {
		return nil, fmt.Errorf("#%s requires a Reference Object", dest)
	}
	e.placements[src] = append(e.placements[src], dest)
	out := *n
	out.Anchor = ""
	out.Content = nil
	if n.Kind == yaml.MappingNode {
		seen := make(map[string]bool, len(n.Content)/2)
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i]
			if key.Kind != yaml.ScalarNode || key.Tag != "!!str" || seen[key.Value] {
				return nil, fmt.Errorf("invalid or duplicate key at #%s", dest)
			}
			seen[key.Value] = true
			if kind == "schemas" && (key.Value == "$id" || key.Value == "$anchor") {
				continue
			}
			childKind := childKind(kind, key.Value)
			if strings.HasPrefix(key.Value, "x-") && !strings.HasPrefix(kind, "map:") {
				childKind = "extension"
				if e.excludeExtensions {
					childKind = "opaque"
				}
			}
			if kind == "schemas" && field(n, "schemaFormat") != nil && key.Value == "schema" {
				format := field(n, "schemaFormat").Value
				childKind = "opaque"
				if jsonSchemaFormat(format) {
					childKind = "schemas"
				}
			}
			// A schema resource changes reference scope. Until its relocation can
			// be proven safe, fail rather than silently change its meaning.
			if kind == "schemas" && (key.Value == "$dynamicRef" || key.Value == "$recursiveRef") {
				return nil, fmt.Errorf("schema keyword %s at #%s requires resource-aware bundling", key.Value, dest)
			}
			segment := "/" + escape(key.Value)
			child, err := e.walk(target{src.doc, src.pointer + segment}, dest+segment, childKind, depth+1)
			if err != nil {
				return nil, err
			}
			out.Content = append(out.Content, clone(key), child)
		}
	} else if n.Kind == yaml.SequenceNode {
		for i := range n.Content {
			segment := fmt.Sprintf("/%d", i)
			child, err := e.walk(target{src.doc, src.pointer + segment}, dest+segment, strings.TrimPrefix(kind, "list:"), depth+1)
			if err != nil {
				return nil, err
			}
			out.Content = append(out.Content, child)
		}
	}
	return &out, nil
}

func (e *engine) schemaTerminal(t target) (target, error) {
	seen := make(map[target]bool)
	for {
		if seen[t] {
			return target{}, fmt.Errorf("reference-only schema cycle at %s#%s", t.doc.uri, t.pointer)
		}
		seen[t] = true
		if err := e.loader.scan(t, "schemas", t.doc.uri, 0); err != nil {
			return target{}, err
		}
		n, err := t.node()
		if err != nil {
			return target{}, err
		}
		r := field(n, "$ref")
		if r == nil {
			return t, nil
		}
		if r.Kind != yaml.ScalarNode || r.Tag != "!!str" || r.Value == "" {
			return target{}, fmt.Errorf("invalid schema $ref")
		}
		t, err = e.loader.schemaReference(t, r.Value)
		if err != nil {
			return target{}, err
		}
	}
}

func (e *engine) compose(src target, kind string, depth int) (string, error) {
	key := componentKey{src, kind}
	if path, ok := e.components[key]; ok {
		return path, nil
	}
	name := componentName(src)
	bucket := ensureMap(ensureMap(e.root, "components"), kind)
	if field(bucket, name) != nil {
		name += e.delimiter + componentName(target{src.doc, ""})
	}
	base := name
	for i := 1; field(bucket, name) != nil; i++ {
		name = fmt.Sprintf("%s%s%d", base, e.delimiter, i)
	}
	set(bucket, name, &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"})
	path := "/components/" + kind + "/" + escape(name)
	e.components[key] = path
	e.active[src] = path
	value, err := e.walk(src, path, kind, depth)
	delete(e.active, src)
	if err != nil {
		return "", err
	}
	e.imports = append(e.imports, componentImport{kind, name, value})
	return path, nil
}

func (e *engine) finishReferences() error {
	// Channel/server targets must settle before message membership is checked.
	for _, messages := range []bool{false, true} {
		for _, pending := range e.pending {
			if (pending.role == "messages") != messages {
				continue
			}
			prefix := "/" + pending.role + "/"
			if messages {
				container := pending.path[:strings.LastIndex(pending.path, "/messages/")]
				channel, err := pointer(e.root, container+"/channel/$ref")
				if err != nil {
					return fmt.Errorf("#%s: messages require a channel", pending.path)
				}
				u, err := fragment(channel.Value)
				if err != nil {
					return err
				}
				prefix = u + "/messages/"
			}
			candidates := make(map[string]bool)
			for ancestor := pending.target; ; {
				suffix := strings.TrimPrefix(pending.target.pointer, ancestor.pointer)
				for _, place := range e.placements[ancestor] {
					paths, err := e.project(place + suffix)
					if err != nil {
						return err
					}
					for _, path := range paths {
						if strings.HasPrefix(path, prefix) && !strings.Contains(strings.TrimPrefix(path, prefix), "/") {
							if _, err := e.outputTarget(path); err == nil {
								candidates[path] = true
							}
						}
					}
				}
				if ancestor.pointer == "" {
					break
				}
				ancestor.pointer = ancestor.pointer[:strings.LastIndex(ancestor.pointer, "/")]
			}
			// Component operations/channels can reference reusable objects too.
			rootUse, err := e.rootUse(pending.path)
			if err != nil {
				return err
			}
			if len(candidates) == 0 && !rootUse && strings.HasPrefix(pending.path, "/components/") && !messages {
				for _, place := range e.placements[pending.target] {
					candidates[place] = true
				}
			}
			if candidates[pending.preferred] {
				candidates = map[string]bool{pending.preferred: true}
			}
			if len(candidates) != 1 {
				return fmt.Errorf("#%s: expected one local %s target, found %d", pending.path, pending.role, len(candidates))
			}
			for path := range candidates {
				field(pending.node, "$ref").Value = localRef(path)
			}
		}
	}
	return nil
}

func (e *engine) rootUse(path string) (bool, error) {
	paths, err := e.project(path)
	if err != nil {
		return false, err
	}
	for _, p := range paths {
		if strings.HasPrefix(p, "/operations/") || strings.HasPrefix(p, "/channels/") {
			return true, nil
		}
	}
	return false, nil
}

// project exposes logical paths through composed component aliases.
func (e *engine) project(path string) ([]string, error) {
	paths := []string{path}
	seen := map[string]bool{path: true}
	history := map[string][]string{path: nil}
	for i := 0; i < len(paths); i++ {
		if i > 4096 {
			return nil, fmt.Errorf("reference placement expansion limit exceeded at #%s", path)
		}
		p := paths[i]
		for prefix := p; prefix != ""; {
			for _, alias := range e.aliases[prefix] {
				if slices.Contains(history[p], prefix) {
					return nil, fmt.Errorf("cyclic reference placement at #%s", path)
				}
				next := alias + strings.TrimPrefix(p, prefix)
				if !seen[next] {
					seen[next] = true
					history[next] = append(append([]string(nil), history[p]...), prefix)
					paths = append(paths, next)
				}
			}
			j := strings.LastIndex(prefix, "/")
			if j < 0 {
				break
			}
			prefix = prefix[:j]
		}
	}
	return paths, nil
}

func (e *engine) outputTarget(path string) (*yaml.Node, error) {
	seen := make(map[string]bool)
	for !seen[path] {
		seen[path] = true
		// Follow reference ancestors as well as terminal references.
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
		n := e.root
		redirect := false
		for i := 0; i <= len(parts); i++ {
			if r := field(n, "$ref"); r != nil && r.Kind == yaml.ScalarNode && r.Tag == "!!str" {
				p, err := fragment(r.Value)
				if err != nil {
					return nil, err
				}
				path = p
				if i < len(parts) {
					path += "/" + strings.Join(parts[i:], "/")
				}
				redirect = true
				break
			}
			if i == len(parts) {
				return n, nil
			}
			var err error
			n, err = pointer(n, "/"+parts[i])
			if err != nil {
				return nil, err
			}
		}
		if !redirect {
			return n, nil
		}
	}
	return nil, fmt.Errorf("reference loop at #%s", path)
}
