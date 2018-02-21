package event_service

import (
	"errors"
	"fmt"
)

type OrderCreator struct {
	OrderRepository OrderRepository
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder := o.OrderRepository.Find(orderID)
	if existingOrder != nil {
		return errors.New(fmt.Sprintf("Couldn't create duplicate order %s.", orderID))
	}

	o.OrderRepository.Store(Order{
		ID:         orderID,
		LocationID: LocationID,
	})
	return nil
}
