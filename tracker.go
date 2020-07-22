// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

// Tracker contains the state of the Tracker sent to the action server by the
// Rasa engine.
type Tracker struct {
	ConversationID     string       `json:"conversation_id"`
	SenderID           string       `json:"sender_id"` // TODO(ed): verify if this field is ever set
	Slots              SlotMap      `json:"slots,omitempty"`
	LatestMessage      *ParseResult `json:"latest_message,omitempty"`
	LatestActionName   string       `json:"latest_action_name,omitempty"`
	LatestEventTime    Time         `json:"latest_event_time,omitempty"`
	LatestInputChannel string       `json:"latest_input_channel,omitempty"`
	Events             EventList    `json:"events"`
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
func (t *Tracker) LatestEntityValues(entity string) (out []string) {
	if len(t.LatestMessage.Entities) == 0 {
		return
	}

	for _, entry := range t.LatestMessage.Entities {
		entry := entry
		if entry.Entity == entity && entry.Value != "" {
			out = append(out, entry.Value)
		}
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
	Validate       bool         `json:"validate,omitempty"`
	Rejected       bool         `json:"rejected,omitempty"`
	TriggerMessage *ParseResult `json:"trigger_message"`
}

// IsActive returns whether f represents an active Form.
func (f *ActiveForm) IsActive() bool {
	return f != nil && f.Name != ""
}

// Is returns whether f represents a Form with the provided name.
func (f *ActiveForm) Is(name string) bool {
	return f != nil && f.Name == name
}

// ShouldValidate TODO
func (f *ActiveForm) ShouldValidate(defaults bool) bool {
	return f.Validate || defaults
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
}

// SlotMap is a wrapper type around slots.
type SlotMap map[string]interface{}

// Has returns true if the SlotMap contains the requested Slot.
func (m SlotMap) Has(slot string) bool {
	_, exists := m[slot]
	return exists
}
