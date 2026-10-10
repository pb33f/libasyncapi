// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

const header = "asyncapi: 3.1.0\ninfo: {title: Test, version: 1.0.0}\n"

func decode(t *testing.T, b []byte) *yaml.Node {
	t.Helper()
	var n yaml.Node
	require.NoError(t, yaml.Unmarshal(b, &n))
	return n.Content[0]
}

func value(t *testing.T, n *yaml.Node, p string) *yaml.Node {
	t.Helper()
	v, err := pointer(n, p)
	require.NoError(t, err)
	return v
}

func runBundle(spec []byte, cfg *libasyncapi.DocumentConfiguration, composed bool) ([]byte, error) {
	if composed {
		return BundleBytesComposed(spec, cfg, &BundleCompositionConfig{Delimiter: "__"})
	}
	return BundleBytes(spec, cfg)
}

func filesConfig(files map[string]string) *libasyncapi.DocumentConfiguration {
	fs := fstest.MapFS{}
	for name, data := range files {
		fs[name] = &fstest.MapFile{Data: []byte(data)}
	}
	return &libasyncapi.DocumentConfiguration{BasePath: "/specs", SpecFilePath: "root.yaml", AllowFileReferences: true, LocalFS: fs, ExcludeExtensionRefs: true}
}

func TestBundlePortableFixture(t *testing.T) {
	spec, err := os.ReadFile("../test_fixtures/bundler/asyncapi.yaml")
	require.NoError(t, err)
	base, err := filepath.Abs("../test_fixtures/bundler")
	require.NoError(t, err)
	for _, composed := range []bool{false, true} {
		t.Run(fmt.Sprint(composed), func(t *testing.T) {
			cfg := &libasyncapi.DocumentConfiguration{BasePath: base, SpecFilePath: "asyncapi.yaml", ExcludeExtensionRefs: true}
			before := append([]byte(nil), spec...)
			result, err := runBundle(spec, cfg, composed)
			require.NoError(t, err)
			assert.Equal(t, before, spec)
			assert.Equal(t, "asyncapi.yaml", cfg.SpecFilePath)
			n := decode(t, result)
			assert.Equal(t, "#/channels/orders", value(t, n, "/operations/receiveOrder/channel/$ref").Value)
			assert.Equal(t, "#/channels/orders/messages/created", value(t, n, "/operations/receiveOrder/messages/0/$ref").Value)
			assert.Empty(t, value(t, n, "/operations/sendOrder/messages").Content)
			assert.Equal(t, "literal-example.yaml", value(t, n, "/x-data/$ref").Value)
			e := &engine{root: n}
			channel, err := e.outputTarget("/channels/orders")
			require.NoError(t, err)
			assert.Equal(t, "#/servers/production", value(t, channel, "/servers/0/$ref").Value)
			message, err := e.outputTarget("/channels/orders/messages/created")
			require.NoError(t, err)
			assert.Equal(t, "literal-payload.yaml", value(t, message, "/examples/0/payload/$ref").Value)
			if composed {
				assert.Equal(t, "#/components/schemas/Order__schemas", value(t, message, "/payload/$ref").Value)
				assert.Equal(t, "#/components/schemas/Order__schemas", value(t, n, "/components/schemas/Order__schemas/properties/previous/$ref").Value)
			} else {
				assert.Equal(t, "#/channels/orders/messages/created/payload", value(t, message, "/payload/properties/previous/$ref").Value)
			}
			// The second run has no filesystem and must preserve a portable bundle.
			again, err := runBundle(result, &libasyncapi.DocumentConfiguration{ExcludeExtensionRefs: true}, composed)
			require.NoError(t, err)
			assert.Equal(t, string(result), string(again))
		})
	}
}

func TestBundleRequiredReferencesThroughRootAliases(t *testing.T) {
	for _, composed := range []bool{false, true} {
		cfg := filesConfig(map[string]string{"channel.yaml": "address: a\nmessages:\n  event: {name: Event}\n", "operation.yaml": "action: send\nchannel: {$ref: './channel.yaml'}\nmessages:\n - {$ref: './channel.yaml#/messages/event'}\n"})
		result, err := runBundle([]byte(header+"channels:\n  a: {$ref: channel.yaml}\noperations:\n  send: {$ref: operation.yaml}\n"), cfg, composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, result)}
		op, err := e.outputTarget("/operations/send")
		require.NoError(t, err)
		assert.Equal(t, "#/channels/a", value(t, op, "/channel/$ref").Value)
		assert.Equal(t, "#/channels/a/messages/event", value(t, op, "/messages/0/$ref").Value)
	}
}

func TestBundleEveryComponentCategory(t *testing.T) {
	for _, kind := range []string{"schemas", "servers", "channels", "operations", "messages", "securitySchemes", "serverVariables", "parameters", "correlationIds", "replies", "replyAddresses", "externalDocs", "tags", "operationTraits", "messageTraits", "serverBindings", "channelBindings", "operationBindings", "messageBindings"} {
		t.Run(kind, func(t *testing.T) {
			cfg := filesConfig(map[string]string{"object.yaml": "description: retained\nx-custom: {value: true}\n"})
			spec := []byte(header + "components:\n  " + kind + ":\n    imported: {$ref: object.yaml}\n")
			result, err := BundleBytesComposed(spec, cfg, nil)
			require.NoError(t, err)
			n := decode(t, result)
			assert.Equal(t, "#/components/"+kind+"/object", value(t, n, "/components/"+kind+"/imported/$ref").Value)
			assert.Equal(t, "retained", value(t, n, "/components/"+kind+"/object/description").Value)
		})
	}
}

func TestBundleErrors(t *testing.T) {
	tests := []struct {
		name, spec string
		files      map[string]string
		want       string
	}{
		{"yaml", "[", nil, "parse AsyncAPI"},
		{"version", "asyncapi: 2.6.0", nil, "2"},
		{"missing", header + "channels: {a: {$ref: missing.yaml}}", nil, "load"},
		{"pointer", header + "channels: {a: {$ref: '#/absent'}}", nil, "not found"},
		{"badRef", header + "channels: {a: {$ref: []}}", nil, "invalid $ref"},
		{"loop", header + "channels: {a: {$ref: a.yaml}}", map[string]string{"a.yaml": "$ref: b.yaml", "b.yaml": "$ref: a.yaml"}, "reference-only cycle"},
		{"channelObject", header + "operations: {a: {channel: {address: a}}}", nil, "requires a Reference"},
		{"unplacedChannel", header + "operations: {a: {channel: {$ref: a.yaml}}}", map[string]string{"a.yaml": "address: a"}, "expected one local"},
		{"ambiguous", header + "channels: {a: {$ref: a.yaml}, b: {$ref: a.yaml}}\noperations: {a: {channel: {$ref: a.yaml}}}", map[string]string{"a.yaml": "address: a"}, "found 2"},
		{"duplicate", header + "channels: {a: {}, a: {}}", nil, "duplicate"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, composed := range []bool{false, true} {
				out, err := runBundle([]byte(tt.spec), filesConfig(tt.files), composed)
				require.Error(t, err)
				assert.Nil(t, out)
				assert.Contains(t, err.Error(), tt.want)
			}
		})
	}
}

func TestBundleRemoteAndDisabledAccess(t *testing.T) {
	spec := []byte(header + "channels: {a: {$ref: 'https://example.com/channel.yaml'}}")
	calls := 0
	cfg := &libasyncapi.DocumentConfiguration{AllowRemoteReferences: true, RemoteURLHandler: func(location string) (*http.Response, error) {
		calls++
		assert.Equal(t, "https://example.com/channel.yaml", location)
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("address: a")), Header: http.Header{"Content-Type": []string{"application/yaml"}}}, nil
	}}
	out, err := BundleBytes(spec, cfg)
	require.NoError(t, err)
	assert.Contains(t, string(out), "address: a")
	assert.Equal(t, 1, calls)
	cfg.AllowRemoteReferences = false
	out, err = BundleBytes(spec, cfg)
	require.Error(t, err)
	assert.Nil(t, out)
	assert.Equal(t, 1, calls)
}

func TestBundleEscapedPointersAndOpaqueSchemas(t *testing.T) {
	cfg := filesConfig(map[string]string{"schemas.yaml": "a/b~c+%:\n  schemaFormat: application/vnd.apache.avro+json;version=1.9.0\n  schema:\n    type: record\n    fields: [{name: x, type: string}]\n"})
	spec := []byte(header + "components:\n  schemas:\n    event: {$ref: 'schemas.yaml#/a~1b~0c+%25'}\n")
	out, err := BundleBytes(spec, cfg)
	require.NoError(t, err)
	assert.Equal(t, "record", value(t, decode(t, out), "/components/schemas/event/schema/type").Value)
	assert.Contains(t, string(out), "schemaFormat:")
}

func TestBundleDelimiterAndJSON(t *testing.T) {
	for _, delimiter := range []string{"/", "#", " ", "\u00a0"} {
		_, err := BundleBytesComposed([]byte(header), nil, &BundleCompositionConfig{Delimiter: delimiter})
		require.Error(t, err)
	}
	input := []byte(`{"asyncapi":"3.0.0","info":{"title":"JSON","version":"1"}}`)
	out, err := BundleBytes(input, nil)
	require.NoError(t, err)
	require.True(t, bytes.Contains(out, []byte("JSON")))
}

func TestBundleSchemaResourceScopes(t *testing.T) {
	for _, composed := range []bool{false, true} {
		cfg := filesConfig(map[string]string{
			"schema.yaml": "$id: 'https://schemas.example.com/order'\ntype: object\ndefinitions:\n  value: {$id: '#value', type: string}\nproperties:\n  id: {$ref: '#value'}\n  parent: {$ref: '#'}\n",
		})
		out, err := runBundle([]byte(header+"components:\n  schemas:\n    Order: {$ref: schema.yaml}\n"), cfg, composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, out)}
		schema, err := e.outputTarget("/components/schemas/Order")
		require.NoError(t, err)
		assert.Nil(t, field(schema, "$id"))
		id, err := e.outputTarget("/components/schemas/Order/properties/id")
		require.NoError(t, err)
		assert.Equal(t, "string", field(id, "type").Value)
		_, err = runBundle(out, nil, composed)
		require.NoError(t, err)
	}
}

func BenchmarkBundle(b *testing.B) {
	files := map[string]string{"message.yaml": "name: Event\npayload: {type: object, properties: {id: {type: string}}}"}
	var spec strings.Builder
	spec.WriteString(header + "channels:\n")
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&spec, "  channel%d:\n    messages: {event: {$ref: message.yaml}}\n", i)
	}
	input := []byte(spec.String())
	for _, composed := range []bool{false, true} {
		b.Run(fmt.Sprint(composed), func(b *testing.B) {
			cfg := filesConfig(files)
			b.ReportAllocs()
			for b.Loop() {
				if _, err := runBundle(input, cfg, composed); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
