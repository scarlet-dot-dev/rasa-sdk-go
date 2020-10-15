// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

// Slots TODO
type Slots struct {
	KBAttribute      string   `json:"attribute,omitempty"`
	KBLastObject     string   `json:"knowledge_base_last_object,omitempty"`
	KBLastObjectType string   `json:"knowledge_base_last_object_type,omitempty"`
	KBListedObjects  []string `json:"knowledge_base_listed_objects,omitempty"`
	KBSlotMention    string   `json:"mention,omitempty"`
	KBObjectType     string   `json:"object_type,omitempty"`
}

//
const (
	SlotAttribute      = "attribute"
	SlotLastObject     = "knowledge_base_last_object"
	SlotLastObjectType = "knowledge_base_last_object_type"
	SlotListedObjects  = "knowledge_base_listed_objects"
	SlotMention        = "mention"
	SlotObjectType     = "object_type"
)
