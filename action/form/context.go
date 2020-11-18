// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

import (
	"time"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// ValidatorContext TODO
type ValidatorContext struct {
	action.Context
	validator Validator
}

// ValidateSlots TODO
func (c *ValidatorContext) ValidateSlots(
	disp *action.CollectingDispatcher,
	slots rasa.Slots,
) (events rasa.Events, err error) {
	sc := make(rasa.Slots)

	for slot, value := range slots {
		validator := c.validator.Validator(slot)

		var sset rasa.Slots
		if sset, err = validator.Validate(c, disp, value); err != nil {
			return
		}
		sc.Update(sset)
	}

	// turn validated slots into SlotSet events
	timestamp := rasa.Time(time.Now())
	for key := range sc {
		events = append(events, rasa.SlotSet{
			Timestamp: timestamp,
			Key:       key,
			Value:     sc[key],
		})
	}
	return
}

// RequestSlot returns a rasa event to mark the provided slot as a required
// slot.
func (c *ValidatorContext) RequestSlot(slot string) rasa.Event {
	return &rasa.SlotSet{
		Timestamp: c.Now(),
		Key:       RequestedSlot,
		Value:     slot,
	}
}
