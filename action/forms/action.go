package form

import (
	"fmt"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// Action TODO
type Action struct {
	Handler Handler
}

// ensure interface
var _ action.Handler = (*Action)(nil)

// ActionName implements action.Handler.
func (a *Action) ActionName() string {
	return a.Handler.FormName()
}

// Run implements action.Handler.
//
// Run executes the side effects of this form.
//  Steps:
//  - activate if needed
//  - validate user input if needed
//  - set validated slots
//  - utter_ask_{slot} template with the next required slot
//  - submit the form if all required slots are set
//  - deactivate the form
func (a *Action) Run(
	actx *action.Context,
	disp *action.CollectingDispatcher,
) (events []rasa.Event, err error) {
	var add []rasa.Event

	// form context for Handler methods
	fctx := Context{actx}

	// first ensure the form is activated
	if add, err = a.activateIfRequired(&fctx, disp); err != nil {
		return
	} else if len(add) > 0 {
		events = append(events, add...)
	}

	// validate if needed
	if add, err = a.validateIfRequired(&fctx, disp); err != nil {
		return
	} else if len(add) > 0 {
		events = append(events, add...)
	}

	// if validation caused the form to be deactivated, abort
	if containsFormDeactivate(events) {
		return
	}

	// perform remaining actions with updated tracker
	a.applySlotSets(fctx.Tracker, events)

	// get the next slot request events
	if add, err = a.requestNextSlot(&fctx, disp); err != nil {
		return
	} else if len(add) > 0 {
		events = append(events, add...)
		return
	}

	// no new events - submit the form

	//

	return
}

// String implements fmt.Stringer.
//
// String returns the name of the form.
func (a *Action) String() string {
	return fmt.Sprintf("FormAction(%s)", a.ActionName())
}

//
//
//

// Activate form if the form is called for the fist timne.
//
// When activating, required slots will be validated if they were filled in
// prior to the form's activation, the rasa.Form event will be returned with the
// name of the form, and any SlotSet events from validation of pre-filled slots
// will be returned.
func (a *Action) activateIfRequired(
	ctx *Context,
	dispatcher *action.CollectingDispatcher,
) (events []rasa.Event, err error) {
	// logger := action.
	// UNIMPLEMENTED
	return
}

// validateIfRequired will perform validation on all existing slots that are required by the form.
func (a *Action) validateIfRequired(ctx *Context, disp *action.CollectingDispatcher) (events []rasa.Event, err error) {
	// UNIMPLEMENTED
	return nil, nil
}

//
func (a *Action) requestNextSlot(
	ctx *Context,
	disp *action.CollectingDispatcher,
) (events []rasa.Event, err error) {
	// if implemented by the handler, defer call
	if nsr, ok := a.Handler.(SlotRequester); ok {
		events, err = nsr.RequestNextSlot(ctx, disp)
		return
	}

	// manual call
	return
}

//
func (a *Action) applySlotSets(t *rasa.Tracker, events []rasa.Event) {
	for i := range events {
		if event, ok := events[i].(*rasa.SlotSet); ok {
			t.Slots[event.Key] = event.Value
		}
	}
}

//
func containsFormDeactivate(events []rasa.Event) bool {
	for _, event := range events {
		if e, ok := event.(*rasa.Form); ok {
			if e.Name == "" {
				return true
			}
		}
	}
	return false
}
