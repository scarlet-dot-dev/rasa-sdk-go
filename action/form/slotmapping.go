// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package form

// Mapper represents a serializable slot mapping object.
type Mapper interface {
	//
	Desires(intent string) bool

	//
	Extract(ctx *Context) interface{}
}

// ensure interfaces
var _ Mapper = (*FromEntity)(nil)
var _ Mapper = (*FromTriggerIntent)(nil)
var _ Mapper = (*FromIntent)(nil)
var _ Mapper = (*FromText)(nil)

// SlotMapping is a list of Mappers to map an individual slot.
type SlotMapping []Mapper

// SlotMappings is an alias type wrapping a mapping of slots to their
// SlotMapping configurations.
type SlotMappings map[string]SlotMapping

// Mapping returns the SlotMapping configuration for the provided slot.
func (m SlotMappings) Mapping(slot string) SlotMapping {
	if m != nil {
		if ret, ok := m[slot]; ok && len(ret) > 0 {
			return ret
		}
	}

	// return the default
	return SlotMapping{FromEntity{Entity: slot}}
}

// FromEntity TODO
type FromEntity struct {
	Entity  string
	Intents IntentFilter
	Role    string
	Group   string
}

// Desires implements Mapper.
func (m FromEntity) Desires(intent string) bool {
	return m.Intents == nil || m.Intents.Desires(intent)
}

// Extract implements Mapper.
func (m FromEntity) Extract(ctx *Context) interface{} {
	return ctx.EntityValue(m.Entity, m.Role, m.Group)
}

// FromTriggerIntent TODO
type FromTriggerIntent struct {
	Value   interface{}
	Intents IntentFilter
}

// Desires implements Mapper.
func (m FromTriggerIntent) Desires(intent string) bool {
	return m.Intents == nil || m.Intents.Desires(intent)
}

// Extract implements Mapper.
func (m FromTriggerIntent) Extract(ctx *Context) interface{} {
	// return nothing - Extract is used for requested slots, trigger_intent is
	// used on form activation only.
	return nil
}

// FromIntent TODO
type FromIntent struct {
	Value   interface{}
	Intents IntentFilter
}

// Desires implements Mapper.
func (m FromIntent) Desires(intent string) bool {
	return m.Intents == nil || m.Intents.Desires(intent)
}

// Extract implements Mapper.
func (m FromIntent) Extract(ctx *Context) interface{} {
	return m.Value
}

// FromText TODO
type FromText struct {
	Intents IntentFilter
}

// Desires implements Mapper.
func (m FromText) Desires(intent string) bool {
	return m.Intents == nil || m.Intents.Desires(intent)
}

// Extract implements Mapper.
func (m FromText) Extract(ctx *Context) interface{} {
	return ctx.Tracker.LatestMessage.Text
}
