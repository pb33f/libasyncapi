package yamlparse

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/testify/assert"
	"github.com/pb33f/testify/require"
)

func TestJSONEscapesPreserveSourcePositions(t *testing.T) {
	for _, source := range []string{
		`{"title":"\ud83d\ude00","version":"1.0","path":"a\/b","after":true}`,
		"{\n  \"title\":\"\\ud83d\\ude00\",\"version\":\"1.0\",\n  \"after\":true\n}",
		`["\ud83d\ude00", {"version":"1.0"}]`,
		`{"literal":"é","title":"\ud83d\ude00","version":"1.0"}`,
		"\t{\"title\"\n: \"\\ud83d\\ude00\", \"version\":\"1\"}",
	} {
		root, err := Parse([]byte(source))
		require.NoError(t, err)
		var visit func(*yaml.Node)
		visit = func(n *yaml.Node) {
			if n.Value == "version" {
				line := strings.Split(source, "\n")[n.Line-1]
				want := utf8.RuneCountInString(line[:strings.Index(line, `"version"`)]) + 1
				assert.Equal(t, want, n.Column)
			}
			for _, c := range n.Content {
				visit(c)
			}
		}
		visit(root)
	}
}

func TestJSONStringValuesAndNumberTypes(t *testing.T) {
	for _, value := range []string{"one\u0085two", "one\u2028two", "one\u2029two"} {
		root, err := Parse([]byte(`{"value":"` + value + `","number":9007199254740993,"float":1e3,"bool":true,"null":null}`))
		require.NoError(t, err)
		nodes := root.Content[0].Content
		assert.Equal(t, value, nodes[1].Value)
		assert.Equal(t, "9007199254740993", nodes[3].Value)
		assert.Equal(t, "!!int", nodes[3].Tag)
		assert.Equal(t, "!!float", nodes[5].Tag)
		assert.Equal(t, "!!bool", nodes[7].Tag)
		assert.Equal(t, "!!null", nodes[9].Tag)
	}
}

func TestParseRejectsInvalidAndMultipleDocuments(t *testing.T) {
	for _, source := range []string{`{"title":"\ud83d\ude00",`, "---\na: b\n---\nc: d", `{"title":"\ud83d\ude00","bad":"\q"}`} {
		_, err := Parse([]byte(source))
		require.Error(t, err)
	}
}
