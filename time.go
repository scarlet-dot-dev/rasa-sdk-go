// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package rasa

import (
	"encoding/json"
	"time"
)

// zt is the undefined Zero Time variable
var zt time.Time

// Time provides a type alias around time.Time which serializes as an int64 for
// proper interaction with Rasa.
//
// During both serializing and deserializing, all sub-second precision is
// dropped.
type Time time.Time

// ensure interfaces
var _ json.Marshaler = (Time{})
var _ json.Unmarshaler = (*Time)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var seconds float64
	if err = json.Unmarshal(data, &seconds); err != nil {
		return
	}

	if seconds == 0 {
		// init to zero value
		*t = Time{}
		return
	}

	*t = Time(time.Unix(int64(seconds), 0))
	return nil
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	// only set if t has been set to a non-zero value
	if tt := time.Time(t); !tt.IsZero() && tt != zt {
		return json.Marshal(tt.Unix())
	}

	// default serialize as 0, in stead of "nothing".
	return json.Marshal(0)
}
