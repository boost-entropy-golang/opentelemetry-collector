// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pmetricotlp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/collector/pdata/internal"
	otlpcollectormetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/collector/metrics/v1"
)

func TestExportPartialSuccess_MoveTo(t *testing.T) {
	ms := generateTestExportPartialSuccess()
	dest := NewExportPartialSuccess()
	ms.MoveTo(dest)
	assert.Equal(t, NewExportPartialSuccess(), ms)
	assert.Equal(t, generateTestExportPartialSuccess(), dest)
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() {
		ms.MoveTo(newExportPartialSuccess(&otlpcollectormetrics.ExportMetricsPartialSuccess{}, &sharedState))
	})
	assert.Panics(t, func() {
		newExportPartialSuccess(&otlpcollectormetrics.ExportMetricsPartialSuccess{}, &sharedState).MoveTo(dest)
	})
}

func TestExportPartialSuccess_CopyTo(t *testing.T) {
	ms := NewExportPartialSuccess()
	orig := NewExportPartialSuccess()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	orig = generateTestExportPartialSuccess()
	orig.CopyTo(ms)
	assert.Equal(t, orig, ms)
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() {
		ms.CopyTo(newExportPartialSuccess(&otlpcollectormetrics.ExportMetricsPartialSuccess{}, &sharedState))
	})
}

func TestExportPartialSuccess_RejectedDataPoints(t *testing.T) {
	ms := NewExportPartialSuccess()
	assert.Equal(t, int64(0), ms.RejectedDataPoints())
	ms.SetRejectedDataPoints(int64(13))
	assert.Equal(t, int64(13), ms.RejectedDataPoints())
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() {
		newExportPartialSuccess(&otlpcollectormetrics.ExportMetricsPartialSuccess{}, &sharedState).SetRejectedDataPoints(int64(13))
	})
}

func TestExportPartialSuccess_ErrorMessage(t *testing.T) {
	ms := NewExportPartialSuccess()
	assert.Empty(t, ms.ErrorMessage())
	ms.SetErrorMessage("error message")
	assert.Equal(t, "error message", ms.ErrorMessage())
	sharedState := internal.StateReadOnly
	assert.Panics(t, func() {
		newExportPartialSuccess(&otlpcollectormetrics.ExportMetricsPartialSuccess{}, &sharedState).SetErrorMessage("error message")
	})
}

func generateTestExportPartialSuccess() ExportPartialSuccess {
	tv := NewExportPartialSuccess()
	fillTestExportPartialSuccess(tv)
	return tv
}

func fillTestExportPartialSuccess(tv ExportPartialSuccess) {
	tv.orig.RejectedDataPoints = int64(13)
	tv.orig.ErrorMessage = "error message"
}
