// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
)

// Context contains request specific context, such as a logger for the specific
// request with tracing information.
type Context struct {
	Logger Logger
}

// unique context key for Request values.
type requestContextKey struct{}

// WithContext returns a new context with the additional Request information.
func WithContext(ctx context.Context, r *Context) context.Context {
	return context.WithValue(ctx, requestContextKey{}, r)
}

// ContextFrom returns the additional Request information from the Context.
func ContextFrom(ctx context.Context) *Context {
	r, _ := ctx.Value(requestContextKey{}).(*Context)
	return r
}
