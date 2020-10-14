// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import "math/rand"

// OrdinalMapper TODO
type OrdinalMapper interface {
	Map(mention string) func(list []string) string
	Has(mention string) bool
}

//
func get(list []string, index int) string {
	if len(list) <= index || index < 0 {
		return ""
	}
	return list[index]
}

// DefaultOrdinalMapper TODO
func DefaultOrdinalMapper() OrdinalMapper {
	return &defaultOrdinalMapper{
		"1":    func(list []string) string { return get(list, 0) },
		"2":    func(list []string) string { return get(list, 1) },
		"3":    func(list []string) string { return get(list, 2) },
		"4":    func(list []string) string { return get(list, 3) },
		"5":    func(list []string) string { return get(list, 4) },
		"6":    func(list []string) string { return get(list, 5) },
		"7":    func(list []string) string { return get(list, 6) },
		"8":    func(list []string) string { return get(list, 7) },
		"9":    func(list []string) string { return get(list, 8) },
		"10":   func(list []string) string { return get(list, 9) },
		"LAST": func(list []string) string { return get(list, len(list)-1) },
		"ANY":  func(list []string) string { return get(list, rand.Intn(len(list))) },
	}
}

// defaultOrdinalMapper is the simplest implementation of OrdinalMapper.
type defaultOrdinalMapper map[string]func(list []string) string

// ensure interface
var _ OrdinalMapper = (defaultOrdinalMapper)(nil)

// Map implements OrdinalMapper
func (m defaultOrdinalMapper) Map(mention string) func([]string) string {
	return m[mention]
}

// Has implements ordinalMapper
func (m defaultOrdinalMapper) Has(mention string) bool {
	_, ok := m[mention]
	return ok
}
