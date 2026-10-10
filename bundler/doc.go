// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

// Package bundler combines AsyncAPI 3 documents and their referenced files.
// Bundling preserves root-local references and fields that require Reference
// Objects. It does not validate application semantics or expand traits.
package bundler
