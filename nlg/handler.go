// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package nlg

import (
	"net/http"

	"go.scarlet.dev/rasa"
)

// Request defines the request input received by the /nlg endpoint.
type Request struct {
	Template  string        `json:"template"`
	Argumants rasa.JSONMap  `json:"arguments"`
	Tracker   *rasa.Tracker `json:"tracker"`
	Channel   Channel       `json:"channel"`
}

// Channel contains a channel description.
type Channel struct {
	Name string `json:"name"`
}

// Response defines the response input sent back by the /nlg endpoint.
type Response rasa.Message

// Handler implements a request handler capable of responding to HTTP requests
// to the nlg endpoint.
type Handler struct {
	// TODO
}

// ensure interface
var _ http.Handler = (*Handler)(nil)

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO(ed): implement
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("TODO"))
}
