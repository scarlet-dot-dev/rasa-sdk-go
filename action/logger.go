// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package action

import "context"

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

// // debug alias for a potentially nil logger.
// func ldebug(l Logger, format string, args ...interface{}) {
// 	if l != nil {
// 		l.Debugf(format, args...)
// 	}
// }

// // info alias for a potentially nil logger.
// func linfo(l Logger, format string, args ...interface{}) {
// 	if l != nil {
// 		l.Infof(format, args...)
// 	}
// }

// // warn alias for a potentially nil logger.
// func lwarn(l Logger, format string, args ...interface{}) {
// 	if l != nil {
// 		l.Warnf(format, args...)
// 	}
// }

// // error alias for a potentially nil logger.
// func lerror(l Logger, format string, args ...interface{}) {
// 	if l != nil {
// 		l.Errorf(format, args...)
// 	}
// }

// LogDebugf will extract the request logger from the Context, and log the
// message if a logger is found.
//
// It will panic if the Context does not contain a Request.
func LogDebugf(ctx context.Context, format string, args ...interface{}) {
	r := ContextFrom(ctx)
	if r == nil {
		panic("no Request logger present")
	}

	if l := r.Logger; l != nil {
		l.Debugf(format, args...)
	}
}

// LogInfof will extract the request logger from the Context, and log the
// message if a logger is found.
//
// It will panic if the Context does not contain a Request.
func LogInfof(ctx context.Context, format string, args ...interface{}) {
	r := ContextFrom(ctx)
	if r == nil {
		panic("no Request logger present")
	}

	if l := r.Logger; l != nil {
		l.Infof(format, args...)
	}
}

// LogWarnf will extract the request logger from the Context, and log the
// message if a logger is found.
//
// It will panic if the Context does not contain a Request.
func LogWarnf(ctx context.Context, format string, args ...interface{}) {
	r := ContextFrom(ctx)
	if r == nil {
		panic("no Request logger present")
	}

	if l := r.Logger; l != nil {
		l.Warnf(format, args...)
	}
}

// LogErrorf will extract the request logger from the Context, and log the
// message if a logger is found.
//
// It will panic if the Context does not contain a Request.
func LogErrorf(ctx context.Context, format string, args ...interface{}) {
	r := ContextFrom(ctx)
	if r == nil {
		panic("no Request logger present")
	}

	if l := r.Logger; l != nil {
		l.Errorf(format, args...)
	}
}
