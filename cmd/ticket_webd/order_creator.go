package main

import (
	"fmt"

	"github.com/timhugh/ticket_service/root"
)

type OrderCreator struct {
	OrderRepository
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder, _ := o.OrderRepository.Find(LocationID, orderID)
	if existingOrder != nil {
		return fmt.Errorf("Couldn't create duplicate order %s.", orderID)
	}

	return o.OrderRepository.Create(root.Order{
		ID:         orderID,
		LocationID: LocationID,
	})
}
