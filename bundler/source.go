// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/libopenapi/index"
)

type source struct {
	uri  string
	root *yaml.Node
}
type target struct {
	doc     *source
	pointer string
}

func (t target) node() (*yaml.Node, error) { return pointer(t.doc.root, t.pointer) }

type loader struct {
	root      *source
	rolodex   *index.Rolodex
	documents map[string]*source
	resolved  map[target]target
	resolving int
	resources map[string]target
	scopes    map[target]string
	scanned   map[target]bool
}

func newLoader(root *yaml.Node, config *libasyncapi.DocumentConfiguration) (*loader, error) {
	base, err := filepath.Abs(config.BasePath)
	if err != nil {
		return nil, err
	}
	rootPath := config.SpecFilePath
	if rootPath == "" {
		rootPath = filepath.Join(base, "__asyncapi_entry__.yaml")
	} else if !filepath.IsAbs(rootPath) {
		rootPath = filepath.Join(base, rootPath)
	}
	uriPath := filepath.ToSlash(rootPath)
	// A drive letter belongs in the URL path, not its authority.
	if !strings.HasPrefix(uriPath, "/") {
		uriPath = "/" + uriPath
	}
	rootURI := (&url.URL{Scheme: "file", Path: uriPath}).String()
	if config.BaseURL != nil {
		baseURL := *config.BaseURL
		if !strings.HasSuffix(baseURL.Path, "/") {
			baseURL.Path += "/"
		}
		rootURI = baseURL.ResolveReference(&url.URL{Path: filepath.Base(rootPath)}).String()
	}
	cfg := index.CreateClosedAPIIndexConfig()
	cfg.SkipDocumentCheck = true
	cfg.AvoidBuildIndex = true
	cfg.AvoidCircularReferenceCheck = true
	cfg.SkipExternalRefResolution = true
	cfg.ExcludeExtensionRefs = true
	cfg.AllowFileLookup = config.AllowFileReferences || config.BasePath != ""
	cfg.AllowRemoteLookup = config.AllowRemoteReferences
	cfg.AllowUnknownExtensionContentDetection = true
	cfg.Logger = config.Logger
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	rolo := index.NewRolodex(cfg)
	// The loader owns the source tree. No typed model is built, which avoids
	// discarding unsupported schema formats and does not eagerly fetch examples.
	rolo.SetRootNode(root)
	rolo.SetRootIndex(index.NewSpecIndexWithConfig(&yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}, cfg))
	if cfg.AllowFileLookup {
		if config.LocalFS != nil {
			rolo.AddLocalFS(base, config.LocalFS)
		} else {
			volumeRoot := filepath.VolumeName(base) + string(filepath.Separator)
			rolo.AddLocalFS(volumeRoot, os.DirFS(volumeRoot))
		}
	}
	if cfg.AllowRemoteLookup {
		remote, err := index.NewRemoteFSWithConfig(cfg)
		if err != nil {
			return nil, err
		}
		if config.RemoteURLHandler != nil {
			remote.RemoteHandlerFunc = config.RemoteURLHandler
		}
		rolo.AddRemoteFS("default", remote)
	}
	src := &source{rootURI, root.Content[0]}
	l := &loader{root: src, rolodex: rolo, documents: map[string]*source{rootURI: src}, resolved: make(map[target]target), resources: make(map[string]target), scopes: make(map[target]string), scanned: make(map[target]bool)}
	if err := l.scan(target{src, ""}, "root", src.uri, 0); err != nil {
		return nil, err
	}
	return l, nil
}

// scan registers schema resources without interpreting examples or payloads in
// other languages. References are resolved only when the bundler visits them.
func (l *loader) scan(t target, kind, base string, depth int) error {
	if kind == "opaque" || kind == "schemaContent" {
		return nil
	}
	if depth > 512 {
		return fmt.Errorf("schema scope depth exceeded")
	}
	n, err := t.node()
	if err != nil {
		return err
	}
	if kind == "schemas" {
		if l.scanned[t] {
			return nil
		}
		l.scanned[t] = true
		if id := field(n, "$id"); id != nil {
			u, err := url.Parse(id.Value)
			if err != nil {
				return err
			}
			b, _ := url.Parse(base)
			base = b.ResolveReference(u).String()
			if prior, ok := l.resources[base]; ok && prior != t {
				return fmt.Errorf("duplicate schema resource %q", base)
			}
			l.resources[base] = t
		}
		l.scopes[t] = base
		if anchor := field(n, "$anchor"); anchor != nil {
			u, _ := url.Parse(base)
			u.Fragment = anchor.Value
			l.resources[u.String()] = t
		}
	}
	if n.Kind == yaml.MappingNode {
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i].Value
			child := childKind(kind, key)
			if kind == "schemas" && key == "schema" && field(n, "schemaFormat") != nil && !jsonSchemaFormat(field(n, "schemaFormat").Value) {
				child = "schemaContent"
			}
			if err := l.scan(target{t.doc, t.pointer + "/" + escape(key)}, child, base, depth+1); err != nil {
				return err
			}
		}
	} else if n.Kind == yaml.SequenceNode {
		for i := range n.Content {
			if err := l.scan(target{t.doc, fmt.Sprintf("%s/%d", t.pointer, i)}, strings.TrimPrefix(kind, "list:"), base, depth+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (l *loader) schemaReference(from target, ref string) (target, error) {
	base := l.scopes[from]
	if base == "" {
		base = from.doc.uri
	}
	u, err := url.Parse(ref)
	if err != nil {
		return target{}, err
	}
	b, _ := url.Parse(base)
	u = b.ResolveReference(u)
	if t, ok := l.resources[u.String()]; ok {
		return t, nil
	}
	frag := u.Fragment
	u.Fragment = ""
	u.RawFragment = ""
	if resource, ok := l.resources[u.String()]; ok && (frag == "" || strings.HasPrefix(frag, "/")) {
		t := target{resource.doc, resource.pointer + frag}
		_, err := t.node()
		return t, err
	}
	// Load the containing resource before selecting a fragment. Its root $id
	// also defines the scope of schemas reached through a JSON pointer.
	document, err := l.resolve(from, u.String())
	if err != nil {
		return target{}, err
	}
	kind := "schemas"
	if field(document.doc.root, "asyncapi") != nil {
		kind = "root"
	}
	if err := l.scan(target{document.doc, ""}, kind, document.doc.uri, 0); err != nil {
		return target{}, err
	}
	if frag != "" && !strings.HasPrefix(frag, "/") {
		base := l.scopes[target{document.doc, ""}]
		if base == "" {
			base = document.doc.uri
		}
		anchor, _ := url.Parse(base)
		anchor.Fragment = frag
		if resolved, ok := l.resources[anchor.String()]; ok {
			return resolved, nil
		}
		return target{}, fmt.Errorf("schema anchor %q not found", frag)
	}
	u.Fragment = frag
	return l.resolve(from, u.String())
}

func (l *loader) resolve(from target, ref string) (target, error) {
	l.resolving++
	defer func() { l.resolving-- }()
	if l.resolving > 128 {
		return target{}, fmt.Errorf("reference ancestry cycle resolving %q", ref)
	}
	u, err := url.Parse(ref)
	if err != nil {
		return target{}, fmt.Errorf("invalid reference %q: %w", ref, err)
	}
	base, _ := url.Parse(from.doc.uri)
	u = base.ResolveReference(u)
	ptr := u.Fragment
	if ptr != "" && !strings.HasPrefix(ptr, "/") {
		return target{}, fmt.Errorf("unsupported reference anchor %q", ref)
	}
	u.Fragment = ""
	u.RawFragment = ""
	uri := u.String()
	doc := l.documents[uri]
	if doc == nil {
		location := uri
		switch u.Scheme {
		case "file":
			if u.Host != "" && u.Host != "localhost" {
				return target{}, fmt.Errorf("unsupported file host in %q", ref)
			}
			path := u.Path
			if strings.HasPrefix(path, "/") && filepath.VolumeName(path[1:]) != "" {
				path = path[1:]
			}
			location = filepath.FromSlash(path)
		case "http", "https":
		default:
			return target{}, fmt.Errorf("unsupported reference scheme %q", u.Scheme)
		}
		file, err := l.rolodex.Open(location)
		if err != nil {
			return target{}, fmt.Errorf("load %q: %w", uri, err)
		}
		if file == nil {
			return target{}, fmt.Errorf("empty reference document %q", uri)
		}
		n, err := parse([]byte(file.GetContent()))
		if err != nil {
			return target{}, fmt.Errorf("parse %q: %w", uri, err)
		}
		if n == nil || len(n.Content) == 0 {
			return target{}, fmt.Errorf("empty reference document %q", uri)
		}
		if n.Kind == yaml.DocumentNode {
			n = n.Content[0]
		}
		doc = &source{uri, n}
		l.documents[uri] = doc
	}
	t := target{doc, ""}
	if ptr == "" {
		return t, nil
	}
	for _, part := range strings.Split(ptr[1:], "/") {
		n, err := t.node()
		if err != nil {
			return target{}, err
		}
		if r := field(n, "$ref"); r != nil && r.Kind == yaml.ScalarNode && r.Tag == "!!str" {
			t, err = l.terminal(t, false)
			if err != nil {
				return target{}, err
			}
		}
		t.pointer += "/" + part
		if _, err = t.node(); err != nil {
			return target{}, fmt.Errorf("resolve %q: %w", ref, err)
		}
	}
	return t, nil
}

// bareOnly preserves metadata inside foreign schemas. Those results must not
// share the cache used when every $ref field is interpreted as a reference.
func (l *loader) terminal(t target, bareOnly bool) (target, error) {
	if cached, ok := l.resolved[t]; ok && !bareOnly {
		return cached, nil
	}
	seen := make(map[target]bool)
	for {
		if seen[t] {
			return target{}, fmt.Errorf("reference-only cycle at %s#%s", t.doc.uri, t.pointer)
		}
		seen[t] = true
		n, err := t.node()
		if err != nil {
			return target{}, err
		}
		r := field(n, "$ref")
		if bareOnly && !referenceOnly(n) {
			r = nil
		}
		if r == nil {
			if !bareOnly {
				for start := range seen {
					l.resolved[start] = t
				}
			}
			return t, nil
		}
		if r.Kind != yaml.ScalarNode || r.Tag != "!!str" || r.Value == "" {
			return target{}, fmt.Errorf("invalid $ref at %s#%s", t.doc.uri, t.pointer)
		}
		t, err = l.resolve(t, r.Value)
		if err != nil {
			return target{}, err
		}
	}
}
