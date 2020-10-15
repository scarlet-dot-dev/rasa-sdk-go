// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

// Domain contains the configuration of the AI's domain as sent to the action
// server by Rasa's engine.
type Domain struct {
	Config    DomainConfig                     `json:"config"`
	Intents   []IntentDescription              `json:"intents"`
	Entities  []string                         `json:"entities"`
	Slots     map[string]SlotDescription       `json:"slots"`
	Responses map[string][]TemplateDescription `json:"responses"`
	Actions   []string                         `json:"actions"`
}

// DomainConfig contains domain settings.
type DomainConfig struct {
	StoreEntitiesAsSlots bool `json:"store_entities_as_slots"`
}

// IntentDescription contains a domain intent description.
type IntentDescription struct {
	UseEntities bool `json:"use_entities"`
}

// SlotDescription contains a domain slot description.
type SlotDescription struct {
	AutoFill     bool     `json:"auto_fill"`
	Type         string   `json:"type"`
	InitialValue string   `json:"initial_value,omitempty"`
	Values       []string `json:"values,omitempty"`
}

// TemplateDescription contains a domain template description.
type TemplateDescription struct {
	// the template text
	Text string
}
