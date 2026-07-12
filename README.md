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

The public Go API is documented at
[pkg.go.dev/github.com/pb33f/libasyncapi](https://pkg.go.dev/github.com/pb33f/libasyncapi).
Full project documentation and a dedicated logo have not been published yet.

Need help, have a question, or want to share what you are building?
[Join the pb33f Discord](https://discord.gg/x7VACVuEGP).

## License

`libasyncapi` is released under the [MIT License](LICENSE).
