// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

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

// WithKwargs adds the free-form kwargs to the m.Kwargs.
func (m *Message) WithKwargs(kwargs map[string]interface{}) *Message {
	for key := range kwargs {
		m.Kwargs[key] = kwargs[key]
	}
	return m
}
