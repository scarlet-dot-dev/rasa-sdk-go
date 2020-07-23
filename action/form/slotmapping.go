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

// Mappers implements JSON Unmarshaling of Rasa's slot mapper types.
type Mappers []Mapper

// FromEntity TODO
type FromEntity struct {
	Entity string
	Intent IntentFilter
	Role   string
	Group  string
}

// Desires implements Mapper.
func (m FromEntity) Desires(intent string) bool {
	return m.Intent == nil || m.Intent.Desires(intent)
}

// Extract implements Mapper.
func (m FromEntity) Extract(ctx *Context) interface{} {
	return ctx.GetEntityValue(m.Entity, m.Role, m.Group)
}

// FromTriggerIntent TODO
type FromTriggerIntent struct {
	Value  interface{}
	Intent IntentFilter
}

// Desires implements Mapper.
func (m FromTriggerIntent) Desires(intent string) bool {
	return m.Intent == nil || m.Intent.Desires(intent)
}

// Extract implements Mapper.
func (m FromTriggerIntent) Extract(ctx *Context) interface{} {
	// return nothing - Extract is used for requested slots, trigger_intent is
	// used on form activation only.
	return nil
}

// FromIntent TODO
type FromIntent struct {
	Value  interface{}
	Intent IntentFilter
}

// Desires implements Mapper.
func (m FromIntent) Desires(intent string) bool {
	return m.Intent == nil || m.Intent.Desires(intent)
}

// Extract implements Mapper.
func (m FromIntent) Extract(ctx *Context) interface{} {
	return m.Value
}

// FromText TODO
type FromText struct {
	Intent IntentFilter
}

// Desires implements Mapper.
func (m FromText) Desires(intent string) bool {
	return m.Intent == nil || m.Intent.Desires(intent)
}

// Extract implements Mapper.
func (m FromText) Extract(ctx *Context) interface{} {
	return ctx.Tracker.LatestMessage.Text
}
