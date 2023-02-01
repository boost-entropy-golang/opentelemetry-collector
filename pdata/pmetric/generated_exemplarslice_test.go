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

package pmetric

import (
	"testing"

	"github.com/stretchr/testify/assert"

	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

func TestExemplarSlice(t *testing.T) {
	es := NewExemplarSlice()
	assert.Equal(t, 0, es.Len())
	es = newExemplarSlice(&[]otlpmetrics.Exemplar{})
	assert.Equal(t, 0, es.Len())

	es.EnsureCapacity(7)
	emptyVal := newExemplar(&otlpmetrics.Exemplar{})
	testVal := generateTestExemplar()
	assert.Equal(t, 7, cap(*es.orig))
	for i := 0; i < es.Len(); i++ {
		el := es.AppendEmpty()
		assert.Equal(t, emptyVal, el)
		fillTestExemplar(el)
		assert.Equal(t, testVal, el)
	}
}

func TestExemplarSlice_CopyTo(t *testing.T) {
	dest := NewExemplarSlice()
	// Test CopyTo to empty
	NewExemplarSlice().CopyTo(dest)
	assert.Equal(t, NewExemplarSlice(), dest)

	// Test CopyTo larger slice
	generateTestExemplarSlice().CopyTo(dest)
	assert.Equal(t, generateTestExemplarSlice(), dest)

	// Test CopyTo same size slice
	generateTestExemplarSlice().CopyTo(dest)
	assert.Equal(t, generateTestExemplarSlice(), dest)
}

func TestExemplarSlice_EnsureCapacity(t *testing.T) {
	es := generateTestExemplarSlice()
	// Test ensure smaller capacity.
	const ensureSmallLen = 4
	expectedEs := make(map[*otlpmetrics.Exemplar]bool)
	for i := 0; i < es.Len(); i++ {
		expectedEs[es.At(i).orig] = true
	}
	assert.Equal(t, es.Len(), len(expectedEs))
	es.EnsureCapacity(ensureSmallLen)
	assert.Less(t, ensureSmallLen, es.Len())
	foundEs := make(map[*otlpmetrics.Exemplar]bool, es.Len())
	for i := 0; i < es.Len(); i++ {
		foundEs[es.At(i).orig] = true
	}
	assert.Equal(t, expectedEs, foundEs)

	// Test ensure larger capacity
	const ensureLargeLen = 9
	oldLen := es.Len()
	assert.Equal(t, oldLen, len(expectedEs))
	es.EnsureCapacity(ensureLargeLen)
	assert.Equal(t, ensureLargeLen, cap(*es.orig))
}

func TestExemplarSlice_MoveAndAppendTo(t *testing.T) {
	// Test MoveAndAppendTo to empty
	expectedSlice := generateTestExemplarSlice()
	dest := NewExemplarSlice()
	src := generateTestExemplarSlice()
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestExemplarSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo empty slice
	src.MoveAndAppendTo(dest)
	assert.Equal(t, generateTestExemplarSlice(), dest)
	assert.Equal(t, 0, src.Len())
	assert.Equal(t, expectedSlice.Len(), dest.Len())

	// Test MoveAndAppendTo not empty slice
	generateTestExemplarSlice().MoveAndAppendTo(dest)
	assert.Equal(t, 2*expectedSlice.Len(), dest.Len())
	for i := 0; i < expectedSlice.Len(); i++ {
		assert.Equal(t, expectedSlice.At(i), dest.At(i))
		assert.Equal(t, expectedSlice.At(i), dest.At(i+expectedSlice.Len()))
	}
}

func TestExemplarSlice_RemoveIf(t *testing.T) {
	// Test RemoveIf on empty slice
	emptySlice := NewExemplarSlice()
	emptySlice.RemoveIf(func(el Exemplar) bool {
		t.Fail()
		return false
	})

	// Test RemoveIf
	filtered := generateTestExemplarSlice()
	pos := 0
	filtered.RemoveIf(func(el Exemplar) bool {
		pos++
		return pos%3 == 0
	})
	assert.Equal(t, 5, filtered.Len())
}

func generateTestExemplarSlice() ExemplarSlice {
	tv := NewExemplarSlice()
	fillTestExemplarSlice(tv)
	return tv
}

func fillTestExemplarSlice(tv ExemplarSlice) {
	*tv.orig = make([]otlpmetrics.Exemplar, 7)
	for i := 0; i < 7; i++ {
		fillTestExemplar(newExemplar(&(*tv.orig)[i]))
	}
}
