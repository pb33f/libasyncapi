package libasyncapi

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestDocumentResolvesChannelAliasMessages(t *testing.T) {
	for _, missing := range []bool{false, true} {
		for _, key := range []string{"c", "c~1d~0e"} {
			name := strings.ReplaceAll(strings.ReplaceAll(key, "~1", "/"), "~0", "~")
			source := fmt.Sprintf(`asyncapi: 3.1.0
info: {title: aliases, version: '1'}
channels:
  %s: {$ref: '#/components/channels/alias'}
operations:
  receive:
    action: receive
    channel: {$ref: '#/channels/%s'}
    messages: [{$ref: '#/channels/%s/messages/m'}]
components:
  channels:
    alias: {$ref: '#/components/channels/target'}
    target: {messages: {m: {name: event}}}
`, name, key, key)
			if missing {
				source += "  messages: {bad: {$ref: '#/channels/" + key + "/messages/missing'}}\n"
			}
			doc, err := NewDocument([]byte(source))
			require.NoError(t, err)
			if missing {
				require.NotEmpty(t, doc.Errors())
				for _, issue := range doc.Errors() {
					assert.Contains(t, issue.Error(), "missing")
				}
			} else {
				assert.Empty(t, doc.Errors())
			}
			target := doc.Index().FindComponentInRoot(t.Context(), "#/components/channels/target/messages/m")
			resolved, _, _ := doc.Index().SearchIndexForReferenceWithContext(t.Context(), "#/channels/"+key+"/messages/m")
			require.NotNil(t, target)
			require.NotNil(t, resolved)
			assert.Same(t, target.Node, resolved.Node)
			var before yaml.Node
			require.NoError(t, yaml.Unmarshal([]byte(source), &before))
			assert.Equal(t, &before, doc.RootNode())
		}
	}
}

func TestDocumentJSONSurrogateEscapes(t *testing.T) {
	source := []byte(`{"asyncapi":"3.0.0","info":{"title":"\ud83c\udf03","version":"1"}}`)
	version, err := DetectAsyncAPIVersion(source)
	require.NoError(t, err)
	assert.Equal(t, "3.0.0", version)
	doc, err := NewDocument(source)
	require.NoError(t, err)
	assert.Empty(t, doc.Errors())
	assert.Equal(t, "🌃", doc.Model().Info.Title)
}

func TestDocumentAliasFailuresRemainVisible(t *testing.T) {
	for _, tc := range []struct{ name, channel, messages string }{
		{"invalid escape", "{$ref: '#/components/channels/target'}", "[{$ref: '#/channels/c~2d/messages/m'}]"},
		{"duplicate ref", "{$ref: '#/components/channels/target'}", "[{$ref: '#/channels/c~02d/messages/m', $ref: '#/channels/c~02d/messages/missing'}]"},
		{"cycle", "{$ref: '#/channels/c~02d'}", "[{$ref: '#/channels/c~02d/messages/m'}]"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := fmt.Sprintf("asyncapi: 3.1.0\nchannels: {'c~2d': %s}\noperations: {send: {action: send, channel: {$ref: '#/channels/c~02d'}, messages: %s}}\ncomponents: {channels: {target: {messages: {m: {name: event}}}}}", tc.channel, tc.messages)
			doc, err := NewDocument([]byte(source))
			require.NoError(t, err)
			require.NotEmpty(t, doc.Errors())
			if tc.name == "duplicate ref" {
				assert.Contains(t, fmt.Sprint(doc.Errors()), "missing")
			}
		})
	}
}

func TestDocumentExternalChannelAlias(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "channel.yaml"), []byte("messages: {m: {name: external}}\n"), 0600))
	config := NewDocumentConfiguration()
	config.BasePath = dir
	source := "asyncapi: 3.1.0\nchannels: {c: {$ref: channel.yaml}}\noperations: {send: {action: send, channel: {$ref: '#/channels/c'}, messages: [{$ref: '#/channels/c/messages/m'}]}}\n"
	doc, err := NewDocumentWithConfiguration([]byte(source), config)
	require.NoError(t, err)
	assert.Empty(t, doc.Errors())
	ref, owner, _ := doc.Index().SearchIndexForReferenceWithContext(t.Context(), "#/channels/c/messages/m")
	require.NotNil(t, ref)
	assert.NotSame(t, doc.Index(), owner)
	assert.Equal(t, "external", ref.Node.Content[1].Value)
}

func TestForeignSchemaMetadataIsNotAReference(t *testing.T) {
	source := []byte(`asyncapi: 3.1.0
channels: {c: {messages: {m: {$ref: '#/components/messages/m'}}}}
components:
  messages:
    m:
      payload:
        schemaFormat: application/vnd.apache.avro+json;version=1.9.0
        schema: {type: record, name: Event, fields: [], $ref: 'https://example.test/metadata'}
`)
	calls := 0
	config := NewDocumentConfiguration()
	config.AllowRemoteReferences = true
	config.RemoteURLHandler = func(string) (*http.Response, error) { calls++; return nil, fmt.Errorf("metadata was fetched") }
	doc, err := NewDocumentWithConfiguration(source, config)
	require.NoError(t, err)
	assert.Empty(t, doc.Errors())
	assert.Zero(t, calls)
	assert.Same(t, doc.RootNode(), doc.Index().GetRootNode())
	var expected yaml.Node
	require.NoError(t, yaml.Unmarshal(source, &expected))
	assert.Equal(t, &expected, doc.RootNode())
	ref, _, _ := doc.Index().SearchIndexForReferenceWithContext(t.Context(), "#/components/messages/m")
	target := doc.RootNode().Content[0].Content[5].Content[1].Content[1]
	require.NotNil(t, ref)
	assert.Same(t, target, ref.Node)

	missing := []byte("asyncapi: 3.1.0\ncomponents: {schemas: {S: {schemaFormat: application/vnd.apache.avro+json, schema: {$ref: missing.avsc}}}}")
	doc, err = NewDocument(missing)
	require.NoError(t, err)
	assert.NotEmpty(t, doc.Errors())
}
