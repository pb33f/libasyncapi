package what_changed_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pb33f/libasyncapi"
	wc "github.com/pb33f/libasyncapi/what-changed"
	"github.com/pb33f/testify/require"
)

// a repaired alias is registered by its raw local reference, so the indexed
// location of a base path qualified reference does not find it.
func TestWrappedSchemaReferenceThroughRepairedAlias(t *testing.T) {
	spec := `asyncapi: 3.0.0
info:
  title: Alias
  version: 1.0.0
channels:
  orders:
    address: orders
    messages:
      created:
        payload:
          schemaFormat: application/vnd.apache.avro;version=1.9.0
          schema:
            $ref: '#/components/schemas/Alias/properties/x'
components:
  schemas:
    Alias:
      $ref: '#/components/schemas/Target'
    Target:
      type: object
      properties:
        x:
          type: string
`
	load := func(name, content string) libasyncapi.Document {
		dir := filepath.Join(t.TempDir(), name)
		require.NoError(t, os.MkdirAll(dir, 0o755))
		path := filepath.Join(dir, "asyncapi.yaml")
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
		doc, err := libasyncapi.NewDocumentWithConfiguration([]byte(content), &libasyncapi.DocumentConfiguration{
			BasePath:            dir,
			SpecFilePath:        path,
			AllowFileReferences: true,
		})
		require.NoError(t, err)
		require.False(t, doc.IsPartial(), "%v", doc.Errors())
		return doc
	}

	left := load("left", spec)
	right := load("right", strings.Replace(spec, "type: string", "type: integer", 1))
	changes := wc.CompareDocuments(left, right)
	require.NotNil(t, changes)
	wrapperChanged := false
	for _, change := range changes.GetAllChanges() {
		wrapperChanged = wrapperChanged || change.Property == "schema"
	}
	require.True(t, wrapperChanged, "the wrapper's own change was lost")
}
