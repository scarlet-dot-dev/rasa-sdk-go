// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

// Message TODO
type Message struct {
	Text        string    `json:"text,omitempty"`
	Image       string    `json:"image,omitempty"`
	JSONMessage JSONMap   `json:"json_message,omitempty"`
	Template    string    `json:"template,omitempty"`
	Attachment  string    `json:"attachment,omitempty"`
	Buttons     []JSONMap `json:"buttons,omitempty"`
	Elements    []JSONMap `json:"elements,omitempty"`
}

// JSONMap TODO
type JSONMap map[string]interface{}

// CollectingDispatcher implements a response collector.
//
// The implementation differs from Rasa's Python SDK by not supporting `kwargs`.
type CollectingDispatcher []Message

// Utter will add the Message to the response list.
func (d *CollectingDispatcher) Utter(msg *Message) {
	*d = append(*d, *msg)
}

// Clear will empty the CollectingDispatcher.
func (d *CollectingDispatcher) Clear() {
	*d = nil
}
