// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

// Logger provides an interface for logging provided by the package itself.
//
// The Logger is provided
type Logger interface {
	// Debugf is used to output debug (verbose) information. Logs of debug level
	// are useful during development, but should be disabled in a production
	// environment.
	Debugf(format string, args ...interface{})

	// Infof is used to log general information. It is the default logging
	// level.
	Infof(format string, args ...interface{})

	// Warnf is used to log warnings.
	Warnf(format string, args ...interface{})

	// Errorf is used to log errors.
	Errorf(format string, args ...interface{})
}
