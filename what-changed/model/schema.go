package model

import (
	"github.com/pb33f/go-yaml"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/internal/schemaformat"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// compareSchema preserves native schema semantics and unwraps Multi Format
// Schemas. Foreign schema content is opaque and uses the multiFormatSchema rule.
func compareSchema(left, right *base.SchemaProxy, config *BreakingRulesConfig) *wcmodel.SchemaChanges {
	if left.IsReference() && right.IsReference() && left.GetReference() != right.GetReference() {
		return wcmodel.CompareSchemas(left, right)
	}
	lf, rf := lowasync.SchemaFormatNode(left.GetValueNode()), lowasync.SchemaFormatNode(right.GetValueNode())
	if lf == nil && rf == nil {
		return wcmodel.CompareSchemas(left, right)
	}
	var changes []*Change
	if lf == nil || rf == nil || lf.Value != rf.Value {
		kind := Modified
		if lf == nil {
			kind = PropertyAdded
		} else if rf == nil {
			kind = PropertyRemoved
		}
		CreateChange(&changes, kind, lowasync.SchemaFormatLabel, lf, rf, schemaFormatBreaking(kind, config), left, right)
	}
	ls, li := lowasync.SchemaContent(left)
	rs, ri := lowasync.SchemaContent(right)
	jsonLeft, jsonRight := lf == nil || schemaformat.IsJSON(lf.Value), rf == nil || schemaformat.IsJSON(rf.Value)
	var result *wcmodel.SchemaChanges
	if jsonLeft && jsonRight && ls != nil && rs != nil && ls.Kind == yaml.MappingNode && rs.Kind == yaml.MappingNode {
		lp, rp := left, right
		if lf != nil {
			lp = lowasync.SchemaContentProxy(left)
		}
		if rf != nil {
			rp = lowasync.SchemaContentProxy(right)
		}
		result = wcmodel.CompareSchemas(lp, rp)
	} else if (ls == nil) != (rs == nil) || (ls != nil && rs != nil && !low.CompareYAMLNodes(ls, rs)) {
		kind := Modified
		if ls == nil {
			kind = PropertyAdded
		} else if rs == nil {
			kind = PropertyRemoved
		}
		CreateChangeWithEncoding(&changes, kind, lowasync.SchemaLabel, ls, rs, schemaFormatBreaking(kind, config), &schemaSource{ls, li}, &schemaSource{rs, ri})
	}
	if result == nil {
		if len(changes) == 0 {
			return nil
		}
		result = &wcmodel.SchemaChanges{PropertyChanges: NewPropertyChanges(changes)}
	} else {
		result.Changes = append(result.Changes, changes...)
	}
	return result
}

func schemaFormatBreaking(kind int, config *BreakingRulesConfig) bool {
	switch kind {
	case PropertyAdded:
		return BreakingAdded(CompMultiFormatSchema, "", config)
	case PropertyRemoved:
		return BreakingRemoved(CompMultiFormatSchema, "", config)
	default:
		return BreakingModified(CompMultiFormatSchema, "", config)
	}
}

type schemaSource struct {
	node  *yaml.Node
	index *index.SpecIndex
}

func (s *schemaSource) GetRootNode() *yaml.Node    { return s.node }
func (s *schemaSource) GetIndex() *index.SpecIndex { return s.index }
