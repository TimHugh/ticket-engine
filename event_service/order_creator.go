package event_service

import (
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/lib/repos"
)

type OrderCreator struct {
	OrderRepository repos.OrderRepository
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder := o.OrderRepository.Find(orderID)
	if existingOrder != nil {
		return errors.New(fmt.Sprintf("Couldn't create duplicate order %s.", orderID))
	}

	o.OrderRepository.Store(repos.Order{
		ID:         orderID,
		LocationID: LocationID,
	})
	return nil
}
