// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package webhooks

import (
	"context"

	"go.scarlet.dev/rasa"
)

// CallbackInput TODO
type CallbackInput struct {
	Sender   string       `json:"sender"`
	Message  string       `json:"message"`
	Metadata rasa.JSONMap `json:"metadata,omitempty"`
}

// Callback implements a client to the Callback webhook.
type Callback struct {
	*baseClient
}

// NewCallback creates a new client to the Callback webhook.
func NewCallback(opts *ClientOpts) *Callback {
	return &Callback{
		baseClient: &baseClient{
			ClientOpts: opts,
			channel:    "callback",
		},
	}
}

// Send TODO
func (c *Callback) Send(ctx context.Context, in *CallbackInput, opts ...RequestOption) (err error) {
	err = c.baseClient.request(ctx, in, nil, opts...)
	return
}
