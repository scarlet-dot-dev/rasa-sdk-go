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
type Context struct {

	//
	Tracker *rasa.Tracker

	//
	Domain *rasa.Domain

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
