// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import (
	"fmt"
	"reflect"
)

// Tracker contains the state of the Tracker sent to the action server by the
// Rasa engine.
type Tracker struct {
	SenderID         string       `json:"sender_id"`
	Slots            Slots        `json:"slots,omitempty"`
	LatestMessage    *ParseResult `json:"latest_message,omitempty"`
	LatestActionName string       `json:"latest_action_name,omitempty"`
	Events           Events       `json:"events"`
	Paused           bool         `json:"paused"`
	FollowupAction   string       `json:"followup_action,omitempty"`
	ActiveLoop       *TActiveLoop `json:"active_loop,omitempty"`
}

// HasSlots returns whether there are any Slots present in the Tracker.
func (t *Tracker) HasSlots() bool {
	return len(t.Slots) > 0
}

// HasActiveForm returns whether the Tracker state represents an active Form.
func (t *Tracker) HasActiveForm() bool {
	return t.HasActiveLoop()
}

// HasActiveLoop returns whether the Tracker state represents an active Loop.
func (t *Tracker) HasActiveLoop() bool {
	return t.ActiveLoop.IsActive()
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

// Entity returns the current value of the requested entity as an EntityValue.
func (t *Tracker) Entity(name, role, group string) (value EntityValue) {
	raw := t.EntityValues(name, role, group)
	if len(raw) == 1 {
		value = StringValue(raw[0])
		return
	}
	value = SliceValue(raw)
	return
}

// EntityValues returns the current value of the requested entity as a
// slice. The slice may be empty, and may contain 0 or more entries.
func (t *Tracker) EntityValues(
	entity, role, group string,
) (values []string) {
	values = t.LatestEntityValues(entity, role, group)
	return
}

// Slot returns the value of the slot as an interface. The `ok` flag
// indicates whether the slot was present.
func (t *Tracker) Slot(name string) (val interface{}, ok bool) {
	val, ok = t.Slots[name]
	return
}

// SlotAs attempts to assign the value of the slot to the `dst` pointer. The
// `ok` flag indicates whether the slot was present, and successfully
// assigned to `dst`.
//
// Passing a non-assignable type to `dst` leads to undefined behaviour.
func (t *Tracker) SlotAs(name string, dst interface{}) (ok bool) {
	val, exists := t.Slot(name)
	if !exists {
		return
	}

	rdst := reflect.ValueOf(dst)
	rval := reflect.ValueOf(val)
	if ok = rdst.CanSet() && rval.Type().AssignableTo(rdst.Type()); !ok {
		return
	}

	rdst.Set(reflect.ValueOf(val))
	return
}

// SlotsToValidate returns the slots which were recently set.
//
// This can be used to validate form slots after they were extracted.
func (t *Tracker) SlotsToValidate() (slots Slots) {
	slots = make(map[string]interface{})
	events := t.Events

	// look at the newest events in the tracker
	for i := len(events) - 1; i >= 0; i-- {
		// The `FormAction` in Rasa Open Source will append all slot candidates
		// at the end of the tracker events.
		if se, ok := events[i].(*SlotSet); ok {
			slots[se.Key] = se.Value
			continue
		}

		// found a different event type - stop the loop
		break
	}

	// return the found slots
	return
}

// TActiveLoop holds a ActiveLoop description in the Tracker.
type TActiveLoop struct {
	Name           string       `json:"name"`
	Validate       *bool        `json:"validate,omitempty"`
	Rejected       bool         `json:"rejected,omitempty"`
	TriggerMessage *ParseResult `json:"trigger_message"`
}

// ShouldValidate returns whether the form should validate itself. It will
// return true unless validation has been explicitely disabled.
func (f *TActiveLoop) ShouldValidate() bool {
	if f.Validate == nil {
		return true
	}
	return *f.Validate
}

// IsActive returns whether f represents an active Form.
func (f *TActiveLoop) IsActive() bool {
	return f != nil && f.Name != ""
}

// Is returns whether f represents a Form with the provided name.
func (f *TActiveLoop) Is(name string) bool {
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

// EntityValue specifies an interface over extracted entity values that can be
// either a single string or a list of strings.
type EntityValue interface {
	AsString() string
	AsSlice() []string
	Count() int
}

// ensure interfaces
var _ EntityValue = (StringValue)("")
var _ EntityValue = (SliceValue)(nil)

// StringValue implements the EntityValue interface for entities that only
// contain a single value.
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

// SliceValue implements the EntityValue interface for entities that contain
// more than one value.
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
