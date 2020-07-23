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

// ValidateSlots TODO
func (c *Context) ValidateSlots(
	disp *action.CollectingDispatcher,
	slots rasa.Slots,
) (events rasa.Events, err error) {
	sc := make(rasa.Slots)

	for slot, value := range slots {
		validator := c.handler.Validator(slot)

		var sset rasa.Slots
		if sset, err = validator.Validate(c, disp, value); err != nil {
			return
		}
		sc.Update(sset)
	}

	// turn validated slots into SlotSet events
	timestamp := rasa.Time(time.Now())
	for key := range sc {
		events = append(events, rasa.SlotSet{
			Timestamp: timestamp,
			Key:       key,
			Value:     sc[key],
		})
	}
	return
}

// Check whether form action should request given slot.
func (c *Context) shouldRequestSlot(slot string) bool {
	return !c.Tracker.Slots.Has(slot)
}

//
func (c *Context) mappersForSlot(slotToFill string) (m Mappers) {
	mappers := c.handler.SlotMappers()
	if m = mappers[slotToFill]; len(mappers) > 0 {
		return
	}

	// fallback to default mapper
	m = Mappers{FromEntity{Entity: slotToFill}}
	return nil
}

// GetEntityValue TODO
func (c *Context) GetEntityValue(
	name, role, group string,
) (value rasa.EntityValue) {
	raw := c.GetLatestEntityValues(name, role, group)
	if len(raw) == 1 {
		value = rasa.StringValue(raw[0])
		return
	}
	value = rasa.SliceValue(raw)
	return
}

// GetLatestEntityValues TODO
func (c *Context) GetLatestEntityValues(
	entity, role, group string,
) (values []string) {
	values = c.Tracker.LatestEntityValues(entity, role, group)
	return
}

//
func (c *Context) extractRequestedSlot() (values rasa.Slots) {
	slotToFill, _ := c.Tracker.Slots[RequestedSlot].(string)
	c.Debugf("trying to extract requested slot [%s]", slotToFill)

	// get mappers
	mappers := c.mappersForSlot(slotToFill)
	for _, mapper := range mappers {
		c.Debugf("got mapping %v", mapper)

		if c.IntentIsDesired(mapper) {
			if value := mapper.Extract(c); value != nil {
				values = rasa.Slots{slotToFill: value}
				return
			}
		}
	}

	// test the mappers for a match
	c.Debugf("failed to extract requested slot [%s]", slotToFill)
	return
}

//
func (c *Context) extractOtherSlots() (values rasa.Slots) {
	slotToFill := c.Tracker.Slots[RequestedSlot]
	values = rasa.Slots{}

	requiredSlots := c.handler.RequiredSlots(c)
	for _, slot := range requiredSlots {
		if slot == slotToFill {
			continue // skip
		}

		mappings := c.mappersForSlot(slot)
		for _, mapping := range mappings {
			_ = mapping
			panic("unimplemented")
		}
	}
	return
}

// IntentIsDesired checks whether user intent matches intent conditions.
func (c *Context) IntentIsDesired(mapping Mapper) bool {
	intent := c.Tracker.LatestMessage.Intent.Name
	return mapping.Desires(intent)
}

// EntityIsDesired checks whether otherSlot should be filled by an entity in the
// input or not.
func (c *Context) EntityIsDesired(mapping Mapper, otherSlot string) bool {
	panic("UNIMPLEMENTED")
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

//  Validation is required if:
//	- the form is active
//	- the form is called after `action_listen`
//	- form validation was not cancelled
func (c *Context) shouldValidate() bool {
	return c.Tracker.LatestActionName == "action_listen" &&
		c.Tracker.HasActiveForm() &&
		c.Tracker.ActiveForm.ShouldValidate()
}
