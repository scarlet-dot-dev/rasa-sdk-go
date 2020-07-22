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
	ec := (*eventCapture)(&events)

	// form context for Handler methods
	fctx := Context{actx}

	// first ensure the form is activated
	if _, err = ec.capture(a.activateIfRequired(&fctx, disp)); err != nil {
		return
	}

	// validate if needed
	if _, err = ec.capture(a.validateIfRequired(&fctx, disp)); err != nil {
		return
	}

	// if validation caused the form to be deactivated, abort
	if containsFormDeactivate(*ec) {
		return
	}

	// perform remaining actions with updated tracker
	a.applySlotSets(fctx.Tracker, events)

	// get the next slot request events
	var added bool
	if added, err = ec.capture(a.Handler.RequestNextSlot(&fctx, disp)); err != nil || !added {
		return
	}

	// no new events - submit the form AND deactivate
	if _, err = ec.capture(a.Handler.Submit(&fctx, disp)); err != nil {
		return
	}
	if _, err = ec.capture(a.Handler.Deactivate()); err != nil {
		return
	}

	// done
	return
}

//
type eventCapture []rasa.Event

//
func (c *eventCapture) capture(e []rasa.Event, er error) (added bool, err error) {
	if err = er; err == nil && len(e) > 0 {
		*c = append(*c)
		added = true
	}
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
