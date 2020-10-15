// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

//
type testHandler struct {
	DefaultEmbed

	_FormName      func() string
	_RequiredSlots func(*Context) []string
	_Submit        func(*Context, *action.CollectingDispatcher) (rasa.Events, error)
	_SlotMappings  func() SlotMappings
}

func (h *testHandler) FormName() string {
	if h._FormName != nil {
		return h._FormName()
	}
	return "action_test_form"
}

func (h *testHandler) RequiredSlots(ctx *Context) []string {
	if h._RequiredSlots != nil {
		return h._RequiredSlots(ctx)
	}
	return []string{}
}

func (h *testHandler) Submit(ctx *Context, disp *action.CollectingDispatcher) (events rasa.Events, err error) {
	if h._Submit != nil {
		events, err = h._Submit(ctx, disp)
	}
	return
}

func (h *testHandler) SlotMappings() SlotMappings {
	if h._SlotMappings != nil {
		return h._SlotMappings()
	}
	return h.DefaultEmbed.SlotMappings()
}

// TestAction provides tests for methods on the base form.Action type.
func TestAction(t *testing.T) {

	t.Run("extract requested slot default", func(t *testing.T) {
		tracker := &rasa.Tracker{
			SenderID: "default",
			Slots: rasa.Slots{
				RequestedSlot: "some_slot",
			},
			LatestMessage: &rasa.ParseResult{
				Entities: []rasa.Entity{{
					Entity: "some_slot",
					Value:  "some_value",
				}},
			},
			Events:           rasa.Events{},
			Paused:           false,
			FollowupAction:   "",
			LatestActionName: "action_listen",
		}

		actx := Context{
			&action.Context{Tracker: tracker},
			&testHandler{},
		}
		slotValues := actx.extractRequestedSlot("some_slot")

		expect := rasa.Slots{"some_slot": rasa.StringValue("some_value")}
		require.Equal(t, expect, slotValues)
	})

	t.Run("extract requested slot from entity no intent", func(t *testing.T) {
		handler := &testHandler{
			_FormName: func() string { return "some_form" },
			_SlotMappings: func() SlotMappings {
				return SlotMappings{
					"some_slot": {FromEntity{
						Entity: "some_entity",
					}},
				}
			},
		}

		tracker := &rasa.Tracker{
			SenderID: "default",
			Slots: rasa.Slots{
				RequestedSlot: "some_slot",
			},
			LatestMessage: &rasa.ParseResult{
				Entities: []rasa.Entity{{
					Entity: "some_entity",
					Value:  "some_value",
				}},
			},
			LatestActionName: "action_listen",
		}

		actx := Context{
			&action.Context{Tracker: tracker},
			handler,
		}
		slotValues := actx.extractRequestedSlot("some_slot")

		expect := rasa.Slots{"some_slot": rasa.StringValue("some_value")}
		require.Equal(t, expect, slotValues)
	})
}
