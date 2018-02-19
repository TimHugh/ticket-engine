package event_service

import (
	"errors"
	"fmt"
)

type Order struct {
	ID         string
	LocationID string
}

type OrderRepository interface {
	Store(Order)
	Find(string) *Order
}

type OrderCreator struct {
	OrderRepository OrderRepository
}

func (o OrderCreator) Create(order Order) error {
	existingOrder := o.OrderRepository.Find(order.ID)
	if existingOrder != nil {
		return errors.New(fmt.Sprintf("Couldn't create duplicate order %s.", order.ID))
	}
	o.OrderRepository.Store(order)
	return nil
}
