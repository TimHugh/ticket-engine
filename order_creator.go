package main

import (
	"fmt"

	"github.com/timhugh/ticket_service/common"
)

type OrderCreator struct {
	orderRepository orderRepository
}

type orderRepository interface {
	Find(string) (*common.Order, error)
	Store(common.Order) error
}

func (o OrderCreator) Create(orderID string, LocationID string) error {
	existingOrder, _ := o.orderRepository.Find(orderID)
	if existingOrder != nil {
		return fmt.Errorf("Couldn't create duplicate order %s.", orderID)
	}

	return o.orderRepository.Store(common.Order{
		ID:         orderID,
		LocationID: LocationID,
	})
}
