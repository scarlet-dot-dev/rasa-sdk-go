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

// AskForSlotAction implements action.Handler
type AskForSlotAction struct {
	Asker SlotAsker
}

// ensure interface
var _ action.Handler = (*AskForSlotAction)(nil)

// ActionName implements action.Handler
func (a *AskForSlotAction) ActionName() string {
	form, slot := a.Asker.Form(), a.Asker.Slot()
	if form == "" {
		return fmt.Sprintf("action_ask_%s", slot)
	}
	return fmt.Sprintf("action_ask_%s__%s", form, slot)
}

// Run implements action.Handler.
func (a *AskForSlotAction) Run(
	ctx action.Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	return a.Asker.Ask(ctx, disp)
}

// SlotAsker TODO
type SlotAsker interface {
	// Form returns the for which this slot asker is intended.
	//
	// If Form returns nothing, the asker is considered global for the slot.
	Form() string

	// Slot returns the slot for which this slot asked is intender.
	Slot() string

	// Ask should execute the action handling.
	Ask(ctx action.Context, disp *action.CollectingDispatcher) (events rasa.Events, err error)
}
