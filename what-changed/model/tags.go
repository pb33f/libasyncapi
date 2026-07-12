// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"fmt"

	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// TagChanges represents changes made to a single AsyncAPI Tag object.
type TagChanges struct {
	*PropertyChanges
	ExternalDocChanges *ExternalDocChanges `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	ExtensionChanges   *ExtensionChanges   `json:"extensions,omitempty" yaml:"extensions,omitempty"`
}

// GetAllChanges returns a slice of all changes made between Tag objects.
func (t *TagChanges) GetAllChanges() []*Change {
	if t == nil {
		return nil
	}
	var changes []*Change
	changes = append(changes, t.Changes...)
	if t.ExternalDocChanges != nil {
		changes = append(changes, t.ExternalDocChanges.GetAllChanges()...)
	}
	if t.ExtensionChanges != nil {
		changes = append(changes, t.ExtensionChanges.GetAllChanges()...)
	}
	return changes
}

// TotalChanges returns a count of everything that changed.
func (t *TagChanges) TotalChanges() int {
	if t == nil {
		return 0
	}
	c := t.PropertyChanges.TotalChanges()
	if t.ExternalDocChanges != nil {
		c += t.ExternalDocChanges.TotalChanges()
	}
	if t.ExtensionChanges != nil {
		c += t.ExtensionChanges.TotalChanges()
	}
	return c
}

// TotalBreakingChanges returns the number of breaking changes made.
func (t *TagChanges) TotalBreakingChanges() int {
	if t == nil {
		return 0
	}
	c := t.PropertyChanges.TotalBreakingChanges()
	if t.ExternalDocChanges != nil {
		c += t.ExternalDocChanges.TotalBreakingChanges()
	}
	if t.ExtensionChanges != nil {
		c += t.ExtensionChanges.TotalBreakingChanges()
	}
	return c
}

// CompareTags compares two AsyncAPI Tag objects and returns a pointer to TagChanges,
// or nil if nothing changed.
func CompareTags(l, r *lowasync.Tag, configs ...*BreakingRulesConfig) *TagChanges {
	config := comparisonConfig(configs)
	if l == nil || r == nil {
		return nil
	}
	if low.AreEqual(l, r) {
		return nil
	}

	var changes []*Change
	props := []*PropertyCheck{
		NewPropertyCheck(CompTag, PropName,
			l.Name.ValueNode, r.Name.ValueNode, lowasync.NameLabel, &changes, l, r, config),
		NewPropertyCheck(CompTag, PropDescription,
			l.Description.ValueNode, r.Description.ValueNode, lowasync.DescriptionLabel, &changes, l, r, config),
	}
	CheckProperties(props, config)

	tc := new(TagChanges)
	compareNestedObject(l.ExternalDocs, r.ExternalDocs, lowasync.ExternalDocsLabel,
		CompTag, PropExternalDocs, &changes, configuredNestedCompare(config, CompareExternalDocs), &tc.ExternalDocChanges, config)
	tc.ExtensionChanges = CheckExtensions(l, r)
	tc.PropertyChanges = NewPropertyChanges(changes)
	if tc.TotalChanges() <= 0 {
		return nil
	}
	return tc
}

// CompareTagSlices compares two slices of Tag value references by name and occurrence.
// Added and removed tags are recorded against the supplied component and property rules,
// and modified tags are compared with CompareTags. Repeated and unnamed tags retain every
// occurrence instead of collapsing in a map. Modified duplicates use a stable "name#N"
// result key; unique names keep their original key.
func CompareTagSlices(l, r []low.ValueReference[*lowasync.Tag],
	changes *[]*Change, component, property string,
	configs ...*BreakingRulesConfig,
) map[string]*TagChanges {
	config := comparisonConfig(configs)
	if len(l) == 0 && len(r) == 0 {
		return nil
	}
	breakingAdded := BreakingAdded(component, property, config)
	breakingRemoved := BreakingRemoved(component, property, config)

	lTags, lOrder := groupTagsByName(l)
	rTags, rOrder := groupTagsByName(r)

	tagChanges := make(map[string]*TagChanges)
	for _, name := range lOrder {
		left, right := lTags[name], rTags[name]
		common := min(len(left), len(right))
		for i := 0; i < common; i++ {
			lt, rt := left[i], right[i]
			switch {
			case lt.Value == nil && rt.Value == nil:
				continue
			case lt.Value == nil:
				CreateChange(changes, ObjectAdded, lowasync.TagsLabel,
					nil, rt.ValueNode, breakingAdded, nil, rt.Value)
			case rt.Value == nil:
				CreateChange(changes, ObjectRemoved, lowasync.TagsLabel,
					lt.ValueNode, nil, breakingRemoved, lt.Value, nil)
			default:
				if tc := CompareTags(lt.Value, rt.Value, config); tc != nil {
					key := tagOccurrenceKey(name, i, len(left), len(right), tagChanges)
					tagChanges[key] = tc
				}
			}
		}
		for i := common; i < len(left); i++ {
			lt := left[i]
			CreateChange(changes, ObjectRemoved, lowasync.TagsLabel,
				lt.ValueNode, nil, breakingRemoved, lt.Value, nil)
		}
		for i := common; i < len(right); i++ {
			rt := right[i]
			CreateChange(changes, ObjectAdded, lowasync.TagsLabel,
				nil, rt.ValueNode, breakingAdded, nil, rt.Value)
		}
	}
	for _, name := range rOrder {
		if _, seenLeft := lTags[name]; seenLeft {
			continue
		}
		for _, rt := range rTags[name] {
			CreateChange(changes, ObjectAdded, lowasync.TagsLabel,
				nil, rt.ValueNode, breakingAdded, nil, rt.Value)
		}
	}
	if len(tagChanges) == 0 {
		return nil
	}
	return tagChanges
}

func groupTagsByName(tags []low.ValueReference[*lowasync.Tag]) (
	map[string][]low.ValueReference[*lowasync.Tag], []string,
) {
	grouped := make(map[string][]low.ValueReference[*lowasync.Tag], len(tags))
	order := make([]string, 0, len(tags))
	for _, tag := range tags {
		name := ""
		if tag.Value != nil {
			name = tag.Value.Name.Value
		}
		if _, exists := grouped[name]; !exists {
			order = append(order, name)
		}
		grouped[name] = append(grouped[name], tag)
	}
	return grouped, order
}

func tagOccurrenceKey(name string, occurrence, leftCount, rightCount int,
	existing map[string]*TagChanges,
) string {
	if leftCount == 1 && rightCount == 1 {
		if _, used := existing[name]; !used {
			return name
		}
	}
	base := name
	if base == "" {
		base = "<unnamed>"
	}
	for suffix := occurrence + 1; ; suffix++ {
		key := fmt.Sprintf("%s#%d", base, suffix)
		if _, used := existing[key]; !used {
			return key
		}
	}
}
