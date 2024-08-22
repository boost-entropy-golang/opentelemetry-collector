// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package e2e

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/connector/connectortest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/internal/sharedcomponent"
	"go.opentelemetry.io/collector/processor/processortest"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/service"
	"go.opentelemetry.io/collector/service/extensions"
	"go.opentelemetry.io/collector/service/pipelines"
	"go.opentelemetry.io/collector/service/telemetry"
)

var nopType = component.MustNewType("nop")

func Test_ComponentStatusReporting_SharedInstance(t *testing.T) {
	eventsReceived := make(map[*componentstatus.InstanceID][]*componentstatus.Event)

	set := service.Settings{
		BuildInfo:     component.NewDefaultBuildInfo(),
		CollectorConf: confmap.New(),
		ReceiversConfigs: map[component.ID]component.Config{
			component.NewID(component.MustNewType("test")): &receiverConfig{},
		},
		ReceiversFactories: map[component.Type]receiver.Factory{
			component.MustNewType("test"): newReceiverFactory(),
		},
		Processors: processortest.NewNopBuilder(),
		Exporters:  exportertest.NewNopBuilder(),
		Connectors: connectortest.NewNopBuilder(),
		Extensions: extension.NewBuilder(
			map[component.ID]component.Config{
				component.NewID(component.MustNewType("watcher")): &extensionConfig{eventsReceived},
			},
			map[component.Type]extension.Factory{
				component.MustNewType("watcher"): newExtensionFactory(),
			}),
	}
	set.BuildInfo = component.BuildInfo{Version: "test version", Command: "otelcoltest"}

	cfg := service.Config{
		Telemetry: telemetry.Config{
			Logs: telemetry.LogsConfig{
				Level:       zapcore.InfoLevel,
				Development: false,
				Encoding:    "console",
				Sampling: &telemetry.LogsSamplingConfig{
					Enabled:    true,
					Tick:       10 * time.Second,
					Initial:    100,
					Thereafter: 100,
				},
				OutputPaths:       []string{"stderr"},
				ErrorOutputPaths:  []string{"stderr"},
				DisableCaller:     false,
				DisableStacktrace: false,
				InitialFields:     map[string]any(nil),
			},
			Metrics: telemetry.MetricsConfig{
				Level: configtelemetry.LevelNone,
			},
		},
		Pipelines: pipelines.Config{
			component.MustNewID("traces"): {
				Receivers: []component.ID{component.NewID(component.MustNewType("test"))},
				Exporters: []component.ID{component.NewID(nopType)},
			},
			component.MustNewID("metrics"): {
				Receivers: []component.ID{component.NewID(component.MustNewType("test"))},
				Exporters: []component.ID{component.NewID(nopType)},
			},
		},
		Extensions: extensions.Config{component.NewID(component.MustNewType("watcher"))},
	}

	s, err := service.New(context.Background(), set, cfg)
	require.NoError(t, err)

	err = s.Start(context.Background())
	require.NoError(t, err)
	time.Sleep(15 * time.Second)
	err = s.Shutdown(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 5, len(eventsReceived))

	for instanceID, events := range eventsReceived {
		if instanceID.ComponentID() == component.NewID(component.MustNewType("test")) {
			for i, e := range events {
				if i == 0 {
					assert.Equal(t, componentstatus.StatusStarting, e.Status())
				}
				if i == 1 {
					assert.Equal(t, componentstatus.StatusRecoverableError, e.Status())
				}
				if i == 2 {
					assert.Equal(t, componentstatus.StatusOK, e.Status())
				}
				if i == 3 {
					assert.Equal(t, componentstatus.StatusStopping, e.Status())
				}
				if i == 4 {
					assert.Equal(t, componentstatus.StatusStopped, e.Status())
				}
				if i >= 5 {
					assert.Fail(t, "received too many events")
				}
			}
		}
	}
}

func newReceiverFactory() receiver.Factory {
	return receiver.NewFactory(
		component.MustNewType("test"),
		createDefaultReceiverConfig,
		receiver.WithTraces(createTraces, component.StabilityLevelStable),
		receiver.WithMetrics(createMetrics, component.StabilityLevelStable),
	)
}

type testReceiver struct{}

func (t *testReceiver) Start(_ context.Context, host component.Host) error {
	if statusReporter, ok := host.(componentstatus.Reporter); ok {
		statusReporter.Report(componentstatus.NewRecoverableErrorEvent(errors.New("test recoverable error")))
		go func() {
			statusReporter.Report(componentstatus.NewEvent(componentstatus.StatusOK))
		}()
	}
	return nil
}

func (t *testReceiver) Shutdown(_ context.Context) error {
	return nil
}

type receiverConfig struct{}

func createDefaultReceiverConfig() component.Config {
	return &receiverConfig{}
}

func createTraces(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	_ consumer.Traces,
) (receiver.Traces, error) {
	oCfg := cfg.(*receiverConfig)
	r, err := receivers.LoadOrStore(
		oCfg,
		func() (*testReceiver, error) {
			return &testReceiver{}, nil
		},
		&set.TelemetrySettings,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func createMetrics(
	_ context.Context,
	set receiver.Settings,
	cfg component.Config,
	_ consumer.Metrics,
) (receiver.Metrics, error) {
	oCfg := cfg.(*receiverConfig)
	r, err := receivers.LoadOrStore(
		oCfg,
		func() (*testReceiver, error) {
			return &testReceiver{}, nil
		},
		&set.TelemetrySettings,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

var receivers = sharedcomponent.NewMap[*receiverConfig, *testReceiver]()

func newExtensionFactory() extension.Factory {
	return extension.NewFactory(
		component.MustNewType("watcher"),
		createDefaultExtensionConfig,
		createExtension,
		component.StabilityLevelStable,
	)
}

func createExtension(_ context.Context, _ extension.Settings, cfg component.Config) (extension.Extension, error) {
	oCfg := cfg.(*extensionConfig)
	return &testExtension{
		eventsReceived: oCfg.eventsReceived,
	}, nil
}

type testExtension struct {
	eventsReceived map[*componentstatus.InstanceID][]*componentstatus.Event
}

type extensionConfig struct {
	eventsReceived map[*componentstatus.InstanceID][]*componentstatus.Event
}

func createDefaultExtensionConfig() component.Config {
	return &extensionConfig{}
}

// Start implements the component.Component interface.
func (t *testExtension) Start(_ context.Context, _ component.Host) error {
	return nil
}

// Shutdown implements the component.Component interface.
func (t *testExtension) Shutdown(_ context.Context) error {
	return nil
}

// ComponentStatusChanged implements the extension.StatusWatcher interface.
func (t *testExtension) ComponentStatusChanged(
	source *componentstatus.InstanceID,
	event *componentstatus.Event,
) {
	t.eventsReceived[source] = append(t.eventsReceived[source], event)
}

// NotifyConfig implements the extension.ConfigWatcher interface.
func (t *testExtension) NotifyConfig(_ context.Context, _ *confmap.Conf) error {
	return nil
}

// Ready implements the extension.PipelineWatcher interface.
func (t *testExtension) Ready() error {
	return nil
}

// NotReady implements the extension.PipelineWatcher interface.
func (t *testExtension) NotReady() error {
	return nil
}
