package form

import (
	"fmt"
	"time"

	"go.scarlet.dev/rasa"
	"go.scarlet.dev/rasa/action"
)

// Context is provided to the Form action. It wraps an action.Context for
// additional convenience methods.
type Context struct {
	*action.Context

	handler Handler
}

// Check whether form action should request given slot.
func (c *Context) shouldRequestSlot(slot string) bool {
	return !c.Tracker.Slots.Has(slot)
}

//
func (c *Context) mappingsForSlot(slotToFill string) []map[string]interface{} {
	// TODO
	return nil
}

// requestNextSlot implements the default routine for determining the next slot
// request.
func (c *Context) requestNextSlot(
	dispatcher *action.CollectingDispatcher,
) (events rasa.Events, err error) {
	required := c.handler.RequiredSlots(c)
	for _, slot := range required {
		if c.shouldRequestSlot(slot) {
			c.Debugf("request next slot [%s]", slot)
			dispatcher.Utter(&action.Message{
				Template: fmt.Sprintf("utter_ask_%s", slot),
				Kwargs:   c.Tracker.Slots,
			})
			events = rasa.Events{
				&rasa.SlotSet{
					Timestamp: rasa.Time(time.Now()),
					Key:       RequestedSlot,
					Value:     slot,
				},
			}
			return
		}
	}

	// no more required slots to fill
	return
}
