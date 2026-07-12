// Copyright 2026 Princess Beef Heavy Industries / Dave Shanley
// SPDX-License-Identifier: MIT

package what_changed

import (
	"os"
	"testing"

	libasyncapi "github.com/pb33f/libasyncapi"
	"github.com/pb33f/libasyncapi/what-changed/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadDocument(t *testing.T, path string) libasyncapi.Document {
	t.Helper()
	spec, err := os.ReadFile(path)
	require.NoError(t, err)
	doc, err := libasyncapi.NewDocument(spec)
	require.NoError(t, err)
	require.NotNil(t, doc.Model())
	return doc
}

func compareFixtures(t *testing.T) *model.DocumentChanges {
	t.Helper()
	original := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	modified := loadDocument(t, "test_fixtures/comprehensive-bindings-modified.yaml")
	return CompareDocuments(original, modified)
}

func TestCompareAsyncAPIDocuments(t *testing.T) {
	original := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	modified := loadDocument(t, "test_fixtures/comprehensive-bindings-modified.yaml")

	changes := CompareAsyncAPIDocuments(original.GoLow(), modified.GoLow())
	require.NotNil(t, changes)

	// the modified fixture makes 19 deliberate changes; the message-level edits
	// (amqp contentEncoding, kafka payload schema) are $refs resolved into both the
	// channel and components subtrees, so each counts twice.
	assert.Equal(t, 19, changes.TotalChanges())
	assert.Equal(t, 16, changes.TotalBreakingChanges())
	assert.Len(t, changes.GetAllChanges(), 19)
}

func TestCompareDocuments_InfoChange(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// info.title was modified: non-breaking.
	require.NotNil(t, changes.InfoChanges)
	assert.Equal(t, 1, changes.InfoChanges.TotalChanges())
	assert.Equal(t, 0, changes.InfoChanges.TotalBreakingChanges())
	ch := changes.InfoChanges.Changes[0]
	assert.Equal(t, model.Modified, ch.ChangeType)
	assert.Equal(t, "Comprehensive Bindings Test", ch.Original)
	assert.Equal(t, "Comprehensive Bindings Test v2", ch.New)
}

func TestCompareDocuments_ServerHostChange(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// kafka-server host was modified: breaking.
	require.Contains(t, changes.ServerChanges, "kafka-server")
	sc := changes.ServerChanges["kafka-server"]
	assert.Equal(t, 1, sc.TotalChanges())
	assert.Equal(t, 1, sc.TotalBreakingChanges())
	assert.Equal(t, "kafka2.example.com:9092", sc.Changes[0].New)

	// mqtt-server keepAlive was modified inside the mqtt server binding: breaking.
	require.Contains(t, changes.ServerChanges, "mqtt-server")
	mqtt := changes.ServerChanges["mqtt-server"]
	require.NotNil(t, mqtt.BindingsChanges)
	require.NotNil(t, mqtt.BindingsChanges.MQTTChanges)
	assert.Equal(t, 1, mqtt.BindingsChanges.MQTTChanges.TotalBreakingChanges())
}

func TestCompareDocuments_ChannelBindingChanges(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// kafka-channel binding topic + partitions were modified: both breaking.
	require.Contains(t, changes.ChannelChanges, "kafka-channel")
	kafka := changes.ChannelChanges["kafka-channel"]
	require.NotNil(t, kafka.BindingsChanges)
	require.NotNil(t, kafka.BindingsChanges.KafkaChanges)
	assert.Equal(t, 2, kafka.BindingsChanges.KafkaChanges.TotalChanges())
	assert.Equal(t, 2, kafka.BindingsChanges.KafkaChanges.TotalBreakingChanges())

	// amqp-channel exchange type was modified: breaking.
	require.Contains(t, changes.ChannelChanges, "amqp-channel")
	amqp := changes.ChannelChanges["amqp-channel"]
	require.NotNil(t, amqp.BindingsChanges)
	require.NotNil(t, amqp.BindingsChanges.AMQPChanges)
	assert.Equal(t, 1, amqp.BindingsChanges.AMQPChanges.TotalBreakingChanges())

	// sqs-channel queue name was modified: breaking.
	require.Contains(t, changes.ChannelChanges, "sqs-channel")
	sqs := changes.ChannelChanges["sqs-channel"]
	require.NotNil(t, sqs.BindingsChanges)
	require.NotNil(t, sqs.BindingsChanges.SQSChanges)
	assert.Equal(t, 1, sqs.BindingsChanges.SQSChanges.TotalBreakingChanges())
}

func TestCompareDocuments_ChannelAddedAndRemoved(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// ws-channel was removed (breaking), audit-channel was added (not breaking).
	var added, removed *model.Change
	for _, ch := range changes.Changes {
		switch ch.ChangeType {
		case model.ObjectAdded:
			added = ch
		case model.ObjectRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
	assert.Equal(t, "audit-channel", added.New)
	assert.False(t, added.Breaking)
	assert.Equal(t, "ws-channel", removed.Original)
	assert.True(t, removed.Breaking)
}

func TestCompareDocuments_OperationChanges(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// mqtt-operation changed its channel $ref (ref-only comparison) and its qos.
	require.Contains(t, changes.OperationChanges, "mqtt-operation")
	mqtt := changes.OperationChanges["mqtt-operation"]
	assert.Equal(t, 2, mqtt.TotalChanges())
	assert.Equal(t, 2, mqtt.TotalBreakingChanges())
	var refChange *model.Change
	for _, ch := range mqtt.Changes {
		if ch.Property == "channel" {
			refChange = ch
		}
	}
	require.NotNil(t, refChange)
	assert.Equal(t, model.Modified, refChange.ChangeType)
	assert.Equal(t, "#/channels/kafka-channel", refChange.Original)
	assert.Equal(t, "#/channels/http-channel", refChange.New)

	// http-operation method POST -> PUT: breaking.
	require.Contains(t, changes.OperationChanges, "http-operation")
	http := changes.OperationChanges["http-operation"]
	require.NotNil(t, http.BindingsChanges)
	require.NotNil(t, http.BindingsChanges.HTTPChanges)
	assert.Equal(t, 1, http.BindingsChanges.HTTPChanges.TotalBreakingChanges())
}

func TestCompareDocuments_DelegatedSchemaChange(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// the kafkaMessage payload schema changed type string -> integer; the schema
	// comparison is delegated to libopenapi and shows up under components.
	require.NotNil(t, changes.ComponentsChanges)
	require.Contains(t, changes.ComponentsChanges.MessageChanges, "kafkaMessage")
	km := changes.ComponentsChanges.MessageChanges["kafkaMessage"]
	require.NotNil(t, km.PayloadChanges)
	assert.Positive(t, km.PayloadChanges.TotalChanges())
	assert.Positive(t, km.PayloadChanges.TotalBreakingChanges())

	// the amqpMessage contentEncoding change lands in the message binding.
	require.Contains(t, changes.ComponentsChanges.MessageChanges, "amqpMessage")
	am := changes.ComponentsChanges.MessageChanges["amqpMessage"]
	require.NotNil(t, am.BindingsChanges)
	require.NotNil(t, am.BindingsChanges.AMQPChanges)
	assert.Equal(t, 1, am.BindingsChanges.AMQPChanges.TotalBreakingChanges())
}

func TestCompareDocuments_ChannelServerRefSwap(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// kafka-channel swapped its server ref from kafka-server to mqtt-server: one
	// breaking removal plus one non-breaking addition, with ref strings rendered.
	require.Contains(t, changes.ChannelChanges, "kafka-channel")
	kafka := changes.ChannelChanges["kafka-channel"]
	var added, removed *model.Change
	for _, ch := range kafka.Changes {
		if ch.Property != "servers" {
			continue
		}
		switch ch.ChangeType {
		case model.ObjectAdded:
			added = ch
		case model.ObjectRemoved:
			removed = ch
		}
	}
	require.NotNil(t, added)
	require.NotNil(t, removed)
	assert.Equal(t, "#/servers/mqtt-server", added.New)
	assert.False(t, added.Breaking)
	assert.Equal(t, "#/servers/kafka-server", removed.Original)
	assert.True(t, removed.Breaking)
}

func TestCompareDocuments_SecuritySchemeTypeChange(t *testing.T) {
	changes := compareFixtures(t)
	require.NotNil(t, changes)

	// the saslScram security scheme changed type scramSha256 -> scramSha512: breaking.
	require.NotNil(t, changes.ComponentsChanges)
	require.Contains(t, changes.ComponentsChanges.SecuritySchemeChanges, "saslScram")
	ss := changes.ComponentsChanges.SecuritySchemeChanges["saslScram"]
	assert.Equal(t, 1, ss.TotalChanges())
	assert.Equal(t, 1, ss.TotalBreakingChanges())
	assert.Equal(t, "scramSha256", ss.Changes[0].Original)
	assert.Equal(t, "scramSha512", ss.Changes[0].New)
}

func TestCompareDocuments_SelfCompareIsNil(t *testing.T) {
	original := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	same := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	assert.Nil(t, CompareDocuments(original, same))
}

func TestCompareDocuments_NilDocuments(t *testing.T) {
	original := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	assert.Nil(t, CompareDocuments(nil, nil))
	assert.Nil(t, CompareDocuments(original, nil))
	assert.Nil(t, CompareDocuments(nil, original))
}

func TestCompareDocuments_CustomBreakingRules(t *testing.T) {
	// by default the info.title change is not breaking: 14 breaking changes.
	changes := compareFixtures(t)
	require.NotNil(t, changes)
	assert.Equal(t, 16, changes.TotalBreakingChanges())

	// make info.title modifications breaking: 15 breaking changes.
	custom := model.NewDefaultBreakingRulesConfig()
	custom.Merge(&model.BreakingRulesConfig{
		Info: &model.InfoRules{
			Title: &model.BreakingChangeRule{Modified: boolRef(true)},
		},
	})
	original := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	modified := loadDocument(t, "test_fixtures/comprehensive-bindings-modified.yaml")
	changes = CompareDocumentsWithConfig(original, modified, custom)
	require.NotNil(t, changes)
	assert.Equal(t, 17, changes.TotalBreakingChanges())

	changes = CompareAsyncAPIDocumentsWithConfig(original.GoLow(), modified.GoLow(), custom)
	require.NotNil(t, changes)
	assert.Equal(t, 17, changes.TotalBreakingChanges())
}

func TestCompareDocumentsWithConfig_NilDocuments(t *testing.T) {
	doc := loadDocument(t, "test_fixtures/comprehensive-bindings.yaml")
	assert.Nil(t, CompareDocumentsWithConfig(nil, doc, nil))
	assert.Nil(t, CompareDocumentsWithConfig(doc, nil, nil))
}

func boolRef(b bool) *bool {
	return &b
}

func Benchmark_CompareAsyncAPIDocuments(b *testing.B) {
	o, _ := os.ReadFile("test_fixtures/comprehensive-bindings.yaml")
	m, _ := os.ReadFile("test_fixtures/comprehensive-bindings-modified.yaml")
	od, _ := libasyncapi.NewDocument(o)
	md, _ := libasyncapi.NewDocument(m)
	_ = od.Model()
	_ = md.Model()
	lo, lm := od.GoLow(), md.GoLow()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CompareAsyncAPIDocuments(lo, lm)
	}
}
