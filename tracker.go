// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import "fmt"

// Tracker contains the state of the Tracker sent to the action server by the
// Rasa engine.
type Tracker struct {
	ConversationID     string       `json:"conversation_id"`
	SenderID           string       `json:"sender_id"` // TODO(ed): verify if this field is ever set
	Slots              Slots        `json:"slots,omitempty"`
	LatestMessage      *ParseResult `json:"latest_message,omitempty"`
	LatestActionName   string       `json:"latest_action_name,omitempty"`
	LatestEventTime    Time         `json:"latest_event_time,omitempty"`
	LatestInputChannel string       `json:"latest_input_channel,omitempty"`
	Events             Events       `json:"events"`
	Paused             bool         `json:"paused"`
	FollowupAction     string       `json:"followup_action,omitempty"`
	ActiveForm         *ActiveForm  `json:"active_form,omitempty"`
}

// HasSlots returns whether there are any Slots present in the Tracker.
func (t *Tracker) HasSlots() bool {
	return len(t.Slots) > 0
}

// LatestEntityValues returns the entity values found for the passed entity name
// in the latest message.
func (t *Tracker) LatestEntityValues(entity, role, group string) (values []string) {
	if len(t.LatestMessage.Entities) == 0 {
		return nil
	}

	entities := t.LatestMessage.Entities
	for i := range entities {
		val := entities[i]
		if val.Entity != entity || val.Value == "" {
			continue
		}
		if role != "" && role != val.Role {
			continue
		}
		if group != "" && group != val.Group {
			continue
		}

		values = append(values, val.Value)
	}
	return
}

// HasActiveForm returns whether the Tracker state represents an active Form.
func (t *Tracker) HasActiveForm() bool {
	return t.ActiveForm.IsActive()
}

// ActiveForm holds a ActiveForm description in the Tracker.
type ActiveForm struct {
	Name           string       `json:"name"`
	Validate       *bool        `json:"validate,omitempty"`
	Rejected       bool         `json:"rejected,omitempty"`
	TriggerMessage *ParseResult `json:"trigger_message"`
}

//
func (f *ActiveForm) ShouldValidate() bool {
	if f.Validate == nil {
		return true
	}
	return *f.Validate
}

// IsActive returns whether f represents an active Form.
func (f *ActiveForm) IsActive() bool {
	return f != nil && f.Name != ""
}

// Is returns whether f represents a Form with the provided name.
func (f *ActiveForm) Is(name string) bool {
	return f != nil && f.Name == name
}

// ParseResult holds a processed (parsed) message description.
type ParseResult struct {
	Intent        Intent   `json:"intent"`
	IntentRanking []Intent `json:"intent_ranking,omitempty"`
	Entities      []Entity `json:"entities,omitempty"`
	Text          string   `json:"text,omitempty"`
}

// Intent describes an intent and its detected confidence.
type Intent struct {
	Confidence float64 `json:"confidence,omitempty"`
	Name       string  `json:"name,omitempty"`
}

// Entity describes an entity and its detected location, value, and confidence.
type Entity struct {
	Start      int     `json:"start"`
	End        int     `json:"end"`
	Value      string  `json:"value"`
	Entity     string  `json:"entity"`
	Confidence float64 `json:"confidence"`
	Group      string  `json:"group,omitempty"`
	Role       string  `json:"role,omitempty"`
}

// Slots is a wrapper type around slots.
type Slots map[string]interface{}

// Has returns true if the Slots contains the requested Slot.
func (m Slots) Has(slot string) bool {
	val, exists := m[slot]
	return exists && val != nil
}

// Update will copy the values present in s into m.
func (m Slots) Update(s Slots) {
	for key := range s {
		m[key] = s[key]
	}
}

// String implements fmt.Stringer.
//
// String returns a simple json-like representation of the map's values.
func (m Slots) String() string {
	return fmt.Sprintf("%#v", m)
}

// EntityValues TODO
type EntityValue interface {
	AsString() string
	AsSlice() []string
	Count() int
}

// ensure interfaces
var _ EntityValue = (StringValue)("")
var _ EntityValue = (SliceValue)(nil)

// EntityValue TODO
type StringValue string

// AsString implements EntityValue
func (v StringValue) AsString() string {
	return string(v)
}

// AsSlice implements EntityValue.
func (v StringValue) AsSlice() []string {
	return []string{string(v)}
}

// Count implements EntityValue.
func (v StringValue) Count() int {
	return 1
}

// SliceValue TODO
type SliceValue []string

// AsString implements EntityValue
func (v SliceValue) AsString() string {
	if len(v) > 0 {
		return v[0]
	}
	return ""
}

// AsSlice implements EntityValue.
func (v SliceValue) AsSlice() []string {
	return v
}

// Count implements EntityValue.
func (v SliceValue) Count() int {
	return len(v)
}
