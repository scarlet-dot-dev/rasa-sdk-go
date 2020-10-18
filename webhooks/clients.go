// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package webhooks

import (
	"fmt"
	"net/http"
	"strings"

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

// WebhookURL returns the url for the webhook of the provided channel.
func (c *ClientOpts) WebhookURL(channel string) string {
	return fmt.Sprintf(
		"%s/webhooks/%s/webhook",
		c.endpoint(),
		channel,
	)
}

// HTTPClient returns the configured http.Client, or the default if none is
// provided.
func (c *ClientOpts) HTTPClient() *http.Client {
	if c == nil || c.Client == nil {
		return http.DefaultClient
	}
	return c.Client
}
