// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package visitor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/pb33f/libasyncapi"
	"github.com/pb33f/libasyncapi/visitor"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// PathRecorder is a simple visitor that records all visited paths.
type PathRecorder struct {
	Paths []string
}

func (r *PathRecorder) Visit(ctx context.Context, node any) error {
	r.Paths = append(r.Paths, visitor.Path(ctx))
	return nil
}

// SchemaCollector implements SchemaVisitor to track schema visits and circular refs.
type SchemaCollector struct {
	PathRecorder
	References []string
	Circular   []string
	Entered    []string
	Left       []string
}

func (c *SchemaCollector) EnterSchema(ctx context.Context, proxy *highbase.SchemaProxy) error {
	path := visitor.Path(ctx)
	c.Entered = append(c.Entered, path)
	if proxy.IsReference() {
		c.References = append(c.References, proxy.GetReference())
	}
	return nil
}

func (c *SchemaCollector) LeaveSchema(ctx context.Context, proxy *highbase.SchemaProxy, schema *highbase.Schema, err error) {
	c.Left = append(c.Left, visitor.Path(ctx))
}

func (c *SchemaCollector) SkipCircularRef(ctx context.Context, proxy *highbase.SchemaProxy, ref string) error {
	c.Circular = append(c.Circular, ref)
	return nil
}

// PolymorphicRecorder extends SchemaCollector with polymorphic tracking.
type PolymorphicRecorder struct {
	SchemaCollector
	AllOfEntered []string
	AllOfLeft    []string
	OneOfEntered []string
	OneOfLeft    []string
	AnyOfEntered []string
	AnyOfLeft    []string
}

func (p *PolymorphicRecorder) EnterAllOf(ctx context.Context, schema *highbase.Schema, count int) error {
	p.AllOfEntered = append(p.AllOfEntered, visitor.Path(ctx))
	return nil
}

func (p *PolymorphicRecorder) LeaveAllOf(ctx context.Context, schema *highbase.Schema, err error) {
	p.AllOfLeft = append(p.AllOfLeft, visitor.Path(ctx))
}

func (p *PolymorphicRecorder) EnterOneOf(ctx context.Context, schema *highbase.Schema, count int) error {
	p.OneOfEntered = append(p.OneOfEntered, visitor.Path(ctx))
	return nil
}

func (p *PolymorphicRecorder) LeaveOneOf(ctx context.Context, schema *highbase.Schema, err error) {
	p.OneOfLeft = append(p.OneOfLeft, visitor.Path(ctx))
}

func (p *PolymorphicRecorder) EnterAnyOf(ctx context.Context, schema *highbase.Schema, count int) error {
	p.AnyOfEntered = append(p.AnyOfEntered, visitor.Path(ctx))
	return nil
}

func (p *PolymorphicRecorder) LeaveAnyOf(ctx context.Context, schema *highbase.Schema, err error) {
	p.AnyOfLeft = append(p.AnyOfLeft, visitor.Path(ctx))
}

// ErroringVisitor returns an error at a specific path.
type ErroringVisitor struct {
	SchemaCollector
	ErrorOnPath    string
	LeaveWasCalled bool
}

func (e *ErroringVisitor) Visit(ctx context.Context, node any) error {
	path := visitor.Path(ctx)
	e.Paths = append(e.Paths, path)
	if path == e.ErrorOnPath {
		return errors.New("intentional error")
	}
	return nil
}

func (e *ErroringVisitor) LeaveSchema(ctx context.Context, proxy *highbase.SchemaProxy, schema *highbase.Schema, err error) {
	e.LeaveWasCalled = true
	e.Left = append(e.Left, visitor.Path(ctx))
}

// TestPathRootIsEmptyString verifies root path is empty string per RFC 6901.
func TestPathRootIsEmptyString(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Root node should have empty path
	require.NotEmpty(t, recorder.Paths)
	assert.Equal(t, "", recorder.Paths[0], "Root path should be empty string")

	// Info should be /info, not //info
	assert.Contains(t, recorder.Paths, "/info")
	for _, p := range recorder.Paths {
		assert.NotContains(t, p, "//", "Path should not contain double slashes: %s", p)
	}
}

// TestStackBasedCycleDetection verifies same ref in different branches is visited twice.
func TestStackBasedCycleDetection(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
components:
  schemas:
    Shared:
      type: string
    Parent:
      type: object
      properties:
        child1:
          $ref: '#/components/schemas/Shared'
        child2:
          $ref: '#/components/schemas/Shared'
`))
	require.NoError(t, err)

	collector := &SchemaCollector{}
	walker := visitor.NewWalker(collector)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Shared should be visited TWICE (once per branch), not skipped
	sharedCount := 0
	for _, ref := range collector.References {
		if ref == "#/components/schemas/Shared" {
			sharedCount++
		}
	}
	assert.Equal(t, 2, sharedCount, "Shared ref should be visited twice (stack-based, not global visited)")
	assert.Empty(t, collector.Circular, "No circular refs in this spec")
}

// TestTrueCycleDetection verifies self-referencing schemas are detected as cycles.
func TestTrueCycleDetection(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
components:
  schemas:
    Node:
      type: object
      properties:
        children:
          type: array
          items:
            $ref: '#/components/schemas/Node'
`))
	require.NoError(t, err)

	collector := &SchemaCollector{}
	walker := visitor.NewWalker(collector)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Should detect cycle when Node references itself while still on stack
	assert.Contains(t, collector.Circular, "#/components/schemas/Node",
		"Self-referencing schema should be detected as circular")
}

// TestLeaveSchemaOnError verifies LeaveSchema fires even when child walk fails.
func TestLeaveSchemaOnError(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
components:
  schemas:
    Test:
      type: object
      properties:
        name:
          type: string
`))
	require.NoError(t, err)

	ev := &ErroringVisitor{
		// Error during child traversal (on a property schema), not on the parent schema.
		// This tests that LeaveSchema fires even when walking children fails.
		ErrorOnPath: "/components/schemas/Test/properties/name",
	}
	walker := visitor.NewWalker(ev)
	_ = walker.Walk(context.Background(), doc.Model())

	// LeaveSchema should have been called for Test schema despite error in child
	assert.True(t, ev.LeaveWasCalled, "LeaveSchema must fire even on error")
}

// TestEmptySegmentIgnored verifies empty segments don't create trailing slashes.
func TestEmptySegmentIgnored(t *testing.T) {
	ctx := context.Background()
	ctx = visitor.WithPath(ctx, "/channels")

	// Empty segment should be a no-op
	ctx = visitor.AppendPath(ctx, "")
	assert.Equal(t, "/channels", visitor.Path(ctx), "Empty segment should not add trailing slash")

	// Non-empty still works
	ctx = visitor.AppendPath(ctx, "users")
	assert.Equal(t, "/channels/users", visitor.Path(ctx))
}

// TestJSONPointerEscaping verifies / and ~ are properly escaped.
func TestJSONPointerEscaping(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  user/signup:
    address: user/signup
    messages:
      msg~1:
        payload:
          type: string
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// / escaped as ~1, ~ escaped as ~0
	assert.Contains(t, recorder.Paths, "/channels/user~1signup",
		"Forward slash in key should be escaped as ~1")

	// Check for message key with ~
	found := false
	for _, p := range recorder.Paths {
		if p == "/channels/user~1signup/messages/msg~01" {
			found = true
			break
		}
	}
	assert.True(t, found, "Tilde in key should be escaped as ~0")
}

// TestPolymorphicVisitor verifies allOf/oneOf/anyOf callbacks work.
func TestPolymorphicVisitor(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
components:
  schemas:
    Combined:
      allOf:
        - type: object
          properties:
            name:
              type: string
        - type: object
          properties:
            age:
              type: integer
`))
	require.NoError(t, err)

	recorder := &PolymorphicRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// allOf should have been entered and left
	assert.NotEmpty(t, recorder.AllOfEntered, "EnterAllOf should have been called")
	assert.NotEmpty(t, recorder.AllOfLeft, "LeaveAllOf should have been called")
	assert.Equal(t, len(recorder.AllOfEntered), len(recorder.AllOfLeft),
		"Enter and Leave counts should match")
}

// TestDepthTracking verifies depth increases as we descend through typed nodes.
func TestDepthTracking(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  events:
    address: events
    messages:
      UserCreated:
        payload:
          type: object
`))
	require.NoError(t, err)

	depths := make(map[string]int)
	walker := visitor.NewWalker(&depthRecorder{depths: depths})
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Root should be depth 0
	assert.Equal(t, 0, depths[""], "Root should be at depth 0")

	// Info should be depth 1 (direct child of doc)
	assert.Equal(t, 1, depths["/info"], "Info should be at depth 1")

	// Channel entries are visited at depth 1 (same as other doc children)
	assert.Equal(t, 1, depths["/channels/events"], "Channel should be at depth 1")

	// Message entries are visited at depth 2 (child of channel)
	assert.Equal(t, 2, depths["/channels/events/messages/UserCreated"],
		"Message should be at depth 2")
}

type depthRecorder struct {
	depths map[string]int
}

func (d *depthRecorder) Visit(ctx context.Context, node any) error {
	d.depths[visitor.Path(ctx)] = visitor.Depth(ctx)
	return nil
}

// TestSchemaPropertiesWalked verifies schema properties are walked correctly.
func TestSchemaPropertiesWalked(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
components:
  schemas:
    User:
      type: object
      properties:
        name:
          type: string
        email:
          type: string
`))
	require.NoError(t, err)

	collector := &SchemaCollector{}
	walker := visitor.NewWalker(collector)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Should have entered schemas for User, name, and email
	assert.GreaterOrEqual(t, len(collector.Entered), 3,
		"Should enter at least 3 schemas (User + 2 properties)")
}

// TestMessagePayloadWalked verifies message payload schemas are walked.
func TestMessagePayloadWalked(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  events:
    address: events
    messages:
      UserCreated:
        payload:
          type: object
          properties:
            userId:
              type: string
`))
	require.NoError(t, err)

	collector := &SchemaCollector{}
	walker := visitor.NewWalker(collector)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Should have entered the payload schema and its property
	assert.GreaterOrEqual(t, len(collector.Entered), 2,
		"Should enter payload schema and userId property")
}

// TestChannelServersReferencesWalked verifies channel server refs are visited.
func TestChannelServersReferencesWalked(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
servers:
  production:
    host: api.example.com
    protocol: https
  staging:
    host: staging.example.com
    protocol: https
channels:
  events:
    address: events
    servers:
      - $ref: '#/servers/production'
      - $ref: '#/servers/staging'
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Channel server references should be visited
	assert.Contains(t, recorder.Paths, "/channels/events/servers/0",
		"First server reference should be visited")
	assert.Contains(t, recorder.Paths, "/channels/events/servers/1",
		"Second server reference should be visited")
}

// TestOperationChannelAndMessagesReferencesWalked verifies operation refs are visited.
func TestOperationChannelAndMessagesReferencesWalked(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  events:
    address: events
    messages:
      UserCreated:
        payload:
          type: string
operations:
  publishEvent:
    action: send
    channel:
      $ref: '#/channels/events'
    messages:
      - $ref: '#/channels/events/messages/UserCreated'
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Operation channel reference should be visited
	assert.Contains(t, recorder.Paths, "/operations/publishEvent/channel",
		"Operation channel reference should be visited")
	// Operation messages references should be visited
	assert.Contains(t, recorder.Paths, "/operations/publishEvent/messages/0",
		"Operation message reference should be visited")
}

// TestOperationReplyReferencesWalked verifies operation reply refs are visited.
func TestOperationReplyReferencesWalked(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  requests:
    address: requests
    messages:
      Request:
        payload:
          type: string
  replies:
    address: replies
    messages:
      Response:
        payload:
          type: string
operations:
  sendRequest:
    action: send
    channel:
      $ref: '#/channels/requests'
    reply:
      channel:
        $ref: '#/channels/replies'
      messages:
        - $ref: '#/channels/replies/messages/Response'
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	err = walker.Walk(context.Background(), doc.Model())
	require.NoError(t, err)

	// Operation reply channel reference should be visited
	assert.Contains(t, recorder.Paths, "/operations/sendRequest/reply/channel",
		"Reply channel reference should be visited")
	// Operation reply messages references should be visited
	assert.Contains(t, recorder.Paths, "/operations/sendRequest/reply/messages/0",
		"Reply message reference should be visited")
}

// StopOnPathVisitor stops traversal when it reaches a specific path.
type StopOnPathVisitor struct {
	PathRecorder
	StopOnPath string
	Stopped    bool
}

func (s *StopOnPathVisitor) Visit(ctx context.Context, node any) error {
	path := visitor.Path(ctx)
	s.Paths = append(s.Paths, path)
	if path == s.StopOnPath {
		s.Stopped = true
		return visitor.ErrStopTraversal
	}
	return nil
}

// TestErrStopTraversal verifies clean early termination of traversal.
func TestErrStopTraversal(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: Test
  version: "1.0.0"
channels:
  events:
    address: events
  users:
    address: users
`))
	require.NoError(t, err)

	// Stop when we reach the events channel
	stopper := &StopOnPathVisitor{StopOnPath: "/channels/events"}
	walker := visitor.NewWalker(stopper)
	err = walker.Walk(context.Background(), doc.Model())

	// Should return nil (not an error) despite early stop
	require.NoError(t, err, "ErrStopTraversal should result in nil error")
	assert.True(t, stopper.Stopped, "Visitor should have stopped")

	// Should have visited up to and including /channels/events
	assert.Contains(t, stopper.Paths, "/channels/events")

	// Should NOT have visited /channels/users (comes after events in insertion order)
	assert.NotContains(t, stopper.Paths, "/channels/users",
		"Traversal should have stopped before reaching /channels/users")
}

func TestWalkerVisitsSQSBindings(t *testing.T) {
	doc, err := libasyncapi.NewDocument([]byte(`asyncapi: "3.0.0"
info:
  title: SQS Test
  version: "1.0.0"
servers:
  sqsServer:
    host: sqs.us-east-1.amazonaws.com
    protocol: sqs
    bindings:
      sqs: {}
channels:
  workerQueue:
    address: worker
    bindings:
      sqs:
        queue:
          name: worker
          fifoQueue: false
        bindingVersion: 0.3.0
operations:
  receiveWorkerQueue:
    action: receive
    channel:
      $ref: '#/channels/workerQueue'
    bindings:
      sqs:
        queues:
          - name: worker
            fifoQueue: false
        bindingVersion: 0.3.0
components:
  messages:
    WorkerMessage:
      payload:
        type: object
      bindings:
        sqs: {}
`))
	require.NoError(t, err)

	recorder := &PathRecorder{}
	walker := visitor.NewWalker(recorder)
	require.NoError(t, walker.Walk(context.Background(), doc.Model()))

	assert.Contains(t, recorder.Paths, "/servers/sqsServer/bindings/sqs")
	assert.Contains(t, recorder.Paths, "/channels/workerQueue/bindings/sqs")
	assert.Contains(t, recorder.Paths, "/operations/receiveWorkerQueue/bindings/sqs")
	assert.Contains(t, recorder.Paths, "/components/messages/WorkerMessage/bindings/sqs")
}
