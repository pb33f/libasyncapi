package asyncapi

import (
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/utils"
)

// SchemaHash includes Multi Format Schema content that a JSON Schema proxy does
// not model. Quick hashing a reference hashes only its target, so the reference
// itself is added to keep a move between identical targets visible.
func SchemaHash(proxy *base.SchemaProxy) string {
	if proxy == nil {
		return low.GenerateHashString(proxy)
	}
	hash := proxy.GetReference()
	format := SchemaFormatNode(proxy.GetValueNode())
	if format == nil {
		return hash + low.GenerateHashString(proxy)
	}
	hash += format.Value
	content, _ := SchemaContent(proxy)
	if schemaformat.IsJSON(format.Value) && content != nil && content.Kind == yaml.MappingNode {
		if nested := SchemaContentProxy(proxy); nested != nil {
			hash += low.GenerateHashString(nested)
		}
		return hash
	}
	return hash + low.GenerateHashString(content)
}

type schemaContentKey struct{ node *yaml.Node }

// SchemaContentProxy reuses the parsed inner JSON Schema within its owning
// index. Index cleanup releases it; no process-wide wrapper cache is retained.
func SchemaContentProxy(proxy *base.SchemaProxy) *base.SchemaProxy {
	if SchemaFormatNode(proxy.GetValueNode()) == nil {
		return proxy
	}
	key := schemaContentKey{proxy.GetValueNode()}
	idx := proxy.GetIndex()
	if idx != nil && idx.GetHighCache() != nil {
		if value, ok := idx.GetHighCache().Load(key); ok {
			return value.(*base.SchemaProxy)
		}
	}
	node, owner := SchemaContent(proxy)
	if node == nil {
		return nil
	}
	nested := new(base.SchemaProxy)
	if nested.Build(proxy.GetContext(), nil, node, owner) != nil {
		return nil
	}
	if idx != nil && idx.GetHighCache() != nil {
		idx.GetHighCache().Store(key, nested)
	}
	return nested
}

// SchemaContent returns the wrapper's content and owning index. Only a bare
// schema reference is resolved; reference-like keys inside opaque data stay data.
func SchemaContent(proxy *base.SchemaProxy) (*yaml.Node, *index.SpecIndex) {
	root, idx := proxy.GetValueNode(), proxy.GetIndex()
	if SchemaFormatNode(root) == nil {
		return root, idx
	}
	var content *yaml.Node
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == SchemaLabel {
			content = utils.NodeAlias(root.Content[i+1])
			break
		}
	}
	var seen map[*yaml.Node]bool
	for content != nil && content.Kind == yaml.MappingNode && len(content.Content) == 2 && content.Content[0].Value == RefLabel && idx != nil {
		if seen == nil {
			seen = make(map[*yaml.Node]bool)
		}
		if seen[content] {
			break
		}
		seen[content] = true
		ref, next := searchReference(idx, content.Content[1])
		if ref == nil || ref.Node == nil || next == nil {
			break
		}
		content, idx = ref.Node, next
	}
	return content, idx
}

// searchReference resolves a $ref value from the location recorded while indexing.
// Repaired aliases are registered by their raw reference, so the raw value remains
// the fallback when the indexed location finds nothing.
func searchReference(idx *index.SpecIndex, value *yaml.Node) (*index.Reference, *index.SpecIndex) {
	if ref := indexedReference(idx, value); ref != nil {
		if found, owner := idx.SearchIndexForReferenceByReference(ref); found != nil && found.Node != nil {
			return found, owner
		}
	}
	return idx.SearchIndexForReference(value.Value)
}

type referenceSitesKey struct{}

// indexedReference finds the index entry for a $ref value. A relative value is
// only meaningful to the file that holds it, so searching for the raw string can
// resolve it against the wrong base. Sites are cached on the index because schema
// hashing looks one up for every wrapped reference.
func indexedReference(idx *index.SpecIndex, value *yaml.Node) *index.Reference {
	cache := idx.GetHighCache()
	if cache != nil {
		if sites, ok := cache.Load(referenceSitesKey{}); ok {
			return sites.(map[*yaml.Node]*index.Reference)[value]
		}
	}
	refs := idx.GetRawReferencesSequenced()
	sites := make(map[*yaml.Node]*index.Reference, len(refs))
	for _, ref := range refs {
		if _, ok := sites[ref.KeyNode]; !ok {
			sites[ref.KeyNode] = ref
		}
	}
	if cache != nil {
		cache.Store(referenceSitesKey{}, sites)
	}
	return sites[value]
}

// SchemaFormatNode returns the format field of a Multi Format Schema wrapper.
func SchemaFormatNode(root *yaml.Node) *yaml.Node {
	if root == nil || root.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == SchemaFormatLabel {
			return root.Content[i+1]
		}
	}
	return nil
}
