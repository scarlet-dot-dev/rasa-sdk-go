// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import "encoding/json"

// JSONMap is a descriptive type alias for a free-form JSON Object.
type JSONMap = map[string]interface{}

// Message TODO
type Message struct {
	Text        string    `json:"text,omitempty"`
	Image       string    `json:"image,omitempty"`
	JSONMessage JSONMap   `json:"json_message,omitempty"`
	Template    string    `json:"template,omitempty"`
	Attachment  string    `json:"attachment,omitempty"`
	Buttons     []Button  `json:"buttons,omitempty"`
	Elements    []JSONMap `json:"elements,omitempty"`

	// Kwargs holds additional fields at the root level of the object that are
	// not otherwise provided in the default Message struct.
	//
	// Kwargs are used internally by the action server. Prefer using the
	// JSONMessage field if custom payloads are required to avoid conflicts with
	// the library.
	Kwargs JSONMap `json:"-"` // FIXME implement custom ser/de
}

// ensure interfaces.
var _ json.Marshaler = (*Message)(nil)

// TODO(ed): is Unmarshaler needed? it seems Message is a one-way dto.
// var _ json.Unmarshaler = (*Message)(nil)

// WithKwargs adds the free-form kwargs to the m.Kwargs.
func (m *Message) WithKwargs(kwargs JSONMap) *Message {
	for key := range kwargs {
		m.Kwargs[key] = kwargs[key]
	}
	return m
}

// MarshalJSON implements json.Marshaler.
func (m *Message) MarshalJSON() (data []byte, err error) {
	if m == nil {
		return
	}

	// standard fields
	raw := make(JSONMap)
	if m.Text != "" {
		raw["text"] = m.Text
	}
	if m.Image != "" {
		raw["image"] = m.Image
	}
	if m.JSONMessage != nil {
		raw["json_message"] = m.JSONMessage
	}
	if m.Template != "" {
		raw["template"] = m.Template
	}
	if m.Attachment != "" {
		raw["attachment"] = m.Attachment
	}
	if len(m.Buttons) > 0 {
		raw["buttons"] = m.Buttons
	}
	if len(m.Elements) > 0 {
		raw["elements"] = m.Elements
	}

	// copy KWargs
	for key := range m.Kwargs {
		raw[key] = m.Kwargs[key]
	}

	return json.Marshal(raw)
}

// Button defines the structure of a Button response.
//
// A button can be clicked by the user in a conversation.
type Button struct {
	// Title holds teh text on the button.
	Title string `json:"title"`
	// Payload holds the payload being sent if the button is pressed.
	Payload string `json:"payload"`
}
