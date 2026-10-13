// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package schemaformat identifies schema languages whose references use JSON Schema semantics.
package schemaformat

import "mime"

// IsJSON reports whether format uses the supported JSON Schema reference syntax.
func IsJSON(format string) bool {
	mediaType, _, err := mime.ParseMediaType(format)
	if err != nil {
		return false
	}
	switch mediaType {
	case "application/vnd.aai.asyncapi+json", "application/vnd.aai.asyncapi+yaml",
		"application/schema+json", "application/schema+yaml",
		"application/vnd.oai.openapi", "application/vnd.oai.openapi+json", "application/vnd.oai.openapi+yaml":
		return true
	}
	return false
}
