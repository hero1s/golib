// Package event is lightweight event manager and dispatcher implements by Go.
//
package event

import "github.com/hero1s/golib/stringutils"

// Event interface
type Event interface {
	Name() string
	// Target() interface{}
	Get(key string) interface{}
	Add(key string, val interface{})
	Set(key string, val interface{})
	Data() map[string]interface{}
	SetData(M) Event
	Abort(bool)
	IsAborted() bool
}

// BasicEvent a basic event struct define.
type BasicEvent struct {
	// event name
	name string
	// user data.
	data map[string]interface{}
	// target
	target interface{}
	// mark is aborted
	aborted bool
}

// NewBasic new an basic event instance
func NewBasic(name string, data M) *BasicEvent {
	if data == nil {
		data = make(map[string]interface{})
	}

	return &BasicEvent{
		name: name,
		data: data,
	}
}

// Abort abort event loop exec
func (e *BasicEvent) Abort(abort bool) {
	e.aborted = abort
}

// Fill event data
func (e *BasicEvent) Fill(target interface{}, data M) *BasicEvent {
	if data != nil {
		e.data = data
	}

	e.target = target
	return e
}

// AttachTo add current event to the event manager.
func (e *BasicEvent) AttachTo(em ManagerFace) {
	em.AddEvent(e)
}

// Get get data by index
func (e *BasicEvent) Get(key string) interface{} {
	if v, ok := e.data[key]; ok {
		return v
	}

	return nil
}

// 便捷转换
func (e *BasicEvent) GetInt(key string) int {
	if v, ok := e.data[key]; ok {
		return stringutils.GetInt(v)
	}
	return 0
}

func (e *BasicEvent) GetString(key string) string {
	if v, ok := e.data[key]; ok {
		return stringutils.GetString(v)
	}
	return ""
}

// Add value by key
func (e *BasicEvent) Add(key string, val interface{}) {
	if _, ok := e.data[key]; !ok {
		e.Set(key, val)
	}
}

// Set value by key
func (e *BasicEvent) Set(key string, val interface{}) {
	if e.data == nil {
		e.data = make(map[string]interface{})
	}

	e.data[key] = val
}

// Name get event name
func (e *BasicEvent) Name() string {
	return e.name
}

// Data get all data
func (e *BasicEvent) Data() map[string]interface{} {
	return e.data
}

// IsAborted check.
func (e *BasicEvent) IsAborted() bool {
	return e.aborted
}

// Target get target
func (e *BasicEvent) Target() interface{} {
	return e.target
}

// SetName set event name
func (e *BasicEvent) SetName(name string) *BasicEvent {
	e.name = name
	return e
}

// SetData set data to the event
func (e *BasicEvent) SetData(data M) Event {
	if data != nil {
		e.data = data
	}
	return e
}

// SetTarget set event target
func (e *BasicEvent) SetTarget(target interface{}) *BasicEvent {
	e.target = target
	return e
}
