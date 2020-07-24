// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import "go.scarlet.dev/rasa/action"

// QueryAction implements the action.Handler interface for KnowledgeBaseAction.
type QueryAction struct {
}

// Context TODO
//
// Context provides default implementations for methods of the KnowledgeBase interface, if
type Context struct {
	*action.Context
}
