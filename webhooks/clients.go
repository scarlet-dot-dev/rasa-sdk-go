// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package webhooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"go.scarlet.dev/rasa"
)

// ClientOpts TODO
type ClientOpts struct {
	// Endpoint points to the root of the Rasa API.
	//
	// Defaults to `http://localhost:5005`
	Endpoint string

	// Client is the http.Client instance to use for calls to the Rasa webhooks.
	//
	// If empty, is will default to http.DefaultClient.
	Client *http.Client
}

//
func (c *ClientOpts) endpoint() string {
	if c == nil || c.Endpoint == "" {
		return rasa.DefaultRasaEndpoint
	}
	return strings.TrimSuffix(c.Endpoint, "/")
}

//
func (c *ClientOpts) webhookURL(channel string) string {
	return fmt.Sprintf(
		"%s/webhooks/%s/webhook",
		c.endpoint(),
		channel,
	)
}

//
func (c *ClientOpts) httpClient() *http.Client {
	if c == nil || c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}

// baseClient
type baseClient struct {
	*ClientOpts
	channel string
}

// request TODO
func (c *baseClient) request(ctx context.Context, dst, body interface{}, opts ...RequestOption) (err error) {
	data, err := json.Marshal(body)
	if err != nil {
		return
	}

	r, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.webhookURL(c.channel),
		bytes.NewBuffer(data),
	)
	if err != nil {
		return
	}

	// apply middleware
	for i := range opts {
		opts[i].Apply(ctx, r)
	}

	resp, err := c.httpClient().Do(r)
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
