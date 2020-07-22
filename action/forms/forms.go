package form

import (
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// RequestedSlot is used to store information needed to do the form handling.
const (
	RequestedSlot = "requested_slot"
)

// Handler TODO
type Handler interface {
	// Name is the unique identifier of the form. This name is the action that
	// triggers the form.
	FormName() string

	// RequestedSlots is a list of required slots that the form has to fill.
	//
	// The Tracker can be used to request a different set of slots depending on
	// the state of the dialogue.
	RequiredSlots(ctx *Context) []string

	// SlotMappers should return a map of (slot, mapper) values to map required
	// slots.
	//
	// Returning a nil or empty map is converted to a mapping of the slot to the
	// extracted entity with the same name.
	SlotMappers() Mappers

	// Validator returns a custom validator for the slot. If no custom
	// validation is needed, this method should return nil.
	Validator(slot string) Validator

	// Submit should process the action resulting from the Form being finalized.
	//
	// The Context passed to Run is the context of the HTTP request sent by
	// Rasa's engine.
	Submit(
		ctx *Context,
		dispatcher *action.CollectingDispatcher,
	) (events []rasa.Event, err error)
}

// SlotRequester is an OPTIONAL interface for Handler. When implemented, it will
// be called instead of the standard implementation of FormHandler's
// RequestNextSlot.
type SlotRequester interface {
	// RequestNextSlot TODO
	RequestNextSlot(ctx *Context, dispatcher *action.CollectingDispatcher) ([]rasa.Event, error)
}

// Deactivater is an OPTIONAL interface for Handler. When implemented, it will
// be called instead of the standard implementation of form.Action's deactivate
// implementation.
//
// Context provides a default implementation for this interface called
// DeactivateForm.
type Deactivater interface {
	Deactivate(ctx *Context, dispatcher *action.CollectingDispatcher) ([]rasa.Event, error)
}
