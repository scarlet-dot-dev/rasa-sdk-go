package form

import (
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// Validator specifies the interface for a custom Slot Validator.
type Validator interface {
	// Validate will receive a slot value to perform validation on.
	//
	// The returned slots map may contain any set of (slot, value) mappings to
	// pass back to Rasa. It is up to the developer to ensure these slots are
	// all valid.
	Validate(
		ctx *Context,
		dispatcher *action.CollectingDispatcher,
		value interface{},
	) (slots rasa.SlotMap, err error)
}

// ValidatorFunc is a wrapper type for functions that may be used as Validator
// implementations.
type ValidatorFunc func(
	ctx *Context,
	dispatcher *action.CollectingDispatcher,
	value interface{},
) (rasa.SlotMap, error)

// ensure interface
var _ Validator = (ValidatorFunc)(nil)

// Validate implements Validator.
func (fn ValidatorFunc) Validate(ctx *Context, dispatcher *action.CollectingDispatcher, value interface{}) (rasa.SlotMap, error) {
	return fn(ctx, dispatcher, value)
}
