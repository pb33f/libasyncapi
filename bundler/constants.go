// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package bundler

import lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"

// Traversal kinds distinguish document structure from content copied verbatim.
const (
	kindRoot          = "root"
	kindOpaque        = "opaque"
	kindSchemaContent = "schemaContent"
	kindExtension     = "extension"
	mapKindPrefix     = "map:"
	listKindPrefix    = "list:"
	refKindPrefix     = "ref:"
	bindingKindPrefix = "binding:"
	bindingsSuffix    = "Bindings"
)

const (
	componentsPathPrefix = "/" + lowasync.ComponentsLabel + "/"
	channelsPathPrefix   = "/" + lowasync.ChannelsLabel + "/"
	operationsPathPrefix = "/" + lowasync.OperationsLabel + "/"
	messagesPathSegment  = "/" + lowasync.MessagesLabel + "/"
	channelRefPath       = "/" + lowasync.ChannelLabel + "/" + lowasync.RefLabel
	extensionPrefix      = "x-"
	defaultDelimiter     = "__"
	defaultComponentName = "imported"
)

// Legacy JSON Schema keywords are not exposed by libopenapi's label constants.
const (
	definitionsLabel     = "definitions"
	dependenciesLabel    = "dependencies"
	additionalItemsLabel = "additionalItems"
	recursiveRefLabel    = "$recursiveRef"
)

const (
	fileScheme       = "file"
	httpScheme       = "http"
	httpsScheme      = "https"
	localFileHost    = "localhost"
	defaultEntryFile = "__asyncapi_entry__.yaml"
	defaultRemoteFS  = "default"
)
