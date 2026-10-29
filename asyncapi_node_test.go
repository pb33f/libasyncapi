package libasyncapi

import (
	"io/fs"
	"testing"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/testify/require"
)

func TestNewDocumentFromNode(t *testing.T) {
	var root yaml.Node
	require.NoError(t, yaml.Unmarshal([]byte("asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nchannels: {}\n"), &root))
	doc, err := NewDocumentFromNode(&root, nil)
	require.NoError(t, err)
	require.False(t, doc.IsPartial())
	require.Same(t, &root, doc.RootNode())
	require.Equal(t, "Events", doc.Model().Info.Title)
	_, err = NewDocumentFromNode(nil, nil)
	require.ErrorIs(t, err, ErrInvalidYAML)
}

type revisionTestFS struct{ opened bool }

func (f *revisionTestFS) Open(name string) (fs.File, error) {
	f.opened = true
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
func (f *revisionTestFS) GetFiles() map[string]index.RolodexFile { return nil }

func TestDocumentUsesProvidedRolodexFS(t *testing.T) {
	files := new(revisionTestFS)
	config := NewDocumentConfiguration()
	config.LocalFS = files
	config.AllowFileReferences = true
	config.BasePath = t.TempDir()
	doc, err := NewDocumentWithConfiguration([]byte("asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\n"), config)
	require.NoError(t, err)
	require.False(t, doc.IsPartial())
	require.False(t, files.opened, "a revision filesystem must not be directory-walked")
}
