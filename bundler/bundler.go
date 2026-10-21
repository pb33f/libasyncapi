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
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/internal/yamlparse"
	"github.com/pb33f/libopenapi/datamodel/low/base"
)

// BundleCompositionConfig controls component names in composed output.
type BundleCompositionConfig struct {
	// Delimiter separates colliding names. The default is "__".
	Delimiter string
}

// BundleBytes inlines external references where AsyncAPI permits an object.
// Root-local references are preserved. Recursive and reference-only fields use
// local references. External channels and servers required by root operations
// are imported at the root when no existing destination is available.
// On failure no output is returned. Input and configuration
// are not modified. Output is YAML, including when the input is JSON.
func BundleBytes(spec []byte, config *libasyncapi.DocumentConfiguration) ([]byte, error) {
	return bundle(spec, config, false, defaultDelimiter)
}

// BundleBytesComposed places external objects in their AsyncAPI component
// sections and rewrites references to local targets. Reference-only fields keep
// the channel, message and server locations required by AsyncAPI. Missing root
// destinations for required external channels and servers are added as aliases.
func BundleBytesComposed(spec []byte, config *libasyncapi.DocumentConfiguration, composition *BundleCompositionConfig) ([]byte, error) {
	delimiter := defaultDelimiter
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
	inlineSchemas     map[target]string
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
	e.inlineSchemas = make(map[target]string)
	e.placements = make(map[target][]string)
	// Reserve all existing component names before imports are allocated.
	e.root = &yaml.Node{Kind: yaml.MappingNode, Tag: yamlparse.MappingTag}
	if components := field(input.Content[0], lowasync.ComponentsLabel); components != nil {
		set(e.root, lowasync.ComponentsLabel, clone(components))
	}
	if composed {
		if err := e.reserveComponentTargets(); err != nil {
			return nil, err
		}
	}
	result, err := e.walk(target{l.root, ""}, "", kindRoot, 0)
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
		if p.role == lowasync.MessagesLabel || rootUse || !strings.HasPrefix(p.path, componentsPathPrefix) || len(e.placements[p.target]) > 0 {
			continue
		}
		if _, err := e.compose(p.target, p.role, 0); err != nil {
			return nil, err
		}
	}
	// Imported components were collected separately so traversal never sees its
	// own output as input and root key order remains stable.
	for _, item := range e.imports {
		components := ensureMap(result, lowasync.ComponentsLabel)
		set(ensureMap(components, item.kind), item.name, item.node)
	}
	e.root = result
	if err := e.materializeRootTargets(); err != nil {
		return nil, err
	}
	if err := e.finishReferences(); err != nil {
		return nil, err
	}
	output := *input
	output.Content = []*yaml.Node{result}
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(indentation(input))
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
	// Resolve the outer Reference Object of a foreign schema, but do not
	// interpret its contents as JSON Schema once the resource is loaded.
	if kind == kindOpaque || (kind == kindSchemaContent && !referenceOnly(n)) {
		return clone(n), nil
	}
	if (strings.HasPrefix(kind, mapKindPrefix) || kind == lowasync.ComponentsLabel || (componentKind(kind) && kind != lowasync.SchemasLabel)) && n.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("#%s: %s requires an object", dest, kind)
	}
	if strings.HasPrefix(kind, listKindPrefix) && n.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("#%s requires an array", dest)
	}
	if kind == lowasync.SchemasLabel {
		if err := e.loader.scan(src, kind, src.doc.uri, 0); err != nil {
			return nil, err
		}
	}
	if ref := field(n, lowasync.RefLabel); ref != nil && !strings.HasPrefix(kind, mapKindPrefix) && kind != kindRoot && kind != lowasync.ComponentsLabel {
		if ref.Kind != yaml.ScalarNode || ref.Tag != yamlparse.StringTag || ref.Value == "" {
			return nil, fmt.Errorf("invalid $ref at %s#%s", src.doc.uri, src.pointer)
		}
		resolved, err := e.resolveReference(src, kind, ref.Value)
		if err != nil {
			return nil, err
		}
		resolvedNode, err := resolved.node()
		if err != nil {
			return nil, err
		}
		if kind != kindSchemaContent && resolved.pointer == "" && resolvedNode.Kind == yaml.ScalarNode && resolvedNode.Tag != yamlparse.BooleanTag {
			return nil, fmt.Errorf("reference target %s#%s must contain an object, array, or boolean schema", resolved.doc.uri, resolved.pointer)
		}
		if strings.HasPrefix(kind, refKindPrefix) {
			out := reference(n, ref.Value)
			preferred := ""
			if src.doc == e.loader.root && strings.HasPrefix(ref.Value, "#") {
				preferred, _ = fragment(ref.Value)
			}
			e.pending = append(e.pending, pendingReference{out, resolved, dest, strings.TrimPrefix(kind, refKindPrefix), preferred})
			return out, nil
		}
		if kind == lowasync.SchemasLabel && src.doc == e.loader.root && e.loader.scopes[src] == src.doc.uri && strings.HasPrefix(ref.Value, "#/") {
			return reference(n, ref.Value), nil
		}
		if kind == lowasync.SchemasLabel && resolved.doc == e.loader.root {
			return reference(n, localRef(resolved.pointer)), nil
		}
		if src.doc == e.loader.root && strings.HasPrefix(ref.Value, "#") && (kind != lowasync.SchemasLabel || e.loader.scopes[src] == src.doc.uri) {
			e.placements[resolved] = append(e.placements[resolved], dest)
			if kind != lowasync.SchemasLabel {
				// Preserve the local alias even when its terminal target is external.
				local, _ := fragment(ref.Value)
				if local != dest {
					e.aliases[local] = append(e.aliases[local], dest)
				}
			}
			return reference(n, ref.Value), nil
		}
		e.placements[resolved] = append(e.placements[resolved], dest)
		if path, ok := e.active[resolved]; ok {
			return reference(n, localRef(path)), nil
		}
		if !e.composed && kind == lowasync.SchemasLabel {
			// Emit each external schema once. Shared graphs otherwise expand
			// exponentially even though they contain no reference cycle.
			if path, ok := e.inlineSchemas[resolved]; ok {
				return reference(n, localRef(path)), nil
			}
			e.inlineSchemas[resolved] = dest
		}
		if e.composed && componentKind(kind) && e.components[componentKey{resolved, kind}] != dest {
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
	if strings.HasPrefix(kind, refKindPrefix) {
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
			if key.Kind != yaml.ScalarNode || seen[key.Value] {
				return nil, fmt.Errorf("invalid or duplicate key at #%s", dest)
			}
			seen[key.Value] = true
			if kind == lowasync.SchemasLabel && (key.Value == base.IdLabel || key.Value == base.AnchorLabel) {
				continue
			}
			childKind := childKind(kind, key.Value)
			if strings.HasPrefix(key.Value, extensionPrefix) && !strings.HasPrefix(kind, mapKindPrefix) {
				childKind = kindExtension
				if e.excludeExtensions {
					childKind = kindOpaque
				}
			}
			if kind == lowasync.SchemasLabel && field(n, lowasync.SchemaFormatLabel) != nil && key.Value == lowasync.SchemaLabel {
				format := field(n, lowasync.SchemaFormatLabel).Value
				childKind = kindSchemaContent
				if jsonSchemaFormat(format) {
					childKind = lowasync.SchemasLabel
				}
			}
			// A schema resource changes reference scope. Until its relocation can
			// be proven safe, fail rather than silently change its meaning.
			if kind == lowasync.SchemasLabel && (key.Value == base.DynamicRefLabel || key.Value == recursiveRefLabel) {
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
			child, err := e.walk(target{src.doc, src.pointer + segment}, dest+segment, strings.TrimPrefix(kind, listKindPrefix), depth+1)
			if err != nil {
				return nil, err
			}
			out.Content = append(out.Content, child)
		}
	}
	return &out, nil
}

func (e *engine) resolveReference(src target, kind, ref string) (target, error) {
	var resolved target
	var err error
	if kind == lowasync.SchemasLabel {
		resolved, err = e.loader.schemaReference(src, ref)
	} else {
		resolved, err = e.loader.resolve(src, ref)
	}
	if err != nil {
		return target{}, fmt.Errorf("%s#%s: %w", src.doc.uri, src.pointer, err)
	}
	// Reject reference-only loops before treating any local reference as safe.
	if kind == lowasync.SchemasLabel {
		return e.schemaTerminal(resolved)
	}
	return e.loader.terminal(resolved, kind == kindSchemaContent)
}

// Reserve explicit component destinations before walking uses that may precede
// them in the document. The first bare external reference owns the target.
func (e *engine) reserveComponentTargets() error {
	components := field(e.loader.root.root, lowasync.ComponentsLabel)
	if components == nil || components.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(components.Content); i += 2 {
		kind, bucket := components.Content[i].Value, components.Content[i+1]
		if !componentKind(kind) || bucket.Kind != yaml.MappingNode {
			continue
		}
		for j := 0; j < len(bucket.Content); j += 2 {
			n := bucket.Content[j+1]
			ref := field(n, lowasync.RefLabel)
			if len(n.Content) != 2 || ref == nil || ref.Kind != yaml.ScalarNode || ref.Tag != yamlparse.StringTag || ref.Value == "" || strings.HasPrefix(ref.Value, "#") {
				continue
			}
			path := componentsPathPrefix + kind + "/" + escape(bucket.Content[j].Value)
			t, err := e.resolveReference(target{e.loader.root, path}, kind, ref.Value)
			if err != nil {
				return err
			}
			if t.doc != e.loader.root {
				key := componentKey{t, kind}
				if _, exists := e.components[key]; !exists {
					e.components[key] = path
				}
			}
		}
	}
	return nil
}

func (e *engine) schemaTerminal(t target) (target, error) {
	seen := make(map[target]bool)
	for {
		if seen[t] {
			return target{}, fmt.Errorf("reference-only schema cycle at %s#%s", t.doc.uri, t.pointer)
		}
		seen[t] = true
		if err := e.loader.scan(t, lowasync.SchemasLabel, t.doc.uri, 0); err != nil {
			return target{}, err
		}
		n, err := t.node()
		if err != nil {
			return target{}, err
		}
		r := field(n, lowasync.RefLabel)
		if r == nil {
			return t, nil
		}
		if r.Kind != yaml.ScalarNode || r.Tag != yamlparse.StringTag || r.Value == "" {
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
	bucket := ensureMap(ensureMap(e.root, lowasync.ComponentsLabel), kind)
	if field(bucket, name) != nil {
		name += e.delimiter + componentName(target{src.doc, ""})
	}
	base := name
	for i := 1; field(bucket, name) != nil; i++ {
		name = fmt.Sprintf("%s%s%d", base, e.delimiter, i)
	}
	set(bucket, name, &yaml.Node{Kind: yaml.MappingNode, Tag: yamlparse.MappingTag})
	path := componentsPathPrefix + kind + "/" + escape(name)
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

// referenceCandidates finds existing logical destinations, including paths
// reached through component aliases. All reference placement uses this lookup.
func (e *engine) referenceCandidates(t target, prefix string) (map[string]bool, error) {
	candidates := make(map[string]bool)
	for ancestor := t; ; {
		suffix := strings.TrimPrefix(t.pointer, ancestor.pointer)
		for _, place := range e.placements[ancestor] {
			paths, err := e.project(place + suffix)
			if err != nil {
				return nil, err
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
	return candidates, nil
}

// materializeRootTargets imports external channels and servers required by
// root operations and channels. Existing destinations, including ambiguous
// destinations, remain authoritative. Messages must belong to their channel.
func (e *engine) materializeRootTargets() error {
	imported := len(e.imports)
	for {
		changed := false
		for i := 0; i < len(e.pending); i++ {
			p := e.pending[i]
			if p.role != lowasync.ChannelsLabel && p.role != lowasync.ServersLabel || p.target.doc == e.loader.root {
				continue
			}
			rootUse, err := e.rootUse(p.path)
			if err != nil {
				return err
			}
			if !rootUse {
				continue
			}
			candidates, err := e.referenceCandidates(p.target, "/"+p.role+"/")
			if err != nil {
				return err
			}
			if len(candidates) != 0 {
				continue
			}
			bucket := ensureMap(e.root, p.role)
			name := componentName(p.target)
			if field(bucket, name) != nil {
				name += e.delimiter + componentName(target{p.target.doc, ""})
			}
			base := name
			for j := 1; field(bucket, name) != nil; j++ {
				name = fmt.Sprintf("%s%s%d", base, e.delimiter, j)
			}
			path := "/" + p.role + "/" + escape(name)
			var value *yaml.Node
			if e.composed {
				component, err := e.compose(p.target, p.role, 0)
				if err != nil {
					return err
				}
				value = reference(&yaml.Node{}, localRef(component))
				e.aliases[component] = append(e.aliases[component], path)
			} else {
				value, err = e.walk(p.target, path, p.role, 0)
				if err != nil {
					return err
				}
			}
			set(bucket, name, value)
			// New imports may contain references used by the next pending edge.
			for ; imported < len(e.imports); imported++ {
				item := e.imports[imported]
				set(ensureMap(ensureMap(e.root, lowasync.ComponentsLabel), item.kind), item.name, item.node)
			}
			changed = true
		}
		if !changed {
			return nil
		}
	}
}

func (e *engine) finishReferences() error {
	// Channel/server targets must settle before message membership is checked.
	for _, messages := range []bool{false, true} {
		for _, pending := range e.pending {
			if (pending.role == lowasync.MessagesLabel) != messages {
				continue
			}
			prefix := "/" + pending.role + "/"
			if messages {
				container := pending.path[:strings.LastIndex(pending.path, messagesPathSegment)]
				channel, err := pointer(e.root, container+channelRefPath)
				if err != nil {
					return fmt.Errorf("#%s: messages require a channel", pending.path)
				}
				u, err := fragment(channel.Value)
				if err != nil {
					return err
				}
				prefix = u + messagesPathSegment
			}
			candidates, err := e.referenceCandidates(pending.target, prefix)
			if err != nil {
				return err
			}
			// Component operations/channels can reference reusable objects too.
			rootUse, err := e.rootUse(pending.path)
			if err != nil {
				return err
			}
			if len(candidates) == 0 && !rootUse && strings.HasPrefix(pending.path, componentsPathPrefix) && !messages {
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
				field(pending.node, lowasync.RefLabel).Value = localRef(path)
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
		if strings.HasPrefix(p, operationsPathPrefix) || strings.HasPrefix(p, channelsPathPrefix) {
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
			if r := field(n, lowasync.RefLabel); r != nil && r.Kind == yaml.ScalarNode && r.Tag == yamlparse.StringTag {
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
