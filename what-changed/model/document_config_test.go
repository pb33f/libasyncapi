// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"sync"
	"testing"

	libasyncapi "github.com/pb33f/libasyncapi"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func comparisonDocument(t *testing.T, title string) *lowasync.AsyncAPI {
	t.Helper()
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: 3.0.0
info:
  title: ` + title + `
  version: 1.0.0
`))
	require.NoError(t, err)
	return doc.GoLow()
}

func comparisonDocumentFromYAML(t *testing.T, spec string) *lowasync.AsyncAPI {
	t.Helper()
	doc, err := libasyncapi.NewDocument([]byte(spec))
	require.NoError(t, err)
	return doc.GoLow()
}

func titleRuleConfig(breaking bool) *BreakingRulesConfig {
	config := NewDefaultBreakingRulesConfig()
	config.Merge(&BreakingRulesConfig{
		Info: &InfoRules{Title: rule(false, breaking, false)},
	})
	return config
}

func TestCompareDocumentsWithConfig_IsScopedAndRestoresGlobalRules(t *testing.T) {
	left := comparisonDocument(t, "Old")
	right := comparisonDocument(t, "New")

	SetActiveBreakingRulesConfig(titleRuleConfig(true))
	defer ResetActiveBreakingRulesConfig()

	changes := CompareDocumentsWithConfig(left, right, nil)
	require.NotNil(t, changes)
	assert.Zero(t, changes.TotalBreakingChanges(), "nil per-call config must select defaults")
	assert.True(t, BreakingModified(CompInfo, PropTitle), "global config must be restored")
}

func TestCompareDocuments_UsesActiveRulesAndHandlesNil(t *testing.T) {
	left := comparisonDocument(t, "Old")
	right := comparisonDocument(t, "New")
	assert.Nil(t, CompareDocuments(nil, right))
	assert.Nil(t, CompareDocuments(left, nil))

	SetActiveBreakingRulesConfig(titleRuleConfig(true))
	defer ResetActiveBreakingRulesConfig()
	changes := CompareDocuments(left, right)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalBreakingChanges())
}

func TestCompareDocumentsWithConfig_ConcurrentRulesDoNotCrossTalk(t *testing.T) {
	left := comparisonDocument(t, "Old")
	right := comparisonDocument(t, "New")
	breaking, nonBreaking := titleRuleConfig(true), titleRuleConfig(false)

	const comparisons = 32
	start := make(chan struct{})
	results := make(chan error, comparisons)
	var wg sync.WaitGroup
	for i := range comparisons {
		wantBreaking := i%2 == 0
		config := nonBreaking
		if wantBreaking {
			config = breaking
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			changes := CompareDocumentsWithConfig(left, right, config)
			if changes == nil || (changes.TotalBreakingChanges() == 1) != wantBreaking {
				results <- assert.AnError
			}
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	assert.Empty(t, results)
}

func TestDirectComparators_ConcurrentRulesDoNotCrossTalkOrLeakGlobally(t *testing.T) {
	left := buildInfo(t, `title: Old
version: 1.0.0`)
	right := buildInfo(t, `title: New
version: 1.0.0`)
	breaking, nonBreaking := titleRuleConfig(true), titleRuleConfig(false)

	ResetActiveBreakingRulesConfig()
	const comparisons = 64
	start := make(chan struct{})
	results := make(chan error, comparisons)
	var wg sync.WaitGroup
	for i := range comparisons {
		wantBreaking := i%2 == 0
		config := nonBreaking
		if wantBreaking {
			config = breaking
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			changes := CompareInfo(left, right, config)
			if changes == nil || (changes.TotalBreakingChanges() == 1) != wantBreaking {
				results <- assert.AnError
			}
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	assert.Empty(t, results)
	assert.False(t, BreakingModified(CompInfo, PropTitle),
		"comparison-local rules must never become process-global")
}

func TestCompareDocumentsWithConfig_PropagatesRulesThroughNestedComparators(t *testing.T) {
	left := comparisonDocumentFromYAML(t, `asyncapi: 3.0.0
info: {title: API, version: 1.0.0}
servers:
  production:
    host: old.example.com
    protocol: kafka`)
	right := comparisonDocumentFromYAML(t, `asyncapi: 3.0.0
info: {title: API, version: 1.0.0}
servers:
  production:
    host: new.example.com
    protocol: kafka`)

	config := NewDefaultBreakingRulesConfig()
	config.Merge(&BreakingRulesConfig{
		Server: &ServerRules{Host: rule(false, false, false)},
	})
	changes := CompareDocumentsWithConfig(left, right, config)
	require.NotNil(t, changes)
	assert.Equal(t, 1, changes.TotalChanges())
	assert.Zero(t, changes.TotalBreakingChanges())
}
