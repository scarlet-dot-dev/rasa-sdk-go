package action

import "go.scarlet.dev/rasa"

// Handler specifies the interface implemented by action handlers.
type Handler interface {
	// ActionName returns the name of the action for which the instance
	// implements handling.
	ActionName() string

	// Run should process the action.
	//
	// The Context passed to Run is the context of the HTTP request sent by
	// Rasa's engine.
	Run(
		ctx *Context,
		dispatcher *CollectingDispatcher,
	) (events rasa.Events, err error)
}
