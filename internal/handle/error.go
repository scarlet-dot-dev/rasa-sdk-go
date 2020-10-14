// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package handle

// Error will call fn with error if *e != nil.
func Error(e *error, fn func(err error) error) {
	// first check for a panic
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			panic(err)
		}
		*e = fn(err)
		return
	}

	// not a panic / recover
	if err := *e; err != nil {
		*e = fn(err)
	}
}

// Check error.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
