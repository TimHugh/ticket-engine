package main

import (
	"testing"

	"errors"
)

type MockHandler struct {
	Error   bool
	Handled bool
}

func (m *MockHandler) Handle(event Event) error {
	m.Handled = true
	if m.Error {
		return errors.New("error")
	}
	return nil
}

func TestRegisterRoutes(t *testing.T) {
	router := NewEventRouter()
	handler := &MockHandler{Error: false}
	event := Event{
		Type: "event",
	}

	router.Register(event.Type, handler)
	router.Dispatch(event)

	if !handler.Handled {
		t.Error("Expected handler to be registered and called.")
	}
}

func TestRejectsUnknownRoutes(t *testing.T) {
	router := NewEventRouter()
	event := Event{
		Type: "event",
	}

	err := router.Dispatch(event)

	if err == nil {
		t.Error("Expected error on dispatching unknown event")
	}
}

func DelegatesHandlerErrors(t *testing.T) {
	router := NewEventRouter()
	handler := &MockHandler{Error: true}
	event := Event{
		Type: "event",
	}

	router.Register(event.Type, handler)
	err := router.Dispatch(event)

	if err == nil {
		t.Error("Expected router to delegate error back from handler")
	}
}
