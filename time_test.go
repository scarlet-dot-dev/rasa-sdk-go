// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestTime
func TestTime(t *testing.T) {
	t.Run("JSON", func(t *testing.T) {
		cases := []struct {
			input  time.Time // marshal input
			output time.Time // unmarshal output
			json   []byte    // marshal result
		}{
			{
				input:  time.Time{},
				output: time.Time{},
				json:   []byte("0"),
			},
			{
				// the 0-time is treated the same as Time{}, considering there
				// will be no new events from 1970.
				input:  time.Unix(0, 0),
				output: time.Time{},
				json:   []byte("0"),
			},
			{
				// see comment of previous case
				input:  time.Unix(0, 123456789),
				output: time.Time{},
				json:   []byte("0"),
			},
			{
				input:  time.Unix(1, 0),
				output: time.Unix(1, 0),
				json:   []byte("1"),
			},
			{
				input:  time.Unix(1234567890, 0),
				output: time.Unix(1234567890, 0),
				json:   []byte("1234567890"),
			},
			{
				input:  time.Unix(1234567890, 123456789),
				output: time.Unix(1234567890, 0),
				json:   []byte("1234567890"),
			},
		}

		for i := range cases {
			entry := cases[i]

			input := Time(entry.input)
			ser, err := json.Marshal(input)
			require.NoErrorf(t, err, "failed on %d", i)
			require.Equalf(t, entry.json, ser, "failed on %d", i)

			var result Time
			err = json.Unmarshal(ser, &result)
			require.NoErrorf(t, err, "failed on %d", i)
			require.Equalf(t, Time(entry.output), result, "failed on %d", i)
		}
	})
}
