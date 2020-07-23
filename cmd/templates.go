// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

import (
	"strings"
	"text/template"
)

//
type tmplConstContext struct {
	Package       string
	Actions       []string
	LicenseHeader string
}

type tmplActioContext struct {
	Package       string
	Action        string
	LicenseHeader string
}

type tmplServeContext struct {
	Package       string
	Actions       []string
	LicenseHeader string
}

// templates
var tmplConst *template.Template
var tmplActio *template.Template
var tmplServe *template.Template

func init() {
	funcs := template.FuncMap{
		"ToCamelCase": ToCamelCase,
		"ToHandlerTypeName": func(s string) string {
			return ToCamelCase(strings.TrimPrefix(s, "action_"))
		},
		"ToFormHandler": func(s string) string {
			return ToCamelCase(strings.TrimPrefix(s, "action_"))
		},
	}
	tmplConst = template.Must(template.New("const").Funcs(funcs).Parse(constants))
	tmplActio = template.Must(template.New("actio").Funcs(funcs).Parse(actionHandler))
	tmplServe = template.Must(template.New("serve").Funcs(funcs).Parse(serverCreate))
}

const (
	actionHandler = `{{ .LicenseHeader }}

package {{ .Package }}

import (
	"context"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// {{ .Action | ToHandlerTypeName }} implements the action.Handler for the
// action "{{ .Action }}".
type {{ .Action | ToHandlerTypeName }} struct {
	// TODO fields here if needed
}

// ensure interface
var _ action.Handler = (*{{ .Action | ToHandlerTypeName }})(nil)

// ActionName implements action.Handler
func ({{ .Action | ToHandlerTypeName }}) ActionName() string {
	return {{ .Action | ToCamelCase }}
}

// Run implements action.Handler.
func (h *{{ .Action | ToHandlerTypeName}}) Run(
	ctx *actions.Context,
	dispatcher *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	// TODO implement the Action Handler
	return
}
`

	constants = `{{ .LicenseHeader }}

package {{ .Package }}

// Action constants
const (
{{- range $element := .Actions }}
{{ $element | ToCamelCase }} = "{{ $element -}}"
{{- end }}
)

// Template constants
// TODO
`

	serverCreate = `{{ .LicenseHeader }}

package {{ .Package }}

import (
	"go.scarlet.dev/rasa/action"
)

// RegisterHandlers will add all generated handlers to the provided action server
func RegisterHandlers(s *action.Server) *action.Server {
	{{ range $element := .Actions -}}
	s.RegisterActions(&{{ $element | ToHandlerTypeName }}{})
	{{ end }}
	return s
}

// NewActionServer creates a new action.Server with all generated Action
// handlers already registered.
func NewActionServer() *action.Server {
	return RegisterHandlers(action.NewServer())
}
`
)
