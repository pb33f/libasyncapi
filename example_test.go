// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package libasyncapi_test

import (
	"fmt"
	"log"

	"github.com/pb33f/libasyncapi"
)

func ExampleNewDocument() {
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

	// Output:
	// Orders API uses AsyncAPI 3.0.0
	// channel orders: /orders
}
