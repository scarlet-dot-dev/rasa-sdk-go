// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import (
	"time"

	"go.scarlet.dev/rasa"
)

// Get the name of the object the user referred to. Either the NER detected the
// object and stored its name in the corresponding slot (e.g. "PastaBar"
// is detected as "restaurant") or the user referred to the object by any kind of
// mention, such as "first one" or "it".
func getObjectName(ctx *Context, omapper OrdinalMapper, useLastObjectMention bool) (val string, resolved bool) {
	var mention string
	mentionOK := ctx.SlotAs(SlotMention, &mention)
	var objType string
	_ = ctx.SlotAs(SlotObjectType, &objType)

	if mentionOK && mention != "" {
		val, resolved = resolveMention(ctx, omapper)
		return
	}

	var objName string
	objNameOK := ctx.SlotAs(objType, &objName)
	if objNameOK && objName != "" {
		val, resolved = objName, true
		return
	}

	if useLastObjectMention {
		resolved = ctx.SlotAs(SlotLastObject, &val)
		return
	}

	return
}

// Resolve the given mention to the name of the actual object.
//
// Different kind of mentions exist. We distinguish between ordinal mentions and
// all others for now.
// For ordinal mentions we resolve the mention of an object, such as 'the first
// one', to the actual object name. If multiple objects are listed during the
// conversation, the objects are stored in the slot 'knowledge_base_listed_objects'
// as a list. We resolve the mention, such as 'the first one', to the list index
// and retrieve the actual object (using the 'ordinal_mention_mapping').
// For any other mention, such as 'it' or 'that restaurant', we just assume the
// user is referring to the last mentioned object in the conversation.
func resolveMention(ctx *Context, omapper OrdinalMapper) (value string, resolved bool) {
	var mention string
	mentionOK := ctx.SlotAs(SlotMention, &mention)
	var listedItems []string
	listedItemsOK := ctx.SlotAs(SlotListedObjects, listedItems)
	var lastObj string
	lastObjOK := ctx.SlotAs(SlotLastObject, &lastObj)
	lastObjType, _ := ctx.Slot(SlotLastObjectType)
	currObjType, _ := ctx.Slot(SlotObjectType)

	if !mentionOK || mention == "" {
		return
	}

	if listedItemsOK && len(listedItems) > 0 && omapper.Has(mention) {
		resolved = true
		value = omapper.Map(mention)(listedItems)
		return
	}

	if currObjType == lastObjType {
		value, resolved = lastObj, lastObjOK
		return
	}

	return
}

// If the user mentioned one or multiple attributes of the provided object_type in
// an utterance, we extract all attribute values from the tracker and put them
// in a list. The list is used later on to filter a list of objects.
//
// For example: The user says 'What Italian restaurants do you know?'.
// The NER should detect 'Italian' as 'cuisine'.
// We know that 'cuisine' is an attribute of the object type 'restaurant'.
// Thus, this method returns [{'name': 'cuisine', 'value': 'Italian'}] as
// list of attributes for the object type 'restaurant'.
func getAttributeSlots(ctx *Context, attrs []string) []Attribute {
	vals := []Attribute{}

	for i := range attrs {
		attr := attrs[i]
		attrVal, ok := ctx.Tracker.Slots[attr]
		if sval, sok := attrVal.(string); ok && sok && sval != "" {
			vals = append(vals, Attribute{
				Name:  attr,
				Value: sval,
			})
		}
	}

	return vals
}

// Reset all attribute slots of the current object type.
//
// If the user is saying something like "Show me all restaurants with Italian
// cuisine.", the NER should detect "restaurant" as "object_type" and "Italian" as
// "cuisine" object. So, we should filter the restaurant objects in the
// knowledge base by their cuisine (= Italian). When listing objects, we check
// what attributes are detected by the NER. We take all attributes that are set,
// e.g. cuisine = Italian. If we don't reset the attribute slots after the request
// is done and the next utterance of the user would be, for example, "List all
// restaurants that have wifi.", we would have two attribute slots set: "wifi" and
// "cuisine". Thus, we would filter all restaurants for two attributes now:
// wifi = True and cuisine = Italian. However, the user did not specify any
// cuisine in the request. To avoid that we reset the attribute slots once the
// request is done.
func resetAttributeSlots(ctx *Context, attrs []string) (events rasa.Events) {
	for i := range attrs {
		attr := attrs[i]
		attrVal, ok := ctx.Tracker.Slots[attr]
		if ok && attrVal != nil {
			events = append(events, rasa.SlotSet{
				Key:       attr,
				Value:     nil,
				Timestamp: rasa.Time(time.Now()),
			})
		}
	}

	return
}

//
func sliceContains(slice []string, value string) bool {
	for i := range slice {
		if slice[i] == value {
			return true
		}
	}
	return false
}
