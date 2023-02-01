// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package ptrace

import (
	"testing"

	"github.com/stretchr/testify/assert"

	otlptrace "go.opentelemetry.io/collector/pdata/internal/data/protogen/trace/v1"
)

func TestSpanEventSlice(t *testing.T) {
	es := NewSpanEventSlice()
	assert.Equal(t, 0, es.Len())
	es = newSpanEventSlice(&[]*otlptrace.Span_Event{})
	assert.Equal(t, 0, es.Len())

	es.EnsureCapacity(7)
	emptyVal := newSpanEvent(&otlptrace.Span_Event{})
	testVal := generateTestSpanEvent()
	assert.Equal(t, 7, cap(*es.orig))
	for i := 0; i < es.Len(); i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, el)
		fillTestSpanEvent(el)
		assert.Equal(t, testVal, el)
	}
}

func TestSpanEventSlice_CopyTo(t *testing.T) {
	dest := NewSpanEventSlice()
	// Test CopyTo to empty
	NewSpanEventSlice().CopyTo(dest)
	assert.Equal(t, NewSpanEventSlice(), dest)

	// Test CopyTo larger slice
	generateTestSpanEventSlice().CopyTo(dest)
	assert.Equal(t, generateTestSpanEventSlice(), dest)

	// Test CopyTo same size slice
	generateTestSpanEventSlice().CopyTo(dest)
	assert.Equal(t, generateTestSpanEventSlice(), dest)
}

func TestSpanEventSlice_EnsureCapacity(t *testing.T) {
	es := generateTestSpanEventSlice()
	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	expectedEs := make(map[*otlptrace.Span_Event]bool)
	for i := 0; i < es.Len(); i++ {
		expectedEs[es.At(i).orig] = true
	}
	assert.Equal(t, es.Len(), len(expectedEs))
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	foundEs := make(map[*otlptrace.Span_Event]bool, es.Len())
	for i := 0; i < es.Len(); i++ {
		foundEs[es.At(i).orig] = true
	}
	assert.Equal(t, expectedEs, foundEs)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	oldLen := es.Len()
	expectedEs = make(map[*otlptrace.Span_Event]bool, oldLen)
	for i := 0; i < oldLen; i++ {
		expectedEs[es.At(i).orig] = true
	}
	assert.Equal(t, oldLen, len(expectedEs))
	es.EnsureCapacity(ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
	foundEs = make(map[*otlptrace.Span_Event]bool, oldLen)
	for i := 0; i < oldLen; i++ {
		foundEs[es.At(i).orig] = true
	}
	assert.Equal(t, expectedEs, foundEs)
}

func TestSpanEventSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestSpanEventSlice()
	dest := NewSpanEventSlice()
	src := generateTestSpanEventSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestSpanEventSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestSpanEventSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestSpanEventSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestSpanEventSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewSpanEventSlice()
	emptySlice.RemoveIf(func(el SpanEvent) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestSpanEventSlice()
	pos := 0
	filtered.RemoveIf(func(el SpanEvent) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func generateTestSpanEventSlice() SpanEventSlice {
	tv := NewSpanEventSlice()
	fillTestSpanEventSlice(tv)
	return tv
}

func fillTestSpanEventSlice(tv SpanEventSlice) {
	*tv.orig = make([]*otlptrace.Span_Event, 7)
	for i := 0; i < 7; i++ {
		(*tv.orig)[i] = &otlptrace.Span_Event{}
		fillTestSpanEvent(newSpanEvent((*tv.orig)[i]))
	}
}
