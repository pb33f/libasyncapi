// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pb33f/testify/require"
)

func TestDocumentAliasExpansionLimit(t *testing.T) {
	var spec strings.Builder
	spec.WriteString("asyncapi: 3.0.0\ninfo: {title: Events, version: 1.0.0}\nx-data:\n  a0: &a0 {value: data}\n")
	for i := 1; i <= 4; i++ {
		aliases := strings.TrimSuffix(strings.Repeat(fmt.Sprintf("*a%d, ", i-1), 8), ", ")
		fmt.Fprintf(&spec, "  a%d: &a%d [%s]\n", i, i, aliases)
	}
	for _, allow := range []bool{false, true} {
		t.Run(fmt.Sprint(allow), func(t *testing.T) {
			config := NewDocumentConfiguration()
			config.AllowExcessiveAliases = allow
			doc, err := NewDocumentWithConfiguration([]byte(spec.String()), config)
			require.NoError(t, err)
			defer doc.Rolodex().Release()
			require.Equal(t, allow, doc.Index().GetConfig().AllowExcessiveAliases)
			require.Equal(t, allow, config.ToLibOpenAPIConfig().AllowExcessiveAliases)
			if allow {
				require.False(t, doc.IsPartial(), "%v", doc.Errors())
			} else {
				require.True(t, doc.IsPartial())
				require.Contains(t, fmt.Sprint(doc.Errors()), "excessive aliasing")
			}
		})
	}
}
