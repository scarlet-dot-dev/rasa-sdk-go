// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nlg

import (
	"context"

	"go.scarlet.dev/rasa"
)

// Selector TODO
type Selector interface {
	// Select should return the proper Template for the provided context.
	Select(ctx *Context) (Template, error)
}

// Template TODO
type Template interface {
	// Generate
	Generate(ctx *Context) (rasa.Message, error)
}

// Context TODO
type Context struct {
	context.Context
	Name    string
	Args    rasa.JSONMap
	Tracker *rasa.Tracker
}

// SimpleSelector TODO
type SimpleSelector map[string]Template

// SimpleTemplate TODO
type SimpleTemplate string
