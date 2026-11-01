package what_changed_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/pb33f/libasyncapi"
	wc "github.com/pb33f/libasyncapi/what-changed"
	"github.com/pb33f/libasyncapi/what-changed/model"
	"github.com/pb33f/testify/require"
)

func TestMultiFormatSchemaChanges(t *testing.T) {
	for _, location := range []string{"payload", "component"} {
		for _, format := range []string{"application/schema+json;version=draft-07", "application/vnd.apache.avro;version=1.9.0"} {
			t.Run(location+"/"+format, func(t *testing.T) {
				old := wrappedSchemaSpec(location, format, "string")
				updated := wrappedSchemaSpec(location, format, "int")
				if strings.Contains(format, "schema+json") {
					updated = wrappedSchemaSpec(location, format, "integer")
				}
				left, err := libasyncapi.NewDocument([]byte(old))
				require.NoError(t, err)
				right, err := libasyncapi.NewDocument([]byte(updated))
				require.NoError(t, err)
				require.False(t, left.IsPartial())
				require.False(t, right.IsPartial())
				changes := wc.CompareDocuments(left, right)
				require.NotNil(t, changes)
				require.Equal(t, 1, changes.TotalChanges())
				require.Equal(t, 1, changes.TotalBreakingChanges())
				change := changes.GetAllChanges()[0]
				require.NotNil(t, change.Context.OriginalLine)
				require.NotNil(t, change.Context.NewLine)
				if strings.Contains(format, "schema+json") {
					require.Equal(t, "type", change.Property)
				} else {
					require.Equal(t, "schema", change.Property)
					require.Contains(t, change.OriginalEncoded, "string")
				}
				require.Nil(t, wc.CompareDocuments(left, left))
			})
		}
	}
}

func TestMultiFormatSchemaRulesAndFormatChange(t *testing.T) {
	old := wrappedSchemaSpec("payload", "application/vnd.apache.avro;version=1.9.0", "string")
	updated := strings.Replace(old, "type: string", "type: int", 1)
	left, err := libasyncapi.NewDocument([]byte(old))
	require.NoError(t, err)
	right, err := libasyncapi.NewDocument([]byte(updated))
	require.NoError(t, err)
	config := model.NewDefaultBreakingRulesConfig()
	no := false
	config.MultiFormatSchema.Modified = &no
	changes := wc.CompareDocumentsWithConfig(left, right, config)
	require.NotNil(t, changes)
	require.Zero(t, changes.TotalBreakingChanges())
	right, err = libasyncapi.NewDocument([]byte(strings.Replace(old, "version=1.9.0", "version=1.10.0", 1)))
	require.NoError(t, err)
	changes = wc.CompareDocuments(left, right)
	require.NotNil(t, changes)
	require.Equal(t, "schemaFormat", changes.GetAllChanges()[0].Property)
}

func TestMultiFormatSchemaIgnoresFormatting(t *testing.T) {
	spec := wrappedSchemaSpec("payload", "application/vnd.apache.avro;version=1.9.0", "string")
	left, err := libasyncapi.NewDocument([]byte(spec))
	require.NoError(t, err)
	right, err := libasyncapi.NewDocument([]byte("\n" + spec))
	require.NoError(t, err)
	require.Nil(t, wc.CompareDocuments(left, right))
}

func TestMultiFormatSchemaExternalContent(t *testing.T) {
	for _, format := range []string{"application/schema+json;version=draft-07", "application/vnd.apache.avro+json;version=1.9.0"} {
		t.Run(format, func(t *testing.T) {
			spec := wrappedSchemaSpec("payload", format, "string")
			spec = strings.Replace(spec, "type: string", "$ref: payload.yaml", 1)
			var docs [2]libasyncapi.Document
			for i, kind := range []string{"string", "integer"} {
				dir := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(dir, "payload.yaml"), []byte("type: "+kind+"\n"), 0600))
				config := libasyncapi.NewDocumentConfiguration()
				config.BasePath = dir
				config.AllowFileReferences = true
				var err error
				docs[i], err = libasyncapi.NewDocumentWithConfiguration([]byte(spec), config)
				require.NoError(t, err)
				require.False(t, docs[i].IsPartial())
			}
			changes := wc.CompareDocuments(docs[0], docs[1])
			require.NotNil(t, changes)
			require.Equal(t, 1, changes.TotalChanges())
			if strings.Contains(format, "schema+json") {
				require.Equal(t, "type", changes.GetAllChanges()[0].Property)
			} else {
				require.Equal(t, "schema", changes.GetAllChanges()[0].Property)
			}
			require.True(t, strings.HasSuffix(changes.GetAllChanges()[0].Context.DocumentLocation, "payload.yaml"))
		})
	}
}

func TestNativeExternalPayloadAfterSchemaAccess(t *testing.T) {
	for _, warm := range []bool{false, true} {
		t.Run(fmt.Sprint(warm), func(t *testing.T) {
			spec := "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nchannels:\n  events:\n    address: events\n    messages:\n      event:\n        payload:\n          $ref: payload.yaml\n"
			var docs [2]libasyncapi.Document
			for i, kind := range []string{"string", "integer"} {
				dir := t.TempDir()
				require.NoError(t, os.WriteFile(filepath.Join(dir, "payload.yaml"), []byte("type: "+kind+"\n"), 0600))
				config := libasyncapi.NewDocumentConfiguration()
				config.BasePath = dir
				config.AllowFileReferences = true
				var err error
				docs[i], err = libasyncapi.NewDocumentWithConfiguration([]byte(spec), config)
				require.NoError(t, err)
				if warm {
					for _, channel := range docs[i].Model().Channels.FromOldest() {
						for _, message := range channel.Messages.FromOldest() {
							require.NotNil(t, message.Payload.Schema())
						}
					}
				}
			}
			changes := wc.CompareDocuments(docs[0], docs[1])
			require.NotNil(t, changes)
			require.Equal(t, 1, changes.TotalChanges())
		})
	}
}

func wrappedSchemaSpec(location, format, kind string) string {
	root := "asyncapi: 3.0.0\ninfo:\n  title: Synthetic events\n  version: 1.0.0\n"
	if location == "component" {
		return root + fmt.Sprintf("components:\n  schemas:\n    Event:\n      schemaFormat: %s\n      schema:\n        type: %s\n", format, kind)
	}
	return root + fmt.Sprintf("channels:\n  events:\n    address: events\n    messages:\n      event:\n        payload:\n          schemaFormat: %s\n          schema:\n            type: %s\n", format, kind)
}

func TestSkipCircularCheckSelfReference(t *testing.T) {
	root := "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n"
	schema := "components:\n  schemas:\n    Node:\n      type: object\n      properties:\n        next: {$ref: '#/components/schemas/Node'}\n"
	mutual := "components:\n  schemas:\n    A:\n      type: object\n      properties:\n        b: {$ref: '#/components/schemas/B'}\n    B:\n      type: object\n      properties:\n        a: {$ref: '#/components/schemas/A'}\n"
	message := "channels:\n  events:\n    address: events\n    messages:\n      event:\n        payload: {$ref: '#/components/schemas/Node'}\n"
	for name, spec := range map[string]string{"component": root + schema, "mutual": root + mutual, "payload": root + message + schema} {
		t.Run(name, func(t *testing.T) {
			config := libasyncapi.NewDocumentConfiguration()
			config.SkipCircularReferenceCheck = true
			left, err := libasyncapi.NewDocumentWithConfiguration([]byte(spec), config)
			require.NoError(t, err)
			right, err := libasyncapi.NewDocumentWithConfiguration([]byte(strings.Replace(spec, "type: object", "type: object\n      title: changed", 1)), config)
			require.NoError(t, err)
			changes := wc.CompareDocuments(left, right)
			require.NotNil(t, changes)
			require.Equal(t, 1, changes.TotalChanges())
			require.Equal(t, "title", changes.GetAllChanges()[0].Property)
		})
	}
}

func TestMultiFormatHeaders(t *testing.T) {
	const avro, jsonSchema = "application/vnd.apache.avro;version=1.9.0", "application/schema+json;version=draft-07"
	// Traits have no identity, so an edited list entry is one removal plus one addition.
	locations := map[string]struct {
		path string
		list bool
	}{
		"message":           {path: "channels:\n  events:\n    address: events\n    messages:\n      event:\n        headers:\n%s"},
		"component":         {path: "components:\n  messages:\n    event:\n      headers:\n%s"},
		"trait":             {path: "components:\n  messageTraits:\n    common:\n      headers:\n%s"},
		"component channel": {path: "components:\n  channels:\n    events:\n      address: events\n      messages:\n        event:\n          headers:\n%s"},
		"inline trait":      {path: "components:\n  messages:\n    event:\n      traits:\n        - headers:\n%s", list: true},
		"channel trait":     {path: "channels:\n  events:\n    address: events\n    messages:\n      event:\n        traits:\n          - headers:\n%s", list: true},
	}
	for name, location := range locations {
		spec := func(format, schema string) string {
			body := "              schemaFormat: " + format + "\n              schema: " + schema + "\n"
			return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n" + fmt.Sprintf(location.path, body)
		}
		compare := func(t *testing.T, old, updated, property string) {
			left, err := libasyncapi.NewDocument([]byte(old))
			require.NoError(t, err)
			right, err := libasyncapi.NewDocument([]byte(updated))
			require.NoError(t, err)
			require.False(t, left.IsPartial(), "%v", left.Errors())
			require.False(t, right.IsPartial(), "%v", right.Errors())
			changes := wc.CompareDocuments(left, right)
			require.NotNil(t, changes)
			if location.list {
				property = "traits"
			}
			for _, change := range changes.GetAllChanges() {
				require.Equal(t, property, change.Property)
			}
			require.Nil(t, wc.CompareDocuments(left, left))
		}
		t.Run(name+"/avro", func(t *testing.T) {
			compare(t, spec(avro, "{type: string}"), spec(avro, "{type: int}"), "schema")
		})
		t.Run(name+"/json", func(t *testing.T) {
			compare(t, spec(jsonSchema, "{type: string}"), spec(jsonSchema, "{type: integer}"), "type")
		})
		t.Run(name+"/avro named reference", func(t *testing.T) {
			record := "{type: record, name: Headers, fields: [{name: %s, type: {$ref: com.x.Y}}]}"
			compare(t, spec(avro, fmt.Sprintf(record, "a")), spec(avro, fmt.Sprintf(record, "b")), "schema")
		})
	}
	t.Run("single change", func(t *testing.T) {
		spec := "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nchannels:\n  events:\n    address: events\n    messages:\n      event:\n        headers:\n          schemaFormat: " + avro + "\n          schema: {type: string}\n"
		left, err := libasyncapi.NewDocument([]byte(spec))
		require.NoError(t, err)
		right, err := libasyncapi.NewDocument([]byte(strings.Replace(spec, "string", "int", 1)))
		require.NoError(t, err)
		changes := wc.CompareDocuments(left, right)
		require.NotNil(t, changes)
		require.Equal(t, 1, changes.TotalChanges())
		require.Equal(t, 1, changes.TotalBreakingChanges())
	})
}

func TestJSONSchemaWrapperNonObjectContent(t *testing.T) {
	const format = "application/schema+json;version=draft-07"
	spec := func(content string) string {
		return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nchannels:\n  events:\n    address: events\n    messages:\n      event:\n        payload:\n          schemaFormat: " + format + "\n          schema: " + content + "\n"
	}
	for _, test := range [][2]string{{"true", "false"}, {"true", "{type: string}"}, {"{type: string}", "false"}} {
		t.Run(test[0]+" to "+test[1], func(t *testing.T) {
			left, err := libasyncapi.NewDocument([]byte(spec(test[0])))
			require.NoError(t, err)
			right, err := libasyncapi.NewDocument([]byte(spec(test[1])))
			require.NoError(t, err)
			changes := wc.CompareDocuments(left, right)
			require.NotNil(t, changes)
			require.Equal(t, 1, changes.TotalChanges())
			require.Equal(t, "schema", changes.GetAllChanges()[0].Property)
			require.Nil(t, wc.CompareDocuments(left, left))
		})
	}
	left, err := libasyncapi.NewDocument([]byte(spec("true")))
	require.NoError(t, err)
	right, err := libasyncapi.NewDocument([]byte("\n" + spec("true")))
	require.NoError(t, err)
	require.Nil(t, wc.CompareDocuments(left, right))
}

func TestJSONSchemaWrapperAliasedContent(t *testing.T) {
	spec := func(kind string) string {
		return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nx-shared: &shared {type: " + kind + "}\nchannels:\n  events:\n    address: events\n    messages:\n      event:\n        payload:\n          schemaFormat: application/schema+json;version=draft-07\n          schema: *shared\n"
	}
	left, err := libasyncapi.NewDocument([]byte(spec("string")))
	require.NoError(t, err)
	right, err := libasyncapi.NewDocument([]byte(spec("integer")))
	require.NoError(t, err)
	var properties []string
	for _, change := range wc.CompareDocuments(left, right).GetAllChanges() {
		properties = append(properties, change.Property)
	}
	require.Contains(t, properties, "type")
}

func TestMultiFormatCircularJSONContent(t *testing.T) {
	spec := func(title string) string {
		return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\ncomponents:\n  schemas:\n    W:\n      schemaFormat: application/schema+json;version=draft-07\n      schema:\n        title: " + title + "\n        type: object\n        properties:\n          next: {$ref: '#/components/schemas/W/schema'}\n"
	}
	for _, skip := range []bool{false, true} {
		t.Run(fmt.Sprint(skip), func(t *testing.T) {
			config := libasyncapi.NewDocumentConfiguration()
			config.SkipCircularReferenceCheck = skip
			left, _ := libasyncapi.NewDocumentWithConfiguration([]byte(spec("old")), config)
			right, _ := libasyncapi.NewDocumentWithConfiguration([]byte(spec("new")), config)
			changes := wc.CompareDocuments(left, right)
			require.NotNil(t, changes)
			require.Equal(t, 1, changes.TotalChanges())
		})
	}
}

func TestSchemaReferenceRename(t *testing.T) {
	const components = "components:\n  schemas:\n    Alias: {$ref: '#/components/schemas/%s'}\n    S: {type: string}\n    S2: {type: string}\n    W:\n      schemaFormat: application/vnd.apache.avro;version=1.9.0\n      schema: {type: string}\n    W2:\n      schemaFormat: application/vnd.apache.avro;version=1.9.0\n      schema: {type: string}\n"
	message := func(body string) string {
		return "channels:\n  events:\n    address: events\n    messages:\n      event:\n" + body
	}
	for name, test := range map[string]struct {
		spec, old, updated string
	}{
		"payload":      {message("        payload: {$ref: '#/components/schemas/%s'}\n"), "S", "S2"},
		"headers":      {message("        headers: {$ref: '#/components/schemas/%s'}\n"), "S", "S2"},
		"wrapper":      {message("        payload: {$ref: '#/components/schemas/%s'}\n"), "W", "W2"},
		"kafka key":    {message("        bindings:\n          kafka:\n            key: {$ref: '#/components/schemas/%s'}\n"), "S", "S2"},
		"http headers": {message("        bindings:\n          http:\n            headers: {$ref: '#/components/schemas/%s'}\n"), "S", "S2"},
	} {
		t.Run(name, func(t *testing.T) {
			spec := func(target string) string {
				return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n" + fmt.Sprintf(test.spec, target) + fmt.Sprintf(components, "S")
			}
			for _, skip := range []bool{false, true} {
				config := libasyncapi.NewDocumentConfiguration()
				config.SkipCircularReferenceCheck = skip
				left, err := libasyncapi.NewDocumentWithConfiguration([]byte(spec(test.old)), config)
				require.NoError(t, err)
				right, err := libasyncapi.NewDocumentWithConfiguration([]byte(spec(test.updated)), config)
				require.NoError(t, err)
				changes := wc.CompareDocuments(left, right)
				require.NotNil(t, changes)
				require.Equal(t, 1, changes.TotalChanges())
				require.Equal(t, "$ref", changes.GetAllChanges()[0].Property)
				require.Nil(t, wc.CompareDocuments(left, left))
			}
		})
	}
	t.Run("trait", func(t *testing.T) {
		spec := func(target string) string {
			return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n" + message("        traits:\n          - headers: {$ref: '#/components/schemas/"+target+"'}\n") + fmt.Sprintf(components, "S")
		}
		left, err := libasyncapi.NewDocument([]byte(spec("S")))
		require.NoError(t, err)
		right, err := libasyncapi.NewDocument([]byte(spec("S2")))
		require.NoError(t, err)
		changes := wc.CompareDocuments(left, right)
		require.NotNil(t, changes)
		require.Equal(t, "traits", changes.GetAllChanges()[0].Property)
	})
	t.Run("component alias", func(t *testing.T) {
		spec := func(target string) string {
			return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n" + fmt.Sprintf(components, target)
		}
		left, err := libasyncapi.NewDocument([]byte(spec("S")))
		require.NoError(t, err)
		right, err := libasyncapi.NewDocument([]byte(spec("S2")))
		require.NoError(t, err)
		changes := wc.CompareDocuments(left, right)
		require.NotNil(t, changes)
		require.Equal(t, 1, changes.TotalChanges())
		require.Equal(t, "$ref", changes.GetAllChanges()[0].Property)
	})
	t.Run("inline to reference with same content", func(t *testing.T) {
		spec := func(payload string) string {
			return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n" + message("        payload: "+payload+"\n") + fmt.Sprintf(components, "S")
		}
		left, err := libasyncapi.NewDocument([]byte(spec("{type: string}")))
		require.NoError(t, err)
		right, err := libasyncapi.NewDocument([]byte(spec("{$ref: '#/components/schemas/S'}")))
		require.NoError(t, err)
		require.Nil(t, wc.CompareDocuments(left, right))
	})
}

func TestSkipCircularCheckReferenceRename(t *testing.T) {
	spec := func(target string) string {
		return "asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nchannels:\n  events:\n    address: events\n    messages:\n      event:\n        payload: {$ref: '#/components/schemas/" + target + "'}\ncomponents:\n  schemas:\n    A:\n      type: object\n      properties:\n        next: {$ref: '#/components/schemas/A'}\n    B:\n      type: object\n      properties:\n        next: {$ref: '#/components/schemas/B'}\n"
	}
	config := libasyncapi.NewDocumentConfiguration()
	config.SkipCircularReferenceCheck = true
	left, err := libasyncapi.NewDocumentWithConfiguration([]byte(spec("A")), config)
	require.NoError(t, err)
	right, err := libasyncapi.NewDocumentWithConfiguration([]byte(spec("B")), config)
	require.NoError(t, err)
	changes := wc.CompareDocuments(left, right)
	require.NotNil(t, changes)
	require.Equal(t, 1, changes.TotalChanges())
	require.Equal(t, "$ref", changes.GetAllChanges()[0].Property)
}

func TestRemoteReferencesStayUnderBaseURL(t *testing.T) {
	const avro, jsonSchema = "application/vnd.apache.avro;version=1.9.0", "application/schema+json;version=draft-07"
	for name, test := range map[string]struct {
		spec string
	}{
		"avro":    {strings.Replace(wrappedSchemaSpec("payload", avro, "string"), "type: string", "$ref: schemas/s.yaml", 1)},
		"json":    {strings.Replace(wrappedSchemaSpec("payload", jsonSchema, "string"), "type: string", "$ref: schemas/s.yaml", 1)},
		"alias":   {strings.Replace(wrappedSchemaSpec("payload", avro, "string"), "type: string", "$ref: 'schemas/lib.yaml#/alias/schema'", 1)},
		"chained": {strings.Replace(wrappedSchemaSpec("payload", avro, "string"), "type: string", "$ref: schemas/chain.yaml", 1)},
	} {
		for _, root := range []string{"", "file", "dir", "both"} {
			t.Run(name+"/"+root, func(t *testing.T) {
				var docs [2]libasyncapi.Document
				var mu sync.Mutex
				var paths []string
				for i, kind := range []string{"string", "integer"} {
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						mu.Lock()
						paths = append(paths, r.URL.Path)
						mu.Unlock()
						switch r.URL.Path {
						case "/specs/schemas/s.yaml":
							_, _ = w.Write([]byte("type: " + kind + "\n"))
						case "/specs/schemas/chain.yaml":
							_, _ = w.Write([]byte("$ref: s.yaml\n"))
						case "/specs/schemas/lib.yaml":
							_, _ = w.Write([]byte("alias: {$ref: 'other.yaml#/target'}\n"))
						case "/specs/schemas/other.yaml":
							_, _ = w.Write([]byte("target:\n  schema: {type: " + kind + "}\n"))
						default:
							http.NotFound(w, r)
						}
					}))
					defer server.Close()
					base, err := url.Parse(server.URL + "/specs/")
					require.NoError(t, err)
					config := libasyncapi.NewDocumentConfiguration()
					config.BaseURL = base
					config.AllowRemoteReferences = true
					switch root {
					case "file":
						config.SpecFilePath = filepath.Join(t.TempDir(), "root.yaml")
					case "dir":
						config.BasePath = t.TempDir()
					case "both":
						config.BasePath = t.TempDir()
						config.SpecFilePath = filepath.Join(config.BasePath, "root.yaml")
					}
					docs[i], err = libasyncapi.NewDocumentWithConfiguration([]byte(test.spec), config)
					require.NoError(t, err)
					require.False(t, docs[i].IsPartial(), "%v", docs[i].Errors())
				}
				changes := wc.CompareDocuments(docs[0], docs[1])
				mu.Lock()
				defer mu.Unlock()
				for _, path := range paths {
					require.True(t, strings.HasPrefix(path, "/specs/"), "requested %s", path)
				}
				require.NotNil(t, changes)
				require.Equal(t, 1, changes.TotalChanges())
			})
		}
	}
}
