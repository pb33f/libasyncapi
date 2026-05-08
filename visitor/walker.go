// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package visitor

import (
	"context"
	"errors"
	"reflect"

	"github.com/pb33f/libasyncapi/datamodel/high/asyncapi"
	highbase "github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Walker traverses AsyncAPI documents in pre-order (visit node before children).
type Walker struct {
	visitor Visitor
}

// NewWalker creates a new Walker with the given visitor.
func NewWalker(v Visitor) *Walker {
	return &Walker{visitor: v}
}

// Walk traverses the document starting from the root.
// Returns an error if the visitor returns an error.
// ErrStopTraversal is converted to nil, allowing clean early termination.
func (w *Walker) Walk(ctx context.Context, doc *asyncapi.AsyncAPI) error {
	ctx = WithDepth(ctx, 0)
	ctx = WithPath(ctx, "") // Root is empty string per RFC 6901
	ctx = WithStack(ctx, make(map[schemaKey]bool))
	err := w.walkNode(ctx, doc)
	if errors.Is(err, ErrStopTraversal) {
		return nil // Clean stop, not an error
	}
	return err
}

// walkNode is the central dispatch that ensures pre-order visitation (Visit is
// called before descending into children) and routes to type-specific walkers.
// ErrStopTraversal is propagated up to stop siblings; Walk() converts it to nil.
func (w *Walker) walkNode(ctx context.Context, node any) error {
	if node == nil {
		return nil
	}

	// Pre-order: visit node before children
	if err := w.visitor.Visit(ctx, node); err != nil {
		return err // Propagate all errors including ErrStopTraversal
	}

	// Dispatch based on type
	switch n := node.(type) {
	case *asyncapi.AsyncAPI:
		return w.walkAsyncAPI(ctx, n)
	case *asyncapi.Info:
		return w.walkInfo(ctx, n)
	case *highbase.Contact:
		return nil // leaf node
	case *highbase.License:
		return nil // leaf node
	case *asyncapi.Server:
		return w.walkServer(ctx, n)
	case *asyncapi.ServerVariable:
		return nil // leaf node
	case *asyncapi.ServerBindings:
		return w.walkServerBindings(ctx, n)
	case *asyncapi.Channel:
		return w.walkChannel(ctx, n)
	case *asyncapi.ChannelBindings:
		return w.walkChannelBindings(ctx, n)
	case *asyncapi.Operation:
		return w.walkOperation(ctx, n)
	case *asyncapi.OperationBindings:
		return w.walkOperationBindings(ctx, n)
	case *asyncapi.OperationTrait:
		return w.walkOperationTrait(ctx, n)
	case *asyncapi.OperationReply:
		return w.walkOperationReply(ctx, n)
	case *asyncapi.OperationReplyAddress:
		return nil // leaf node
	case *asyncapi.Components:
		return w.walkComponents(ctx, n)
	case *asyncapi.Message:
		return w.walkMessage(ctx, n)
	case *asyncapi.MessageBindings:
		return w.walkMessageBindings(ctx, n)
	case *asyncapi.MessageTrait:
		return w.walkMessageTrait(ctx, n)
	case *asyncapi.MessageExample:
		return nil // leaf node
	case *asyncapi.CorrelationID:
		return nil // leaf node
	case *asyncapi.Tag:
		return w.walkTag(ctx, n)
	case *asyncapi.ExternalDoc:
		return nil // leaf node
	case *asyncapi.SecurityScheme:
		return w.walkSecurityScheme(ctx, n)
	case *asyncapi.OAuthFlows:
		return w.walkOAuthFlows(ctx, n)
	case *asyncapi.OAuthFlow:
		return nil // leaf node
	case *asyncapi.Parameter:
		return nil // leaf node
	case *highbase.SchemaProxy:
		return w.walkSchemaProxy(ctx, n)
	case *highbase.Schema:
		return w.walkSchema(ctx, n)
	// Binding types - visit but no children to walk
	case *asyncapi.HTTPServerBinding, *asyncapi.KafkaServerBinding, *asyncapi.MQTTServerBinding, *asyncapi.SQSServerBinding:
		return nil
	case *asyncapi.HTTPChannelBinding, *asyncapi.WebSocketChannelBinding, *asyncapi.KafkaChannelBinding, *asyncapi.AMQPChannelBinding, *asyncapi.SQSChannelBinding:
		return nil
	case *asyncapi.HTTPOperationBinding, *asyncapi.KafkaOperationBinding, *asyncapi.AMQPOperationBinding, *asyncapi.MQTTOperationBinding, *asyncapi.SQSOperationBinding:
		return nil
	case *asyncapi.HTTPMessageBinding, *asyncapi.KafkaMessageBinding, *asyncapi.AMQPMessageBinding, *asyncapi.MQTTMessageBinding, *asyncapi.SQSMessageBinding:
		return nil
	// Reference types - leaf nodes (resolved items are visited in their collections)
	case *low.Reference:
		return nil
	}

	return nil
}

// walkAsyncAPI walks the root document.
func (w *Walker) walkAsyncAPI(ctx context.Context, doc *asyncapi.AsyncAPI) error {
	childCtx := withChild(ctx, doc)

	if doc.Info != nil {
		if err := w.walkNode(AppendPath(childCtx, "info"), doc.Info); err != nil {
			return err
		}
	}
	if err := w.walkServersMap(childCtx, "servers", doc.Servers); err != nil {
		return err
	}
	if err := w.walkChannelsMap(childCtx, "channels", doc.Channels); err != nil {
		return err
	}
	if err := w.walkOperationsMap(childCtx, "operations", doc.Operations); err != nil {
		return err
	}
	if doc.Components != nil {
		if err := w.walkNode(AppendPath(childCtx, "components"), doc.Components); err != nil {
			return err
		}
	}
	return nil
}

// walkInfo walks the info object.
func (w *Walker) walkInfo(ctx context.Context, info *asyncapi.Info) error {
	childCtx := withChild(ctx, info)

	if info.Contact != nil {
		if err := w.walkNode(AppendPath(childCtx, "contact"), info.Contact); err != nil {
			return err
		}
	}
	if info.License != nil {
		if err := w.walkNode(AppendPath(childCtx, "license"), info.License); err != nil {
			return err
		}
	}
	if info.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), info.ExternalDocs); err != nil {
			return err
		}
	}
	if err := w.walkTagsSlice(childCtx, "tags", info.Tags); err != nil {
		return err
	}
	return nil
}

// walkServer walks a server object.
func (w *Walker) walkServer(ctx context.Context, server *asyncapi.Server) error {
	childCtx := withChild(ctx, server)

	if err := w.walkServerVariablesMap(childCtx, "variables", server.Variables); err != nil {
		return err
	}
	if err := w.walkSecuritySchemesSlice(childCtx, "security", server.Security); err != nil {
		return err
	}
	if err := w.walkTagsSlice(childCtx, "tags", server.Tags); err != nil {
		return err
	}
	if server.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), server.ExternalDocs); err != nil {
			return err
		}
	}
	if server.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), server.Bindings); err != nil {
			return err
		}
	}
	return nil
}

// walkServerBindings walks server bindings.
func (w *Walker) walkServerBindings(ctx context.Context, bindings *asyncapi.ServerBindings) error {
	childCtx := withChild(ctx, bindings)

	if bindings.HTTP != nil {
		if err := w.walkNode(AppendPath(childCtx, "http"), bindings.HTTP); err != nil {
			return err
		}
	}
	if bindings.Kafka != nil {
		if err := w.walkNode(AppendPath(childCtx, "kafka"), bindings.Kafka); err != nil {
			return err
		}
	}
	if bindings.MQTT != nil {
		if err := w.walkNode(AppendPath(childCtx, "mqtt"), bindings.MQTT); err != nil {
			return err
		}
	}
	if bindings.SQS != nil {
		if err := w.walkNode(AppendPath(childCtx, "sqs"), bindings.SQS); err != nil {
			return err
		}
	}
	return nil
}

// walkChannel walks a channel object.
func (w *Walker) walkChannel(ctx context.Context, channel *asyncapi.Channel) error {
	childCtx := withChild(ctx, channel)

	if err := w.walkReferencesSlice(childCtx, "servers", channel.Servers); err != nil {
		return err
	}
	if err := w.walkMessagesMap(childCtx, "messages", channel.Messages); err != nil {
		return err
	}
	if err := w.walkParametersMap(childCtx, "parameters", channel.Parameters); err != nil {
		return err
	}
	if err := w.walkTagsSlice(childCtx, "tags", channel.Tags); err != nil {
		return err
	}
	if channel.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), channel.ExternalDocs); err != nil {
			return err
		}
	}
	if channel.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), channel.Bindings); err != nil {
			return err
		}
	}
	return nil
}

// walkChannelBindings walks channel bindings.
func (w *Walker) walkChannelBindings(ctx context.Context, bindings *asyncapi.ChannelBindings) error {
	childCtx := withChild(ctx, bindings)

	if bindings.HTTP != nil {
		if err := w.walkNode(AppendPath(childCtx, "http"), bindings.HTTP); err != nil {
			return err
		}
	}
	if bindings.WebSocket != nil {
		if err := w.walkNode(AppendPath(childCtx, "ws"), bindings.WebSocket); err != nil {
			return err
		}
	}
	if bindings.Kafka != nil {
		if err := w.walkNode(AppendPath(childCtx, "kafka"), bindings.Kafka); err != nil {
			return err
		}
	}
	if bindings.AMQP != nil {
		if err := w.walkNode(AppendPath(childCtx, "amqp"), bindings.AMQP); err != nil {
			return err
		}
	}
	if bindings.SQS != nil {
		if err := w.walkNode(AppendPath(childCtx, "sqs"), bindings.SQS); err != nil {
			return err
		}
	}
	return nil
}

// walkOperation walks an operation object.
func (w *Walker) walkOperation(ctx context.Context, op *asyncapi.Operation) error {
	childCtx := withChild(ctx, op)

	if op.Channel != nil {
		if err := w.walkNode(AppendPath(childCtx, "channel"), op.Channel); err != nil {
			return err
		}
	}
	if err := w.walkReferencesSlice(childCtx, "messages", op.Messages); err != nil {
		return err
	}
	if op.Reply != nil {
		if err := w.walkNode(AppendPath(childCtx, "reply"), op.Reply); err != nil {
			return err
		}
	}
	if err := w.walkTagsSlice(childCtx, "tags", op.Tags); err != nil {
		return err
	}
	if op.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), op.ExternalDocs); err != nil {
			return err
		}
	}
	if op.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), op.Bindings); err != nil {
			return err
		}
	}
	if err := w.walkOperationTraitsSlice(childCtx, "traits", op.Traits); err != nil {
		return err
	}
	if err := w.walkSecuritySchemesSlice(childCtx, "security", op.Security); err != nil {
		return err
	}
	return nil
}

// walkOperationBindings walks operation bindings.
func (w *Walker) walkOperationBindings(ctx context.Context, bindings *asyncapi.OperationBindings) error {
	childCtx := withChild(ctx, bindings)

	if bindings.HTTP != nil {
		if err := w.walkNode(AppendPath(childCtx, "http"), bindings.HTTP); err != nil {
			return err
		}
	}
	if bindings.Kafka != nil {
		if err := w.walkNode(AppendPath(childCtx, "kafka"), bindings.Kafka); err != nil {
			return err
		}
	}
	if bindings.AMQP != nil {
		if err := w.walkNode(AppendPath(childCtx, "amqp"), bindings.AMQP); err != nil {
			return err
		}
	}
	if bindings.MQTT != nil {
		if err := w.walkNode(AppendPath(childCtx, "mqtt"), bindings.MQTT); err != nil {
			return err
		}
	}
	if bindings.SQS != nil {
		if err := w.walkNode(AppendPath(childCtx, "sqs"), bindings.SQS); err != nil {
			return err
		}
	}
	return nil
}

// walkOperationTrait walks an operation trait.
func (w *Walker) walkOperationTrait(ctx context.Context, trait *asyncapi.OperationTrait) error {
	childCtx := withChild(ctx, trait)

	if err := w.walkSecuritySchemesSlice(childCtx, "security", trait.Security); err != nil {
		return err
	}
	if err := w.walkTagsSlice(childCtx, "tags", trait.Tags); err != nil {
		return err
	}
	if trait.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), trait.ExternalDocs); err != nil {
			return err
		}
	}
	if trait.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), trait.Bindings); err != nil {
			return err
		}
	}
	return nil
}

// walkOperationReply walks an operation reply.
func (w *Walker) walkOperationReply(ctx context.Context, reply *asyncapi.OperationReply) error {
	childCtx := withChild(ctx, reply)

	if reply.Address != nil {
		if err := w.walkNode(AppendPath(childCtx, "address"), reply.Address); err != nil {
			return err
		}
	}
	if reply.Channel != nil {
		if err := w.walkNode(AppendPath(childCtx, "channel"), reply.Channel); err != nil {
			return err
		}
	}
	if err := w.walkReferencesSlice(childCtx, "messages", reply.Messages); err != nil {
		return err
	}
	return nil
}

// walkComponents walks the components object.
func (w *Walker) walkComponents(ctx context.Context, comp *asyncapi.Components) error {
	childCtx := withChild(ctx, comp)

	// Walk all component collections
	if err := w.walkSchemasMap(childCtx, "schemas", comp.Schemas); err != nil {
		return err
	}
	if err := w.walkServersMap(childCtx, "servers", comp.Servers); err != nil {
		return err
	}
	if err := w.walkChannelsMap(childCtx, "channels", comp.Channels); err != nil {
		return err
	}
	if err := w.walkOperationsMap(childCtx, "operations", comp.Operations); err != nil {
		return err
	}
	if err := w.walkMessagesMap(childCtx, "messages", comp.Messages); err != nil {
		return err
	}
	if err := w.walkSecuritySchemesMap(childCtx, "securitySchemes", comp.SecuritySchemes); err != nil {
		return err
	}
	if err := w.walkServerVariablesMap(childCtx, "serverVariables", comp.ServerVariables); err != nil {
		return err
	}
	if err := w.walkParametersMap(childCtx, "parameters", comp.Parameters); err != nil {
		return err
	}
	if err := w.walkCorrelationIDsMap(childCtx, "correlationIds", comp.CorrelationIDs); err != nil {
		return err
	}
	if err := w.walkRepliesMap(childCtx, "replies", comp.Replies); err != nil {
		return err
	}
	if err := w.walkReplyAddressesMap(childCtx, "replyAddresses", comp.ReplyAddresses); err != nil {
		return err
	}
	if err := w.walkExternalDocsMap(childCtx, "externalDocs", comp.ExternalDocs); err != nil {
		return err
	}
	if err := w.walkTagsMap(childCtx, "tags", comp.Tags); err != nil {
		return err
	}
	if err := w.walkOperationTraitsMap(childCtx, "operationTraits", comp.OperationTraits); err != nil {
		return err
	}
	if err := w.walkMessageTraitsMap(childCtx, "messageTraits", comp.MessageTraits); err != nil {
		return err
	}
	if err := w.walkServerBindingsMap(childCtx, "serverBindings", comp.ServerBindings); err != nil {
		return err
	}
	if err := w.walkChannelBindingsMap(childCtx, "channelBindings", comp.ChannelBindings); err != nil {
		return err
	}
	if err := w.walkOperationBindingsMap(childCtx, "operationBindings", comp.OperationBindings); err != nil {
		return err
	}
	if err := w.walkMessageBindingsMap(childCtx, "messageBindings", comp.MessageBindings); err != nil {
		return err
	}
	return nil
}

// walkMessage walks a message object.
func (w *Walker) walkMessage(ctx context.Context, msg *asyncapi.Message) error {
	childCtx := withChild(ctx, msg)

	if msg.Headers != nil {
		if err := w.walkNode(AppendPath(childCtx, "headers"), msg.Headers); err != nil {
			return err
		}
	}
	if msg.Payload != nil {
		if err := w.walkNode(AppendPath(childCtx, "payload"), msg.Payload); err != nil {
			return err
		}
	}
	if msg.CorrelationID != nil {
		if err := w.walkNode(AppendPath(childCtx, "correlationId"), msg.CorrelationID); err != nil {
			return err
		}
	}
	if err := w.walkTagsSlice(childCtx, "tags", msg.Tags); err != nil {
		return err
	}
	if msg.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), msg.ExternalDocs); err != nil {
			return err
		}
	}
	if msg.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), msg.Bindings); err != nil {
			return err
		}
	}
	if err := w.walkMessageExamplesSlice(childCtx, "examples", msg.Examples); err != nil {
		return err
	}
	if err := w.walkMessageTraitsSlice(childCtx, "traits", msg.Traits); err != nil {
		return err
	}
	return nil
}

// walkMessageBindings walks message bindings.
func (w *Walker) walkMessageBindings(ctx context.Context, bindings *asyncapi.MessageBindings) error {
	childCtx := withChild(ctx, bindings)

	if bindings.HTTP != nil {
		if err := w.walkNode(AppendPath(childCtx, "http"), bindings.HTTP); err != nil {
			return err
		}
	}
	if bindings.Kafka != nil {
		if err := w.walkNode(AppendPath(childCtx, "kafka"), bindings.Kafka); err != nil {
			return err
		}
	}
	if bindings.AMQP != nil {
		if err := w.walkNode(AppendPath(childCtx, "amqp"), bindings.AMQP); err != nil {
			return err
		}
	}
	if bindings.MQTT != nil {
		if err := w.walkNode(AppendPath(childCtx, "mqtt"), bindings.MQTT); err != nil {
			return err
		}
	}
	if bindings.SQS != nil {
		if err := w.walkNode(AppendPath(childCtx, "sqs"), bindings.SQS); err != nil {
			return err
		}
	}
	return nil
}

// walkMessageTrait walks a message trait object.
func (w *Walker) walkMessageTrait(ctx context.Context, trait *asyncapi.MessageTrait) error {
	childCtx := withChild(ctx, trait)

	if trait.Headers != nil {
		if err := w.walkNode(AppendPath(childCtx, "headers"), trait.Headers); err != nil {
			return err
		}
	}
	if trait.CorrelationID != nil {
		if err := w.walkNode(AppendPath(childCtx, "correlationId"), trait.CorrelationID); err != nil {
			return err
		}
	}
	if err := w.walkTagsSlice(childCtx, "tags", trait.Tags); err != nil {
		return err
	}
	if trait.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), trait.ExternalDocs); err != nil {
			return err
		}
	}
	if trait.Bindings != nil {
		if err := w.walkNode(AppendPath(childCtx, "bindings"), trait.Bindings); err != nil {
			return err
		}
	}
	if err := w.walkMessageExamplesSlice(childCtx, "examples", trait.Examples); err != nil {
		return err
	}
	return nil
}

// walkTag walks a tag object.
func (w *Walker) walkTag(ctx context.Context, tag *asyncapi.Tag) error {
	childCtx := withChild(ctx, tag)

	if tag.ExternalDocs != nil {
		if err := w.walkNode(AppendPath(childCtx, "externalDocs"), tag.ExternalDocs); err != nil {
			return err
		}
	}
	return nil
}

// walkSecurityScheme walks a security scheme object.
func (w *Walker) walkSecurityScheme(ctx context.Context, sec *asyncapi.SecurityScheme) error {
	childCtx := withChild(ctx, sec)

	if sec.Flows != nil {
		if err := w.walkNode(AppendPath(childCtx, "flows"), sec.Flows); err != nil {
			return err
		}
	}
	return nil
}

// walkOAuthFlows walks OAuth flows.
func (w *Walker) walkOAuthFlows(ctx context.Context, flows *asyncapi.OAuthFlows) error {
	childCtx := withChild(ctx, flows)

	if flows.Implicit != nil {
		if err := w.walkNode(AppendPath(childCtx, "implicit"), flows.Implicit); err != nil {
			return err
		}
	}
	if flows.Password != nil {
		if err := w.walkNode(AppendPath(childCtx, "password"), flows.Password); err != nil {
			return err
		}
	}
	if flows.ClientCredentials != nil {
		if err := w.walkNode(AppendPath(childCtx, "clientCredentials"), flows.ClientCredentials); err != nil {
			return err
		}
	}
	if flows.AuthorizationCode != nil {
		if err := w.walkNode(AppendPath(childCtx, "authorizationCode"), flows.AuthorizationCode); err != nil {
			return err
		}
	}
	return nil
}

// walkSchemaProxy handles SchemaProxy with stack-based cycle detection.
// EnterSchema and LeaveSchema are called for SchemaVisitor implementations.
// LeaveSchema is guaranteed to fire via defer only if EnterSchema succeeded.
// If EnterSchema returns an error, LeaveSchema is NOT called.
func (w *Walker) walkSchemaProxy(ctx context.Context, proxy *highbase.SchemaProxy) (retErr error) {
	if proxy == nil {
		return nil
	}

	sv, isSchemaVisitor := w.visitor.(SchemaVisitor)
	var schema *highbase.Schema
	var enteredSchema bool // Track if EnterSchema was called

	// Set up panic-safe defer for LeaveSchema - only fires if EnterSchema was called
	if isSchemaVisitor {
		defer func() {
			if enteredSchema {
				sv.LeaveSchema(ctx, proxy, schema, retErr)
			}
		}()
	}

	// Build unique key for cycle detection.
	// For references, use only the ref string (not pointer) because different
	// SchemaProxy instances may point to the same $ref.
	// For inline schemas, use the pointer address.
	var key schemaKey
	if proxy.IsReference() {
		key = schemaKey{ref: proxy.GetReference()}
	} else {
		key = schemaKey{ptr: reflect.ValueOf(proxy).Pointer()}
	}

	// Check if currently on stack (true cycle - would cause infinite recursion)
	stack := Stack(ctx)
	if stack[key] {
		if isSchemaVisitor {
			retErr = sv.SkipCircularRef(ctx, proxy, key.ref)
		}
		return retErr
	}

	// Push onto stack before descending
	stack[key] = true
	defer delete(stack, key) // Always pop from stack

	// Notify SchemaVisitor of entry
	if isSchemaVisitor {
		if err := sv.EnterSchema(ctx, proxy); err != nil {
			retErr = err
			return retErr
		}
		enteredSchema = true // Mark that EnterSchema succeeded
	}

	// Resolve the schema
	schema = proxy.Schema()
	if schema != nil {
		// Walk the resolved schema with proxy as parent
		retErr = w.walkNode(withChild(ctx, proxy), schema)
	}

	return retErr
}

// walkSchema walks a resolved schema, including all schema subtrees.
func (w *Walker) walkSchema(ctx context.Context, schema *highbase.Schema) error {
	childCtx := withChild(ctx, schema)

	// Walk polymorphic schemas (allOf, oneOf, anyOf)
	if err := w.walkPolymorphic(childCtx, schema, "allOf", schema.AllOf); err != nil {
		return err
	}
	if err := w.walkPolymorphic(childCtx, schema, "oneOf", schema.OneOf); err != nil {
		return err
	}
	if err := w.walkPolymorphic(childCtx, schema, "anyOf", schema.AnyOf); err != nil {
		return err
	}

	// Walk properties
	if err := w.walkSchemasMap(childCtx, "properties", schema.Properties); err != nil {
		return err
	}

	// Walk items (array schema)
	if schema.Items != nil {
		if items := schema.Items.A; items != nil {
			if err := w.walkNode(AppendPath(childCtx, "items"), items); err != nil {
				return err
			}
		}
	}

	// Walk additionalProperties
	if schema.AdditionalProperties != nil {
		if addProps := schema.AdditionalProperties.A; addProps != nil {
			if err := w.walkNode(AppendPath(childCtx, "additionalProperties"), addProps); err != nil {
				return err
			}
		}
	}

	// Walk not
	if schema.Not != nil {
		if err := w.walkNode(AppendPath(childCtx, "not"), schema.Not); err != nil {
			return err
		}
	}

	// Walk 3.1 specific properties
	if err := w.walkSchemaProxiesSlice(childCtx, "prefixItems", schema.PrefixItems); err != nil {
		return err
	}
	if schema.Contains != nil {
		if err := w.walkNode(AppendPath(childCtx, "contains"), schema.Contains); err != nil {
			return err
		}
	}
	if schema.If != nil {
		if err := w.walkNode(AppendPath(childCtx, "if"), schema.If); err != nil {
			return err
		}
	}
	if schema.Then != nil {
		if err := w.walkNode(AppendPath(childCtx, "then"), schema.Then); err != nil {
			return err
		}
	}
	if schema.Else != nil {
		if err := w.walkNode(AppendPath(childCtx, "else"), schema.Else); err != nil {
			return err
		}
	}
	if err := w.walkSchemasMap(childCtx, "dependentSchemas", schema.DependentSchemas); err != nil {
		return err
	}
	if err := w.walkSchemasMap(childCtx, "patternProperties", schema.PatternProperties); err != nil {
		return err
	}
	if schema.PropertyNames != nil {
		if err := w.walkNode(AppendPath(childCtx, "propertyNames"), schema.PropertyNames); err != nil {
			return err
		}
	}
	if schema.UnevaluatedItems != nil {
		if err := w.walkNode(AppendPath(childCtx, "unevaluatedItems"), schema.UnevaluatedItems); err != nil {
			return err
		}
	}
	if schema.UnevaluatedProperties != nil {
		if unevalProps := schema.UnevaluatedProperties.A; unevalProps != nil {
			if err := w.walkNode(AppendPath(childCtx, "unevaluatedProperties"), unevalProps); err != nil {
				return err
			}
		}
	}
	if schema.ContentSchema != nil {
		if err := w.walkNode(AppendPath(childCtx, "contentSchema"), schema.ContentSchema); err != nil {
			return err
		}
	}

	return nil
}

// walkPolymorphic handles allOf/oneOf/anyOf with guaranteed Leave callbacks.
func (w *Walker) walkPolymorphic(ctx context.Context, schema *highbase.Schema, name string, proxies []*highbase.SchemaProxy) (retErr error) {
	if len(proxies) == 0 {
		return nil
	}

	polyCtx := AppendPath(ctx, name)
	pv, isPoly := w.visitor.(PolymorphicVisitor)

	// Enter callback
	if isPoly {
		var err error
		switch name {
		case "allOf":
			err = pv.EnterAllOf(polyCtx, schema, len(proxies))
		case "oneOf":
			err = pv.EnterOneOf(polyCtx, schema, len(proxies))
		case "anyOf":
			err = pv.EnterAnyOf(polyCtx, schema, len(proxies))
		}
		if err != nil {
			return err
		}
	}

	// Leave callback - ALWAYS fires via defer
	if isPoly {
		defer func() {
			switch name {
			case "allOf":
				pv.LeaveAllOf(polyCtx, schema, retErr)
			case "oneOf":
				pv.LeaveOneOf(polyCtx, schema, retErr)
			case "anyOf":
				pv.LeaveAnyOf(polyCtx, schema, retErr)
			}
		}()
	}

	// Walk children
	for i, proxy := range proxies {
		if err := w.walkNode(AppendIndex(polyCtx, i), proxy); err != nil {
			retErr = err
			return retErr // Leave callback will still fire via defer
		}
	}

	return nil
}

// walkSchemaProxiesSlice walks a slice of schema proxies.
func (w *Walker) walkSchemaProxiesSlice(ctx context.Context, name string, proxies []*highbase.SchemaProxy) error {
	if len(proxies) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, proxy := range proxies {
		if err := w.walkNode(AppendIndex(sliceCtx, i), proxy); err != nil {
			return err
		}
	}
	return nil
}

// Map iteration helpers - all iterate in insertion order via FromOldest()

func (w *Walker) walkServersMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Server]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkChannelsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Channel]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkOperationsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Operation]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkMessagesMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Message]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkMessageTraitsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.MessageTrait]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkServerVariablesMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.ServerVariable]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkParametersMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Parameter]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkSecuritySchemesMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.SecurityScheme]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkCorrelationIDsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.CorrelationID]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkRepliesMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.OperationReply]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkReplyAddressesMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.OperationReplyAddress]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkExternalDocsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.ExternalDoc]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkTagsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.Tag]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkOperationTraitsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.OperationTrait]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkServerBindingsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.ServerBindings]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkChannelBindingsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.ChannelBindings]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkOperationBindingsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.OperationBindings]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkMessageBindingsMap(ctx context.Context, name string, m *orderedmap.Map[string, *asyncapi.MessageBindings]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

// walkSchemasMap walks an ordered map of schema proxies.
// Calls walkNode for SchemaProxy so basic visitors see the proxy in pre-order.
func (w *Walker) walkSchemasMap(ctx context.Context, name string, m *orderedmap.Map[string, *highbase.SchemaProxy]) error {
	if m == nil {
		return nil
	}
	mapCtx := AppendPath(ctx, name)
	for k, v := range m.FromOldest() {
		if err := w.walkNode(AppendPath(mapCtx, k), v); err != nil {
			return err
		}
	}
	return nil
}

// Slice iteration helpers

func (w *Walker) walkTagsSlice(ctx context.Context, name string, items []*asyncapi.Tag) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkMessageTraitsSlice(ctx context.Context, name string, items []*asyncapi.MessageTrait) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkOperationTraitsSlice(ctx context.Context, name string, items []*asyncapi.OperationTrait) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkSecuritySchemesSlice(ctx context.Context, name string, items []*asyncapi.SecurityScheme) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkMessageExamplesSlice(ctx context.Context, name string, items []*asyncapi.MessageExample) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}

func (w *Walker) walkReferencesSlice(ctx context.Context, name string, items []*low.Reference) error {
	if len(items) == 0 {
		return nil
	}
	sliceCtx := AppendPath(ctx, name)
	for i, item := range items {
		if err := w.walkNode(AppendIndex(sliceCtx, i), item); err != nil {
			return err
		}
	}
	return nil
}
