// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"regexp"
	"strings"
)

var regex = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

// ToCamelCase is a utility method to turn snake_case identifiers into CamelCase.
func ToCamelCase(s string) string {
	return regex.ReplaceAllStringFunc(s, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
