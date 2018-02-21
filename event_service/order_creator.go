package event_service

import (
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

type OrderCreator struct {
	OrderRepository common.OrderRepository
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder := o.OrderRepository.Find(orderID)
	if existingOrder != nil {
		return errors.New(fmt.Sprintf("Couldn't create duplicate order %s.", orderID))
	}

	o.OrderRepository.Store(common.Order{
		ID:         orderID,
		LocationID: LocationID,
	})
	return nil
}
