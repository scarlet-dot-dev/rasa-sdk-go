// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"reflect"
	"time"

	"go.scarlet.dev/rasa"
)

// Context contains request specific context, such as a logger for the specific
// request with tracing information.
type Context struct {

	//
	Tracker *rasa.Tracker

	//
	Domain *rasa.Domain

	//
	nowfn func() time.Time

	// internal fields
	logger  Logger
	context context.Context
}

// ensure interfaces.
var _ Logger = (*Context)(nil)

// Context returns the Context of the action request that triggered the
// execution of the action handler.
func (c *Context) Context() context.Context {
	return c.context
}

// Debugf impements Logger.
func (c *Context) Debugf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Debugf(format, args...)
	}
}

// Infof implements Logger.
func (c *Context) Infof(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Infof(format, args...)
	}
}

// Warnf implements Logger.
func (c *Context) Warnf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Warnf(format, args...)
	}
}

// Errorf implements Logger.
func (c *Context) Errorf(format string, args ...interface{}) {
	if c.logger != nil {
		c.logger.Errorf(format, args...)
	}
}

// SetSlot TODO
func (c *Context) SetSlot(name string, val interface{}) rasa.SlotSet {
	return rasa.SlotSet{
		Timestamp: rasa.Time(time.Now()),
		Key:       name,
		Value:     val,
	}
}

// ResetSlot TODO
func (c *Context) ResetSlot(name string) rasa.SlotSet {
	return rasa.SlotSet{
		Timestamp: rasa.Time(time.Now()),
		Key:       name,
		Value:     nil,
	}
}

// Now returns the current time in the wrapper type rasa.Time.
func (c *Context) Now() rasa.Time {
	if c.nowfn != nil {
		return rasa.Time(c.nowfn())
	}
	return rasa.Time(time.Now())
}

// EntityValue TODO
func (c *Context) EntityValue(
	name, role, group string,
) (value rasa.EntityValue) {
	raw := c.EntityValues(name, role, group)
	if len(raw) == 1 {
		value = rasa.StringValue(raw[0])
		return
	}
	value = rasa.SliceValue(raw)
	return
}

// EntityValues TODO
func (c *Context) EntityValues(
	entity, role, group string,
) (values []string) {
	values = c.Tracker.LatestEntityValues(entity, role, group)
	return
}

// Slot TODO
func (c *Context) Slot(name string) (val interface{}, ok bool) {
	val, ok = c.Tracker.Slots[name]
	return
}

// SlotAs TODO
func (c *Context) SlotAs(name string, dst interface{}) (ok bool) {
	val, exists := c.Slot(name)
	if !exists {
		return
	}

	rdst := reflect.ValueOf(dst)
	rval := reflect.ValueOf(val)
	if ok = rdst.CanSet() && rval.Type().AssignableTo(rdst.Type()); !ok {
		return
	}

	rdst.Set(reflect.ValueOf(val))
	return
}
