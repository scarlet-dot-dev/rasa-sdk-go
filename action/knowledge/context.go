// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import (
	"fmt"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// Utterer TODO
type Utterer interface {
	UtterObjects(
		disp *action.CollectingDispatcher,
		otype ObjectType,
		objs []Object,
	)

	UtterAttribute(
		disp *action.CollectingDispatcher,
		obj Object,
		attr Attribute,
	)
}

// Context TODO
//
// Context provides default implementations for methods of the KnowledgeBase
// interface.
type Context struct {
	*action.Context
}

//
var _ Utterer = (*Context)(nil)

// UtterObjects implements a default for the Utterer interface.
func (c *Context) UtterObjects(
	disp *action.CollectingDispatcher,
	otype ObjectType,
	objs []Object,
) {
	if len(objs) == 0 {
		disp.Utter(&rasa.Message{
			Text: fmt.Sprintf("I could not find any objects of type '%s'", otype.TypeName()),
		})
		return
	}

	disp.Utter(&rasa.Message{
		Text: fmt.Sprintf("Found the following objects of type '%s'", otype.TypeName()),
	})

	for i := range objs {
		disp.Utter(&rasa.Message{
			Text: fmt.Sprintf("%d: %s", i+1, objs[i].Representation()),
		})
	}

	return
}

// UtterAttribute utters a response that informs the user about the attribute
// value of the attribute of interest.
func (c *Context) UtterAttribute(
	disp *action.CollectingDispatcher,
	obj Object,
	attr Attribute,
) {
	if len(attr.Value) > 0 {
		disp.Utter(&rasa.Message{
			Text: fmt.Sprintf(
				"'%s' has the value '%s' for attribute '%s'",
				obj.Key(),
				attr.Value,
				attr.Name,
			),
		})
		return
	}

	disp.Utter(&rasa.Message{
		Text: fmt.Sprintf(
			"Did not find a valid value for attribute '%s' for object '%s'.",
			attr.Name,
			obj.Key(),
		),
	})
	return
}
