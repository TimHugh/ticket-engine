package event_service

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

func TestRegistersRoutes(t *testing.T) {
	router := NewEventRouter()
	handler := &MockHandler{Error: false}
	event := EventType("event")

	router.Register(event, handler)

	if router.Routes[event] != handler {
		t.Error("Expected handler to be registered")
	}
}

func TestRejectsUnknownRoutes(t *testing.T) {
	router := &EventRouter{
		Routes: RouteList{
			EventType("known event"): &MockHandler{Error: false},
		},
	}

	mockEvent := Event{
		EventType: EventType("unknown event"),
	}

	err := router.Dispatch(mockEvent)
	if err == nil {
		t.Error("Expected error on dispatching unknown event")
	}
}

func DelegatesHandlerErrors(t *testing.T) {
	event := Event{
		EventType: EventType("event"),
	}

	router := &EventRouter{
		Routes: RouteList{
			event.EventType: &MockHandler{Error: true},
		},
	}

	err := router.Dispatch(event)
	if err == nil {
		t.Error("Expected router to delegate error back from handler")
	}
}

func TestDispatchesEvents(t *testing.T) {
	event := Event{
		EventType: EventType("event"),
	}

	handler := MockHandler{Error: false}
	router := &EventRouter{
		Routes: RouteList{
			event.EventType: &handler,
		},
	}

	router.Dispatch(event)
	if handler.Handled != true {
		t.Error("Expected event to be dispatched to handler.")
	}
}
