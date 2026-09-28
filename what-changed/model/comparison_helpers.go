// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"reflect"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
	wcmodel "github.com/pb33f/libopenapi/what-changed/model"
)

// IMPORTANT: the rule-aware helpers in this file are deliberate local re-implementations
// of their libopenapi counterparts. The imported wcmodel versions resolve breaking status
// against libopenapi's package-global (OpenAPI shaped) breaking-rules configuration, which
// knows nothing about AsyncAPI components — so AsyncAPI rules would be silently ignored
// and colliding key names (info.title, tag.name, securityScheme.type) would resolve to
// OpenAPI rules. The versions here consult THIS package's active configuration, and lean
// on the config-free wcmodel primitives (CheckForRemoval/Addition/Modification,
// CreateChange, SetReferenceIfExists) for the mechanics.

func configuredCompare[T, R any](config *BreakingRulesConfig,
	compare func(l, r T, configs ...*BreakingRulesConfig) R,
) func(l, r T) R {
	return func(l, r T) R {
		return compare(l, r, config)
	}
}

func configuredNestedCompare[T, R any](config *BreakingRulesConfig,
	compare func(l, r T, configs ...*BreakingRulesConfig) *R,
) func(l, r T) *R {
	return func(l, r T) *R {
		return compare(l, r, config)
	}
}

// NewPropertyCheck creates a new PropertyCheck for the given component and property,
// resolving breaking status from this package's active AsyncAPI breaking-rules config.
func NewPropertyCheck(component, property string, leftNode, rightNode *yaml.Node,
	label string, changes *[]*Change, original, updated any, configs ...*BreakingRulesConfig,
) *PropertyCheck {
	return &PropertyCheck{
		LeftNode:  leftNode,
		RightNode: rightNode,
		Label:     label,
		Changes:   changes,
		Breaking:  BreakingModified(component, property, configs...),
		Component: component,
		Property:  property,
		Original:  original,
		New:       updated,
	}
}

// CheckProperties checks a slice of PropertyCheck objects for additions, removals and
// modifications, resolving per change-type breaking status from this package's active
// AsyncAPI breaking-rules configuration.
func CheckProperties(properties []*PropertyCheck, configs ...*BreakingRulesConfig) {
	// cache config once outside the loop (avoids repeated mutex operations)
	config := comparisonConfig(configs)

	for _, n := range properties {
		var breakingAdded, breakingModified, breakingRemoved bool

		if n.Component != "" {
			if r := config.GetRule(n.Component, n.Property); r != nil {
				breakingAdded = r.Added != nil && *r.Added
				breakingModified = r.Modified != nil && *r.Modified
				breakingRemoved = r.Removed != nil && *r.Removed
			} else {
				breakingAdded = n.Breaking
				breakingModified = n.Breaking
				breakingRemoved = n.Breaking
			}
		} else {
			breakingAdded = n.Breaking
			breakingModified = n.Breaking
			breakingRemoved = n.Breaking
		}

		wcmodel.CheckForRemoval(n.LeftNode, n.RightNode, n.Label, n.Changes, breakingRemoved, n.Original, n.New)
		wcmodel.CheckForAddition(n.LeftNode, n.RightNode, n.Label, n.Changes, breakingAdded, n.Original, n.New)
		wcmodel.CheckForModification(n.LeftNode, n.RightNode, n.Label, n.Changes, breakingModified, n.Original, n.New)
	}
}

// CheckMapForChangesWithRules checks a left and right low level map for any additions,
// subtractions or modifications to values, resolving breaking status for additions and
// removals from this package's active AsyncAPI breaking-rules configuration.
func CheckMapForChangesWithRules[T any, R any](expLeft, expRight *orderedmap.Map[low.KeyReference[string], low.ValueReference[T]],
	changes *[]*Change, label string, compareFunc func(l, r T) R, component, property string,
	configs ...*BreakingRulesConfig,
) map[string]R {
	config := comparisonConfig(configs)
	return checkMapForChangesInternal(expLeft, expRight, changes, label, compareFunc, true,
		BreakingAdded(component, property, config), BreakingRemoved(component, property, config))
}

// checkMapForChangesInternal is the core implementation that checks a left and right low
// level map for additions, subtractions or modifications. The breakingAdded and
// breakingRemoved parameters control whether additions and removals are marked breaking.
func checkMapForChangesInternal[T any, R any](expLeft, expRight *orderedmap.Map[low.KeyReference[string], low.ValueReference[T]],
	changes *[]*Change, label string, compareFunc func(l, r T) R, compare bool,
	breakingAdded, breakingRemoved bool,
) map[string]R {
	lHashes := make(map[string]string)
	rHashes := make(map[string]string)
	lValues := make(map[string]low.ValueReference[T])
	rValues := make(map[string]low.ValueReference[T])
	lNodes := make(map[string]*yaml.Node)
	rNodes := make(map[string]*yaml.Node)
	var lOrder, rOrder []string

	if expLeft != nil {
		for k, v := range expLeft.FromOldest() {
			lOrder = append(lOrder, k.Value)
			lValues[k.Value] = v
			lNodes[k.Value] = changeNode(v.ValueNode, k.KeyNode)
			if !isNil(v.Value) {
				lHashes[k.Value] = low.GenerateHashString(v.Value)
			}
		}
	}

	if expRight != nil {
		for k, v := range expRight.FromOldest() {
			rOrder = append(rOrder, k.Value)
			rValues[k.Value] = v
			rNodes[k.Value] = changeNode(v.ValueNode, k.KeyNode)
			if !isNil(v.Value) {
				rHashes[k.Value] = low.GenerateHashString(v.Value)
			}
		}
	}

	expChanges := make(map[string]R)
	for _, k := range lOrder {
		lv := lValues[k]
		rv, existsRight := rValues[k]
		if !existsRight {
			CreateChange(changes, ObjectRemoved, label,
				lNodes[k], nil, breakingRemoved, lv.Value, nil)
			continue
		}

		lNil, rNil := isNil(lv.Value), isNil(rv.Value)
		switch {
		case lNil && rNil:
			continue
		case lNil:
			CreateChange(changes, ObjectAdded, label,
				nil, rNodes[k], breakingAdded, nil, rv.Value)
			continue
		case rNil:
			CreateChange(changes, ObjectRemoved, label,
				lNodes[k], nil, breakingRemoved, lv.Value, nil)
			continue
		}

		if lHashes[k] == rHashes[k] || !compare {
			continue
		}

		ch := compareFunc(lv.Value, rv.Value)
		if !reflect.ValueOf(&ch).Elem().IsZero() {
			expChanges[k] = ch
			var cr any = ch
			wcmodel.SetReferenceIfExists(&lv, cr)
		}
	}

	for _, k := range rOrder {
		if _, existsLeft := lValues[k]; existsLeft {
			continue
		}
		rv := rValues[k]
		CreateChange(changes, ObjectAdded, label,
			nil, rNodes[k], breakingAdded, nil, rv.Value)
	}
	return expChanges
}

// changeNode returns the value node when it contains a useful scalar value, otherwise it
// falls back to the map key node. Mapping nodes have an empty Value field, so using the key
// gives additions and removals a useful Original/New value without mutating parsed input.
func changeNode(valueNode, keyNode *yaml.Node) *yaml.Node {
	if valueNode != nil && valueNode.Value != "" {
		return valueNode
	}
	if keyNode != nil {
		return keyNode
	}
	return valueNode
}

// isNil reports whether a generic value is nil, handling typed nil pointers that would
// otherwise compare non-nil when boxed into an interface.
func isNil[T any](v T) bool {
	rv := reflect.ValueOf(v)
	return !rv.IsValid() || (rv.Kind() == reflect.Ptr && rv.IsNil())
}

// compareNestedObject handles the recurring pattern for singular nested objects: when both
// sides are present the comparator runs and its result lands in target; when only one side
// is present an ObjectAdded or ObjectRemoved change is recorded using the breaking rules
// for the supplied component and property.
func compareNestedObject[T any, R any](l, r low.NodeReference[T], label, component, property string,
	changes *[]*Change, compareFunc func(l, r T) *R, target **R,
	configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	lNil := isNil(l.Value)
	rNil := isNil(r.Value)
	switch {
	case !lNil && !rNil:
		*target = compareFunc(l.Value, r.Value)
	case lNil && !rNil:
		CreateChange(changes, ObjectAdded, label, nil, r.ValueNode,
			BreakingAdded(component, property, config), nil, r.Value)
	case !lNil && rNil:
		CreateChange(changes, ObjectRemoved, label, l.ValueNode, nil,
			BreakingRemoved(component, property, config), l.Value, nil)
	}
}

// refValueNode returns the yaml node holding the reference string itself, so changes
// render with the actual $ref value in Original/New rather than an empty mapping value.
// Falls back to the supplied node when the reference carries no node.
func refValueNode(ref *low.Reference, fallback *yaml.Node) *yaml.Node {
	if ref != nil {
		if n := ref.GetReferenceNode(); n != nil {
			return n
		}
	}
	return fallback
}

// compareSingleRef compares two reference-only fields (like an operation's channel) by
// their reference strings. AsyncAPI stores these as raw $ref pointers, so there is no
// resolved object to recurse into — a different target is a modification.
func compareSingleRef(l, r low.NodeReference[*low.Reference], label, component, property string,
	changes *[]*Change, configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	var lRef, rRef string
	if l.Value != nil {
		lRef = l.Value.GetReference()
	}
	if r.Value != nil {
		rRef = r.Value.GetReference()
	}
	switch {
	case lRef == rRef:
		return
	case lRef == "":
		CreateChange(changes, PropertyAdded, label, nil, refValueNode(r.Value, r.ValueNode),
			BreakingAdded(component, property, config), nil, r.Value)
	case rRef == "":
		CreateChange(changes, PropertyRemoved, label, refValueNode(l.Value, l.ValueNode), nil,
			BreakingRemoved(component, property, config), l.Value, nil)
	default:
		CreateChange(changes, Modified, label,
			refValueNode(l.Value, l.ValueNode), refValueNode(r.Value, r.ValueNode),
			BreakingModified(component, property, config), l.Value, r.Value)
	}
}

// compareRefSlices compares two slices of reference-only entries (like a channel's servers
// or an operation's messages) as sets of reference strings, recording additions and
// removals using the breaking rules for the supplied component and property. Duplicate
// references collapse to a single set entry.
func compareRefSlices(l, r []low.ValueReference[*low.Reference], label, component, property string,
	changes *[]*Change, configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	lRefs := make(map[string]low.ValueReference[*low.Reference], len(l))
	rRefs := make(map[string]low.ValueReference[*low.Reference], len(r))
	var lOrder, rOrder []string
	for _, v := range l {
		if v.Value != nil {
			ref := v.Value.GetReference()
			if _, exists := lRefs[ref]; !exists {
				lOrder = append(lOrder, ref)
			}
			lRefs[ref] = v
		}
	}
	for _, v := range r {
		if v.Value != nil {
			ref := v.Value.GetReference()
			if _, exists := rRefs[ref]; !exists {
				rOrder = append(rOrder, ref)
			}
			rRefs[ref] = v
		}
	}
	for _, ref := range lOrder {
		v := lRefs[ref]
		if _, ok := rRefs[ref]; !ok {
			CreateChange(changes, ObjectRemoved, label, refValueNode(v.Value, v.ValueNode), nil,
				breakingRemoved, v.Value, nil)
		}
	}
	for _, ref := range rOrder {
		v := rRefs[ref]
		if _, ok := lRefs[ref]; !ok {
			CreateChange(changes, ObjectAdded, label, nil, refValueNode(v.Value, v.ValueNode),
				breakingAdded, nil, v.Value)
		}
	}
}

// compareUnkeyedSlices compares two slices of objects that have no natural identity key
// (like traits or security entries) by hashing each entry. Entries present on only one
// side are recorded as additions or removals; an in-place edit therefore reports as one
// removal plus one addition (and inherits the removal's breaking status, even for a
// cosmetic edit). Identical duplicate entries collapse to a single set entry, so
// reducing [A, A] to [A] reports nothing.
func compareUnkeyedSlices[T any](l, r []low.ValueReference[T], label, component, property string,
	changes *[]*Change, configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	lHashes := make(map[string]low.ValueReference[T], len(l))
	rHashes := make(map[string]low.ValueReference[T], len(r))
	var lOrder, rOrder []string
	for _, v := range l {
		if !isNil(v.Value) {
			hash := low.GenerateHashString(v.Value)
			if _, exists := lHashes[hash]; !exists {
				lOrder = append(lOrder, hash)
			}
			lHashes[hash] = v
		}
	}
	for _, v := range r {
		if !isNil(v.Value) {
			hash := low.GenerateHashString(v.Value)
			if _, exists := rHashes[hash]; !exists {
				rOrder = append(rOrder, hash)
			}
			rHashes[hash] = v
		}
	}
	for _, hash := range lOrder {
		v := lHashes[hash]
		if _, ok := rHashes[hash]; !ok {
			CreateChange(changes, ObjectRemoved, label, v.ValueNode, nil,
				breakingRemoved, v.Value, nil)
		}
	}
	for _, hash := range rOrder {
		v := rHashes[hash]
		if _, ok := lHashes[hash]; !ok {
			CreateChange(changes, ObjectAdded, label, nil, v.ValueNode,
				breakingAdded, nil, v.Value)
		}
	}
}

// compareIndexedSlices compares two slices of objects positionally: entries sharing an
// index are compared with compareFunc, surplus entries on either side are recorded as
// object additions or removals using the breaking rules for the supplied component and
// property. Used for ordered collections without a natural identity key, like SQS
// policy statements, where position carries meaning.
func compareIndexedSlices[T any, R any](l, r []low.ValueReference[T], label, component, property string,
	changes *[]*Change, compareFunc func(l, r T) *R, configs ...*BreakingRulesConfig,
) []*R {
	config := comparisonConfig(configs)
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	var results []*R
	maxLen := max(len(l), len(r))
	for i := 0; i < maxLen; i++ {
		switch {
		case i < len(l) && i < len(r):
			if res := compareFunc(l[i].Value, r[i].Value); res != nil {
				results = append(results, res)
			}
		case i < len(l): // entry removed
			CreateChange(changes, ObjectRemoved, label, l[i].ValueNode, nil,
				breakingRemoved, l[i].Value, nil)
		default: // entry added
			CreateChange(changes, ObjectAdded, label, nil, r[i].ValueNode,
				breakingAdded, nil, r[i].Value)
		}
	}
	return results
}

// compareRawNode compares two raw yaml.Node valued properties (like message example
// payloads or SQS policy statement principals) structurally, recording a modification,
// addition or removal with encoded values for rendering.
func compareRawNode(l, r low.NodeReference[*yaml.Node], label, component, property string,
	changes *[]*Change, configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	lNil := l.Value == nil
	rNil := r.Value == nil
	switch {
	case lNil && rNil:
		return
	case lNil:
		CreateChangeWithEncoding(changes, PropertyAdded, label, nil, r.ValueNode,
			BreakingAdded(component, property, config), nil, r.Value)
	case rNil:
		CreateChangeWithEncoding(changes, PropertyRemoved, label, l.ValueNode, nil,
			BreakingRemoved(component, property, config), l.Value, nil)
	default:
		if !low.CompareYAMLNodes(l.Value, r.Value) {
			CreateChangeWithEncoding(changes, Modified, label, l.ValueNode, r.ValueNode,
				BreakingModified(component, property, config), l.Value, r.Value)
		}
	}
}

// ExtractStringValueSliceChangesWithRules compares two low level string slices for changes,
// resolving breaking status from this package's active AsyncAPI breaking-rules
// configuration. Slices are compared as value sets, so duplicate entries collapse to one
// and reordering reports nothing.
func ExtractStringValueSliceChangesWithRules(lParam, rParam []low.ValueReference[string],
	changes *[]*Change, label string, component, property string,
	configs ...*BreakingRulesConfig,
) {
	config := comparisonConfig(configs)
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	lValues := make(map[string]low.ValueReference[string], len(lParam))
	rValues := make(map[string]low.ValueReference[string], len(rParam))
	var lOrder, rOrder []string
	for i := range lParam {
		value := lParam[i].Value
		if _, exists := lValues[value]; !exists {
			lOrder = append(lOrder, value)
		}
		lValues[value] = lParam[i]
	}
	for i := range rParam {
		value := rParam[i].Value
		if _, exists := rValues[value]; !exists {
			rOrder = append(rOrder, value)
		}
		rValues[value] = rParam[i]
	}
	for _, k := range lOrder {
		if _, ok := rValues[k]; !ok {
			CreateChange(changes, PropertyRemoved, label,
				lValues[k].ValueNode, nil, breakingRemoved,
				lValues[k].Value, nil)
		}
	}
	for _, k := range rOrder {
		if _, ok := lValues[k]; !ok {
			CreateChange(changes, PropertyAdded, label,
				nil, rValues[k].ValueNode, breakingAdded,
				nil, rValues[k].Value)
		}
	}
}
