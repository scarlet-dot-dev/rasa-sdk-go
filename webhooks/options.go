// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package webhooks

import (
	"context"
	"fmt"
	"net/http"
)

// RequestOption can be used to inject middleware into the clients / requests.
type RequestOption interface {
	//
	// Inject will be called after the request is built up by the client and
	// before it is sent over the wire.
	//
	// Inject can be used to provide credentials.
	Apply(ctx context.Context, r *http.Request)
}

// ensure interface
var _ RequestOption = (*BearerAuthOption)(nil)

// BearerAuthOption will apply a JWT to the request.
//
// The header "Authorization" will be used, with "Bearer <JWT>" as its value.
type BearerAuthOption struct {
	// Multivalue allows the header to not override, but instead append to the
	// Authorization header.
	Multivalue bool

	// Token should hold the JWT token.
	Token string
}

// Apply implements RequestOption.
func (o *BearerAuthOption) Apply(ctx context.Context, r *http.Request) {
	if o.Multivalue {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", o.Token))
	} else {
		r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.Token))
	}
}
