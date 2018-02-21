package event_service

import (
	"errors"
	"fmt"
)

type Event struct {
	OrderID    string
	Type       string
	LocationID string
}

type EventHandler interface {
	Handle(event Event) error
}

type RouteList map[string]EventHandler

type EventRouter struct {
	routes RouteList
}

func NewEventRouter() *EventRouter {
	router := EventRouter{
		routes: make(RouteList),
	}
	return &router
}

func (e EventRouter) Register(event string, handler EventHandler) {
	e.routes[event] = handler
}

func (e EventRouter) Dispatch(event Event) error {
	handler, exists := e.routes[event.Type]
	if !exists {
		return errors.New(fmt.Sprintf("No route registered for event %s", event.Type))
	}
	return handler.Handle(event)
}
