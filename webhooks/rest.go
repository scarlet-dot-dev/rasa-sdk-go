// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package webhooks

import (
	"context"

	"go.scarlet.dev/rasa"
)

// RestInput TODO
type RestInput struct {
	Sender   string       `json:"sender"`
	Message  string       `json:"message"`
	Metadata rasa.JSONMap `json:"metadata,omitempty"`
}

// Rest implements a client to the Rest webhook.
type Rest struct {
	*baseClient
}

// NewRest creates a new client to the Rest webhook.
func NewRest(opts *ClientOpts) *Rest {
	return &Rest{
		baseClient: &baseClient{
			ClientOpts: opts,
			channel:    "rest",
		},
	}
}

// Send TODO
func (c *Rest) Send(ctx context.Context, in *RestInput, opts ...RequestOption) (out []rasa.Message, err error) {
	err = c.baseClient.request(ctx, &out, in, opts...)
	return
}
