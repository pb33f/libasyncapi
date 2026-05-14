# AGENTS.md

Operational guide for AI agents working in `github.com/pb33f/libasyncapi`.

## TL;DR

- Purpose: AsyncAPI 3.0+ parser for Go with high-level and low-level APIs.
- Entry point: `libasyncapi.NewDocument([]byte) (Document, error)`.
- High-level model: `doc.Model() *asyncapi.AsyncAPI`.
- Low-level model: `doc.GoLow() *lowasync.AsyncAPI`.
- AsyncAPI 2.x is intentionally unsupported and returns an error.
- This file is the single source of truth. `CLAUDE.md` should only point here.

## Verify Changes

Run the smallest useful set first, then broaden when the change touches shared
behavior.

```bash
go test ./...
go vet ./...
go build ./...
```

For visitor changes, also run:

```bash
go test -race ./visitor/...
```

For formatting after Go edits:

```bash
gofmt -w <changed-go-files>
```

## Repo Map

```text
asyncapi.go                       package docs and NewDocument entry point
document.go                       Document interface
document_config.go                parser configuration
version.go                        version detection and validation
errors.go                         public error types
datamodel/high/asyncapi/          high-level user-facing models
datamodel/low/asyncapi/           low-level YAML-aware models
datamodel/low/asyncapi/constants.go field labels
datamodel/low/asyncapi/create_document.go document builder and index setup
visitor/                          visitor interfaces, walker, context helpers
test_fixtures/                    integration and binding fixtures
```

High-level and low-level files mostly mirror each other by object name:
`server.go`, `channel.go`, `operation.go`, `message.go`, `components.go`,
`security_scheme.go`, and `bindings.go`.

## Core Model Rules

- Prefer high-level models for normal application behavior.
- Use low-level models when line numbers, raw YAML nodes, source structure, or
  render fidelity matter.
- Every high-level type should expose `GoLow()` and `GoLowUntyped()`.
- Low-level fields are wrapped in `low.NodeReference[T]`; use `IsEmpty()` to
  test whether a field was present.
- Preserve YAML order with `orderedmap.Map[K, V]`. Do not replace parsed maps
  with plain Go maps unless the source model already uses one intentionally.
- Optional scalar fields may be pointers at the high level. Check for nil before
  dereferencing, for example `channel.Address != nil`.
- Schemas come from `libopenapi` (`highbase.SchemaProxy` and
  `highbase.Schema`). Always handle `proxy.Schema() == nil`.

## Common Usage

Parse a document:

```go
spec, err := os.ReadFile("asyncapi.yaml")
if err != nil {
    return err
}

doc, err := libasyncapi.NewDocument(spec)
if err != nil {
    return err
}

model := doc.Model()
fmt.Println(model.Info.Title)
```

Handle partial parses:

```go
doc, err := libasyncapi.NewDocument(spec)
if err != nil {
    return err
}

if doc.IsPartial() {
    for _, parseErr := range doc.Errors() {
        log.Printf("parse warning: %v", parseErr)
    }
}
```

Iterate ordered maps:

```go
for name, channel := range model.Channels.FromOldest() {
    if channel.Address != nil {
        fmt.Printf("%s: %s\n", name, *channel.Address)
    }
}
```

Access source metadata:

```go
lowModel := doc.GoLow()
fmt.Println(lowModel.Info.Value.Title.KeyNode.Line)
```

Configure multi-file references:

```go
config := libasyncapi.NewDocumentConfiguration()
config.BasePath = "/path/to/specs"
config.AllowFileReferences = true
config.AllowRemoteReferences = false

doc, err := libasyncapi.NewDocumentWithConfiguration(spec, config)
```

## Document Access

| Need | API |
|------|-----|
| High-level model | `doc.Model()` |
| Low-level model | `doc.GoLow()` |
| Version | `doc.GetVersion()` |
| Spec info | `doc.GetSpecInfo()` |
| Reference index | `doc.Index()` |
| Multi-file manager | `doc.Rolodex()` |
| Root YAML node | `doc.RootNode()` |
| Parse errors | `doc.Errors()` |
| Partial parse state | `doc.IsPartial()` |

## Low-Level Patterns

Most low-level types follow this shape:

- `Build(ctx, keyNode, root, idx)` parses from YAML.
- `Hash()` returns a stable content hash.
- `GetRootNode()` and `GetKeyNode()` expose YAML nodes.
- `GetExtensions()` exposes `x-*` extensions.
- `GetIndex()` and `GetContext()` preserve parser context.

Binding low-level types in `datamodel/low/asyncapi/bindings.go` embed
`BaseBinding`. Use `initBuild()` for shared metadata/index/context setup and
`hashExtensions()` at the end of binding `Hash()` implementations.

Use structured extraction helpers instead of ad hoc YAML traversal:

```go
title := low.ExtractValueNode[string](ctx, "title", root)
info, err := low.ExtractObject[*Info](ctx, "info", root, idx)
servers, keyNode, valueNode, err := low.ExtractMap[*Server](ctx, "servers", root, idx)
tags, keyNode, valueNode, err := low.ExtractArray[*Tag](ctx, "tags", root, idx)
```

## Rendering Rules

- High-level `MarshalYAML()` should use `high.NewNodeBuilder(value, value.low)`
  when low-level metadata must be preserved.
- Do not patch low-level structs to "fix" high-level rendering unless parsing
  metadata is actually wrong.
- For binding render fidelity, add `MarshalYAML()` on high-level bindings and
  nested helper structs that must preserve extensions, key order, or raw nodes.

Example:

```go
func (s *SQSChannelBinding) MarshalYAML() (interface{}, error) {
    nb := high.NewNodeBuilder(s, s.low)
    return nb.Render(), nil
}
```

## Visitor Rules

Basic visitor:

```go
type Recorder struct {
    Paths []string
}

func (r *Recorder) Visit(ctx context.Context, node any) error {
    r.Paths = append(r.Paths, visitor.Path(ctx))
    return nil
}

rec := &Recorder{}
walker := visitor.NewWalker(rec)
err := walker.Walk(context.Background(), doc.Model())
```

Context helpers:

```go
visitor.Path(ctx)
visitor.Depth(ctx)
visitor.Stack(ctx)
visitor.Parent(ctx)
visitor.AppendPath(ctx, "segment")
visitor.AppendIndex(ctx, 5)
```

Schema visitors can implement `EnterSchema`, `LeaveSchema`, and
`SkipCircularRef`. Polymorphic visitors can implement `EnterAllOf`,
`LeaveAllOf`, `EnterOneOf`, `LeaveOneOf`, `EnterAnyOf`, and `LeaveAnyOf`.

When adding a new traversable type, update `visitor/walker.go` and add coverage
in `visitor/examples_test.go` or a focused walker test.

## Bindings

Binding availability by protocol:

| Protocol | Server | Channel | Operation | Message |
|----------|--------|---------|-----------|---------|
| HTTP | `HTTPServerBinding` | `HTTPChannelBinding` | `HTTPOperationBinding` | `HTTPMessageBinding` |
| Kafka | `KafkaServerBinding` | `KafkaChannelBinding` | `KafkaOperationBinding` | `KafkaMessageBinding` |
| WebSocket | - | `WebSocketChannelBinding` | - | - |
| AMQP | - | `AMQPChannelBinding` | `AMQPOperationBinding` | `AMQPMessageBinding` |
| MQTT | `MQTTServerBinding` | - | `MQTTOperationBinding` | `MQTTMessageBinding` |
| SQS | `SQSServerBinding` | `SQSChannelBinding` | `SQSOperationBinding` | `SQSMessageBinding` |

Access pattern:

```go
if server.Bindings != nil && server.Bindings.Kafka != nil {
    kafka := server.Bindings.Kafka
    fmt.Println(kafka.SchemaRegistryURL)
}
```

### SQS

SQS has typed bindings and queue helper types:

| Type | Key fields |
|------|------------|
| `SQSServerBinding` | Extensions only |
| `SQSChannelBinding` | `Queue`, `DeadLetterQueue`, `BindingVersion` |
| `SQSOperationBinding` | `Queues`, `BindingVersion` |
| `SQSMessageBinding` | Extensions only |
| `SQSQueue` | `Name`, `ARN`, `FifoQueue`, `DeduplicationScope`, `FifoThroughputLimit`, `DeliveryDelay`, `VisibilityTimeout`, `ReceiveMessageWaitTime`, `MessageRetentionPeriod`, `RedrivePolicy`, `Policy`, `Tags` |
| `SQSRedrivePolicy` | `DeadLetterQueue`, `MaxReceiveCount` |
| `SQSPolicy` | `Statements` |
| `SQSPolicyStatement` | `Effect`, `Principal`, `Action`, `Resource`, `Condition` |

```go
if channel.Bindings != nil && channel.Bindings.SQS != nil {
    sqs := channel.Bindings.SQS
    if sqs.Queue != nil {
        fmt.Println(sqs.Queue.Name)
        fmt.Println(sqs.Queue.VisibilityTimeout)
    }
}

if op.Bindings != nil && op.Bindings.SQS != nil {
    for _, queue := range op.Bindings.SQS.Queues {
        fmt.Println(queue.Name)
    }
}
```

SQS policy `Principal`, `Action`, `Resource`, and `Condition` fields are raw
`*yaml.Node` values because the binding allows scalar, sequence, or mapping
forms. Inspect the node shape or marshal the node when typed access is needed.

The walker visits SQS bindings at paths like `/channels/<name>/bindings/sqs`
and `/operations/<name>/bindings/sqs`. Queue, redrive policy, and policy helper
objects are data on the binding and are not walked as separate visitor nodes.

## Adding New AsyncAPI Objects

Use nearby objects as the primary template. The shortest safe checklist is:

1. Add labels in `datamodel/low/asyncapi/constants.go`.
2. Add the low-level type with `Build`, `Hash`, metadata accessors, and
   extension handling.
3. Add the high-level type with a `NewType(low *lowasync.Type)` constructor,
   `GoLow()`, `GoLowUntyped()`, `Render()`, and `MarshalYAML()` when render
   fidelity matters.
4. Wire the type into its parent low-level `Build()` and high-level
   constructor.
5. Update `visitor/walker.go` if the type should be traversed.
6. Add tests for high-level fields, low-level metadata, hashes where relevant,
   rendering, and visitor paths when applicable.

## Adding New Binding Types

Bindings have extra wiring:

1. Add protocol and field labels in `constants.go`.
2. Add low-level binding structs in `datamodel/low/asyncapi/bindings.go`.
3. Add high-level binding structs and constructors in
   `datamodel/high/asyncapi/bindings.go`.
4. Add fields to the relevant containers: `ServerBindings`, `ChannelBindings`,
   `OperationBindings`, and `MessageBindings`.
5. Update `visitor/walker.go` so the binding is visited under
   `/.../bindings/<protocol>`.
6. Add or extend `bindings_test.go`.
7. Add representative examples to `test_fixtures/comprehensive-bindings.yaml`.

For bindings with nested helper types or extensions, include render round-trip
coverage. SQS bindings are the reference pattern for that round-trip wiring.

## Fixtures

- `test_fixtures/streetlights-kafka.yaml`: full Kafka example with bindings.
- `test_fixtures/multi-protocol.yaml`: HTTP, Kafka, and WebSocket.
- `test_fixtures/comprehensive-bindings.yaml`: binding coverage, including SQS
  queue/redrive/policy cases.
- `test_fixtures/shared-message.yaml`: component references.

## Common Mistakes

- Direct map iteration loses YAML order. Use `.FromOldest()` or `.FromNewest()`.
- Optional high-level pointer fields can be nil.
- Low-level fields need `IsEmpty()` checks.
- Schema resolution can fail; check `proxy.Schema() != nil`.
- Do not mutate low-level models as application state. Treat them as parsed
  source metadata.
- Do not bypass tests or pre-commit hooks. Fix the underlying issue.

## Agent Notes

Rules for any AI agent working in this repo:

- Keep changes scoped to the request and existing code patterns.
- Preserve user edits already present in the worktree.
- Prefer focused tests first, then broader tests when behavior is shared.
- If verification cannot be run, say exactly why.

### Codex

- Use `apply_patch` for manual file edits.
- Do not run destructive git commands unless explicitly requested.
- Do not commit unless the user explicitly asks for a commit.
- When editing docs, keep text ASCII unless the file already requires otherwise.
- Report changed files and verification in the final response.

### Claude

- Claude should read this file through the repository `CLAUDE.md` redirect.
- Claude in the pb33f org is not permitted to run `git commit`; stage changes
  and leave committing to the user.
