package eventemitter

import (
	"reflect"
)

type Response struct {
	// Name of the Event
	EventName string

	// Slice of all the handler's return values
	Ret []interface{}
}

type EventEmitter struct {
	events map[string][]Listener
}

type Listener struct {
	Callback reflect.Value
	ID       interface{}
}

func New() *EventEmitter {
	e := new(EventEmitter)
	e.Init()

	return e
}

// Allocates the EventEmitters memory. Has to be called when
// embedding an EventEmitter in another Type.
func (self *EventEmitter) Init() {
	self.events = make(map[string][]Listener)
}

func (self *EventEmitter) Listeners(event string) []Listener {
	return self.events[event]
}

// Alias to AddListener.
func (self *EventEmitter) On(event string, listener interface{}, id interface{}) {
	self.AddListener(event, listener, id)
}

// AddListener adds an event listener on the given event name. id is used to
// keep track of different listeners and distinguish one from the other. If a
// listener with a duplicate id is added, the second one is ignored.
func (self *EventEmitter) AddListener(event string, listener interface{}, id interface{}) {
	// Check if the event exists, otherwise initialize the list
	// of handlers for this event.
	if _, exists := self.events[event]; !exists {
		self.events[event] = []Listener{}
	}

	// Check if the listener id exists.
	for _, x := range self.events[event] {
		if x.ID == id { // match found
			return
		}
	}

	var l reflect.Value
	l, ok := listener.(reflect.Value)
	if !ok {
		l = reflect.ValueOf(listener)
	}
	x := Listener{
		Callback: l,
		ID:       id,
	}
	self.events[event] = append(self.events[event], x)
}

// RemoveListener removes the listener listening to `event' with the specified
// id.
func (self *EventEmitter) RemoveListener(event string, id interface{}) {
	listeners := self.events[event]
	for i, x := range listeners {
		if x.ID == id { // we got a match, so remove
			self.events[event] = append(listeners[:i], listeners[i+1:]...)
		}
		return
	}
}

// RemoveListeners removes all listeners from the given event.
func (self *EventEmitter) RemoveListeners(event string) {
	delete(self.events, event)
}

// Emits the given event. Puts all arguments following the event name
// into the Event's `Argv` member. Returns a channel if listeners were
// called, nil otherwise.
func (self *EventEmitter) Emit(event string, argv ...interface{}) <-chan *Response {
	listeners, exists := self.events[event]

	if !exists {
		return nil
	}

	var callArgv []reflect.Value
	c := make(chan *Response)

	for _, a := range argv {
		callArgv = append(callArgv, reflect.ValueOf(a))
	}

	for _, listener := range listeners {
		go func(listener reflect.Value) {
			retVals := listener.Call(callArgv)
			var ret []interface{}

			for _, r := range retVals {
				ret = append(ret, r.Interface())
			}

			c <- &Response{EventName: event, Ret: ret}
		}(listener.Callback)
	}

	return c
}
