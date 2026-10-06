// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"testing"

	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestDetectAsyncAPIVersion_Valid300(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"`)
	version, err := DetectAsyncAPIVersion(spec)
	require.NoError(t, err)
	assert.Equal(t, "3.0.0", version)
}

func TestDetectAsyncAPIVersion_Valid301(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.1"`)
	version, err := DetectAsyncAPIVersion(spec)
	require.NoError(t, err)
	assert.Equal(t, "3.0.1", version)
}

func TestDetectAsyncAPIVersion_UnquotedVersion(t *testing.T) {
	spec := []byte(`asyncapi: 3.0.0`)
	version, err := DetectAsyncAPIVersion(spec)
	require.NoError(t, err)
	assert.Equal(t, "3.0.0", version)
}

func TestDetectAsyncAPIVersion_WithOtherFields(t *testing.T) {
	spec := []byte(`asyncapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"`)
	version, err := DetectAsyncAPIVersion(spec)
	require.NoError(t, err)
	assert.Equal(t, "3.0.0", version)
}

func TestDetectAsyncAPIVersion_MissingAsyncAPIField(t *testing.T) {
	spec := []byte(`openapi: "3.0.0"`)
	_, err := DetectAsyncAPIVersion(spec)
	assert.ErrorIs(t, err, ErrNoAsyncAPIVersion)
}

func TestDetectAsyncAPIVersion_EmptySpec(t *testing.T) {
	spec := []byte(``)
	_, err := DetectAsyncAPIVersion(spec)
	assert.ErrorIs(t, err, ErrInvalidYAML)
}

func TestDetectAsyncAPIVersion_InvalidYAML(t *testing.T) {
	spec := []byte(`{{{invalid`)
	_, err := DetectAsyncAPIVersion(spec)
	assert.ErrorIs(t, err, ErrInvalidYAML)
}

func TestDetectAsyncAPIVersion_EmptyValue(t *testing.T) {
	spec := []byte(`asyncapi: ""`)
	_, err := DetectAsyncAPIVersion(spec)
	assert.ErrorIs(t, err, ErrNoAsyncAPIVersion)
}

func TestParseAsyncAPIVersion_Valid300(t *testing.T) {
	specInfo, err := ParseAsyncAPIVersion("3.0.0")
	require.NoError(t, err)
	require.NotNil(t, specInfo)
	assert.Equal(t, "3.0.0", specInfo.Version)
	assert.Equal(t, "asyncapi", specInfo.SpecType)
	assert.Equal(t, uint64(3), specInfo.VersionParsed.Major())
	assert.Equal(t, uint64(0), specInfo.VersionParsed.Minor())
	assert.Equal(t, uint64(0), specInfo.VersionParsed.Patch())
}

func TestParseAsyncAPIVersion_Valid301(t *testing.T) {
	specInfo, err := ParseAsyncAPIVersion("3.0.1")
	require.NoError(t, err)
	require.NotNil(t, specInfo)
	assert.Equal(t, uint64(3), specInfo.VersionParsed.Major())
	assert.Equal(t, uint64(0), specInfo.VersionParsed.Minor())
	assert.Equal(t, uint64(1), specInfo.VersionParsed.Patch())
}

func TestParseAsyncAPIVersion_Valid310(t *testing.T) {
	specInfo, err := ParseAsyncAPIVersion("3.1.0")
	require.NoError(t, err)
	require.NotNil(t, specInfo)
	assert.Equal(t, uint64(3), specInfo.VersionParsed.Major())
	assert.Equal(t, uint64(1), specInfo.VersionParsed.Minor())
}

func TestParseAsyncAPIVersion_WithVPrefix(t *testing.T) {
	specInfo, err := ParseAsyncAPIVersion("v3.0.0")
	require.NoError(t, err)
	require.NotNil(t, specInfo)
	assert.Equal(t, "v3.0.0", specInfo.Version)
	assert.Equal(t, uint64(3), specInfo.VersionParsed.Major())
}

func TestParseAsyncAPIVersion_AsyncAPI2NotSupported(t *testing.T) {
	_, err := ParseAsyncAPIVersion("2.6.0")
	assert.ErrorIs(t, err, ErrAsyncAPI2NotSupported)
}

func TestParseAsyncAPIVersion_AsyncAPI200NotSupported(t *testing.T) {
	_, err := ParseAsyncAPIVersion("2.0.0")
	assert.ErrorIs(t, err, ErrAsyncAPI2NotSupported)
}

func TestParseAsyncAPIVersion_InvalidVersionString(t *testing.T) {
	_, err := ParseAsyncAPIVersion("invalid")
	assert.ErrorIs(t, err, ErrInvalidAsyncAPIVersion)
}

func TestParseAsyncAPIVersion_Version1NotSupported(t *testing.T) {
	_, err := ParseAsyncAPIVersion("1.0.0")
	assert.ErrorIs(t, err, ErrInvalidAsyncAPIVersion)
}

func TestIsAsyncAPI3_True(t *testing.T) {
	assert.True(t, IsAsyncAPI3("3.0.0"))
	assert.True(t, IsAsyncAPI3("3.0.1"))
	assert.True(t, IsAsyncAPI3("3.1.0"))
	assert.True(t, IsAsyncAPI3("v3.0.0"))
}

func TestIsAsyncAPI3_False(t *testing.T) {
	assert.False(t, IsAsyncAPI3("2.6.0"))
	assert.False(t, IsAsyncAPI3("2.0.0"))
	assert.False(t, IsAsyncAPI3("1.0.0"))
	assert.False(t, IsAsyncAPI3("invalid"))
	assert.False(t, IsAsyncAPI3(""))
}

func TestDefaultSchemaFormat(t *testing.T) {
	assert.Equal(t, "application/vnd.aai.asyncapi+json;version=3.0.0", DefaultSchemaFormat("3.0.0"))
	assert.Equal(t, "application/vnd.aai.asyncapi+json;version=3.0.1", DefaultSchemaFormat("3.0.1"))
	assert.Equal(t, "application/vnd.aai.asyncapi+json;version=3.1.0", DefaultSchemaFormat("3.1.0"))
	// Verify v-prefix is stripped
	assert.Equal(t, "application/vnd.aai.asyncapi+json;version=3.0.0", DefaultSchemaFormat("v3.0.0"))
}
