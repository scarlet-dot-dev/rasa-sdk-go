// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvent(t *testing.T) {
	t.Run("JSON unmarshal", func(t *testing.T) {

	})

	t.Run("JSON marshal", func(t *testing.T) {

	})
}

func TestEventMarshalJSON(t *testing.T) {
	// TODO(ed): test if for all Event types, the `Event` field is properly set
	// for serialization.
}

// TestEventList
func TestEventList(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cases := []Events{
			{},
			{&SlotSet{}},
		}

		for i := range cases {
			entry := cases[i]

			ser, err := json.Marshal(entry)
			require.NoError(t, err)

			var result Events
			err = json.Unmarshal(ser, &result)
			require.NoError(t, err)
			require.Equal(t, len(entry), len(result))

			for i := range result {
				require.EqualValues(t, entry[i], result[i])
			}
		}
	})
}
