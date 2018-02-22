package main

import (
	"testing"

	"errors"
	"strings"
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

func TestRouting(t *testing.T) {
	router := NewEventRouter()
	goodHandler := &MockHandler{}
	errorHandler := &MockHandler{Error: true}
	router.Register("event", goodHandler)
	router.Register("error_event", errorHandler)

	goodEvent := Event{Type: "event"}
	err := router.Dispatch(goodEvent)
	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}
	if !goodHandler.Handled {
		t.Error("Expected handler to be registered and called.")
	}

	unknownEvent := Event{Type: "unknown_event"}
	err = router.Dispatch(unknownEvent)
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Errorf("Expected unknown event error")
	}

	errorEvent := Event{Type: "error_event"}
	err = router.Dispatch(errorEvent)
	if err == nil || !strings.Contains(err.Error(), "error") {
		t.Error("Expected to receive error from handler")
	}
}
