package main

import (
	"fmt"
)

type Event struct {
	OrderID    string `json:"entity_id"`
	Type       string `json:"event_type"`
	LocationID string `json:"location_id"`
}

type EventHandler interface {
	Handle(event Event) error
}

type NoopHandler struct{}

func (h NoopHandler) Handle(event Event) error {
	return nil
}

type RouteList map[string]EventHandler

type EventRouter struct {
	routes RouteList
}

func NewEventRouter() EventRouter {
	return EventRouter{
		routes: make(RouteList),
	}
}

func (e EventRouter) Register(event string, handler EventHandler) {
	e.routes[event] = handler
}

func (e EventRouter) Dispatch(event Event) error {
	handler, exists := e.routes[event.Type]
	if !exists {
		return fmt.Errorf("Recieved unknown event type '%s'", event.Type)
	}
	return handler.Handle(event)
}
