package event_service

import (
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/lib/repos"
)

type PaymentUpdateHandler struct {
	OrderCreator       OrderCreator
	LocationRepository repos.LocationRepository
}

func (p PaymentUpdateHandler) Handle(event Event) error {

	// TODO: this is implemented in the SquareRequestValidator now
	location := p.LocationRepository.Find(event.LocationID)
	if location == nil {
		return errors.New(fmt.Sprintf("Received event for unknown location ID %s", event.LocationID))
	}

	return p.OrderCreator.Create(event.OrderID, event.LocationID)
}
