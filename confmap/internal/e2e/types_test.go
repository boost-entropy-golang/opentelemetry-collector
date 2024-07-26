// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package e2etest

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/internal/globalgates"
)

type TargetField string

const (
	TargetFieldInt          TargetField = "int_field"
	TargetFieldString       TargetField = "string_field"
	TargetFieldBool         TargetField = "bool_field"
	TargetFieldInlineString TargetField = "inline_string_field"
)

type Test struct {
	value        string
	targetField  TargetField
	expected     any
	resolveErr   string
	unmarshalErr string
}

type TargetConfig[T any] struct {
	Field T `mapstructure:"field"`
}

func NewResolver(t testing.TB, path string) *confmap.Resolver {
	resolver, err := confmap.NewResolver(confmap.ResolverSettings{
		URIs: []string{filepath.Join("testdata", path)},
		ProviderFactories: []confmap.ProviderFactory{
			fileprovider.NewFactory(),
			envprovider.NewFactory(),
		},
	})
	require.NoError(t, err)
	return resolver
}

func AssertExpectedMatch[T any](t *testing.T, tt Test, conf *confmap.Conf, cfg *TargetConfig[T]) {
	err := conf.Unmarshal(cfg)
	if tt.unmarshalErr != "" {
		require.ErrorContains(t, err, tt.unmarshalErr)
		return
	}
	require.NoError(t, err)
	require.Equal(t, tt.expected, cfg.Field)
}

func AssertResolvesTo(t *testing.T, resolver *confmap.Resolver, tt Test) {
	conf, err := resolver.Resolve(context.Background())
	if tt.resolveErr != "" {
		require.ErrorContains(t, err, tt.resolveErr)
		return
	}
	require.NoError(t, err)

	switch tt.targetField {
	case TargetFieldInt:
		var cfg TargetConfig[int]
		AssertExpectedMatch(t, tt, conf, &cfg)
	case TargetFieldString, TargetFieldInlineString:
		var cfg TargetConfig[string]
		AssertExpectedMatch(t, tt, conf, &cfg)
	case TargetFieldBool:
		var cfg TargetConfig[bool]
		AssertExpectedMatch(t, tt, conf, &cfg)
	default:
		t.Fatalf("unexpected target field %q", tt.targetField)
	}
}

func TestTypeCasting(t *testing.T) {
	values := []Test{
		{
			value:       "123",
			targetField: TargetFieldInt,
			expected:    123,
		},
		{
			value:       "123",
			targetField: TargetFieldString,
			expected:    "123",
		},
		{
			value:       "123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 123 expansion",
		},
		{
			value:       "0123",
			targetField: TargetFieldInt,
			expected:    83,
		},
		{
			value:       "0123",
			targetField: TargetFieldString,
			expected:    "83",
		},
		{
			value:       "0123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 83 expansion",
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldInt,
			expected:    3735928559,
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldString,
			expected:    "3735928559",
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 3735928559 expansion",
		},
		{
			value:       "\"0123\"",
			targetField: TargetFieldString,
			expected:    "0123",
		},
		{
			value:       "\"0123\"",
			targetField: TargetFieldInt,
			expected:    83,
		},
		{
			value:       "\"0123\"",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 0123 expansion",
		},
		{
			value:       "!!str 0123",
			targetField: TargetFieldString,
			expected:    "0123",
		},
		{
			value:       "!!str 0123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 0123 expansion",
		},
		{
			value:       "'!!str 0123'",
			targetField: TargetFieldString,
			expected:    "!!str 0123",
		},
		{
			value:       "\"!!str 0123\"",
			targetField: TargetFieldInlineString,
			expected:    "inline field with !!str 0123 expansion",
		},
		{
			value:       "''",
			targetField: TargetFieldString,
			expected:    "",
		},
		{
			value:       "\"\"",
			targetField: TargetFieldInlineString,
			expected:    "inline field with  expansion",
		},
		{
			value:       "t",
			targetField: TargetFieldBool,
			expected:    true,
		},
		{
			value:       "23",
			targetField: TargetFieldBool,
			expected:    true,
		},
		{
			value:       "foo\nbar",
			targetField: TargetFieldString,
			expected:    "foo bar",
		},
		{
			value:       "foo\nbar",
			targetField: TargetFieldInlineString,
			expected:    "inline field with foo bar expansion",
		},
		{
			value:       "\"1111:1111:1111:1111:1111::\"",
			targetField: TargetFieldString,
			expected:    "1111:1111:1111:1111:1111::",
		},
		{
			value:       "\"1111:1111:1111:1111:1111::\"",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 1111:1111:1111:1111:1111:: expansion",
		},
		{
			value:       "2006-01-02T15:04:05Z07:00",
			targetField: TargetFieldString,
			expected:    "2006-01-02T15:04:05Z07:00",
		},
		{
			value:       "2006-01-02T15:04:05Z07:00",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 2006-01-02T15:04:05Z07:00 expansion",
		},
	}

	previousValue := globalgates.StrictlyTypedInputGate.IsEnabled()
	err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, false)
	require.NoError(t, err)
	defer func() {
		err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, previousValue)
		require.NoError(t, err)
	}()

	for _, tt := range values {
		t.Run(tt.value+"/"+string(tt.targetField), func(t *testing.T) {
			testFile := "types_expand.yaml"
			if tt.targetField == TargetFieldInlineString {
				testFile = "types_expand_inline.yaml"
			}
			resolver := NewResolver(t, testFile)
			t.Setenv("ENV", tt.value)
			AssertResolvesTo(t, resolver, tt)
		})
	}
}

func TestStrictTypeCasting(t *testing.T) {
	values := []Test{
		{
			value:       "123",
			targetField: TargetFieldInt,
			expected:    123,
		},
		{
			value:       "123",
			targetField: TargetFieldString,
			expected:    "123",
		},
		{
			value:       "123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 123 expansion",
		},
		{
			value:       "0123",
			targetField: TargetFieldInt,
			expected:    83,
		},
		{
			value:       "0123",
			targetField: TargetFieldString,
			expected:    "0123",
		},
		{
			value:       "0123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 0123 expansion",
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldInt,
			expected:    3735928559,
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldString,
			expected:    "0xdeadbeef",
		},
		{
			value:       "0xdeadbeef",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 0xdeadbeef expansion",
		},
		{
			value:       "\"0123\"",
			targetField: TargetFieldString,
			expected:    "\"0123\"",
		},
		{
			value:        "\"0123\"",
			targetField:  TargetFieldInt,
			unmarshalErr: "'field' expected type 'int', got unconvertible type 'string', value: '\"0123\"'",
		},
		{
			value:       "\"0123\"",
			targetField: TargetFieldInlineString,
			expected:    "inline field with \"0123\" expansion",
		},
		{
			value:       "!!str 0123",
			targetField: TargetFieldString,
			expected:    "!!str 0123",
		},
		{
			value:       "!!str 0123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with !!str 0123 expansion",
		},
		{
			value:        "t",
			targetField:  TargetFieldBool,
			unmarshalErr: "'field' expected type 'bool', got unconvertible type 'string', value: 't'",
		},
		{
			value:        "23",
			targetField:  TargetFieldBool,
			unmarshalErr: "'field' expected type 'bool', got unconvertible type 'int', value: '23'",
		},
		{
			value:       "{\"field\": 123}",
			targetField: TargetFieldInlineString,
			expected:    "inline field with {\"field\": 123} expansion",
		},
		{
			value:       "1111:1111:1111:1111:1111::",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 1111:1111:1111:1111:1111:: expansion",
		},
		{
			value:       "1111:1111:1111:1111:1111::",
			targetField: TargetFieldString,
			expected:    "1111:1111:1111:1111:1111::",
		},
		{
			value:       "2006-01-02T15:04:05Z07:00",
			targetField: TargetFieldString,
			expected:    "2006-01-02T15:04:05Z07:00",
		},
		{
			value:       "2006-01-02T15:04:05Z07:00",
			targetField: TargetFieldInlineString,
			expected:    "inline field with 2006-01-02T15:04:05Z07:00 expansion",
		},
	}

	previousValue := globalgates.StrictlyTypedInputGate.IsEnabled()
	err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, true)
	require.NoError(t, err)
	defer func() {
		err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, previousValue)
		require.NoError(t, err)
	}()

	for _, tt := range values {
		t.Run(tt.value+"/"+string(tt.targetField)+"/"+"direct", func(t *testing.T) {
			testFile := "types_expand.yaml"
			if tt.targetField == TargetFieldInlineString {
				testFile = "types_expand_inline.yaml"
			}
			resolver := NewResolver(t, testFile)
			t.Setenv("ENV", tt.value)
			AssertResolvesTo(t, resolver, tt)
		})

		t.Run(tt.value+"/"+string(tt.targetField)+"/"+"indirect", func(t *testing.T) {
			testFile := "types_expand.yaml"
			if tt.targetField == TargetFieldInlineString {
				testFile = "types_expand_inline.yaml"
			}

			resolver := NewResolver(t, testFile)
			t.Setenv("ENV", "${env:ENV2}")
			t.Setenv("ENV2", tt.value)
			AssertResolvesTo(t, resolver, tt)
		})
	}
}

func TestRecursiveInlineString(t *testing.T) {
	values := []Test{
		{
			value:       "123",
			targetField: TargetFieldString,
			expected:    "The value The value 123 is wrapped is wrapped",
		},
		{
			value:       "123",
			targetField: TargetFieldInlineString,
			expected:    "inline field with The value The value 123 is wrapped is wrapped expansion",
		},
		{
			value:       "opentelemetry",
			targetField: TargetFieldString,
			expected:    "The value The value opentelemetry is wrapped is wrapped",
		},
		{
			value:       "opentelemetry",
			targetField: TargetFieldInlineString,
			expected:    "inline field with The value The value opentelemetry is wrapped is wrapped expansion",
		},
	}

	previousValue := globalgates.StrictlyTypedInputGate.IsEnabled()
	err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, true)
	require.NoError(t, err)
	defer func() {
		err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, previousValue)
		require.NoError(t, err)
	}()

	for _, tt := range values {
		t.Run(tt.value+"/"+string(tt.targetField), func(t *testing.T) {
			testFile := "types_expand.yaml"
			if tt.targetField == TargetFieldInlineString {
				testFile = "types_expand_inline.yaml"
			}

			resolver := NewResolver(t, testFile)
			t.Setenv("ENV", "The value ${env:ENV2} is wrapped")
			t.Setenv("ENV2", "The value ${env:ENV3} is wrapped")
			t.Setenv("ENV3", tt.value)
			AssertResolvesTo(t, resolver, tt)
		})
	}
}

func TestRecursiveMaps(t *testing.T) {
	value := "{value: 123}"

	previousValue := globalgates.StrictlyTypedInputGate.IsEnabled()
	err := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, true)
	require.NoError(t, err)
	defer func() {
		seterr := featuregate.GlobalRegistry().Set(globalgates.StrictlyTypedInputID, previousValue)
		require.NoError(t, seterr)
	}()

	resolver := NewResolver(t, "types_expand.yaml")
	t.Setenv("ENV", `{env: "${env:ENV2}", inline: "inline ${env:ENV2}"}`)
	t.Setenv("ENV2", `{env2: "${env:ENV3}"}`)
	t.Setenv("ENV3", value)
	conf, err := resolver.Resolve(context.Background())
	require.NoError(t, err)

	type Value struct {
		Value int `mapstructure:"value"`
	}
	type ENV2 struct {
		Env2 Value `mapstructure:"env2"`
	}
	type ENV struct {
		Env    ENV2   `mapstructure:"env"`
		Inline string `mapstructure:"inline"`
	}
	type Target struct {
		Field ENV `mapstructure:"field"`
	}

	var cfg Target
	err = conf.Unmarshal(&cfg)
	require.NoError(t, err)
	require.Equal(t,
		Target{Field: ENV{
			Env: ENV2{
				Env2: Value{
					Value: 123,
				}},
			Inline: "inline {env2: \"{value: 123}\"}",
		}},
		cfg,
	)

	confStr, err := resolver.Resolve(context.Background())
	require.NoError(t, err)
	var cfgStr TargetConfig[string]
	err = confStr.Unmarshal(&cfgStr)
	require.NoError(t, err)
	require.Equal(t, `{env: "{env2: "{value: 123}"}", inline: "inline {env2: "{value: 123}"}"}`,
		cfgStr.Field,
	)
}
