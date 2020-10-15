// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

const (
	// DefaultServerPort is the default port for the Rasa Action Server.
	//
	// This port is used for both the `/webhook` and `/nlg` endpoints.
	DefaultServerPort = "5055"

	// DefaultRasaAPIPort is the default port for the Rasa API.
	//
	// Used by the webhooks package.
	DefaultRasaAPIPort = "5005"

	// DefaultRasaEndpoint is the default endpoint used by webhook clients.
	DefaultRasaEndpoint = "http://localhost:5005"
)
