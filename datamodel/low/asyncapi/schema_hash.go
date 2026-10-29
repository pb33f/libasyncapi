package asyncapi

import (
	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
)

// SchemaHash includes Multi Format Schema content that a JSON Schema proxy does
// not model. Native schemas retain libopenapi's reference and hashing behavior.
func SchemaHash(proxy *base.SchemaProxy) string {
	if proxy != nil {
		root := proxy.GetValueNode()
		if format := SchemaFormatNode(root); format != nil {
			hash := format.Value
			if !schemaformat.IsJSON(format.Value) {
				content, _ := SchemaContent(proxy)
				return hash + low.GenerateHashString(content)
			}
			if schemaformat.IsJSON(format.Value) {
				if nested := SchemaContentProxy(proxy); nested != nil {
					hash += low.GenerateHashString(nested)
				}
			}
			return hash
		}
	}
	return low.GenerateHashString(proxy)
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
			content = root.Content[i+1]
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
		ref, next := idx.SearchIndexForReference(content.Content[1].Value)
		if ref == nil || ref.Node == nil || next == nil {
			break
		}
		content, idx = ref.Node, next
	}
	return content, idx
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
