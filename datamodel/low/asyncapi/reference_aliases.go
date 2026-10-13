// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/index"
)

type referenceLocation struct {
	index *index.SpecIndex
	ref   string
}

// repairReferenceAliases registers pointers that pass through Reference Objects.
// It preserves source nodes and leaves the index's raw diagnostics available.
// Only errors for references proven to resolve are removed from Document.Errors.
func repairReferenceAliases(rolodex *index.Rolodex) map[*index.IndexingError]bool {
	resolved := make(map[*index.IndexingError]bool)
	for _, idx := range append(rolodex.GetIndexes(), rolodex.GetRootIndex()) {
		if idx == nil {
			continue
		}
		for _, raw := range idx.GetReferenceIndexErrors() {
			issue, ok := raw.(*index.IndexingError)
			if !ok || issue.KeyNode == nil {
				continue
			}
			ref := issue.KeyNode.Value
			target := resolveThroughAliases(idx, ref, make(map[referenceLocation]bool))
			if target == nil {
				continue
			}
			alias := *target
			alias.Definition, alias.FullDefinition = ref, ref
			mapped := idx.GetMappedReferences()
			if mapped == nil {
				mapped = make(map[string]*index.Reference)
			}
			mapped[ref] = &alias
			idx.SetMappedReferences(mapped)
			resolved[issue] = true
		}
	}
	return resolved
}

func removeRepairedAliasErrors(errs []error, resolved map[*index.IndexingError]bool) []error {
	var prune func(error) error
	prune = func(err error) error {
		if issue, ok := err.(*index.IndexingError); ok && resolved[issue] {
			return nil
		}
		if joined, ok := err.(interface{ Unwrap() []error }); ok {
			var remaining []error
			for _, child := range joined.Unwrap() {
				if child = prune(child); child != nil {
					remaining = append(remaining, child)
				}
			}
			return errors.Join(remaining...)
		}
		return err
	}
	var remaining []error
	for _, err := range errs {
		if err = prune(err); err != nil {
			remaining = append(remaining, err)
		}
	}
	return remaining
}

func resolveThroughAliases(idx *index.SpecIndex, ref string, active map[referenceLocation]bool) *index.Reference {
	key := referenceLocation{idx, ref}
	if active[key] || len(active) > 512 {
		return nil
	}
	active[key] = true
	defer delete(active, key)
	base, fragment, ok := strings.Cut(ref, "#")
	if !ok {
		return nil
	}
	pointer, err := url.PathUnescape(fragment)
	if err != nil || !strings.HasPrefix(pointer, "/") {
		return nil
	}
	parts := strings.Split(pointer[1:], "/")
	for _, part := range parts {
		for i := 0; i < len(part); i++ {
			if part[i] == '~' {
				if i+1 == len(part) || (part[i+1] != '0' && part[i+1] != '1') {
					return nil
				}
				i++
			}
		}
	}
	for cut := 1; cut < len(parts); cut++ {
		prefix := base + (&url.URL{Fragment: "/" + strings.Join(parts[:cut], "/")}).String()
		var target *index.Reference
		owner := idx
		if base == "" {
			target = idx.FindComponentInRoot(context.Background(), prefix)
		} else {
			target, owner, _ = idx.SearchIndexForReferenceWithContext(context.Background(), prefix)
		}
		if target == nil {
			continue
		}
		if owner == nil {
			owner = idx
		}
		node := target.Node
		valid := true
		seen := make(map[*yaml.Node]bool)
		for _, token := range parts[cut:] {
			for node != nil {
				next := referenceField(node)
				if next == "" {
					break
				}
				if seen[node] {
					node = nil
					break
				}
				seen[node] = true
				resolved, nextOwner, _ := owner.SearchIndexForReferenceWithContext(context.Background(), next)
				if resolved == nil {
					resolved = resolveThroughAliases(owner, next, active)
					if resolved != nil {
						nextOwner = resolved.Index
					}
				}
				if resolved == nil {
					node = nil
					break
				}
				node = resolved.Node
				if nextOwner != nil {
					owner = nextOwner
				}
			}
			token = strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
			node = pointerChild(node, token)
			if node == nil || !valid {
				valid = false
				break
			}
		}
		if valid {
			return &index.Reference{Node: node, Index: owner, Resolved: true}
		}
	}
	return nil
}

func pointerChild(node *yaml.Node, token string) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.MappingNode {
		for i := 0; i+1 < len(node.Content); i += 2 {
			if node.Content[i].Value == token {
				return node.Content[i+1]
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		i, err := strconv.Atoi(token)
		if err == nil && i >= 0 && i < len(node.Content) && strconv.Itoa(i) == token {
			return node.Content[i]
		}
	}
	return nil
}

func referenceField(node *yaml.Node) string {
	ref := pointerChild(node, "$ref")
	if ref != nil && ref.Kind == yaml.ScalarNode && ref.Tag == "!!str" {
		return ref.Value
	}
	return ""
}
