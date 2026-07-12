// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package reports

import (
	"os"
	"testing"

	libasyncapi "github.com/pb33f/libasyncapi"
	lowasync "github.com/pb33f/libasyncapi/datamodel/low/asyncapi"
	"github.com/pb33f/libasyncapi/what-changed/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateOverallReport_NilChanges(t *testing.T) {
	report := CreateOverallReport(nil)
	require.NotNil(t, report)
	assert.Empty(t, report.ChangeReport)
}

func TestCreateOverallReport(t *testing.T) {
	o, err := os.ReadFile("../test_fixtures/comprehensive-bindings.yaml")
	require.NoError(t, err)
	m, err := os.ReadFile("../test_fixtures/comprehensive-bindings-modified.yaml")
	require.NoError(t, err)

	od, err := libasyncapi.NewDocument(o)
	require.NoError(t, err)
	md, err := libasyncapi.NewDocument(m)
	require.NoError(t, err)
	require.NotNil(t, od.Model())
	require.NotNil(t, md.Model())

	changes := model.CompareDocuments(od.GoLow(), md.GoLow())
	require.NotNil(t, changes)

	report := CreateOverallReport(changes)
	require.NotNil(t, report)

	// info: title change, non-breaking.
	require.Contains(t, report.ChangeReport, lowasync.InfoLabel)
	assert.Equal(t, 1, report.ChangeReport[lowasync.InfoLabel].Total)
	assert.Equal(t, 0, report.ChangeReport[lowasync.InfoLabel].Breaking)

	// servers: host + keepAlive changes, both breaking.
	require.Contains(t, report.ChangeReport, lowasync.ServersLabel)
	assert.Equal(t, 2, report.ChangeReport[lowasync.ServersLabel].Total)
	assert.Equal(t, 2, report.ChangeReport[lowasync.ServersLabel].Breaking)

	// channels and operations and components are all represented.
	require.Contains(t, report.ChangeReport, lowasync.ChannelsLabel)
	require.Contains(t, report.ChangeReport, lowasync.OperationsLabel)
	require.Contains(t, report.ChangeReport, lowasync.ComponentsLabel)

	// totals across the report match the document totals. the channel add/remove land
	// under the root "channels" property bucket via the property merge.
	total := 0
	breaking := 0
	for _, c := range report.ChangeReport {
		total += c.Total
		breaking += c.Breaking
	}
	assert.Equal(t, changes.TotalChanges(), total)
	assert.Equal(t, changes.TotalBreakingChanges(), breaking)
}
