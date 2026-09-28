// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package model

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pb33f/go-yaml"
)

// BreakingRulesConfig holds all breaking change rules organized by AsyncAPI component.
// Structure mirrors the AsyncAPI 3.0 specification.
type BreakingRulesConfig struct {
	AsyncAPI           *BreakingChangeRule           `json:"asyncapi,omitempty" yaml:"asyncapi,omitempty"`
	ID                 *BreakingChangeRule           `json:"id,omitempty" yaml:"id,omitempty"`
	DefaultContentType *BreakingChangeRule           `json:"defaultContentType,omitempty" yaml:"defaultContentType,omitempty"`
	Servers            *BreakingChangeRule           `json:"servers,omitempty" yaml:"servers,omitempty"`
	Channels           *BreakingChangeRule           `json:"channels,omitempty" yaml:"channels,omitempty"`
	Operations         *BreakingChangeRule           `json:"operations,omitempty" yaml:"operations,omitempty"`
	Info               *InfoRules                    `json:"info,omitempty" yaml:"info,omitempty"`
	Contact            *ContactRules                 `json:"contact,omitempty" yaml:"contact,omitempty"`
	License            *LicenseRules                 `json:"license,omitempty" yaml:"license,omitempty"`
	Server             *ServerRules                  `json:"server,omitempty" yaml:"server,omitempty"`
	ServerVariable     *ServerVariableRules          `json:"serverVariable,omitempty" yaml:"serverVariable,omitempty"`
	Channel            *ChannelRules                 `json:"channel,omitempty" yaml:"channel,omitempty"`
	Parameter          *ParameterRules               `json:"parameter,omitempty" yaml:"parameter,omitempty"`
	Operation          *OperationRules               `json:"operation,omitempty" yaml:"operation,omitempty"`
	OperationTrait     *OperationTraitRules          `json:"operationTrait,omitempty" yaml:"operationTrait,omitempty"`
	OperationReply     *OperationReplyRules          `json:"operationReply,omitempty" yaml:"operationReply,omitempty"`
	OperationReplyAddr *OperationReplyAddressRules   `json:"operationReplyAddress,omitempty" yaml:"operationReplyAddress,omitempty"`
	Message            *MessageRules                 `json:"message,omitempty" yaml:"message,omitempty"`
	MessageTrait       *MessageTraitRules            `json:"messageTrait,omitempty" yaml:"messageTrait,omitempty"`
	MessageExample     *MessageExampleRules          `json:"messageExample,omitempty" yaml:"messageExample,omitempty"`
	CorrelationID      *CorrelationIDRules           `json:"correlationId,omitempty" yaml:"correlationId,omitempty"`
	Components         *ComponentsRules              `json:"components,omitempty" yaml:"components,omitempty"`
	SecurityScheme     *SecuritySchemeRules          `json:"securityScheme,omitempty" yaml:"securityScheme,omitempty"`
	OAuthFlows         *OAuthFlowsRules              `json:"oauthFlows,omitempty" yaml:"oauthFlows,omitempty"`
	OAuthFlow          *OAuthFlowRules               `json:"oauthFlow,omitempty" yaml:"oauthFlow,omitempty"`
	Tag                *TagRules                     `json:"tag,omitempty" yaml:"tag,omitempty"`
	ExternalDocs       *ExternalDocsRules            `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`
	ServerBindings     *ServerBindingsRules          `json:"serverBindings,omitempty" yaml:"serverBindings,omitempty"`
	ChannelBindings    *ChannelBindingsRules         `json:"channelBindings,omitempty" yaml:"channelBindings,omitempty"`
	OperationBindings  *OperationBindingsRules       `json:"operationBindings,omitempty" yaml:"operationBindings,omitempty"`
	MessageBindings    *MessageBindingsRules         `json:"messageBindings,omitempty" yaml:"messageBindings,omitempty"`
	HTTPOperation      *HTTPOperationBindingRules    `json:"httpOperationBinding,omitempty" yaml:"httpOperationBinding,omitempty"`
	HTTPMessage        *HTTPMessageBindingRules      `json:"httpMessageBinding,omitempty" yaml:"httpMessageBinding,omitempty"`
	KafkaServer        *KafkaServerBindingRules      `json:"kafkaServerBinding,omitempty" yaml:"kafkaServerBinding,omitempty"`
	KafkaChannel       *KafkaChannelBindingRules     `json:"kafkaChannelBinding,omitempty" yaml:"kafkaChannelBinding,omitempty"`
	KafkaTopicConfig   *KafkaTopicConfigurationRules `json:"kafkaTopicConfiguration,omitempty" yaml:"kafkaTopicConfiguration,omitempty"`
	KafkaOperation     *KafkaOperationBindingRules   `json:"kafkaOperationBinding,omitempty" yaml:"kafkaOperationBinding,omitempty"`
	KafkaMessage       *KafkaMessageBindingRules     `json:"kafkaMessageBinding,omitempty" yaml:"kafkaMessageBinding,omitempty"`
	WebSocketChannel   *WebSocketChannelBindingRules `json:"wsChannelBinding,omitempty" yaml:"wsChannelBinding,omitempty"`
	AMQPChannel        *AMQPChannelBindingRules      `json:"amqpChannelBinding,omitempty" yaml:"amqpChannelBinding,omitempty"`
	AMQPExchange       *AMQPExchangeRules            `json:"amqpExchange,omitempty" yaml:"amqpExchange,omitempty"`
	AMQPQueue          *AMQPQueueRules               `json:"amqpQueue,omitempty" yaml:"amqpQueue,omitempty"`
	AMQPOperation      *AMQPOperationBindingRules    `json:"amqpOperationBinding,omitempty" yaml:"amqpOperationBinding,omitempty"`
	AMQPMessage        *AMQPMessageBindingRules      `json:"amqpMessageBinding,omitempty" yaml:"amqpMessageBinding,omitempty"`
	MQTTServer         *MQTTServerBindingRules       `json:"mqttServerBinding,omitempty" yaml:"mqttServerBinding,omitempty"`
	MQTTLastWill       *MQTTLastWillRules            `json:"mqttLastWill,omitempty" yaml:"mqttLastWill,omitempty"`
	MQTTOperation      *MQTTOperationBindingRules    `json:"mqttOperationBinding,omitempty" yaml:"mqttOperationBinding,omitempty"`
	MQTTMessage        *MQTTMessageBindingRules      `json:"mqttMessageBinding,omitempty" yaml:"mqttMessageBinding,omitempty"`
	SQSChannel         *SQSChannelBindingRules       `json:"sqsChannelBinding,omitempty" yaml:"sqsChannelBinding,omitempty"`
	SQSOperation       *SQSOperationBindingRules     `json:"sqsOperationBinding,omitempty" yaml:"sqsOperationBinding,omitempty"`
	SQSQueue           *SQSQueueRules                `json:"sqsQueue,omitempty" yaml:"sqsQueue,omitempty"`
	SQSIdentifier      *SQSIdentifierRules           `json:"sqsIdentifier,omitempty" yaml:"sqsIdentifier,omitempty"`
	SQSRedrivePolicy   *SQSRedrivePolicyRules        `json:"sqsRedrivePolicy,omitempty" yaml:"sqsRedrivePolicy,omitempty"`
	SQSPolicy          *SQSPolicyRules               `json:"sqsPolicy,omitempty" yaml:"sqsPolicy,omitempty"`
	SQSPolicyStatement *SQSPolicyStatementRules      `json:"sqsPolicyStatement,omitempty" yaml:"sqsPolicyStatement,omitempty"`

	ruleCache map[string]*BreakingChangeRule
	cacheOnce sync.Once
}

var (
	configType = reflect.TypeOf(BreakingRulesConfig{})
	ruleType   = reflect.TypeOf((*BreakingChangeRule)(nil))
)

// Merge applies user overrides to the configuration. Only non-nil values from
// the override config replace the current values. Uses reflection to reduce boilerplate.
func (c *BreakingRulesConfig) Merge(override *BreakingRulesConfig) {
	if override == nil {
		return
	}

	cVal := reflect.ValueOf(c).Elem()
	oVal := reflect.ValueOf(override).Elem()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		cField := cVal.Field(i)
		oField := oVal.Field(i)

		if !cField.CanSet() {
			continue
		}

		if field.Type == ruleType {
			cField.Set(reflect.ValueOf(mergeRule(
				cField.Interface().(*BreakingChangeRule),
				oField.Interface().(*BreakingChangeRule),
			)))
			continue
		}

		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if oField.IsNil() {
				continue
			}
			if cField.IsNil() {
				cField.Set(reflect.New(field.Type.Elem()))
			}
			mergeRulesStruct(cField.Elem(), oField.Elem())
		}
	}

	c.invalidateCache()
}

// IsBreaking looks up whether a change is breaking based on the component, property, and change type.
// Returns the configured breaking status, or false if the rule is not found.
func (c *BreakingRulesConfig) IsBreaking(component, property, changeType string) bool {
	rule := c.GetRule(component, property)
	if rule == nil {
		return false
	}

	switch changeType {
	case ChangeTypeAdded:
		if rule.Added != nil {
			return *rule.Added
		}
	case ChangeTypeModified:
		if rule.Modified != nil {
			return *rule.Modified
		}
	case ChangeTypeRemoved:
		if rule.Removed != nil {
			return *rule.Removed
		}
	}
	return false
}

// GetRule returns the BreakingChangeRule for a given component and property.
// Returns nil if no rule is defined. Uses internal cache for O(1) lookups.
func (c *BreakingRulesConfig) GetRule(component, property string) *BreakingChangeRule {
	c.cacheOnce.Do(func() {
		c.ruleCache = c.buildRuleCache()
	})
	if property == "" {
		return c.ruleCache[component]
	}
	return c.ruleCache[component+"."+property]
}

// buildRuleCache creates a flat map of all rules for O(1) lookups using reflection.
func (c *BreakingRulesConfig) buildRuleCache() map[string]*BreakingChangeRule {
	cache := make(map[string]*BreakingChangeRule, 300)
	cVal := reflect.ValueOf(c).Elem()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fVal := cVal.Field(i)

		compName := jsonTagName(field)
		if compName == "" || !fVal.CanInterface() {
			continue
		}

		if field.Type == ruleType {
			cache[compName] = fVal.Interface().(*BreakingChangeRule)
			continue
		}

		if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if fVal.IsNil() {
				continue
			}
			addRulesToCache(cache, compName, fVal.Elem())
		}
	}
	return cache
}

// invalidateCache resets the cache so it will be rebuilt on next access.
func (c *BreakingRulesConfig) invalidateCache() {
	c.cacheOnce = sync.Once{}
	c.ruleCache = nil
}

// jsonTagName extracts the field name from a JSON struct tag.
func jsonTagName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx]
	}
	return tag
}

// mergeRulesStruct merges all *BreakingChangeRule fields from override into base.
func mergeRulesStruct(base, override reflect.Value) {
	rulesType := base.Type()
	for i := 0; i < rulesType.NumField(); i++ {
		field := rulesType.Field(i)
		if field.Type != ruleType {
			continue
		}
		bField := base.Field(i)
		oField := override.Field(i)
		if bField.CanSet() {
			bField.Set(reflect.ValueOf(mergeRule(
				bField.Interface().(*BreakingChangeRule),
				oField.Interface().(*BreakingChangeRule),
			)))
		}
	}
}

// addRulesToCache adds all *BreakingChangeRule fields from a rule struct to the cache.
func addRulesToCache(cache map[string]*BreakingChangeRule, compName string, rulesVal reflect.Value) {
	rulesType := rulesVal.Type()
	for i := 0; i < rulesType.NumField(); i++ {
		field := rulesType.Field(i)
		fVal := rulesVal.Field(i)

		propName := jsonTagName(field)
		if propName == "" || field.Type != ruleType {
			continue
		}
		cache[compName+"."+propName] = fVal.Interface().(*BreakingChangeRule)
	}
}

// mergeRule merges an override rule into a base rule.
// nil values in override are ignored, non-nil values replace the base.
func mergeRule(base, override *BreakingChangeRule) *BreakingChangeRule {
	if override == nil {
		return base
	}
	if base == nil {
		return override
	}
	result := &BreakingChangeRule{
		Added:    base.Added,
		Modified: base.Modified,
		Removed:  base.Removed,
	}
	if override.Added != nil {
		result.Added = override.Added
	}
	if override.Modified != nil {
		result.Modified = override.Modified
	}
	if override.Removed != nil {
		result.Removed = override.Removed
	}
	return result
}

// --- Config Validation ---

// ConfigValidationError represents a single validation issue in a breaking rules config.
type ConfigValidationError struct {
	// Message is a human-readable description of the issue.
	Message string

	// Path is the YAML path where the issue was found (e.g., "channel.address").
	Path string

	// Line is the 1-based line number in the YAML source (0 if unknown).
	Line int

	// Column is the 1-based column number in the YAML source (0 if unknown).
	Column int

	// FoundKey is the misplaced key that was detected.
	FoundKey string

	// SuggestedPath is where the key should be placed instead.
	SuggestedPath string
}

// Error implements the error interface.
func (e *ConfigValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("%s (line %d, column %d)", e.Message, e.Line, e.Column)
	}
	return e.Message
}

// ConfigValidationResult holds the results of validating a breaking rules config.
type ConfigValidationResult struct {
	// Errors contains all validation issues found.
	Errors []*ConfigValidationError
}

// HasErrors returns true if any validation errors were found.
func (r *ConfigValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

// Error implements the error interface, joining all errors.
func (r *ConfigValidationResult) Error() string {
	if !r.HasErrors() {
		return ""
	}
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

// validTopLevelComponents is the set of valid top-level keys in a breaking rules config.
// Built from BreakingRulesConfig struct field tags at init time.
var validTopLevelComponents = buildValidComponentSet()

// buildValidComponentSet creates a set of valid top-level component names
// by reflecting on the BreakingRulesConfig struct tags.
func buildValidComponentSet() map[string]bool {
	result := make(map[string]bool)
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		name := jsonTagName(field)
		if name != "" && name != "-" {
			result[name] = true
		}
	}
	return result
}

// breakingRuleFields are the valid fields for BreakingChangeRule
var breakingRuleFields = map[string]bool{
	"added":    true,
	"modified": true,
	"removed":  true,
}

// simpleRuleComponents are components that are directly BreakingChangeRule (not nested)
// For these, added/modified/removed at depth 1 is correct (e.g., "asyncapi.modified: false")
var simpleRuleComponents = map[string]bool{
	"asyncapi":           true,
	"id":                 true,
	"defaultContentType": true,
	"servers":            true,
	"channels":           true,
	"operations":         true,
}

// componentProperties maps each component to its valid property names
// Built from reflection on BreakingRulesConfig struct
var componentProperties = buildComponentPropertiesMap()

func buildComponentPropertiesMap() map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		componentName := jsonTagName(field)
		if componentName == "" || componentName == "-" {
			continue
		}
		// Get the field type (pointer to rules struct)
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}
		if fieldType.Kind() != reflect.Struct {
			continue
		}
		// Build property set for this component
		props := make(map[string]bool)
		for j := 0; j < fieldType.NumField(); j++ {
			propField := fieldType.Field(j)
			propName := jsonTagName(propField)
			if propName != "" && propName != "-" {
				props[propName] = true
			}
		}
		result[componentName] = props
	}
	return result
}

// isValidPropertyForComponent checks if a property is valid for the given component
func isValidPropertyForComponent(component, property string) bool {
	if props, ok := componentProperties[component]; ok {
		return props[property]
	}
	return false
}

// ValidateBreakingRulesConfigYAML validates raw YAML bytes for a breaking rules config.
// It detects misplaced nested configurations (e.g., "channel.tag" should be just "tag"
// at the top level) and returns all validation errors found.
// Returns nil if the configuration is valid.
func ValidateBreakingRulesConfigYAML(yamlBytes []byte) *ConfigValidationResult {
	var rootNode yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &rootNode); err != nil {
		return &ConfigValidationResult{
			Errors: []*ConfigValidationError{{
				Message: fmt.Sprintf("invalid YAML: %v", err),
			}},
		}
	}

	result := &ConfigValidationResult{}
	validateConfigNode(&rootNode, "", result)

	if result.HasErrors() {
		return result
	}
	return nil
}

// validateConfigNode recursively walks the YAML tree looking for misplaced configurations.
func validateConfigNode(node *yaml.Node, path string, result *ConfigValidationResult) {
	validateConfigNodeWithDepthAndProperty(node, path, 0, "", "", result)
}

// validateConfigNodeWithDepthAndProperty is the internal recursive validator.
// depth 0 = root, depth 1 = under a component (e.g., "channel"), depth 2 = under a property
// parentComponent tracks the root-level component we're under, parentProperty the property
// within that component.
func validateConfigNodeWithDepthAndProperty(node *yaml.Node, path string, depth int, parentComponent, parentProperty string, result *ConfigValidationResult) {
	// Document nodes contain a single content node
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		validateConfigNodeWithDepthAndProperty(node.Content[0], path, depth, parentComponent, parentProperty, result)
		return
	}

	// Only process mapping nodes (objects)
	if node.Kind != yaml.MappingNode {
		return
	}

	// Process key-value pairs in the mapping
	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		if keyNode.Kind != yaml.ScalarNode {
			continue
		}

		key := keyNode.Value
		currentPath := buildConfigPath(path, key)

		// Track the parent component when we enter a top-level component
		currentParent := parentComponent
		currentProperty := parentProperty
		if depth == 0 && validTopLevelComponents[key] {
			currentParent = key
			currentProperty = ""
		} else if depth == 1 && parentComponent != "" {
			// At depth 1, the key is a property name under a component
			currentProperty = key
		}

		// If we're already nested under a component and find another top-level component name,
		// this is a misplacement error UNLESS the key is a valid property of the parent component.
		// For example, "channel.servers" is valid because ChannelRules has a Servers property,
		// even though "servers" is also a top-level component.
		if path != "" && validTopLevelComponents[key] && !isValidPropertyForComponent(parentComponent, key) {
			result.Errors = append(result.Errors, &ConfigValidationError{
				Message:       fmt.Sprintf("found '%s' nested under '%s'; '%s' should be a top-level key", key, path, key),
				Path:          currentPath,
				Line:          keyNode.Line,
				Column:        keyNode.Column,
				FoundKey:      key,
				SuggestedPath: key,
			})
		}

		// Check for breaking rule fields (added/modified/removed) at wrong depth
		// depth 1 = directly under a component (e.g., "channel.added" is wrong)
		// These should only appear at depth 2 (e.g., "channel.address.added" is correct)
		// Exception: "simple" components like asyncapi, servers, channels are directly
		// BreakingChangeRule, so "asyncapi.modified: false" is correct
		if depth == 1 && breakingRuleFields[key] && !simpleRuleComponents[parentComponent] {
			componentName := path
			result.Errors = append(result.Errors, &ConfigValidationError{
				Message:       fmt.Sprintf("'%s' found directly under '%s'; breaking rules must be nested under a property name (e.g., '%s.<property>.%s' not '%s.%s')", key, componentName, componentName, key, componentName, key),
				Path:          currentPath,
				Line:          keyNode.Line,
				Column:        keyNode.Column,
				FoundKey:      key,
				SuggestedPath: fmt.Sprintf("%s.<property>.%s", componentName, key),
			})
		}

		// At depth 2+, check if we're trying to add invalid keys under a BreakingChangeRule.
		// A BreakingChangeRule (like channel.tags) can only have added/modified/removed.
		// If we find anything else (like "name"), it's someone trying to configure
		// a component's sub-rules under the wrong parent.
		if depth >= 2 && parentProperty != "" && !breakingRuleFields[key] && validTopLevelComponents[parentProperty] {
			if props, ok := componentProperties[parentProperty]; ok && props[key] {
				result.Errors = append(result.Errors, &ConfigValidationError{
					Message:       fmt.Sprintf("'%s' is incorrectly nested under '%s'; move your '%s:' block to the top level of your config (current: %s.%s.%s, should be: %s.%s)", parentProperty, parentComponent, parentProperty, parentComponent, parentProperty, key, parentProperty, key),
					Path:          currentPath,
					Line:          keyNode.Line,
					Column:        keyNode.Column,
					FoundKey:      parentProperty,
					SuggestedPath: parentProperty,
				})
			}
		}

		// Recurse into nested mappings
		if valueNode.Kind == yaml.MappingNode {
			validateConfigNodeWithDepthAndProperty(valueNode, currentPath, depth+1, currentParent, currentProperty, result)
		}
	}
}

// buildConfigPath constructs a dotted path from parent and child components.
func buildConfigPath(parent, child string) string {
	if parent == "" {
		return child
	}
	return parent + "." + child
}
