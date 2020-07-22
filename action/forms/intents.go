package form

import "encoding/json"

// IntentList is a slice type of strings, where each strings refers to a named
// intent.
type IntentList []string

// IntentList implements IntentLister.
func (l IntentList) IntentList() IntentList {
	return l
}

// Intent is a string type wrapper that implements IntentLister to allow a
// single intent to be passed where a list of intents is requested.
type Intent string

// IntentList implements IntentLister.
func (l Intent) IntentList() IntentList {
	return IntentList{string(l)}
}

// MarshalJSON implements json.Marshaler.
func (l Intent) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.IntentList())
}

// IntentLister is an interface to provide a union over Intent and IntentList,
// to allow treating both as an IntentList.
type IntentLister interface {
	IntentList() IntentList
}

// ensure interface
var _ json.Marshaler = (*Intent)(nil)
