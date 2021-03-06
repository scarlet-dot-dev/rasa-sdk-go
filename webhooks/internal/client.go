// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"go.scarlet.dev/rasa/webhooks"
)

// BaseClient implements an embeddable base for webhook clients.
type BaseClient struct {
	*webhooks.ClientOpts
	Channel string
}

// Request performs a single POST request to a webhook endpoint.
func (c *BaseClient) Request(ctx context.Context, dst, body interface{}, opts ...webhooks.RequestOption) (err error) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.WebhookURL(c.Channel),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return
	}

	// apply middleware
	for i := range opts {
		opts[i].Apply(ctx, r)
	}

	resp, err := c.HTTPClient().Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("request code [%s]", resp.Status)
		return
	}

	if dst == nil {
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(dst); err != nil {
		return
	}

	return
}
