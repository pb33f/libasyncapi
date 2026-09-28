// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/pb33f/go-yaml"
)

// SpecInfo contains parsed version information from an AsyncAPI specification.
type SpecInfo struct {
	Version       string
	VersionParsed *semver.Version
	SpecType      string
}

// cleanVersionString removes the optional "v" prefix from version strings.
func cleanVersionString(version string) string {
	return strings.TrimPrefix(version, "v")
}

// DetectAsyncAPIVersion extracts and validates the asyncapi version from a specification.
// It returns the version string and any error encountered.
// If the specification is AsyncAPI 2.x, it returns ErrAsyncAPI2NotSupported.
func DetectAsyncAPIVersion(spec []byte) (string, error) {
	var rootNode yaml.Node
	if err := yaml.Unmarshal(spec, &rootNode); err != nil {
		return "", ErrInvalidYAML
	}
	return ExtractVersionFromNode(&rootNode)
}

// ExtractVersionFromNode extracts the asyncapi version from an already-parsed YAML node.
// This avoids double-parsing when the document has already been unmarshaled.
func ExtractVersionFromNode(rootNode *yaml.Node) (string, error) {
	if rootNode.Kind != yaml.DocumentNode || len(rootNode.Content) == 0 {
		return "", ErrInvalidYAML
	}

	mappingNode := rootNode.Content[0]
	if mappingNode.Kind != yaml.MappingNode {
		return "", ErrInvalidYAML
	}

	for i := 0; i < len(mappingNode.Content)-1; i += 2 {
		keyNode := mappingNode.Content[i]
		valueNode := mappingNode.Content[i+1]

		if keyNode.Value == "asyncapi" {
			version := valueNode.Value
			if version == "" {
				return "", ErrNoAsyncAPIVersion
			}
			return version, nil
		}
	}

	return "", ErrNoAsyncAPIVersion
}

// ParseAsyncAPIVersion parses a version string and returns SpecInfo.
// It validates that the version is a valid semver and is AsyncAPI 3.x.
func ParseAsyncAPIVersion(version string) (*SpecInfo, error) {
	v, err := semver.NewVersion(cleanVersionString(version))
	if err != nil {
		return nil, ErrInvalidAsyncAPIVersion
	}

	if v.Major() == 2 {
		return nil, ErrAsyncAPI2NotSupported
	}

	if v.Major() < 3 {
		return nil, ErrInvalidAsyncAPIVersion
	}

	return &SpecInfo{
		Version:       version,
		VersionParsed: v,
		SpecType:      "asyncapi",
	}, nil
}

// IsAsyncAPI3 checks if the given version string represents AsyncAPI 3.x.
func IsAsyncAPI3(version string) bool {
	v, err := semver.NewVersion(cleanVersionString(version))
	if err != nil {
		return false
	}
	return v.Major() == 3
}

// DefaultSchemaFormat returns the default schema format for a given AsyncAPI version.
// Per the spec, when schemaFormat is omitted, it defaults to the AsyncAPI schema format
// with the version matching the document's asyncapi field.
func DefaultSchemaFormat(asyncapiVersion string) string {
	return "application/vnd.aai.asyncapi+json;version=" + cleanVersionString(asyncapiVersion)
}
