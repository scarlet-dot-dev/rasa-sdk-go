package form

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"go.scarlet.dev/rasa/handle"
)

// MapperType TODO
type MapperType string

// TODO
const (
	MapperTypeFromEntity        = MapperType("from_entity")
	MapperTypeFromTriggerIntent = MapperType("from_trigger_intent")
	MapperTypeFromIntent        = MapperType("from_intent")
	MapperTypeFromText          = MapperType("from_text")
)

// Mapper represents a serializable slot mapping object.
type Mapper interface {
	// Type returns the constant representing the slot mapping type in a
	// marshalled JSON object.
	Type() MapperType
}

// Mappers implements JSON Unmarshaling of Rasa's slot mapper types.
type Mappers []Mapper

// ensure interface
var _ json.Unmarshaler = (*Mappers)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (l *Mappers) UnmarshalJSON(data []byte) (err error) {
	defer handle.Error(&err, func(err error) error {
		return errors.WithMessage(err, "unable to unmarshal Mappers")
	})

	// initialize as empty mappers list
	*l = Mappers{}

	// get the Raw messages
	var mappers []json.RawMessage
	if err = json.Unmarshal(data, &mappers); err != nil {
		return
	}

	for i := range mappers {
		var mapper Mapper
		if mapper, err = unmarshalSlotMapper(mappers[i]); err != nil {
			return
		}
		*l = append(*l, mapper)
	}

	return
}

// unmarshalSlotMapper will unmarshal the provided serialized JSON as an Event.
func unmarshalSlotMapper(data []byte) (mapper Mapper, err error) {
	var marker struct {
		Type MapperType `json:"type"`
	}
	if err = json.Unmarshal(data, &marker); err != nil {
		return
	}

	switch marker.Type {
	default:
		// error case - unknown of unsupported event type
		err = fmt.Errorf("invalid slot mapper type [%s]", marker.Type)
		return
	case MapperTypeFromEntity:
		mapper = new(FromEntity)
	case MapperTypeFromTriggerIntent:
		mapper = new(FromTriggerIntent)
	case MapperTypeFromIntent:
		mapper = new(FromIntent)
	case MapperTypeFromText:
		mapper = new(FromText)
	}

	// unmarshal the JSON event
	err = json.Unmarshal(data, mapper)
	return
}

// structToMap implements a single-layer conversion from a struct to a
// string-indexed map.
func structToMap(s interface{}) map[string]interface{} {
	// TODO(ed): replace this with a self-maintained implementation?
	st := structs.New(s)
	st.TagName = "json"
	return st.Map()
}

// marshalMapper will marshal the provided Mapper value into its JSON
// representation.
func marshalMapper(e Mapper) ([]byte, error) {
	// turn the event into a map
	result := structToMap(e)

	// add the event type property
	result["type"] = e.Type()
	return json.Marshal(result)
}

// FromEntity TODO
type FromEntity struct {
	Entity    string     `json:"entity"`
	Intent    IntentList `json:"intent,omitempty"`
	NotIntent IntentList `json:"not_intent,omitempty"`
	Role      string     `json:"role,omitempty"`
	Group     string     `json:"group,omitempty"`
}

// FromTriggerIntent TODO
type FromTriggerIntent struct {
	Value     interface{} `json:"value"`
	Intent    IntentList  `json:"intent,omitempty"`
	NotIntent IntentList  `json:"not_intent,omitempty"`
}

// FromIntent TODO
type FromIntent struct {
	Value     interface{} `json:"value"`
	Intent    IntentList  `json:"intent,omitempty"`
	NotIntent IntentList  `json:"not_intent,omitempty"`
}

// FromText TODO
type FromText struct {
	Intent    IntentList `json:"intent,omitempty"`
	NotIntent IntentList `json:"not_intent,omitempty"`
}

// Type impements SlotMapper.
func (FromEntity) Type() MapperType { return MapperTypeFromEntity }

// Type impements SlotMapper.
func (FromTriggerIntent) Type() MapperType { return MapperTypeFromTriggerIntent }

// Type impements SlotMapper.
func (FromIntent) Type() MapperType { return MapperTypeFromIntent }

// Type impements SlotMapper.
func (FromText) Type() MapperType { return MapperTypeFromText }

// MarshalJSON implements json.Marshaler.
func (m *FromEntity) MarshalJSON() ([]byte, error) { return marshalMapper(m) }

// MarshalJSON implements json.Marshaler.
func (m *FromTriggerIntent) MarshalJSON() ([]byte, error) { return marshalMapper(m) }

// MarshalJSON implements json.Marshaler.
func (m *FromIntent) MarshalJSON() ([]byte, error) { return marshalMapper(m) }

// MarshalJSON implements json.Marshaler.
func (m *FromText) MarshalJSON() ([]byte, error) { return marshalMapper(m) }

// ensure interfaces
var _ Mapper = (*FromEntity)(nil)
var _ Mapper = (*FromTriggerIntent)(nil)
var _ Mapper = (*FromIntent)(nil)
var _ Mapper = (*FromText)(nil)

// ensure interfaces
var _ json.Marshaler = (*FromEntity)(nil)
var _ json.Marshaler = (*FromTriggerIntent)(nil)
var _ json.Marshaler = (*FromIntent)(nil)
var _ json.Marshaler = (*FromText)(nil)
