// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/handle"
)

// Request is the request body for the webhook request.
type Request struct {
	NextAction string        `json:"next_action"`
	SenderID   string        `json:"sender_id"`
	Tracker    *rasa.Tracker `json:"tracker"`
	Domain     *rasa.Domain  `json:"domain"`
	Version    string        `json:"version"`
}

// Response is the response body for the action server in case of
// successful handling of the request.
type Response struct {
	Events    rasa.Events `json:"events"`
	Responses []Message   `json:"responses"`
}

// Server implements http.Handler for the Action Server endpoint.
type Server struct {
	Handlers   map[string]Handler
	PrettyJSON bool
	Logger     Logger
}

// ensure interface
var _ http.Handler = (*Server)(nil)

// NewServer creates a new Server isntance with zero or more initial handlers.
func NewServer(handlers ...Handler) *Server {
	return NewServerWithHandlers(make(map[string]Handler, len(handlers))).RegisterActions(handlers...)
}

// NewServerWithHandlers creates a new Server instance with the provided map as
// it's initial handlers.
func NewServerWithHandlers(handlers map[string]Handler) (s *Server) {
	// verify the sanity of the input
	for action := range handlers {
		if real := handlers[action].ActionName(); real != action {
			panic(fmt.Sprintf("illegal handler, found [%s], expected [%s]", real, action))
		}
	}
	s = &Server{
		Handlers: handlers,
	}
	return
}

// RegisterAction will add the provided Action handler to the server.
//
// Every action should only have a single handler, so the method will panic if
// an attempt is made to register more than one handler for the same action.
//
// This method should only be called *before* the server is started.
func (s *Server) RegisterAction(action Handler) *Server {
	if _, exists := s.Handlers[action.ActionName()]; exists {
		panic(fmt.Sprintf("handler for action [%s] already exists", action.ActionName()))
	}
	s.Handlers[action.ActionName()] = action
	return s
}

// RegisterActions will call RegisterAction for all actions provided to it.
func (s *Server) RegisterActions(actions ...Handler) *Server {
	for i := range actions {
		s.RegisterAction(actions[i])
	}
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/actions"):
		s.withLogs(s.handleActions, w, r)
	case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/health"):
		s.withLogs(s.handleHealth, w, r)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/webhook"):
		s.withLogs(s.handleWebhook, w, r)
	default:
		http.NotFound(w, r) // TODO(ed): log this?
	}
}

// handleWebhook implements the HTTP handler for the /webhook endpoint of the
// action server.
func (s *Server) handleWebhook(ctx context.Context, r *http.Request) (response interface{}, err error) {
	// TODO
	req := Request{}
	defer r.Body.Close()
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		err = &UnmarshalError{cause: err}
		return
	}

	// log action
	s.debugf("[sender: %s - action: %s]", req.SenderID, req.NextAction)

	action := req.NextAction
	handler, exists := s.Handlers[action]
	if !exists || handler == nil {
		err = &MissingHandlerError{action}
		return
	}

	// handle action
	disp := CollectingDispatcher{} // non-nil
	events, err := handler.Run(
		&Context{
			context: ctx,
			logger:  s.Logger,
			Tracker: req.Tracker,
			Domain:  req.Domain,
		},
		&disp,
	)
	if err != nil {
		err = &HandlerError{action, err}
		return
	}
	if events == nil {
		events = rasa.Events{} // non-nil
	}

	// respond
	response = &Response{
		Events:    events,
		Responses: disp,
	}
	return
}

// handleHealth implements the HTTP handler for the /health endpoint of the
// action server.
func (s *Server) handleHealth(ctx context.Context, r *http.Request) (interface{}, error) {
	return struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}, nil
}

// handleActions implements the HTTP handler for the /actions endpoint of the
// action server.
func (s *Server) handleActions(ctx context.Context, r *http.Request) (interface{}, error) {
	// respond
	resp := []string{}
	for key := range s.Handlers {
		resp = append(resp, key)
	}
	sort.Strings(resp)
	return resp, nil
}

// serverError will try to serve an error status and body based on the error, if
// any.
func (s *Server) serveError(w http.ResponseWriter, errp *error) {
	err := *errp
	if err == nil {
		return
	}

	if e, ok := err.(respErr); ok {
		_ = s.serveJSON(w, e.respCode(), struct {
			Error string `json:"error"`
		}{
			Error: e.respBody(),
		})
		return
	}

	// generic error response
	_ = s.serveJSON(w, http.StatusInternalServerError, struct {
		ActionName string `json:"action_name,omitempty"`
		Error      string `json:"error"`
	}{
		// TODO include action_name in the response
		Error: err.Error(),
	})
}

// serveJSON will try to serve the provided src value as a json response.
func (s *Server) serveJSON(w http.ResponseWriter, status int, src interface{}) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	if s.PrettyJSON {
		enc.SetIndent("", "    ")
	}
	err = enc.Encode(src)
	return
}

//
func (s *Server) withLogs(
	fn func(ctx context.Context, r *http.Request) (interface{}, error),
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := s.requestContext(r.Context(), r)
	defer cancel()

	// log
	s.infof("%s - %s START", r.Method, r.URL.String())
	defer s.infof("%s - %s END", r.Method, r.URL.String())

	// ensure error handling
	var err error
	defer s.serveError(w, &err)
	defer handle.Error(&err, func(err error) error {
		s.errorf("%s - %s : %s", r.Method, r.URL.String(), err.Error())
		return err
	})

	resp, err := fn(ctx, r)
	if err != nil {
		return // don't respond, send error
	}

	err = s.serveJSON(w, http.StatusOK, resp)
}

//
func (s *Server) requestContext(ctx context.Context, r *http.Request) (context.Context, context.CancelFunc) {
	// TODO(ed): enhance rasactx.Request using r
	return context.WithTimeout(ctx, time.Second*10) // TODO(ed): extract timeout constant
}

// debugf
func (s *Server) debugf(format string, args ...interface{}) {
	if s.Logger != nil {
		s.Logger.Debugf(format, args...)
	}
}

// infof
func (s *Server) infof(format string, args ...interface{}) {
	if s.Logger != nil {
		s.Logger.Infof(format, args...)
	}
}

// errorf
func (s *Server) errorf(format string, args ...interface{}) {
	if s.Logger != nil {
		s.Logger.Errorf(format, args...)
	}
}
