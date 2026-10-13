// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBundleExternalMessagesThroughComponentChannel(t *testing.T) {
	for _, composed := range []bool{false, true} {
		for _, channelRef := range []string{"#/components/channels/shared", "#/components/channels/alias"} {
			t.Run(fmt.Sprintf("composed=%t/%s", composed, channelRef), func(t *testing.T) {
				cfg := filesConfig(map[string]string{
					"channels/events.yaml": "address: events\nmessages: {event: {$ref: '../messages/event.yaml'}}\n",
					"messages/event.yaml":  "name: Event\npayload: {type: string}\n",
				})
				spec := header + fmt.Sprintf(`channels:
  events: {$ref: '%s'}
operations:
  receive:
    action: receive
    channel: {$ref: '#/channels/events'}
    messages: [{$ref: '#/channels/events/messages/event'}]
    reply:
      channel: {$ref: '#/channels/events'}
      messages: [{$ref: '#/channels/events/messages/event'}]
components:
  channels:
    alias: {$ref: '#/components/channels/shared'}
    shared: {$ref: channels/events.yaml}
`, channelRef)
				out, err := runBundle([]byte(spec), cfg, composed)
				require.NoError(t, err)
				n := decode(t, out)
				for _, path := range []string{"/operations/receive/messages/0/$ref", "/operations/receive/reply/messages/0/$ref"} {
					assert.Equal(t, "#/channels/events/messages/event", value(t, n, path).Value)
				}
				again, err := runBundle(out, nil, composed)
				require.NoError(t, err)
				assert.Equal(t, string(out), string(again))
			})
		}
	}
}

func TestBundleSharedSchemaGraphStaysSmall(t *testing.T) {
	files := make(map[string]string)
	for i := 0; i < 26; i++ {
		schema := "type: object\nproperties:\n"
		for j := 1; j <= 3 && i+j < 26; j++ {
			schema += fmt.Sprintf("  p%d: {$ref: s%d.yaml}\n", j, i+j)
		}
		if i == 25 {
			schema = "type: string\n"
		}
		files[fmt.Sprintf("s%d.yaml", i)] = schema
	}
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(header+"components: {schemas: {Root: {$ref: s0.yaml}}}\n"), filesConfig(files), composed)
		require.NoError(t, err)
		assert.Less(t, len(out), 50_000)
		again, err := runBundle(out, nil, composed)
		require.NoError(t, err)
		assert.Equal(t, string(out), string(again))
	}
}

func TestBundleJSONSurrogateEscapes(t *testing.T) {
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(`{"asyncapi":"3.0.0","info":{"title":"\ud83c\udf03","version":"1"},"components":{"schemas":{"S":{"$ref":"s.json"}}}}`), filesConfig(map[string]string{"s.json": `{"type":"string","description":"\ud83c\udf03"}`}), composed)
		require.NoError(t, err)
		assert.Equal(t, "🌃", value(t, decode(t, out), "/info/title").Value)
		_, err = runBundle(out, nil, composed)
		require.NoError(t, err)
	}
}

func TestBundleProtocolBindingSchemas(t *testing.T) {
	tests := []struct {
		kind, protocol, field string
	}{
		{"operationBindings", "kafka", "groupId"},
		{"operationBindings", "kafka", "clientId"},
		{"serverBindings", "mqtt", "sessionExpiryInterval"},
		{"serverBindings", "mqtt", "maximumPacketSize"},
		{"operationBindings", "mqtt", "messageExpiryInterval"},
		{"messageBindings", "mqtt", "correlationData"},
		{"messageBindings", "mqtt", "responseTopic"},
	}
	for _, tt := range tests {
		for _, composed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%s/composed=%t", tt.protocol, tt.field, composed), func(t *testing.T) {
				cfg := filesConfig(map[string]string{
					"bindings/shared.yaml": fmt.Sprintf("%s:\n  %s: {$ref: '../schemas/value.yaml'}\n", tt.protocol, tt.field),
					"schemas/value.yaml":   "type: string\nenum: [events]\n",
				})
				spec := []byte(header + fmt.Sprintf("components:\n  %s:\n    shared: {$ref: bindings/shared.yaml}\n", tt.kind))
				out, err := runBundle(spec, cfg, composed)
				require.NoError(t, err)
				e := &engine{root: decode(t, out)}
				schema, err := e.outputTarget("/components/" + tt.kind + "/shared/" + tt.protocol + "/" + tt.field)
				require.NoError(t, err)
				assert.Equal(t, "string", value(t, schema, "/type").Value)
				assert.Equal(t, "events", value(t, schema, "/enum/0").Value)
				again, err := runBundle(out, nil, composed)
				require.NoError(t, err)
				assert.Equal(t, string(out), string(again))
				missing := filesConfig(map[string]string{
					"bindings/shared.yaml": fmt.Sprintf("%s:\n  %s: {$ref: '../schemas/missing.yaml'}\n", tt.protocol, tt.field),
				})
				out, err = runBundle(spec, missing, composed)
				require.ErrorContains(t, err, "missing.yaml")
				assert.Nil(t, out)
			})
		}
	}
}

func TestBundleBindingDataRemainsOpaque(t *testing.T) {
	spec := []byte(header + `components:
  serverBindings:
    mqtt:
      mqtt:
        clientId: literal-client
        sessionExpiryInterval: 60
        maximumPacketSize: 1200
  operationBindings:
    custom:
      custom:
        groupId: {$ref: literal.yaml}
        correlationData: {$ref: literal.yaml}
        responseTopic: {$ref: literal.yaml}
  messageBindings:
    mqtt:
      mqtt:
        responseTopic: events/reply
        x-example: {$ref: literal.yaml}
`)
	for _, composed := range []bool{false, true} {
		out, err := runBundle(spec, &libasyncapi.DocumentConfiguration{ExcludeExtensionRefs: true}, composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "literal-client", value(t, n, "/components/serverBindings/mqtt/mqtt/clientId").Value)
		assert.Equal(t, "60", value(t, n, "/components/serverBindings/mqtt/mqtt/sessionExpiryInterval").Value)
		assert.Equal(t, "events/reply", value(t, n, "/components/messageBindings/mqtt/mqtt/responseTopic").Value)
		assert.Equal(t, "literal.yaml", value(t, n, "/components/operationBindings/custom/custom/groupId/$ref").Value)
		assert.Equal(t, "literal.yaml", value(t, n, "/components/messageBindings/mqtt/mqtt/x-example/$ref").Value)
	}
}

func TestBundleExternalMultiFormatSchemaContent(t *testing.T) {
	for _, content := range []string{
		`{"type":"record","name":"Event","fields":[{"name":"value","type":"string","default":{"$ref":"literal.json"}}]}`,
		`{"type":"record","name":"Event","fields":[],"$ref":"literal.json"}`,
		`["null", "string"]`,
		`"string"`,
	} {
		for _, composed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/composed=%t", content, composed), func(t *testing.T) {
				cfg := filesConfig(map[string]string{
					"messages/event.yaml": "payload:\n  schemaFormat: application/vnd.apache.avro+json;version=1.9.0\n  schema: {$ref: '../schemas/event.avsc'}\n",
					"schemas/event.avsc":  content,
				})
				spec := []byte(header + "channels: {events: {messages: {event: {$ref: messages/event.yaml}}}}\n")
				out, err := runBundle(spec, cfg, composed)
				require.NoError(t, err)
				e := &engine{root: decode(t, out)}
				message, err := e.outputTarget("/channels/events/messages/event")
				require.NoError(t, err)
				var want, got any
				require.NoError(t, decode(t, []byte(content)).Decode(&want))
				require.NoError(t, value(t, message, "/payload/schema").Decode(&got))
				assert.Equal(t, want, got)
				assert.NotContains(t, string(out), "event.avsc")
				again, err := runBundle(out, nil, composed)
				require.NoError(t, err)
				assert.Equal(t, string(out), string(again))
			})
		}
	}
}

func TestBundleScalarExtensionFragments(t *testing.T) {
	for _, ref := range []string{"#/x-data/value", "data.yaml#/value"} {
		for _, composed := range []bool{false, true} {
			cfg := filesConfig(map[string]string{"data.yaml": "value: metadata\n"})
			cfg.ExcludeExtensionRefs = false
			spec := []byte(header + fmt.Sprintf("x-data: {value: metadata}\nx-copy: {$ref: '%s'}\n", ref))
			out, err := runBundle(spec, cfg, composed)
			require.NoError(t, err)
			e := &engine{root: decode(t, out)}
			copied, err := e.outputTarget("/x-copy")
			require.NoError(t, err)
			assert.Equal(t, "metadata", copied.Value)
		}
	}
}

func TestBundleInlineForeignSchemaMetadata(t *testing.T) {
	spec := []byte(header + `components:
  schemas:
    Event:
      schemaFormat: application/vnd.apache.avro+json;version=1.9.0
      schema: {type: record, name: Event, fields: [], $ref: literal.json}
`)
	for _, composed := range []bool{false, true} {
		out, err := runBundle(spec, nil, composed)
		require.NoError(t, err)
		assert.Equal(t, "literal.json", value(t, decode(t, out), "/components/schemas/Event/schema/$ref").Value)
		doc, err := libasyncapi.NewDocument(out)
		require.NoError(t, err)
		assert.Empty(t, doc.Errors())
	}
}

func TestBundleMultiFormatRemoteAccess(t *testing.T) {
	spec := []byte(header + `components:
  schemas:
    Event:
      schemaFormat: application/vnd.apache.avro+json;version=1.9.0
      schema: {$ref: 'https://example.com/event.avsc'}
`)
	for _, composed := range []bool{false, true} {
		calls := 0
		cfg := &libasyncapi.DocumentConfiguration{RemoteURLHandler: func(location string) (*http.Response, error) {
			calls++
			assert.Equal(t, "https://example.com/event.avsc", location)
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"type":"record","name":"Event","fields":[]}`)), Header: http.Header{"Content-Type": []string{"application/json"}}}, nil
		}}
		out, err := runBundle(spec, cfg, composed)
		require.Error(t, err)
		assert.Nil(t, out)
		assert.Zero(t, calls)
		cfg.AllowRemoteReferences = true
		out, err = runBundle(spec, cfg, composed)
		require.NoError(t, err)
		assert.Positive(t, calls)
		assert.Equal(t, "record", value(t, decode(t, out), "/components/schemas/Event/schema/type").Value)
	}
}

func TestBundleMultiFormatReferenceErrors(t *testing.T) {
	for _, reference := range []string{"missing.avsc", "loop.avsc", "[]", "''"} {
		for _, composed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/composed=%t", reference, composed), func(t *testing.T) {
				spec := []byte(header + fmt.Sprintf("components:\n  schemas:\n    Event:\n      schemaFormat: application/vnd.apache.avro+json;version=1.9.0\n      schema: {$ref: %s}\n", reference))
				out, err := runBundle(spec, filesConfig(map[string]string{"loop.avsc": "$ref: loop.avsc"}), composed)
				require.Error(t, err)
				assert.Nil(t, out)
			})
		}
	}
}
