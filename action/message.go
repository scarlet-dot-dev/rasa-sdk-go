// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import "encoding/json"

// Message TODO
type Message struct {
	Text        string                 `json:"text,omitempty"`
	Image       string                 `json:"image,omitempty"`
	JSONMessage JSONMap                `json:"json_message,omitempty"`
	Template    string                 `json:"template,omitempty"`
	Attachment  string                 `json:"attachment,omitempty"`
	Buttons     []JSONMap              `json:"buttons,omitempty"`
	Elements    []JSONMap              `json:"elements,omitempty"`
	Kwargs      map[string]interface{} `json:"-"` // FIXME implement custom ser/de
}

// ensure interfaces.
var _ json.Marshaler = (*Message)(nil)

// TODO(ed): is Unmarshaler needed? it seems Message is a one-way dto.
// var _ json.Unmarshaler = (*Message)(nil)

// WithKwargs adds the free-form kwargs to the m.Kwargs.
func (m *Message) WithKwargs(kwargs map[string]interface{}) *Message {
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
