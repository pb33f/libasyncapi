// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/pb33f/libopenapi/datamodel"
)

// DocumentConfiguration holds configuration options for parsing AsyncAPI documents.
type DocumentConfiguration struct {
	// SpecFilePath identifies the entry file, allowing references back to it.
	// Relative paths are resolved against BasePath.
	SpecFilePath string

	// ExcludeExtensionRefs leaves $ref values inside x-* extensions untouched.
	ExcludeExtensionRefs bool

	// AllowExcessiveAliases disables the indexing expansion limit for trusted YAML.
	// Alias cycles are still rejected. The default is false.
	AllowExcessiveAliases bool

	// BasePath is the base path for resolving relative file references.
	// If empty, the current working directory is used.
	// Bundling applies this base to entry-document imports while preserving
	// SpecFilePath as the entry's identity and fragment-only reference target.
	BasePath string

	// BaseURL is the base URL for resolving relative URL references.
	// When bundling, it takes precedence over BasePath for entry-document imports.
	BaseURL *url.URL

	// AllowFileReferences enables resolution of file:// references.
	// Defaults to false.
	AllowFileReferences bool

	// AllowRemoteReferences enables resolution of http:// and https:// references.
	// Defaults to false.
	AllowRemoteReferences bool

	// RemoteURLHandler is a custom function for fetching remote references.
	// If nil, the default http client is used.
	RemoteURLHandler func(url string) (*http.Response, error)

	// LocalFS is a custom filesystem for resolving local file references.
	// If nil, the OS filesystem is used.
	LocalFS fs.FS

	// SkipCircularReferenceCheck disables circular reference detection.
	// This can improve performance but may cause infinite loops if circular
	// references exist in the specification.
	SkipCircularReferenceCheck bool

	// ExtractRefsSequentially processes references one at a time instead of
	// in parallel. This is slower but can be useful for debugging.
	ExtractRefsSequentially bool

	// Logger is an optional structured logger for diagnostic messages.
	Logger *slog.Logger
}

// NewDocumentConfiguration creates a new DocumentConfiguration with default values.
func NewDocumentConfiguration() *DocumentConfiguration {
	return &DocumentConfiguration{
		Logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		})),
	}
}

// ToLibOpenAPIConfig converts this configuration to libopenapi's configuration format.
func (c *DocumentConfiguration) ToLibOpenAPIConfig() *datamodel.DocumentConfiguration {
	if c == nil {
		return datamodel.NewDocumentConfiguration()
	}

	return &datamodel.DocumentConfiguration{
		SpecFilePath:               c.SpecFilePath,
		ExcludeExtensionRefs:       c.ExcludeExtensionRefs,
		AllowExcessiveAliases:      c.AllowExcessiveAliases,
		BasePath:                   c.BasePath,
		BaseURL:                    c.BaseURL,
		AllowFileReferences:        c.AllowFileReferences,
		AllowRemoteReferences:      c.AllowRemoteReferences,
		SkipCircularReferenceCheck: c.SkipCircularReferenceCheck,
		ExtractRefsSequentially:    c.ExtractRefsSequentially,
		Logger:                     c.Logger,
		RemoteURLHandler:           c.RemoteURLHandler,
		LocalFS:                    c.LocalFS,
	}
}
