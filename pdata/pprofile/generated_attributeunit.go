// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Code generated by "pdata/internal/cmd/pdatagen/main.go". DO NOT EDIT.
// To regenerate this file run "make genpdata".

package pprofile

import (
	"go.opentelemetry.io/collector/pdata/internal"
	otlpprofiles "go.opentelemetry.io/collector/pdata/internal/data/protogen/profiles/v1development"
)

// AttributeUnit Represents a mapping between Attribute Keys and Units.
//
// This is a reference type, if passed by value and callee modifies it the
// caller will see the modification.
//
// Must use NewAttributeUnit function to create new instances.
// Important: zero-initialized instance is not valid for use.
type AttributeUnit struct {
	orig  *otlpprofiles.AttributeUnit
	state *internal.State
}

func newAttributeUnit(orig *otlpprofiles.AttributeUnit, state *internal.State) AttributeUnit {
	return AttributeUnit{orig: orig, state: state}
}

// NewAttributeUnit creates a new empty AttributeUnit.
//
// This must be used only in testing code. Users should use "AppendEmpty" when part of a Slice,
// OR directly access the member if this is embedded in another struct.
func NewAttributeUnit() AttributeUnit {
	state := internal.StateMutable
	return newAttributeUnit(&otlpprofiles.AttributeUnit{}, &state)
}

// MoveTo moves all properties from the current struct overriding the destination and
// resetting the current instance to its zero value
func (ms AttributeUnit) MoveTo(dest AttributeUnit) {
	ms.state.AssertMutable()
	dest.state.AssertMutable()
	*dest.orig = *ms.orig
	*ms.orig = otlpprofiles.AttributeUnit{}
}

// AttributeKeyStrindex returns the attributekeystrindex associated with this AttributeUnit.
func (ms AttributeUnit) AttributeKeyStrindex() int32 {
	return ms.orig.AttributeKeyStrindex
}

// SetAttributeKeyStrindex replaces the attributekeystrindex associated with this AttributeUnit.
func (ms AttributeUnit) SetAttributeKeyStrindex(v int32) {
	ms.state.AssertMutable()
	ms.orig.AttributeKeyStrindex = v
}

// UnitStrindex returns the unitstrindex associated with this AttributeUnit.
func (ms AttributeUnit) UnitStrindex() int32 {
	return ms.orig.UnitStrindex
}

// SetUnitStrindex replaces the unitstrindex associated with this AttributeUnit.
func (ms AttributeUnit) SetUnitStrindex(v int32) {
	ms.state.AssertMutable()
	ms.orig.UnitStrindex = v
}

// CopyTo copies all properties from the current struct overriding the destination.
func (ms AttributeUnit) CopyTo(dest AttributeUnit) {
	dest.state.AssertMutable()
	dest.SetAttributeKeyStrindex(ms.AttributeKeyStrindex())
	dest.SetUnitStrindex(ms.UnitStrindex())
}
