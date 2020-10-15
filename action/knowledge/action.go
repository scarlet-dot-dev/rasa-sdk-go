// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import (
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

//
const (
	TmplAskRephrase = "utter_ask_rephrase"
)

// QueryAction implements the action.Handler interface for
// ActionQueryKnowledgebase.
type QueryAction struct {
	KnowledgeBase           Storage
	Utterer                 Utterer
	IgnoreLastObjectMention bool
}

//
var _ action.Handler = (*QueryAction)(nil)

// ActionName implements action.Handler.
func (a *QueryAction) ActionName() string {
	return "action_query_knowledge_base"
}

// Run implements action.Handler.
//
// Run executes this action. If the user asks a question about an attribute,
// the knowledge base is queried for that attribute. Otherwise, if no
// attribute was detected in the request or the user is talking about a new
// object type, multiple objects of the requested type are returned from the
// knowledge base.
func (a *QueryAction) Run(
	ctx *action.Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	kctx := &Context{ctx}

	// get slots
	objType, hasObjType := ctx.Tracker.Slots[SlotObjectType].(string)
	lastObjType, hasLastObjType := ctx.Tracker.Slots[SlotLastObjectType].(string)
	attr, hasAttr := ctx.Tracker.Slots[SlotAttribute].(string)

	// be sure to reject empty as non-existent
	hasObjType = hasObjType && len(objType) > 0
	hasLastObjType = hasLastObjType && len(lastObjType) > 0
	hasAttr = hasAttr && len(attr) > 0

	isNewReq := objType != lastObjType

	if !hasObjType {
		disp.Utter(&rasa.Message{Template: TmplAskRephrase})
		return
	}

	if !hasAttr || isNewReq {
		return a.queryObjects(kctx, disp)
	}
	if hasAttr {
		return a.queryAttribute(kctx, disp)
	}

	disp.Utter(&rasa.Message{Template: TmplAskRephrase})
	return
}

// queryObjects queries the knowledge base for objects of the requested object
// type and outputs those to the user. The objects are filtered by any attribute
// the user mentioned in the request.
func (a *QueryAction) queryObjects(
	ctx *Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	//
	var objType string
	_ = ctx.SlotAs(SlotObjectType, &objType)

	//
	ts, err := a.KnowledgeBase.ForType(objType)
	if err != nil {
		return
	}

	//
	objAttrs := ts.Attributes()
	attrs := getAttributeSlots(ctx, objAttrs)
	objects := ts.GetObjects(attrs, 5)
	if len(objects) == 0 {
		events = resetAttributeSlots(ctx, objAttrs)
		return
	}

	//
	var lastObjID string
	if len(objects) > 0 {
		objects[0].Key()
	}

	//
	var listedObjectIDs []string
	for i := range objects {
		listedObjectIDs = append(listedObjectIDs, objects[i].Key())
	}

	// build the return events
	events = append(
		events,
		ctx.SetSlot(SlotObjectType, objType),
		ctx.ResetSlot(SlotMention),
		ctx.ResetSlot(SlotAttribute),
		ctx.SetSlot(SlotLastObject, lastObjID),
		ctx.SetSlot(SlotLastObjectType, objType),
		ctx.SetSlot(SlotListedObjects, listedObjectIDs),
	)
	events = append(events, resetAttributeSlots(ctx, objAttrs)...)
	return
}

// Queries the knowledge base for the value of the requested attribute of the
// mentioned object and outputs it to the user.
func (a *QueryAction) queryAttribute(
	ctx *Context,
	disp *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	//
	var objType string
	_ = ctx.SlotAs(SlotObjectType, &objType)
	var attr string
	attrOK := ctx.SlotAs(SlotAttribute, &attr)

	otype, err := a.KnowledgeBase.ForType(objType)
	if err != nil {
		return
	}

	oname, resolved := getObjectName(ctx, a.KnowledgeBase.OrdinalMapping(), !a.IgnoreLastObjectMention)
	if !resolved || !attrOK {
		a.utterAskRephrase(disp)
		events = append(events, ctx.ResetSlot(SlotMention))
		return
	}

	obj := otype.GetObject(oname)
	if obj == nil || !sliceContains(otype.Attributes(), attr) {
		a.utterAskRephrase(disp)
		events = append(events, ctx.ResetSlot(SlotMention))
		return
	}

	panic("TODO: finish implementation")
}

//
func (a *QueryAction) utterObjects(
	ctx *Context,
	disp *action.CollectingDispatcher,
	otype ObjectType,
	objs []Object,
) {
	if a.Utterer == nil {
		ctx.UtterObjects(disp, otype, objs)
		return
	}

	a.Utterer.UtterObjects(disp, otype, objs)
}

//
func (a *QueryAction) utterAttribute(
	ctx *Context,
	disp *action.CollectingDispatcher,
	otype Object,
	attr Attribute,
) {
	if a.Utterer == nil {
		ctx.UtterAttribute(disp, otype, attr)
		return
	}

	a.Utterer.UtterAttribute(disp, otype, attr)
}

//
func (a *QueryAction) utterAskRephrase(disp *action.CollectingDispatcher) {
	disp.Utter(&rasa.Message{
		Template: TmplAskRephrase,
	})
}
