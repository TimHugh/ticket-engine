package main

import (
	"errors"
	"fmt"

	"github.com/timhugh/ticket-engine/common"
)

type OrderCreator struct {
	orderRepository common.OrderRepository
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder := o.orderRepository.Find(orderID)
	if existingOrder != nil {
		return errors.New(fmt.Sprintf("Couldn't create duplicate order %s.", orderID))
	}

	o.orderRepository.Store(common.Order{
		ID:         orderID,
		LocationID: LocationID,
	})

	return nil
}
