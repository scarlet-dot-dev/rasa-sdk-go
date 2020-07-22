package form

import (
	"time"

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

	// Submit should process the action resulting from the Form being finalized.
	//
	// The Context passed to Run is the context of the HTTP request sent by
	// Rasa's engine.
	Submit(
		ctx *Context,
		dispatcher *action.CollectingDispatcher,
	) (events rasa.Events, err error)

	// RequestNextSlot TODO
	//
	// A default implementation of RequestNextSlot can be provided by embedding
	// `HandlerDefaults`.
	RequestNextSlot(ctx *Context, dispatcher *action.CollectingDispatcher) (rasa.Events, error)

	// SlotMappers should return a map of (slot, mapper) values to map required
	// slots.
	//
	// Returning a nil or empty map is converted to a mapping of the slot to the
	// extracted entity with the same name.
	//
	// A default implementation of SlotMappers can be provided by embedding
	// `HandlerDefaults`.
	SlotMappers() Mappers

	// Validator returns a custom validator for the slot. If no custom
	// validation is needed, this method should return nil.
	//
	// A default implementation of Validator can be provided by embedding
	// `HandlerDefaults`.
	Validator(slot string) Validator

	// Deactivate TODO
	//
	// A default implementation of Deactivate can be provided by embedding
	// `HandlerDefaults`.
	Deactivate() (events rasa.Events, err error)
}

// SlotRequester is an OPTIONAL interface for Handler. When implemented, it will
// be called instead of the standard implementation of FormHandler's
// RequestNextSlot.
type SlotRequester interface {
	// RequestNextSlot TODO
	RequestNextSlot(ctx *Context, dispatcher *action.CollectingDispatcher) (rasa.Events, error)
}

// DefaultEmbed provides default implementations for optional methods of the
// Handler interface.
type DefaultEmbed struct{}

// RequestNextSlot implements Handler.
//
//
func (DefaultEmbed) RequestNextSlot(
	ctx *Context,
	dispatcher *action.CollectingDispatcher,
) (rasa.Events, error) {
	return ctx.requestNextSlot(dispatcher)
}

// SlotMappers should return a map of (slot, mapper) values to map required
// slots.
//
// Returning a nil or empty map is converted to a mapping of the slot to the
// extracted entity with the same name.
func (DefaultEmbed) SlotMappers() Mappers {
	return nil
}

// Validator returns a custom validator for the slot. If no custom
// validation is needed, this method should return nil.
func (DefaultEmbed) Validator(slot string) Validator {
	return nil
}

// Deactivate TODO
func (DefaultEmbed) Deactivate() (events rasa.Events, err error) {
	events = rasa.Events{
		rasa.Form{
			Timestamp: rasa.Time(time.Now()),
			Name:      "",
		},
		rasa.SlotSet{
			Timestamp: rasa.Time(time.Now()),
			Key:       RequestedSlot,
			Value:     nil,
		},
	}
	return
}
