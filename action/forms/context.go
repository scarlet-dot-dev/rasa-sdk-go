package form

import (
	"time"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// Context is provided to the Form action. It wraps an action.Context for
// additional convenience methods.
type Context struct {
	*action.Context
}

// Check whether form action should request given slot.
func (c *Context) shouldRequestSlot(slot string) bool {
	_, ok := c.Tracker.Slots[slot]
	return !ok
}

//
func (c *Context) mappingsForSlot(slotToFill string) []map[string]interface{} {
	// TODO
	return nil
}

// DeactivateForm returns a list of events that tell Rasa to deactivate the
// current form.
func (c *Context) DeactivateForm() []rasa.Event {
	return []rasa.Event{
		rasa.Form{
			Timestamp: rasa.Time(time.Now()),
			Name:      "",
		},
		rasa.SlotSet{
			Timestamp: rasa.Time(time.Now()),
			Key:       RequestedSlot,
			Value:     nil,
		},
	}
}
