package rasa

import (
	"context"
	"net/http"
)

// baseClient
type baseClient struct {
	host   string
	client *http.Client
}

// request TODO
func (c *baseClient) request(ctx context.Context, body, dst interface{}, opts ...RequestOption) {
	var r *http.Request

	// apply middleware
	for i := range opts {
		opts[i].Apply(ctx, r)
	}
}

// RequestOption can be used to inject middleware into the clients / requests.
type RequestOption interface {
	//
	// Inject will be called after the request is built up by the client and
	// before it is sent over the wire.
	//
	// Inject can be used to provide credentials.
	Apply(ctx context.Context, r *http.Request)
}

// RestInput TODO
type RestInput struct {
	Sender   string                 `json:"sender"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CallbackInput TODO
type CallbackInput struct {
	Sender   string                 `json:"sender"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Rest implements a client to the Rest webhook.
type Rest struct {
	*baseClient
}

// Send TODO
func (c *Rest) Send(ctx context.Context, in *RestInput, opts ...RequestOption) (out []Message, err error) {
	c.baseClient.request(ctx, in, &out, opts...)
	return
}

// Callback implements a client to the Callback webhook.
type Callback struct {
	*baseClient
}

// Send TODO
func (c *Callback) Send(ctx context.Context, in *CallbackInput, opts ...RequestOption) (err error) {
	c.baseClient.request(ctx, in, nil, opts...)
	return
}
