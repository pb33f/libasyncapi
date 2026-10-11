// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBundleScalarKeys(t *testing.T) {
	spec := []byte(header + `channels:
  404:
    bindings:
      kafka:
        1: retained
    messages:
      1:
        payload:
          type: object
          properties:
            1: {$ref: scalar.yaml}
            true: {$ref: scalar.yaml}
operations:
  send:
    action: send
    channel: {$ref: '#/channels/404'}
    messages: [{$ref: '#/channels/404/messages/1'}]
`)
	for _, composed := range []bool{false, true} {
		out, err := runBundle(spec, filesConfig(map[string]string{"scalar.yaml": "type: string"}), composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "retained", value(t, n, "/channels/404/bindings/kafka/1").Value)
		e := &engine{root: n}
		for _, key := range []string{"1", "true"} {
			schema, err := e.outputTarget("/channels/404/messages/1/payload/properties/" + key)
			require.NoError(t, err)
			assert.Equal(t, "string", field(schema, "type").Value)
		}
		again, err := runBundle(out, nil, composed)
		require.NoError(t, err)
		assert.Equal(t, string(out), string(again))
		for _, keys := range []string{"404: {}\n  '404': {}", "? [a, b]\n  : {}"} {
			out, err := runBundle([]byte(header+"channels:\n  "+keys), nil, composed)
			require.ErrorContains(t, err, "invalid or duplicate key")
			assert.Nil(t, out)
		}
	}
}

func TestBundleComposedUsesDeclaredComponents(t *testing.T) {
	components := `components:
  schemas:
    User: {$ref: './schemas/User.yaml'}
    Duplicate: {$ref: './schemas/User.yaml'}
    Local: {$ref: '#/components/schemas/User'}
  messages:
    UserCreated: {$ref: './messages/UserCreated.yaml'}
`
	channels := `channels:
  users:
    messages:
      created: {$ref: './messages/UserCreated.yaml'}
operations:
  receive:
    action: receive
    channel: {$ref: '#/channels/users'}
    messages: [{$ref: './messages/UserCreated.yaml'}]
`
	cfg := filesConfig(map[string]string{
		"schemas/User.yaml":         "type: object\nproperties:\n  previous: {$ref: '#'}\n",
		"messages/UserCreated.yaml": "name: UserCreated\npayload: {$ref: '../schemas/User.yaml'}\n",
	})
	for _, body := range []string{channels + components, components + channels} {
		out, err := BundleBytesComposed([]byte(header+body), cfg, nil)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Len(t, value(t, n, "/components/schemas").Content, 6)
		assert.Len(t, value(t, n, "/components/messages").Content, 2)
		assert.Equal(t, "object", value(t, n, "/components/schemas/User/type").Value)
		for _, p := range []string{"/components/schemas/User/properties/previous", "/components/schemas/Duplicate", "/components/schemas/Local", "/components/messages/UserCreated/payload"} {
			assert.Equal(t, "#/components/schemas/User", value(t, n, p+"/$ref").Value)
		}
		assert.Equal(t, "UserCreated", value(t, n, "/components/messages/UserCreated/name").Value)
		assert.Equal(t, "#/components/messages/UserCreated", value(t, n, "/channels/users/messages/created/$ref").Value)
		assert.Equal(t, "#/channels/users/messages/created", value(t, n, "/operations/receive/messages/0/$ref").Value)
		again, err := BundleBytesComposed(out, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, string(out), string(again))
	}
}

func TestBundleRejectsReferencedTextAndKeepsBooleanSchemas(t *testing.T) {
	for _, composed := range []bool{false, true} {
		for _, text := range []string{"root:x:0:0:root:/root:/bin/sh\n", "42", "null", "'true'", "true", "false"} {
			out, err := runBundle([]byte(header+"components: {schemas: {A: {$ref: data.yaml}}}"), filesConfig(map[string]string{"data.yaml": text}), composed)
			if text == "true" || text == "false" {
				require.NoError(t, err)
				assert.Equal(t, text, value(t, decode(t, out), "/components/schemas/A").Value)
			} else {
				require.ErrorContains(t, err, "must contain an object, array, or boolean schema")
				assert.Nil(t, out)
			}
		}
		out, err := runBundle([]byte(header+"components: {schemas: {A: {$ref: '../outside.yaml'}}}"), filesConfig(map[string]string{"outside.yaml": "type: string"}), composed)
		require.Error(t, err)
		assert.Nil(t, out)
	}
}

func TestBundleOpenAPIFormatSchemas(t *testing.T) {
	for _, format := range []string{"application/vnd.oai.openapi;version=3.0.0", "application/vnd.oai.openapi+json;version=3.0.3", "application/vnd.oai.openapi+yaml;version=3.1.0"} {
		for _, composed := range []bool{false, true} {
			spec := []byte(header + fmt.Sprintf("channels:\n  pets:\n    messages:\n      pet:\n        payload:\n          schemaFormat: %s\n          schema: {$ref: pet.yaml}\n", format))
			cfg := filesConfig(map[string]string{"pet.yaml": "type: object\nproperties:\n  name: {$ref: name.yaml}\n", "name.yaml": "type: string"})
			out, err := runBundle(spec, cfg, composed)
			require.NoError(t, err)
			e := &engine{root: decode(t, out)}
			n, err := e.outputTarget("/channels/pets/messages/pet/payload/schema/properties/name")
			require.NoError(t, err)
			assert.Equal(t, "string", field(n, "type").Value)
			assert.NotContains(t, string(out), ".yaml")
			again, err := runBundle(out, nil, composed)
			require.NoError(t, err)
			assert.Equal(t, string(out), string(again))
		}
	}
}

func TestAliasBudgetCountsOnlyExpansion(t *testing.T) {
	var n yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("values: [1, 2, 3, 4, 5]\n"), &n))
	budget := 1
	_, err := expandAliases(&n, make(map[*yaml.Node]bool), 0, &budget, false)
	require.NoError(t, err)
	assert.Equal(t, 1, budget)

	require.NoError(t, yaml.Unmarshal([]byte("base: &base [1, 2, 3]\ncopy: *base\n"), &n))
	budget = 4 // The copied sequence plus its three scalar children.
	out, err := expandAliases(&n, make(map[*yaml.Node]bool), 0, &budget, false)
	require.NoError(t, err)
	assert.Zero(t, budget)
	assert.Len(t, field(out.Content[0], "copy").Content, 3)
	budget = 3
	_, err = expandAliases(&n, make(map[*yaml.Node]bool), 0, &budget, false)
	require.ErrorContains(t, err, "document expansion limit exceeded")

	require.NoError(t, yaml.Unmarshal([]byte("base: &base [1, 2]\nnext: &next [*base, *base]\ncopy: *next\n"), &n))
	budget = 12
	_, err = expandAliases(&n, make(map[*yaml.Node]bool), 0, &budget, false)
	require.ErrorContains(t, err, "document expansion limit exceeded")
}

func TestBundleIndentationIgnoresScalarContent(t *testing.T) {
	for _, prefix := range []string{
		"x-note: {text: 'wrapped\n continuation'}\n",
		"x-note: |\n scalar content\n",
		"# a comment\n # indented comment\n",
	} {
		spec := []byte(header + prefix + "channels:\n    events:\n        address: events\n")
		out, err := BundleBytes(spec, nil)
		require.NoError(t, err)
		assert.Contains(t, string(out), "channels:\n    events:\n        address: events\n")
	}
	var n yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("{\n \"asyncapi\": \"3.0.0\",\n \"info\": {\"title\": \"JSON\", \"version\": \"1\"}\n}"), &n))
	assert.Equal(t, 2, indentation(&n))
	assert.False(t, jsonSchemaFormat("application/not-asyncapi"))
	assert.False(t, jsonSchemaFormat("invalid; format"))
	assert.True(t, jsonSchemaFormat("application/vnd.aai.asyncapi+json;version=3.0.0"))
	assert.True(t, jsonSchemaFormat("application/schema+json;version=draft-07"))
	assert.False(t, jsonSchemaFormat(strings.Repeat(" ", 3)))
}
