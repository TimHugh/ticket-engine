package event_service

import (
	"errors"
	"fmt"
)

type Event struct {
	OrderID    string
	EventType  EventType
	LocationID string
}

type EventHandler interface {
	Handle(event Event) error
}

type EventType string

type RouteList map[EventType]EventHandler

type EventRouter struct {
	Routes RouteList
}

func NewEventRouter() *EventRouter {
	router := EventRouter{
		Routes: make(RouteList),
	}
	return &router
}

func (e EventRouter) Register(event EventType, handler EventHandler) {
	e.Routes[event] = handler
}

func (e EventRouter) Dispatch(event Event) error {
	handler, exists := e.Routes[event.EventType]
	if !exists {
		return errors.New(fmt.Sprintf("No route registered for event %s", event.EventType))
	}
	return handler.Handle(event)
}
