// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBundleFileURIWithSpaces(t *testing.T) {
	base := filepath.Join(t.TempDir(), "spec files")
	require.NoError(t, os.Mkdir(base, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(base, "schema file.yaml"), []byte("type: string\n"), 0600))
	cfg := &libasyncapi.DocumentConfiguration{BasePath: base, SpecFilePath: "root file.yaml"}
	spec := []byte(header + "components:\n  schemas:\n    Text: {$ref: 'schema%20file.yaml'}\n")
	root, err := parse(spec)
	require.NoError(t, err)
	loader, err := newLoader(root, cfg)
	require.NoError(t, err)
	u, err := url.Parse(loader.root.uri)
	require.NoError(t, err)
	assert.Empty(t, u.Host)
	assert.Equal(t, "file", u.Scheme)
	assert.True(t, strings.HasPrefix(u.Path, "/"))
	for _, composed := range []bool{false, true} {
		out, err := runBundle(spec, cfg, composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, out)}
		schema, err := e.outputTarget("/components/schemas/Text")
		require.NoError(t, err)
		assert.Equal(t, "string", field(schema, "type").Value)
	}
}

func TestBundleReusableOperation(t *testing.T) {
	for _, composed := range []bool{false, true} {
		cfg := filesConfig(map[string]string{"channel.yaml": "address: a\nmessages: {event: {name: Event}}"})
		spec := header + "components:\n  operations:\n    reusable:\n      action: send\n      channel: {$ref: channel.yaml}\n      messages: [{$ref: 'channel.yaml#/messages/event'}]\n"
		out, err := runBundle([]byte(spec), cfg, composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "#/components/channels/channel", value(t, n, "/components/operations/reusable/channel/$ref").Value)
		assert.Equal(t, "#/components/channels/channel/messages/event", value(t, n, "/components/operations/reusable/messages/0/$ref").Value)
	}
}

func TestBundleExplicitRootReferenceDisambiguates(t *testing.T) {
	for _, composed := range []bool{false, true} {
		cfg := filesConfig(map[string]string{"channel.yaml": "address: a"})
		out, err := runBundle([]byte(header+"channels: {a: {$ref: channel.yaml}, b: {$ref: channel.yaml}}\noperations: {send: {action: send, channel: {$ref: '#/channels/a'}}}"), cfg, composed)
		require.NoError(t, err)
		assert.Equal(t, "#/channels/a", value(t, decode(t, out), "/operations/send/channel/$ref").Value)
	}
}

func TestBundleNamesAreNotKeywords(t *testing.T) {
	spec := header + "channels:\n  x-events:\n    messages:\n      event:\n        payload:\n          type: object\n          properties:\n            $ref: {$ref: scalar.yaml}\n            x-name: {$ref: scalar.yaml}\n"
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), filesConfig(map[string]string{"scalar.yaml": "type: string"}), composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, out)}
		for _, name := range []string{"$ref", "x-name"} {
			n, err := e.outputTarget("/channels/x-events/messages/event/payload/properties/" + name)
			require.NoError(t, err)
			assert.Equal(t, "string", field(n, "type").Value)
		}
	}
}

func TestBundleTraitsRepliesAndBindings(t *testing.T) {
	cfg := filesConfig(map[string]string{
		"schema.yaml":      "type: string",
		"tag.yaml":         "name: events",
		"docs.yaml":        "url: https://example.com",
		"security.yaml":    "type: userPassword",
		"variable.yaml":    "default: localhost",
		"parameter.yaml":   "description: id",
		"correlation.yaml": "location: '$message.header#/id'",
		"address.yaml":     "location: '$message.header#/replyTo'",
		"trait.yaml":       "summary: shared",
	})
	spec := header + `servers:
  broker:
    host: '{host}'
    protocol: kafka
    variables: {host: {$ref: variable.yaml}}
    bindings: {kafka: {schemaRegistryUrl: 'https://example.com'}}
channels:
  events:
    address: events
    parameters: {id: {$ref: parameter.yaml}}
    bindings: {ws: {query: {$ref: schema.yaml}}}
    messages:
      event:
        name: Event
        traits: [{$ref: trait.yaml}]
        headers: {$ref: schema.yaml}
        correlationId: {$ref: correlation.yaml}
operations:
  send:
    action: send
    channel: {$ref: '#/channels/events'}
    traits: [{$ref: trait.yaml}]
    tags: [{$ref: tag.yaml}]
    externalDocs: {$ref: docs.yaml}
    security: [{$ref: security.yaml}]
    bindings: {http: {query: {$ref: schema.yaml}}}
    reply:
      address: {$ref: address.yaml}
      channel: {$ref: '#/channels/events'}
      messages: [{$ref: '#/channels/events/messages/event'}]
`
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), cfg, composed)
		require.NoError(t, err)
		assert.NotContains(t, string(out), ".yaml")
		again, err := runBundle(out, nil, composed)
		require.NoError(t, err)
		assert.Equal(t, string(out), string(again))
	}
}

func TestBundleExtensionsAndBackReferences(t *testing.T) {
	spec := header + "channels: {event: {$ref: channel.yaml}}\ncomponents: {schemas: {ID: {type: string}}}\nx-import: {$ref: extension.yaml}\n"
	cfg := filesConfig(map[string]string{"channel.yaml": "messages: {event: {payload: {$ref: 'root.yaml#/components/schemas/ID'}}}", "extension.yaml": "custom: value"})
	cfg.ExcludeExtensionRefs = false
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), cfg, composed)
		require.NoError(t, err)
		assert.Equal(t, "value", value(t, decode(t, out), "/x-import/custom").Value)
		assert.Contains(t, string(out), "#/components/schemas/ID")
	}
}

func TestBundleInvalidInputs(t *testing.T) {
	for _, spec := range []string{
		"", "info: {}", "asyncapi: nonsense", "asyncapi: 4.0.0",
		header + "---\nasyncapi: 3.0.0", header + "channels: []", header + "components: []",
		header + "channels: {a: &a {nested: *a}}",
		header + "channels: {a: {$ref: '%zz'}}", header + "channels: {a: {$ref: 'ftp://example.com/a'}}",
		header + "channels: {a: {$ref: 'file://host/a.yaml'}}", header + "channels: {a: {$ref: '#/channels/~2'}}",
		header + "components: {schemas: {A: {$ref: '#/components/schemas/A'}}}",
		header + "components: {schemas: {A: {$dynamicRef: '#a'}}}",
		header + "components: {schemas: {A: {$id: 'https://example.com/a'}, B: {$id: 'https://example.com/a'}}}",
		header + "operations: {send: {messages: [{$ref: '#/components/messages/M'}]}}\ncomponents: {messages: {M: {name: M}}}",
	} {
		t.Run(fmt.Sprintf("%q", spec), func(t *testing.T) {
			out, err := BundleBytes([]byte(spec), nil)
			require.Error(t, err)
			assert.Nil(t, out)
		})
	}
	for _, external := range []string{"", "[", "one\n---\ntwo", "42"} {
		_, err := BundleBytes([]byte(header+"channels: {a: {$ref: a.yaml}}"), filesConfig(map[string]string{"a.yaml": external}))
		require.Error(t, err)
	}
	_, err := BundleBytes([]byte(header+"channels: {a: {$ref: a.yaml}}"), &libasyncapi.DocumentConfiguration{})
	require.Error(t, err)
}

func TestBundleAliasesAndMergePrecedence(t *testing.T) {
	spec := header + "channels:\n  a: &a\n    address: a\n    messages: {event: {payload: {$ref: schema.yaml}}}\n  b:\n    <<: *a\n    address: b\n"
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), filesConfig(map[string]string{"schema.yaml": "type: string"}), composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "a", value(t, n, "/channels/a/address").Value)
		assert.Equal(t, "b", value(t, n, "/channels/b/address").Value)
		assert.NotContains(t, string(out), "<<:")
		assert.NotContains(t, string(out), "*a")
		_, err = runBundle(out, nil, composed)
		require.NoError(t, err)
	}
}

func TestBundleConcurrentCalls(t *testing.T) {
	cfg := filesConfig(map[string]string{"schema.yaml": "type: object\nproperties: {next: {$ref: '#'}}"})
	spec := []byte(header + "components: {schemas: {A: {$ref: schema.yaml}}}")
	want, err := BundleBytesComposed(spec, cfg, nil)
	require.NoError(t, err)
	var wg sync.WaitGroup
	for range 16 {
		wg.Go(func() {
			out, err := BundleBytesComposed(spec, cfg, nil)
			assert.NoError(t, err)
			assert.Equal(t, string(want), string(out))
		})
	}
	wg.Wait()
}

func TestBundleDeepDocumentFails(t *testing.T) {
	spec := header + "x-deep: " + strings.Repeat("[", 520) + "0" + strings.Repeat("]", 520)
	_, err := BundleBytes([]byte(spec), nil)
	require.Error(t, err)
}

func TestBundleRecursivePlacementFails(t *testing.T) {
	spec := header + `servers: {s: {host: localhost, protocol: kafka}}
components:
  channels:
    A:
      servers: [{$ref: '#/servers/s'}]
      messages:
        recursive: {$ref: '#/components/channels/A'}
`
	_, err := BundleBytes([]byte(spec), nil)
	require.ErrorContains(t, err, "cyclic reference placement")
}

func TestBundlePreservesRootReferenceChain(t *testing.T) {
	spec := header + "components:\n  schemas:\n    A: {$ref: '#/components/schemas/B'}\n    B: {$ref: '#/components/schemas/C'}\n    C: {type: string}\n"
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), nil, composed)
		require.NoError(t, err)
		assert.Equal(t, "#/components/schemas/B", value(t, decode(t, out), "/components/schemas/A/$ref").Value)
	}
}

func TestBundleExternalSchemaFragmentKeepsParentScope(t *testing.T) {
	cfg := filesConfig(map[string]string{"schema.yaml": `$id: https://schemas.example.com/container
definitions:
  Text: {$anchor: text, type: string}
  Value:
    type: object
    properties:
      text: {$ref: '#text'}
`})
	for _, composed := range []bool{false, true} {
		spec := header + "components:\n  schemas:\n    Value: {$ref: 'schema.yaml#/definitions/Value'}\n    Text: {$ref: 'schema.yaml#text'}\n"
		out, err := runBundle([]byte(spec), cfg, composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, out)}
		for _, path := range []string{"/components/schemas/Text", "/components/schemas/Value/properties/text"} {
			n, err := e.outputTarget(path)
			require.NoError(t, err)
			assert.Equal(t, "string", field(n, "type").Value)
		}
		_, err = runBundle(out, nil, composed)
		require.NoError(t, err)
	}
}

func TestBundleRootChannelThroughExternalComponent(t *testing.T) {
	files := map[string]string{
		"reusable/channels.yaml": `lights:
  address: lights/{id}
  parameters:
    id: {description: light id}
  messages:
    on: {$ref: './messages.yaml#/on'}
    off: {$ref: './messages.yaml#/off'}
replies:
  address: replies
  messages:
    off: {$ref: './messages.yaml#/off'}
`,
		"reusable/messages.yaml":   "on: {name: on, payload: {type: object, properties: {cmd: {const: 'on'}}}}\noff: {name: off, payload: {type: object, properties: {cmd: {const: 'off'}}}}\n",
		"reusable/operations.yaml": "turnOn:\n  action: send\n  channel: {$ref: './channels.yaml#/lights'}\n  messages: [{$ref: './channels.yaml#/lights/messages/on'}]\n",
	}
	channels := "channels:\n  lights: {$ref: '#/components/channels/lights'}\n  replies: {$ref: '#/components/channels/replies'}\n"
	components := "components:\n  channels:\n    lights: {$ref: './reusable/channels.yaml#/lights'}\n    replies: {$ref: './reusable/channels.yaml#/replies'}\n"
	reusable := "operations:\n  turnOn: {$ref: '#/components/operations/turnOn'}\n"
	reusableRefs := map[string]string{
		"/operations/turnOn":                       "#/components/operations/turnOn",
		"/components/operations/turnOn/channel":    "#/channels/lights",
		"/components/operations/turnOn/messages/0": "#/channels/lights/messages/on",
	}
	tests := []struct {
		name, operations, components string
		refs                         map[string]string
	}{
		{"operation", `operations:
  turnOn:
    action: send
    channel: {$ref: '#/channels/lights'}
    messages: [{$ref: '#/channels/lights/messages/on'}]
`, "", map[string]string{
			"/operations/turnOn/channel":    "#/channels/lights",
			"/operations/turnOn/messages/0": "#/channels/lights/messages/on",
		}},
		{"reply", `operations:
  turnOn:
    action: send
    channel: {$ref: '#/channels/lights'}
    reply:
      channel: {$ref: '#/channels/replies'}
      messages: [{$ref: '#/channels/replies/messages/off'}]
`, "", map[string]string{
			"/operations/turnOn/channel":          "#/channels/lights",
			"/operations/turnOn/reply/channel":    "#/channels/replies",
			"/operations/turnOn/reply/messages/0": "#/channels/replies/messages/off",
		}},
		{"componentOperation", reusable, `  operations:
    turnOn:
      action: send
      channel: {$ref: '#/channels/lights'}
      messages: [{$ref: '#/channels/lights/messages/on'}]
`, reusableRefs},
		{"externalComponentOperation", reusable, "  operations:\n    turnOn: {$ref: './reusable/operations.yaml#/turnOn'}\n", reusableRefs},
	}
	for _, tt := range tests {
		for _, composed := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/%t", tt.name, composed), func(t *testing.T) {
				out, err := runBundle([]byte(header+channels+tt.operations+components+tt.components), filesConfig(files), composed)
				require.NoError(t, err)
				n := decode(t, out)
				for p, want := range tt.refs {
					assert.Equal(t, want, value(t, n, p+"/$ref").Value, p)
				}
				assert.Equal(t, "#/components/channels/lights", value(t, n, "/channels/lights/$ref").Value)
				assert.Equal(t, "#/components/channels/replies", value(t, n, "/channels/replies/$ref").Value)
				assert.Equal(t, "lights/{id}", value(t, n, "/components/channels/lights/address").Value)
				assert.Equal(t, "light id", value(t, n, "/components/channels/lights/parameters/id/description").Value)
				assert.Equal(t, "replies", value(t, n, "/components/channels/replies/address").Value)
				if composed {
					assert.Equal(t, "#/components/messages/on", value(t, n, "/components/channels/lights/messages/on/$ref").Value)
				} else {
					assert.Equal(t, "on", value(t, n, "/components/channels/lights/messages/on/payload/properties/cmd/const").Value)
				}
				e := &engine{root: n}
				for p, want := range map[string]string{"/channels/lights/messages/on": "on", "/channels/lights/messages/off": "off", "/channels/replies/messages/off": "off"} {
					message, err := e.outputTarget(p)
					require.NoError(t, err)
					assert.Equal(t, want, field(message, "name").Value, p)
				}
				assert.NotContains(t, string(out), ".yaml")
				again, err := runBundle(out, nil, composed)
				require.NoError(t, err)
				assert.Equal(t, string(out), string(again))
			})
		}
	}
}

func TestBundleExternalComponentChannelAmbiguity(t *testing.T) {
	files := map[string]string{
		"channel.yaml": "address: a\nmessages:\n  on: {$ref: message.yaml}\n  again: {$ref: message.yaml}\n",
		"message.yaml": "name: on\n",
	}
	shared := "channels:\n  a: {$ref: '#/components/channels/shared'}\n  b: {$ref: '#/components/channels/shared'}\n"
	component := "components:\n  channels:\n    shared: {$ref: channel.yaml}\n"
	tests := []struct{ name, spec, want string }{
		{"sharedComponent", shared + "operations: {send: {action: send, channel: {$ref: channel.yaml}}}\n" + component, "#/operations/send/channel: expected one local channels target, found 2"},
		{"componentReference", shared + "operations: {send: {action: send, channel: {$ref: '#/components/channels/shared'}}}\n" + component, "#/operations/send/channel: expected one local channels target, found 2"},
		{"twoComponents", "channels:\n  a: {$ref: '#/components/channels/x'}\n  b: {$ref: '#/components/channels/y'}\noperations: {send: {action: send, channel: {$ref: channel.yaml}}}\ncomponents:\n  channels:\n    x: {$ref: channel.yaml}\n    y: {$ref: channel.yaml}\n", "#/operations/send/channel: expected one local channels target, found 2"},
		{"messages", shared + "operations: {send: {action: send, channel: {$ref: '#/channels/a'}, messages: [{$ref: message.yaml}]}}\n" + component, "#/operations/send/messages/0: expected one local messages target, found 2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, composed := range []bool{false, true} {
				out, err := runBundle([]byte(header+tt.spec), filesConfig(files), composed)
				require.ErrorContains(t, err, tt.want)
				assert.Nil(t, out)
			}
		})
	}
	// Explicit root references still select one of the shared uses.
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(header+shared+"operations: {send: {action: send, channel: {$ref: '#/channels/b'}, messages: [{$ref: '#/channels/b/messages/again'}]}}\n"+component), filesConfig(files), composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "#/channels/b", value(t, n, "/operations/send/channel/$ref").Value)
		assert.Equal(t, "#/channels/b/messages/again", value(t, n, "/operations/send/messages/0/$ref").Value)
	}
}
