package form

import "encoding/json"

// Intents is a slice type of strings, where each strings refers to a named
// intent.
type Intents []string

// IntentList implements IntentLister.
func (l Intents) IntentList() Intents {
	return l
}

// Contains TODO
func (l Intents) Contains(intent string) bool {
	for i := range l {
		if l[i] == intent {
			return true
		}
	}
	return false
}

// Intent is a string type wrapper that implements IntentLister to allow a
// single intent to be passed where a list of intents is requested.
type Intent string

// IntentList implements IntentLister.
func (l Intent) IntentList() Intents {
	return Intents{string(l)}
}

// MarshalJSON implements json.Marshaler.
func (l Intent) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.IntentList())
}

// IntentLister is an interface to provide a union over Intent and IntentList,
// to allow treating both as an IntentList.
type IntentLister interface {
	IntentList() Intents
}

// ensure interface
var _ json.Marshaler = (*Intent)(nil)

// IntentFilter TODO
type IntentFilter interface {
	// Desires TODO
	Desires(intent string) bool
}

// AllowIntents implements an IntentFilter based on an allowlist.
type AllowIntents struct {
	IntentLister
}

// Desires implements IntentFilter.
func (a AllowIntents) Desires(intent string) bool {
	return a.IntentList().Contains(intent)
}

// BlockIntents implements an IntentFilter based on a blocklist.
type BlockIntents struct {
	IntentLister
}

// Desires implements IntentFilter.
func (b BlockIntents) Desires(intent string) bool {
	return !b.IntentList().Contains(intent)
}