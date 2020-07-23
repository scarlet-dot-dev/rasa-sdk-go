// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"fmt"
	"net/http"
)

// respErr
type respErr interface {
	respCode() int
	respBody() string
}

// InvalidRequestError indicates that the request errored due to it being
// malformed or otherwise incorrect (such as invalid JSON).
type InvalidRequestError struct {
	Cause error
}

var _ error = (*InvalidRequestError)(nil)
var _ respErr = (*InvalidRequestError)(nil)

// Error implements builtin.error.
func (e *InvalidRequestError) Error() string {
	return e.Cause.Error()
}

// Unwrap implements errors.Unwrap.
func (e *InvalidRequestError) Unwrap() error {
	return e.Cause
}

//
func (e *InvalidRequestError) respCode() int {
	return http.StatusBadRequest
}

//
func (e *InvalidRequestError) respBody() string {
	return "invalid request"
}

// MissingHandlerError indicates
type MissingHandlerError struct {
	Action string
}

var _ error = (*MissingHandlerError)(nil)
var _ respErr = (*MissingHandlerError)(nil)

// Error implements builtin.error.
func (e *MissingHandlerError) Error() string {
	return fmt.Sprintf(
		"received request for action [%s] with no configured handler",
		e.Action,
	)
}

//
func (e *MissingHandlerError) respCode() int {
	return http.StatusInternalServerError
}

//
func (e *MissingHandlerError) respBody() string {
	return "invalid request"
}

// HandlerError is the type used to wrap errors occuring inside action handlers.
type HandlerError struct {
	Action string
	Cause  error
}

var _ error = (*HandlerError)(nil)
var _ respErr = (*HandlerError)(nil)

// Error implements builtin.error.
func (e *HandlerError) Error() string {
	return fmt.Sprintf(
		"Error occured when handling action [%s]: %s",
		e.Action,
		e.Cause.Error(),
	)
}

// Unwrap implements errors.Unwrap.
func (e *HandlerError) Unwrap() error {
	return e.Cause
}

//
func (e *HandlerError) respCode() int {
	return http.StatusInternalServerError
}

//
func (e *HandlerError) respBody() string {
	return "error handling the action"
}

// UnmarshalError indicates an error resulting from unmarshalling invalid JSON.
type UnmarshalError struct {
	cause error
}

// ensure interface
var _ error = (*UnmarshalError)(nil)
var _ respErr = (*UnmarshalError)(nil)

// Unwrap implements errors.Unwraper.
func (e *UnmarshalError) Unwrap() error {
	return e.cause
}

// Error implements error.
func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("invalid JSON: %s", e.cause.Error())
}

// respCode implements respErr.
func (e *UnmarshalError) respCode() int {
	return http.StatusBadRequest
}

// respBody implements respErr.
func (e *UnmarshalError) respBody() string {
	return fmt.Sprintf("invalid JSON received: %s", e.cause.Error())
}

// ExecutionRejection implements error for errors that should stop Rasa from
// executing an action.
type ExecutionRejection struct {
	Action string
	Reason string
}

// Error implements builtin.error.
func (e *ExecutionRejection) Error() string {
	return fmt.Sprintf(
		"rejected execution of [%s] for reason [%s]",
		e.Action,
		e.Reason,
	)
}
