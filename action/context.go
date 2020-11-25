// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"

	"go.scarlet.dev/rasa"
)

// Context contains request specific context, such as a logger for the specific
// request with tracing information.
type Context interface {
	// Context returns the Context of the action request that triggered the
	// execution of the action handler.
	Context() context.Context

	// Tracker returns the tracker state associated with the Context.
	Tracker() *rasa.Tracker

	// Domain returns the domain specification associated with the Context.
	Domain() *rasa.Domain

	// Logger returns the logger associated with the context. Logger should
	// never return nil. If no logger is available, a non-nil no-op logger
	// should be returned.
	Logger() Logger
}

// contextImpl implements the Context interface for the
type contextImpl struct {

	//
	tracker *rasa.Tracker

	//
	domain *rasa.Domain

	// internal fields
	logger  Logger
	context context.Context
}

// ensure interfaces.
var _ Logger = (*contextImpl)(nil)
var _ Context = (*contextImpl)(nil)

// Logger implements Context.
//
// Logger will return the contextImpl, which wraps the logger field with nil
// checks.
func (c *contextImpl) Logger() Logger {
	return c
}

// Context implements Context.
func (c *contextImpl) Context() context.Context {
	return c.context
}

// Tracker implements Context.
func (c *contextImpl) Tracker() *rasa.Tracker {
	return c.tracker
}

// Domain implements Context.
func (c *contextImpl) Domain() *rasa.Domain {
	return c.domain
}

// Debugf impements Logger.
func (c *contextImpl) Debugf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Debugf(format, args...)
	}
}

// Infof implements Logger.
func (c *contextImpl) Infof(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(format, args...)
	}
}

// Warnf implements Logger.
func (c *contextImpl) Warnf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(format, args...)
	}
}

// Errorf implements Logger.
func (c *contextImpl) Errorf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Errorf(format, args...)
	}
}

// SetSlot is a constructor for the rasa.SlotSet event. It returns a rasa.Event
// that sets the slot to the provided value.
//
// SetSlot is a simple alias function, equivalent to constructing
// rasa.SlotSet{Key: slot, Value: val}
func SetSlot(slot string, val interface{}) rasa.SlotSet {
	return rasa.SlotSet{
		Key:   slot,
		Value: val,
	}
}

// ResetSlot is a constructor for the rasa.SlotSet event. It returns a
// rasa.Event that resets the slot to nil (`None`).
func ResetSlot(slot string) rasa.SlotSet {
	return rasa.SlotSet{
		Key:   slot,
		Value: nil,
	}
}
