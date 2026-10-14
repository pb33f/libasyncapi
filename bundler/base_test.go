// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestBundleBaseOverridePreservesEntryIdentity(t *testing.T) {
	for _, remote := range []bool{false, true} {
		for _, composed := range []bool{false, true} {
			t.Run(fmt.Sprintf("remote=%t/composed=%t", remote, composed), func(t *testing.T) {
				dir := t.TempDir()
				entry := filepath.Join(dir, "v2", "asyncapi.yaml")
				base := filepath.Join(dir, "common")
				rootRef := (&url.URL{Scheme: "file", Path: "/" + strings.TrimPrefix(filepath.ToSlash(entry), "/")}).String() + "#/components/messages/m"
				files := map[string]string{
					"ch.yaml":           "messages:\n  imported: {$ref: './asyncapi.yaml#/components/messages/m'}\n  entry: {$ref: '" + rootRef + "'}\n",
					"asyncapi.yaml":     "components: {messages: {m: {payload: {type: integer}}}}\n",
					"schema.yaml":       "type: boolean\n",
					"nested/value.yaml": "type: number\n",
				}
				cfg := &libasyncapi.DocumentConfiguration{SpecFilePath: entry, AllowFileReferences: true}
				if remote {
					cfg.BaseURL, _ = url.Parse("https://example.com/common")
					cfg.AllowRemoteReferences = true
					cfg.RemoteURLHandler = func(location string) (*http.Response, error) {
						name := strings.TrimPrefix(location, "https://example.com/common/")
						content, ok := files[name]
						if !ok {
							return nil, fmt.Errorf("unexpected URL %s", location)
						}
						return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(content)), Header: http.Header{}}, nil
					}
				} else {
					cfg.BasePath = base
					for name, content := range files {
						path := filepath.Join(base, name)
						require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
						require.NoError(t, os.WriteFile(path, []byte(content), 0600))
					}
				}
				spec := header + `channels:
  imported: {$ref: ch.yaml}
  local:
    messages:
      own: {$ref: '#/components/messages/m'}
operations:
  send:
    action: send
    channel: {$ref: '#/channels/local'}
components:
  messages:
    m: {payload: {type: string}}
  schemas:
    Imported: {$ref: schema.yaml}
    WithID:
      $id: nested/container.yaml
      type: object
      properties:
        value: {$ref: value.yaml}
    Local: {$ref: '#/components/schemas/Own'}
    Own: {type: 'null'}
`
				out, err := runBundle([]byte(spec), cfg, composed)
				require.NoError(t, err)
				e := &engine{root: decode(t, out)}
				for path, kind := range map[string]string{
					"/channels/imported/messages/imported/payload": "integer",
					"/channels/imported/messages/entry/payload":    "string",
					"/channels/local/messages/own/payload":         "string",
					"/components/schemas/Imported":                 "boolean",
					"/components/schemas/WithID/properties/value":  "number",
					"/components/schemas/Local":                    "null",
				} {
					n, err := e.outputTarget(path)
					require.NoError(t, err, path)
					require.NotNil(t, field(n, "type"), path)
					assert.Equal(t, kind, field(n, "type").Value, path)
				}
				assert.Equal(t, "#/channels/local", value(t, decode(t, out), "/operations/send/channel/$ref").Value)
			})
		}
	}
}

func TestBundleBaseOverridePreservesExplicitSchemaID(t *testing.T) {
	dir := t.TempDir()
	entry := filepath.Join(dir, "v2", "asyncapi.yaml")
	base := filepath.Join(dir, "common")
	for path, kind := range map[string]string{
		filepath.Join(filepath.Dir(entry), "value.yaml"):           "string",
		filepath.Join(base, "value.yaml"):                          "integer",
		filepath.Join(filepath.Dir(entry), "nested", "value.yaml"): "boolean",
		filepath.Join(base, "nested", "value.yaml"):                "number",
	} {
		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
		require.NoError(t, os.WriteFile(path, []byte("type: "+kind), 0600))
	}
	id := (&url.URL{Scheme: "file", Path: "/" + strings.TrimPrefix(filepath.ToSlash(entry), "/")}).String()
	spec := []byte(header + "components:\n  schemas:\n    Scoped:\n      $id: " + id + "\n      type: object\n      properties:\n        value: {$ref: value.yaml}\n        alias: {$ref: '#/properties/value'}\n        nested:\n          $id: nested/container.yaml\n          properties:\n            value: {$ref: value.yaml}\n")
	for _, composed := range []bool{false, true} {
		out, err := runBundle(spec, &libasyncapi.DocumentConfiguration{SpecFilePath: entry, BasePath: base, AllowFileReferences: true}, composed)
		require.NoError(t, err)
		e := &engine{root: decode(t, out)}
		n, err := e.outputTarget("/components/schemas/Scoped/properties/value")
		require.NoError(t, err)
		assert.Equal(t, "string", field(n, "type").Value)
		n, err = e.outputTarget("/components/schemas/Scoped/properties/alias")
		require.NoError(t, err)
		assert.Equal(t, "string", field(n, "type").Value)
		n, err = e.outputTarget("/components/schemas/Scoped/properties/nested/properties/value")
		require.NoError(t, err)
		assert.Equal(t, "boolean", field(n, "type").Value)
	}
}
