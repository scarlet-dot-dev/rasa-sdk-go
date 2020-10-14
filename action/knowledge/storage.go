// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package knowledge

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// Storage TODO
type Storage interface {
	// // GetAttributesOfObject returns a list of all attributes that belong to the
	// // provided object type.
	// Attributes(objectType string) []string

	// // GetKeyAttributeOfObject returns the key attribute for the given object
	// // type.
	// KeyAttribute(objectType string) string

	// OrdinalMapping
	OrdinalMapping() OrdinalMapper

	// // GetObjects queries the knowledge base for objects of the given type.
	// // Restrict the objects by the provided attributes, if any attributes are
	// // given.
	// GetObjects(objectType string, attributes []map[string]string, limit int)

	// //
	// GetObject(objectType string, identifier string)

	// ForType TODO
	ForType(object string) (ObjectType, error)
}

// Attribute TODO
type Attribute struct {
	Name  string
	Value string
}

// Object TODO
type Object interface {
	// Type returns the name of the object type.
	// Type() string

	// Key returns the value of the id attribute.
	Key() string

	// Representation returns the human readable name of the object.
	Representation() string

	// Attribute extracts the attribute with the given name from the Object.
	Attribute(name string) string
}

// ObjectType TODO
type ObjectType interface {
	TypeName() string
	Attributes() []string
	KeyAttribute() string
	GetObjects(attributes []Attribute, limit int) []Object
	GetObject(id string) Object
	// OrdinalMapper() OrdinalMapper
}

// Representer specifies the interface to be implemented by objects returned by
// the KnowledgeBase which need to be represented by the Agent.
type Representer interface {
	Represent() string
}

// ObjectMap implements Representer for simple dynamic objects.
type ObjectMap map[string]interface{}

// Entry implements Object for the InMemory knowledgebase
type Entry map[string]string

// TODO(ed):
var _ Object = (Entry)(nil)

// Key implements Object.
func (e Entry) Key() string {
	return e["id"]
}

// Representation implements Object.
func (e Entry) Representation() string {
	if v, ok := e["name"]; ok {
		return v
	}
	return fmt.Sprintf("%v", e)
}

// Attribute implements Object.
func (e Entry) Attribute(name string) string {
	return e[name]
}

// InMemory implements Base for a simple, in-memory knowledge base.
type InMemory struct {
	// objects holds the
	objects  map[string][]Entry
	ordinals OrdinalMapper
}

// NewInMemory creates a new InMemory knowledge base by parsing the provided
// io.Reader as a stream of JSON.
func NewInMemory(r io.Reader) (kb *InMemory, err error) {
	kb = &InMemory{}
	if err = json.NewDecoder(r).Decode(&kb.objects); err != nil {
		return
	}

	// for each type, sort slice
	for otype := range kb.objects {
		entries := kb.objects[otype]
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Key() < entries[j].Key()
		})
	}

	return
}

// WithOrdinalMapper TODO
func (m *InMemory) WithOrdinalMapper(om OrdinalMapper) *InMemory {
	m.ordinals = om
	return m
}

// OrdinalMapper implements Base.
func (m *InMemory) OrdinalMapper() OrdinalMapper {
	return m.ordinals
}

// InMemoryType implements ObjectType for InMemory storage.
type InMemoryType struct {
	store *InMemory
	otype string
}
