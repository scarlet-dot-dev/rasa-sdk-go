// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToCamelCase(t *testing.T) {
	cases := []struct {
		input  string
		expect string
	}{
		{"", ""},
		{"action_handler_name", "ActionHandlerName"},
		{"simple", "Simple"},
		{"Equal", "Equal"},
		{"EqualEqual", "EqualEqual"},
		{"sOmes_t_UFf", "SOmesTUFf"},
	}

	for i := range cases {
		entry := cases[i]
		result := ToCamelCase(entry.input)
		require.Equal(t, entry.expect, result)
	}
}
