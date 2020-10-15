// Copyright (c) 2020 Eddy <eddy@scarlet.dev>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package cmd

// TODO(ed): override UnmarshalYAML for SlotYaml and TemplateYaml

// DomainYaml specifies the format of domain.yaml.
type DomainYaml struct {
	Actions  []string
	Entities []string
	Forms    []string
	// Intents  []string
	// Slots     map[string]SlotYaml
	// Templates map[string]TemplateYaml
}

// SlotYaml specifies the format of the slot values in domain.yaml.
type SlotYaml struct {
	Type   string
	Values string
}

// TemplateYaml specifies the format of the template values in domain.yaml.
type TemplateYaml struct {
	Text string
}
