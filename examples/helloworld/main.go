// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

// Package main contains an implementation of a "hello world" example to demo
// the usage of the Rasa SDK.
package main

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// HelloWorld is an action which utters a "hello world" message to the user when
// it's triggered during a conversation.
type HelloWorld struct{}

// ensure interface
var _ action.Handler = (*HelloWorld)(nil)

// ActionName implements action.Handler.
func (HelloWorld) ActionName() string {
	return "action_hello_world"
}

// Run implements action.Handler.
func (h *HelloWorld) Run(
	// actx contains the "request context" for this handler. It can be used to
	// access information about the user, the request lifetime, the conversation
	// state, and the logger.
	actx action.Context,

	// dispatcher can be used to send messages to the user.
	dispatcher *action.CollectingDispatcher,
) (
	// events will contain any events that should be sent back to the user. In
	// this example, this will be left empty - we will only use the dispatcher
	// to utter a message.
	events rasa.Events,

	// err should return a non-nil error in the event of a non-recoverable or
	// unexpected error. The error will be caught by the action handler, and
	// will be turned into an error response to be sent back to the Rasa Core
	// server.
	err error,
) {
	// Any action code goes here - go wild!

	// use the "context" package as you would in any other HTTP server.
	ctx, cancel := context.WithTimeout(actx.Context(), time.Minute)
	defer cancel()

	// The context.Context is linked to the HTTP request that triggered the action handler.
	_ = ctx

	// You can access the logger through the context
	actx.Logger().Infof("Hello world log!")

	// dispatch a text message to the user.
	dispatcher.Utter(&rasa.Message{
		Text: "Hello world from Rasa! Check out https://go.scarlet.dev/rasa for more documentation.",
		// NOTE: a rasa.Message can also contain buttons, attachments,
		// templates, and more!
	})

	// using an unamed return will (in this case) send back the default values
	// for `events` and `err`.
	return
}

// main starts the action server.
func main() {
	// create the action.Server.
	serv := action.NewServer(
		&HelloWorld{},
	)

	// setting PrettyJSON is useful for debugging as it will format the JSON
	// returned by the API with indentation.
	serv.PrettyJSON = true

	// Logger can hold a logger.
	serv.Logger = logrus.New()
	serv.Logger.Infof(
		"To see the Action Server in action, try visiting http://localhost:%s/actions",
		rasa.DefaultServerPort,
	)

	// create the HTTP server and use the action.Server as handler.
	http.ListenAndServe(
		// listen on the default port, at 0.0.0.0:5055
		":"+rasa.DefaultServerPort,

		// serv implements http.Handler, so it can be passed as the only
		// argument for quick setups.
		serv,
	)

	// And we're done!
	// The endpoint will now respond to the basic /health, /actions, and
	// /webhook urls. Try it: http://localhost:5055/actions.
}
