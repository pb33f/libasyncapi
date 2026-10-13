// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package schemaformat identifies schema languages whose references use JSON Schema semantics.
package schemaformat

import "mime"

const (
	asyncAPIJSONMediaType   = "application/vnd.aai.asyncapi+json"
	asyncAPIYAMLMediaType   = "application/vnd.aai.asyncapi+yaml"
	jsonSchemaJSONMediaType = "application/schema+json"
	jsonSchemaYAMLMediaType = "application/schema+yaml"
	openAPIMediaType        = "application/vnd.oai.openapi"
	openAPIJSONMediaType    = "application/vnd.oai.openapi+json"
	openAPIYAMLMediaType    = "application/vnd.oai.openapi+yaml"
)

// IsJSON reports whether format uses the supported JSON Schema reference syntax.
func IsJSON(format string) bool {
	mediaType, _, err := mime.ParseMediaType(format)
	if err != nil {
		return false
	}
	switch mediaType {
	case asyncAPIJSONMediaType, asyncAPIYAMLMediaType,
		jsonSchemaJSONMediaType, jsonSchemaYAMLMediaType,
		openAPIMediaType, openAPIJSONMediaType, openAPIYAMLMediaType:
		return true
	}
	return false
}
