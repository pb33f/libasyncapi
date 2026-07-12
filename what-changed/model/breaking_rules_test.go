// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateDefaultBreakingRules_CachedAndComplete(t *testing.T) {
	a := GenerateDefaultBreakingRules()
	b := GenerateDefaultBreakingRules()
	assert.Same(t, a, b)

	// spot check defaults across categories.
	assert.True(t, a.IsBreaking(CompAsyncAPI, "", ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompServer, PropHost, ChangeTypeModified))
	assert.False(t, a.IsBreaking(CompInfo, PropTitle, ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompChannel, PropAddress, ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompKafkaChannelBinding, PropTopic, ChangeTypeModified))
	assert.False(t, a.IsBreaking(CompKafkaChannelBinding, PropBindingVersion, ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompSQSQueue, PropName, ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompMQTTOperationBinding, PropQoS, ChangeTypeModified))
	assert.True(t, a.IsBreaking(CompChannels, "", ChangeTypeRemoved))
	assert.False(t, a.IsBreaking(CompChannels, "", ChangeTypeAdded))
}

func TestBreakingRulesConfig_Merge(t *testing.T) {
	cfg := NewDefaultBreakingRulesConfig()
	require.False(t, cfg.IsBreaking(CompInfo, PropTitle, ChangeTypeModified))

	cfg.Merge(&BreakingRulesConfig{
		Info: &InfoRules{Title: rule(false, true, false)},
	})
	assert.True(t, cfg.IsBreaking(CompInfo, PropTitle, ChangeTypeModified))
	assert.False(t, cfg.IsBreaking(CompInfo, PropTitle, ChangeTypeAdded))

	// unrelated rules remain at defaults after the merge.
	assert.True(t, cfg.IsBreaking(CompServer, PropHost, ChangeTypeModified))

	// merging nil is a no-op.
	cfg.Merge(nil)
	assert.True(t, cfg.IsBreaking(CompInfo, PropTitle, ChangeTypeModified))

	// the shared default cache must not have been mutated.
	assert.False(t, GenerateDefaultBreakingRules().IsBreaking(CompInfo, PropTitle, ChangeTypeModified))
}

func TestActiveBreakingRulesConfig_SetAndReset(t *testing.T) {
	defer ResetActiveBreakingRulesConfig()

	require.False(t, BreakingModified(CompInfo, PropTitle))

	custom := NewDefaultBreakingRulesConfig()
	custom.Merge(&BreakingRulesConfig{Info: &InfoRules{Title: rule(false, true, false)}})
	SetActiveBreakingRulesConfig(custom)
	assert.True(t, BreakingModified(CompInfo, PropTitle))

	ResetActiveBreakingRulesConfig()
	assert.False(t, BreakingModified(CompInfo, PropTitle))
}

func TestBreakingRulesConfig_GetRule(t *testing.T) {
	cfg := GenerateDefaultBreakingRules()
	assert.NotNil(t, cfg.GetRule(CompAsyncAPI, ""))
	assert.NotNil(t, cfg.GetRule(CompServer, PropHost))
	assert.Nil(t, cfg.GetRule("not-a-component", "nope"))
	assert.Nil(t, cfg.GetRule(CompServer, "not-a-property"))
}

func TestValidateBreakingRulesConfigYAML_Valid(t *testing.T) {
	valid := []byte(`
asyncapi:
  modified: false
server:
  host:
    modified: false
kafkaChannelBinding:
  topic:
    modified: false
    removed: true
`)
	assert.Nil(t, ValidateBreakingRulesConfigYAML(valid))
}

func TestValidateBreakingRulesConfigYAML_MisplacedTopLevelKey(t *testing.T) {
	// "tag" is a top-level component, nesting it under "channel" is a misplacement.
	bad := []byte(`
channel:
  tag:
    name:
      modified: true
`)
	result := ValidateBreakingRulesConfigYAML(bad)
	require.NotNil(t, result)
	require.True(t, result.HasErrors())
	assert.Contains(t, result.Error(), "'tag'")
}

func TestValidateBreakingRulesConfigYAML_RuleFieldAtWrongDepth(t *testing.T) {
	// added/modified/removed must sit under a property, not directly under a
	// nested component.
	bad := []byte(`
channel:
  modified: true
`)
	result := ValidateBreakingRulesConfigYAML(bad)
	require.NotNil(t, result)
	require.True(t, result.HasErrors())
	assert.Contains(t, result.Errors[0].Message, "must be nested under a property name")
	assert.Positive(t, result.Errors[0].Line)
}

func TestValidateBreakingRulesConfigYAML_SimpleComponentRule(t *testing.T) {
	// simple components carry rules directly, so this is valid.
	ok := []byte(`
channels:
  removed: false
`)
	assert.Nil(t, ValidateBreakingRulesConfigYAML(ok))
}

func TestValidateBreakingRulesConfigYAML_InvalidYAML(t *testing.T) {
	result := ValidateBreakingRulesConfigYAML([]byte("a: [unclosed"))
	require.NotNil(t, result)
	assert.Contains(t, result.Error(), "invalid YAML")
}
