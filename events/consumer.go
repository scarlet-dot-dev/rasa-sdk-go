// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package events provides types for consumers of events passed to event
// brokers.
//
// See https://rasa.com/docs/rasa/event-brokers for more information.
package events

import "go.scarlet.dev/rasa"

// Entry is a single entry posted to an event broker.
type Entry struct {
	SenderID  string         `json:"sender_id"`
	Timestamp rasa.Time      `json:"timestamp"`
	Event     rasa.EventType `json:"event"`
	Text      string         `json:"text"`
	Data      string         `json:"data,omitempty"`
	Metadata  rasa.JSONMap   `json:"metadata,omitempty"`
}
