// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callback

import (
	"context"
	"encoding/json"
	"net/http"

	"go.scarlet.dev/rasa"
)

// Receiver is the interface to be implemented by the application to handle
// received outputs.
type Receiver interface {
	// Receive will be called with the unmarshalled payload from the POST call.
	Receive(ctx context.Context, messages []rasa.Message) (err error)
}

// ensure interface
var _ Receiver = (ReceiverFn)(nil)
var _ Receiver = (*NoopReceiver)(nil)

// ReceiverFn implements the Receiver interface for functions.
type ReceiverFn func(context.Context, []rasa.Message) error

// Receive implements Receiver.
func (fn ReceiverFn) Receive(ctx context.Context, messages []rasa.Message) error {
	return fn(ctx, messages)
}

// WithLogging TODO
func WithLogging(r Receiver, logger func(m rasa.Message)) Receiver {
	return ReceiverFn(func(c context.Context, m []rasa.Message) error {
		for i := range m {
			logger(m[i])
		}
		return r.Receive(c, m)
	})
}

// NoopReceiver implements the Receiver interface with a no-op implementation or
// Receive.
//
// NoopReceiver can be used during debugging or development, such as combining
// it with `WithLogging` to output all received messages.
type NoopReceiver struct{}

// Receive implements Receiver.
func (NoopReceiver) Receive(context.Context, []rasa.Message) error {
	return nil
}

// Handler implements the http.Handler interface for the callback endpoint.
//
// For more information, see
// https://rasa.com/docs/rasa/connectors/your-own-website#callbackinput
type Handler struct {
	Receiver Receiver
}

// ensure interface
var _ http.Handler = (*Handler)(nil)

// ServeHTTP
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	var err error
	defer func() {
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}()

	// require POST
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// decode JSON
	var messages []rasa.Message
	if err = json.NewDecoder(r.Body).Decode(&messages); err != nil {
		return
	}

	// handle the received messages
	if err = h.Receiver.Receive(r.Context(), messages); err != nil {
		return
	}

	// return "success" for parity with the python sdk
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
