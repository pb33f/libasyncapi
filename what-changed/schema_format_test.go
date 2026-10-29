package what_changed_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
