// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package callback

import (
	"context"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/webhooks"
	"go.scarlet.dev/rasa/webhooks/internal"
)

// Input TODO
type Input struct {
	Sender   string       `json:"sender"`
	Message  string       `json:"message"`
	Metadata rasa.JSONMap `json:"metadata,omitempty"`
}

// Client implements a client to the Client webhook.
type Client struct {
	*internal.BaseClient
}

// NewCallback creates a new client to the Callback webhook.
func NewCallback(opts *webhooks.ClientOpts) *Client {
	return &Client{
		BaseClient: &internal.BaseClient{
			ClientOpts: opts,
			Channel:    "callback",
		},
	}
}

// Send TODO
func (c *Client) Send(ctx context.Context, in *Input, opts ...webhooks.RequestOption) (err error) {
	err = c.BaseClient.Request(ctx, in, nil, opts...)
	return
}
