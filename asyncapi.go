// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package libasyncapi is a library containing tools for reading and manipulating AsyncAPI 3+
// specifications into strongly typed documents. These documents have two APIs, a high level
// (porcelain) and a low level (plumbing).
//
// Every single type has a 'GoLow()' method that drops down from the high API to the low API.
// Once in the low API, the entire original document data is available, including all comments,
// line and column numbers for keys and values.
//
// # Supported Versions
//
// This library supports AsyncAPI 3.0 and later. AsyncAPI 2.x specifications are not supported
// and will return ErrAsyncAPI2NotSupported.
//
// # Basic Usage
//
//	spec, err := os.ReadFile("asyncapi.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	doc, err := libasyncapi.NewDocument(spec)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	model := doc.Model()
//	fmt.Printf("API: %s\n", model.Info.Title)
//
// # Error Handling
//
// The library uses lenient parsing - it accumulates errors and continues parsing where possible.
// Use doc.IsPartial() to check if errors occurred, and doc.Errors() to retrieve them.
// The model may still be usable even with errors.
//
// # Multi-File Specifications
//
// For specifications split across multiple files, use NewDocumentWithConfiguration and set the
// BasePath to enable file reference resolution:
//
//	config := libasyncapi.NewDocumentConfiguration()
//	config.BasePath = "/path/to/spec/directory"
//	config.AllowFileReferences = true
//
//	doc, err := libasyncapi.NewDocumentWithConfiguration(spec, config)
package libasyncapi

import (
	"errors"
	"fmt"

	"github.com/pb33f/go-yaml"
	highasync "github.com/pb33f/libasyncapi/datamodel/high/asyncapi"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
)

// NewDocument will create a new AsyncAPI instance from an AsyncAPI specification []byte array.
// If anything goes wrong when parsing, reading or processing the specification, there will be
// no document returned, instead an error will be returned that explains what failed.
//
// This function will NOT automatically follow (meaning load) any file or remote references.
// If you need to load external references, use NewDocumentWithConfiguration instead.
func NewDocument(spec []byte) (Document, error) {
	return NewDocumentWithConfiguration(spec, nil)
}

// NewDocumentWithConfiguration is the same as NewDocument, except it accepts a configuration
// that allows control over file and remote reference resolution.
func NewDocumentWithConfiguration(spec []byte, config *DocumentConfiguration) (Document, error) {
	if config == nil {
		config = NewDocumentConfiguration()
	}

	doc := &document{
		config: config,
		errors: make([]error, 0, 4),
	}

	// Parse YAML
	var rootNode yaml.Node
	if err := yaml.Unmarshal(spec, &rootNode); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidYAML, err)
	}
	doc.rootNode = &rootNode

	// Extract and validate version
	version, err := ExtractVersionFromNode(&rootNode)
	if err != nil {
		return nil, err
	}
	doc.version = version

	specInfo, err := ParseAsyncAPIVersion(version)
	if err != nil {
		return nil, err
	}
	doc.specInfo = specInfo

	// Convert to low-level document configuration
	lowConfig := &lowasync.DocumentConfiguration{
		BasePath:                   config.BasePath,
		BaseURL:                    config.BaseURL,
		AllowFileReferences:        config.AllowFileReferences,
		AllowRemoteReferences:      config.AllowRemoteReferences,
		SkipCircularReferenceCheck: config.SkipCircularReferenceCheck,
		ExtractRefsSequentially:    config.ExtractRefsSequentially,
		Logger:                     config.Logger,
		LocalFS:                    config.LocalFS,
	}
	if config.RemoteURLHandler != nil {
		lowConfig.RemoteURLHandler = lowasync.RemoteURLHandler(config.RemoteURLHandler)
	}

	// Build low-level document
	lowDoc, buildErr := lowasync.CreateDocumentWithConfig(&rootNode, lowConfig)
	if buildErr != nil {
		// Check if it's a multi-error
		var docErr *lowasync.DocumentBuildError
		if errors.As(buildErr, &docErr) {
			doc.errors = append(doc.errors, docErr.Errors...)
		} else {
			doc.addError(buildErr)
		}
	}

	// If we have a low-level document, build the high-level model
	if lowDoc != nil {
		doc.lowModel = lowDoc
		doc.specIndex = lowDoc.Index
		doc.rolodex = lowDoc.Rolodex

		// Build high-level model
		doc.highModel = highasync.NewAsyncAPI(lowDoc)
	}

	// Return document even if there were errors (lenient parsing)
	if len(doc.errors) > 0 && doc.lowModel == nil {
		return nil, errors.Join(doc.errors...)
	}

	return doc, nil
}
