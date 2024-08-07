// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package status // import "go.opentelemetry.io/collector/service/internal/status"

import (
	"errors"
	"fmt"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// onTransitionFunc receives a component.StatusEvent on a successful state transition
type onTransitionFunc func(*component.StatusEvent)

// errInvalidStateTransition is returned for invalid state transitions
var errInvalidStateTransition = errors.New("invalid state transition")

// fsm is a finite state machine that models transitions for component status
type fsm struct {
	current      *component.StatusEvent
	transitions  map[component.Status]map[component.Status]struct{}
	onTransition onTransitionFunc
}

// transition will attempt to execute a state transition. If it's successful, it calls the
// onTransitionFunc with a StatusEvent representing the new state. Returns an error if the arguments
// result in an invalid status, or if the state transition is not valid.
func (m *fsm) transition(ev *component.StatusEvent) error {
	if _, ok := m.transitions[m.current.Status()][ev.Status()]; !ok {
		return fmt.Errorf(
			"cannot transition from %s to %s: %w",
			m.current.Status(),
			ev.Status(),
			errInvalidStateTransition,
		)
	}
	m.current = ev
	m.onTransition(ev)
	return nil
}

// newFSM creates a state machine with all valid transitions for component.Status.
// The initial state is set to component.StatusNone.
func newFSM(onTransition onTransitionFunc) *fsm {
	return &fsm{
		current:      component.NewStatusEvent(component.StatusNone),
		onTransition: onTransition,
		transitions: map[component.Status]map[component.Status]struct{}{
			component.StatusNone: {
				component.StatusStarting: {},
			},
			component.StatusStarting: {
				component.StatusOK:               {},
				component.StatusRecoverableError: {},
				component.StatusPermanentError:   {},
				component.StatusFatalError:       {},
				component.StatusStopping:         {},
			},
			component.StatusOK: {
				component.StatusRecoverableError: {},
				component.StatusPermanentError:   {},
				component.StatusFatalError:       {},
				component.StatusStopping:         {},
			},
			component.StatusRecoverableError: {
				component.StatusOK:             {},
				component.StatusPermanentError: {},
				component.StatusFatalError:     {},
				component.StatusStopping:       {},
			},
			component.StatusPermanentError: {},
			component.StatusFatalError:     {},
			component.StatusStopping: {
				component.StatusRecoverableError: {},
				component.StatusPermanentError:   {},
				component.StatusFatalError:       {},
				component.StatusStopped:          {},
			},
			component.StatusStopped: {},
		},
	}
}

// NotifyStatusFunc is the receiver of status events after successful state transitions
type NotifyStatusFunc func(*component.InstanceID, *component.StatusEvent)

// InvalidTransitionFunc is the receiver of invalid transition errors
type InvalidTransitionFunc func(error)

// ServiceStatusFunc is the expected type of ReportStatus
type ServiceStatusFunc func(*component.InstanceID, *component.StatusEvent)

// ErrStatusNotReady is returned when trying to report status before service start
var ErrStatusNotReady = errors.New("report component status is not ready until service start")

// Reporter handles component status reporting
type Reporter interface {
	Ready()
	ReportStatus(id *component.InstanceID, ev *component.StatusEvent)
	ReportOKIfStarting(id *component.InstanceID)
}

type reporter struct {
	mu                  sync.Mutex
	ready               bool
	fsmMap              map[*component.InstanceID]*fsm
	onStatusChange      NotifyStatusFunc
	onInvalidTransition InvalidTransitionFunc
}

// NewReporter returns a reporter that will invoke the NotifyStatusFunc when a component's status
// has changed.
func NewReporter(onStatusChange NotifyStatusFunc, onInvalidTransition InvalidTransitionFunc) Reporter {
	return &reporter{
		fsmMap:              make(map[*component.InstanceID]*fsm),
		onStatusChange:      onStatusChange,
		onInvalidTransition: onInvalidTransition,
	}
}

// Ready enables status reporting
func (r *reporter) Ready() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ready = true
}

// ReportStatus reports status for the given InstanceID
func (r *reporter) ReportStatus(
	id *component.InstanceID,
	ev *component.StatusEvent,
) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.ready {
		r.onInvalidTransition(ErrStatusNotReady)
	} else {
		if err := r.componentFSM(id).transition(ev); err != nil {
			r.onInvalidTransition(err)
		}
	}
}

func (r *reporter) ReportOKIfStarting(id *component.InstanceID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.ready {
		r.onInvalidTransition(ErrStatusNotReady)
	}
	fsm := r.componentFSM(id)
	if fsm.current.Status() == component.StatusStarting {
		if err := fsm.transition(component.NewStatusEvent(component.StatusOK)); err != nil {
			r.onInvalidTransition(err)
		}
	}
}

// Note: a lock must be acquired before calling this method.
func (r *reporter) componentFSM(id *component.InstanceID) *fsm {
	fsm, ok := r.fsmMap[id]
	if !ok {
		fsm = newFSM(func(ev *component.StatusEvent) { r.onStatusChange(id, ev) })
		r.fsmMap[id] = fsm
	}
	return fsm
}

// NewReportStatusFunc returns a function to be used as ReportStatus for component.TelemetrySettings
func NewReportStatusFunc(
	id *component.InstanceID,
	srvStatus ServiceStatusFunc,
) func(*component.StatusEvent) {
	return func(ev *component.StatusEvent) {
		srvStatus(id, ev)
	}
}
