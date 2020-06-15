// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//
func init() {
	rootCmd.AddCommand(versionCmd)
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rasagen - v%s\n", rootCmd.Version)
			fmt.Printf("rasa core - >= 1.9.x\n")
		},
	}
)
