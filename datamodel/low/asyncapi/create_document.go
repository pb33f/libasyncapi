// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package asyncapi

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"time"

	"github.com/pb33f/go-yaml"
	"github.com/pb33f/libopenapi/datamodel/low"
	"github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"github.com/pb33f/libopenapi/utils"
)

// ErrNoAsyncAPIVersion is returned when the asyncapi version field is missing from the specification.
// This matches the naming convention used in the top-level package for consistency.
var ErrNoAsyncAPIVersion = errors.New("no asyncapi version found in document")

// modelCtxKey is the context key used by libopenapi for model context.
// This must match the string key used in libopenapi/datamodel/low/base/context.go
const modelCtxKey = "modelCtx"

// RemoteURLHandler is a function type for custom remote URL fetching.
type RemoteURLHandler = func(url string) (*http.Response, error)

// DocumentConfiguration provides options for loading and parsing AsyncAPI documents.
type DocumentConfiguration struct {
	// BasePath is the base path for resolving relative file references.
	BasePath string

	// BaseURL is the base URL for resolving relative remote references.
	BaseURL *url.URL

	// AllowFileReferences enables loading local file references.
	AllowFileReferences bool

	// AllowRemoteReferences enables loading remote references.
	AllowRemoteReferences bool

	// SkipCircularReferenceCheck skips checking for circular references.
	SkipCircularReferenceCheck bool

	// ExtractRefsSequentially extracts references one at a time instead of in parallel.
	ExtractRefsSequentially bool

	// RemoteURLHandler is a custom function for fetching remote URLs.
	RemoteURLHandler RemoteURLHandler

	// LocalFS is a custom filesystem for resolving local file references.
	// If nil, the OS filesystem is used.
	LocalFS fs.FS

	// Logger is the logger to use for debug output.
	Logger *slog.Logger
}

// NewDocumentConfiguration creates a new DocumentConfiguration with default values.
func NewDocumentConfiguration() *DocumentConfiguration {
	return &DocumentConfiguration{}
}

// CreateDocument creates a new AsyncAPI low-level document from the provided root node.
func CreateDocument(rootNode *yaml.Node) (*AsyncAPI, error) {
	return CreateDocumentWithConfig(rootNode, NewDocumentConfiguration())
}

// CreateDocumentWithConfig creates a new AsyncAPI low-level document with the provided configuration.
func CreateDocumentWithConfig(rootNode *yaml.Node, config *DocumentConfiguration) (*AsyncAPI, error) {
	if config == nil {
		config = NewDocumentConfiguration()
	}

	// Extract asyncapi version
	_, labelNode, versionNode := utils.FindKeyNodeFull(AsyncAPILabel, rootNode.Content)
	if versionNode == nil {
		return nil, ErrNoAsyncAPIVersion
	}

	doc := &AsyncAPI{
		AsyncAPI: low.NodeReference[string]{
			Value:     versionNode.Value,
			KeyNode:   labelNode,
			ValueNode: versionNode,
		},
	}
	doc.Nodes = low.ExtractNodes(nil, rootNode.Content[0])

	// Create index configuration
	idxConfig := index.CreateClosedAPIIndexConfig()
	idxConfig.SkipDocumentCheck = true // AsyncAPI documents don't have openapi/swagger fields
	idxConfig.AvoidCircularReferenceCheck = config.SkipCircularReferenceCheck
	idxConfig.ExtractRefsSequentially = config.ExtractRefsSequentially

	// Handle base URL
	baseURL := config.BaseURL
	idxConfig.BaseURL = urlWithoutTrailingSlash(baseURL)
	idxConfig.BasePath = config.BasePath

	if config.Logger != nil {
		idxConfig.Logger = config.Logger
	}

	// Create rolodex for reference resolution
	rolodex := index.NewRolodex(idxConfig)
	rolodex.SetRootNode(rootNode)
	doc.Rolodex = rolodex

	// If basePath is provided, add a local filesystem to the rolodex
	if idxConfig.BasePath != "" || config.AllowFileReferences {
		cwd, err := filepath.Abs(config.BasePath)
		if err != nil {
			return nil, errors.New("failed to resolve base path: " + err.Error())
		}

		// Create a local filesystem config
		localFSConf := index.LocalFSConfig{
			BaseDirectory: cwd,
			IndexConfig:   idxConfig,
		}

		// Use custom LocalFS if provided
		if config.LocalFS != nil {
			localFSConf.DirFS = config.LocalFS
		}

		fileFS, err := index.NewLocalFSWithConfig(&localFSConf)
		if err != nil {
			return nil, errors.New("failed to create local filesystem: " + err.Error())
		}
		idxConfig.AllowFileLookup = true

		// Add the filesystem to the rolodex
		rolodex.AddLocalFS(cwd, fileFS)
	}

	// If base URL is provided, add a remote filesystem to the rolodex
	if idxConfig.BaseURL != nil || config.AllowRemoteReferences {
		// Create a remote filesystem
		remoteFS, err := index.NewRemoteFSWithConfig(idxConfig)
		if err != nil {
			return nil, errors.New("failed to create remote filesystem: " + err.Error())
		}
		if config.RemoteURLHandler != nil {
			remoteFS.RemoteHandlerFunc = config.RemoteURLHandler
		}
		idxConfig.AllowRemoteLookup = true

		// Add to the rolodex
		u := "default"
		if config.BaseURL != nil {
			u = config.BaseURL.String()
		}
		rolodex.AddRemoteFS(u, remoteFS)
	}

	// Index all references
	var errs []error
	if config.Logger != nil {
		config.Logger.Debug("indexing rolodex")
	}
	now := time.Now()
	if indexErr := rolodex.IndexTheRolodex(context.Background()); indexErr != nil {
		errs = append(errs, indexErr)
	}
	if config.Logger != nil {
		config.Logger.Debug("rolodex indexed", "ms", time.Since(now).Milliseconds())
	}

	if !config.SkipCircularReferenceCheck {
		if config.Logger != nil {
			config.Logger.Debug("checking for circular references")
		}
		circNow := time.Now()
		rolodex.CheckForCircularReferences()
		if config.Logger != nil {
			config.Logger.Debug("circular check completed", "ms", time.Since(circNow).Milliseconds())
		}
	}

	// Extract rolodex errors
	roloErrs := rolodex.GetCaughtErrors()
	if roloErrs != nil {
		errs = append(errs, roloErrs...)
	}

	// Set root index
	doc.Index = rolodex.GetRootIndex()

	// Create context with schema cache
	// Note: must use string key "modelCtx" to match libopenapi's GetModelContext()
	var cacheMap sync.Map
	modelContext := base.ModelContext{SchemaCache: &cacheMap}
	ctx := context.WithValue(context.Background(), modelCtxKey, &modelContext)

	// Extract ID field
	_, idLabel, idNode := utils.FindKeyNodeFull(IDLabel, rootNode.Content)
	if idNode != nil {
		doc.ID = low.NodeReference[string]{
			Value:     idNode.Value,
			KeyNode:   idLabel,
			ValueNode: idNode,
		}
	}

	// Extract defaultContentType field
	_, dctLabel, dctNode := utils.FindKeyNodeFull(DefaultContentTypeLabel, rootNode.Content)
	if dctNode != nil {
		doc.DefaultContentType = low.NodeReference[string]{
			Value:     dctNode.Value,
			KeyNode:   dctLabel,
			ValueNode: dctNode,
		}
	}

	// Build the document model
	err := doc.Build(ctx, nil, rootNode.Content[0], doc.Index)
	if err != nil {
		errs = append(errs, err)
	}

	// Return document with any accumulated errors
	if len(errs) > 0 {
		return doc, &DocumentBuildError{Errors: errs}
	}

	return doc, nil
}

// urlWithoutTrailingSlash removes trailing slash from URL if present.
func urlWithoutTrailingSlash(u *url.URL) *url.URL {
	if u == nil {
		return nil
	}
	if len(u.Path) > 0 && u.Path[len(u.Path)-1] == '/' {
		newURL := *u
		newURL.Path = u.Path[:len(u.Path)-1]
		return &newURL
	}
	return u
}

// DocumentBuildError wraps multiple errors that occurred during document building.
type DocumentBuildError struct {
	Errors []error
}

func (e *DocumentBuildError) Error() string {
	if len(e.Errors) == 0 {
		return "unknown document build error"
	}
	return e.Errors[0].Error()
}

func (e *DocumentBuildError) Unwrap() []error {
	return e.Errors
}
