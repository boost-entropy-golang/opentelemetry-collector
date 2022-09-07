// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package plog

import (
	"testing"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	otlplogs "go.opentelemetry.io/collector/pdata/internal/data/protogen/logs/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

var logsOTLP = func() Logs {
	ld := NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	rl.Resource().Attributes().UpsertString("host.name", "testHost")
	rl.Resource().SetDroppedAttributesCount(1)
	rl.SetSchemaUrl("testSchemaURL")
	il := rl.ScopeLogs().AppendEmpty()
	il.Scope().SetName("name")
	il.Scope().SetVersion("version")
	il.Scope().SetDroppedAttributesCount(1)
	il.SetSchemaUrl("ScopeLogsSchemaURL")
	lg := il.LogRecords().AppendEmpty()
	lg.SetSeverityNumber(SeverityNumber(otlplogs.SeverityNumber_SEVERITY_NUMBER_ERROR))
	lg.SetSeverityText("Error")
	lg.SetDroppedAttributesCount(1)
	lg.SetFlags(LogRecordFlags(otlplogs.LogRecordFlags_LOG_RECORD_FLAG_UNSPECIFIED))
	traceID := pcommon.TraceID([16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10})
	spanID := pcommon.SpanID([8]byte{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18})
	lg.SetTraceID(traceID)
	lg.SetSpanID(spanID)
	lg.Body().SetStringVal("hello world")
	t := pcommon.NewTimestampFromTime(time.Now())
	lg.SetTimestamp(t)
	lg.SetObservedTimestamp(t)
	lg.Attributes().UpsertString("sdkVersion", "1.0.1")
	return ld
}

func TestReadLogsDataUnknownField(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := readLogsData(iter)
	assert.NoError(t, iter.Error)
	assert.EqualValues(t, otlplogs.LogsData{}, value)
}

func TestReadResourceLogsUnknownField(t *testing.T) {
	jsonStr := `{"extra":"","resource":{"extra":""}}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := readResourceLogs(iter)
	assert.NoError(t, iter.Error)
	assert.EqualValues(t, &otlplogs.ResourceLogs{}, value)
}

func TestReadScopeLogs(t *testing.T) {
	jsonStr := `{"extra":""}`
	iter := jsoniter.ConfigFastest.BorrowIterator([]byte(jsonStr))
	defer jsoniter.ConfigFastest.ReturnIterator(iter)
	value := readScopeLogs(iter)
	assert.NoError(t, iter.Error)
	assert.EqualValues(t, &otlplogs.ScopeLogs{}, value)
}

func TestLogsJSON(t *testing.T) {
	type args struct {
		logFunc func() Logs
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "otlp",
			args: args{
				logFunc: logsOTLP,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, opEnumsAsInts := range []bool{true, false} {
				for _, opEmitDefaults := range []bool{true, false} {
					for _, opOrigName := range []bool{true, false} {
						marshaller := &jsonMarshaler{
							delegate: jsonpb.Marshaler{
								EnumsAsInts:  opEnumsAsInts,
								EmitDefaults: opEmitDefaults,
								OrigName:     opOrigName,
							}}
						ld := tt.args.logFunc()
						jsonBuf, err := marshaller.MarshalLogs(ld)
						assert.NoError(t, err)
						decoder := NewJSONUnmarshaler()
						var got interface{}
						got, err = decoder.UnmarshalLogs(jsonBuf)
						assert.NoError(t, err)
						assert.EqualValues(t, ld, got)
					}
				}
			}
		})
	}
}

var logsJSON = `{"resourceLogs":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"testHost"}}]},"scopeLogs":[{"scope":{"name":"name","version":"version"},"logRecords":[{"severityText":"Error","body":{},"traceId":"","spanId":""}]}]}]}`

func TestLogsJSON_WithoutTraceIdAndSpanId(t *testing.T) {
	decoder := NewJSONUnmarshaler()
	_, err := decoder.UnmarshalLogs([]byte(logsJSON))
	assert.NoError(t, err)
}

var logsJSONWrongTraceID = `{"resourceLogs":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"testHost"}}]},"scopeLogs":[{"scope":{"name":"name","version":"version"},"logRecords":[{"severityText":"Error","body":{},"traceId":"--","spanId":""}]}]}]}`

func TestLogsJSON_WrongTraceID(t *testing.T) {
	decoder := NewJSONUnmarshaler()
	_, err := decoder.UnmarshalLogs([]byte(logsJSONWrongTraceID))
	assert.Error(t, err)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "parse trace_id")
	}
}

var LogsJSONWrongSpanID = `{"resourceLogs":[{"resource":{"attributes":[{"key":"host.name","value":{"stringValue":"testHost"}}]},"scopeLogs":[{"scope":{"name":"name","version":"version"},"logRecords":[{"severityText":"Error","body":{},"traceId":"","spanId":"--"}]}]}]}`

func TestLogsJSON_WrongSpanID(t *testing.T) {
	decoder := NewJSONUnmarshaler()
	_, err := decoder.UnmarshalLogs([]byte(LogsJSONWrongSpanID))
	assert.Error(t, err)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "parse span_id")
	}
}

func TestLogsJSON_MarshalJSON(t *testing.T) {
	marshal := NewJSONMarshaler()
	assert.NotNil(t, marshal)
}
