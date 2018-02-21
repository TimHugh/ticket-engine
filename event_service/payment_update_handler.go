package event_service

import (
	"errors"
	"fmt"
)

type PaymentUpdateHandler struct {
	OrderCreator       OrderCreator
	LocationRepository LocationRepository
}

func (p PaymentUpdateHandler) Handle(event Event) error {
	location := p.LocationRepository.Find(event.LocationID)
	if location == nil {
		return errors.New(fmt.Sprintf("Received event for unknown location ID %s", event.LocationID))
	}
	return p.OrderCreator.Create(event.OrderID, event.LocationID)
}