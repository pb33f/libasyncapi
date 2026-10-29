# libasyncapi - high-performance AsyncAPI tools for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/pb33f/libasyncapi.svg)](https://pkg.go.dev/github.com/pb33f/libasyncapi)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`libasyncapi` parses AsyncAPI 3.0+ documents into strongly typed Go models. It provides a
high-level API for normal application work and a low-level, YAML-aware API for tooling that
needs source nodes, line and column numbers, key ordering, references, or render fidelity.

AsyncAPI 2.x is intentionally unsupported.

---

## Features

- AsyncAPI 3.0+ parsing with high-level and low-level models.
- Lenient partial parsing with accumulated diagnostics.
- Local and remote multi-file reference resolution.
- Ordered maps and access to the original YAML node tree.
- Model rendering and original-document serialization.
- Visitor-based document traversal, including schemas and protocol bindings.
- Breaking and non-breaking change detection through `what-changed`.
- Typed HTTP, Kafka, WebSocket, AMQP, MQTT, and SQS bindings.
- Inline and composed AsyncAPI 3 bundling through `bundler`.

## Bundling

Use `github.com/pb33f/libasyncapi/bundler` to produce a portable YAML document:

```go
config := &libasyncapi.DocumentConfiguration{
    BasePath: "/path/to/specs",
    SpecFilePath: "asyncapi.yaml",
    AllowFileReferences: true,
    ExcludeExtensionRefs: true,
}
bundled, err := bundler.BundleBytes(spec, config)
```

`BundleBytesComposed(spec, config, &bundler.BundleCompositionConfig{Delimiter: "__"})`
moves imported objects into AsyncAPI component sections. Collisions add the source
file name and, when needed, a numeric suffix. Both modes preserve root-local
references and fields that require Reference Objects. Cycles use local references.
Operation messages continue to reference their channel's messages; operations and
channel servers continue to reference the required root locations. An operations-only
actor can reference channels in an external domain file. When no root destination
exists, the bundler imports the required external channels and servers at the root.
Existing root destinations take precedence; ambiguous targets still return an error.
Inline mode
emits each shared external schema once and uses local references for later uses
to prevent exponential output growth.

The bundler reads the source tree directly and does not require a typed model.
It preserves comments, key order, explicit empty values, examples, bindings, and
multi-format schema content. YAML aliases and merge keys are expanded. It does not
validate the complete AsyncAPI contract or merge traits. Input bytes and
configuration are not modified. An error returns no output, including when a
reference cannot be loaded.

Remote loading requires `AllowRemoteReferences`. `LocalFS` and `RemoteURLHandler`
can supply custom sources. `SpecFilePath` identifies the entry document for
references back to it. A nonempty `BasePath` enables local loading and supplies the
entry document's external reference directory, even when `SpecFilePath` is in a
different directory. Fragment-only references retain the entry document's identity;
loaded files resolve relative paths from their own locations. Set
`ExcludeExtensionRefs` to keep `$ref` values in extensions as data;
example values are always treated as data.

`BasePath` is a starting directory, not a file access boundary: parent paths and
absolute file references can read outside it, and their contents enter the bundle.
For untrusted input, supply a `LocalFS` that enforces the required boundary
(including symlink access). Referenced scalar documents are rejected except for
boolean schemas and content inside a Multi Format Schema Object.

Schema resource IDs and anchors are resolved before relocation and removed from
the output so local references use the bundled document's scope. OpenAPI-format
schemas also resolve their references. A Reference Object in the `schema` field
of a Multi Format Schema Object loads the referenced JSON or YAML value, including
Avro records, unions, and primitive types. The bundler preserves the schema content
without interpreting language-specific imports or nested data as JSON Schema.
Kafka `groupId` and `clientId`, and MQTT schema-valued fields, also resolve their
references. Circular YAML aliases, multiple YAML documents, and dynamic or recursive JSON Schema reference
keywords return errors. These features need separate semantics and are not silently
expanded.

The document parser resolves message pointers through channel Reference Objects,
including composed bundles, while retaining the original source nodes. Successful
alias repairs are removed from `Document.Errors()`; raw index diagnostics remain
available through libopenapi's index accessors. JSON input uses JSON string and
whitespace semantics and retains source line and column positions.

## Install

```bash
go get github.com/pb33f/libasyncapi
```

## Quick start

```go
package main

import (
	"fmt"
	"log"

	"github.com/pb33f/libasyncapi"
)

func main() {
	spec := []byte(`asyncapi: 3.0.0
info:
  title: Orders API
  version: 1.0.0
channels:
  orders:
    address: /orders
`)

	doc, err := libasyncapi.NewDocument(spec)
	if err != nil {
		log.Fatal(err)
	}

	model := doc.Model()
	fmt.Printf("%s uses AsyncAPI %s\n", model.Info.Title, doc.GetVersion())
	for name, channel := range model.Channels.FromOldest() {
		if channel.Address != nil {
			fmt.Printf("channel %s: %s\n", name, *channel.Address)
		}
	}
}
```

Output:

```text
Orders API uses AsyncAPI 3.0.0
channel orders: /orders
```

## Partial documents and diagnostics

Parsing continues where possible. A document can therefore contain a useful partial model
and one or more diagnostics:

```go
doc, err := libasyncapi.NewDocument(spec)
if err != nil {
	log.Fatal(err) // The document could not be created at all.
}

if doc.IsPartial() {
	for _, parseErr := range doc.Errors() {
		log.Printf("AsyncAPI diagnostic: %v", parseErr)
	}
}
```

## High-level and low-level models

Use `Model()` for ordinary application logic:

```go
model := doc.Model()
fmt.Println(model.Info.Title)
```

Use `GoLow()` when source information matters:

```go
lowModel := doc.GoLow()
fmt.Println(lowModel.Info.Value.Title.KeyNode.Line)
```

Every high-level model type exposes `GoLow()` and `GoLowUntyped()` for moving back to its
source-aware representation.

## Multi-file documents

External references are disabled by default. Enable only the reference types your
application needs:

```go
config := libasyncapi.NewDocumentConfiguration()
config.BasePath = "/path/to/asyncapi"
config.AllowFileReferences = true
config.AllowRemoteReferences = false

doc, err := libasyncapi.NewDocumentWithConfiguration(spec, config)
```

## Rendering and serialization

`Render()` writes the current high-level model. `Serialize()` writes the original YAML tree:

```go
rendered, err := doc.Render()
original, err := doc.Serialize()
```

## Compare documents

The `what-changed` package produces a hierarchical change model with breaking-change
classification:

```go
import what_changed "github.com/pb33f/libasyncapi/what-changed"

changes := what_changed.CompareDocuments(originalDoc, updatedDoc)
if changes != nil {
	fmt.Printf("%d changes, %d breaking\n",
		changes.TotalChanges(), changes.TotalBreakingChanges())
}
```

## Documentation and support

Multi Format Schemas preserve their `schemaFormat` and content during comparison.
JSON Schema wrappers use the existing schema comparison rules. Other schema
languages are compared as opaque content; format and content changes are breaking
by default and can be configured with the `multiFormatSchema` rule. A bare
`schema: {$ref: ...}` resolves its target, while reference-like fields inside
foreign schema content remain data. Formatting and source line changes do not
count as content changes.

Tools that already parse YAML or JSON can call `NewDocumentFromNode(root, config)`
to avoid parsing the same input again. This transfers tree ownership to the
document; do not mutate the tree while the document or its models are in use.

The public Go API is documented at
[pkg.go.dev/github.com/pb33f/libasyncapi](https://pkg.go.dev/github.com/pb33f/libasyncapi).
Full project documentation and a dedicated logo have not been published yet.

Need help, have a question, or want to share what you are building?
[Join the pb33f Discord](https://discord.gg/x7VACVuEGP).

## License

`libasyncapi` is released under the [MIT License](LICENSE).
