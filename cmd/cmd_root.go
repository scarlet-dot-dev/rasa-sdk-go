// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCtx struct {
		// Used for flags.
		DomainFile string
	}

	rootCmd = &cobra.Command{
		Version: "1.0.0",
		Use:     "rasagen",
		Short:   "A code generator for Rasa based Applications",
		Long: `Rasagen is a command line utility for code generation of constants
and boilerplate for applications written with the rasa-sdk-go.

Note that rasagen was written for the purpose of reducing the time spent
writing boilerplate code for custom action handlers. It is not written to
integrate with go:generate.`,
	}
)

//
func init() {
	rootCmd.PersistentFlags().StringVar(
		&rootCtx.DomainFile,
		"domain",
		"domain.yml",
		"Path to the rasa domain.yaml configuration file",
	)

	rootCmd.MarkPersistentFlagFilename("domain")
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
