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

package pmetric // import "go.opentelemetry.io/collector/pdata/pmetric"

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpcollectormetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/collector/metrics/v1"
	otlpmetrics "go.opentelemetry.io/collector/pdata/internal/data/protogen/metrics/v1"
)

// Metrics is the top-level struct that is propagated through the metrics pipeline.
// Use NewMetrics to create new instance, zero-initialized instance is not valid for use.
type Metrics internal.Metrics

func newMetrics(orig *otlpcollectormetrics.ExportMetricsServiceRequest) Metrics {
	return Metrics(internal.NewMetrics(orig))
}

func (ms Metrics) getOrig() *otlpcollectormetrics.ExportMetricsServiceRequest {
	return internal.GetOrigMetrics(internal.Metrics(ms))
}

// NewMetrics creates a new Metrics struct.
func NewMetrics() Metrics {
	return newMetrics(&otlpcollectormetrics.ExportMetricsServiceRequest{})
}

// Clone returns a copy of MetricData.
func (ms Metrics) Clone() Metrics {
	cloneMd := NewMetrics()
	ms.ResourceMetrics().CopyTo(cloneMd.ResourceMetrics())
	return cloneMd
}

// MoveTo moves all properties from the current struct to dest
// resetting the current instance to its zero value.
func (ms Metrics) MoveTo(dest Metrics) {
	*dest.getOrig() = *ms.getOrig()
	*ms.getOrig() = otlpcollectormetrics.ExportMetricsServiceRequest{}
}

// ResourceMetrics returns the ResourceMetricsSlice associated with this Metrics.
func (ms Metrics) ResourceMetrics() ResourceMetricsSlice {
	return newResourceMetricsSlice(&ms.getOrig().ResourceMetrics)
}

// MetricCount calculates the total number of metrics.
func (ms Metrics) MetricCount() int {
	metricCount := 0
	rms := ms.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			metricCount += ilm.Metrics().Len()
		}
	}
	return metricCount
}

// DataPointCount calculates the total number of data points.
func (ms Metrics) DataPointCount() (dataPointCount int) {
	rms := ms.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			ilm := ilms.At(j)
			ms := ilm.Metrics()
			for k := 0; k < ms.Len(); k++ {
				m := ms.At(k)
				switch m.Type() {
				case MetricTypeGauge:
					dataPointCount += m.Gauge().DataPoints().Len()
				case MetricTypeSum:
					dataPointCount += m.Sum().DataPoints().Len()
				case MetricTypeHistogram:
					dataPointCount += m.Histogram().DataPoints().Len()
				case MetricTypeExponentialHistogram:
					dataPointCount += m.ExponentialHistogram().DataPoints().Len()
				case MetricTypeSummary:
					dataPointCount += m.Summary().DataPoints().Len()
				}
			}
		}
	}
	return
}

// Deprecated: [v0.61.0] use MetricType.
type MetricDataType = MetricType

const (
	// Deprecated: [v0.61.0] use MetricTypeNone.
	MetricDataTypeNone = MetricTypeNone
	// Deprecated: [v0.61.0] use MetricTypeGauge.
	MetricDataTypeGauge = MetricTypeGauge
	// Deprecated: [v0.61.0] use MetricTypeSum.
	MetricDataTypeSum = MetricTypeSum
	// Deprecated: [v0.61.0] use MetricTypeHistogram.
	MetricDataTypeHistogram = MetricTypeHistogram
	// Deprecated: [v0.61.0] use MetricTypeExponentialHistogram.
	MetricDataTypeExponentialHistogram = MetricTypeExponentialHistogram
	// Deprecated: [v0.61.0] use MetricTypeSummary.
	MetricDataTypeSummary = MetricTypeSummary
)

// Deprecated: [v0.61.0] use Metric.Type().
func (ms Metric) DataType() MetricDataType {
	return ms.Type()
}

// MetricType specifies the type of data in a Metric.
type MetricType int32

const (
	MetricTypeNone MetricType = iota
	MetricTypeGauge
	MetricTypeSum
	MetricTypeHistogram
	MetricTypeExponentialHistogram
	MetricTypeSummary
)

// String returns the string representation of the MetricType.
func (mdt MetricType) String() string {
	switch mdt {
	case MetricTypeNone:
		return "None"
	case MetricTypeGauge:
		return "Gauge"
	case MetricTypeSum:
		return "Sum"
	case MetricTypeHistogram:
		return "Histogram"
	case MetricTypeExponentialHistogram:
		return "ExponentialHistogram"
	case MetricTypeSummary:
		return "Summary"
	}
	return ""
}

// MetricAggregationTemporality defines how a metric aggregator reports aggregated values.
// It describes how those values relate to the time interval over which they are aggregated.
type MetricAggregationTemporality int32

const (
	// MetricAggregationTemporalityUnspecified is the default MetricAggregationTemporality, it MUST NOT be used.
	MetricAggregationTemporalityUnspecified = MetricAggregationTemporality(otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED)
	// MetricAggregationTemporalityDelta is a MetricAggregationTemporality for a metric aggregator which reports changes since last report time.
	MetricAggregationTemporalityDelta = MetricAggregationTemporality(otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA)
	// MetricAggregationTemporalityCumulative is a MetricAggregationTemporality for a metric aggregator which reports changes since a fixed start time.
	MetricAggregationTemporalityCumulative = MetricAggregationTemporality(otlpmetrics.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE)
)

// String returns the string representation of the MetricAggregationTemporality.
func (at MetricAggregationTemporality) String() string {
	return otlpmetrics.AggregationTemporality(at).String()
}

// NumberDataPointValueType specifies the type of NumberDataPoint value.
type NumberDataPointValueType int32

const (
	NumberDataPointValueTypeNone NumberDataPointValueType = iota
	NumberDataPointValueTypeInt
	NumberDataPointValueTypeDouble
)

// String returns the string representation of the NumberDataPointValueType.
func (nt NumberDataPointValueType) String() string {
	switch nt {
	case NumberDataPointValueTypeNone:
		return "None"
	case NumberDataPointValueTypeInt:
		return "Int"
	case NumberDataPointValueTypeDouble:
		return "Double"
	}
	return ""
}

// ExemplarValueType specifies the type of Exemplar measurement value.
type ExemplarValueType int32

const (
	ExemplarValueTypeNone ExemplarValueType = iota
	ExemplarValueTypeInt
	ExemplarValueTypeDouble
)

// String returns the string representation of the ExemplarValueType.
func (nt ExemplarValueType) String() string {
	switch nt {
	case ExemplarValueTypeNone:
		return "None"
	case ExemplarValueTypeInt:
		return "Int"
	case ExemplarValueTypeDouble:
		return "Double"
	}
	return ""
}

// Deprecated: [v0.61.0] not used, will be deleted in next release.
type OptionalType int32

const (
	// Deprecated: [v0.61.0] not used, will be deleted in next release.
	OptionalTypeNone OptionalType = iota
	// Deprecated: [v0.61.0] not used, will be deleted in next release.
	OptionalTypeDouble
)

// Deprecated: [v0.61.0] not used, will be deleted in next release.
func (ot OptionalType) String() string {
	switch ot {
	case OptionalTypeNone:
		return "None"
	case OptionalTypeDouble:
		return "Double"
	}
	return ""
}

// Deprecated: [v0.61.0] use NumberDataPoint.DoubleValue().
func (ms NumberDataPoint) DoubleVal() float64 {
	return ms.DoubleValue()
}

// Deprecated: [v0.61.0] use NumberDataPoint.SetDoubleValue().
func (ms NumberDataPoint) SetDoubleVal(v float64) {
	ms.SetDoubleValue(v)
}

// Deprecated: [v0.61.0] use NumberDataPoint.IntValue().
func (ms NumberDataPoint) IntVal() int64 {
	return ms.IntValue()
}

// Deprecated: [v0.61.0] use NumberDataPoint.SetIntValue().
func (ms NumberDataPoint) SetIntVal(v int64) {
	ms.SetIntValue(v)
}

// Deprecated: [v0.61.0] use Exemplar.DoubleValue().
func (ms Exemplar) DoubleVal() float64 {
	return ms.DoubleValue()
}

// Deprecated: [v0.61.0] use Exemplar.SetDoubleValue().
func (ms Exemplar) SetDoubleVal(v float64) {
	ms.SetDoubleValue(v)
}

// Deprecated: [v0.61.0] use Exemplar.IntValue().
func (ms Exemplar) IntVal() int64 {
	return ms.IntValue()
}

// Deprecated: [v0.61.0] use Exemplar.SetIntValue().
func (ms Exemplar) SetIntVal(v int64) {
	ms.SetIntValue(v)
}
