// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.scarlet.dev/rasa"
)

type testHandler1 struct{}
type testHandlerNoEvent struct{}
type testHandlerNoDispatch struct{}
type testHandlerErr struct{}

func (testHandler1) ActionName() string          { return "action_test" }
func (testHandlerNoEvent) ActionName() string    { return "action_no_event" }
func (testHandlerNoDispatch) ActionName() string { return "action_no_dispatch" }
func (testHandlerErr) ActionName() string        { return "action_error" }

func (testHandler1) Run(ctx Context, dispatcher *CollectingDispatcher) (events rasa.Events, err error) {
	dispatcher.Utter(&rasa.Message{
		Text: "test string",
	})
	events = append(events, &rasa.SlotSet{
		Key:   "test",
		Value: "420",
	})
	return
}

func (testHandlerNoEvent) Run(ctx Context, dispatcher *CollectingDispatcher) (events rasa.Events, err error) {
	dispatcher.Utter(&rasa.Message{
		Text: "test string",
	})
	return
}

func (testHandlerNoDispatch) Run(ctx Context, dispatcher *CollectingDispatcher) (events rasa.Events, err error) {
	events = append(events, &rasa.SlotSet{
		Key:   "test",
		Value: "420",
	})
	return
}

func (testHandlerErr) Run(ctx Context, dispatcher *CollectingDispatcher) (events rasa.Events, err error) {
	err = errors.New("test error")
	return
}

func TestServer(t *testing.T) {
	handler := NewServer(
		&testHandler1{},
		&testHandlerNoEvent{},
		&testHandlerNoDispatch{},
		&testHandlerErr{},
	)
	testRequest := func(
		t *testing.T,
		method, url string,
		status int,
		body interface{},
		result interface{},
	) {
		var bodyReader io.Reader
		if body != nil {
			var buff bytes.Buffer
			err := json.NewEncoder(&buff).Encode(&body)
			require.NoError(t, err)
			bodyReader = &buff
		}

		req := httptest.NewRequest(method, url, bodyReader)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		if !assert.Equal(t, status, resp.StatusCode) {
			t.Logf("body: [%s]", w.Body.String())
			t.FailNow()
		}
		require.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		if result != nil {
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		}
	}

	//
	t.Run("endpoint /actions", func(t *testing.T) {
		expect := []string{
			"action_test",
			"action_no_event",
			"action_no_dispatch",
			"action_error",
		}

		var result []string
		testRequest(
			t,
			"GET",
			"https://example.com/actions",
			http.StatusOK,
			nil,
			&result,
		)

		// compare equality on the sorted slices
		require.ElementsMatch(t, expect, result)
	})

	//
	t.Run("endpoint /health", func(t *testing.T) {
		type status struct {
			Status string `json:"status"`
		}
		expect := status{"OK"}

		var result status
		testRequest(
			t,
			"GET",
			"https://example.com/health",
			http.StatusOK,
			nil,
			&result,
		)

		require.Equal(t, expect, result)
	})

	//
	t.Run("endpoint /webhook", func(t *testing.T) {
		t.Run("action_test", func(t *testing.T) {
			expect := Response{
				Responses: []rasa.Message{{
					Text: "test string",
				}},
				Events: rasa.Events{
					&rasa.SlotSet{
						Key:   "test",
						Value: "420",
					},
				},
			}

			body := &Request{
				NextAction: "action_test",
				SenderID:   "TODO",
				Domain:     nil,    // TODO
				Tracker:    nil,    // TODO
				Version:    "TODO", // TODO
			}

			var result Response
			testRequest(
				t,
				"POST",
				"https://example.com/webhook",
				http.StatusOK,
				body,
				&result,
			)

			// compare equality on the sorted slices
			require.EqualValues(t, expect, result)
		})

		t.Run("action_no_event", func(t *testing.T) {
			expect := Response{
				Responses: []rasa.Message{{
					Text: "test string",
				}},
				Events: rasa.Events{}, // Rasa expects an empty, non-nil list of events.
			}

			body := &Request{
				NextAction: "action_no_event",
				SenderID:   "TODO",
				Domain:     nil,    // TODO
				Tracker:    nil,    // TODO
				Version:    "TODO", // TODO
			}

			var result Response
			testRequest(
				t,
				"POST",
				"https://example.com/webhook",
				http.StatusOK,
				body,
				&result,
			)

			// compare equality on the sorted slices
			require.Equal(t, expect, result)
			require.NotNil(t, result.Events)
		})
		t.Run("action_no_dispatch", func(t *testing.T) {
			expect := Response{
				Responses: []rasa.Message{}, // Rasa expects non-nil values
				Events: rasa.Events{
					&rasa.SlotSet{
						Key:   "test",
						Value: "420",
					},
				},
			}

			body := &Request{
				NextAction: "action_no_dispatch",
				SenderID:   "TODO",
				Domain:     nil,    // TODO
				Tracker:    nil,    // TODO
				Version:    "TODO", // TODO
			}

			var result Response
			testRequest(
				t,
				"POST",
				"https://example.com/webhook",
				http.StatusOK,
				body,
				&result,
			)

			// compare equality on the sorted slices
			require.EqualValues(t, expect, result)
		})
		t.Run("action_error", func(t *testing.T) {
		})

		t.Run("missing handler", func(t *testing.T) {
			body := &Request{
				NextAction: "action_does_not_exist",
				SenderID:   "TODO",
				Domain:     nil,    // TODO
				Tracker:    nil,    // TODO
				Version:    "TODO", // TODO
			}

			testRequest(
				t,
				"POST",
				"https://example.com/webhook",
				http.StatusInternalServerError,
				body,
				nil,
			)
		})
	})
}
