// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.scarlet.dev/rasa/internal/handle"
	"gopkg.in/yaml.v2"
)

//
func init() {
	actionsCmd.Flags().StringVar(
		&actionsCtx.Package,
		"package",
		"actions",
		"name of the go package",
	)
	actionsCmd.Flags().StringVar(
		&actionsCtx.OutDir,
		"outdir",
		".",
		"output directory for the generated files",
	)
	actionsCmd.Flags().StringVar(
		&actionsCtx.License,
		"license",
		"",
		"string containing the license header for the generated files",
	)

	actionsCmd.MarkFlagDirname("outdir")
	rootCmd.AddCommand(actionsCmd)
}

var (
	actionsCtx struct {
		Package string
		OutDir  string
		License string
	}

	actionsCmd = &cobra.Command{
		Use:   "actions",
		Short: "generate boilerplate for a custom action server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer handle.Error(&err, func(err error) error {
				fmt.Println("encountered an error: " + err.Error())
				return nil
			})

			// Load the domain file
			var config DomainYaml
			func() {
				log.Printf("loading domain from %s\n", rootCtx.DomainFile)
				file, err := os.Open(rootCtx.DomainFile)
				handle.Check(err)
				defer file.Close()

				err = yaml.NewDecoder(file).Decode(&config)
				handle.Check(err)
			}()

			// TODO(ed): make sure ALL config slices are sorted
			sort.Sort(sort.StringSlice(config.Actions))

			// createFile is a small utility closure for executing templates
			createFile := func(tmpl *template.Template, filename string, ctx interface{}) {
				log.Printf("creating file [%s]\n", filename)

				var buff bytes.Buffer

				err = tmpl.Execute(&buff, &ctx)
				handle.Check(err)

				output, err := format.Source(buff.Bytes())
				err = errors.WithMessage(err, "formatting failed for "+filename)
				handle.Check(err)

				err = ioutil.WriteFile(filepath.Join(actionsCtx.OutDir, filename), output, 0666)
				handle.Check(err)
			}

			// generate constants
			log.Printf("printing constants to %s/consts.go\n", actionsCtx.OutDir)
			createFile(
				tmplConst,
				"constants.go",
				&tmplConstContext{
					Package:       actionsCtx.Package,
					Actions:       config.Actions,
					LicenseHeader: actionsCtx.License,
				},
			)

			filtered := make([]string, 0)
			for _, action := range config.Actions {
				if strings.HasPrefix(action, "action_") {
					filtered = append(filtered, action)
				}
			}

			// generate the action handlers
			for _, action := range filtered {
				createFile(tmplActio, action+".go", &tmplActioContext{
					Package:       actionsCtx.Package,
					Action:        action,
					LicenseHeader: actionsCtx.License,
				})
			}

			// generate the server boilerplate
			createFile(tmplServe, "server.go", &tmplServeContext{
				Package:       actionsCtx.Package,
				Actions:       filtered,
				LicenseHeader: actionsCtx.License,
			})

			return
		},
	}
)
