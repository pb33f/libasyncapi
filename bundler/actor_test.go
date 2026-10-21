// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"testing"

	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBundleOperationsOnlyActor(t *testing.T) {
	files := map[string]string{
		"domain.yaml": `servers:
  broker: {host: broker.example.com, protocol: kafka}
channels:
  events:
    address: events
    servers: [{$ref: '#/servers/broker'}]
    messages:
      created: {$ref: '#/components/messages/created'}
components:
  securitySchemes:
    OAuth2:
      type: oauth2
      flows: &flows
        clientCredentials: {tokenUrl: 'https://auth.example.com', availableScopes: {read: Read}}
  operationTraits:
    reader:
      security: [{type: oauth2, flows: *flows}]
  messages:
    created:
      name: Created
      payload: {$ref: 'payload.yaml#/Event'}
`,
		"payload.yaml": "Event: {type: object, properties: {previous: {$ref: '#/Event'}}}\n",
	}
	operation := `action: receive
channel: {$ref: 'domain.yaml#/channels/events'}
messages: [{$ref: 'domain.yaml#/channels/events/messages/created'}]
traits: [{$ref: 'domain.yaml#/components/operationTraits/reader'}]
`
	files["operation.yaml"] = operation
	for _, composed := range []bool{false, true} {
		for _, reusable := range []bool{false, true} {
			t.Run(fmt.Sprintf("composed=%t/reusable=%t", composed, reusable), func(t *testing.T) {
				spec := header + "operations: {receive: {$ref: operation.yaml}, again: {$ref: operation.yaml}}\n"
				if reusable {
					spec = header + "operations: {receive: {$ref: '#/components/operations/receive'}, again: {$ref: '#/components/operations/receive'}}\ncomponents:\n  operations:\n    receive: {$ref: operation.yaml}\n"
				}
				out, err := runBundle([]byte(spec), filesConfig(files), composed)
				require.NoError(t, err)
				n := decode(t, out)
				e := &engine{root: n}
				assert.Len(t, value(t, n, "/channels").Content, 2)
				assert.Len(t, value(t, n, "/servers").Content, 2)
				for _, name := range []string{"receive", "again"} {
					op, err := e.outputTarget("/operations/" + name)
					require.NoError(t, err)
					assert.Equal(t, "#/channels/events", value(t, op, "/channel/$ref").Value)
					assert.Equal(t, "#/channels/events/messages/created", value(t, op, "/messages/0/$ref").Value)
					trait, err := e.outputTarget("/operations/" + name + "/traits/0")
					require.NoError(t, err)
					assert.Equal(t, "https://auth.example.com", value(t, trait, "/security/0/flows/clientCredentials/tokenUrl").Value)
				}
				channel, err := e.outputTarget("/channels/events")
				require.NoError(t, err)
				assert.Equal(t, "events", field(channel, "address").Value)
				assert.Equal(t, "#/servers/broker", value(t, channel, "/servers/0/$ref").Value)
				payload, err := e.outputTarget("/channels/events/messages/created/payload/properties/previous")
				require.NoError(t, err)
				assert.Equal(t, "object", field(payload, "type").Value)
				assert.NotContains(t, string(out), ".yaml")
				assert.NotContains(t, string(out), "OAuth2:")
				again, err := runBundle(out, nil, composed)
				require.NoError(t, err)
				assert.Equal(t, string(out), string(again))
			})
		}
	}
}

func TestBundleActorRootNameCollision(t *testing.T) {
	for _, composed := range []bool{false, true} {
		t.Run(fmt.Sprint(composed), func(t *testing.T) {
			files := map[string]string{
				"one/domain.yaml": "channels: {events: {address: one, servers: [{$ref: '#/servers/broker'}]}}\nservers: {broker: {host: one, protocol: kafka}}\n",
				"two/domain.yaml": "channels: {events: {address: two, servers: [{$ref: '#/servers/broker'}]}}\nservers: {broker: {host: two, protocol: kafka}}\n",
			}
			spec := header + `channels: {events: {address: existing}}
servers: {broker: {host: existing, protocol: kafka}}
operations:
  one: {action: send, channel: {$ref: 'one/domain.yaml#/channels/events'}}
  two: {action: send, channel: {$ref: 'two/domain.yaml#/channels/events'}}
`
			out, err := runBundle([]byte(spec), filesConfig(files), composed)
			require.NoError(t, err)
			n := decode(t, out)
			e := &engine{root: n}
			assert.Equal(t, "existing", value(t, n, "/channels/events/address").Value)
			assert.Equal(t, "existing", value(t, n, "/servers/broker/host").Value)
			for _, tt := range []struct{ op, suffix string }{{"one", ""}, {"two", "__1"}} {
				path := "/channels/events__domain" + tt.suffix
				assert.Equal(t, "#"+path, value(t, n, "/operations/"+tt.op+"/channel/$ref").Value)
				channel, err := e.outputTarget(path)
				require.NoError(t, err)
				assert.Equal(t, tt.op, field(channel, "address").Value)
				assert.Equal(t, "#/servers/broker__domain"+tt.suffix, value(t, channel, "/servers/0/$ref").Value)
			}
		})
	}
}

func TestBundleActorDoesNotImportUnrelatedMessages(t *testing.T) {
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(header+"operations: {send: {channel: {$ref: channel.yaml}, messages: [{$ref: message.yaml}]}}\n"), filesConfig(map[string]string{
			"channel.yaml": "address: events\nmessages: {event: {name: Event}}\n",
			"message.yaml": "name: Other\n",
		}), composed)
		require.ErrorContains(t, err, "expected one local messages target, found 0")
		assert.Nil(t, out)
	}
}

func TestBundleActorReusesComponentsAndImportsReply(t *testing.T) {
	files := map[string]string{
		"domain.yaml": "channels: {events: {address: events, servers: [{$ref: '#/servers/broker'}]}, replies: {address: replies}}\nservers: {broker: {host: broker, protocol: kafka}}\n",
	}
	// Components precede root operations. Their server edges only become root
	// uses after the referenced channel has been imported.
	spec := header + `components:
  channels:
    shared: {$ref: 'domain.yaml#/channels/events'}
operations:
  send:
    action: send
    channel: {$ref: '#/components/channels/shared'}
    reply:
      channel: {$ref: 'domain.yaml#/channels/replies'}
`
	for _, composed := range []bool{false, true} {
		out, err := runBundle([]byte(spec), filesConfig(files), composed)
		require.NoError(t, err)
		n := decode(t, out)
		assert.Equal(t, "#/channels/events", value(t, n, "/operations/send/channel/$ref").Value)
		assert.Equal(t, "#/channels/replies", value(t, n, "/operations/send/reply/channel/$ref").Value)
		e := &engine{root: n}
		channel, err := e.outputTarget("/channels/events")
		require.NoError(t, err)
		assert.Equal(t, "#/servers/broker", value(t, channel, "/servers/0/$ref").Value)
		if composed {
			assert.Equal(t, "#/components/channels/shared", value(t, n, "/channels/events/$ref").Value)
		}
	}
}
