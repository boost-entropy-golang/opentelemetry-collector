// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package ptrace

import (
	"go.opentelemetry.io/collector/pdata/internal"
	"go.opentelemetry.io/collector/pdata/internal/data"
	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

// Span represents a single operation within a trace.
// See Span definition in OTLP: https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/trace/v1/trace.proto
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewSpan function to create new instances.
// Important: zero-initialized instance is not valid for use.
type Span struct {
	orig *otlptrace.Span
}

func newSpan(orig *otlptrace.Span) Span {
	return Span{orig}
}

// NewSpan creates a new empty Span.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewSpan() Span {
	return newSpan(&otlptrace.Span{})
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms Span) MoveTo(dest Span) {
	*dest.orig = *ms.orig
	*ms.orig = otlptrace.Span{}
}

// TraceID returns the traceid associated with this Span.
func (ms Span) TraceID() pcommon.TraceID {
	return pcommon.TraceID(ms.orig.TraceId)
}

// SetTraceID replaces the traceid associated with this Span.
func (ms Span) SetTraceID(v pcommon.TraceID) {
	ms.orig.TraceId = data.TraceID(v)
}

// SpanID returns the spanid associated with this Span.
func (ms Span) SpanID() pcommon.SpanID {
	return pcommon.SpanID(ms.orig.SpanId)
}

// SetSpanID replaces the spanid associated with this Span.
func (ms Span) SetSpanID(v pcommon.SpanID) {
	ms.orig.SpanId = data.SpanID(v)
}

// TraceState returns the tracestate associated with this Span.
func (ms Span) TraceState() pcommon.TraceState {
	return pcommon.TraceState(internal.NewTraceState(&ms.orig.TraceState))
}

// ParentSpanID returns the parentspanid associated with this Span.
func (ms Span) ParentSpanID() pcommon.SpanID {
	return pcommon.SpanID(ms.orig.ParentSpanId)
}

// SetParentSpanID replaces the parentspanid associated with this Span.
func (ms Span) SetParentSpanID(v pcommon.SpanID) {
	ms.orig.ParentSpanId = data.SpanID(v)
}

// Name returns the name associated with this Span.
func (ms Span) Name() string {
	return ms.orig.Name
}

// SetName replaces the name associated with this Span.
func (ms Span) SetName(v string) {
	ms.orig.Name = v
}

// Kind returns the kind associated with this Span.
func (ms Span) Kind() SpanKind {
	return SpanKind(ms.orig.Kind)
}

// SetKind replaces the kind associated with this Span.
func (ms Span) SetKind(v SpanKind) {
	ms.orig.Kind = otlptrace.Span_SpanKind(v)
}

// StartTimestamp returns the starttimestamp associated with this Span.
func (ms Span) StartTimestamp() pcommon.Timestamp {
	return pcommon.Timestamp(ms.orig.StartTimeUnixNano)
}

// SetStartTimestamp replaces the starttimestamp associated with this Span.
func (ms Span) SetStartTimestamp(v pcommon.Timestamp) {
	ms.orig.StartTimeUnixNano = uint64(v)
}

// EndTimestamp returns the endtimestamp associated with this Span.
func (ms Span) EndTimestamp() pcommon.Timestamp {
	return pcommon.Timestamp(ms.orig.EndTimeUnixNano)
}

// SetEndTimestamp replaces the endtimestamp associated with this Span.
func (ms Span) SetEndTimestamp(v pcommon.Timestamp) {
	ms.orig.EndTimeUnixNano = uint64(v)
}

// Attributes returns the Attributes associated with this Span.
func (ms Span) Attributes() pcommon.Map {
	return pcommon.Map(internal.NewMap(&ms.orig.Attributes))
}

// DroppedAttributesCount returns the droppedattributescount associated with this Span.
func (ms Span) DroppedAttributesCount() uint32 {
	return ms.orig.DroppedAttributesCount
}

// SetDroppedAttributesCount replaces the droppedattributescount associated with this Span.
func (ms Span) SetDroppedAttributesCount(v uint32) {
	ms.orig.DroppedAttributesCount = v
}

// Events returns the Events associated with this Span.
func (ms Span) Events() SpanEventSlice {
	return newSpanEventSlice(&ms.orig.Events)
}

// DroppedEventsCount returns the droppedeventscount associated with this Span.
func (ms Span) DroppedEventsCount() uint32 {
	return ms.orig.DroppedEventsCount
}

// SetDroppedEventsCount replaces the droppedeventscount associated with this Span.
func (ms Span) SetDroppedEventsCount(v uint32) {
	ms.orig.DroppedEventsCount = v
}

// Links returns the Links associated with this Span.
func (ms Span) Links() SpanLinkSlice {
	return newSpanLinkSlice(&ms.orig.Links)
}

// DroppedLinksCount returns the droppedlinkscount associated with this Span.
func (ms Span) DroppedLinksCount() uint32 {
	return ms.orig.DroppedLinksCount
}

// SetDroppedLinksCount replaces the droppedlinkscount associated with this Span.
func (ms Span) SetDroppedLinksCount(v uint32) {
	ms.orig.DroppedLinksCount = v
}

// Status returns the status associated with this Span.
func (ms Span) Status() Status {
	return newStatus(&ms.orig.Status)
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms Span) CopyTo(dest Span) {
	dest.SetTraceID(ms.TraceID())
	dest.SetSpanID(ms.SpanID())
	ms.TraceState().CopyTo(dest.TraceState())
	dest.SetParentSpanID(ms.ParentSpanID())
	dest.SetName(ms.Name())
	dest.SetKind(ms.Kind())
	dest.SetStartTimestamp(ms.StartTimestamp())
	dest.SetEndTimestamp(ms.EndTimestamp())
	ms.Attributes().CopyTo(dest.Attributes())
	dest.SetDroppedAttributesCount(ms.DroppedAttributesCount())
	ms.Events().CopyTo(dest.Events())
	dest.SetDroppedEventsCount(ms.DroppedEventsCount())
	ms.Links().CopyTo(dest.Links())
	dest.SetDroppedLinksCount(ms.DroppedLinksCount())
	ms.Status().CopyTo(dest.Status())
}
