// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

import (
	"fmt"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

const (
	// RequestedSlot is used to store information needed to do the form handling.
	RequestedSlot = "requested_slot"

	// LoopInterruptedKey is used to detect if the loop is interrupted
	LoopInterruptedKey = "is_interrupted"
)

// ValidatorAction implements a base action for slot validation used by forms.
//
// It is based on FormValidationAction.
type ValidatorAction struct {
	Validator Validator
}

// ensure interface
var _ action.Handler = (*ValidatorAction)(nil)

// ActionName implements action.Handler.
func (a *ValidatorAction) ActionName() string {
	return fmt.Sprintf("validate_%s", a.Validator.Form())
}

// Run implements action.Handler.
func (a *ValidatorAction) Run(
	actx *action.Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	// form context for Handler methods
	ctx := ValidatorContext{actx, a.Validator}
	events, err = a.Validator.Validate(&ctx, disp)
	return
}

// String implements fmt.Stringer.
//
// String returns the name of the form.
func (a *ValidatorAction) String() string {
	return fmt.Sprintf("FormValidationAction(%s)", a.ActionName())
}

// Validator TODO
type Validator interface {
	// Form returns the name of the form this validator is intended for.
	Form() string

	// Validate extracts and validates the value of requested slot.
	//
	// If nothing was extracted, reject execution of the form action.
	//
	// A default implementation of Validate can be provided by embedding
	// `HandlerDefaults`.
	Validate(
		ctx *ValidatorContext,
		disp *action.CollectingDispatcher,
	) (events rasa.Events, err error)

	// Validator returns the validator for the provided slot.
	Validator(slot string) SlotValidator
}

// ValidatorEmbed provides default implementations for optional methods of
// the ValidationHandler interface.
type ValidatorEmbed struct{}

// Validate implements Validator.
func (ValidatorEmbed) Validate(ctx *ValidatorContext, disp *action.CollectingDispatcher) (events rasa.Events, err error) {
	return ctx.ValidateSlots(disp, ctx.Tracker.SlotsToValidate())
}
