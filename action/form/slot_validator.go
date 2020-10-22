// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

import (
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// SlotValidator specifies the interface for a custom Slot SlotValidator.
type SlotValidator interface {
	// Validate will receive a slot value to perform validation on.
	//
	// The returned slots map may contain any set of (slot, value) mappings to
	// pass back to Rasa. It is up to the developer to ensure these slots are
	// all valid.
	Validate(
		ctx *ValidatorContext,
		dispatcher *action.CollectingDispatcher,
		value interface{},
	) (slots rasa.Slots, err error)
}

// ensure interface
var _ SlotValidator = (SlotValidatorFunc)(nil)
var _ SlotValidator = (DefaultValidator)("")

// SlotValidatorFunc is a wrapper type for functions that may be used as Validator
// implementations.
type SlotValidatorFunc func(
	ctx *ValidatorContext,
	dispatcher *action.CollectingDispatcher,
	value interface{},
) (rasa.Slots, error)

// Validate implements Validator.
func (fn SlotValidatorFunc) Validate(ctx *ValidatorContext, dispatcher *action.CollectingDispatcher, value interface{}) (rasa.Slots, error) {
	return fn(ctx, dispatcher, value)
}

// DefaultValidator implements the Validator interface for simple cases where
// custom validation is not required.
type DefaultValidator string

// Validate implements Validator.
func (v DefaultValidator) Validate(
	ctx *ValidatorContext,
	dispatcher *action.CollectingDispatcher,
	value interface{},
) (slots rasa.Slots, err error) {
	slots = rasa.Slots{
		string(v): value,
	}
	return
}
