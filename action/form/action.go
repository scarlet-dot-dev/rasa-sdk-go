// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

import (
	"fmt"
	"time"

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
) (events rasa.Events, err error) {
	ec := (*eventCapture)(&events)

	// form context for Handler methods
	fctx := Context{actx, a.Handler}

	// first ensure the form is activated
	if _, err = ec.capture(a.activateIfRequired(&fctx, disp)); err != nil {
		return
	}

	// validate if needed
	if _, err = ec.capture(a.validateIfRequired(&fctx, disp)); err != nil {
		return
	}

	// if validation caused the form to be deactivated, abort
	if ec.containsLoopDeactivate() {
		return
	}

	// perform remaining actions with updated tracker
	ec.applySlotSets(fctx.Tracker)

	// get the next slot request events
	var added int
	if added, err = ec.capture(a.Handler.RequestNextSlot(&fctx, disp)); err != nil || added == 0 {
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

// String implements fmt.Stringer.
//
// String returns the name of the form.
func (a *Action) String() string {
	return fmt.Sprintf("FormAction(%s)", a.ActionName())
}

//
//
//

// Activate loop if the loop is called for the fist timne.
//
// When activating, required slots will be validated if they were filled in
// prior to the form's activation, the rasa.Form event will be returned with the
// name of the form, and any SlotSet events from validation of pre-filled slots
// will be returned.
func (a *Action) activateIfRequired(
	ctx *Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	ec := (*eventCapture)(&events)

	if ctx.Tracker.HasActiveForm() {
		ctx.Debugf("the loop [%s] is active", ctx.Tracker.ActiveLoop.Name)
	} else {
		ctx.Debugf("there is no active loop", ctx.Tracker.ActiveLoop.Name)
	}

	if ctx.Tracker.ActiveLoop.Is(a.Handler.FormName()) {
		// we are active - nothing to do
		return
	}

	// activate the form
	ec.append(rasa.ActiveLoop{
		Timestamp: rasa.Time(time.Now()),
		Name:      a.Handler.FormName(),
	})

	//
	prefilledSlots := make(rasa.Slots)
	requiredSlots := a.Handler.RequiredSlots(ctx)
	for _, slot := range requiredSlots {
		if !ctx.shouldRequestSlot(slot) {
			prefilledSlots[slot] = ctx.Tracker.Slots[slot]
		}
	}

	// if there are no prefilled slots, we are done
	if len(prefilledSlots) == 0 {
		ctx.Debugf("no pre-filled required slots to validate")
		return
	}

	// validate the prefilled slots
	ctx.Debugf("validating pre-filled required slots: %s", prefilledSlots)
	if _, err = ec.capture(ctx.ValidateSlots(disp, prefilledSlots)); err != nil {
		return
	}

	return
}

// validateIfRequired will perform validation on all existing slots that are
// required by the form, if required.
//
//  Validation is required if:
//	- the form is active
//	- the form is called after `action_listen`
//	- form validation was not cancelled
func (a *Action) validateIfRequired(ctx *Context, disp *action.CollectingDispatcher) (events rasa.Events, err error) {
	if !ctx.shouldValidate() {
		ctx.Debugf("Skipping validation")
		return
	}

	//
	ctx.Debugf("validating user input [%s]", ctx.Tracker.LatestMessage)
	events, err = a.Handler.Validate(ctx, disp)
	return
}

//
//
//

//
type eventCapture rasa.Events

// appemd
func (c *eventCapture) append(e ...rasa.Event) {
	*c = append(*c, e...)
}

//
func (c *eventCapture) capture(e rasa.Events, ce error) (added int, err error) {
	added, err = len(e), ce
	if err == nil && len(e) > 0 {
		c.append(e...)
	}
	return
}

// containsLoopDeactivate
func (c *eventCapture) containsLoopDeactivate() bool {
	for _, event := range *c {
		if e, ok := event.(*rasa.ActiveLoop); ok {
			if e.Name == "" {
				return true
			}
		}
	}
	return false
}

// applySlotSets
func (c *eventCapture) applySlotSets(t *rasa.Tracker) {
	for i := range *c {
		if event, ok := (*c)[i].(*rasa.SlotSet); ok {
			t.Slots[event.Key] = event.Value
		}
	}
}
