// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"go.scarlet.dev/rasa"
)

// TestMessage
func TestMessage(t *testing.T) {
	t.Run("JSON marshal", func(t *testing.T) {
		cases := []struct {
			json []byte
			msg  *rasa.Message
		}{
			// normal
			{
				json: []byte(`{"text":"test"}`),
				msg: &rasa.Message{
					Text: "test",
				},
			},
		}

		for i := range cases {
			entry := cases[i]

			marshaled, err := json.Marshal(entry.msg)
			require.NoError(t, err)

			for _, str := range [][]byte{entry.json, marshaled} {
				var result rasa.Message
				err := json.Unmarshal([]byte(str), &result)
				require.NoError(t, err)

				require.Equal(t, entry.msg, &result)
			}
		}
	})
}

// TestCollectingDispatcher
func TestCollectingDispatcher(t *testing.T) {
	t.Run("utter", func(t *testing.T) {
		var disp CollectingDispatcher
		disp.Utter(&rasa.Message{Text: "test"})

		require.Equal(t, CollectingDispatcher{{Text: "test"}}, disp)
		require.EqualValues(t, []rasa.Message{{Text: "test"}}, disp)
	})
}
