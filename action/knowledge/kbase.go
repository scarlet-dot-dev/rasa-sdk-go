package knowledge

import (
	"fmt"
	"math/rand"
)

// Base TODO
type Base interface {
	// GetAttributesOfObject returns a list of all attributes that belong to the
	// provided object type.
	Attributes(objectType string) []string

	// GetKeyAttributeOfObject returns the key attribute for the given object
	// type.
	KeyAttribute(objectType string) string

	// OrdinalMapping
	OrdinalMapping() OrdinalMapper

	// GetObjects queries the knowledge base for objects of the given type.
	// Restrict the objects by the provided attributes, if any attributes are
	// given.
	GetObjects(objectType string, attributes []map[string]string, limit int)

	//
	GetObject(objectType string, identifier string)

	// ForType TODO
	// ForType(object string) Type
}

// Type TODO
type Type interface {
	Attributes() []string
	KeyAttribute() string
	GetObjects(attributes []map[string]string, limit int)
	GetObject(id string)
}

// OrdinalMapper TODO
type OrdinalMapper interface {
	Map(mention string) func(list []string) string
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

// Representer specifies the interface to be implemented by objects returned by
// the KnowledgeBase which need to be represented by the Agent.
type Representer interface {
	Represent() string
}

// ObjectMap implements Representer for simple dynamic objects.
type ObjectMap map[string]interface{}

// Represent implements Representer.
func (m ObjectMap) Represent() string {
	if v, ok := m["name"].(string); ok {
		return v
	}
	return fmt.Sprintf("%v", m)
}

// Entry is the interface implemented by knowledge base entries.
type Entry interface {
	Represent() string

	// Matches
	Matches(key string) bool

	// ID returns the Object ID
	ID() string
}

// InMemory implements Base for a simple, in-memory knowledge base.
type InMemory struct {
	objects  map[string][]Entry
	ordinals OrdinalMapper
}

// OrdinalMapper implements Base.
func (m *InMemory) OrdinalMapper() OrdinalMapper {
	return m.ordinals
}
