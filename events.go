// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"go.scarlet.dev/rasa/internal/handle"
)

// EventType TODO
type EventType string

// TODO
const (
	EventTypeActionExecuted          = EventType("action")
	EventTypeActionExecutionRejected = EventType("action_execution_rejected")
	EventTypeActionReverted          = EventType("undo")
	EventTypeAgentUttered            = EventType("agent")
	EventTypeAllSlotsReset           = EventType("reset_slots")
	EventTypeBotUttered              = EventType("bot")
	EventTypeConversationPaused      = EventType("pause")
	EventTypeConversationResumed     = EventType("resume")
	EventTypeFollowupAction          = EventType("followup")
	EventTypeActiveLoop              = EventType("form")
	EventTypeLoopInterrupted         = EventType("form_validation")
	EventTypeReminderCancelled       = EventType("cancel_reminder")
	EventTypeReminderScheduled       = EventType("reminder")
	EventTypeRestarted               = EventType("restart")
	EventTypeSessionStarted          = EventType("session_started")
	EventTypeSlotSet                 = EventType("slot")
	EventTypeStoryExported           = EventType("export")
	EventTypeUserUtteranceReverted   = EventType("rewind")
	EventTypeUserUttered             = EventType("user")
)

// Event represents a serializable event object.
type Event interface {
	// Type returns the constant representing the event's type in a marshalled
	// JSON object.
	Type() EventType
}

// Events implements JSON Unmarshaling of Rasa's Event types.
type Events []Event

// ensure interface
var _ json.Unmarshaler = (*Events)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (l *Events) UnmarshalJSON(data []byte) (err error) {
	defer handle.Error(&err, func(err error) error {
		return errors.WithMessage(err, "unable to unmarshal EventList")
	})

	// initialize as empty eventlist
	*l = Events{}

	// get the Raw event messages
	var events []json.RawMessage
	if err = json.Unmarshal(data, &events); err != nil {
		return
	}

	for i := range events {
		var evt Event
		if evt, err = unmarshalEvent(events[i]); err != nil {
			return
		}
		*l = append(*l, evt)
	}

	return
}

// unmarshalEvent will unmarshal the provided serialized JSON as an Event.
func unmarshalEvent(data []byte) (evt Event, err error) {
	// marker is used to determine the type of the serialized JSON object.
	//
	// TODO(ed): This method of extracting the event before the type as a single
	// field struct requires each event to be parsed twice. Depending on the
	// potential size of event lists and individual events, it may be worth
	// benchmarking different methods of performing this serialization.
	var marker struct {
		Event EventType `json:"event"`
	}
	if err = json.Unmarshal(data, &marker); err != nil {
		return
	}

	switch marker.Event {
	default:
		// error case - unknown of unsupported event type
		err = fmt.Errorf("invalid event type [%s]", marker.Event)
		return
	case EventTypeActionExecuted:
		evt = new(ActionExecuted)
	case EventTypeActionExecutionRejected:
		evt = new(ActionExecutionRejected)
	case EventTypeActionReverted:
		evt = new(ActionReverted)
	case EventTypeAgentUttered:
		evt = new(AgentUttered)
	case EventTypeAllSlotsReset:
		evt = new(AllSlotsReset)
	case EventTypeBotUttered:
		evt = new(BotUttered)
	case EventTypeConversationPaused:
		evt = new(ConversationPaused)
	case EventTypeConversationResumed:
		evt = new(ConversationResumed)
	case EventTypeFollowupAction:
		evt = new(FollowupAction)
	case EventTypeReminderCancelled:
		evt = new(ReminderCancelled)
	case EventTypeReminderScheduled:
		evt = new(ReminderScheduled)
	case EventTypeRestarted:
		evt = new(Restarted)
	case EventTypeSessionStarted:
		evt = new(SessionStarted)
	case EventTypeSlotSet:
		evt = new(SlotSet)
	case EventTypeStoryExported:
		evt = new(StoryExported)
	case EventTypeUserUtteranceReverted:
		evt = new(UserUtteranceReverted)
	case EventTypeUserUttered:
		evt = new(UserUttered)
	}

	// unmarshal the JSON event
	err = json.Unmarshal(data, evt)
	return
}

// structToMap implements a single-layer conversion from a struct to a
// string-indexed map.
func structToMap(s interface{}) map[string]interface{} {
	// TODO(ed): replace this with a self-maintained implementation?
	st := structs.New(s)
	st.TagName = "json"
	return st.Map()
}

// marshalEvent will marshal the provided Event value into its JSON
// representation.
func marshalEvent(e Event) ([]byte, error) {
	// turn the event into a map
	result := structToMap(e)

	// add the event type property
	result["event"] = e.Type()
	return json.Marshal(result)
}

// ActionExecuted TODO
type ActionExecuted struct {
	Timestamp  Time    `json:"timestamp,omitempty"`
	ActionName string  `json:"name"`
	Policy     string  `json:"policy,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

// ActionExecutionRejected TODO
type ActionExecutionRejected struct {
	Timestamp  Time    `json:"timestamp,omitempty"`
	ActionName string  `json:"name"`
	Policy     string  `json:"policy,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

// ActionReverted TODO
type ActionReverted struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// AgentUttered TODO
type AgentUttered struct {
	Timestamp Time    `json:"timestamp,omitempty"`
	Text      string  `json:"test,omitempty"`
	Data      JSONMap `json:"data,omitempty"`
}

// AllSlotsReset TODO
type AllSlotsReset struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// BotUttered TODO
type BotUttered struct {
	Timestamp Time    `json:"timestamp,omitempty"`
	Text      string  `json:"text,omitempty"`
	Data      JSONMap `json:"data,omitempty"`
	Metadata  JSONMap `json:"metadata,omitempty"`
}

// ConversationPaused TODO
type ConversationPaused struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// ConversationResumed TODO
type ConversationResumed struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// FollowupAction TODO
type FollowupAction struct {
	Timestamp  Time   `json:"timestamp,omitempty"`
	ActionName string `json:"name"`
}

// ActiveLoop TODO
type ActiveLoop struct {
	Timestamp Time   `json:"timestamp,omitempty"`
	Name      string `json:"name,omitempty"`
}

// LoopInterrupted TODO
type LoopInterrupted struct {
	Timestamp     Time `json:"timestamp,omitempty"`
	IsInterrupted bool `json:"is_interrupted,omitempty"`
}

// ReminderCancelled TODO
type ReminderCancelled struct {
	Timestamp  Time      `json:"timestamp,omitempty"`
	Name       string    `json:"name,omitempty"`
	IntentName string    `json:"intent,omitempty"`
	Entities   []JSONMap `json:"entities,omitempty"`
}

// ReminderScheduled TODO
type ReminderScheduled struct {
	Timestamp         Time      `json:"timestamp,omitempty"`
	Name              string    `json:"name,omitempty"`
	IntentName        string    `json:"intent,omitempty"`
	Entities          []JSONMap `json:"entities,omitempty"`
	DateTime          time.Time `json:"date_time"`
	KillOnUserMessage bool      `json:"kill_on_user_msg"`
}

// Restarted TODO
type Restarted struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// SessionStarted TODO
type SessionStarted struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// SlotSet TODO
type SlotSet struct {
	Timestamp Time        `json:"timestamp,omitempty"`
	Key       string      `json:"name"`
	Value     interface{} `json:"value,omitempty"`
}

// StoryExported TODO
type StoryExported struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// UserUtteranceReverted TODO
type UserUtteranceReverted struct {
	Timestamp Time `json:"timestamp,omitempty"`
}

// UserUttered TODO
type UserUttered struct {
	Timestamp    Time    `json:"timestamp,omitempty"`
	Text         string  `json:"text,omitempty"`
	ParseData    JSONMap `json:"parse_data,omitempty"`
	InputChannel string  `json:"input_channel,omitempty"`
}

// Type implements Event.
func (ActionExecuted) Type() EventType { return EventTypeActionExecuted }

// Type implements Event.
func (ActionExecutionRejected) Type() EventType { return EventTypeActionExecutionRejected }

// Type implements Event.
func (ActionReverted) Type() EventType { return EventTypeActionReverted }

// Type implements Event.
func (AgentUttered) Type() EventType { return EventTypeAgentUttered }

// Type implements Event.
func (AllSlotsReset) Type() EventType { return EventTypeAllSlotsReset }

// Type implements Event.
func (BotUttered) Type() EventType { return EventTypeBotUttered }

// Type implements Event.
func (ConversationPaused) Type() EventType { return EventTypeConversationPaused }

// Type implements Event.
func (ConversationResumed) Type() EventType { return EventTypeConversationResumed }

// Type implements Event.
func (FollowupAction) Type() EventType { return EventTypeFollowupAction }

// Type implements Event.
func (ActiveLoop) Type() EventType { return EventTypeActiveLoop }

// Type implements Event.
func (LoopInterrupted) Type() EventType { return EventTypeLoopInterrupted }

// Type implements Event.
func (ReminderCancelled) Type() EventType { return EventTypeReminderCancelled }

// Type implements Event.
func (ReminderScheduled) Type() EventType { return EventTypeReminderScheduled }

// Type implements Event.
func (Restarted) Type() EventType { return EventTypeRestarted }

// Type implements Event.
func (SessionStarted) Type() EventType { return EventTypeSessionStarted }

// Type implements Event.
func (SlotSet) Type() EventType { return EventTypeSlotSet }

// Type implements Event.
func (StoryExported) Type() EventType { return EventTypeStoryExported }

// Type implements Event.
func (UserUtteranceReverted) Type() EventType { return EventTypeUserUtteranceReverted }

// Type implements Event.
func (UserUttered) Type() EventType { return EventTypeUserUttered }

// MarshalJSON implements json.Marshaler.
func (e *ActionExecuted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ActionExecutionRejected) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ActionReverted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *AgentUttered) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *AllSlotsReset) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *BotUttered) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ConversationPaused) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ConversationResumed) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *FollowupAction) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ActiveLoop) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *LoopInterrupted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ReminderCancelled) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *ReminderScheduled) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *Restarted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *SessionStarted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *SlotSet) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *StoryExported) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *UserUtteranceReverted) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// MarshalJSON implements json.Marshaler.
func (e *UserUttered) MarshalJSON() ([]byte, error) { return marshalEvent(e) }

// ensure interfaces
var _ Event = (*ActionExecuted)(nil)
var _ Event = (*ActionExecutionRejected)(nil)
var _ Event = (*ActionReverted)(nil)
var _ Event = (*AgentUttered)(nil)
var _ Event = (*AllSlotsReset)(nil)
var _ Event = (*BotUttered)(nil)
var _ Event = (*ConversationPaused)(nil)
var _ Event = (*ConversationResumed)(nil)
var _ Event = (*FollowupAction)(nil)
var _ Event = (*ActiveLoop)(nil)
var _ Event = (*LoopInterrupted)(nil)
var _ Event = (*ReminderCancelled)(nil)
var _ Event = (*ReminderScheduled)(nil)
var _ Event = (*Restarted)(nil)
var _ Event = (*SessionStarted)(nil)
var _ Event = (*SlotSet)(nil)
var _ Event = (*StoryExported)(nil)
var _ Event = (*UserUtteranceReverted)(nil)
var _ Event = (*UserUttered)(nil)

// ensure interfaces
var _ json.Marshaler = (*ActionExecuted)(nil)
var _ json.Marshaler = (*ActionExecutionRejected)(nil)
var _ json.Marshaler = (*ActionReverted)(nil)
var _ json.Marshaler = (*AgentUttered)(nil)
var _ json.Marshaler = (*AllSlotsReset)(nil)
var _ json.Marshaler = (*BotUttered)(nil)
var _ json.Marshaler = (*ConversationPaused)(nil)
var _ json.Marshaler = (*ConversationResumed)(nil)
var _ json.Marshaler = (*FollowupAction)(nil)
var _ json.Marshaler = (*ActiveLoop)(nil)
var _ json.Marshaler = (*LoopInterrupted)(nil)
var _ json.Marshaler = (*ReminderCancelled)(nil)
var _ json.Marshaler = (*ReminderScheduled)(nil)
var _ json.Marshaler = (*Restarted)(nil)
var _ json.Marshaler = (*SessionStarted)(nil)
var _ json.Marshaler = (*SlotSet)(nil)
var _ json.Marshaler = (*StoryExported)(nil)
var _ json.Marshaler = (*UserUtteranceReverted)(nil)
var _ json.Marshaler = (*UserUttered)(nil)
